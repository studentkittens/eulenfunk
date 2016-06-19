package ui

import (
	"github.com/studentkittens/eulenfunk/display"
)

type Menu struct {
	Name    string
	Entries []Entry
	Cursor  int

	lw *display.LineWriter
}

func NewMenu(name string, lw *display.LineWriter) (*Menu, error) {
	return &Menu{
		Name: name,
		lw:   lw,
	}, nil
}

func (mn *Menu) ActiveEntryName() string {
	if len(mn.Entries) == 0 {
		return ""
	}

	return mn.Entries[mn.Cursor].Name()
}

func (mn *Menu) scrollToNextSelectable(up bool) {
	for mn.Cursor >= 0 && mn.Cursor < len(mn.Entries) && !mn.Entries[mn.Cursor].Selectable() {
		if up {
			mn.Cursor++
		} else {
			mn.Cursor--
		}
	}

	// Clamp value if nothing suitable was found in that direction:
	if mn.Cursor < 0 {
		mn.Cursor = 0
	}

	if mn.Cursor >= len(mn.Entries) {
		mn.Cursor = len(mn.Entries) - 1
	}
}

func (mn *Menu) Scroll(move int) {
	mn.Cursor += move

	up := move >= 0
	mn.scrollToNextSelectable(up)

	// Check if we succeeded:
	if !mn.Entries[mn.Cursor].Selectable() {
		mn.scrollToNextSelectable(!up)
	}
}

func (mn *Menu) Display(width int) error {
	for pos, ClickEntry := range mn.Entries {
		line := ClickEntry.Render(width, pos == mn.Cursor)

		if _, err := mn.lw.Printf("line %s %d %s", mn.Name, pos, line); err != nil {
			return err
		}
	}

	return nil
}

func (mn *Menu) Click() error {
	if len(mn.Entries) == 0 {
		return nil
	}

	return mn.Entries[mn.Cursor].Action()
}
