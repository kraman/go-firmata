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

type SerialSubCommand byte

// Configure a builtin or soft serial port. This command must be called before sending serial data.
// Set txPin and rxPin to 0x00 for builtin serial ports.
func (c *FirmataClient) SerialConfig(port SerialPort, baud int, txPin byte, rxPin byte) (err error) {
	baudBytes := intto7Bit(baud)
	bufferSize := intto7Bit(1024)
	termChar := to7Bit('\n')
	c.serialChan = make(chan string, 10)

	err = c.sendSysEx(Serial, byte(SerialConfig)|byte(port),
		baudBytes[0], baudBytes[1], baudBytes[2],
		bufferSize[0], bufferSize[1], bufferSize[2],
		termChar[0], termChar[1])
	return
}

// Get channel for incoming serial data
func (c *FirmataClient) GetSerialData() <-chan string {
	return c.serialChan
}

func (c *FirmataClient) parseSerialResponse(data7bit []byte) {

	data := make([]byte, 0)
	for i := 1; i < len(data7bit); i = i + 2 {
		data = append(data, byte(from7Bit(data7bit[i], data7bit[i+1])))
	}
	select {
	case c.serialChan <- string(data):
	default:
		c.Log.Print("Serial data buffer overflow. No listener?")
	}
}
