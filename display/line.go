package display

import (
	"sync"
	"time"
	"unicode/utf8"
)

// Line is a fixed width buffer with scrolling support
// It also supports special names for special symbols.
type Line struct {
	sync.Mutex
	Pos         int
	ScrollDelay time.Duration

	text []rune
	buf  []rune

	// current offset mod len(buf)
	scrollPos int
}

// NewLine returns a new line at `pos`, `w` runes long.
func NewLine(pos int, w int) *Line {
	ln := &Line{
		Pos:  pos,
		text: []rune{},
		buf:  make([]rune, w),
	}

	// Initial render:
	ln.Lock()
	ln.redraw()
	ln.Unlock()

	go func() {
		var delay time.Duration

		for {
			ln.Lock()
			{
				if ln.ScrollDelay == 0 {
					delay = 200 * time.Millisecond
					ln.scrollPos = 0
				} else {
					delay = ln.ScrollDelay
					if len(ln.text) > 0 {
						ln.scrollPos = (ln.scrollPos + 1) % len(ln.text)
					}
				}

				ln.redraw()
			}
			ln.Unlock()

			time.Sleep(delay)
		}
	}()

	return ln
}

func (ln *Line) redraw() {
	scroll(ln.buf, ln.text, ln.scrollPos)
}

// Redraw makes sure the line is up-to-date.
// It can be called if events happeneded that are out of reach of `Line`.
func (ln *Line) Redraw() {
	ln.Lock()
	defer ln.Unlock()

	ln.redraw()
}

// SetText sets and updates the text of `Line`.  If `useEncoding` is false the
// text is not converted to the special one-rune encoding of the LCD which is
// useful for debugging on a normal terminal.
func (ln *Line) SetText(text string, useEncoding bool) {
	ln.Lock()
	defer ln.Unlock()

	// Add a nice separtor symbol in between scroll borders:
	if utf8.RuneCountInString(text) > len(ln.buf) {
		text += "   ━❤━   "
	}

	var encodedText []rune
	if useEncoding {
		encodedText = encode(text)
	} else {
		// Just take the incoming encoding,
		// might render weirdly on LCD though.
		encodedText = []rune(text)
	}

	// Check if we need to re-render...
	if string(encodedText) != string(ln.text) {
		ln.scrollPos = 0
	}

	ln.text = encodedText
	ln.redraw()
}

// SetScrollDelay sets the scroll speed of the line (i.e. the delay between one
// "shift"). Shorter delay means faster scrolling.
func (ln *Line) SetScrollDelay(delay time.Duration) {
	ln.Lock()
	defer ln.Unlock()

	if delay == 0 {
		ln.scrollPos = 0
	}

	ln.ScrollDelay = delay
	ln.redraw()
}

func scroll(buf []rune, text []rune, m int) {
	for i := 0; i < len(buf); i++ {
		buf[i] = 0
	}

	if len(text) < len(buf) {
		copy(buf, text)
	} else {
		// Scrolling needed:
		n := copy(buf, text[m:])

		if n < len(buf) {
			// Some space left, copy from front text:
			copy(buf[n:], text[:m])
		}
	}
}

// Render returns the current line contents with fixed width
func (ln *Line) Render() []rune {
	ln.Lock()
	defer ln.Unlock()

	return ln.buf
}
