# go-firmata
A Golang wrapper for [Firmata](https://www.arduino.cc/en/reference/firmata) on [Arduino](https://www.arduino.cc/) 

[![GoDoc](http://godoc.org/github.com/kraman/go-firmata?status.svg)](http://godoc.org/github.com/kraman/go-firmata)

## Installation

```bash
	go get github.com/kraman/go-firmata
```

## Usage

```go
package main

import (
	"github.com/kraman/go-firmata"
	"time"
)

var led uint8 = 13

func main() {
	arduino, err := firmata.NewClient("COM1", 57600)
	if err != nil {
		panic(err)
	}

	// arduino.Verbose = true

	myDelay := time.Millisecond * 250

	// Set led pin as output
	arduino.SetPinMode(led, firmata.Output)

	// Blink led 10 times
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

```