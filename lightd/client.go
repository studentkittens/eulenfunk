package lightd

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/studentkittens/eulenfunk/util"
)

// Send one or more effects to lightd.
func Send(cfg *Config, effects ...string) error {
	if len(effects) == 0 {
		return nil
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		log.Printf("Unable to connect to `lightd`: %v", err)
		return err
	}

	defer util.Closer(conn)

	for _, effectSpec := range effects {
		effectLine := effectSpec + "\n"
		if _, err := conn.Write([]byte(effectLine)); err != nil {
			log.Printf("Cannot send effect `%s`: %v", effectSpec, err)
			return err
		}
	}

	return nil
}

// Locker is a utility to hold a lock on the LED resource
type Locker struct {
	conn net.Conn
}

// NewLocker will create a new Locker connected to the lightd at `cfg.Host` and
// `cfg.Port`.
func NewLocker(cfg *Config) (*Locker, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		log.Printf("Unable to connect to `lightd`: %v", err)
		return nil, err
	}

	return &Locker{conn}, nil
}

// Lock will give you exclusive access to the LED or waits until it can be locked.
func (lk *Locker) Lock() error {
	return lk.send("!lock\n")
}

// Unlock returns the exclusive access right to the next waiting or lightd.
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

// Close will close the connection used by Locker
// NOTE: No Unlock is done!
func (lk *Locker) Close() error {
	return lk.conn.Close()
}
