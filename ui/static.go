package ui

import (
	"log"
	"unicode/utf8"

	"github.com/studentkittens/eulenfunk/display"
)

var (
	STARTUP_SCREEN = []string{
		"/ / / / / / / / / / / / / / / /",
		"WELCOME TO EULENFUNK",
		" GUT. ECHT. ANDERS. ",
		"/ / / / / / / / / / / / / / / /",
	}

	SHUTDOWN_SCREEN = []string{
		"SHUTTING DOWN - BYE!",
		"                    ",
		"PLEASE WAIT 1 MINUTE",
		"BEFORE POWERING OFF!",
	}

	ABOUT_SCREEN = []string{
		"EULENFUNK IS MADE BY",
		"  Christoph <qitta> Piechula (christoph@nullcat.de)",
		"  Christopher <sahib> Pahl (sahib@online.de)",
		"  Susanne <Trüffelkauz> Kießling (aggro@thene.org)",
	}

	SCREENS = map[string][]string{
		"startup":  STARTUP_SCREEN,
		"shutdown": SHUTDOWN_SCREEN,
		"about":    ABOUT_SCREEN,
	}
)

func drawBlock(lw *display.LineWriter, window string, block []string) error {
	for idx, line := range block {
		if _, err := lw.Formatf("line %s %d %s", window, idx, line); err != nil {
			return err
		}

		// Needs scrolling for proper display:
		if utf8.RuneCountInString(line) > 20 {
			if _, err := lw.Formatf("scroll %s %d 300ms", window, idx); err != nil {
				return err
			}
		}
	}

	return nil
}

func drawStaticScreens(lw *display.LineWriter) error {
	for window, block := range SCREENS {
		if err := drawBlock(lw, window, block); err != nil {
			return err
		}
	}

	return nil
}

func switchToStatic(lw *display.LineWriter, window string) {
	if _, ok := SCREENS[window]; !ok {
		log.Printf("No such static window `%s`", window)
	}

	if _, err := lw.Formatf("switch %s", window); err != nil {
		log.Printf("Failed to switch to static screen: %v", err)
	}
}
