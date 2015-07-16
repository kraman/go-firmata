package main

import (
	"github.com/argandas/go-firmata"
	"time"
)

var led uint8 = 13

func main() {
	arduino, err := firmata.NewClient("COM24", 57600)
	if err != nil {
		panic(err)
	}

	// arduino.Verbose = true

	myDelay := time.Millisecond * 250

	arduino.SetPinMode(led, firmata.Output)
	for x := 0; x < 10; x++ {
		// Turn ON led
		arduino.DigitalWrite(led, true)
		arduino.Delay(myDelay)
		// Turn OFF led
		arduino.DigitalWrite(led, false)
		arduino.Delay(myDelay)

	}
	arduino.Close()
}
