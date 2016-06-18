all: radio-led radio-lcd radio-rotary

# Im too lazy for make right now:
clean:
	rm -f radio-*

radio-led: clean
	gcc led-driver.c softPwm.c -o radio-led -lwiringPi -pthread -Os -Wall -Wextra

radio-lcd: clean
	gcc lcd-driver.c -o radio-lcd -lwiringPi -lwiringPiDev -Os -Wall -Wextra

radio-rotary: clean
	gcc rot-driver.c -o radio-rotary -lwiringPi -Os -Wall -Wextra -lm
