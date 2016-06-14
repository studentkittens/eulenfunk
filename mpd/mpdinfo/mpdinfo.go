package mpdinfo

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/fhs/gompd/mpd"
	"golang.org/x/net/context"
)

type Config struct {
	Host        string
	Port        int
	DisplayHost string
	DisplayPort int
}

func displayInfo(conn net.Conn, block []string) error {
	for idx, line := range block {
		if _, err := conn.Write([]byte(fmt.Sprintf("line mpd %d %s\n", idx, line))); err != nil {
			log.Printf("Failed to send line to display server: %v", err)
			return err
		}
	}

	return nil
}

func displayPlaylists(conn net.Conn, playlists []mpd.Attrs) error {
	if _, err := conn.Write([]byte("truncate playlists 0")); err != nil {
		return err
	}

	for idx, playlist := range playlists {
		line := fmt.Sprintf("line playlists %d %02d %s", idx, idx+1, playlist["playlist"])
		if _, err := conn.Write([]byte(line)); err != nil {
			return err
		}
	}

	return nil
}

func displayStats(conn net.Conn, stats mpd.Attrs) error {
	dbPlaytimeSecs, err := strconv.Atoi(stats["db_playtime"])
	if err != nil {
		return err
	}

	dbPlaytimeDays := float64(dbPlaytimeSecs) / (60 * 60 * 24)

	block := []string{
		fmt.Sprintf("%8s: %s", "Artists", stats["artist"]),
		fmt.Sprintf("%8s: %s", "Albums", stats["albums"]),
		fmt.Sprintf("%8s: %s", "Songs", stats["songs"]),
		fmt.Sprintf("%8s: %.2f", "Playtime", dbPlaytimeDays),
	}

	for idx, line := range block {
		if _, err := conn.Write([]byte(fmt.Sprintf("line stats %d %s\n", idx, line))); err != nil {
			log.Printf("Failed to send line to display server: %v", err)
			return err
		}
	}

	return nil
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

func commandHandler(client *mpd.Client, cmdCh <-chan string) {
	for cmd := range cmdCh {
		var err error

		switch cmd {
		case "next":
			err = client.Next()
		case "prev":
			err = client.Previous()
		case "play":
			err = client.Pause(false)
		case "pause":
			err = client.Pause(true)
		case "toggle":
			status, err := client.Status()
			if err != nil {
				log.Printf("Failed to fetch status: %v", err)
				continue
			}

			switch status["state"] {
			case "play":
				err = client.Pause(true)
			case "pause":
				err = client.Pause(false)
			case "stop":
				err = client.Play(0)
			}
		case "stop":
			err = client.Stop()
		}

		if err != nil {
			log.Printf("Executung `%s` failed: %v", cmd, err)
		}
	}
}

func Run(cfg *Config, ctx context.Context, cmdCh chan string) error {
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

	if _, err := dispConn.Write([]byte("switch mpd\n")); err != nil {
		log.Printf("Failed to send hello to display server: %v", err)
	}

	// Make the first 3 lines scrolling:
	for idx := 0; idx < 3; idx++ {
		cmd := fmt.Sprintf("scroll mpd %d 400ms\n", idx)
		if _, err := dispConn.Write([]byte(cmd)); err != nil {
			log.Printf("Failed to set scroll: %v", err)
		}
	}

	go commandHandler(client, cmdCh)

	// Make sure the mpd connection survives long timeouts:
	go func() {
		ticker := time.NewTicker(1 * time.Minute)

		select {
		case <-ctx.Done():
			break
		case <-ticker.C:
			client.Ping()
		}
	}()

	updateCh := make(chan string)

	// sync extra every few seconds:
	go func() {
		// Do an initial update:
		updateCh <- "player"
		updateCh <- "stored_playlist"

		ticker := time.NewTicker(1 * time.Second)

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				updateCh <- "player"
			}
		}
	}()

	go func() {
		updateCh <- "stats"

		ticker := time.NewTicker(time.Minute)

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				updateCh <- "stats"
			}
		}
	}()

	// Also sync on every mpd event:
	go func() {
		w, err := mpd.NewWatcher(
			"tcp",
			fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			"",
			"player",
			"stored_playlist",
		)

		if err != nil {
			log.Printf("Failed to create watcher: %v", err)
			return
		}

		defer w.Close()

		for {
			select {
			case <-ctx.Done():
				return
			case ev := <-w.Event:
				updateCh <- ev
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case ev := <-updateCh:
			switch ev {
			case "stored_playlist":
				spl, err := client.ListPlaylists()
				if err != nil {
					log.Printf("Failed to list stored playlists: %v", err)
					continue
				}

				if err := displayPlaylists(dispConn, spl); err != nil {
					log.Printf("Failed to display playlists: %v", err)
					continue
				}
			case "stats":
				stats, err := client.Stats()
				if err != nil {
					log.Printf("Failed to fetch statistics: %v", err)
					continue
				}

				if err := displayStats(dispConn, stats); err != nil {
					log.Printf("Failed to display playlists: %v", err)
					continue
				}
			case "player":
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

				if err := displayInfo(dispConn, block); err != nil {
					log.Printf("Failed to display status info: %v", err)
					continue
				}
			}
		}
	}

	return nil
}
