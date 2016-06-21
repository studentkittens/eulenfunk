package ui

import (
	"log"
	"time"
	"unicode/utf8"

	"github.com/studentkittens/eulenfunk/display"
)

var (
	startupScreen = []string{
		"/ / / / / / / / / / / / / / / /",
		"WELCOME TO EULENFUNK",
		" GUT. ECHT. ANDERS. ",
		"/ / / / / / / / / / / / / / / /",
	}
	shutdownScreen = []string{
		"SHUTTING DOWN - BYE!",
		"                    ",
		"PLEASE WAIT 1 MINUTE",
		"BEFORE POWERING OFF!",
	}
	aboutNamesScreen = []string{
		"━━ MADE WITH ❤ BY ━━",
		"  Christoph Piechula",
		"  Christopher Pahl  ",
		"  Susanne Kießling  ",
	}
	aboutNicksScreen = []string{
		"━━ MADE WITH ❤ BY ━━",
		"christoph@nullcat.de",
		"     sahib@online.de",
		"  susanne@nullcat.de",
	}
	aboutFuncScreen = []string{
		"━━ MADE WITH ❤ BY ━━",
		" ψ Hardware Engineer",
		" ▶ Software Engineer",
		" ✓ Design Micrathene",
	}
	screens = map[string][][]string{
		"startup":  [][]string{startupScreen},
		"shutdown": [][]string{shutdownScreen},
		"about":    [][]string{aboutNamesScreen, aboutNicksScreen, aboutFuncScreen},
	}
)

func drawBlock(lw *display.LineWriter, window string, block []string) error {
	for idx, line := range block {
		if _, err := lw.Printf("line %s %d %s", window, idx, line); err != nil {
			return err
		}

		// Needs scrolling for proper display:
		if utf8.RuneCountInString(line) > 20 {
			if _, err := lw.Printf("scroll %s %d 500ms", window, idx); err != nil {
				return err
			}
		}
	}

	return nil
}

func pageThrough(lw *display.LineWriter, window string, blocks [][]string) {
	i := 0
	for {
		if err := drawBlock(lw, window, blocks[i]); err != nil {
			log.Printf("Failed to page through %s: %v", window, err)
		}

		time.Sleep(2 * time.Second)
		i = (i + 1) % len(blocks)
	}
}

func drawStaticScreens(lw *display.LineWriter) error {
	for window, blocks := range screens {
		if len(blocks) == 1 {
			// One static draw is enough:
			if err := drawBlock(lw, window, blocks[0]); err != nil {
				return err
			}
		} else {
			go pageThrough(lw, window, blocks)
		}
	}

	return nil
}

func switchToStatic(lw *display.LineWriter, window string) {
	if _, ok := screens[window]; !ok {
		log.Printf("No such static window `%s`", window)
	}

	if _, err := lw.Printf("switch %s", window); err != nil {
		log.Printf("Failed to switch to static screen: %v", err)
	}
}
