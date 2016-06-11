package lightd

import (
	"fmt"
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

func Lock(cfg *Config) error {
	return Send(cfg, "!lock")
}

func Unlock(cfg *Config) error {
	return Send(cfg, "!unlock")
}
