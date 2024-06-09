package mfrc522

// TagCommand is an RFID tag command.
// They are not specified in the MFRC522 datasheet, but in various other documents.
type TagCommand = byte

const (
	// Commands for various tags (ISO 14443-3, Type A, section 6.4)

	// RequestACmd requests the tag to go from IDLE state to READY and prepare for anti-collision or selection.
	RequestACmd TagCommand = 0x26

	// HaltACmd instructs an active tag to go into the HALT state.
	HaltACmd TagCommand = 0x50

	// WakeUpACmd requests the tag to go from IDLE or HALT state to READY and prepare for anti-collision or selection.
	WakeUpACmd TagCommand = 0x52

	// CascadeTagCmd is used during anti-collision.
	CascadeTagCmd TagCommand = 0x88

	// AntiCollSelect1Cmd sets cascade level 1 during anti-collision.
	AntiCollSelect1Cmd TagCommand = 0x93

	// AntiCollSelect2Cmd sets cascade level 2 during anti-collision.
	AntiCollSelect2Cmd TagCommand = 0x95

	// AntiCollSelect3Cmd sets cascade level 3 during anti-collision.
	AntiCollSelect3Cmd TagCommand = 0x97

	// RequestAnswerToResetCmd is a request command for Answer to Reset.
	RequestAnswerToResetCmd TagCommand = 0xE0

	// Commands for MIFARE Classic (Mifare Classic 1K data sheet, Section 9)

	// AuthKeyACmd is used to authenticate a block using key A.
	AuthKeyACmd TagCommand = 0x60

	// AuthKeyBCmd is used to authenticate a block using key B.
	AuthKeyBCmd TagCommand = 0x61

	// ReadBlockCmd is used to read a 16-byte block from an authenticated sector of a tag.
	ReadBlockCmd TagCommand = 0x30

	// WriteBlockCmd is used to write a 16-byte block to an authenticated sector of a tag.
	WriteBlockCmd TagCommand = 0xA0

	// DecrementBlockCmd is used to decrement a block value and store the result in the internal data register.
	DecrementBlockCmd TagCommand = 0xC0

	// IncrementBlockCmd is used to increment a block value and store the result in the internal data register.
	IncrementBlockCmd TagCommand = 0xC1

	// RestoreBlockCmd reads the contents of a block into the internal data register.
	RestoreBlockCmd TagCommand = 0xC2

	// TransferBlockCmd writes the contents of the internal data register to a block.
	TransferBlockCmd TagCommand = 0xB0
)
