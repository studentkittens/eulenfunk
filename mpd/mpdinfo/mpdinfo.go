package mpdinfo

import (
	"fmt"
	"log"
	"net"
	"strconv"
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
		if _, err := conn.Write([]byte(fmt.Sprintf("scroll %d 400ms\n", idx))); err != nil {
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

func isRadio(currSong mpd.Attrs) bool {
	_, ok := currSong["Name"]
	return ok
}

func format(currSong, status mpd.Attrs) ([]string, error) {
	if isRadio(currSong) {
		return formatRadio(currSong, status)
	}

	return formatSong(currSong, status)
}

func formatTimeSpec(tm time.Duration) string {
	h, m, s := int(tm.Hours()), int(tm.Minutes())%60, int(tm.Seconds())%60

	f := fmt.Sprintf("%02d:%02d", m, s)
	if h == 0 {
		return f
	}

	return fmt.Sprintf("%02d:", h) + f
}

func formatStatusLine(currSong, status mpd.Attrs) string {
	state := "[" + status["state"] + "]"
	elapsedStr := status["elapsed"]

	elapsedSec, err := strconv.ParseFloat(elapsedStr, 64)
	if err != nil {
		return state
	}

	state += " "
	state += formatTimeSpec(time.Duration(elapsedSec*1000) * time.Millisecond)

	// Append total time if available:
	if timeStr, ok := currSong["Time"]; ok {
		if totalSec, err := strconv.Atoi(timeStr); err == nil {
			state += "/" + formatTimeSpec(time.Duration(totalSec)*time.Second)
		}
	}

	return state
}

func formatRadio(currSong, status mpd.Attrs) ([]string, error) {
	block := []string{
		currSong["Title"],
		fmt.Sprintf("Radio: %s", currSong["Name"]),
		fmt.Sprintf("Bitrate: %s Kbit/s", status["bitrate"]),
		formatStatusLine(currSong, status),
	}

	return block, nil
}

func formatSong(currSong, status mpd.Attrs) ([]string, error) {
	block := []string{
		currSong["Artist"],
		fmt.Sprintf("%s (Genre: %s)", currSong["Album"], currSong["Genre"]),
		fmt.Sprintf("%s %s", currSong["Title"], currSong["Track"]),
		formatStatusLine(currSong, status),
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
