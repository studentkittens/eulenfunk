package ui

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"golang.org/x/net/context"

	"github.com/studentkittens/eulenfunk/ambilight"
	"github.com/studentkittens/eulenfunk/display"
	"github.com/studentkittens/eulenfunk/ui/mpd"
	"github.com/studentkittens/eulenfunk/util"
)

type Action func() error

// TODO: find better name
type MenuLine interface {
	Render(w int, active bool) string
	Name() string
	Action() error
	Selectable() bool
}

type Separator struct {
	Title string
}

func (sp *Separator) Render(w int, active bool) string {
	return util.Center(strings.ToUpper(" "+sp.Title+" "), w, '━')
}

func (sp *Separator) Name() string {
	return sp.Title
}

func (sp *Separator) Action() error {
	return nil
}

func (sp *Separator) Selectable() bool {
	return false
}

type Entry struct {
	Text       string
	ActionFunc Action
	State      string
}

func (en *Entry) Render(w int, active bool) string {
	prefix := "  "
	if active {
		prefix = "❤ "
	}

	state := en.State
	if state != "" {
		state = " [" + state + "]"
	}

	m := w - utf8.RuneCountInString(state) - utf8.RuneCountInString(prefix)
	return fmt.Sprintf("%s%-*s%s", prefix, m, en.Text, state)
}

func (en *Entry) Name() string {
	return en.Text
}

func (en *Entry) Action() error {
	if en.ActionFunc != nil {
		return en.ActionFunc()
	}

	return nil
}

func (en *Entry) Selectable() bool {
	return true
}

////////////////////////

type Menu struct {
	Name    string
	Entries []MenuLine
	Cursor  int

	lw *display.LineWriter
}

func NewMenu(name string, lw *display.LineWriter) (*Menu, error) {
	return &Menu{
		Name: name,
		lw:   lw,
	}, nil
}

func (mn *Menu) ActiveEntryName() string {
	if len(mn.Entries) == 0 {
		return ""
	}

	return mn.Entries[mn.Cursor].Name()
}

func (mn *Menu) scrollToNextSelectable(up bool) {
	for mn.Cursor >= 0 && mn.Cursor < len(mn.Entries) && !mn.Entries[mn.Cursor].Selectable() {
		if up {
			mn.Cursor++
		} else {
			mn.Cursor--
		}
	}

	// Clamp value if nothing suitable was found in that direction:
	if mn.Cursor < 0 {
		mn.Cursor = 0
	}

	if mn.Cursor >= len(mn.Entries) {
		mn.Cursor = len(mn.Entries) - 1
	}
}

func (mn *Menu) Scroll(move int) {
	mn.Cursor += move

	up := move >= 0
	mn.scrollToNextSelectable(up)

	// Check if we succeeded:
	if !mn.Entries[mn.Cursor].Selectable() {
		mn.scrollToNextSelectable(!up)
	}
}

func (mn *Menu) Display(width int) error {
	for pos, entry := range mn.Entries {
		line := entry.Render(width, pos == mn.Cursor)
		log.Printf("ACtive %t %d == %d %s", pos == mn.Cursor, pos, mn.Cursor, line)

		if _, err := mn.lw.Formatf("line %s %d %s", mn.Name, pos, line); err != nil {
			return err
		}
	}

	return nil
}

func (mn *Menu) Click() error {
	if len(mn.Entries) == 0 {
		return nil
	}

	entry := mn.Entries[mn.Cursor]
	if entry.Action == nil {
		return nil
	}

	return entry.Action()
}

////////////////////////

const (
	// No movement (initial)
	DirectionNone = iota
	DirectionRight
	DirectionLeft
)

type MenuManager struct {
	sync.Mutex
	Config *Config

	// TODO: cleanup
	Active       *Menu
	Menus        map[string]*Menu
	TimedActions map[time.Duration]Action

	lw                   *display.LineWriter
	rotateActions        []Action
	releaseActions       []Action
	currValue, lastValue int
	activeWindow         string
	rotary               *util.Rotary
}

func NewMenuManager(cfg *Config, lw *display.LineWriter, initialWin string) (*MenuManager, error) {
	rty, err := util.NewRotary()
	if err != nil {
		return nil, err
	}

	// Switch to mpd initially:
	if _, err := lw.Formatf("switch %s", initialWin); err != nil {
		return nil, err
	}

	mgr := &MenuManager{
		Menus:        make(map[string]*Menu),
		TimedActions: make(map[time.Duration]Action),
		activeWindow: initialWin,
		Config:       cfg,
		lw:           lw,
		rotary:       rty,
	}

	timedActionExecd := make(map[time.Duration]bool)

	go func() {
		for state := range rty.Button {
			if mgr.Active == nil {
				continue
			}

			if !state {
				fmt.Println("Button released")
				mgr.Lock()
				timedActionExecd = make(map[time.Duration]bool)
				mgr.Unlock()

				for idx, action := range mgr.releaseActions {
					if err := action(); err != nil {
						log.Printf("release action %d failed: %v", idx, err)
					}
				}

				continue
			}

			fmt.Println("Button pressed")

			mgr.Lock()
			active := mgr.Active
			mgr.Unlock()

			// Check if we're actually in a menu:
			if mgr.ActiveWindow() == active.Name {
				// if strings.HasPrefix(mgr.ActiveWindow(), "menu-") {
				if err := active.Click(); err != nil {
					name := active.ActiveEntryName()
					log.Printf("Action for menu entry `%s` failed: %v", name, err)
				}
			}
		}
	}()

	go func() {
		for duration := range rty.Pressed {
			log.Printf("Pressed for %s", duration)
			mgr.Lock()

			// Find the action with smallest non-negative diff:
			var diff time.Duration
			var action Action
			var actionTime time.Duration

			for after, timedAction := range mgr.TimedActions {
				newDiff := duration - after
				if after <= duration && (action == nil || newDiff < diff) {
					diff = duration - after
					action = timedAction
					actionTime = after
				}
			}

			if action != nil && !timedActionExecd[actionTime] {
				// Call action() unlocked:
				mgr.Unlock()
				action()
				mgr.Lock()

				timedActionExecd[actionTime] = true
			}

			mgr.Unlock()
		}
	}()

	go func() {
		for value := range rty.Value {
			mgr.Lock()
			if mgr.Active == nil {
				mgr.Unlock()
				continue
			}

			mgr.lastValue = mgr.currValue
			mgr.currValue = value
			diff := mgr.currValue - mgr.lastValue
			name := mgr.Active.Name
			mgr.Unlock()

			log.Printf("Value: %d Diff %d\n", value, diff)

			mgr.Active.Scroll(diff)
			if _, err := lw.Formatf("move %s %d", name, diff); err != nil {
				log.Printf("move failed: %v", err)
			}

			for idx, action := range mgr.rotateActions {
				if err := action(); err != nil {
					log.Printf("Rotate action %d failed: %v", idx, err)
				}
			}

			mgr.Active.Display(mgr.Config.Width)
		}
	}()

	return mgr, nil
}

func (mgr *MenuManager) Display() error {
	mgr.Lock()
	defer mgr.Unlock()

	return mgr.Active.Display(mgr.Config.Width)
}

func (mgr *MenuManager) ActiveWindow() string {
	mgr.Lock()
	defer mgr.Unlock()

	return mgr.activeWindow
}

func (mgr *MenuManager) Direction() int {
	mgr.Lock()
	defer mgr.Unlock()

	switch {
	case mgr.lastValue < mgr.currValue:
		return DirectionRight
	case mgr.lastValue > mgr.currValue:
		return DirectionLeft
	default:
		return DirectionNone
	}
}

func (mgr *MenuManager) Value() int {
	mgr.Lock()
	defer mgr.Unlock()

	return mgr.currValue
}

func (mgr *MenuManager) RotateAction(a Action) {
	mgr.Lock()
	defer mgr.Unlock()

	mgr.rotateActions = append(mgr.rotateActions, a)
}

func (mgr *MenuManager) ReleaseAction(a Action) {
	mgr.Lock()
	defer mgr.Unlock()

	mgr.releaseActions = append(mgr.releaseActions, a)
}

func (mgr *MenuManager) SwitchTo(name string) error {
	mgr.Lock()
	defer mgr.Unlock()

	if menu, ok := mgr.Menus[name]; ok {
		mgr.Active = menu
		mgr.Active.Display(mgr.Config.Width)
	}

	if _, err := mgr.lw.Formatf("switch %s", name); err != nil {
		log.Printf("switch failed: %v", err)
		return err
	}

	mgr.activeWindow = name
	return nil
}

func (mgr *MenuManager) AddTimedAction(after time.Duration, action Action) {
	mgr.Lock()
	defer mgr.Unlock()

	mgr.TimedActions[after] = action
}

func (mgr *MenuManager) AddMenu(name string, entries []MenuLine) error {
	mgr.Lock()
	defer mgr.Unlock()

	menu, err := NewMenu(name, mgr.lw)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		menu.Entries = append(menu.Entries, entry)
	}

	if mgr.Active == nil {
		mgr.Active = menu
		mgr.Active.Display(mgr.Config.Width)
	}

	// Pre-select first selectable entry:
	menu.Scroll(0)

	mgr.Menus[name] = menu
	return nil
}

func (mgr *MenuManager) Close() error {
	return mgr.rotary.Close()
}

//////////////////////////////////////

func sysCommand(name string, args ...string) func() error {
	return func() error {
		return exec.Command(name, args...).Run()
	}
}

func boolToGlyph(b bool) string {
	if b {
		return "✓"
	} else {
		return ""
	}
}

func createPartyModeEntry(cfg *Config, mgr *MenuManager) (*Entry, error) {
	client, err := ambilight.NewClient(&ambilight.Config{
		Host: cfg.AmbilightHost,
		Port: cfg.AmbilightPort,
	})

	if err != nil {
		return nil, err
	}

	enabled, err := client.Enabled()
	if err != nil {
		return nil, err
	}

	partyModeEntry := &Entry{
		Text:  "Party!",
		State: boolToGlyph(enabled),
	}

	partyModeEntry.ActionFunc = func() error {
		defer client.Close()

		enabled, err := client.Enabled()
		if err != nil {
			return err
		}

		if err := client.Enable(!enabled); err != nil {
			return err
		}

		partyModeEntry.State = boolToGlyph(enabled)

		// TODO: Just all DIsplay() after each action automatically?
		return mgr.Display()
	}

	return partyModeEntry, nil
}

func createOutputEntry(mgr *MenuManager, MPD *mpd.Client) (*Entry, error) {
	outputs, err := MPD.Outputs()
	if err != nil {
		return nil, err
	}

	initialOutput := ""
	if len(outputs) > 0 {
		initialOutput = outputs[0]
	}

	outputEntry := &Entry{
		Text:  "Output",
		State: initialOutput,
	}

	outputEntry.ActionFunc = func() error {
		outputs, err := MPD.Outputs()
		if err != nil {
			return err
		}

		idx := 0
		for id, output := range outputs {
			if output == outputEntry.State {
				idx = id
				break
			}
		}

		newOutput := outputs[(idx+1)%len(outputs)]
		outputEntry.State = newOutput

		if err := MPD.SwitchToOutput(newOutput); err != nil {
			return err
		}

		return mgr.Display()
	}

	return outputEntry, nil
}

func createPlaybackEntry(mgr *MenuManager, MPD *mpd.Client) (*Entry, error) {
	playbackEntry := &Entry{
		Text:  "Playback",
		State: mpd.StateToUnicode(MPD.CurrentState()),
	}

	order := []string{"play", "pause", "stop"}

	playbackEntry.ActionFunc = func() error {
		idx := 0
		currState := MPD.CurrentState()
		for orderIdx, state := range order {
			if state == currState {
				idx = orderIdx
				break
			}
		}

		newState := order[(idx+1)%len(order)]
		playbackEntry.State = mpd.StateToUnicode(newState)

		var err error
		switch newState {
		case "play":
			err = MPD.Play()
		case "pause":
			err = MPD.Pause()
		case "stop":
			err = MPD.Stop()
		}

		if err != nil {
			return err
		}

		return mgr.Display()
	}

	return playbackEntry, nil
}

func createRandomEntry(mgr *MenuManager, MPD *mpd.Client) (*Entry, error) {
	randomEntry := &Entry{
		Text:  "Random",
		State: boolToGlyph(MPD.IsRandom()),
	}

	randomEntry.ActionFunc = func() error {
		enable := !MPD.IsRandom()
		randomEntry.State = boolToGlyph(enable)
		if err := MPD.EnableRandom(enable); err != nil {
			return err
		}

		return mgr.Display()
	}

	return randomEntry, nil
}

/////////////////////////

type Config struct {
	Width  int
	Height int

	DisplayHost string
	DisplayPort int

	MPDHost string
	MPDPort int

	AmbilightHost string
	AmbilightPort int
}

/////////////////////////
// MENU MAINLOOP LOGIC //
/////////////////////////

func Run(cfg *Config, ctx context.Context) error {
	log.Printf("Connecting to displayd...")
	lw, err := display.Connect(&display.Config{
		Host: cfg.DisplayHost,
		Port: cfg.DisplayPort,
	})

	if err != nil {
		return err
	}

	defer lw.Close()

	if err := drawStaticScreens(lw); err != nil {
		log.Printf("Failed to draw static screens: %v", err)
		return err
	}

	log.Printf("Creating menus...")
	mgr, err := NewMenuManager(cfg, lw, "startup")
	if err != nil {
		return err
	}

	// Some flags to coordinate actions:
	ignoreRelease := false

	switcher := func(name string) func() error {
		return func() error {
			ignoreRelease = true
			return mgr.SwitchTo(name)
		}
	}

	// Start auxillary services:
	log.Printf("Starting background services...")
	MPD, err := mpd.NewClient(&mpd.Config{
		Host:        cfg.MPDHost,
		Port:        cfg.MPDPort,
		DisplayHost: cfg.DisplayHost,
		DisplayPort: cfg.DisplayPort,
	})

	if err != nil {
		log.Printf("Failed to create mpd client: %v", err)
		return err
	}

	go MPD.Run(ctx)
	go RunClock(lw, cfg.Width, ctx)
	go RunSysinfo(lw, cfg.Width, ctx)

	// Create some special entries with extended logic:

	outputEntry, err := createOutputEntry(mgr, MPD)
	if err != nil {
		log.Printf("Failed to create output entry: %v", err)
		return err
	}

	partyModeEntry, err := createPartyModeEntry(cfg, mgr)
	if err != nil {
		log.Printf("Failed to create party-mode entry: %v", err)
		return err
	}

	playbackEntry, err := createPlaybackEntry(mgr, MPD)
	if err != nil {
		log.Printf("Failed to create playback entry: %v", err)
		return err
	}

	randomEntry, err := createRandomEntry(mgr, MPD)
	if err != nil {
		log.Printf("Failed to create random entry: %v", err)
		return err
	}

	// Define the menu structure:

	mainMenu := []MenuLine{
		&Separator{"MODES"},
		&Entry{
			Text:       "Music status",
			ActionFunc: switcher("mpd"),
		},
		&Entry{
			Text: "Playlists",
			ActionFunc: func() error {
				return showPlaylistWindow(lw, MPD)
			},
		},
		&Entry{
			Text:       "Clock",
			ActionFunc: switcher("clock"),
		},
		&Entry{
			Text:       "System info",
			ActionFunc: switcher("sysinfo"),
		},
		&Entry{
			Text:       "Statistics",
			ActionFunc: switcher("stats"),
		},
		&Separator{"OPTIONS"},
		partyModeEntry,
		outputEntry,
		playbackEntry,
		randomEntry,
		&Separator{"SYSTEM"},
		&Entry{
			Text:       "Powermenu",
			ActionFunc: switcher("menu-power"),
		},
		&Entry{
			Text:       "About",
			ActionFunc: switcher("about"),
		},
	}

	powerMenu := []MenuLine{
		&Entry{
			Text: "Poweroff",
			ActionFunc: func() error {
				switchToStatic(lw, "shutdown")
				return sysCommand("systemctl", "poweroff")()
			},
		},
		&Entry{
			Text: "Reboot",
			ActionFunc: func() error {
				switchToStatic(lw, "shutdown")
				return sysCommand("systemctl", "reboot")()
			},
		},
		&Entry{
			Text:       "Exit",
			ActionFunc: switcher("menu-main"),
		},
	}

	if err := mgr.AddMenu("menu-main", mainMenu); err != nil {
		log.Printf("Add main-menu failed: %v", err)
		return err
	}

	if err := mgr.AddMenu("menu-power", powerMenu); err != nil {
		log.Printf("Add main-power failed: %v", err)
		return err
	}

	mgr.AddTimedAction(10*time.Millisecond, func() error {
		return nil
	})

	mgr.AddTimedAction(600*time.Millisecond, func() error {
		ignoreRelease = true
		return mgr.SwitchTo("menu-main")
	})

	mgr.AddTimedAction(3*time.Second, func() error {
		ignoreRelease = true
		return mgr.SwitchTo("menu-power")
	})

	mgr.AddTimedAction(8*time.Second, func() error {
		ignoreRelease = true
		cmd := sysCommand("aplay", "/root/hoot.wav")
		go cmd()
		return nil
	})

	mgr.ReleaseAction(func() error {
		if ignoreRelease {
			ignoreRelease = false
			return nil
		}

		switch currWin := mgr.ActiveWindow(); currWin {
		case "mpd":
			if err := MPD.TogglePlayback(); err != nil {
				log.Printf("Failed to toggle playback: %v", err)
			}
		default:
			// This is a bit of a hack:
			// Enable "click to exit window" on most non-menu windows:
			if !strings.Contains(currWin, "menu") {
				return mgr.SwitchTo("menu-main")
			}
		}

		return nil
	})

	mgr.RotateAction(func() error {
		if mgr.ActiveWindow() != "mpd" {
			return nil
		}

		log.Printf("rotate action")
		switch mgr.Direction() {
		case DirectionRight:
			log.Printf("Play next")
			if err := MPD.Next(); err != nil {
				log.Printf("Failed to skip to next: %v", err)
			}
		case DirectionLeft:
			log.Printf("Play prev")
			if err := MPD.Next(); err != nil {
				log.Printf("Failed to skip to prev: %v", err)
			}
		}

		return nil
	})

	log.Printf("Waiting for a silent death...")
	<-ctx.Done()

	return mgr.Close()
}
