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

///////////////////////////
// SERVER IMPLEMENTATION //
///////////////////////////

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
		if _, err := srv.DriverPipe.Write([]byte(fmt.Sprintf("%d ", idx))); err != nil {
			log.Printf("Failed to write to driver: %v", err)
		}

		// Convert []rune to []byte manually:
		binaryLine := make([]byte, len(line))
		for off, rn := range line {
			binaryLine[off] = byte(rn)
		}

		binaryLine = append(binaryLine, '\n')

		if _, err := srv.DriverPipe.Write(binaryLine); err != nil {
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

func parseMoveTruncate(line string) (string, int, error) {
	name, dummy, pos := "", "", 0
	if _, err := fmt.Sscanf(line, "%s %s %d", &dummy, &name, &pos); err != nil {
		log.Printf("Failed to parse move command `%s`: %v", line, err)
		return "", 0, err
	}

	return name, pos, nil
}

func handleMove(srv *server, line string) {
	name, pos, err := parseMoveTruncate(line)
	if err != nil {
		return
	}

	srv.Move(name, pos)
}

func handleTruncate(srv *server, line string) {
	name, pos, err := parseMoveTruncate(line)
	if err != nil {
		return
	}

	srv.Truncate(name, pos)
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

func handleSingle(srv *server, line string, conn io.Writer) bool {
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
		return false
	case "quit":
		srv.Quit <- true
		return false
	default:
		log.Printf("Ignoring unknown command `%s`", line)
	}

	return true
}

func handleAll(srv *server, conn io.ReadWriteCloser) {
	scanner := bufio.NewScanner(conn)
	defer util.Closer(conn)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}

		if !handleSingle(srv, line, conn) {
			break
		}
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
