package display

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/net/context"
)

// LineWriter is good at sending single formatted line to displayd.
// It's not an abstraction over the text protocol, but an easy access to it.
// It's methods will block if the displayd server does not exist yet or
// crashed. It then attempts a reconnect in the background.
type LineWriter struct {
	host string
	port int

	conn   net.Conn
	ctx    context.Context
	cancel context.CancelFunc
}

// Write sends arbitrary bytes to displayd. You should use Printf() instead.
func (lw *LineWriter) Write(p []byte) (int, error) {
	if !bytes.HasSuffix(p, []byte("\n")) {
		p = append(p, '\n')
	}

	for {
		select {
		case <-lw.ctx.Done():
			return 0, nil
		default:
		}

		n, err := lw.conn.Write(p)
		if err != nil {
			lw.retryUntilSuccesfull()
			continue
		}

		return n, err
	}
}

// Read reads a response from the display server.
// This is only needed by the debugging dumping program.
func (lw *LineWriter) Read(p []byte) (int, error) {
	for {
		select {
		case <-lw.ctx.Done():
			return 0, nil
		default:
		}

		n, err := lw.conn.Read(p)

		if err != nil {
			lw.retryUntilSuccesfull()
			continue
		}

		return n, err
	}
}

// Printf formats and sends a message to displayd in a fmt.Printf like fashion.
func (lw *LineWriter) Printf(format string, args ...interface{}) (int, error) {
	return lw.Write([]byte(fmt.Sprintf(format, args...)))
}

// Line writes a line in `window` at lineno `pos` consisting of `text`
func (lw *LineWriter) Line(window string, pos int, text string) error {
	_, err := lw.Printf("line %s %d %s", window, pos, text)
	return err
}

// ScrollDelay sets the delay between a scroll increment of the line in the
// window `window` at position `pos` to `delay`.
func (lw *LineWriter) ScrollDelay(window string, pos int, delay time.Duration) error {
	_, err := lw.Printf("scroll %s %d %s", window, pos, delay.String())
	return err
}

// Switch makes `window` the active window.
func (lw *LineWriter) Switch(window string) error {
	_, err := lw.Printf("switch %s", window)
	return err
}

// Move moves the window `window` down by `plus` lines.
// `plus` may be negative to go up again.
// Think of it as vertical scrolling.
func (lw *LineWriter) Move(window string, plus int) error {
	_, err := lw.Printf("move %s %d", window, plus)
	return err
}

// Truncate cuts off the window contents of `window` at the
// absolute offset `cutoff`. Lines above will be cleared.
func (lw *LineWriter) Truncate(window string, cutoff int) error {
	_, err := lw.Printf("truncate %s %d", window, cutoff)
	return err
}

// Quit makes displayd quit.
func (lw *LineWriter) Quit() error {
	_, err := lw.Printf("quit")
	return err
}

// Render returns a display of the current active window.
func (lw *LineWriter) Render() ([]byte, error) {
	buf := &bytes.Buffer{}

	if _, err := lw.Printf("render"); err != nil {
		return nil, err
	}

	sizeBuf := make([]byte, 8)
	n, err := lw.Read(sizeBuf)
	if n != 8 {
		log.Printf("Bad size header (%d != 8)", n)
		return nil, fmt.Errorf("Bad size header")
	}

	if err != nil {
		log.Printf("Reading size header failed: %v", err)
		return nil, err
	}

	size := binary.LittleEndian.Uint64(sizeBuf)
	if _, err := io.CopyN(buf, lw, int64(size)); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Close cancels all pending operations and frees resources.
func (lw *LineWriter) Close() error {
	lw.cancel()
	return lw.conn.Close()
}

func (lw *LineWriter) reconnect() error {
	addr := fmt.Sprintf("%s:%d", lw.host, lw.port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	lw.conn = conn
	return nil
}

func (lw *LineWriter) retryUntilSuccesfull() {
	for {
		select {
		case <-lw.ctx.Done():
			return
		default:
		}

		if err := lw.reconnect(); err != nil {
			log.Printf("Failed to connect to displayd: %v", err)
			log.Printf("Retry in 5 seconds")
			time.Sleep(5 * time.Second)
			continue
		}

		break
	}
}

// Connect creates a LineWriter pointed to the displayd at $(cfg.Host):$(cfg.Port)
func Connect(cfg *Config, ctx context.Context) (*LineWriter, error) {
	subCtx, cancel := context.WithCancel(ctx)

	lw := &LineWriter{
		host:   cfg.Host,
		port:   cfg.Port,
		ctx:    subCtx,
		cancel: cancel,
	}

	lw.retryUntilSuccesfull()
	return lw, nil
}

func cancelled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

// DumpClient is a client that renders the content of `window` onto stdout.
// Optionally it will clear & update the screen if `update` is true.
func DumpClient(cfg *Config, ctx context.Context, window string, update bool) error {
	lw, err := Connect(cfg, ctx)
	if err != nil {
		return err
	}

	if _, err := lw.Printf("switch %s", window); err != nil {
		return err
	}

	for !cancelled(ctx) {
		if update {
			// Clear the screen:
			fmt.Println("\033[H\033[2J")
		}

		matrix, err := lw.Render()
		if err != nil {
			return err
		}

		if _, err := io.Copy(os.Stdout, bytes.NewReader(matrix)); err != nil {
			return err
		}

		if update {
			time.Sleep(50 * time.Millisecond)
			continue
		}

		break
	}

	return nil
}

// InputClient is a dump client that can be used to send arbitrary displayd
// lines in a netcat or telnet like fashion from the commandline.
func InputClient(cfg *Config, ctx context.Context, quit bool, window string) error {
	lw, err := Connect(cfg, ctx)
	if err != nil {
		return err
	}

	if err := lw.Switch(window); err != nil {
		return err
	}

	if quit {
		return lw.Quit()
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		if _, err := lw.Write([]byte(scanner.Text())); err != nil {
			return err
		}
	}

	return scanner.Err()
}
