package util

import (
	"bufio"
	"io"
	"log"
	"os/exec"
	"strconv"
	"time"
)

type Rotary struct {
	// Button is a channel that gets triggered when the button was on (true) and
	// when it has been released again (false)
	Button chan bool

	// Pressed is a channel that gets triggered when the user presses the rotary
	// switch. The duration passed is the time how long he holds it.
	Pressed chan time.Duration

	// Value is triggerd when the knob was turned.
	Value chan int

	stdout io.ReadCloser
}

func NewRotary() (*Rotary, error) {
	cmd := exec.Command("radio-rotary")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	rty := &Rotary{
		Button:  make(chan bool),
		Pressed: make(chan time.Duration),
		Value:   make(chan int),
		stdout:  stdout,
	}

	go func() {
		scn := bufio.NewScanner(stdout)

		for scn.Scan() {
			line := scn.Text()

			// Some bogus line:
			if len(line) < 3 {
				continue
			}

			switch line[0] {
			case 'v':
				// Value change:
				v, err := strconv.Atoi(line[2:])
				if err != nil {
					log.Printf("Bad rotary value `%s`: %v", line[2:], err)
					continue
				}

				rty.Value <- v
			case 'p':
				// Press state change:
				b, err := strconv.Atoi(line[2:])
				if err != nil {
					log.Printf("Bad press state `%s`: %v", line[2:], err)
					continue
				}

				if b == 0 {
					rty.Button <- false
				} else {
					rty.Button <- true
				}
			case 't':
				// Press time changed:
				t, err := strconv.ParseFloat(line[2:], 64)
				if err != nil {
					log.Printf("Bad press duration `%s`: %v", line[2:], err)
					continue
				}

				rty.Pressed <- time.Duration(t)
			}
		}

		if err := scn.Err(); err != nil && err != io.EOF {
			log.Printf("Failed to scan driver input: %v", err)
		}
	}()

	return rty, nil
}

func (rty *Rotary) Close() error {
	close(rty.Button)
	close(rty.Pressed)
	close(rty.Value)
	return rty.stdout.Close()
}
