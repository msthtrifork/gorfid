package mfrc522

import (
	"errors"
	"time"
)

// ReadRegisterBytes allows reading multiple bytes from a register.
func (m *MFRC522) ReadRegisterBytes(reg Register, readLen int) ([]byte, error) {
	if readLen < 1 {
		return nil, nil
	}

	data := make([]byte, 0, readLen+1)
	for range readLen {
		data = append(data, 0x80|reg)
	}
	data = append(data, 0)

	res := make([]byte, len(data))
	if err := m.spi.Tx(data, res); err != nil {
		return nil, err
	}

	return res[1:], nil
}

// WriteRegisterBytes allows writing multiple bytes to a register.
func (m *MFRC522) WriteRegisterBytes(reg Register, val []byte) error {
	data := append([]byte{reg}, val...)

	if err := m.spi.Tx(data, nil); err != nil {
		return err
	}

	return nil
}

// WriteSequence is a convenience function for writing predefined
// sequences of commands to registers.
func (m *MFRC522) WriteSequence(commands []WriteCommand) error {
	for _, cmd := range commands {
		if err := m.WriteRegister(cmd.Register, cmd.RegisterCommand); err != nil {
			return errors.New("failed to write command" + string(cmd.RegisterCommand) + "to register" +
				string(cmd.Register) + ":" + err.Error())
		}
	}

	return nil
}

// selectCard sets the detected card as selected in the reader and returns its UUID.
func (m *MFRC522) selectCard() ([]byte, error) {
	defer func() { _ = m.ClearIRQ() }()

	if err := m.WaitForInterrupt(m.irqTimeout); err != nil {
		return nil, err
	}

	// Needed to set bit-framing register and clear FIFO buffer
	if _, err := m.numTagBlocks(); err != nil {
		return nil, err
	}

	uuid, err := m.antiCollision()
	if err != nil {
		return nil, err
	}

	_, err = m.selectUUID(uuid)
	if err != nil {
		return nil, err
	}

	if uuid[0] == 0x88 {
		// Some tags have longer UIDs, so the remaining bytes need to be read separately,
		// but this is not supported for this lab exercise.
		return nil, errors.New("tag not supported")
	}

	return uuid[:len(uuid)-1], nil
}

// authenticate authenticates an address (sector+block) for the selected tag.
func (m *MFRC522) authenticate(authMode, addr byte, key, uuid []byte) (AuthStatus, error) {
	data := append([]byte{authMode, addr}, key...)
	data = append(data, uuid[:4]...)

	_, err := m.writeTagCommand(MFAuthentCmd, data)
	if err != nil {
		return AuthReadFail, err
	}

	val, err := m.ReadRegister(Status2Reg)
	if err != nil || val&0x08 == 0 {
		return AuthFail, err
	}

	return AuthOk, nil
}

// crc calculates the CRC of the given data (on the reader).
func (m *MFRC522) crc(data []byte) ([]byte, error) {
	if err := m.WriteSequence([]WriteCommand{
		{CommandReg, IdleCmd},
		{DivIEnReg, 0x04},
		{FIFOLevelReg, 0x80},
	}); err != nil {
		return nil, err
	}
	if err := m.WriteRegisterBytes(FIFODataReg, data); err != nil {
		return nil, err
	}
	if err := m.WriteRegister(CommandReg, CalcCRCCmd); err != nil {
		return nil, err
	}

	for range 100 {
		val, err := m.ReadRegister(DivIrqReg)
		if err != nil {
			return nil, err
		}

		if val&0x04 != 0 {
			if err = m.WriteRegister(CommandReg, IdleCmd); err != nil {
				return nil, err
			}

			crc := make([]byte, 2)
			crc[0], err = m.ReadRegister(CRCResultLowReg)
			if err != nil {
				return nil, err
			}

			crc[1], err = m.ReadRegister(CRCResultHighReg)
			if err != nil {
				return nil, err
			}

			return crc, nil
		}

		time.Sleep(1 * time.Millisecond)
	}

	return nil, errors.New("timed out while calculating CRC")
}

// verifyCRC calculates the CRC and sends it to the tag for verification.
func (m *MFRC522) verifyCRC(cmd, addr byte) ([]byte, error) {
	crc, err := m.crc([]byte{cmd, addr})
	if err != nil {
		return nil, err
	}

	return m.writeTagCommand(TransceiveCmd, []byte{cmd, addr, crc[0], crc[1]})
}

// readTag reads the address (sector+block) from the selected tag.
func (m *MFRC522) readTag(addr byte) ([]byte, error) {
	data, err := m.verifyCRC(ReadBlockCmd, addr)
	if err != nil {
		return nil, err
	}
	if len(data) != 16 {
		return nil, errors.New("invalid data length, expected 16 bytes")
	}

	return data, nil
}

// writeTag writes data to the address (sector+block) on the selected tag.
func (m *MFRC522) writeTag(addr byte, data []byte) error {
	data, err := m.verifyCRC(WriteBlockCmd, addr)
	if err != nil {
		return err
	}
	if data[0]&0x0F != 0x0A {
		return errors.New("couldn't authorize write operation")
	}

	crc, err := m.crc(data[:16])
	if err != nil {
		return err
	}

	writeData := append(data[:16], crc...)
	data, err = m.writeTagCommand(WriteBlockCmd, writeData)
	if err != nil {
		return err
	}

	if data[0]&0x0F != 0x0A {
		return errors.New("write operation failed")
	}

	return nil
}

// writeTagCommand writes a command to the tag and returns the response.
func (m *MFRC522) writeTagCommand(cmd RegisterCommand, data []byte) ([]byte, error) {
	var irqEn, irqWait byte
	switch cmd {
	case MFAuthentCmd:
		irqEn = 0x12
		irqWait = 0x10
	case TransceiveCmd:
		irqEn = 0x77
		irqWait = 0x30
	}

	if err := m.WriteRegister(ComIEnReg, irqEn|0x80); err != nil {
		return nil, err
	}
	if err := m.ClearBitmask(ComIrqReg, 0x80); err != nil {
		return nil, err
	}
	if err := m.SetBitmask(FIFOLevelReg, 0x80); err != nil {
		return nil, err
	}
	if err := m.WriteRegister(CommandReg, IdleCmd); err != nil {
		return nil, err
	}
	if err := m.WriteRegisterBytes(FIFODataReg, data); err != nil {
		return nil, err
	}
	if err := m.WriteRegister(CommandReg, cmd); err != nil {
		return nil, err
	}

	if cmd == TransceiveCmd {
		if err := m.SetBitmask(BitFramingReg, 0x80); err != nil {
			return nil, err
		}
	}

	// Wait for data to be sent
	for range 2000 {
		val, err := m.ReadRegister(ComIrqReg)
		if err != nil {
			return nil, err
		}

		if val&(irqWait|0x01) != 0x00 {
			break
		}
	}

	if err := m.ClearBitmask(BitFramingReg, 0x80); err != nil {
		return nil, err
	}

	errStatus, err := m.ReadRegister(ErrorReg)
	if err != nil || errStatus&0x1B != 0 {
		return nil, errors.New("error during command execution")
	}

	if errStatus&irqEn&0x01 == 1 {
		return nil, errors.New("invalid interrupt request")
	}

	if cmd == TransceiveCmd {
		level, err := m.ReadRegister(FIFOLevelReg)
		if err != nil {
			return nil, err
		}

		reg, err := m.ReadRegister(ControlReg)
		if err != nil {
			return nil, err
		}

		var dataLen int
		if reg&0x07 != 0x00 {
			dataLen = (int(level)-1)*8 + int(reg&0x07)
		} else {
			dataLen = int(level) * 8
		}

		if level == 0 {
			level = 1
		} else if level > 16 {
			level = 16
		}

		res, err := m.ReadRegisterBytes(FIFODataReg, dataLen)
		if err != nil {
			return nil, err
		}

		return res, nil
	}

	return nil, nil
}

// numTagBlocks returns the number of blocks in the tag.
func (m *MFRC522) numTagBlocks() (int, error) {
	if err := m.WriteRegister(BitFramingReg, 0x07); err != nil {
		return -1, err
	}

	data, err := m.writeTagCommand(TransceiveCmd, []byte{0x26})
	if err != nil {
		return -1, err
	}

	return len(data), nil
}

// antiCollision performs the anti-collision procedure and returns the UUID of the selected tag.
func (m *MFRC522) antiCollision() ([]byte, error) {
	if err := m.WriteRegister(BitFramingReg, 0x00); err != nil {
		return nil, err
	}

	data, err := m.writeTagCommand(TransceiveCmd, []byte{AntiCollSelect1Cmd, 0x20})
	if err != nil {
		return nil, err
	}

	if len(data) != 5 {
		return nil, errors.New("invalid data length, expected 5 bytes")
	}

	var crc byte
	for i := 0; i < 4; i++ {
		crc = crc ^ data[i]
	}
	if crc != data[4] {
		return nil, errors.New("CRC mismatch")
	}

	return data, nil
}

// selectUUID selects the tag with the given UUID.
func (m *MFRC522) selectUUID(uuid []byte) (byte, error) {
	data := append([]byte{AntiCollSelect1Cmd, 0x70}, uuid...)

	crc, err := m.crc(data)
	if err != nil {
		return 0, err
	}

	data = append(data, crc...)
	res, err := m.writeTagCommand(TransceiveCmd, data)
	if err != nil {
		return 0, err
	}

	var selectData byte
	if len(res) != 0x18 {
		selectData = res[0]
	} else {
		selectData = 0
	}

	return selectData, nil
}

/*
	Not-implemented section

	These functions are not needed for the basic functionality of the RFID reader,
	so they were not implemented for the lab exercise, but might be in the future.
*/

// antiCollision2 performs the anti-collision procedure for tags with longer UIDs.
func (m *MFRC522) antiCollision2() ([]byte, error) {
	return nil, errors.New("not implemented")
}
