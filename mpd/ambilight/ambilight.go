// moodymusic is a small mpd client that uses the `moodbar` utility
// and a multicolor LED to create an ambient light
package ambilight

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	// External dependencies:
	"github.com/fhs/gompd/mpd"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/studentkittens/eulenfunk/lightd"
)

type Config struct {
	// Host of the mpd server
	Host string

	// Port of the mpd server
	Port int

	// MusicDir is the root path of the mpd database
	MusicDir string

	// MoodDir contains all moodfiles for certain files (if any)
	MoodDir string

	// UpdateMoodDatabase makes the client update the db and exit afterwards.
	UpdateMoodDatabase bool

	// Name of the RGB LED driver binary (`catlight` for my desktop)
	BinaryName string
}

// Current status of the MPD player:
type MPDEvent struct {
	Path        string
	SongChanged bool
	ElapsedMs   float64
	TotalMs     float64
	IsPlaying   bool
	IsStopped   bool
}

// Info needed to render a moodbar:
type MoodInfo struct {
	MusicFile string
	MoodPath  string
}

// RGB Color that stays for a certain duration:
type TimedColor struct {
	R, G, B  uint8
	Duration time.Duration
}

// checkForMoodbar tries to execute `moodbar` and prints a helpful message before exitting otherwise:
func checkForMoodbar() {
	cmd := exec.Command("moodbar", "--help")
	if err := cmd.Run(); err != nil {
		log.Printf("Could not execute `moodbar --help` - is it installed?")
		log.Printf("Error was: %v", err)
		os.Exit(-1)
	}
}

// Walk over all music files and create a .mood file for each in mood-dir.
func updateMoodDatabase(client *mpd.Client, cfg *Config) error {
	if cfg.MoodDir == "" {
		return fmt.Errorf("No mood bar directory given (-mood-dir)")
	}

	if err := os.MkdirAll(cfg.MoodDir, 0777); err != nil {
		return err
	}

	paths, err := client.GetFiles()
	if err != nil {
		return fmt.Errorf("Cannot get all files from mpd: %v", err)
	}

	// Use up to N threads:
	N := 8
	wg := &sync.WaitGroup{}
	wg.Add(N)

	moodChan := make(chan *MoodInfo, N)
	for i := 0; i < N; i++ {
		go func() {
			for pair := range moodChan {
				log.Printf("Processing: %s", pair.MusicFile)
				cmd := exec.Command("moodbar", pair.MusicFile, "-o", pair.MoodPath)
				if err := cmd.Run(); err != nil {
					log.Printf("Failed to execute moodbar on `%s`: %v", pair.MusicFile, err)
				}
			}

			wg.Done()
		}()
	}

	for _, path := range paths {
		moodName := strings.Replace(path, string(filepath.Separator), "|", -1)
		moodPath := filepath.Join(cfg.MoodDir, moodName)

		if _, err := os.Stat(moodPath); err == nil {
			// Already exists, Skipping.
			continue
		}

		dataPath := filepath.Join(cfg.MusicDir, path)

		moodChan <- &MoodInfo{
			MusicFile: dataPath,
			MoodPath:  moodPath,
		}
	}

	close(moodChan)
	wg.Wait()

	return nil
}

// Read the .mood file at `path`; return a 1000-element slice of TimedColor.
// The duration of each color will be 0.
func readMoodbarFile(path string) ([]TimedColor, error) {
	results := []TimedColor{}

	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer fd.Close()

	rgb := make([]byte, 3)
	for {
		_, err := fd.Read(rgb)
		if err != nil && err != io.EOF {
			return nil, err
		}

		results = append(results, TimedColor{
			uint8(rgb[0]), uint8(rgb[1]), uint8(rgb[2]), 0,
		})

		if err == io.EOF {
			break
		}
	}

	return results, nil
}

// Create a HCL Gradient between c1 and c2 using N steps.
// Returns the gradient as slice of individual colors.
func createBlend(c1, c2 TimedColor, N int) []TimedColor {
	// Do nothing if it's the same color:
	if c1.R == c2.R && c1.G == c2.G && c1.B == c2.B {
		return []TimedColor{c1}
	}

	cc1 := colorful.Color{
		float64(c1.R) / 255.,
		float64(c1.G) / 255.,
		float64(c1.B) / 255.,
	}

	cc2 := colorful.Color{
		float64(c2.R) / 255.,
		float64(c2.G) / 255.,
		float64(c2.B) / 255.,
	}

	colors := []TimedColor{}

	for i := 0; i < N; i++ {
		mix := cc1.BlendHcl(cc2, float64(i)/float64(N)).Clamped()
		h, c, l := mix.Hcl()

		// Increase chroma slighly for mid values:
		c = (-((c-1)*(c-1))+1)/2 + c/2

		// Decrease lumninance slighly for mid values:
		l = (l*l)/2 + (l / 2)

		// Convert back to (gamma corrected) RGB for catlight:
		r, g, b := colorful.Hcl(h, c, l).FastLinearRgb()
		colors = append(colors, TimedColor{
			uint8(r * 255), uint8(g * 255), uint8(b * 255),
			((c1.Duration + c2.Duration) / 2) / time.Duration(N),
		})
	}

	return colors
}

// MoodbarRunner sets the current color and blends to it
// by remembering the last color and calculating a gradient between both.
func MoodbarRunner(cfg *Config, colors <-chan TimedColor) {
	cmd := exec.Command(cfg.BinaryName, "cat")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Printf("Cannot fork catlight: %v", err)
		return
	}

	if err := cmd.Start(); err != nil {
		log.Printf("Failed to start `catlight cat`: %v", err)
		return
	}

	defer stdin.Close()

	// First color is always black.
	var lastColor TimedColor

	// Blend between lastColor and current color.
	var blend []TimedColor

	for {
		select {
		case color, ok := <-colors:
			if !ok {
				return
			}

			fmt.Println("recv", color)

			blendInterval := 2
			if color.Duration > 20*time.Millisecond {
				blendInterval = int(math.Sqrt(float64(color.Duration/time.Millisecond))) / 2
			}

			blend = createBlend(lastColor, color, blendInterval)
			lastColor = color
		default:
			if len(blend) > 0 {
				color := blend[0]
				blend = blend[1:]

				colorValue := fmt.Sprintf("%d %d %d\n", color.R, color.G, color.B)
				stdin.Write([]byte(colorValue))
				time.Sleep(color.Duration)
			}
		}
	}
}

// MoodbarAdjuster tried to synchronize the music to the moodbar.
// It will send the correct current color to MoodbarRunner.
func MoodbarAdjuster(eventCh <-chan MPDEvent, colorsCh chan<- TimedColor) {
	var currIdx int
	var colors []TimedColor
	var currEv *MPDEvent

	initialSend := true

	// lightdConfig := &lightd.Config{
	// 	Host: "localhost",
	// 	Port: 3333,
	// }

	sendColor := func(col TimedColor) {
		if err := lightd.Lock(lightdConfig); err != nil {
			log.Printf("Failed to acquire lock (sending anyways): %v", err)
		}

		// Do not crash when colorsCh is closed:
		colorsCh <- col
		time.Sleep(col.Duration)

		if err := lightd.Unlock(lightdConfig); err != nil {
			log.Printf("Failed to unlock: %v", err)
		}
	}

	defer func() {
		close(colorsCh)
	}()

	for {
		select {

		// A new event happened, we need to adjust or even load a new moodbar file:
		case ev, ok := <-eventCh:
			if !ok {
				return
			}

			// Only required if the song changed:
			if ev.SongChanged {
				data, err := readMoodbarFile(ev.Path)
				if err != nil {
					log.Printf("Failed to read moodbar at `%s`: %v", ev.Path, err)

					// Return to black:
					sendColor(TimedColor{0, 0, 0, 0})
					colors = []TimedColor{}
					currIdx = 0
					continue
				}

				colors = data
			}

			// Adjust the moodbar seek offset (1000 samples per total time)
			if ev.TotalMs > 0 {
				currIdx = int((ev.ElapsedMs / ev.TotalMs) * 1000)
			} else {
				// Probably stop or some error:
				currIdx = 0
			}

			currEv = &ev
		default:
			if currIdx >= len(colors) || currEv == nil {
				continue
			}

			// Figure out how much time is needed for one color:
			colors[currIdx].Duration = time.Millisecond * time.Duration(currEv.TotalMs/1000)

			if currEv.IsStopped {
				// Black out on stop, but wait a bit to save cpu time:
				sendColor(TimedColor{0, 0, 0, 500 * time.Millisecond})
			} else if currEv.IsPlaying || initialSend {
				// Send the color to the fader:
				sendColor(colors[currIdx])
				initialSend = false
			}

			// No need to go forth on "pause" or "stop":
			if currEv.IsPlaying {
				currIdx++
			}
		}
	}
}

// StatusUpdater is triggerd on "player" events and fetches the current state infos
// needed for the moodbar sync. It then proceeds to generate a MPDEvent that
// MoodbarAdjuster will receive.
func StatusUpdater(client *mpd.Client, cfg *Config, updateCh <-chan bool, eventCh chan<- MPDEvent) {
	lastSongID := ""

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

		// Check if the song changed compared to last time:
		// (always true for the first iteration)
		songChanged := false
		if status["songid"] != lastSongID {
			lastSongID = status["songid"]
			songChanged = true
		}

		// Find out how much progress we did in the current song.
		// These atteributes might be empty for the stopped state.
		elapsedMs, err := strconv.ParseFloat(status["elapsed"], 64)
		if err != nil && status["elapsed"] != "" {
			log.Printf("Failed to parse elpased (%s): %v", status["elapsed"], err)
		}

		elapsedMs *= 1000

		totalMs, err := strconv.Atoi(song["Time"])
		if err != nil && song["Time"] != "" {
			log.Printf("Failed to parse total (%s): %v", song["time"], err)
		}

		totalMs *= 1000

		// Find out if some music is playing...
		isPlaying, isStopped := false, false
		switch status["state"] {
		case "play":
			isPlaying = true
		case "stop":
			isStopped = true
		}

		// Send the appropiate event:
		eventCh <- MPDEvent{
			Path: filepath.Join(
				cfg.MoodDir,
				strings.Replace(song["file"], string(filepath.Separator), "|", -1),
			),
			SongChanged: songChanged,
			ElapsedMs:   elapsedMs,
			TotalMs:     float64(totalMs),
			IsPlaying:   isPlaying,
			IsStopped:   isStopped,
		}
	}

	close(eventCh)
}

// Watcher instances and connects via channels the go routines that contain the actual logic.
// It also triggers the logic by feeding mpd events to the go routine pipe.
func Watcher(client *mpd.Client, cfg *Config) error {
	// Watch for 'player' events:
	w, err := mpd.NewWatcher("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), "", "player")
	if err != nil {
		log.Fatalf("Failed to create watcher: %v", err)
		return err
	}

	defer w.Close()

	// Log mpd errors, but don't handle them more than that:
	go func() {
		for err := range w.Error {
			log.Println("Error:", err)
		}
	}()

	// Watcher -> StatusUpdater
	updateCh := make(chan bool)

	// StatusUpdater -> MoodbarAdjuster
	eventCh := make(chan MPDEvent)

	// MoodbarAdjuster -> MoodbarRunner
	colorsCh := make(chan TimedColor)

	// Start the respective go routines:
	go MoodbarRunner(cfg, colorsCh)
	go MoodbarAdjuster(eventCh, colorsCh)
	go StatusUpdater(client, cfg, updateCh, eventCh)

	// Also sync extra every few seconds:
	go func() {
		for range time.NewTicker(2 * time.Second).C {
			updateCh <- true
		}
	}()

	// ..but directly react on a changed player event:
	go func() {
		for range w.Event {
			updateCh <- true
		}
	}()

	// Block until something fatal happens:
	sigint := make(chan os.Signal)
	signal.Notify(sigint, os.Interrupt)
	<-sigint

	fmt.Println("\rInterrupted. Bye!")

	// Attempt to clean up, close() should be propagated to the
	// other channels by the respective go routines:
	close(updateCh)
	return nil
}

func RunDaemon(cfg *Config) error {
	client, err := mpd.Dial("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		log.Fatalf("Failed to connect to mpd: %v", err)
		return err
	}

	keepAlivePinger := make(chan bool)

	// Make sure the mpd connection survives long timeouts:
	go func() {
		for range keepAlivePinger {
			client.Ping()
			time.Sleep(1 * time.Minute)
		}
	}()

	// Close pinger and client on exit:
	defer func() {
		close(keepAlivePinger)
		client.Close()
	}()

	if cfg.UpdateMoodDatabase {
		if err := updateMoodDatabase(client, cfg); err != nil {
			log.Fatalf("Failed to update the mood db: %v", err)
		}

		return err
	}

	// Monitor MPD events and sync moodbar appropiately.
	return Watcher(client, cfg)
}
