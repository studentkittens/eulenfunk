package menu

import (
	"log"
	"time"

	"github.com/studentkittens/eulenfunk/display"
)

func RunSysinfo(lw *display.LineWriter, width int, killCh <-chan bool) {
	for {
		select {
		case <-killCh:
			return
		}

		if _, err := lw.Formatf("line sysinfo 2 Im a raspberry. Schuhu!"); err != nil {
			log.Printf("Failed to print sysinfo: %v", err)
		}

		time.Sleep(5 * time.Second)
	}
}
