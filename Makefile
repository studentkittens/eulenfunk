MAKE=make
BIN=/usr/bin

all: install

eulenfunk:
	go install

driver:
	@cd driver && $(MAKE) --no-print-directory

.PHONY: driver

build: eulenfunk driver

install: build
	cp config/scripts/*.sh $(BIN)
	cp config/systemd/*.service /usr/lib/systemd/system
	cp config/udev/*.rules /etc/udev/rules.d 

	systemctl daemon-reload

	cp driver/radio-led $(BIN)
	cp driver/radio-lcd $(BIN)
	cp driver/radio-rotary $(BIN)
