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
enum {
    GLYPH_HBAR = 0,
    GLYPH_PLAY = 1,
    GLYPH_PAUSE= 2,
    GLYPH_HEART= 3,
    GLYPH_CROSS = 4,
    GLYPH_CHECK = 5,
    GLYPH_STOP = 6,
    GLYPH_CACTUS = 7
};

// See also: 
// https://omerk.github.io/lcdchargen/

static uint8_t GlyphDataHBar[8] = {
    0b00000,
    0b00000,
    0b11111,
    0b11111,
    0b11111,
    0b11111,
    0b00000,
    0b00000
};

static uint8_t GlyphDataPlay[8] = {
    0b00000,
    0b10000,
    0b11000,
    0b11100,
    0b11110,
    0b11100,
    0b11000,
    0b10000
};

static uint8_t GlyphDataPause[8] = {
    0b00000,
    0b11011,
    0b11011,
    0b11011,
    0b11011,
    0b11011,
    0b11011,
    0b0000
};

static uint8_t GlyphDataHeart[8] = {
    0b00000,
    0b00000,
    0b01010,
    0b11111,
    0b11111,
    0b01110,
    0b00100,
    0b00000
};

static uint8_t GlyphDataCross[8] = {
    0b10001,
    0b11011,
    0b01110,
    0b00100,
    0b01110,
    0b11011,
    0b10001,
    0b00000
};

static uint8_t GlyphDataCheck[8] = {
    0b00000,
    0b00001,
    0b00011,
    0b10110,
    0b11100,
    0b01000,
    0b00000,
    0b00000
};

static uint8_t GlyphDataStop[8] = {
    0b00000,
    0b00000,
    0b00000,
    0b01110,
    0b01110,
    0b01110,
    0b00000,
    0b00000
};

static uint8_t GlyphDataCactus[8] = {
    0b00000,
    0b00001,
    0b00101,
    0b10101,
    0b10110,
    0b01100,
    0b00100,
    0b00100
};


static int read_from_stdin(int handle) {
    const int n = LCD_WIDTH + 4 + 1;

    char matrix[LCD_HEIGHT][LCD_WIDTH];
    for(int y = 0; y < LCD_HEIGHT; y++) {
        for(int x = 0; x < LCD_WIDTH; x++) {
            matrix[y][x] = 0;
        }
    }

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


        first_space++;
        int i = offset;
        for(; i < LCD_WIDTH && *first_space; i++, first_space++) {
            char c = *first_space;
            if(c != matrix[lineno][i]) {
                lcdPosition(handle, i, lineno);
                lcdPutchar(handle, *first_space);
                matrix[lineno][i] = c;
            }
        }

        if(!offset_given) {
            for(; i < LCD_WIDTH; i++) {
                if(' ' != matrix[lineno][i]) {
                    lcdPosition(handle, i, lineno);
                    lcdPutchar(handle, ' ');
                    matrix[lineno][i] = ' ';
                }
            }
        }
    }

    return EXIT_SUCCESS;
}


int usage(const char *name) {
    fprintf(stderr, "Usage %s [print-charset]", name);
    return EXIT_FAILURE;
}

int print_charset(int handle) {
    int c = 0;

    for(;;) {
        fprintf(stderr, "Printing bytes %d-%d\n", c, c + LCD_WIDTH * LCD_HEIGHT);
        for(int y = 0; y < LCD_HEIGHT; y++) {
            for(int x = 0; x < LCD_WIDTH; x++) {
                lcdPosition(handle, x, y);
                lcdPutchar(handle, c++);
                if(c == 256) {
                    goto done;
                }
            }
        } 

        delay(2000);
        lcdClear(handle);
    }

done:
    return EXIT_SUCCESS;
}

int main(int argc, char **argv) {
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
    lcdCharDef(handle, GLYPH_HBAR,   GlyphDataHBar);
    lcdCharDef(handle, GLYPH_PLAY,   GlyphDataPlay);
    lcdCharDef(handle, GLYPH_PAUSE,  GlyphDataPause);
    lcdCharDef(handle, GLYPH_HEART,  GlyphDataHeart);
    lcdCharDef(handle, GLYPH_CROSS,  GlyphDataCross);
    lcdCharDef(handle, GLYPH_CHECK,  GlyphDataCheck);
    lcdCharDef(handle, GLYPH_STOP,   GlyphDataStop);
    lcdCharDef(handle, GLYPH_CACTUS, GlyphDataCactus);

    if(argc >= 2) {
        if(!strcmp("print-charset", argv[1])) {
            return print_charset(handle);
        } else {
            return usage(argv[0]);
        }

        return EXIT_SUCCESS;
    }

    return read_from_stdin(handle);
}
