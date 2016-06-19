package display

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/net/context"
)

type LineWriter struct {
	host string
	port int

	conn   net.Conn
	ctx    context.Context
	cancel context.CancelFunc
}

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

func (lw *LineWriter) Printf(format string, args ...interface{}) (int, error) {
	return lw.Write([]byte(fmt.Sprintf(format, args...)))
}

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

func RunDumpClient(cfg *Config, ctx context.Context, window string, update bool) error {
	lw, err := Connect(cfg, ctx)
	if err != nil {
		return err
	}

	if _, err := lw.Printf("switch %s", window); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		if _, err := lw.Printf("render"); err != nil {
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

func RunInputClient(cfg *Config, ctx context.Context, quit bool, window string) error {
	lw, err := Connect(cfg, ctx)
	if err != nil {
		return err
	}

	if _, err := lw.Printf("switch %s", window); err != nil {
		return err
	}

	if quit {
		_, err := lw.Printf("quit")
		return err
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		if _, err := lw.Printf(scanner.Text()); err != nil {
			return err
		}
	}

	return scanner.Err()
}
