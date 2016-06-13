#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <stdbool.h>
#include <string.h>

#include <wiringPi.h>
#include <lcd.h>

#define LCD_RS 7
#define LCD_E  8
#define LCD_D4 25
#define LCD_D5 24
#define LCD_D6 23
#define LCD_D7 18

#define LCD_WIDTH 20
#define LCD_HEIGHT 4

// Seems to need an extra prototype to silence a warning...
// (actual implementation is in wiringPiDev)
extern void lcdCharDef(const int fd, int index, unsigned char data [8]);

// Glyph indices (max 8)
#define GLYPH_HOURGLASS 0

// See also: 
// https://omerk.github.io/lcdchargen/

static uint8_t GlyphDataHourglass[8] = {
	0b11111,
	0b00100,
	0b01110,
	0b11011,
	0b11011,
	0b01110,
	0b00100,
	0b11111
};

static int read_from_stdin(int handle) {
    const int n = LCD_WIDTH + 4 + 1;
    char line[n];

    for(;;) {
        memset(line, 0, n);

        if(fgets(line, n, stdin) == NULL) {
            puts("Got EOF, quitting");  // well, probably.
            return EXIT_SUCCESS;
        }

        char *first_space = strchr(line, ' ');
        if(first_space == NULL) {
            continue;
        }

        char *newline = strchr(line, '\n');
        if(newline != NULL) {
            *newline = 0;
        }

        int lineno = 0;
        int offset = 0;
    	bool offset_given = false;

        switch(first_space - line) {
            case 1:
                // Just a line-number (this driver only supports 0-9 lines in theory).
                lineno = *line - '0';
                break;
            case 3:
                offset = 10 * (line[1] - '0') + (line[2] - '0');  
                break;
            case 2:
                offset = line[2] - '0';  
                break;
            default:
                // Some bad formatting going on...
                continue;
        }

        if(lineno >= LCD_HEIGHT) {
            continue;
        }

        lcdPosition(handle, offset, lineno);

	first_space++;
	int i = offset;
	for(; i < LCD_WIDTH && *first_space; i++) {
       	    lcdPutchar(handle, *first_space++);
	}

	if(!offset_given) {
		for(; i < LCD_WIDTH; i++) {
		    lcdPutchar(handle, ' ');
		}
	}
    }

    return EXIT_SUCCESS;
}


int main(void) {
    wiringPiSetupGpio();
    piHiPri(20);

    int handle = lcdInit(
        4, 20,                          // HxW
        4,                              // bit mode
        LCD_RS, LCD_E,                  // RS and Strobe
        LCD_D4, LCD_D5, LCD_D6, LCD_D7, // Data pins
        0, 0, 0, 0
    );

    lcdHome(handle);
    lcdClear(handle);

    // Custom glyph definitions:
    lcdCharDef(handle, GLYPH_HOURGLASS, GlyphDataHourglass);

    return read_from_stdin(handle);
}
