all: build

build:
	go build .
	cd led && make
	cd lcd && make

install: build
	cp eulenfunk /usr/bin
	cp lcd/driver.py /usr/bin/radio-lcd
	cp led/radio-led /usr/bin/radio-led
