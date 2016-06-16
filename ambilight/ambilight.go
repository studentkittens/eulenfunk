// moodymusic is a small mpd client that uses the `moodbar` utility
// and a multicolor LED to create an ambient light
package ambilight

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"

	// External dependencies:
	"github.com/fhs/gompd/mpd"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/studentkittens/eulenfunk/lightd"
)

type Config struct {
	// MPDHost of the mpd server (usually localhost)
	MPDHost string

	// MPDPort of the mpd server (usually 6600)
	MPDPort int

	// Host of the ambilight command server
	Host string

	// Port of the command server
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

type Server struct {
	sync.Mutex
	Config  *Config
	MPD     *mpd.Client
	Context context.Context
	Cancel  context.CancelFunc

	enabled bool
}

func (srv *Server) Enabled() bool {
	srv.Lock()
	defer srv.Unlock()

	return srv.enabled
}

func (srv *Server) Enable(enabled bool) {
	srv.Lock()
	defer srv.Unlock()

	srv.enabled = enabled
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
func updateMoodDatabase(server *Server) error {
	if server.Config.MoodDir == "" {
		return fmt.Errorf("No mood bar directory given (--mood-dir)")
	}

	if err := os.MkdirAll(server.Config.MoodDir, 0777); err != nil {
		return err
	}

	paths, err := server.MPD.GetFiles()
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
		moodPath := filepath.Join(server.Config.MoodDir, moodName)

		if _, err := os.Stat(moodPath); err == nil {
			// Already exists, Skipping.
			continue
		}

		dataPath := filepath.Join(server.Config.MusicDir, path)
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
		r, g, b := colorful.Hcl(h, c, l).LinearRgb()
		colors = append(colors, TimedColor{
			uint8(r * 255), uint8(g * 255), uint8(b * 255),
			((c1.Duration + c2.Duration) / 2) / time.Duration(N),
		})
	}

	return colors
}

func createDriverPipe(cfg *Config) (io.WriteCloser, error) {
	cmd := exec.Command(cfg.BinaryName, "cat")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Printf("Cannot fork catlight: %v", err)
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		log.Printf("Failed to start `catlight cat`: %v", err)
		return nil, err
	}

	return stdin, nil
}

// MoodbarRunner sets the current color and blends to it
// by remembering the last color and calculating a gradient between both.
func MoodbarRunner(server *Server, colors <-chan TimedColor) {
	cfg := server.Config

	stdin, err := createDriverPipe(cfg)
	if err != nil {
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

			if !server.Enabled() {
				continue
			}

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
			} else {
				// Blend is exhausted and no new colors available.
				// Sleep a bit to spare CPU.
				time.Sleep(50 * time.Millisecond)
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

	lightdConfig := &lightd.Config{
		Host: "localhost",
		Port: 3333,
	}

	locker, err := lightd.NewLocker(lightdConfig)
	if err != nil {
		log.Printf("Failed to create locker. Will continue without. lightd running?")
	} else {
		defer locker.Close()
	}

	sendColor := func(col TimedColor) {
		if locker != nil {
			if err := locker.Lock(); err != nil {
				log.Printf("Failed to acquire lock (sending anyways): %v", err)
			}
		}

		// Do not crash when colorsCh is closed:
		colorsCh <- col
		time.Sleep(col.Duration)

		if locker != nil {
			if err := locker.Unlock(); err != nil {
				log.Printf("Failed to unlock: %v", err)
			}
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
func StatusUpdater(server *Server, updateCh <-chan bool, eventCh chan<- MPDEvent) {
	lastSongID := ""

	for range updateCh {
		song, err := server.MPD.CurrentSong()
		if err != nil {
			log.Printf("Unable to fetch current song: %v", err)
			continue
		}

		status, err := server.MPD.Status()
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
				server.Config.MoodDir,
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
func Watcher(server *Server) error {
	// Watch for 'player' events:
	addr := fmt.Sprintf("%s:%d", server.Config.MPDHost, server.Config.MPDPort)

	log.Printf("Watching on %s", addr)
	w, err := mpd.NewWatcher("tcp", addr, "", "player")
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
	go MoodbarRunner(server, colorsCh)
	go MoodbarAdjuster(eventCh, colorsCh)
	go StatusUpdater(server, updateCh, eventCh)

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

	log.Printf("Press CTRL-C to interrupt")

	<-server.Context.Done()
	fmt.Println("\rInterrupted. Bye!")

	// Attempt to clean up, close() should be propagated to the
	// other channels by the respective go routines:
	close(updateCh)
	return nil
}

func handleConn(server *Server, driverStdin io.WriteCloser, conn net.Conn) {
	defer conn.Close()

	scn := bufio.NewScanner(conn)

	for scn.Scan() {
		switch cmd := scn.Text(); cmd {
		case "off":
			log.Printf("Disabling ambilight...")
			server.Enable(false)

			// Wait a short amount to make sure other colors
			// get flushed:
			time.Sleep(100 * time.Millisecond)
			if _, err := driverStdin.Write([]byte("0 0 0\n")); err != nil {
				log.Printf("Failed to turn light off: %v", err)
			}
		case "on":
			log.Printf("Enabling ambilight...")
			server.Enable(true)
		case "state":
			resp := []byte("0\n")
			if server.Enabled() {
				resp = []byte("1\n")
			}

			if _, err := conn.Write(resp); err != nil {
				log.Printf("Failed to write back state response: %v", err)
			}
		case "quit":
			log.Printf("Quitting ambilightd...")
			server.Cancel()
			return
		case "close":
			return
		}
	}

	if err := scn.Err(); err != nil {
		log.Printf("Failed to scan connection: %v", err)
	}
}

func createNetworkListener(server *Server) error {
	addr := fmt.Sprintf("%s:%d", server.Config.Host, server.Config.Port)
	lsn, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	log.Printf("Listening on %v", addr)

	stdin, err := createDriverPipe(server.Config)
	if err != nil {
		return err
	}

	go func() {
		defer stdin.Close()
		defer lsn.Close()

		for {
			select {
			case <-server.Context.Done():
				break
			default:
			}

			conn, err := lsn.Accept()
			if err != nil {
				log.Printf("Failed to accept connection: %v", err)
				break
			}

			log.Printf("Accepting connection from %s", conn.RemoteAddr())
			go handleConn(server, stdin, conn)
		}
	}()

	return nil
}

func RunDaemon(cfg *Config, ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", cfg.MPDHost, cfg.MPDPort)
	mpdClient, err := mpd.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to connect to mpd (%s): %v", addr, err)
		return err
	}

	subCtx, cancel := context.WithCancel(ctx)

	server := &Server{
		Config:  cfg,
		MPD:     mpdClient,
		Context: subCtx,
		Cancel:  cancel,
	}

	server.Enable(true)

	if err := createNetworkListener(server); err != nil {
		return err
	}

	keepAlivePinger := make(chan bool)

	// Make sure the mpd connection survives long timeouts:
	go func() {
		for range keepAlivePinger {
			mpdClient.Ping()
			time.Sleep(1 * time.Minute)
		}
	}()

	// Close pinger and client on exit:
	defer func() {
		close(keepAlivePinger)
		mpdClient.Close()
	}()

	if cfg.UpdateMoodDatabase {
		if err := updateMoodDatabase(server); err != nil {
			log.Fatalf("Failed to update the mood db: %v", err)
		}

		return err
	}

	log.Printf("Starting up...")

	// Monitor MPD events and sync moodbar appropiately.
	return Watcher(server)
}
