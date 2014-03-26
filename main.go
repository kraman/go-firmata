package main

import (
	"fmt"
	"github.com/kraman/go-firmata/firmata"
	"time"
)

func main() {
	c, _ := firmata.NewClient("/dev/tty.usbmodem5d11", 57600)
	c.EnableAnalogInput(65, false)
	c.SerialConfig(38400)

	go func() {
		valueChan := c.GetValues()
		serialChan := c.GetSerialData()
		for {
			select {
			case v := <-valueChan:
				fmt.Printf("value: %v\n", v)
			case v := <-serialChan:
				fmt.Printf("serial: %v\n", v)
			}
		}
	}()

	time.Sleep(time.Minute)
}
