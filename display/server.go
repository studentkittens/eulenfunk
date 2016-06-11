package display

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

/////////////////////////
// LINE IMPLEMENTATION //
/////////////////////////

// Line is a fixed width buffer with scrolling support
// It also supports special names for special symbols.
type Line struct {
	sync.Mutex

	ScrollDelay time.Duration

	text       []byte
	buf        []byte
	scrollPos  int
	driverPipe io.Writer
}

func NewLine(w int, driverPipe io.Writer) *Line {
	ln := &Line{
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

	// TODO: output to driver here
	if _, err := ln.driverPipe.Write(ln.buf); err != nil {
		log.Printf("Failed to write to driver: %v", err)
	}
}

func (ln *Line) SetText(text string) {
	ln.Lock()
	defer ln.Unlock()

	if len(text) > len(ln.buf) {
		text += " -*- "
	}

	// Check if we need to re-render...
	btext := []byte(text)
	if !bytes.Equal(btext, ln.text) {
		ln.scrollPos = 0
	}

	ln.text = btext
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
	Name          string
	Lines         []*Line
	LineOffset    int
	Width, Height int
	DriverPipe    io.Writer
}

func NewWindow(name string, driverPipe io.Writer, w, h int) *Window {
	win := &Window{
		Name:       name,
		Width:      w,
		Height:     h,
		DriverPipe: driverPipe,
	}

	for i := 0; i < h; i++ {
		win.Lines = append(win.Lines, NewLine(w, driverPipe))
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

	x := len(win.Lines)
	// We need to extend:
	if pos >= len(win.Lines) {
		newLines := make([]*Line, pos+1)
		copy(newLines, win.Lines)

		// Create the intermediate lines:
		for i := len(win.Lines); i < len(newLines); i++ {
			newLines[i] = NewLine(win.Width, win.DriverPipe)
		}

		win.Lines = newLines
	}

	fmt.Println("SetLine", pos, len(win.Lines), x)

	win.Lines[pos].SetText(text)
	return nil
}

func (win *Window) SetScrollDelay(pos int, delay time.Duration) error {
	if pos < 0 || pos >= len(win.Lines) {
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

	max := len(win.Lines) - win.Height
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

func (win *Window) Render() []byte {
	buf := &bytes.Buffer{}

	hi := win.LineOffset + win.Height
	if hi > len(win.Lines) {
		hi = len(win.Lines)
	}

	fmt.Println(len(win.Lines), hi, win.LineOffset)
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
}

type Server struct {
	sync.Mutex
	Config     *Config
	Windows    map[string]*Window
	Active     *Window
	Quit       chan bool
	DriverPipe io.Writer
}

func NewServer(cfg *Config) (*Server, error) {
	cmd := exec.Command(cfg.DriverBinary)
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &Server{
		Config:     cfg,
		Windows:    make(map[string]*Window),
		Quit:       make(chan bool, 1),
		DriverPipe: stdinPipe,
	}, nil
}

func (srv *Server) AddWindow(name string) {
	srv.Lock()
	defer srv.Unlock()

	win, ok := srv.Windows[name]

	if !ok {
		log.Printf("Creating new window `%s`", name)
		win = NewWindow(name, srv.DriverPipe, srv.Config.Width, srv.Config.Height)
		srv.Windows[name] = win
	}

	if srv.Active == nil {
		srv.Active = win
	}
}

func (srv *Server) Switch(name string) error {
	srv.Lock()
	defer srv.Unlock()

	win, ok := srv.Windows[name]
	if !ok {
		return fmt.Errorf("No such window: %s", name)
	}

	srv.Active = win
	return nil
}

func (srv *Server) SetLine(pos int, text string) error {
	srv.Lock()
	defer srv.Unlock()

	if srv.Active == nil {
		return fmt.Errorf("No active window")
	}

	return srv.Active.SetLine(pos, text)
}

func (srv *Server) SetScrollDelay(pos int, delay time.Duration) error {
	srv.Lock()
	defer srv.Unlock()

	if srv.Active == nil {
		return fmt.Errorf("No active window")
	}

	return srv.Active.SetScrollDelay(pos, delay)
}

func (srv *Server) Move(window string, n int) error {
	srv.Lock()
	defer srv.Unlock()

	win, ok := srv.Windows[window]
	if !ok {
		return fmt.Errorf("No such window: %s", window)
	}

	win.Move(n)
	return nil
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

		switch split := strings.SplitN(line, " ", 3); split[0] {
		case "window":
			name := ""
			if _, err := fmt.Sscanf(line, "window %s", &name); err != nil {
				log.Printf("Failed to parse window command `%s`: %v", line, err)
				continue
			}

			srv.AddWindow(name)
		case "switch":
			name := ""
			if _, err := fmt.Sscanf(line, "switch %s", &name); err != nil {
				log.Printf("Failed to parse switch command `%s`: %v", line, err)
				continue
			}

			if err := srv.Switch(name); err != nil {
				log.Printf("Unable to switch window: %s", name)
				continue
			}
		case "line":
			text := ""
			if len(split) >= 3 {
				text = split[2]
			}

			pos := 0
			if _, err := fmt.Sscanf(line, "line %d ", &pos); err != nil {
				log.Printf("Failed to parse line command `%s`: %v", line, err)
				continue
			}

			if err := srv.SetLine(pos, text); err != nil {
				log.Printf("Failed to set line: %v", err)
				continue
			}
		case "scroll":
			pos, durationSpec := 0, ""
			if _, err := fmt.Sscanf(line, "scroll %d %s", &pos, &durationSpec); err != nil {
				log.Printf("Failed to parse scroll command `%s`: %v", line, err)
				continue
			}

			duration, err := time.ParseDuration(durationSpec)
			if err != nil {
				log.Printf("Bad duration `%s`: %v", durationSpec, err)
				continue
			}

			if err := srv.SetScrollDelay(pos, duration); err != nil {
				log.Printf("Cannot set scroll: %v", err)
				continue
			}
		case "move":
			name, pos := "", 0
			if _, err := fmt.Sscanf(line, "move %s %d", &name, &pos); err != nil {
				log.Printf("Failed to parse move command `%s`: %v", line, err)
				continue
			}

			if err := srv.Move(name, pos); err != nil {
				log.Printf("Failed to execute move command `%s`: %v", line, err)
				continue
			}
		case "render":
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

func RunDaemon(cfg *Config) error {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	lsn, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	log.Printf("Listening on %s", addr)

	defer lsn.Close()

	// sigint := make(chan os.Signal)
	// signal.Notify(sigint, os.Interrupt)

	srv, err := NewServer(cfg)
	if err != nil {
		return err
	}

	for {
		// Check if we were interrupted:
		select {
		case <-srv.Quit:
			return nil
		// case <-sigint:
		// 	return nil
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

func createClient(cfg *Config, window string) (net.Conn, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	cmd := fmt.Sprintf("window %s\nswitch %s\n", window, window)
	if _, err := conn.Write([]byte(cmd)); err != nil {
		return nil, err
	}

	return conn, nil
}

func RunDumpClient(cfg *Config, window string, update bool) error {
	conn, err := createClient(cfg, window)
	if err != nil {
		return err
	}

	for {
		if _, err := conn.Write([]byte("render\n")); err != nil {
			return err
		}

		if update {
			// Clear the screen:
			fmt.Println("\033[H\033[2J")
		}

		n := int64(cfg.Width*cfg.Height + cfg.Height)
		if _, err := io.CopyN(os.Stdout, conn, n); err != nil {
			return err
		}

		if update {
			time.Sleep(50 * time.Millisecond)
		} else {
			break
		}
	}

	return nil
}

func RunInputClient(cfg *Config, quit bool, window string) error {
	conn, err := createClient(cfg, window)
	if err != nil {
		return err
	}

	if quit {
		_, err := conn.Write([]byte("quit\n"))
		return err
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if _, err := conn.Write([]byte(scanner.Text())); err != nil {
			return err
		}
	}

	return scanner.Err()
}