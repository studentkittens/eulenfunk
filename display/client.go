package display

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

// TODO: find a better name?
type LineWriter struct {
	sync.Mutex

	host string
	port int

	conn net.Conn
	quit bool
}

func (lw *LineWriter) Write(p []byte) (int, error) {
	if !bytes.HasSuffix(p, []byte("\n")) {
		p = append(p, '\n')
	}

	for {
		lw.Lock()
		if lw.quit {
			break
		}
		lw.Unlock()

		n, err := lw.conn.Write(p)
		if err != nil {
			lw.retryUntilSuccesfull()
			continue
		}

		return n, err
	}

	return 0, nil
}

func (lw *LineWriter) Read(p []byte) (int, error) {
	for {
		lw.Lock()
		if lw.quit {
			break
		}

		n, err := lw.conn.Read(p)
		lw.Unlock()

		if err != nil {
			lw.retryUntilSuccesfull()
			continue
		}

		return n, err
	}

	return 0, nil
}

// TODO: Formatf -> Printf
func (lw *LineWriter) Formatf(format string, args ...interface{}) (int, error) {
	return lw.Write([]byte(fmt.Sprintf(format, args...)))
}

func (lw *LineWriter) Close() error {
	lw.Lock()
	defer lw.Unlock()

	lw.quit = true
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
	lw.Lock()
	defer lw.Unlock()

	for !lw.quit {
		if err := lw.reconnect(); err != nil {
			log.Printf("Failed to connect to displayd: %v", err)
			log.Printf("Retry in 5 seconds")
			time.Sleep(5 * time.Second)
			continue
		}

		break
	}
}

func Connect(cfg *Config) (*LineWriter, error) {
	lw := &LineWriter{host: cfg.Host, port: cfg.Port}
	lw.retryUntilSuccesfull()
	return lw, nil
}

func RunDumpClient(cfg *Config, window string, update bool) error {
	lw, err := Connect(cfg)
	if err != nil {
		return err
	}

	if _, err := lw.Formatf("switch %s", window); err != nil {
		return err
	}

	for {
		if _, err := lw.Formatf("render"); err != nil {
			return err
		}

		if update {
			// Clear the screen:
			fmt.Println("\033[H\033[2J")
		}

		n := int64(cfg.Width*cfg.Height + cfg.Height)
		if _, err := io.CopyN(os.Stdout, lw, n); err != nil {
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
	lw, err := Connect(cfg)
	if err != nil {
		return err
	}

	if _, err := lw.Formatf("switch %s", window); err != nil {
		return err
	}

	if quit {
		_, err := lw.Formatf("quit")
		return err
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if _, err := lw.Formatf(scanner.Text()); err != nil {
			return err
		}
	}

	return scanner.Err()
}
