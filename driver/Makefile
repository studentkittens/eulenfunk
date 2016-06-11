all: radio-led radio-lcd radio-rotary

radio-led:
	gcc led-driver.c -o radio-led -lwiringPi -pthread -Os -Wall -Wextra

radio-lcd:
	gcc led-driver.c -o radio-lcd -lwiringPi -lwiringPiDev -Os -Wall -Wextra

radio-rotary:
	gcc rot-driver.c -o radio-rotary -lwiringPi -Os -Wall -Wextra
