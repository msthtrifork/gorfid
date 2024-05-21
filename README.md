# GoRFID

GoRFID is a demo project that uses my Go MFRC522 library to clone MIFARE Classic 1k cards.
The library is a TinyGo driver for the RFID-RC522 (MFRC522) module.

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

### Using with GoLand

I write most of my code in GoLand, so this repo has some things already set up, such as the
Run config for flashing to the board.
If not using Arduino Nano, make sure to update the run configuration.

You can use the following plugins to make your life easier:
- [TinyGo](https://plugins.jetbrains.com/plugin/16915-tinygo)
- [Serial Port Monitor](https://plugins.jetbrains.com/plugin/8031-serial-port-monitor)

They're both from JetBrains, the first one understands TinyGo's `machine` package and how to flash
the boards, and the second one can be used to read the board's serial output.

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

| **MFRC522** | **Arduino Nano** |
|:-----------:|:----------------:|
|     SDA     |       D10        |
|     SCK     |       D13        |
|    MOSI     |       D11        |
|    MISO     |       D12        |
|     IRQ     |                  |
|     GND     |       GND        |
|     RST     |        D5        |
|    3.3V     |       3V3        |

The button was connected to 3V3, then using a 10k Ohm resistor to GND and with D8 for the signal.

<img src="res/breadboard.png" height="400" alt="Setup on a breadboard">

The basic idea of the setup is that once the board is running, it can read and display the UID of
the card that is placed on the module.
The button acts as a mode switch, switching from reading mode, to clone mode, to write mode.
Each of the modes has an associated LED color, so that you always know what mode you're in (using
the built-in LED).
Read mode is green, clone mode is blue, and write mode is red.
If switching to the clone mode, the reader will read the whole contents of the card (so you need to
keep it close until the light stops flashing), and it will remember the data.
If you then switch to the writing mode, it will write the data of the stored card to the new card.
Again, you need to keep the card close until the light stops flashing.
You can freely use reading mode, as the contents of the last cloned card will be stored as long
as the board is powered on.

## Extra info

- [MFRC522 datasheet](https://www.nxp.com/docs/en/data-sheet/MFRC522.pdf)
- [Arduino Nano Pinout](https://docs.arduino.cc/hardware/nano/)

## License

The MFRC522 library is licensed under the [GNU LGPLv3](LICENSE) license, but
the demo project code can be freely used and modified.
