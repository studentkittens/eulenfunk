package ambilight

import (
	"fmt"
	"log"
	"net"
)

type Client struct {
	conn net.Conn
}

func NewClient(cfg *Config) (*Client, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
	}, nil
}

func (cl *Client) send(s string) error {
	if _, err := cl.conn.Write([]byte(s + "\n")); err != nil {
		log.Printf("Failed to send command to ambilightd: %v", err)
		return err
	}

	return nil
}

func (cl *Client) Enabled() (bool, error) {
	if err := cl.send("quit"); err != nil {
		return false, err
	}

	resp := make([]byte, 2)
	if _, err := cl.conn.Read(resp); err != nil {
		return false, err
	}

	return string(resp) == "1\n", nil
}

func (cl *Client) Enable(enable bool) error {
	if enable {
		return cl.send("on")
	} else {
		return cl.send("off")
	}
}

func (cl *Client) Quit() error {
	return cl.send("quit")
}

func (cl *Client) Close() error {
	return cl.send("close")
}
