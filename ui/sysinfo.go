package menu

import (
	"log"
	"time"

	"golang.org/x/net/context"

	"github.com/studentkittens/eulenfunk/display"
)

func RunSysinfo(lw *display.LineWriter, width int, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if _, err := lw.Formatf("line sysinfo 2 I make Schuhu!"); err != nil {
			log.Printf("Failed to print sysinfo: %v", err)
		}

		time.Sleep(5 * time.Second)
	}
}
