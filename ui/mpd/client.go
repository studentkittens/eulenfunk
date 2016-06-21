package mpd

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/fhs/gompd/mpd"
	"github.com/studentkittens/eulenfunk/display"
	"github.com/studentkittens/eulenfunk/util"
	"golang.org/x/net/context"
)

const (
	// PlaybackStop means that there is no current song.
	PlaybackStop = "stop"
	// PlaybackPlay means that there is a current song in progress.
	PlaybackPlay = "play"
	// PlaybackPause means that there is a current song, but no progress.
	PlaybackPause = "pause"
)

// Config defines the connection details the mpd client will use.
type Config struct {
	MPDHost     string
	MPDPort     int
	DisplayHost string
	DisplayPort int
}

// Client is a utility mpd client tailored for the ui's purposes.
// It draws it's status onto the window `mpd`.
type Client struct {
	sync.Mutex

	Config    *Config
	MPD       *ReMPD
	LW        *display.LineWriter
	Status    mpd.Attrs
	CurrSong  mpd.Attrs
	Playlists []string
	Callbacks map[string][]func()

	ctx    context.Context
	cancel context.CancelFunc
}

func displayInfo(lw *display.LineWriter, block []string) error {
	for idx, line := range block {
		if _, err := lw.Printf("line mpd %d %s", idx, line); err != nil {
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
		if _, err := lw.Printf("line stats %d %s", idx, line); err != nil {
			log.Printf("Failed to send line to display server: %v", err)
			return err
		}
	}

	return nil
}

func formatStop(status mpd.Attrs) ([]string, error) {
	return []string{
		"        ⏹          ",
		"Playback stopped.  ",
		"Turn knob to start ",
		"        ⏹          ",
	}, nil
}

func isRadio(currSong mpd.Attrs) bool {
	_, ok := currSong["Name"]
	return ok
}

func displayFormatted(lw *display.LineWriter, currSong, status mpd.Attrs) error {
	var block []string
	var err error

	if status["state"] == PlaybackStop {
		block, err = formatStop(status)
	} else if isRadio(currSong) {
		block, err = formatRadio(currSong, status)
	} else {
		block, err = formatSong(currSong, status)
	}

	if err != nil {
		return err
	}

	if derr := displayInfo(lw, block); derr != nil {
		log.Printf("Failed to display status info: %v", derr)
		return derr
	}

	return nil
}

func formatTimeSpec(tm time.Duration) string {
	h, m, s := int(tm.Hours()), int(tm.Minutes())%60, int(tm.Seconds())%60

	f := fmt.Sprintf("%02d:%02d", m, s)
	if h == 0 {
		return f
	}

	return fmt.Sprintf("%02d:", h) + f
}

// StateToUnicode converts `state` into a nicer unicode glyph.
func StateToUnicode(state string) string {
	switch state {
	case PlaybackPlay:
		return "▶"
	case PlaybackPause:
		return "⏸"
	case PlaybackStop:
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

	length := formatTimeSpec(time.Duration(elapsedSec*1000) * time.Millisecond)

	// Append total time if available:
	if timeStr, ok := currSong["Time"]; ok {
		if totalSec, err := strconv.Atoi(timeStr); err == nil {
			length += "/" + formatTimeSpec(time.Duration(totalSec)*time.Second)
		}
	} else {
		// Pad the elapsed time to the right if no total time available:
		length = "   " + length
	}

	bitrateStr := ""

	if len(status["bitrate"]) > 0 {
		bitrate, err := strconv.Atoi(status["bitrate"])
		if err == nil {
			if bitrate > 999 {
				bitrate = 999
			}

			bitrateStr = fmt.Sprintf("%-3dKBs", bitrate)
		}
	}

	return fmt.Sprintf("%s %s %s", state, bitrateStr, length)
}

func formatRadio(currSong, status mpd.Attrs) ([]string, error) {
	block := []string{
		currSong["Title"],
		fmt.Sprintf("Radio: %s", currSong["Name"]),
		"",
		formatStatusLine(currSong, status),
	}

	return block, nil
}

func formatSong(currSong, status mpd.Attrs) ([]string, error) {
	genre := currSong["Genre"]
	if len(genre) > 0 {
		genre = " (" + genre + ")"
	}

	pos, err := strconv.Atoi(currSong["Pos"])
	if err != nil {
		return nil, err
	}

	block := []string{
		fmt.Sprintf("%s", currSong["Artist"]),
		fmt.Sprintf("%s (#%d)", currSong["Title"], pos+1),
		fmt.Sprintf("%s%s", currSong["Album"], genre),
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

// ListPlaylists returns all stored playlist names.
func (cl *Client) ListPlaylists() []string {
	cl.Lock()
	defer cl.Unlock()

	// Copy slice since it might be modified by updatePlaylists:
	n := make([]string, len(cl.Playlists))
	copy(n, cl.Playlists)
	return n
}

// LoadAndPlayPlaylist subsitutes the queue with the stored playlist `name`
// and immediately starts playing it's first song.
func (cl *Client) LoadAndPlayPlaylist(name string) error {
	cl.Lock()
	defer cl.Unlock()

	if err := cl.MPD.Client().Clear(); err != nil {
		return err
	}

	if err := cl.MPD.Client().PlaylistLoad(name, -1, -1); err != nil {
		return err
	}

	return cl.MPD.Client().Play(0)
}

// TogglePlayback toggles between pause and play.
// If the state is stop, the first queued song is played.
func (cl *Client) TogglePlayback() error {
	cl.Lock()
	defer cl.Unlock()

	mpd := cl.MPD.Client()
	var err error

	switch cl.Status["state"] {
	case PlaybackPlay:
		err = mpd.Pause(true)
	case PlaybackPause:
		err = mpd.Pause(false)
	case PlaybackStop:
		err = mpd.Play(0)
	}

	return err
}

// CurrentState returns the current state ("play", "pause" or "stop")
func (cl *Client) CurrentState() string {
	cl.Lock()
	defer cl.Unlock()

	return cl.Status["state"]
}

// IsRandom returns true when the playback is randomized.
func (cl *Client) IsRandom() bool {
	cl.Lock()
	defer cl.Unlock()

	return cl.Status["random"] == "1"
}

// EnableRandom sets the random state to `enable`.
func (cl *Client) EnableRandom(enable bool) error {
	cl.Lock()
	defer cl.Unlock()

	return cl.MPD.Client().Random(enable)
}

// Next skips to the next song in the queue.
func (cl *Client) Next() error {
	cl.Lock()
	defer cl.Unlock()

	return cl.MPD.Client().Next()
}

// Prev goes to the previous song in the queue.
func (cl *Client) Prev() error {
	cl.Lock()
	defer cl.Unlock()

	return cl.MPD.Client().Previous()
}

// Play unpauses playback or plays the first song if stopped.
func (cl *Client) Play() error {
	cl.Lock()
	defer cl.Unlock()

	if cl.Status["state"] == PlaybackStop {
		return cl.MPD.Client().Play(0)
	}

	return cl.MPD.Client().Pause(false)
}

// Pause sets the playing state to "pause" or leaves "stop".
func (cl *Client) Pause() error {
	cl.Lock()
	defer cl.Unlock()

	return cl.MPD.Client().Pause(true)
}

// Stop sets the state to "stop"
func (cl *Client) Stop() error {
	cl.Lock()
	defer cl.Unlock()

	return cl.MPD.Client().Stop()
}

// Outputs returns a list of outputnames.
// The index of the names are their ids.
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

// ActiveOutput returns the currently selected output.
//
// NOTE: MPD supports more than one active, but our software ignroes that.
//      (German software is excellent at ignoring reality.)
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

// SwitchToOutput enables the output named bt `enableMe`.
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

// NewClient returns a new mpd client that offers a few incomplete, convinience methods
// for altering MPD's state. It also renders the current state to the "mpd" window.
func NewClient(cfg *Config, ctx context.Context) (*Client, error) {
	subCtx, cancel := context.WithCancel(ctx)

	MPD := NewReMPD(cfg.MPDHost, cfg.MPDPort, subCtx)
	lw, err := display.Connect(&display.Config{
		Host: cfg.DisplayHost,
		Port: cfg.DisplayPort,
	}, subCtx)

	if err != nil {
		return nil, err
	}

	if _, err := lw.Printf("switch mpd"); err != nil {
		log.Printf("Failed to send initial switch to display server: %v", err)
	}

	// Make the first 3 lines scrolling:
	for idx := 0; idx < 3; idx++ {
		if _, err := lw.Printf("scroll mpd %d 400ms", idx); err != nil {
			log.Printf("Failed to set scroll: %v", err)
		}
	}

	return &Client{
		Config:    cfg,
		MPD:       MPD,
		LW:        lw,
		Callbacks: make(map[string][]func()),
		ctx:       subCtx,
		cancel:    cancel,
	}, nil
}

// Register remembers a function that will be called the idle event `signal` is received.
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

// Close cancels all client operations
func (cl *Client) Close() error {
	cl.cancel()
	return nil
}

func pinger(ctx context.Context, MPD *ReMPD) {
	ticker := time.NewTicker(1 * time.Minute)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := MPD.Client().Ping(); err != nil {
				log.Printf("ping to MPD failed. Welp. Reason: %v", err)
			}
		}
	}
}

func periodicUpdate(ctx context.Context, MPD *ReMPD, updateCh chan<- string) {
	// Do an initial update:
	updateCh <- "player"
	updateCh <- "stored_playlist"
	updateCh <- "stats"
	updateCh <- "output"

	lo := time.NewTicker(1 * time.Second)
	hi := time.NewTicker(time.Minute)

	for {
		select {
		case <-ctx.Done():
			return
		case <-lo.C:
			updateCh <- "player"
		case <-hi.C:
			updateCh <- "stats"
		}
	}
}

func eventWatcher(ctx context.Context, watcher *ReWatcher, updateCh chan<- string) {
	defer util.Closer(watcher)

	for {
		select {
		case <-ctx.Done():
			return
		case ev := <-watcher.Events:
			updateCh <- ev
		}
	}
}

func (cl *Client) handleUpdate(ev string) {
	switch ev {
	case "stored_playlist":
		if err := cl.updatePlaylists(); err != nil {
			log.Printf("Failed to update playlists: %v", err)
			return
		}
	case "stats":
		stats, err := cl.MPD.Client().Stats()
		if err != nil {
			log.Printf("Failed to fetch statistics: %v", err)
			return
		}

		if err := displayStats(cl.LW, stats); err != nil {
			log.Printf("Failed to display playlists: %v", err)
			return
		}
	case "player", "options":
		song, err := cl.MPD.Client().CurrentSong()
		if err != nil {
			log.Printf("Unable to fetch current song: %v", err)
			return
		}

		status, err := cl.MPD.Client().Status()
		if err != nil {
			log.Printf("Unable to fetch status: %v", err)
			return
		}

		cl.Lock()
		cl.Status = status
		cl.CurrSong = song
		cl.Unlock()

		if err := displayFormatted(cl.LW, song, status); err != nil {
			log.Printf("Failed to format current status: %v", err)
			return
		}
	}

	// Notify observers for all events:
	cl.emit(ev)
}

// Run starts the client operations by keeping the status up-to-date
// and drawing it on the `mpd` window.
func (cl *Client) Run() {
	// Make sure the mpd connection survives long timeouts:
	go pinger(cl.ctx, cl.MPD)

	updateCh := make(chan string)

	// sync extra every few seconds:
	go periodicUpdate(cl.ctx, cl.MPD, updateCh)

	// Also sync on every mpd event:
	watcher := NewReWatcher(cl.Config.MPDHost, cl.Config.MPDPort, cl.ctx)
	go eventWatcher(cl.ctx, watcher, updateCh)

	for {
		select {
		case <-cl.ctx.Done():
			return
		case ev := <-updateCh:
			cl.handleUpdate(ev)
		}
	}
}
