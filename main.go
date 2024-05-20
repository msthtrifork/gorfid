package main

import (
	"time"
)

func main() {
	for {
		println("Random text")
		time.Sleep(time.Millisecond * 1000)
	}
}
