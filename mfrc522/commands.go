package mfrc522

// Command is an MFRC55 command.
// Command definitions are available in Chapter 10 of the MFRC55 datasheet.
type Command byte

var (
	// IdleCmd cancels the current command's execution.
	IdleCmd Command = 0x00

	// MemCmd stores 25 bytes into the internal buffer.
	MemCmd Command = 0x01

	// GenerateRandomIDCmd generates a 10-byte random ID number.
	GenerateRandomIDCmd Command = 0x02

	// CalcCRCCmd activates the CRC coprocessor or performs a self-test.
	CalcCRCCmd Command = 0x03

	// TransmitCmd sends data from the FIFO buffer to the transmitter.
	TransmitCmd Command = 0x04

	// NoCmdChangeCmd can be used to modify the CommandReg register bits without affecting the command execution.
	NoCmdChangeCmd Command = 0x07

	// ReceiveCmd activates the receiver circuits.
	ReceiveCmd Command = 0x08

	// TransceiveCmd transmits data from the FIFO buffer to the transmitter and automatically activates the receiver after transmission.
	TransceiveCmd Command = 0x0C

	// MFAuthentCmd performs the MIFARE standard authentication as a reader.
	MFAuthentCmd Command = 0x0E

	// SoftResetCmd resets the MFRC55.
	SoftResetCmd Command = 0x0F
)
