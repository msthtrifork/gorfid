package mfrc522

// RegisterCommand is an MFRC55 register command.
// RegisterCommand definitions are available in Chapter 10 of the MFRC55 datasheet.
type RegisterCommand = byte

type WriteCommand struct {
	Register
	RegisterCommand
}

// Reader commands
const (
	// IdleCmd cancels the current command's execution.
	IdleCmd RegisterCommand = 0x00

	// MemCmd stores 25 bytes into the internal buffer.
	MemCmd RegisterCommand = 0x01

	// GenerateRandomIDCmd generates a 10-byte random ID number.
	GenerateRandomIDCmd RegisterCommand = 0x02

	// CalcCRCCmd activates the CRC coprocessor or performs a self-test.
	CalcCRCCmd RegisterCommand = 0x03

	// TransmitCmd sends data from the FIFO buffer to the transmitter.
	TransmitCmd RegisterCommand = 0x04

	// NoCmdChangeCmd can be used to modify the CommandReg register bits without affecting the command execution.
	NoCmdChangeCmd RegisterCommand = 0x07

	// ReceiveCmd activates the receiver circuits.
	ReceiveCmd RegisterCommand = 0x08

	// TransceiveCmd transmits data from the FIFO buffer to the transmitter and automatically activates the receiver after transmission.
	TransceiveCmd RegisterCommand = 0x0C

	// MFAuthentCmd performs the MIFARE standard authentication as a reader.
	MFAuthentCmd RegisterCommand = 0x0E

	// SoftResetCmd resets the MFRC55.
	SoftResetCmd RegisterCommand = 0x0F
)

// AuthStatus is the Tag authentication status.
type AuthStatus byte

// Authentication statuses
const (
	AuthOk       AuthStatus = 0x00
	AuthReadFail AuthStatus = 0x01
	AuthFail     AuthStatus = 0x02
)

// Sequences
var (
	// InitSequence is the initialization sequence for the MFRC55.
	InitSequence = []WriteCommand{
		{TxModeReg, 0x00},
		{RxModeReg, 0x00},
		{ModWidthReg, 0x26},
		{TModeReg, 0x80},
		{TPrescalerReg, 0xA9},
		{TReloadHighReg, 0x03},
		{TReloadLowReg, 0xE8},
		{TxASKReg, 0x40},
		{ModeReg, 0x3D},
	}
)
