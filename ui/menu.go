package ui

import (
	"github.com/studentkittens/eulenfunk/display"
)

// Menu is a collection of Entry interfaces grouped in a window.
// It supports the usual menu semantics, i.e. moving down
// and clicking an entry.
type Menu struct {
	// Name is the name of the menu window.
	// Should start with menu-[...]
	Name string

	// Entries is a list of Entry interfaces
	Entries []Entry

	// Cursor is the current offset in the entries
	Cursor int

	lw *display.LineWriter
}

// NewMenu returns a new menu that will use `lw` to display itself
// on the window `name`.
func NewMenu(name string, lw *display.LineWriter) (*Menu, error) {
	return &Menu{
		Name: name,
		lw:   lw,
	}, nil
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

// Scroll moves the menu `move` down (or up if negative)
func (mn *Menu) Scroll(move int) {
	mn.Cursor += move

	up := move >= 0
	mn.scrollToNextSelectable(up)

	// Check if we succeeded:
	if !mn.Entries[mn.Cursor].Selectable() {
		mn.scrollToNextSelectable(!up)
	}
}

// Display draws the menu onto the display
func (mn *Menu) Display(width int) error {
	for pos, ClickEntry := range mn.Entries {
		line := ClickEntry.Render(width, pos == mn.Cursor)
		if err := mn.lw.Line(mn.Name, pos, line); err != nil {
			return err
		}
	}

	return nil
}

// Click executes the action under the cursor
func (mn *Menu) Click() error {
	if len(mn.Entries) == 0 {
		return nil
	}

	return mn.Entries[mn.Cursor].Action()
}
