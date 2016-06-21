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
	gompd "github.com/fhs/gompd/mpd"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/studentkittens/eulenfunk/lightd"
	"github.com/studentkittens/eulenfunk/ui/mpd"
	"github.com/studentkittens/eulenfunk/util"
)

// UseDefaultMoodbar enables a builtin default moodbar if none was found
const UseDefaultMoodbar = true

// Config holds all possible adjusting screws for ambilightd.
type Config struct {
	// MPDHost of the mpd server (usually localhost)
	MPDHost string

	// MPDPort of the mpd server (usually 6600)
	MPDPort int

	// Lightd host (usually localhost)
	LightdHost string

	// Lightd port  (usually 3333)
	LightdPort int

	// Host of the ambilight command server
	AmbiHost string

	// Port of the command server
	AmbiPort int

	// MusicDir is the root path of the mpd database
	MusicDir string

	// MoodDir contains all moodfiles for certain files (if any)
	MoodDir string

	// UpdateMoodDatabase makes the client update the db and exit afterwards.
	UpdateMoodDatabase bool

	// Name of the RGB LED driver binary (`catlight` for my desktop)
	BinaryName string
}

// server holds all runtime info for ambilightd.
type server struct {
	Config *Config
	MPD    *mpd.ReMPD

	// Cancellation:
	Context context.Context
	Cancel  context.CancelFunc

	stateCh chan bool
	enabled bool

	mu sync.Mutex
}

// Current status of the MPD player:
type mpdEvent struct {
	Path        string
	ElapsedMs   float64
	TotalMs     float64
	IsPlaying   bool
	IsStopped   bool
	SongChanged bool
}

// Info needed to render a moodbar:
type moodInfo struct {
	MusicFile string
	MoodPath  string
}

// RGB Color that stays for a certain duration:
type timedColor struct {
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
func updateMoodDatabase(server *server) error {
	checkForMoodbar()

	if server.Config.MoodDir == "" {
		return fmt.Errorf("No mood bar directory given (--mood-dir)")
	}

	if err := os.MkdirAll(server.Config.MoodDir, 0777); err != nil {
		return err
	}

	paths, err := server.MPD.Client().GetFiles()
	if err != nil {
		return fmt.Errorf("Cannot get all files from mpd: %v", err)
	}

	// Use up to N threads:
	N := 8
	wg := &sync.WaitGroup{}
	wg.Add(N)

	moodChan := make(chan *moodInfo, N)
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
		moodChan <- &moodInfo{
			MusicFile: dataPath,
			MoodPath:  moodPath,
		}
	}

	close(moodChan)
	wg.Wait()

	return nil
}

// Read the .mood file at `path`; return a 1000-element slice of timedColor.
// The duration of each color will be 0.
func readMoodbarFile(path string) ([]timedColor, error) {
	results := []timedColor{}

	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer util.Closer(fd)

	rgb := make([]byte, 3)
	for {
		_, err := fd.Read(rgb)
		if err != nil && err != io.EOF {
			return nil, err
		}

		results = append(results, timedColor{
			rgb[0], rgb[1], rgb[2], 0,
		})

		if err == io.EOF {
			break
		}
	}

	return results, nil
}

// Create a HCL Gradient between c1 and c2 using N steps.
// Returns the gradient as slice of individual colors.
func createBlend(c1, c2 timedColor, N int) []timedColor {
	// Do nothing if it's the same color:
	if c1.R == c2.R && c1.G == c2.G && c1.B == c2.B {
		return []timedColor{c1}
	}

	cc1 := colorful.Color{
		R: float64(c1.R) / 255.,
		G: float64(c1.G) / 255.,
		B: float64(c1.B) / 255.,
	}

	cc2 := colorful.Color{
		R: float64(c2.R) / 255.,
		G: float64(c2.G) / 255.,
		B: float64(c2.B) / 255.,
	}

	colors := []timedColor{}

	for i := 0; i < N; i++ {
		mix := cc1.BlendHcl(cc2, float64(i)/float64(N)).Clamped()
		h, c, l := mix.Hcl()

		// Increase chroma slightly for mid values:
		c = (-((c-1)*(c-1))+1)/2 + c/2

		// Decrease lumninance slightly for mid values:
		l = (l*l)/2 + (l / 2)

		// Convert back to (gamma corrected) RGB for catlight:
		r, g, b := colorful.Hcl(h, c, l).FastLinearRgb()
		//hcl := colorful.Hcl(h, c, l)
		//r, g, b := hcl.R, hcl.G, hcl.B
		colors = append(colors, timedColor{
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

// moodbarRunner sets the current color and blends to it
// by remembering the last color and calculating a gradient between both.
func moodbarRunner(server *server, colors <-chan timedColor) {
	cfg := server.Config

	stdin, err := createDriverPipe(cfg)
	if err != nil {
		return
	}

	defer util.Closer(stdin)

	// First color is always black.
	var lastColor timedColor

	// Blend between lastColor and current color.
	var blend []timedColor

	var enabled = true

	for {
		select {
		case newState := <-server.stateCh:
			enabled = newState
			if !enabled {
				blend = []timedColor{timedColor{0, 0, 0, 0}}
			}
		case color, ok := <-colors:
			if !ok {
				return
			}

			if !enabled {
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
				if _, err := stdin.Write([]byte(colorValue)); err != nil {
					log.Printf("Failed to write color to driver: %v", err)
				}

				time.Sleep(color.Duration)
			}
		}
	}
}

func sendColor(locker *lightd.Locker, col timedColor, colorsCh chan<- timedColor) {
	if locker != nil {
		if err := locker.Lock(); err != nil {
			log.Printf("Failed to acquire lock (sending anyways): %v", err)
		}
	}

	colorsCh <- col
	time.Sleep(col.Duration)

	if locker != nil {
		if err := locker.Unlock(); err != nil {
			log.Printf("Failed to unlock: %v", err)
		}
	}
}

func adjust(srv *server, locker *lightd.Locker, ev *mpdEvent, colorsCh chan<- timedColor, colors *[]timedColor) int {
	currIdx := 0

	// Adjust the moodbar seek offset (1000 samples per total time)
	if ev.TotalMs > 0 {
		currIdx = int((ev.ElapsedMs / ev.TotalMs) * 1000)
	}

	if !ev.SongChanged {
		return currIdx
	}

	data, err := readMoodbarFile(ev.Path)
	if err == nil {
		*colors = data
		return currIdx
	}

	log.Printf("Failed to read moodbar at `%s`: %v", ev.Path, err)

	// Return to black:
	if UseDefaultMoodbar {
		*colors = DefaultMoodbar
	} else {
		sendColor(locker, timedColor{0, 0, 0, 0}, colorsCh)
		*colors = []timedColor{}
		return 0
	}

	return currIdx
}

func eatColor(locker *lightd.Locker, currEv *mpdEvent, currCol *timedColor, colorsCh chan<- timedColor, initialSend *bool) {

	// Figure out how much time is needed for one color:
	(*currCol).Duration = time.Millisecond * time.Duration(currEv.TotalMs/1000)

	if currEv.IsStopped {
		// Black out on stop, but wait a bit to save cpu time:
		sendColor(locker, timedColor{0, 0, 0, 500 * time.Millisecond}, colorsCh)
	} else if currEv.IsPlaying || *initialSend {
		// Send the color to the fader:
		sendColor(locker, *currCol, colorsCh)
		*initialSend = false
	}
}

// moodbarAdjuster tried to synchronize the music to the moodbar.
// It will send the correct current color to moodbarRunner.
func moodbarAdjuster(srv *server, eventCh <-chan mpdEvent, colorsCh chan<- timedColor) {
	var (
		currIdx int
		colors  []timedColor
		currEv  *mpdEvent
	)

	initialSend := true

	lightdConfig := &lightd.Config{
		Host: srv.Config.LightdHost,
		Port: srv.Config.LightdPort,
	}

	locker, err := lightd.NewLocker(lightdConfig)
	if err != nil {
		log.Printf("Failed to create locker - will continue without. lightd running?")
	} else {
		defer util.Closer(locker)
	}

	defer func() {
		close(colorsCh)
	}()

	for {
		select {
		case ev, ok := <-eventCh:
			if !ok {
				return
			}

			// A new event happened, we need to adjust or even load a new moodbar file:
			currIdx = adjust(srv, locker, &ev, colorsCh, &colors)
			currEv = &ev
		default:
			if currIdx >= len(colors) || currEv == nil {
				continue
			}

			// Nothing happened, give the led some input:
			eatColor(locker, currEv, &colors[currIdx], colorsCh, &initialSend)

			// No need to go forth on "pause" or "stop":
			if currEv.IsPlaying {
				currIdx++
			}
		}
	}
}

func fetchMPDInfo(client *gompd.Client) (gompd.Attrs, gompd.Attrs, error) {
	song, err := client.CurrentSong()
	if err != nil {
		log.Printf("Unable to fetch current song: %v", err)
		return nil, nil, err
	}

	status, err := client.Status()
	if err != nil {
		log.Printf("Unable to fetch status: %v", err)
		return nil, nil, err
	}

	return song, status, nil
}

// statusUpdater is triggerd on "player" events and fetches the current state infos
// needed for the moodbar sync. It then proceeds to generate a mpdEvent that
// moodbarAdjuster will receive.
func statusUpdater(server *server, updateCh <-chan bool, eventCh chan<- mpdEvent) {
	lastSongID := ""

	for range updateCh {
		song, status, err := fetchMPDInfo(server.MPD.Client())
		if err != nil {
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

		// Send the appropriate event:
		eventCh <- mpdEvent{
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
func Watcher(server *server) error {
	// Watch for 'player' events:
	addr := fmt.Sprintf("%s:%d", server.Config.MPDHost, server.Config.MPDPort)

	log.Printf("Watching on %s", addr)
	watcher := mpd.NewReWatcher(
		server.Config.MPDHost, server.Config.MPDPort,
		server.Context,
		"player",
	)

	defer util.Closer(watcher)

	// Watcher -> statusUpdater
	updateCh := make(chan bool)

	// statusUpdater -> moodbarAdjuster
	eventCh := make(chan mpdEvent)

	// moodbarAdjuster -> moodbarRunner
	colorsCh := make(chan timedColor)

	// Start the respective go routines:
	go moodbarRunner(server, colorsCh)
	go moodbarAdjuster(server, eventCh, colorsCh)
	go statusUpdater(server, updateCh, eventCh)

	// Also sync extra every few seconds:
	go func() {
		for range time.NewTicker(2 * time.Second).C {
			updateCh <- true
		}
	}()

	// ..but directly react on a changed player event:
	go func() {
		for range watcher.Events {
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

func handleConn(server *server, conn net.Conn) {
	defer util.Closer(conn)

	scn := bufio.NewScanner(conn)
	for scn.Scan() {
		switch cmd := scn.Text(); cmd {
		case "off":
			log.Printf("Disabling ambilight...")
			server.stateCh <- false

			server.mu.Lock()
			server.enabled = false
			server.mu.Unlock()
		case "on":
			log.Printf("Enabling ambilight...")
			server.stateCh <- true

			server.mu.Lock()
			server.enabled = true
			server.mu.Unlock()
		case "state":
			resp := []byte("0\n")

			server.mu.Lock()
			if server.enabled {
				resp = []byte("1\n")
			}
			server.mu.Unlock()

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

func createNetworkListener(server *server) error {
	addr := fmt.Sprintf("%s:%d", server.Config.AmbiHost, server.Config.AmbiPort)
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
		defer util.Closer(stdin)
		defer util.Closer(lsn)

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
			go handleConn(server, conn)
		}
	}()

	return nil
}

func keepAlivePinger(MPD *mpd.ReMPD, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if err := MPD.Client().Ping(); err != nil {
			log.Printf("Failed to ping MPD server. Weird: %v", err)
		}

		time.Sleep(1 * time.Minute)
	}
}

// Run starts ambilightd with the settings defined in `cfg`.
// It will stop execution when `ctx` was canceled.
// If something show-stopping occurs on startup an error is returned.
func Run(cfg *Config, ctx context.Context) error {
	subCtx, cancel := context.WithCancel(ctx)
	MPD := mpd.NewReMPD(cfg.MPDHost, cfg.MPDPort, subCtx)

	// Close pinger and client on exit:
	defer util.Closer(MPD.Client())

	server := &server{
		Config:  cfg,
		MPD:     MPD,
		Context: subCtx,
		Cancel:  cancel,
		stateCh: make(chan bool),
		enabled: true,
	}

	if err := createNetworkListener(server); err != nil {
		return err
	}

	// Make sure the mpd connection survives long timeouts:
	go keepAlivePinger(MPD, ctx)

	if cfg.UpdateMoodDatabase {
		if err := updateMoodDatabase(server); err != nil {
			log.Printf("Failed to update the mood db: %v", err)
			return err
		}

		return nil
	}

	log.Printf("Starting up...")

	// Monitor MPD events and sync moodbar appropriately.
	return Watcher(server)
}
