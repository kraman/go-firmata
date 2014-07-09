// Copyright 2014 Krishna Raman
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package firmata

import (
  "code.google.com/p/log4go"
  "github.com/tarm/goserial"

  "fmt"
  "io"
  "time"
)

// Arduino Firmata client for golang
type FirmataClient struct {
  serialDev string
  baud      int
  conn      *io.ReadWriteCloser
  Log       *log4go.Logger

  protocolVersion []byte
  firmwareVersion []int
  firmwareName    string

  ready             bool
  analogMappingDone bool
  capabilityDone    bool

  digitalPinState [8]byte

  analogPinsChannelMap map[int]byte
  analogChannelPinsMap map[byte]int
  pinModes             []map[PinMode]interface{}

  valueChan  chan FirmataValue
  serialChan chan string
  spiChan    chan []byte
}

// Creates a new FirmataClient object and connects to the Arduino board
// over specified serial port. This function blocks till a connection is
// succesfullt established and pin mappings are retrieved.
func NewClient(dev string, baud int) (client *FirmataClient, err error) {
  var conn io.ReadWriteCloser

  c := &serial.Config{Name: dev, Baud: baud}
  conn, err = serial.OpenPort(c)
  if err != nil {
    client.Log.Critical(err)
    return
  }

  logger := make(log4go.Logger) 
  logger.AddFilter("stdout", log4go.INFO, log4go.NewConsoleLogWriter())
  client = &FirmataClient{
    serialDev: dev,
    baud:      baud,
    conn:      &conn,
    Log:       &logger,
  }
  go client.replyReader()

  conn.Write([]byte{byte(SystemReset)})
  t := time.NewTicker(time.Second)

  for !(client.ready && client.analogMappingDone && client.capabilityDone) {
    select {
    case <-t.C:
      //no-op
    case <-time.After(time.Second * 15):
      client.Log.Critical("No response in 30 seconds. Resetting arduino")
      conn.Write([]byte{byte(SystemReset)})
    case <-time.After(time.Second * 30):
      client.Log.Critical("Unable to initialize connection")
      conn.Close()
      client = nil
    }
  }

  client.Log.Info("Client ready to use")

  return
}

// Close the serial connection to properly clean up after ourselves
// Usage: defer client.Close()
func (c *FirmataClient) Close() {
  (*c.conn).Close()
}

// Sets the Pin mode (input, output, etc.) for the Arduino pin
func (c *FirmataClient) SetPinMode(pin byte, mode PinMode) (err error) {
  if c.pinModes[pin][mode] == nil {
    err = fmt.Errorf("Pin mode %v not supported by pin %v", mode, pin)
    return
  }
  cmd := []byte{byte(SetPinMode), (pin & 0x7F), byte(mode)}
  err = c.sendCommand(cmd)
  return
}

// Specified if a digital Pin should be watched for input.
// Values will be streamed back over a channel which can be retrieved by the GetValues() call
func (c *FirmataClient) EnableDigitalInput(pin uint, val bool) (err error) {
  if pin < 0 || pin > uint(len(c.pinModes)) {
    err = fmt.Errorf("Invalid pin number %v\n", pin)
    return
  }
  port := (pin / 8) & 0x7F
  pin = pin % 8

  if val {
    cmd := []byte{byte(EnableDigitalInput) | byte(port), 0x01}
    err = c.sendCommand(cmd)
  } else {
    cmd := []byte{byte(EnableDigitalInput) | byte(port), 0x00}
    err = c.sendCommand(cmd)
  }

  return
}

// Set the value of a digital pin
func (c *FirmataClient) DigitalWrite(pin uint, val bool) (err error) {
  if pin < 0 || pin > uint(len(c.pinModes)) && c.pinModes[pin][Output] != nil {
    err = fmt.Errorf("Invalid pin number %v\n", pin)
    return
  }
  port := (pin / 8) & 0x7F
  portData := &c.digitalPinState[port]
  pin = pin % 8

  if val {
    (*portData) = (*portData) | (1 << pin)
  } else {
    (*portData) = (*portData) & ^(1 << pin)
  }
  data := to7Bit(*(portData))
  cmd := []byte{byte(DigitalMessage) | byte(port), data[0], data[1]}
  err = c.sendCommand(cmd)
  return
}

// Specified if a analog Pin should be watched for input.
// Values will be streamed back over a channel which can be retrieved by the GetValues() call
func (c *FirmataClient) EnableAnalogInput(pin uint, val bool) (err error) {
  if pin < 0 || pin > uint(len(c.pinModes)) && c.pinModes[pin][Analog] != nil {
    err = fmt.Errorf("Invalid pin number %v\n", pin)
    return
  }

  ch := byte(c.analogPinsChannelMap[int(pin)])
  c.Log.Debug("Enable analog inout on pin %v channel %v", pin, ch)
  if val {
    cmd := []byte{byte(EnableAnalogInput) | ch, 0x01}
    err = c.sendCommand(cmd)
  } else {
    cmd := []byte{byte(EnableAnalogInput) | ch, 0x00}
    err = c.sendCommand(cmd)
  }

  return
}

// Set the value of a analog pin
func (c *FirmataClient) AnalogWrite(pin uint, pinData byte) (err error) {
  if pin < 0 || pin > uint(len(c.pinModes)) && c.pinModes[pin][Analog] != nil {
    err = fmt.Errorf("Invalid pin number %v\n", pin)
    return
  }

  data := to7Bit(pinData)
  cmd := []byte{byte(AnalogMessage) | byte(pin), data[0], data[1]}
  err = c.sendCommand(cmd)
  return
}

func (c *FirmataClient) sendCommand(cmd []byte) (err error) {
  bStr := ""
  for _, b := range cmd {
    bStr = bStr + fmt.Sprintf(" %#2x", b)
  }
  c.Log.Trace("Command send%v\n", bStr)

  _, err = (*c.conn).Write(cmd)
  return
}

// Sets the polling interval in milliseconds for analog pin samples
func (c *FirmataClient) SetAnalogSamplingInterval(ms byte) (err error) {
  data := to7Bit(ms)
  err = c.sendSysEx(SamplingInterval, data[0], data[1])
  return
}

// Get the channel to retrieve analog and digital pin values
func (c *FirmataClient) GetValues() <-chan FirmataValue {
  return c.valueChan
}
