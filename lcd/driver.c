#include <stdint.h>
#include <wiringPi.h>
#include <lcd.h>

#define LCD_RS 7
#define LCD_E  8
#define LCD_D4 25
#define LCD_D5 24
#define LCD_D6 23
#define LCD_D7 18
#define LED_ON 15



int main(void) {
    wiringPiSetupGpio();

    int handle = lcdInit(
        4, 20,  4, LCD_RS, LCD_E, 
        LCD_D4, LCD_D5, LCD_D6, LCD_D7,
        0, 0, 0, 0
    );


    lcdHome(handle);
    lcdClear(handle);

    lcdPuts(handle, "Im a sheep, Hurr durr!");
}
