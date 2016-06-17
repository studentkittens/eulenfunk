MAKE=make
BIN=/usr/bin

all: install

build:
	go install
	@cd driver && $(MAKE) --no-print-directory

install: build
	cp config/radio-sysinfo.sh $(BIN)
	cp eulenfunk $(BIN)
	cp driver/radio-led $(BIN)
	cp driver/radio-lcd $(BIN)
	cp driver/radio-rotary $(BIN)
