go-firmata
==========

Arduino Firmata client for golang

Sample usage (blink internal led on pin 13):


    package main

    import (
      "time"
      firmata "github.com/baol/go-firmata"
    )

    func main() {
      c, err := firmata.NewClient("/dev/ttyUSB0", 57600)
      if err != nil {
        panic("Cannot open client")
        panic(err)
      }
      time.Sleep(time.Duration(1) * time.Second)
      c.DigitalWrite(13, true)
      time.Sleep(time.Duration(1) * time.Second)
      c.DigitalWrite(13, false)
    }
