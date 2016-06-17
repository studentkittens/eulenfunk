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

func createClient(cfg *Config, window string) (net.Conn, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	cmd := fmt.Sprintf("switch %s\n", window)
	if _, err := conn.Write([]byte(cmd)); err != nil {
		return nil, err
	}

	return conn, nil
}

type LineWriter struct {
	sync.Mutex
	conn net.Conn
}

func (lw *LineWriter) Write(p []byte) (int, error) {
	lw.Lock()
	defer lw.Unlock()

	if !bytes.HasSuffix(p, []byte("\n")) {
		p = append(p, '\n')
	}

	log.Printf("lw: %s", p)
	return lw.conn.Write(p)
}

func (lw *LineWriter) Formatf(format string, args ...interface{}) (int, error) {
	return lw.Write([]byte(fmt.Sprintf(format, args...)))
}

func (lw *LineWriter) Close() error {
	lw.Lock()
	defer lw.Unlock()

	return lw.conn.Close()
}

func Connect(cfg *Config) (*LineWriter, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &LineWriter{conn: conn}, nil
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
