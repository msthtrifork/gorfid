package main

import (
	"machine"
	"time"

	"gorfid/mfrc522"
)

type state = int

const (
	stateRead state = iota
	stateClone
	stateWrite
)

func main() {
	// Init reader
	rfid, err := mfrc522.Init(machine.D8, machine.D9, 10*time.Second)
	if err != nil {
		println("Failed to initialize MFRC522", err)
		return
	}
	defer rfid.Exit()

	ledRed := machine.D4
	ledGreen := machine.D5
	ledBlue := machine.D6
	button := machine.D7

	ledRed.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ledGreen.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ledBlue.Configure(machine.PinConfig{Mode: machine.PinOutput})
	button.Configure(machine.PinConfig{Mode: machine.PinInput})

	// Set up button to switch states
	var state state
	if err = button.SetInterrupt(machine.PinRising, func(machine.Pin) {
		switch state {
		case stateRead:
			state = stateClone
		case stateClone:
			state = stateWrite
		case stateWrite:
			state = stateRead
		}
		println("State:", state)
	}); err != nil {
		println("Failed to set button interrupt", err)
		return
	}

	// Main loop
	var tagData []byte
	for {
		switch state {
		case stateRead:
			// Green
			ledRed.High()
			ledGreen.Low()
			ledBlue.High()

			uid, err := rfid.ReadTagUUID()
			if err != nil {
				continue
			}

			println("Tag detected:", uid)
		case stateClone:
			// Blue
			ledRed.High()
			ledGreen.High()
			ledBlue.Low()

			uid, err := rfid.ReadTagUUID()
			if err != nil {
				continue
			}
			tagData = uid

			println("Tag cloned:", uid)
		case stateWrite:
			// Red
			ledRed.Low()
			ledGreen.High()
			ledBlue.High()

			if tagData == nil {
				println("No tag data to write")
				time.Sleep(1 * time.Second)
			}

			println("not implemented")
		}
	}
}
