package ui

import (
	"bytes"
	"log"
	"os/exec"
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

		sysinfo, err := exec.Command("radio-sysinfo.sh").Output()
		if err != nil {
			log.Printf("Failed to execute radio-sysinfo.sh: %v", err)
			time.Sleep(30 * time.Second)
			continue
		}

		for idx, line := range bytes.Split(sysinfo, []byte("\n")) {
			if _, err := lw.Formatf("line sysinfo %d %s", idx, line); err != nil {
				log.Printf("Failed to print sysinfo: %v", err)
			}

		}

		time.Sleep(10 * time.Second)
	}
}
