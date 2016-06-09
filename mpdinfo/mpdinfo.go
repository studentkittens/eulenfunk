package mpdinfo

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/fhs/gompd/mpd"
)

type Config struct {
	Host        string
	Port        int
	DisplayHost string
	DisplayPort int
}

func display(conn net.Conn, textCh chan []string) {
	if _, err := conn.Write([]byte("window mpd\nswitch mpd\n")); err != nil {
		log.Printf("Failed to send hello to display server: %v", err)
		return
	}

	// Make the first 3 lines scrolling:
	for idx := 0; idx < 3; idx++ {
		if _, err := conn.Write([]byte(fmt.Sprintf("scroll %d 200ms\n", idx))); err != nil {
			return
		}
	}

	for block := range textCh {
		for idx, line := range block {
			if _, err := conn.Write([]byte(fmt.Sprintf("line %d %s\n", idx, line))); err != nil {
				log.Printf("Failed to send line to display server: %v", err)
			}
		}
	}
}

func format(currSong, status mpd.Attrs) ([]string, error) {
	block := []string{
		currSong["Artist"],
		fmt.Sprintf("%s (Genre: %s)", currSong["Album"], currSong["Genre"]),
		fmt.Sprintf("%s %s", currSong["Title"], currSong["Track"]),
		status["state"],
	}

	return block, nil
}

func Run(cfg *Config) error {
	client, err := mpd.Dial("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		log.Printf("Failed to connect to mpd: %v", err)
		return err
	}

	dispAddr := fmt.Sprintf("%s:%d", cfg.DisplayHost, cfg.DisplayPort)
	dispConn, err := net.Dial("tcp", dispAddr)
	if err != nil {
		return err
	}

	textCh := make(chan []string)

	go display(dispConn, textCh)

	keepAlivePinger := make(chan bool)

	// Make sure the mpd connection survives long timeouts:
	go func() {
		for range keepAlivePinger {
			client.Ping()
			time.Sleep(1 * time.Minute)
		}
	}()

	updateCh := make(chan bool)

	// sync extra every few seconds:
	go func() {
		for range time.NewTicker(1 * time.Second).C {
			updateCh <- true
		}
	}()

	// Also sync on every mpd event:
	go func() {
		w, err := mpd.NewWatcher(
			"tcp",
			fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			"",
			"player",
		)

		if err != nil {
			log.Printf("Failed to create watcher: %v", err)
			return
		}

		defer w.Close()

		for range w.Event {
			updateCh <- true
		}
	}()

	for range updateCh {
		song, err := client.CurrentSong()
		if err != nil {
			log.Printf("Unable to fetch current song: %v", err)
			continue
		}

		status, err := client.Status()
		if err != nil {
			log.Printf("Unable to fetch status: %v", err)
			continue
		}

		block, err := format(song, status)
		if err != nil {
			log.Printf("Failed to format current status: %v", err)
			continue
		}

		textCh <- block
	}

	return nil
}
