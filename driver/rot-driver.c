#include <wiringPi.h> 
#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <stdbool.h>
#include <stdio.h>
#include <time.h>
#include <string.h>
#include <math.h>

#include <fcntl.h>
#include <unistd.h>
#include <sys/time.h>
#include <sys/types.h>
#include <sys/stat.h>

#define MAX_ENCODERS 8

#define ROTARY_PIN_A  12
#define ROTARY_PIN_B  13
#define ROTARY_BUTTON 14

#define ROTATY_INCREMENT (0.25)

struct encoder {
    int pin_a;
    int pin_b;
    int pin_btn;

    volatile double value;
    volatile int lastEncoded;

    volatile struct timeval button_time;
    volatile bool button_was_pressed;
    volatile bool button_was_released;
};

struct encoder encoders[MAX_ENCODERS];

int N_ENCODERS = 0;

void updateButton() {
    struct encoder *encoder = encoders;
    for(; encoder < encoders + N_ENCODERS; encoder++) {
        int button_press = digitalRead(encoder->pin_btn);
	if(button_press == LOW) {
		gettimeofday((struct timeval *)&encoder->button_time, NULL);
		encoder->button_was_pressed = true;
	} else {
		encoder->button_time.tv_sec = 0;
		encoder->button_time.tv_usec = 0;
		encoder->button_was_released = true;
	}
    }
} 

void updateEncoders() {
    struct encoder *encoder = encoders;
    for(; encoder < encoders + N_ENCODERS; encoder++) {
        int MSB = digitalRead(encoder->pin_a);
        int LSB = digitalRead(encoder->pin_b);

        int encoded = (MSB << 1) | LSB;
        int sum = (encoder->lastEncoded << 2) | encoded;

        if(sum == 0b1101 || sum == 0b0100 || sum == 0b0010 || sum == 0b1011) encoder->value -= ROTATY_INCREMENT;
        if(sum == 0b1110 || sum == 0b0111 || sum == 0b0001 || sum == 0b1000) encoder->value += ROTATY_INCREMENT;

        encoder->lastEncoded = encoded;
    }
}

struct encoder *setupencoder(int pin_a, int pin_b, int pin_btn) {
    if (N_ENCODERS > MAX_ENCODERS) {
        printf("Maximum number of encodered exceded: %i\n", MAX_ENCODERS);
        return NULL;
    }

    struct encoder *newencoder = encoders + N_ENCODERS++;
    newencoder->pin_a = pin_a;
    newencoder->pin_b = pin_b;
    newencoder->pin_btn = pin_btn;
    newencoder->value = 0;
    newencoder->lastEncoded = 0;
    newencoder->button_was_pressed = false;
    newencoder->button_was_released = false;

    pinMode(pin_a, INPUT);
    pinMode(pin_b, INPUT);
    pinMode(pin_btn, INPUT);

    pullUpDnControl(pin_a, PUD_UP);
    pullUpDnControl(pin_b, PUD_UP);
    pullUpDnControl(pin_btn, PUD_UP);

    wiringPiISR(pin_a,INT_EDGE_BOTH, updateEncoders);
    wiringPiISR(pin_b,INT_EDGE_BOTH, updateEncoders);
    wiringPiISR(pin_btn,INT_EDGE_BOTH, updateButton);

    return newencoder;
}

int main(void) {
    piHiPri(20);
    wiringPiSetup();

    struct encoder *encoder = setupencoder(
            ROTARY_PIN_A,
            ROTARY_PIN_B,
            ROTARY_BUTTON
    );

    long last_value = -1;
    double hold_time = 0;

    while (1) {
        //updateEncoders();

        // Divide by three since it counts 
        // about 3 value increments per "tick" which is too much.
        long curr = lround(encoder->value);
        if(curr != last_value) {
            printf("v %ld\n", curr);
	    fflush(stdout);
            last_value = curr;
        }

        if(encoder->button_was_pressed) {
            encoder->button_was_pressed = false;
            printf("p 1\n");
	    fflush(stdout);
        }

        if(encoder->button_time.tv_sec > 0) {
            struct timeval now;
            gettimeofday(&now, NULL);

            hold_time = now.tv_sec - encoder->button_time.tv_sec;
            hold_time += (now.tv_usec - encoder->button_time.tv_usec) / (1000.0 * 1000.0);
            printf("t %f\n", hold_time);
	    fflush(stdout);
        } else {
            hold_time = 0;
        }

        if(encoder->button_was_released) {
            encoder->button_was_released = false;
            printf("p 0\n");
	    fflush(stdout);
        }

        // Sleep for a short amount to save cpu-time:
        struct timespec delay = {
            .tv_sec = 0,
            .tv_nsec = 1000 * 1000 * 100,
        };

        nanosleep(&delay, NULL);
    }

    return EXIT_SUCCESS;
}
