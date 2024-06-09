package mfrc522

import (
	"errors"
	"machine"
	"time"
)

// MFRC522 holds the relevant configuration for the MFRC522 RFID reader.
type MFRC522 struct {
	// SPI is used to communicate with the MFRC522 reader.
	// This is technically the host's SPI interface.
	spi machine.SPI

	// rstPin is the reset pin for the MFRC522 reader.
	// It is used to initialize the reader.
	rstPin machine.Pin

	// irqPin is the interrupt pin for the MFRC522 reader.
	// It notifies the host when a card is present.
	irqPin machine.Pin

	// irqTimeout is the maximum time to wait for an interrupt from the reader.
	// The interrupt signals that a card is present.
	irqTimeout time.Duration
}

// Init initializes the MFRC522 reader.
func Init(rstPin, irqPin machine.Pin, irqTimeout time.Duration) (*MFRC522, error) {
	mfrc522 := &MFRC522{
		spi:        machine.SPI0,
		rstPin:     rstPin,
		irqPin:     irqPin,
		irqTimeout: irqTimeout,
	}

	if err := mfrc522.spi.Configure(machine.SPIConfig{Frequency: 1000000}); err != nil {
		return nil, errors.New("failed to configure SPI" + err.Error())
	}

	mfrc522.rstPin.Configure(machine.PinConfig{Mode: machine.PinInput})
	if !mfrc522.rstPin.Get() {
		mfrc522.rstPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
		mfrc522.rstPin.Low()
		time.Sleep(2 * time.Microsecond)
		mfrc522.rstPin.High()
		time.Sleep(50 * time.Microsecond)
	}

	if err := mfrc522.WriteSequence(InitSequence); err != nil {
		return nil, errors.New("Failed to write initialization sequence:" + err.Error())
	}

	if err := mfrc522.AntennaOn(); err != nil {
		return nil, errors.New("Failed to turn on antenna:" + err.Error())
	}

	return mfrc522, nil
}

// Version returns the firmware version of the MFRC522 reader.
func (m *MFRC522) Version() (byte, error) {
	ver, err := m.ReadRegister(VersionReg)
	if err != nil {
		return 0, err
	}

	return ver, nil
}

// ReadTagUUID returns the UUID of the selected RFID tag.
func (m *MFRC522) ReadTagUUID() ([]byte, error) {
	var err error
	defer func() {
		if err == nil {
			err = m.StopCrypto()
		}
	}()

	return m.selectCard()
}

// Reset sends a soft reset command to the MFRC522 reader.
func (m *MFRC522) Reset() error {
	if err := m.WriteRegister(CommandReg, SoftResetCmd); err != nil {
		return err
	}

	time.Sleep(50 * time.Microsecond)
	for range 3 {
		val, err := m.ReadRegister(CommandReg)
		if err != nil {
			return err
		}

		if val&(1<<4) == 0 {
			break
		}

		time.Sleep(50 * time.Microsecond)
	}

	return nil
}

// Exit handles the cleanup of the MFRC522 reader.
func (m *MFRC522) Exit() {
	if err := m.AntennaOff(); err != nil {
		println("Failed to turn off antenna:", err)
	}
}

// WriteRegister writes a byte to the specified register.
func (m *MFRC522) WriteRegister(reg Register, val byte) error {
	return m.WriteRegisterBytes(reg, []byte{val})
}

// ReadRegister reads a byte from the specified register.
func (m *MFRC522) ReadRegister(reg Register) (byte, error) {
	val, err := m.ReadRegisterBytes(reg, 1)
	if err != nil {
		return 0, err
	}

	return val[0], nil
}

// AntennaOn turns on the antenna of the MFRC522 reader.
func (m *MFRC522) AntennaOn() error {
	state, err := m.ReadRegister(TxControlReg)
	if err != nil {
		return err
	}

	if (state & 0x03) != 0x03 {
		return m.WriteRegister(TxControlReg, state|0x03)
	}

	return nil
}

// AntennaOff turns off the antenna of the MFRC522 reader.
func (m *MFRC522) AntennaOff() error {
	return m.ClearBitmask(TxControlReg, 0x03)
}

// SetBitmask sets the specified bits in the register.
func (m *MFRC522) SetBitmask(reg Register, mask byte) error {
	val, err := m.ReadRegister(reg)
	if err != nil {
		return err
	}

	return m.WriteRegister(reg, val|mask)
}

// ClearBitmask clears the specified bits in the register.
func (m *MFRC522) ClearBitmask(reg Register, mask byte) error {
	val, err := m.ReadRegister(reg)
	if err != nil {
		return err
	}

	return m.WriteRegister(reg, val&^mask)
}

// WaitForInterrupt waits for an interrupt from the MFRC522 reader, which
// signals that a tag is present.
func (m *MFRC522) WaitForInterrupt(timeout time.Duration) error {
	irqChan := make(chan bool, 1)
	defer close(irqChan)

	irqFunc := func(p machine.Pin) { irqChan <- true }
	if err := m.irqPin.SetInterrupt(machine.PinToggle, irqFunc); err != nil {
		return err
	}

	if err := m.WriteSequence([]WriteCommand{
		{ComIrqReg, 0x00},
		{ComIEnReg, 0xA0},
	}); err != nil {
		return err
	}

	start := time.Now()
	for start.Before(start.Add(timeout)) {
		if err := m.WriteSequence([]WriteCommand{
			{FIFODataReg, 0x26},
			{CommandReg, TransceiveCmd},
			{BitFramingReg, 0x87},
		}); err != nil {
			return err
		}

		select {
		case <-irqChan:
			return nil
		case <-time.After(100 * time.Millisecond):
		}
	}

	return nil
}

// ClearIRQ clears the interrupt request bits.
func (m *MFRC522) ClearIRQ() error {
	return m.WaitForInterrupt(0)
}

// StopCrypto stops the crypto1 unit, which is needed after entering an authenticated state.
func (m *MFRC522) StopCrypto() error {
	return m.ClearBitmask(Status2Reg, 0x08)
}

// ReadAuthentication reads the tag's authentication data from the specified sector.
func (m *MFRC522) ReadAuthentication(authMode, sector byte, key []byte) ([]byte, error) {
	var err error
	defer func() {
		if err == nil {
			err = m.StopCrypto()
		}
	}()

	uuid, err := m.selectCard()
	if err != nil {
		return nil, err
	}

	addr := sector*4 + 3

	auth, err := m.authenticate(authMode, addr, key, uuid)
	if err != nil {
		return nil, err
	}
	if auth != AuthOk {
		return nil, errors.New("authentication failed")
	}

	return m.readTag(addr)
}

// ReadTagBlock reads a block of data from the specified address (sector+block).
func (m *MFRC522) ReadTagBlock(authMode, sector, block byte, key []byte) ([]byte, error) {
	var err error
	defer func() {
		if err == nil {
			err = m.StopCrypto()
		}
	}()

	uuid, err := m.selectCard()
	if err != nil {
		return nil, err
	}

	auth, err := m.authenticate(authMode, sector*4+block, key, uuid)
	if err != nil {
		return nil, err
	}
	if auth != AuthOk {
		return nil, errors.New("authentication failed")
	}

	return m.readTag(sector*4 + (block % 3))
}

// WriteTag writes data to the specified address (sector+block).
func (m *MFRC522) WriteTag(authMode, sector, block byte, data, key []byte) error {
	var err error
	defer func() {
		if err == nil {
			err = m.StopCrypto()
		}
	}()

	uuid, err := m.selectCard()
	if err != nil {
		return err
	}

	auth, err := m.authenticate(authMode, sector*4+3, key, uuid)
	if err != nil {
		return err
	}
	if auth != AuthOk {
		return errors.New("authentication failed")
	}

	return m.writeTag(sector*4+(block%3), data)
}

/*
	Not-implemented section

	These functions are not needed for the basic functionality of the RFID reader,
	so they were not implemented for the lab exercise, but might be in the future.
*/

// SetAntennaGain sets the antenna gain of the MFRC522 reader.
func (m *MFRC522) SetAntennaGain() error {
	return errors.New("not implemented")
}

// AntennaGain returns the current antenna gain of the MFRC522 reader.
func (m *MFRC522) AntennaGain() error {
	return errors.New("not implemented")
}

// SelfTest performs a self-test on the MFRC522 reader
// (specified in Chapter 16.1.1 of the MFRC55 datasheet).
func (m *MFRC522) SelfTest() (bool, error) {
	return false, errors.New("not implemented")
}

// PowerDown puts the MFRC522 reader into power-down mode.
func (m *MFRC522) PowerDown() error {
	return errors.New("not implemented")
}

// PowerUp wakes the MFRC522 reader from power-down mode.
func (m *MFRC522) PowerUp() error {
	return errors.New("not implemented")
}
