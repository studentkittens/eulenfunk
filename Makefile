MAKE=make
BIN=/usr/bin

all: build

build:
	go build .
	@cd driver && $(MAKE) --no-print-directory

install: build
	cp eulenfunk $(BIN)
	cp driver/radio-led $(BIN)
	cp driver/radio-lcd $(BIN)
	cp driver/radio-rotary $(BIN)
