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
	defer rfid.Exit()
	
	for {
		// Read the UID of a card
		// ReadUID() blocks, so if a board is powerful enough, you can run this in a goroutine
		uid, err := rfid.ReadUID()
		if err != nil {
			println("Error reading UID:", err.Error())
			return
		}
		println("UID:", uid)
    }
}
```

### Using with GoLand

I write most of my code in GoLand, so this repo has some things already set up, such as the
Run config for flashing to the board.
If not using an Arduino Uno, make sure to update the run configuration.

You can use the following plugins to make your life easier:
- [TinyGo](https://plugins.jetbrains.com/plugin/16915-tinygo)
- [Serial Port Monitor](https://plugins.jetbrains.com/plugin/8031-serial-port-monitor)

They're both from JetBrains, the first one understands TinyGo's `machine` package and how to flash
the boards, and the second one can be used to read the board's serial output.

## Example setup

I used the following components:
- Arduino Uno
- RFID-RC522 module
- Two basic RFID cards
- A push button
- An RGB LED
- Four 10k Ohm resistors
- A lot of jumper wires

Any device that is compatible with TinyGo can be used, as long as it supports SPI, UART or I2C,
since that is what the MFRC522 module uses (currently the library only supports SPI).
In some cases, you can also skip the button and resistor if the board includes one.

Since I used an Arduino Uno and set communication over SPI, some pins on the module needed to be
connected to specific pins on the board, because Arduinos have predefined SPI pins.

| **MFRC522** | **Arduino Uno** |
|:-----------:|:---------------:|
|     SDA     |       D10       |
|     SCK     |       D13       |
|    MOSI     |       D11       |
|    MISO     |       D12       |
|     IRQ     |       D9        |
|     GND     |       GND       |
|     RST     |       D8        |
|    3.3V     |       3V3       |

The button was connected to 3V3, then using a 10k Ohm resistor to GND and with D7 for the signal.
The RGB LED was connected to D4 for red, D5 for green, and D6 for blue, with a 10k Ohm resistor
between each pin and the LED leg.
The common leg of the LED was connected to GND.

<img src="res/breadboard.png" height="400" alt="Setup on a breadboard">

The basic idea of the setup is that once the board is running, it can read and display the UID of
the card that is placed on the module.
The button acts as a mode switch, switching from reading mode, to clone mode, to write mode.
Each of the modes has an associated LED color, so that you always know what mode you're in.
Read mode is green, clone mode is blue, and write mode is red.
If switching to the clone mode, the reader will read the contents of the card, and it will remember
the data.
If you then switch to the writing mode, it will write the data of the stored card to the new card.
You can freely use reading mode, as the contents of the last cloned card will be stored as long
as the board is powered on.

## Extra info

- [MFRC522 data sheet](https://www.nxp.com/docs/en/data-sheet/MFRC522.pdf)
- [Arduino Uno Pinout](https://docs.arduino.cc/hardware/uno/)
- [Mifare Classic 1K data sheet](http://www.mouser.com/ds/2/302/MF1S503x-89574.pdf)

## License

The MFRC522 library is licensed under the [GNU LGPLv3](LICENSE) license, but
the demo project code can be freely used and modified.
The code is ported from the public domain
[MFRC522 library by Miguel Balboa](https://github.com/miguelbalboa/rfid).
