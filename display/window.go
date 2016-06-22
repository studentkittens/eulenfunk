package window

import (
	"fmt"
	"time"
)

// Window consists of a fixed number of lines and a handle name
type Window struct {
	// Name of the window
	Name string

	// Lines are all lines of
	// Initially, those are `Height` lines.
	Lines []*Line

	// NLines is the number of lines a window has
	// (might be less than len(Lines) due to op-truncate)
	NLines int

	// LineOffset is the current offset in Lines
	// (as modified by Move and Truncate)
	LineOffset int

	// Width is the number of runes which fits in one line
	Width int

	// Height is the number of lines that can be shown simultaneously.
	Height int

	// UseEncoding defines if a special LCD encoding shall be used.
	UseEncoding bool
}

// NewWindow returns a new window with the dimensions `w`x`h`, named by `name`.
func NewWindow(name string, w, h int, useEncoding bool) *Window {
	win := &Window{
		Name:        name,
		Width:       w,
		Height:      h,
		UseEncoding: useEncoding,
	}

	for i := 0; i < h; i++ {
		ln := NewLine(i, w)
		win.Lines = append(win.Lines, ln)
		win.NLines++
	}

	return win
}

// SetLine sets text of line `pos` to `text`.
// If the line does not exist yet it will be created.
func (win *Window) SetLine(pos int, text string) error {
	if pos < 0 {
		return fmt.Errorf("Bad line position %d", pos)
	}

	// For safety:
	if pos > 1024 {
		return fmt.Errorf("Only up to 1024 lines supported.")
	}

	// We need to extend:
	if pos >= len(win.Lines) {
		newLines := make([]*Line, pos+1)
		copy(newLines, win.Lines)

		// Create the intermediate lines:
		for i := len(win.Lines); i < len(newLines); i++ {
			newLines[i] = NewLine(i, win.Width)
		}

		win.Lines = newLines
	}

	win.NLines = len(win.Lines)
	win.Lines[pos].SetText(text, win.UseEncoding)
	return nil
}

// SetScrollDelay sets the scroll shift delay of line `pos` to `delay`.
func (win *Window) SetScrollDelay(pos int, delay time.Duration) error {
	if pos < 0 || pos >= win.NLines {
		return fmt.Errorf("Bad line position %d", pos)
	}

	win.Lines[pos].SetScrollDelay(delay)
	return nil
}

// Move moves the window contents vertically by `n`.
func (win *Window) Move(n int) {
	if n == 0 {
		// no-op
		return
	}

	max := win.NLines - win.Height

	if win.LineOffset+n > max {
		win.LineOffset = max
	} else {
		win.LineOffset += n
	}

	// Sanity:
	if win.LineOffset < 0 {
		win.LineOffset = 0
	}

	return
}

// Truncate cuts off the window after `n` lines.
func (win *Window) Truncate(n int) int {
	nlines := 0

	switch {
	case n < 0:
		win.LineOffset = 0
		nlines = 0
	case n > len(win.Lines):
		nlines = len(win.Lines)
	default:
		nlines = n
	}

	// Go back if needed:
	diff := nlines - win.NLines
	if diff < 0 {
		win.Move(diff)
	}

	// Clear remaining lines:
	for i := nlines; i < win.NLines; i++ {
		win.Lines[i].SetText("", win.UseEncoding)
	}

	return nlines
}

// Switch makes `win` to the active window.
func (win *Window) Switch() {
	for _, line := range win.Lines {
		line.Redraw()
	}
}

// Render returns the whole current LCD matrix as bytes.
func (win *Window) Render() [][]rune {
	hi := win.LineOffset + win.Height
	if hi > win.NLines {
		hi = win.NLines
	}

	out := [][]rune{}
	for _, line := range win.Lines[win.LineOffset:hi] {
		out = append(out, line.Render())
	}

	return out
}
