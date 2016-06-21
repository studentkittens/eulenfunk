package ui

import (
	"bytes"
	"log"
	"os/exec"
	"time"

	"golang.org/x/net/context"

	"github.com/studentkittens/eulenfunk/display"
)

// RunSysinfo displays system information in the "sysinfo" window.
// The data is obtained from the "radio-sysinfo.sh" script.
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
			if err := lw.Line("sysinfo", idx, string(line)); err != nil {
				log.Printf("Failed to print sysinfo: %v", err)
			}

		}

		time.Sleep(5 * time.Second)
	}
}
