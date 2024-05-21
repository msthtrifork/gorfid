package main

import (
	"gorfid/mfrc522"
)

func main() {
	rfid := mfrc522.Init()
	println(rfid.Version())
}
