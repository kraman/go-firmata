package firmata

import (
	"log"
)

type SerialSubCommand byte

func (c *FirmataClient) SerialConfig(baud int) (err error) {
	baudBytes := IntTo7Bit(baud)
	bufferSize := IntTo7Bit(1024)
	termChar := To7Bit('\n')
	c.serialChan = make(chan string, 10)

	err = c.sendSysEx(Serial, byte(SerialConfig)|byte(HardSerial1),
		baudBytes[0], baudBytes[1], baudBytes[2],
		bufferSize[0], bufferSize[1], bufferSize[2],
		termChar[0], termChar[1])
	return
}

func (c *FirmataClient) GetSerialData() <-chan string {
	return c.serialChan
}

func (c *FirmataClient) parseSerialResponse(data7bit []byte) {

	data := make([]byte, 0)
	for i := 1; i < len(data7bit); i = i + 2 {
		data = append(data, byte(From7Bit(data7bit[i], data7bit[i+1])))
	}
	select {
	case c.serialChan <- string(data):
	default:
		log.Printf("Serial data buffer overflow. No listener?")
	}
}
