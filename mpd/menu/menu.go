package menu

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/studentkittens/eulenfunk/display"
	"github.com/studentkittens/eulenfunk/util"
)

var MAIN_MENU = []string{
	"Exit",
	"Moodbar Lighting",
	"Poweroff",
	"Reboot",
}

type Entry struct {
	Text   string
	Active bool
	Action func()
}

type Menu struct {
	Name    string
	Entries []*Entry
	Cursor  int
}

// menu-create: create and switch to $Name; set all to scrolling
// menu-scroll: prefix curr with *, remove * of prev
// menu-render: output to display

func Run() error {
	cfg := &display.Config{
		Host: "localhost",
		Port: 7778,
	}

	lw, err := display.Connect(cfg)
	if err != nil {
		return err
	}

	defer lw.Close()

	rty, err := util.NewRotary()
	if err != nil {
		return err
	}

	defer rty.Close()

	go func() {
		for state := range rty.Button {
			if state {
				fmt.Println("Button pressed")
			} else {
				fmt.Println("Button released")
				// Activate current action...
			}
		}
	}()

	go func() {
		for duration := range rty.Pressed {
			window := ""

			// TODO:
			switch {
			case duration > 500*time.Millisecond:
			case duration > 2*time.Second:
			case duration > 10*time.Second:
			}

			if window != "" {
				if _, err := lw.Formatf("switch %s", window); err != nil {
					log.Printf("switch failed: %v", err)
				}
			}

		}
	}()

	go func() {
		for value := range rty.Value {
			fmt.Printf("Value: %d\n", value)
			if _, err := lw.Formatf("move menu %d", value); err != nil {
				log.Printf("move failed: %v", err)
			}
		}
	}()

	fmt.Println("Press CTRL-C to shut down")
	ctrlCh := make(chan os.Signal, 1)
	signal.Notify(ctrlCh, os.Interrupt)
	<-ctrlCh

	return nil
}
