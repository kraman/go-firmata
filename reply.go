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
	"bufio"
	"fmt"
)

type FirmataValue struct {
	valueType            FirmataCommand
	value                int
	analogChannelPinsMap map[byte]int
}

func (v FirmataValue) IsAnalog() bool {
	return (v.valueType & 0xF0) == AnalogMessage
}

func (v FirmataValue) GetAnalogValue() (pin int, val int, err error) {
	if !v.IsAnalog() {
		err = fmt.Errorf("Cannot get analog value for digital pin")
		return
	}
	pin = v.analogChannelPinsMap[byte(v.valueType & ^AnalogMessage)]
	val = v.value
	return
}

func (v FirmataValue) GetDigitalValue() (port byte, val map[byte]interface{}, err error) {
	if v.IsAnalog() {
		err = fmt.Errorf("Cannot get digital value for analog pin")
		return
	}

	port = byte(v.valueType & ^DigitalMessage)
	val = make(map[byte]interface{})
	mask := 0x01
	for i := byte(0); i < 8; i++ {
		val[port*8+i] = ((v.value & mask) > 0)
		mask = mask * 2
	}
	return
}

func (v FirmataValue) String() string {
	if v.IsAnalog() {
		p, v, _ := v.GetAnalogValue()
		return fmt.Sprintf("Analog value %v = %v", p, v)
	} else {
		p, v, _ := v.GetAnalogValue()
		return fmt.Sprintf("Digital port %v = %b", p, v)
	}
}

func (c *FirmataClient) replyReader() {
	r := bufio.NewReader(*c.conn)
	c.valueChan = make(chan FirmataValue)
	var init bool

	for {
		b, err := (r.ReadByte())
		if err != nil {
			c.Log.Critical(err)
			return
		}

		cmd := FirmataCommand(b)
    c.Log.Trace("Incoming cmd %v", cmd)
		if !init {
			if cmd != ReportVersion {
				c.Log.Debug("Discarding unexpected command byte %0d (not initialized)\n", b)
				continue
			} else {
				init = true
			}
		}
		
		switch {
		case cmd == ReportVersion:
			c.protocolVersion = make([]byte, 2)
			c.protocolVersion[0], err = r.ReadByte()
			c.protocolVersion[1], err = r.ReadByte()
			c.Log.Info("Protocol version: %d.%d", c.protocolVersion[0], c.protocolVersion[1])
		case cmd == StartSysEx:
			var sysExData []byte
			sysExData, err = r.ReadSlice(byte(EndSysEx))
			if err == nil {
				c.parseSysEx(sysExData[0 : len(sysExData)-1])
			}
		case (cmd&DigitalMessage) > 0 || byte(cmd&AnalogMessage) > 0:
			b1, _ := r.ReadByte()
			b2, _ := r.ReadByte()
			select {
			case c.valueChan <- FirmataValue{cmd, int(from7Bit(b1, b2)), c.analogChannelPinsMap}:
			}
		default:
			c.Log.Debug("Discarding unexpected command byte %0d\n", b)
		}
		if err != nil {
			c.Log.Critical(err)
			return
		}
	}
}
