package ui

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/net/context"

	"github.com/studentkittens/eulenfunk/display"
	"github.com/studentkittens/eulenfunk/util"
)

// RunClock displays the current time in the "clock" window.
func RunClock(lw *display.LineWriter, width int, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		now := time.Now()
		hur, min, sec := now.Clock()
		yer, mon, day := now.Date()

		tm := util.Center(fmt.Sprintf("%02d:%02d:%02d", hur, min, sec), width, ' ')
		dt := util.Center(fmt.Sprintf("%d %s %d", day, mon.String(), yer), width, ' ')

		if err := lw.Line("clock", 1, tm); err != nil {
			log.Printf("Failed to send line 1 of clock: %v", err)
		}

		if err := lw.Line("clock", 2, dt); err != nil {
			log.Printf("Failed to send line 1 of clock: %v", err)
		}

		time.Sleep(1 * time.Second)
	}
}
