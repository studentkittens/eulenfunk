package display

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"golang.org/x/net/context"
)

const (
	GLYPH_HBAR  = 8
	GLYPH_PLAY  = 1
	GLYPH_PAUSE = 2
	GLYPH_HEART = 3
	GLYPH_CROSS = 4
	GLYPH_CHECK = 5
	GLYPH_STOP  = 6
)

var UnicodeToLCDCustom = map[rune]byte{
	'━': GLYPH_HBAR,
	'▶': GLYPH_PLAY,
	'⏸': GLYPH_PAUSE,
	'❤': GLYPH_HEART,
	'×': GLYPH_CROSS,
	'✓': GLYPH_CHECK,
	'⏹': GLYPH_STOP,
	'ä': 132,
	'Ä': 142,
	'ü': 129,
	'Ü': 152,
	'ö': 148,
	'Ö': 153,
	'ß': 224,
}

func encode(s string) []byte {
	// Iterate by rune:
	encoded := []byte{}

	for _, rn := range s {
		b, ok := UnicodeToLCDCustom[rn]
		if !ok {
			if rn > 255 {
				// Multibyte chars would be messed up anyways:
				b = '?'
			} else {
				b = byte(rn)
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

	// TODO: rather use runes:
	text       []byte
	buf        []byte
	scrollPos  int
	driverPipe io.Writer
}

func NewLine(pos int, w int, driverPipe io.Writer) *Line {
	ln := &Line{
		Pos:        pos,
		text:       []byte{},
		buf:        make([]byte, w),
		driverPipe: driverPipe,
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

func (ln *Line) Redraw() {
	ln.Lock()
	defer ln.Unlock()

	ln.redraw()
}

func (ln *Line) SetText(text string, useEncoding bool) {
	ln.Lock()
	defer ln.Unlock()

	if utf8.RuneCountInString(text) > len(ln.buf) {
		text += " ━━ "
	}

	var encodedText []byte

	if useEncoding {
		encodedText = encode(text)
	} else {
		// Just take the incoming encoding,
		// might render weirdly on LCD though.
		encodedText = []byte(text)
	}

	// Check if we need to re-render...
	if !bytes.Equal(encodedText, ln.text) {
		ln.scrollPos = 0
	}

	ln.text = encodedText
	ln.redraw()
}

func (ln *Line) SetScrollDelay(delay time.Duration) {
	ln.Lock()
	defer ln.Unlock()

	if delay == 0 {
		ln.scrollPos = 0
	}

	ln.ScrollDelay = delay
	ln.redraw()
}

func scroll(buf []byte, text []byte, m int) {
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

func (ln *Line) Render() []byte {
	ln.Lock()
	defer ln.Unlock()

	return ln.buf
}

///////////////////////////
// WINDOW IMPLEMENTATION //
///////////////////////////

// Window consists of a fixed number of lines and a handle name
type Window struct {
	Name  string
	Lines []*Line

	// NLines is the number of lines a window has
	// (might be less than len(Lines) due to op-truncate)
	NLines        int
	LineOffset    int
	Width, Height int
	DriverPipe    io.Writer

	UseEncoding bool
}

func NewWindow(name string, driverPipe io.Writer, w, h int, useEncoding bool) *Window {
	win := &Window{
		Name:        name,
		Width:       w,
		Height:      h,
		DriverPipe:  driverPipe,
		UseEncoding: useEncoding,
	}

	for i := 0; i < h; i++ {
		ln := NewLine(i, w, driverPipe)
		win.Lines = append(win.Lines, ln)
		win.NLines++
	}

	return win
}

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
			newLines[i] = NewLine(i, win.Width, win.DriverPipe)
		}

		win.Lines = newLines
		win.NLines = len(win.Lines)
	}

	win.Lines[pos].SetText(text, win.UseEncoding)
	return nil
}

func (win *Window) SetScrollDelay(pos int, delay time.Duration) error {
	if pos < 0 || pos >= win.NLines {
		return fmt.Errorf("Bad line position %d", pos)
	}

	win.Lines[pos].SetScrollDelay(delay)
	return nil
}

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

func (win *Window) Truncate(n int) {
	oldN := win.NLines

	switch {
	case n < 0:
		win.LineOffset = 0
	case n > len(win.Lines):
		win.NLines = len(win.Lines)
	default:
		win.NLines = n
	}

	diff := win.NLines - oldN
	if diff < 0 {
		win.Move(diff)
	}

	// Clear remaining lines:
	for i := win.NLines; i < oldN; i++ {
		win.Lines[i].SetText("", win.UseEncoding)
	}
}

func (win *Window) Switch() {
	for _, line := range win.Lines {
		line.Redraw()
	}
}

func (win *Window) Render() []byte {
	buf := &bytes.Buffer{}

	hi := win.LineOffset + win.Height
	if hi > win.NLines {
		hi = win.NLines
	}

	for _, line := range win.Lines[win.LineOffset:hi] {
		buf.Write(line.Render())
		buf.WriteRune('\n')
	}

	return buf.Bytes()
}

///////////////////////////
// SERVER IMPLEMENTATION //
///////////////////////////

type Config struct {
	Host         string
	Port         int
	Width        int
	Height       int
	DriverBinary string
	NoEncoding   bool
}

type Server struct {
	sync.Mutex
	Config     *Config
	Windows    map[string]*Window
	Active     *Window
	Quit       chan bool
	DriverPipe io.Writer
}

func (srv *Server) renderToDriver() {
	srv.Lock()
	defer srv.Unlock()

	if srv.Active == nil {
		return
	}

	// TODO: Make this nicer, possibly just
	//       render single lines and don't split what Render()
	//       did?

	lines, pos := srv.Active.Render(), 0
	width := srv.Config.Width

	for i := 0; i < len(lines); i += width + 1 {
		lpos := fmt.Sprintf("%d ", pos)

		buf := lines[i : i+width]

		log.Printf("%s%s", lpos, buf)
		if _, err := srv.DriverPipe.Write([]byte(lpos)); err != nil {
			log.Printf("Failed to write to driver: %v", err)
		}

		if _, err := srv.DriverPipe.Write(buf); err != nil {
			log.Printf("Failed to write to driver: %v", err)
		}

		if _, err := srv.DriverPipe.Write([]byte("\n")); err != nil {
			log.Printf("Failed to write to driver: %v", err)
		}

		pos++
	}
}

func NewServer(cfg *Config, ctx context.Context) (*Server, error) {
	cmd := exec.Command(cfg.DriverBinary)
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	srv := &Server{
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

func (srv *Server) createOrLookupWindow(name string) *Window {
	win, ok := srv.Windows[name]

	if !ok {
		log.Printf("Creating new window `%s`", name)
		win = NewWindow(
			name,
			srv.DriverPipe,
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

func (srv *Server) Switch(name string) {
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

func (srv *Server) SetLine(name string, pos int, text string) error {
	srv.Lock()
	defer srv.Unlock()

	return srv.createOrLookupWindow(name).SetLine(pos, text)
}

func (srv *Server) SetScrollDelay(name string, pos int, delay time.Duration) error {
	srv.Lock()
	defer srv.Unlock()

	return srv.createOrLookupWindow(name).SetScrollDelay(pos, delay)
}

func (srv *Server) Move(window string, n int) {
	srv.Lock()
	defer srv.Unlock()

	srv.createOrLookupWindow(window).Move(n)
}

func (srv *Server) Truncate(window string, n int) {
	srv.Lock()
	defer srv.Unlock()

	srv.createOrLookupWindow(window).Truncate(n)
}

func (srv *Server) Render() []byte {
	srv.Lock()
	defer srv.Unlock()

	if srv.Active == nil {
		return nil
	}

	return srv.Active.Render()
}

//////////////////////
// NETWORK HANDLING //
//////////////////////

func handleConn(srv *Server, conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	defer conn.Close()

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}

		switch split := strings.SplitN(line, " ", 4); split[0] {
		case "switch":
			name := ""
			if _, err := fmt.Sscanf(line, "switch %s", &name); err != nil {
				log.Printf("Failed to parse switch command `%s`: %v", line, err)
				continue
			}

			srv.Switch(name)
		case "line":
			text := ""
			if len(split) >= 4 {
				text = split[3]
			}

			win, pos := "", 0

			if _, err := fmt.Sscanf(line, "line %s %d ", &win, &pos); err != nil {
				log.Printf("Failed to parse line command `%s`: %v", line, err)
				continue
			}

			if err := srv.SetLine(win, pos, text); err != nil {
				log.Printf("Failed to set line: %v", err)
				continue
			}
		case "scroll":
			win, pos, durationSpec := "", 0, ""
			if _, err := fmt.Sscanf(line, "scroll %s %d %s", &win, &pos, &durationSpec); err != nil {
				log.Printf("Failed to parse scroll command `%s`: %v", line, err)
				continue
			}

			duration, err := time.ParseDuration(durationSpec)
			if err != nil {
				log.Printf("Bad duration `%s`: %v", durationSpec, err)
				continue
			}

			if err := srv.SetScrollDelay(win, pos, duration); err != nil {
				log.Printf("Cannot set scroll: %v", err)
				continue
			}
		case "move":
			name, pos := "", 0
			if _, err := fmt.Sscanf(line, "move %s %d", &name, &pos); err != nil {
				log.Printf("Failed to parse move command `%s`: %v", line, err)
				continue
			}

			srv.Move(name, pos)
		case "truncate":
			name, pos := "", 0
			if _, err := fmt.Sscanf(line, "truncate %s %d", &name, &pos); err != nil {
				log.Printf("Failed to parse move command `%s`: %v", line, err)
				continue
			}

			srv.Truncate(name, pos)
		case "render":
			// NOTE: This is only used for --dump, not for the actual driver.
			if _, err := conn.Write(srv.Render()); err != nil {
				log.Printf("Failed to respond rendered display: %v", err)
				continue
			}
		case "close":
			return
		case "quit":
			srv.Quit <- true
			return
		default:
			log.Printf("Ignoring unknown command `%s`", line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Reading connection failed: %v", err)
	}
}

func RunDaemon(cfg *Config, ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	lsn, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	log.Printf("Listening on %s", addr)

	defer lsn.Close()

	srv, err := NewServer(cfg, ctx)
	if err != nil {
		return err
	}

	for {
		// Check if we were interrupted:
		select {
		case <-ctx.Done():
			return nil
		case <-srv.Quit:
			return nil
		default:
			// We may continue normally.
		}

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

		go handleConn(srv, conn)
	}

	return nil
}
