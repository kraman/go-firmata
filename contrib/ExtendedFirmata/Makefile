#Copyright 2014 Krishna Raman
#
#Licensed under the Apache License, Version 2.0 (the "License");
#you may not use this file except in compliance with the License.
#You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
#Unless required by applicable law or agreed to in writing, software
#distributed under the License is distributed on an "AS IS" BASIS,
#WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#See the License for the specific language governing permissions and
#limitations under the License.

AVR_GCC=/Applications/Arduino.app/Contents/Resources/Java/hardware/tools/avr/bin/avr-gcc
AVR_GPP=/Applications/Arduino.app/Contents/Resources/Java/hardware/tools/avr/bin/avr-gcc
AVR_AR=/Applications/Arduino.app/Contents/Resources/Java/hardware/tools/avr/bin/avr-ar
AVR_OBJCOPY=/Applications/Arduino.app/Contents/Resources/Java/hardware/tools/avr/bin/avr-objcopy
AVRDUDE=/Applications/Arduino.app/Contents/Resources/Java/hardware/tools/avr/bin/avrdude -C/Applications/Arduino.app/Contents/Resources/Java/hardware/tools/avr/etc/avrdude.conf

ARDUINO_HW_INCLUDE=/Applications/Arduino.app/Contents/Resources/Java/hardware
ARDUINO_LIB_INCLUDE=/Applications/Arduino.app/Contents/Resources/Java/libraries

AVR_LIBC=obj/core/malloc.o obj/core/realloc.o
HW_C_OBJECTS=obj/core/WInterrupts.o obj/core/wiring.o obj/core/wiring_analog.o obj/core/wiring_digital.o obj/core/wiring_pulse.o obj/core/wiring_shift.o 
HW_CPP_OBJECTS=obj/core/CDC.o obj/core/HardwareSerial.o obj/core/HID.o obj/core/IPAddress.o obj/core/main.o obj/core/new.o obj/core/Print.o obj/core/Stream.o \
	obj/core/Tone.o obj/core/USBCore.o obj/core/WMath.o obj/core/WString.o

INO=ExtendedFirmata
SOURCES=obj/ExtendedFirmata.o
LIBS=Servo Wire Firmata SPI

# Programming support using avrdude. Settings and variables.
MCU = atmega2560
F_CPU = 16000000
UPLOAD_RATE = 115200
AVRDUDE_PROGRAMMER=wiring
AVRDUDE_PORT=/dev/tty.usbmodem1411
AVRDUDE_WRITE_FLASH=-U flash:w:obj/$(INO).hex:i
AVRDUDE_FLAGS=-F -p $(MCU) -P $(AVRDUDE_PORT) -c $(AVRDUDE_PROGRAMMER) -b $(UPLOAD_RATE) -D

build: dirs obj/core.a $(SOURCES) libraries cleanup
	$(AVR_GCC) -Os -Wl,--gc-sections,--relax -mmcu=atmega2560 -o obj/$(INO).elf obj/*.o obj/core/core.a -lm
	$(AVR_OBJCOPY) -O ihex -j .eeprom --set-section-flags=.eeprom=alloc,load --no-change-warnings --change-section-lma .eeprom=0 obj/$(INO).elf obj/$(INO).eep
	$(AVR_OBJCOPY) -O ihex -R .eeprom obj/$(INO).elf obj/$(INO).hex
	
dirs:
	mkdir -p obj/core
	ln -sf $(INO).ino $(INO).cpp

cleanup:
	rm $(INO).cpp

upload: build
	$(AVRDUDE) $(AVRDUDE_FLAGS) $(AVRDUDE_WRITE_FLASH)
	
libraries:
	pushd obj; \
	for i in $(LIBS); do \
		echo $$i; \
		CFILES=`find $(ARDUINO_LIB_INCLUDE)/$$i -name \\*.c`; \
		CPPFILES=`find $(ARDUINO_LIB_INCLUDE)/$$i -name \\*.cpp`; \
		INC=""; \
		for j in `find $(ARDUINO_LIB_INCLUDE)/$$i -type d | grep -v example`; do \
			INC="$$INC -I$$j"; \
		done; \
		for j in $$CFILES; do \
			$(AVR_GCC) -c -g -Os -Wall -ffunction-sections -fdata-sections -mmcu=atmega2560 \
				-DF_CPU=16000000L -MMD -DUSB_VID=null -DUSB_PID=null -DARDUINO=105 \
				-I$(ARDUINO_HW_INCLUDE)/arduino/cores/arduino \
				-I$(ARDUINO_HW_INCLUDE)/arduino/variants/mega \
				$$INC $$j; \
		done; \
		for j in $$CPPFILES; do \
			$(AVR_GPP) -c -g -Os -Wall -fno-exceptions -ffunction-sections -fdata-sections -mmcu=atmega2560 \
				-DF_CPU=16000000L -MMD -DUSB_VID=null -DUSB_PID=null -DARDUINO=105 \
				-I$(ARDUINO_HW_INCLUDE)/arduino/cores/arduino \
				-I$(ARDUINO_HW_INCLUDE)/arduino/variants/mega \
				$$INC $$j; \
		done; \
	done; \
	popd

obj/core.a: $(AVR_LIBC) $(HW_C_OBJECTS) $(HW_CPP_OBJECTS)
	for i in $(AVR_LIBC) $(HW_C_OBJECTS) $(HW_CPP_OBJECTS); do \
		$(AVR_AR) rcs obj/core/core.a $$i; \
	done
	
$(SOURCES): obj/%.o: %.cpp
	$(AVR_GPP) -c -g -Os -Wall -fno-exceptions -ffunction-sections -fdata-sections \
		-mmcu=atmega2560 -DF_CPU=16000000L -MMD -DUSB_VID=null -DUSB_PID=null -DARDUINO=105 \
		-I$(ARDUINO_HW_INCLUDE)/arduino/cores/arduino \
		-I$(ARDUINO_HW_INCLUDE)/arduino/variants/mega \
		-I$(ARDUINO_LIB_INCLUDE)/Servo \
		-I$(ARDUINO_LIB_INCLUDE)/Wire \
		-I$(ARDUINO_LIB_INCLUDE)/Firmata \
		-I$(ARDUINO_LIB_INCLUDE)/SPI \
		$< -o $@

$(AVR_LIBC): obj/core/%.o: $(ARDUINO_HW_INCLUDE)/arduino/cores/arduino/avr-libc/%.c
	$(AVR_GCC) -c -g -Os -Wall -ffunction-sections -fdata-sections -fpermissive\
		-mmcu=atmega2560 -DF_CPU=16000000L -MMD -DUSB_VID=null -DUSB_PID=null -DARDUINO=105 \
		-I$(ARDUINO_HW_INCLUDE)/arduino/cores/arduino \
		-I$(ARDUINO_HW_INCLUDE)/arduino/variants/mega \
		$< -o $@		
		
$(HW_C_OBJECTS): obj/core/%.o: $(ARDUINO_HW_INCLUDE)/arduino/cores/arduino/%.c
	$(AVR_GCC) -c -g -Os -Wall -ffunction-sections -fdata-sections \
		-mmcu=atmega2560 -DF_CPU=16000000L -MMD -DUSB_VID=null -DUSB_PID=null -DARDUINO=105 \
		-I$(ARDUINO_HW_INCLUDE)/arduino/cores/arduino \
		-I$(ARDUINO_HW_INCLUDE)/arduino/variants/mega \
		$< -o $@		
	
$(HW_CPP_OBJECTS): obj/core/%.o: $(ARDUINO_HW_INCLUDE)/arduino/cores/arduino/%.cpp
	$(AVR_GPP) -c -g -Os -Wall -fno-exceptions -ffunction-sections -fdata-sections \
		-mmcu=atmega2560 -DF_CPU=16000000L -MMD -DUSB_VID=null -DUSB_PID=null -DARDUINO=105 \
		-I$(ARDUINO_HW_INCLUDE)/arduino/cores/arduino \
		-I$(ARDUINO_HW_INCLUDE)/arduino/variants/mega \
		$< -o $@
