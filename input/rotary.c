#include <wiringPi.h> 
#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <stdio.h>
#include <time.h>
#include <string.h>

#include <fcntl.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/stat.h>

#define MAX_ENCODERS 8

struct encoder {
    int pin_a;
    int pin_b;
    volatile long value;
    volatile int lastEncoded;
};

//Pre-allocate encoder objects on the stack so we don't have to 
//worry about freeing them
struct encoder encoders[MAX_ENCODERS];

/*
  Should be run for every rotary encoder you want to control
  Returns a pointer to the new rotary encoder structer
  The pointer will be NULL is the function failed for any reason
*/
struct encoder *setupencoder(int pin_a, int pin_b); 

int N_ENCODERS = 0;

void updateEncoders() {
    struct encoder *encoder = encoders;
    for(; encoder < encoders + N_ENCODERS; encoder++) {
        int MSB = digitalRead(encoder->pin_a);
        int LSB = digitalRead(encoder->pin_b);

        int encoded = (MSB << 1) | LSB;
        int sum = (encoder->lastEncoded << 2) | encoded;

        if(sum == 0b1101 || sum == 0b0100 || sum == 0b0010 || sum == 0b1011) encoder->value++;
        if(sum == 0b1110 || sum == 0b0111 || sum == 0b0001 || sum == 0b1000) encoder->value--;

        encoder->lastEncoded = encoded;
    }
}

struct encoder *setupencoder(int pin_a, int pin_b) {
    if (N_ENCODERS > MAX_ENCODERS) {
        printf("Maximum number of encodered exceded: %i\n", MAX_ENCODERS);
        return NULL;
    }

    struct encoder *newencoder = encoders + N_ENCODERS++;
    newencoder->pin_a = pin_a;
    newencoder->pin_b = pin_b;
    newencoder->value = 0;
    newencoder->lastEncoded = 0;

    pinMode(pin_a, INPUT);
    pinMode(pin_b, INPUT);
    pullUpDnControl(pin_a, PUD_UP);
    pullUpDnControl(pin_b, PUD_UP);
    wiringPiISR(pin_a,INT_EDGE_BOTH, updateEncoders);
    wiringPiISR(pin_b,INT_EDGE_BOTH, updateEncoders);

    return newencoder;
}

int main(int argc, char **argv) {
    if(argc < 2) {
        printf("Usage: radio-rotary [FIFO_PATH]\n");
        return 3;
    }

    wiringPiSetup();

    if(mkfifo(argv[1], 0644) != 0) {
        printf("Failed to create fifo at %s\n", argv[1]);
        return 4;
    }

    int fd = open(argv[1], O_WRONLY);
    if(fd < 0) {
        printf("Failed to open write fifo");
        return 5;
    }

    struct encoder *encoder = setupencoder(12,13);

    char num_buf[20];
    memset(num_buf, 0, sizeof(num_buf));

    long value;
    while (1) {
        updateEncoders();
        long l = encoder->value;
        if(l!=value) {
            printf("value: %d\n", (void *)l);
            value = l;
        }

        int n = snprintf(num_buf, sizeof(num_buf) - 1, "%ld", l);
        if(write(fd, num_buf, n) == -1) {
            printf("Failed to write to fifo.");
            return 6;
        }

        struct timespec delay = {
            .tv_sec = 0,
            .tv_nsec = 1000 * 1000 * 100,
        };

        nanosleep(&delay, NULL);
    }

    close(fd);
    return EXIT_SUCCESS;
}
