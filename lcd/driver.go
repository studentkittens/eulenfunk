package main

import (
	"log"
	"time"

	rpio "github.com/stianeikeland/go-rpio"
)

const (
	LCD_RS    = 7
	LCD_E     = 8
	LCD_D4    = 25
	LCD_D5    = 24
	LCD_D6    = 23
	LCD_D7    = 18
	LED_ON    = 15
	HIGH_BITS = 4
	LOW_BITS  = 0
)

type LcdData struct {
	DataMode bool
	Data     []byte
}

type LCD struct {
	En, Rw, Rs, D4, D5, D6, D7 rpio.Pin
	queue                      chan []LcdData
}

func byteToBitArray(b byte) [8]uint {
	var a [8]uint
	for i := uint(0); i < 8; i++ {
		v := (b & (1 << i) >> i)
		a[i] = uint(v)
	}
	return a
}

func (l *LCD) enable() {
	l.En.High()
	l.En.Low()
}

func (l *LCD) writeBits(bits [8]uint, row int) {
	var startBit int = row
	dList := []rpio.Pin{l.D4, l.D5, l.D6, l.D7}
	for i := 0; i < 4; i++ {
		// TODO: check
		if int(bits[startBit+i]) > 0 {
			dList[i].High()
		} else {
			dList[i].Low()
		}
	}
}

func (l *LCD) startLcdWorker() {
	for {
		w := <-l.queue
		for _, k := range w {
			if k.DataMode {
				l.Rs.High()
			}
			for _, c := range k.Data {
				l.writeByte(c)
			}
			if k.DataMode {
				l.Rs.Low()
			}

		}

	}
}

func (l *LCD) writeByte(ch byte) {
	bitArr := byteToBitArray(ch)
	l.D4.Low()
	l.D5.Low()
	l.D6.Low()
	l.D7.Low()
	l.writeBits(bitArr, HIGH_BITS)
	l.enable()
	l.writeBits(bitArr, LOW_BITS)
	l.enable()
	time.Sleep(500 * time.Nanosecond)
}

func (l *LCD) writeCommandByte8(ch byte) {
	l.Rs.Low()
	l.D4.Low()
	l.D5.Low()
	l.D6.Low()
	l.D7.Low()

	bitArr := byteToBitArray(ch)
	l.writeBits(bitArr, HIGH_BITS)
	l.enable()
}

// write string to lcd
func (l *LCD) WriteString(text string) {
	l.queue <- []LcdData{LcdData{DataMode: true, Data: []byte(text)}}
}

func main() {
	if err := rpio.Open(); err != nil {
		log.Printf("Failed to open gpio range: %v", err)
		return
	}

	//Rw = rpio.Pin(0) // ?
	lcd := &LCD{}
	lcd.En = rpio.Pin(LCD_E)
	lcd.Rs = rpio.Pin(LCD_RS)
	//lcd.D0 = rpio.Pin(0)
	//lcd.D1 = rpio.Pin(0)
	//lcd.D2 = rpio.Pin(0)
	//lcd.D3 = rpio.Pin(0)
	lcd.D4 = rpio.Pin(LCD_D4)
	lcd.D5 = rpio.Pin(LCD_D5)
	lcd.D6 = rpio.Pin(LCD_D6)
	lcd.D7 = rpio.Pin(LCD_D7)

	// init
	lcd.Rs.Low()
	lcd.writeCommandByte8(0x30)
	lcd.writeCommandByte8(0x30)
	lcd.writeCommandByte8(0x20)
	lcd.writeByte(0x08)
	lcd.writeByte(0x01)
	lcd.writeByte(0x0C)

	lcd.WriteString("Hello World")
	time.Sleep(2 * time.Second)
}
