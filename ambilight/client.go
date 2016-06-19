package ambilight

import (
	"fmt"
	"log"
	"net"

	"github.com/studentkittens/eulenfunk/util"
)

// Client connects to a running ambilightd instance
// and can enable/disable the led playback, check the state
// or quit the daemon remotely.
type Client struct {
	conn net.Conn
}

// NewClient creates a new Client from the cfg.Host and cfg.Port
func NewClient(cfg *Config) (*Client, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &Client{conn: conn}, nil
}

func (cl *Client) send(s string) error {
	if _, err := cl.conn.Write([]byte(s + "\n")); err != nil {
		log.Printf("Failed to send command to ambilightd: %v", err)
		return err
	}

	return nil
}

// Enabled checks if ambilightd is currently doing playingback.
// First return is always false when an error occured.
func (cl *Client) Enabled() (bool, error) {
	if err := cl.send("state"); err != nil {
		return false, err
	}

	resp := make([]byte, 2)
	if _, err := cl.conn.Read(resp); err != nil {
		return false, err
	}

	return string(resp) == "1\n", nil
}

// Enable enables or disables the playback of ambilight.
func (cl *Client) Enable(enable bool) error {
	if enable {
		return cl.send("on")
	}

	return cl.send("off")
}

// Quit attempts to shut down the daemon.
func (cl *Client) Quit() error {
	return cl.send("quit")
}

// Close properly terminates the connection and closes resources.
func (cl *Client) Close() error {
	defer util.Closer(cl.conn)
	return cl.send("close")
}

// WithClient is a convinience function to execute a code snippet with
// an ambilight connection.
func WithClient(host string, port int, fn func(client *Client) error) error {
	client, err := NewClient(&Config{
		Host: host,
		Port: port,
	})

	if err != nil {
		return err
	}

	defer util.Closer(client)
	return fn(client)
}
