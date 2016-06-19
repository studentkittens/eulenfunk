package ui

import (
	"log"
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

	aboutScreen = []string{
		"EULENFUNK IS MADE BY",
		"  Christoph <qitta> Piechula (christoph@nullcat.de)",
		"  Christopher <sahib> Pahl (sahib@online.de)",
		"  Susanne <Trüffelkauz> Kießling (aggro@thene.org)",
	}

	screens = map[string][]string{
		"startup":  startupScreen,
		"shutdown": shutdownScreen,
		"about":    aboutScreen,
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

func drawStaticScreens(lw *display.LineWriter) error {
	for window, block := range screens {
		if err := drawBlock(lw, window, block); err != nil {
			return err
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
