package ui

// ClickEntry is a menu entry that can be clicked
// with the rotary switch.
type ClickEntry struct {
	// Text is the static text displayed in the menu line.
	Text string
	// ActionFunc is called when the entry is clicked.
	ActionFunc Action
}

// Render will display the text and a indicator if the
// entry is currently selected.
func (en *ClickEntry) Render(w int, active bool) string {
	prefix := "  "
	if active {
		prefix = "❤ "
	}

	return prefix + en.Text
}

// Action calls ActionFunc.
func (en *ClickEntry) Action() error {
	if en.ActionFunc != nil {
		return en.ActionFunc()
	}

	return nil
}

// Selectable is always true.
func (en *ClickEntry) Selectable() bool {
	return true
}
