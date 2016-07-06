#include <stdio.h>
#include <stdint.h>
#include <stdbool.h>
#include <stdlib.h>
#include <string.h>
#include <wiringPi.h>
#include "softPwm.h"

#define PIN_RED   7
#define PIN_GREEN 9
#define PIN_BLUE  8

#define CLAMP(x, low, high)  (((x) > (high)) ? (high) : (((x) < (low)) ? (low) : (x)))

static void usage(char *name) {
    fprintf(stderr, "usage:\n");
    fprintf(stderr, "  %s on  ....... turn on LED (white)\n", name);
    fprintf(stderr, "  %s off ....... turn off LED\n", name);
    fprintf(stderr, "  %s cat ....... read rgb tuples from stdin\n", name);
    fprintf(stderr, "  %s rgb  r g b  Set LED color to r,g,b\n", name);
    fprintf(stderr, "  %s hex #FFFFFF Set LED color from hexstring\n", name);
    fprintf(stderr, "  %s fade ...... Show a fade for debugging\n", name);
    exit(1);
}

static int linearize(int v) {
    // Half linear and half quadratic mix:
    return ((v / 256.) * (v / 256.)) * 128 + (v / 128.);
}

static void set_rgb(int r, int g, int b) { 
    softPwmWrite(PIN_RED, linearize(r));
    softPwmWrite(PIN_GREEN, linearize(g));
    softPwmWrite(PIN_BLUE, linearize(b));
}

static int string_to_col(const char * arg) {
    return CLAMP(strtol(arg,NULL,10),0,255);
}

static void set_rgb_from_arr(const char ** arr) {
    int r = string_to_col(arr[0]);
    int g = string_to_col(arr[1]);
    int b = string_to_col(arr[2]);

    set_rgb(r, g, b);
}

static void hexstring_to_rgb(const char * str, unsigned * r, unsigned * g, unsigned * b) {
    char * is_err = NULL;
    if(!(str && r && g && b))
        return;

    union {
        uint32_t val;
        unsigned char arr[sizeof(uint32_t)];
    } rgb_int = {.val=0};

    rgb_int.val = strtoul(str,&is_err, 10);

    if(is_err == NULL || *is_err == '\n' || *is_err == 0) {
        // Little endian:
        *r = rgb_int.arr[2];
        *g = rgb_int.arr[1];
        *b = rgb_int.arr[0];
    }
}

int main(int argc, char **argv) {
    piHiPri(30);
    wiringPiSetup();

    if(argc < 2) {
	    usage(argv[0]);
    }

    softPwmCreate(PIN_RED, 0, 255);
    softPwmCreate(PIN_GREEN, 0, 255);
    softPwmCreate(PIN_BLUE, 0, 255);

    if(strcasecmp(argv[1], "cat") == 0) {
        const int size = 512;
        char buf[size + 1];
        const char * rgb[3] = {0,0,0};
        int i = 0;

        for(;;) {
            memset(buf, 0, size);
            if(fgets(buf,size,stdin) == NULL) {
                puts("Got EOF, quitting");
                break;
            }

            if(buf[0] == '#') {
                unsigned r = 0, g = 0, b = 0;
                hexstring_to_rgb(&buf[1], &r, &g, &b);
                set_rgb(r, g, b);
            } else {
                bool is_valid = true;
                char * node = buf;

                for(i = 0; i < 3; i++) {
                    rgb[i] = node;
                    node = strpbrk(node," \n");
                    if(node != NULL) {
                        *node++ = 0;
                    } else {
                        is_valid = false;
                        break;
                    }
                }

                if(is_valid) {
                    set_rgb_from_arr(rgb);
                }
            }
        }
    } else if(strcasecmp(argv[1], "fade") == 0) {
	for(int i = 0; i < 256; i++) {
		set_rgb(i, 0, 0);
		delay(10);
	}
	for(int i = 0; i < 256; i++) {
		set_rgb(0, i, 0);
		delay(10);
	}
	for(int i = 0; i < 256; i++) {
		set_rgb(0, 0, i);
		delay(10);
	}
	for(int i = 0; i < 256; i++) {
		set_rgb(i, i, i);
		delay(10);
	}
    } else if(strcasecmp(argv[1], "rgb") == 0 && argc > 4) {
        set_rgb_from_arr((const char**)&argv[2]);
    } else if(strcasecmp(argv[1], "hex") == 0) {
        unsigned r = 0, g = 0, b = 0;
        printf("%s\n",&argv[2][1]);
        hexstring_to_rgb(&argv[2][1], &r, &g, &b);
        set_rgb(r, g, b);
    } else if(strcasecmp(argv[1], "off")) {
        set_rgb(0, 0, 0);
    } else if(strcasecmp(argv[1], "on")) {
        set_rgb(0, 0, 0);
    } else {
        usage(argv[0]);
    }

    return EXIT_SUCCESS;
}
