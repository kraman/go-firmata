package firmata

import (
	"bytes"
	"fmt"
	"log"
)

func (c *FirmataClient) parseSysEx(data []byte) {
	var cmd SysExCommand

	cmd = SysExCommand(data[0])
	log.Printf("Processing sysex %v\n", cmd)
	data = data[1:]
	switch {
	case cmd == StringData:
		log.Printf("String data: %v", string(data))
	case cmd == CapabilityResponse:
		dataBuf := bytes.NewBuffer(data)
		c.pinModes = make([]map[PinMode]interface{}, 0)

		pin := 0
		var err error
		var modes []byte
		for ; err == nil; modes, err = dataBuf.ReadBytes(127) {
			pinModes := make(map[PinMode]interface{})
			if len(modes) < 2 {
				continue
			}

			modes = modes[0 : len(modes)-1]
			for i := 0; i < len(modes); i = i + 2 {
				pinModes[PinMode(modes[i])] = modes[i+1]
			}
			c.pinModes = append(c.pinModes, pinModes)
			pin = pin + 1
		}
		log.Printf("Total pins: %v\n", pin-1)
		c.capabilityDone = true
	case cmd == AnalogMappingResponse:
		c.analogPinsChannelMap = make(map[int]byte)
		c.analogChannelPinsMap = make(map[byte]int)
		for pin, channel := range data {
			if channel != 127 {
				c.analogPinsChannelMap[pin] = channel
				c.analogChannelPinsMap[channel] = pin
			}
		}
		log.Printf("pin -> channel: %v\n", c.analogPinsChannelMap)
		c.analogMappingDone = true
	case cmd == ReportFirmware:
		c.firmwareVersion = make([]int, 2)
		c.firmwareVersion[0] = int(data[0])
		c.firmwareVersion[1] = int(data[1])
		data = data[2:]
		//log.Printf("in %v", data[2:])
		c.firmwareName = MultibyteString(data)
		log.Printf("Firmware: %v [%v.%v]", c.firmwareName, c.firmwareVersion[0], c.firmwareVersion[1])
		c.ready = true
		c.sendSysEx(AnalogMappingQuery)
		c.sendSysEx(CapabilityQuery)
	case cmd == Serial:
		c.parseSerialResponse(data)
	default:
		log.Printf("Discarding unexpected SysEx command %v", cmd)
	}
}

func (c *FirmataClient) sendSysEx(cmd SysExCommand, data ...byte) (err error) {
	var b bytes.Buffer

	b.WriteByte(byte(StartSysEx))
	b.WriteByte(byte(cmd))
	b.Write(data)
	b.WriteByte(byte(EndSysEx))

	bStr := ""
	for _, b := range b.Bytes() {
		bStr = bStr + fmt.Sprintf(" %#2x", b)
	}
	log.Printf("SysEx send %v\n", bStr)

	_, err = b.WriteTo(*(c.conn))
	return
}
