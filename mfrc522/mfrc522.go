package mfrc522

import (
	"machine"
)

// MFRC522 holds the relevant configuration for the MFRC522 RFID reader.
type MFRC522 struct {
	// SPI is used to communicate with the MFRC522 reader.
	// This is technically the host's SPI interface.
	spi machine.SPI

	// MFRC522 also supports UART and I2C, but the library currently only supports SPI.
}

func (m *MFRC522) WriteRegister(reg Register, val byte) {
	_, err := m.spi.Transfer(byte(reg))
	if err != nil {
		println(err)
	}

	_, err = m.spi.Transfer(val)
	if err != nil {
		println(err)
	}
}

func (m *MFRC522) ReadRegister(reg Register) byte {
	return 0
}

func (m *MFRC522) AntennaOn() {}

func (m *MFRC522) AntennaOff() {}

func (m *MFRC522) Reset() {
	m.WriteRegister(CommandReg, 0x0F)
}

func (m *MFRC522) Version() byte {
	return m.ReadRegister(VersionReg)
}

func Init() *MFRC522 {
	mfrc522 := &MFRC522{spi: machine.SPI0}

	err := mfrc522.spi.Configure(machine.SPIConfig{Frequency: 9600})
	if err != nil {
		println(err)
	}

	mfrc522.Reset()
	mfrc522.WriteRegister(TModeReg, 0x8D)
	mfrc522.WriteRegister(TPrescalerReg, 0x3E)
	mfrc522.WriteRegister(TReloadHighReg, 30)
	mfrc522.WriteRegister(TReloadLowReg, 0)
	mfrc522.WriteRegister(TxASKReg, 0x40)
	mfrc522.WriteRegister(ModeReg, 0x3D)
	mfrc522.AntennaOn()

	return mfrc522
}

func (m *MFRC522) Exit() {
	m.AntennaOff()
}
