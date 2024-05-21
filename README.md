# GoRFID

GoRFID is a demo project that uses my Go MFRC522 library to clone MIFARE Classic 1k cards.
The library is a TinyGo driver for the RFID-RC522 (MFRC522) module.

## Installation

TODO

## Usage

```go
package main

import (
	"github.com/erazemk/gorfid/mfrc522"
)

func main() {
	rfid := mfrc522.Init()
	println(rfid.Version())
}
```

TODO: More useful example

### Using with GoLand

TODO: Info about how the code can be run in GoLand

## Example setup

I used the following components:
- Arduino Nano
- RFID-RC522 module
- Two basic RFID cards
- A push button
- A 10k Ohm resistor
- Nine jumper wires

Any device that is compatible with TinyGo can be used, as long as it supports SPI, UART or I2C,
since that is what the MFRC522 module uses.
In some cases, you can also skip the button and resistor if the board includes one.

Since I used an Arduino Nano and set communication over SPI, some pins on the module needed to be
connected to specific pins on the board, because the Arduino Nano has a predefined SPI interface.

TODO: Table of pins

TODO: Image of the setup

TODO: Maybe a diagram of the connections

TODO: Explanation of the setup

## Extra info

- [MFRC522 datasheet](https://www.nxp.com/docs/en/data-sheet/MFRC522.pdf)
- [Arduino Nano Pinout](https://docs.arduino.cc/hardware/nano/)

## License

The MFRC522 library is licensed under the [GNU LGPLv3](LICENSE) license, but
the demo project code can be freely used and modified.
