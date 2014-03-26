package firmata

import (
	"github.com/tarm/goserial"

	"fmt"
	"io"
	"log"
	"time"
)

type FirmataClient struct {
	serialDev string
	baud      int
	conn      *io.ReadWriteCloser

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
}

func NewClient(dev string, baud int) (client *FirmataClient, err error) {
	var conn io.ReadWriteCloser

	c := &serial.Config{Name: dev, Baud: baud}
	conn, err = serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
		return
	}

	client = &FirmataClient{
		serialDev: dev,
		baud:      baud,
		conn:      &conn,
	}
	go client.replyReader()

	conn.Write([]byte{byte(SystemReset)})
	t := time.NewTicker(time.Second)

	for !(client.ready && client.analogMappingDone && client.capabilityDone) {
		select {
		case <-t.C:
			//no-op
		case <-time.After(time.Second * 30):
			log.Println("No response in 30 seconds. Resetting arduino")
			conn.Write([]byte{byte(SystemReset)})
		case <-time.After(time.Second * 70):
			log.Fatal("Unable to initialize connection")
			conn.Close()
			client = nil
		}
	}

	log.Println("Client ready to use")

	return
}

func (c *FirmataClient) SetPinMode(pin byte, mode PinMode) (err error) {
	if c.pinModes[pin][mode] == nil {
		err = fmt.Errorf("Pin mode %v not supported by pin %v", mode, pin)
		return
	}
	cmd := []byte{byte(SetPinMode), (pin & 0x7F), byte(mode)}
	err = c.sendCommand(cmd)
	return
}

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
	data := To7Bit(*(portData))
	cmd := []byte{byte(DigitalMessage) | byte(port), data[0], data[1]}
	err = c.sendCommand(cmd)
	return
}

func (c *FirmataClient) EnableAnalogInput(pin uint, val bool) (err error) {
	if pin < 0 || pin > uint(len(c.pinModes)) && c.pinModes[pin][Analog] != nil {
		err = fmt.Errorf("Invalid pin number %v\n", pin)
		return
	}

	ch := byte(c.analogPinsChannelMap[int(pin)])
	log.Printf("pin %v channel %v", pin, ch)
	if val {
		cmd := []byte{byte(EnableAnalogInput) | ch, 0x01}
		err = c.sendCommand(cmd)
	} else {
		cmd := []byte{byte(EnableAnalogInput) | ch, 0x00}
		err = c.sendCommand(cmd)
	}

	return
}

func (c *FirmataClient) AnalogWrite(pin uint, pinData byte) (err error) {
	if pin < 0 || pin > uint(len(c.pinModes)) && c.pinModes[pin][Analog] != nil {
		err = fmt.Errorf("Invalid pin number %v\n", pin)
		return
	}

	data := To7Bit(pinData)
	cmd := []byte{byte(AnalogMessage) | byte(pin), data[0], data[1]}
	err = c.sendCommand(cmd)
	return
}

func (c *FirmataClient) sendCommand(cmd []byte) (err error) {
	bStr := ""
	for _, b := range cmd {
		bStr = bStr + fmt.Sprintf(" %#2x", b)
	}
	log.Printf("Command send%v\n", bStr)

	_, err = (*c.conn).Write(cmd)
	return
}

func (c *FirmataClient) SetAnalogSamplingInterval(ms byte) (err error) {
	data := To7Bit(ms)
	err = c.sendSysEx(SamplingInterval, data[0], data[1])
	return
}

func (c *FirmataClient) GetValues() <-chan FirmataValue {
	return c.valueChan
}
