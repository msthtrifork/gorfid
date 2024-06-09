package mfrc522

// Register is an MFRC55 register.
// Register definitions are available in Chapter 9 of the MFRC55 datasheet.
type Register = byte

var (
	// Reserved00Reg is reserved for future use.
	// Reserved00Reg Register = 0x00

	// CommandReg starts and stops command execution.
	CommandReg Register = 0x01

	// ComIEnReg enables and disables interrupt request control bits.
	ComIEnReg Register = 0x02

	// DivIEnReg enables and disables interrupt request control bits.
	DivIEnReg Register = 0x03

	// ComIrqReg contains the interrupt request bits.
	ComIrqReg Register = 0x04

	// DivIrqReg contains the interrupt request bits.
	DivIrqReg Register = 0x05

	// ErrorReg contains the error status of the last executed command.
	ErrorReg Register = 0x06

	// Status1Reg contains the communication status bits.
	Status1Reg Register = 0x07

	// Status2Reg contains the receiver and transmitter status bits.
	Status2Reg Register = 0x08

	// FIFODataReg is the input and output of a 64-byte FIFO buffer.
	FIFODataReg Register = 0x09

	// FIFOLevelReg contains the number of bytes stored in the FIFO buffer.
	FIFOLevelReg Register = 0x0A

	// WaterLevelReg contains the level for FIFO underflow and overflow warning.
	WaterLevelReg Register = 0x0B

	// ControlReg contains miscellaneous control registers.
	ControlReg Register = 0x0C

	// BitFramingReg contains the adjustment for bit-oriented frames.
	BitFramingReg Register = 0x0D

	// CollReg contains the bit position of the first bit-collision detected on the RF interface.
	CollReg Register = 0x0E

	// Reserved02Reg is reserved for future use.
	// Reserved02Reg Register = 0x0F

	// Reserved03Reg is reserved for future use.
	// Reserved03Reg Register = 0x10

	// ModeReg contains the general modes for transmitting and receiving.
	ModeReg Register = 0x11

	// TxModeReg contains the transmit data rate and framing.
	TxModeReg Register = 0x12

	// RxModeReg contains the reception data rate and framing.
	RxModeReg Register = 0x13

	// TxControlReg controls the logical behavior of the antenna driver pins TX1 and TX2.
	TxControlReg Register = 0x14

	// TxASKReg controls the setting of the transmission modulation.
	TxASKReg Register = 0x15

	// TxSelReg selects the internal sources for the antenna driver.
	TxSelReg Register = 0x16

	// RxSelReg selects the internal receiver settings.
	RxSelReg Register = 0x17

	// RxThresholdReg selects the thresholds for the bit decoder.
	RxThresholdReg Register = 0x18

	// DemodReg controls the demodulator settings.
	DemodReg Register = 0x19

	// Reserved04Reg is reserved for future use.
	// Reserved04Reg Register = 0x1A

	// Reserved05Reg is reserved for future use.
	// Reserved05Reg Register = 0x1B

	// MfTxReg controls some MIFARE communication transmit parameters.
	MfTxReg Register = 0x1C

	// MfRxReg controls some MIFARE communication receive parameters.
	MfRxReg Register = 0x1D

	// Reserved06Reg is reserved for future use.
	// Reserved06Reg Register = 0x1E

	// SerialSpeedReg selects the speed of the serial UART interface.
	SerialSpeedReg Register = 0x1F

	// Reserved07Reg is reserved for future use.
	// Reserved07Reg Register = 0x20

	// CRCResultHighReg contains the upper byte of the CRC calculation.
	CRCResultHighReg Register = 0x21

	// CRCResultLowReg contains the lower byte of the CRC calculation.
	CRCResultLowReg Register = 0x22

	// Reserved08Reg is reserved for future use.
	// Reserved08Reg Register = 0x23

	// ModWidthReg controls the ModWidth setting.
	ModWidthReg Register = 0x24

	// Reserved09Reg is reserved for future use.
	// Reserved09Reg Register = 0x25

	// RFCfgReg configures the receiver gain.
	RFCfgReg Register = 0x26

	// GsNReg configures the conductance of the antenna driver pins TX1 and TX2 for modulation.
	GsNReg Register = 0x27

	// CWGsPReg configures the conductance of the p-driver output during periods of no modulation.
	CWGsPReg Register = 0x28

	// ModGsPReg configures the conductance of the p-driver output during periods of modulation.
	ModGsPReg Register = 0x29

	// TModeReg controls settings for the internal timer.
	TModeReg Register = 0x2A

	// TPrescalerReg controls settings for the internal timer.
	TPrescalerReg Register = 0x2B

	// TReloadHighReg sets the upper byte for the internal timer.
	TReloadHighReg Register = 0x2C

	// TReloadLowReg sets the lower byte for the internal timer.
	TReloadLowReg Register = 0x2D

	// TCounterValHighReg contains the upper byte for the internal timer.
	TCounterValHighReg Register = 0x2E

	// TCounterValLowReg contains the lower byte for the internal timer.
	TCounterValLowReg Register = 0x2F

	// Reserved10Reg is reserved for future use.
	// Reserved10Reg Register = 0x30

	// TestSel1Reg contains the general test signal configuration.
	TestSel1Reg Register = 0x31

	// TestSel2Reg contains the general test signal configuration and PRBS control.
	TestSel2Reg Register = 0x32

	// TestPinEnReg enables the pin output driver on pins D1 to D7.
	TestPinEnReg Register = 0x33

	// TestPinValueReg defines the values for D1 to D7 when it is used as an I/O bus.
	TestPinValueReg Register = 0x34

	// TestBusReg shows the status of the internal test bus.
	TestBusReg Register = 0x35

	// AutoTestReg controls the digital self-test.
	AutoTestReg Register = 0x36

	// VersionReg contains the software version number of the MFRC55.
	VersionReg Register = 0x37

	// AnalogTestReg controls the pins AUX1 and AUX2.
	AnalogTestReg Register = 0x38

	// TestDAC1Reg defines the test value for TestDAC1.
	TestDAC1Reg Register = 0x39

	// TestDAC2Reg defines the test value for TestDAC2.
	TestDAC2Reg Register = 0x3A

	// TestADCReg shows the value of ADC I and Q channels.
	TestADCReg Register = 0x3B

	// Reserved11Reg is reserved for production tests.
	// Reserved11Reg Register = 0x3C

	// Reserved12Reg is reserved for production tests.
	// Reserved12Reg Register = 0x3D

	// Reserved13Reg is reserved for production tests.
	// Reserved13Reg Register = 0x3E

	// Reserved14Reg is reserved for production tests.
	// Reserved14Reg Register = 0x3F
)
