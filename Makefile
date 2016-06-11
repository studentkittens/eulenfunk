MAKE=make

all: build

build:
	go build .
	@cd led && $(MAKE) --no-print-directory
	@cd lcd && $(MAKE) --no-print-directory
	@cd rotary && $(MAKE) --no-print-directory

install: build
	cp eulenfunk /usr/bin
	cp lcd/driver.py /usr/bin/radio-lcd
	cp led/radio-led /usr/bin/radio-led
