package automount

import (
	"fmt"
	"net"

	"github.com/disorganizer/brig/util"
)

// Client is a convinience helper to access the automount text protocol
type Client struct {
	conn net.Conn
}

// NewClient returns a new automountd convinience client.
func NewClient(cfg *Config) (*Client, error) {
	addr := fmt.Sprintf("%s:%d", cfg.AutomountHost, cfg.AutomountPort)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &Client{conn}, nil
}

// Close shuts down the client and frees resources.
func (cl *Client) Close() error {
	return cl.conn.Close()
}

// Mount sends a "mount <device> <label>" message to the automount daemon.
func (cl *Client) Mount(device, label string) error {
	_, err := cl.conn.Write([]byte(fmt.Sprintf("mount %s %s\n", device, label)))
	return err
}

// Unmount sends a "unmount <device>" message to the automountd daemon.
func (cl *Client) Unmount(device string) error {
	_, err := cl.conn.Write([]byte(fmt.Sprintf("unmount %s\n", device)))
	return err
}

// Quit sends a "quit" message to the automountd daemon.
func (cl *Client) Quit() error {
	_, err := cl.conn.Write([]byte("quit\n"))
	return err
}

// WithClient executes fn with a valid client as argument.
// The returned error is passed back to WithClient's caller.
func WithClient(cfg *Config, fn func(cl *Client) error) error {
	client, err := NewClient(cfg)
	if err != nil {
		return err
	}

	defer util.Closer(client)

	return fn(client)
}
