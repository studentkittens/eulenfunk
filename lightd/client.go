package lightd

import (
	"fmt"
	"io"
	"log"
	"net"
)

func Send(cfg *Config, effects ...string) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		log.Printf("Unable to connect to `lightd`: %v", err)
		return err
	}

	defer conn.Close()

	for _, effectSpec := range effects {
		effectLine := effectSpec + "\n"
		if _, err := conn.Write([]byte(effectLine)); err != nil {
			log.Printf("Cannot send effect `%s`: %v", effectSpec, err)
			return err
		}
	}

	return nil
}

type Locker struct {
	conn net.Conn
}

func NewLocker(cfg *Config) (*Locker, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		log.Printf("Unable to connect to `lightd`: %v", err)
		return nil, err
	}

	return &Locker{conn}, nil
}

func (lk *Locker) Lock() error {
	return lk.send("!lock\n")
}

func (lk *Locker) Unlock() error {
	return lk.send("!unlock\n")
}

func (lk *Locker) send(msg string) error {
	if _, err := lk.conn.Write([]byte(msg)); err != nil {
		return err
	}

	ok := make([]byte, 3)
	if _, err := lk.conn.Read(ok); err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (lk *Locker) Close() error {
	return lk.conn.Close()
}
