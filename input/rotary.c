#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <unistd.h>
#include <wiringPi.h>

#define ROTARY_PIN_A 12
#define ROTARY_PIN_B 13
#define ROTARY_BTN   14

//17 pins / 2 pins per encoder = 8 maximum encoders
#define max_encoders 8

struct encoder
{
    int pin_a;
    int pin_b;
    volatile long value;
    volatile int last_value;
};

typedef struct {
    int pin_a, pin_b;
    volatile long last_value, value;
} RotaryEncoder;

RotaryEncoder ROTARY;

void rotary_encoder_update() {
    int MSB = digitalRead(ROTARY.pin_a);
    int LSB = digitalRead(ROTARY.pin_b);

    int encoded = (MSB << 1) | LSB;
    int sum = (ROTARY.last_value << 2) | encoded;

    if(sum == 0b1101 || sum == 0b0100 || sum == 0b0010 || sum == 0b1011) {
        ROTARY.value++;
    }

    if(sum == 0b1110 || sum == 0b0111 || sum == 0b0001 || sum == 0b1000) {
        ROTARY.value--;
    }

    ROTARY.last_value = encoded;
}

int rotary_encoder_setup(int pin_a, int pin_b) {
    ROTARY.pin_a = pin_a;
    ROTARY.pin_b = pin_b;
    ROTARY.value = ROTARY.last_value = 0;

    pinMode(pin_a, INPUT);
    pinMode(pin_b, INPUT);

    pullUpDnControl(pin_a, PUD_UP);
    pullUpDnControl(pin_b, PUD_UP);

    wiringPiISR(pin_a, INT_EDGE_BOTH, rotary_encoder_update);
    wiringPiISR(pin_b, INT_EDGE_BOTH, rotary_encoder_update);
}

int main() {
    if(geteuid() != 0) {
        fprintf(stderr, "Need to be root to run(sudo?)\n");
        exit(1);
    }

    if(wiringPiSetup() == -1) {
        exit(2);
    }

    RotaryEncoder re;
    long value;

    while(1) {
        rotary_encoder_update();
        long l = re.value;

        if(l != value) {
            printf("value: %d\n", (void *)l);
            value = l;
        }
    }

    return EXIT_SUCCESS;
}
