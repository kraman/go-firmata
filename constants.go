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
	"fmt"
)

type FirmataCommand byte
type SysExCommand byte
type PinMode byte
type SerialPort byte

const (
	ProtocolMajorVersion = 2
	ProtocolMinorVersion = 3

	// max number of data bytes in non-Sysex messages
	MaxDataBytes = 32

	// message command bytes (128-255/0x80-0xFF)

	DigitalMessage     FirmataCommand = 0x90 // send data for a digital pin
	AnalogMessage      FirmataCommand = 0xE0 // send data for an analog pin (or PWM)
	EnableAnalogInput  FirmataCommand = 0xC0 // enable analog input by pin #
	EnableDigitalInput FirmataCommand = 0xD0 // enable digital input by port pair
	SetPinMode         FirmataCommand = 0xF4 // set a pin to INPUT/OUTPUT/PWM/etc
	ReportVersion      FirmataCommand = 0xF9 // report protocol version
	SystemReset        FirmataCommand = 0xFF // reset from MIDI
	StartSysEx         FirmataCommand = 0xF0 // start a MIDI Sysex message
	EndSysEx           FirmataCommand = 0xF7 // end a MIDI Sysex message

	// extended command set using sysex (0-127/0x00-0x7F)
	/* 0x00-0x0F reserved for user-defined commands */
	ServoConfig           SysExCommand = 0x70 // set max angle, minPulse, maxPulse, freq
	StringData            SysExCommand = 0x71 // a string message with 14-bits per char
	ShiftData             SysExCommand = 0x75 // a bitstream to/from a shift register
	I2CRequest            SysExCommand = 0x76 // send an I2C read/write request
	I2CReply              SysExCommand = 0x77 // a reply to an I2C read request
	I2CConfig             SysExCommand = 0x78 // config I2C settings such as delay times and power pins
	ExtendedAnalog        SysExCommand = 0x6F // analog write (PWM, Servo, etc) to any pin
	PinStateQuery         SysExCommand = 0x6D // ask for a pin's current mode and value
	PinStateResponse      SysExCommand = 0x6E // reply with pin's current mode and value
	CapabilityQuery       SysExCommand = 0x6B // ask for supported modes and resolution of all pins
	CapabilityResponse    SysExCommand = 0x6C // reply with supported modes and resolution
	AnalogMappingQuery    SysExCommand = 0x69 // ask for mapping of analog to pin numbers
	AnalogMappingResponse SysExCommand = 0x6A // reply with mapping info
	ReportFirmware        SysExCommand = 0x79 // report name and version of the firmware
	SamplingInterval      SysExCommand = 0x7A // set the poll rate of the main loop
	SysExNonRealtime      SysExCommand = 0x7E // MIDI Reserved for non-realtime messages
	SysExRealtime         SysExCommand = 0x7F // MIDI Reserved for realtime messages
	Serial                SysExCommand = 0x60
	SysExSPI              SysExCommand = 0x80

	SerialConfig SerialSubCommand = 0x10
	SerialComm   SerialSubCommand = 0x20
	SerialFlush  SerialSubCommand = 0x30
	SerialClose  SerialSubCommand = 0x40

	SPIConfig SPISubCommand = 0x10
	SPIComm   SPISubCommand = 0x20

	SPI_MODE0 = 0x00
	SPI_MODE1 = 0x04
	SPI_MODE2 = 0x08
	SPI_MODE3 = 0x0C

	SoftSerial  SerialPort = 0x00
	HardSerial1 SerialPort = 0x01
	HardSerial2 SerialPort = 0x02
	HardSerial3 SerialPort = 0x03

	// pin modes
	Input  PinMode = 0x00
	Output PinMode = 0x01
	Analog PinMode = 0x02
	PWM    PinMode = 0x03
	Servo  PinMode = 0x04
	Shift  PinMode = 0x05
	I2C    PinMode = 0x06
	SPI    PinMode = 0x07
)

func (m PinMode) String() string {
	switch {
	case m == Input:
		return "INPUT"
	case m == Output:
		return "OUTPUT"
	case m == Analog:
		return "ANALOG"
	case m == PWM:
		return "PWM"
	case m == Servo:
		return "SERVO"
	case m == Shift:
		return "SHIFT"
	case m == I2C:
		return "I2C"
	}
	return "UNKNOWN"
}

func (c FirmataCommand) String() string {
	switch {
	case (c & 0xF0) == DigitalMessage:
		return fmt.Sprintf("DigitalMessage (0x%x)", byte(c))
	case (c & 0xF0) == AnalogMessage:
		return fmt.Sprintf("AnalogMessage (0x%x)", byte(c))
	case c == EnableAnalogInput:
		return fmt.Sprintf("EnableAnalogInput (0x%x)", byte(c))
	case c == EnableDigitalInput:
		return fmt.Sprintf("EnableDigitalInput (0x%x)", byte(c))
	case c == SetPinMode:
		return fmt.Sprintf("SetPinMode (0x%x)", byte(c))
	case c == ReportVersion:
		return fmt.Sprintf("ReportVersion (0x%x)", byte(c))
	case c == SystemReset:
		return fmt.Sprintf("SystemReset (0x%x)", byte(c))
	case c == StartSysEx:
		return fmt.Sprintf("StartSysEx (0x%x)", byte(c))
	case c == EndSysEx:
		return fmt.Sprintf("EndSysEx (0x%x)", byte(c))
	}
	return fmt.Sprintf("Unexpected command (0x%x)", byte(c))
}

func (c SysExCommand) String() string {
	switch {
	case c == ServoConfig:
		return fmt.Sprintf("ServoConfig (0x%x)", byte(c))
	case c == StringData:
		return fmt.Sprintf("StringData (0x%x)", byte(c))
	case c == ShiftData:
		return fmt.Sprintf("ShiftData (0x%x)", byte(c))
	case c == I2CRequest:
		return fmt.Sprintf("I2CRequest (0x%x)", byte(c))
	case c == I2CReply:
		return fmt.Sprintf("I2CReply (0x%x)", byte(c))
	case c == I2CConfig:
		return fmt.Sprintf("I2CConfig (0x%x)", byte(c))
	case c == ExtendedAnalog:
		return fmt.Sprintf("ExtendedAnalog (0x%x)", byte(c))
	case c == PinStateQuery:
		return fmt.Sprintf("PinStateQuery (0x%x)", byte(c))
	case c == PinStateResponse:
		return fmt.Sprintf("PinStateResponse (0x%x)", byte(c))
	case c == CapabilityQuery:
		return fmt.Sprintf("CapabilityQuery (0x%x)", byte(c))
	case c == CapabilityResponse:
		return fmt.Sprintf("CapabilityResponse (0x%x)", byte(c))
	case c == AnalogMappingQuery:
		return fmt.Sprintf("AnalogMappingQuery (0x%x)", byte(c))
	case c == AnalogMappingResponse:
		return fmt.Sprintf("AnalogMappingResponse (0x%x)", byte(c))
	case c == ReportFirmware:
		return fmt.Sprintf("ReportFirmware (0x%x)", byte(c))
	case c == SamplingInterval:
		return fmt.Sprintf("SamplingInterval (0x%x)", byte(c))
	case c == SysExNonRealtime:
		return fmt.Sprintf("SysExNonRealtime (0x%x)", byte(c))
	case c == SysExRealtime:
		return fmt.Sprintf("SysExRealtime (0x%x)", byte(c))
	case c == Serial:
		return fmt.Sprintf("Serial (0x%x)", byte(c))
	case c == SysExSPI:
		return fmt.Sprintf("SPI (0x%x)", byte(c))
	}
	return fmt.Sprintf("Unexpected SysEx command (0x%x)", byte(c))
}
