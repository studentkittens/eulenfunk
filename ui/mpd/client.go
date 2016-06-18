package mpd

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/fhs/gompd/mpd"
	"github.com/studentkittens/eulenfunk/display"
	"golang.org/x/net/context"
)

type Config struct {
	Host        string
	Port        int
	DisplayHost string
	DisplayPort int
}

type Client struct {
	sync.Mutex

	Config    *Config
	MPD       *ReMPD
	LW        *display.LineWriter
	Status    mpd.Attrs
	CurrSong  mpd.Attrs
	Playlists []string
	Callbacks map[string][]func()
}

func displayInfo(lw *display.LineWriter, block []string) error {
	for idx, line := range block {
		if _, err := lw.Formatf("line mpd %d %s", idx, line); err != nil {
			log.Printf("Failed to send line to display server: %v", err)
			return err
		}
	}

	return nil
}

func displayStats(lw *display.LineWriter, stats mpd.Attrs) error {
	dbPlaytimeSecs, err := strconv.Atoi(stats["db_playtime"])
	if err != nil {
		return err
	}

	dbPlaytimeDays := float64(dbPlaytimeSecs) / (60 * 60 * 24)

	block := []string{
		fmt.Sprintf("%8s: %s", "Artists", stats["artists"]),
		fmt.Sprintf("%8s: %s", "Albums", stats["albums"]),
		fmt.Sprintf("%8s: %s", "Songs", stats["songs"]),
		fmt.Sprintf("%8s: %.2f days", "Playtime", dbPlaytimeDays),
	}

	for idx, line := range block {
		if _, err := lw.Formatf("line stats %d %s", idx, line); err != nil {
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

func StateToUnicode(state string) string {
	switch state {
	case "play":
		return "▶"
	case "pause":
		return "⏸"
	case "stop":
		return "⏹"
	default:
		return "?"
	}
}

func formatStatusLine(currSong, status mpd.Attrs) string {
	state := StateToUnicode(status["state"])
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

func (cl *Client) updatePlaylists() error {
	cl.Lock()
	defer cl.Unlock()

	spl, err := cl.MPD.Client().ListPlaylists()
	if err != nil {
		return err
	}

	cl.Playlists = nil

	for _, playlist := range spl {
		cl.Playlists = append(cl.Playlists, playlist["playlist"])
	}

	return nil
}

func (cl *Client) ListPlaylists() []string {
	cl.Lock()
	defer cl.Unlock()

	// Copy slice since it might be modified by updatePlaylists:
	n := make([]string, len(cl.Playlists))
	copy(n, cl.Playlists)
	return n
}

func (cl *Client) TogglePlayback() error {
	cl.Lock()
	defer cl.Unlock()

	mpd := cl.MPD.Client()
	var err error

	switch cl.Status["state"] {
	case "play":
		err = mpd.Pause(true)
	case "pause":
		err = mpd.Pause(false)
	case "stop":
		err = mpd.Play(0)
	}

	return err
}

func (cl *Client) CurrentState() string {
	cl.Lock()
	defer cl.Unlock()

	return cl.Status["state"]
}

func (cl *Client) IsRandom() bool {
	cl.Lock()
	defer cl.Unlock()

	return cl.Status["random"] == "1"
}

func (cl *Client) EnableRandom(enable bool) error {
	cl.Lock()
	defer cl.Unlock()

	return cl.MPD.Client().Random(enable)
}

func (cl *Client) Next() error {
	cl.Lock()
	defer cl.Unlock()

	return cl.MPD.Client().Next()
}

func (cl *Client) Prev() error {
	cl.Lock()
	defer cl.Unlock()

	return cl.MPD.Client().Previous()
}

func (cl *Client) Play() error {
	cl.Lock()
	defer cl.Unlock()

	return cl.MPD.Client().Pause(false)
}

func (cl *Client) Pause() error {
	cl.Lock()
	defer cl.Unlock()

	return cl.MPD.Client().Pause(true)
}

func (cl *Client) Stop() error {
	cl.Lock()
	defer cl.Unlock()

	return cl.MPD.Client().Stop()
}

func (cl *Client) Outputs() ([]string, error) {
	cl.Lock()
	defer cl.Unlock()

	outputs, err := cl.MPD.Client().ListOutputs()
	if err != nil {
		return nil, err
	}

	names := []string{}

	for _, output := range outputs {
		names = append(names, output["outputname"])
	}

	return names, nil
}

// NOTE: MPD supports more than one active, but our software does not.
func (cl *Client) ActiveOutput() (string, error) {
	cl.Lock()
	defer cl.Unlock()

	outputs, err := cl.MPD.Client().ListOutputs()
	if err != nil {
		return "", err
	}

	for _, output := range outputs {
		if output["outputenabled"] == "1" {
			return output["outputname"], nil
		}
	}

	return "", nil
}

func (cl *Client) SwitchToOutput(enableMe string) error {
	// Disable all other outputs on the way:
	// (one output is enough for our usecase)
	names, err := cl.Outputs()
	if err != nil {
		return err
	}

	cl.Lock()
	defer cl.Unlock()

	for id, name := range names {
		var err error

		if name == enableMe {
			err = cl.MPD.Client().EnableOutput(id)
		} else {
			err = cl.MPD.Client().DisableOutput(id)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func NewClient(cfg *Config) (*Client, error) {
	MPD := NewReMPD(cfg.Host, cfg.Port)
	lw, err := display.Connect(&display.Config{
		Host: cfg.DisplayHost,
		Port: cfg.DisplayPort,
	})

	if err != nil {
		return nil, err
	}

	if _, err := lw.Formatf("switch mpd"); err != nil {
		log.Printf("Failed to send initial switch to display server: %v", err)
	}

	// Make the first 3 lines scrolling:
	for idx := 0; idx < 3; idx++ {
		if _, err := lw.Formatf("scroll mpd %d 400ms", idx); err != nil {
			log.Printf("Failed to set scroll: %v", err)
		}
	}

	return &Client{
		Config:    cfg,
		MPD:       MPD,
		LW:        lw,
		Callbacks: make(map[string][]func()),
	}, nil
}

func (cl *Client) Register(signal string, action func()) {
	cl.Lock()
	defer cl.Unlock()

	cl.Callbacks[signal] = append(cl.Callbacks[signal], action)
}

func (cl *Client) emit(signal string) {
	cl.Lock()
	actions, ok := cl.Callbacks[signal]
	cl.Unlock()

	if !ok {
		return
	}

	for _, action := range actions {
		action()
	}
}

func (cl *Client) Run(ctx context.Context) {
	// Make sure the mpd connection survives long timeouts:
	go func() {
		ticker := time.NewTicker(1 * time.Minute)

		select {
		case <-ctx.Done():
			break
		case <-ticker.C:
			cl.MPD.Client().Ping()
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

	// Update the stats periodically by faking
	// a "stats" event (not a real event):
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
		watcher := NewReWatcher(
			cl.Config.Host, cl.Config.Port,
			"player", "stored_playlist",
		)

		defer watcher.Close()

		for {
			select {
			case <-ctx.Done():
				return
			case ev := <-watcher.Events:
				updateCh <- ev
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case ev := <-updateCh:
			switch ev {
			case "stored_playlist":
				if err := cl.updatePlaylists(); err != nil {
					log.Printf("Failed to update playlists: %v", err)
					continue
				}
			case "stats":
				stats, err := cl.MPD.Client().Stats()
				if err != nil {
					log.Printf("Failed to fetch statistics: %v", err)
					continue
				}

				if err := displayStats(cl.LW, stats); err != nil {
					log.Printf("Failed to display playlists: %v", err)
					continue
				}
			case "player":
				song, err := cl.MPD.Client().CurrentSong()
				if err != nil {
					log.Printf("Unable to fetch current song: %v", err)
					continue
				}

				status, err := cl.MPD.Client().Status()
				if err != nil {
					log.Printf("Unable to fetch status: %v", err)
					continue
				}

				cl.Lock()
				cl.Status = status
				cl.CurrSong = song
				cl.Unlock()

				block, err := format(song, status)
				if err != nil {
					log.Printf("Failed to format current status: %v", err)
					continue
				}

				if err := displayInfo(cl.LW, block); err != nil {
					log.Printf("Failed to display status info: %v", err)
					continue
				}
			}

			// Notify observers for all events:
			cl.emit(ev)
		}
	}
}
