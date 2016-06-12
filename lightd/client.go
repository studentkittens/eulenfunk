package lightd

import (
	"fmt"
	"io"
	"log"
	"net"
)

func send(cfg *Config, wait bool, effects ...string) error {
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

	if wait {
		m := make([]byte, 10)
		if _, err := conn.Read(m); err != nil && err != io.EOF {
			log.Printf("read failed: %v", err)
		}
	}

	return nil
}

func Send(cfg *Config, effects ...string) error {
	return send(cfg, false, effects...)
}

func Lock(cfg *Config) error {
	return send(cfg, true, "!lock")
}

func Unlock(cfg *Config) error {
	return send(cfg, true, "!unlock")
}
