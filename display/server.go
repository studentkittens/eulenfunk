package display

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/studentkittens/eulenfunk/util"

	"golang.org/x/net/context"
)

// NOTE: Custom chars are repeated in 8-15;
//       use 8 instead of 0 (=> nul-byte) therefore.
const (
	GlyphHBar   = 8
	GlyphPlay   = 1
	GlyphPause  = 2
	GlyphHeart  = 3
	GlyphCross  = 4
	GlyphCheck  = 5
	GlyphStop   = 6
	GlyphCactus = 7
)

var unicodeToLCDCustom = map[rune]rune{
	// Real custom characters:
	'━': GlyphHBar,
	'▶': GlyphPlay,
	'⏸': GlyphPause,
	'❤': GlyphHeart,
	'×': GlyphCross,
	'✓': GlyphCheck,
	'⏹': GlyphStop,
	// Existing characters on the LCD:
	'ψ': GlyphCactus,
	'ä': 132,
	'Ä': 142,
	'ü': 129,
	'Ü': 152,
	'ö': 148,
	'Ö': 153,
	'ß': 224,
	'π': 237,
	'৹': 178,
}

func encode(s string) []rune {
	// Iterate by rune:
	encoded := []rune{}

	for _, rn := range s {
		b, ok := unicodeToLCDCustom[rn]
		if !ok {
			if rn > 255 {
				// Multibyte chars would be messed up anyways:
				b = '?'
			} else {
				b = rn
			}
		}

		encoded = append(encoded, b)
	}

	return encoded
}

/////////////////////////
// LINE IMPLEMENTATION //
/////////////////////////

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
// It can be called if events happended that are out of reach of `Line`.
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
		text += " ━❤━ "
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

///////////////////////////
// WINDOW IMPLEMENTATION //
///////////////////////////

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

///////////////////////////
// SERVER IMPLEMENTATION //
///////////////////////////

// Config gives the user to adjust some settings of displayd.
type Config struct {
	// Host of displayd (usually localhost)
	Host string

	// Port of displayd (usually 7777)
	Port int

	// Width is the number of runes per line in the LCD.
	Width int

	// Height is the number of lines on the LCD.
	Height int

	// DriverBinary is the name of the driver to write the output too.
	DriverBinary string

	// NoEncoding disables the special LCD encoding
	NoEncoding bool
}

type server struct {
	sync.Mutex
	Config     *Config
	Windows    map[string]*Window
	Active     *Window
	Quit       chan bool
	DriverPipe io.Writer
}

func (srv *server) renderToDriver() {
	srv.Lock()
	defer srv.Unlock()

	if srv.Active == nil {
		return
	}

	for idx, line := range srv.Active.Render() {
		dline := []byte(fmt.Sprintf("%d %s\n", idx, string(line)))
		if _, err := srv.DriverPipe.Write(dline); err != nil {
			log.Printf("Failed to write to driver: %v", err)
		}
	}
}

// newServer returns a displayd instance based on `cfg` and the cancel context `ctx`.
func newServer(cfg *Config, ctx context.Context) (*server, error) {
	cmd := exec.Command(cfg.DriverBinary)
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	srv := &server{
		Config:     cfg,
		Windows:    make(map[string]*Window),
		Quit:       make(chan bool, 1),
		DriverPipe: stdinPipe,
	}

	go func() {
		// Update the screen with about 7Hz
		ticker := time.NewTicker(150 * time.Millisecond)

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				srv.renderToDriver()
			}
		}
	}()

	return srv, nil
}

func (srv *server) createOrLookupWindow(name string) *Window {
	win, ok := srv.Windows[name]

	if !ok {
		log.Printf("Creating new window `%s`", name)
		win = NewWindow(
			name,
			srv.Config.Width, srv.Config.Height,
			!srv.Config.NoEncoding,
		)
		srv.Windows[name] = win
	}

	if srv.Active == nil {
		srv.Active = win
	}

	srv.Active.Switch()
	return win
}

func (srv *server) Switch(name string) {
	srv.Lock()
	defer srv.Unlock()

	win := srv.createOrLookupWindow(name)

	// Save a redraw just in case:
	if win == srv.Active {
		return
	}

	win.Switch()
	srv.Active = win
	return
}

func (srv *server) SetLine(name string, pos int, text string) error {
	srv.Lock()
	defer srv.Unlock()

	return srv.createOrLookupWindow(name).SetLine(pos, text)
}

func (srv *server) SetScrollDelay(name string, pos int, delay time.Duration) error {
	srv.Lock()
	defer srv.Unlock()

	return srv.createOrLookupWindow(name).SetScrollDelay(pos, delay)
}

func (srv *server) Move(window string, n int) {
	srv.Lock()
	defer srv.Unlock()

	srv.createOrLookupWindow(window).Move(n)
}

func (srv *server) Truncate(window string, n int) {
	srv.Lock()
	win := srv.createOrLookupWindow(window)
	nlines := win.Truncate(n)
	srv.Unlock()

	srv.renderToDriver()

	srv.Lock()
	win.NLines = nlines
	srv.Unlock()

	srv.renderToDriver()
}

func (srv *server) RenderMatrix() []byte {
	srv.Lock()
	defer srv.Unlock()

	if srv.Active == nil {
		return nil
	}

	out := ""
	for _, line := range srv.Active.Render() {
		out += string(line) + "\n"
	}

	return []byte(out)
}

//////////////////////
// NETWORK HANDLING //
//////////////////////

func handleSwitch(srv *server, line string) {
	name := ""
	if _, err := fmt.Sscanf(line, "switch %s", &name); err != nil {
		log.Printf("Failed to parse switch command `%s`: %v", line, err)
		return
	}

	srv.Switch(name)
}

func handleLine(srv *server, line string, split []string) {
	text := ""
	if len(split) >= 4 {
		text = split[3]
	}

	win, pos := "", 0

	if _, err := fmt.Sscanf(line, "line %s %d ", &win, &pos); err != nil {
		log.Printf("Failed to parse line command `%s`: %v", line, err)
		return
	}

	if err := srv.SetLine(win, pos, text); err != nil {
		log.Printf("Failed to set line: %v", err)
		return
	}
}

func handleScroll(srv *server, line string) {
	win, pos, durationSpec := "", 0, ""
	if _, err := fmt.Sscanf(line, "scroll %s %d %s", &win, &pos, &durationSpec); err != nil {
		log.Printf("Failed to parse scroll command `%s`: %v", line, err)
		return
	}

	duration, err := time.ParseDuration(durationSpec)
	if err != nil {
		log.Printf("Bad duration `%s`: %v", durationSpec, err)
		return
	}

	if err := srv.SetScrollDelay(win, pos, duration); err != nil {
		log.Printf("Cannot set scroll: %v", err)
		return
	}
}

func handleMove(srv *server, line string) {
	name, pos := "", 0
	if _, err := fmt.Sscanf(line, "move %s %d", &name, &pos); err != nil {
		log.Printf("Failed to parse move command `%s`: %v", line, err)
		return
	}

	srv.Move(name, pos)
}

func handleTruncate(srv *server, line string) {
	win, limit := "", 0
	if _, err := fmt.Sscanf(line, "truncate %s %d", &win, &limit); err != nil {
		log.Printf("Failed to parse move command `%s`: %v", line, err)
		return
	}

	srv.Truncate(win, limit)
}

func handleRender(srv *server, conn io.Writer) {
	matrix := srv.RenderMatrix()
	sizeBuf := make([]byte, 8)
	binary.LittleEndian.PutUint64(sizeBuf, uint64(len(matrix)))

	if _, err := conn.Write(sizeBuf); err != nil {
		log.Printf("Failed to respond rendered size: %v", err)
		return
	}

	if _, err := conn.Write(matrix); err != nil {
		log.Printf("Failed to respond rendered display: %v", err)
		return
	}
}

func handleSingle(srv *server, line string, conn io.Writer) {
	switch split := strings.SplitN(line, " ", 4); split[0] {
	case "switch":
		handleSwitch(srv, line)
	case "line":
		handleLine(srv, line, split)
	case "scroll":
		handleScroll(srv, line)
	case "move":
		handleMove(srv, line)
	case "truncate":
		handleTruncate(srv, line)
	case "render":
		// NOTE: This is only used for --dump, not for the actual driver.
		handleRender(srv, conn)
	case "close":
		return
	case "quit":
		srv.Quit <- true
		return
	default:
		log.Printf("Ignoring unknown command `%s`", line)
	}
}

func handleAll(srv *server, conn io.ReadWriteCloser) {
	scanner := bufio.NewScanner(conn)
	defer util.Closer(conn)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}

		handleSingle(srv, line, conn)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Reading connection failed: %v", err)
	}
}

func aborted(srv *server, ctx context.Context) bool {
	// Check if we were interrupted:
	select {
	case <-ctx.Done():
		return true
	case <-srv.Quit:
		return true
	default:
		return false
	}
}

// Run starts a new displayd server based on `cfg` and `ctx`.
func Run(cfg *Config, ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	lsn, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	log.Printf("Listening on %s", addr)

	defer util.Closer(lsn)

	srv, err := newServer(cfg, ctx)
	if err != nil {
		return err
	}

	for !aborted(srv, ctx) {
		if tcpLsn, ok := lsn.(*net.TCPListener); ok {
			if err := tcpLsn.SetDeadline(time.Now().Add(2 * time.Second)); err != nil {
				log.Printf("Setting deadline failed: %v", err)
				return err
			}
		}

		conn, err := lsn.Accept()
		if err, ok := err.(*net.OpError); ok && err.Timeout() {
			continue
		}

		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go handleAll(srv, conn)
	}

	return nil
}
