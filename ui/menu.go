package menu

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"

	"github.com/studentkittens/eulenfunk/ambilight"
	"github.com/studentkittens/eulenfunk/display"
	"github.com/studentkittens/eulenfunk/ui/mpdinfo"
	"github.com/studentkittens/eulenfunk/util"
)

type Action func() error

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

	return fmt.Sprintf("%s%-*s%s", prefix, w-len(state)-len(prefix), en.Text, state)
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

func (mn *Menu) Display() error {
	for pos, entry := range mn.Entries {
		// TODO: pass config width
		line := entry.Render(20, pos == mn.Cursor)
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

func NewMenuManager(lw *display.LineWriter, initialWin string) (*MenuManager, error) {
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
		activeWindow: "mpd",
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

			mgr.Active.Display()
		}
	}()

	return mgr, nil
}

func (mgr *MenuManager) Display() error {
	mgr.Lock()
	defer mgr.Unlock()

	return mgr.Active.Display()
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
		mgr.Active.Display()
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
		mgr.Active.Display()
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

func mpdCommand(name string, mpdCmdCh chan<- string) func() error {
	return func() error {
		mpdCmdCh <- name
		return nil
	}
}

func sysCommand(name string, args ...string) func() error {
	return func() error {
		return exec.Command(name, args...).Run()
	}
}

func drawShutdownscreen(lw *display.LineWriter) error {
	startupScreen := []string{
		"SHUTTING DOWN - BYE!",
		"                    ",
		"PLEASE WAIT 1 MINUTE",
		"BEFORE POWERING OFF!",
	}

	for idx, line := range startupScreen {
		if _, err := lw.Formatf("line mpd %d %s", idx, line); err != nil {
			return err
		}
	}

	return nil
}

func drawStartupScreen(lw *display.LineWriter) error {
	startupScreen := []string{
		"/ / / / / / / / / / / / / / / /",
		"WELCOME TO EULENFUNK",
		" GUT. ECHT. ANDERS. ",
		"/ / / / / / / / / / / / / / / /",
	}

	for idx, line := range startupScreen {
		if _, err := lw.Formatf("line mpd %d %s", idx, line); err != nil {
			return err
		}
	}

	if _, err := lw.Formatf("scroll mpd 0 150ms"); err != nil {
		return err
	}

	if _, err := lw.Formatf("scroll mpd 3 200ms"); err != nil {
		return err
	}

	return nil
}

func Run(ctx context.Context) error {
	// TODO: pass config
	cfg := &display.Config{
		Host: "localhost",
		Port: 7778,
	}

	log.Printf("Connecting to displayd...")
	lw, err := display.Connect(cfg)
	if err != nil {
		return err
	}

	defer lw.Close()

	if err := drawStartupScreen(lw); err != nil {
		log.Printf("Failed to draw startup screen: %v", err)
		return err
	}

	log.Printf("Creating menus...")
	mgr, err := NewMenuManager(lw, "mpd")
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

	mpdCmdCh := make(chan string)

	// Start auxillary services:
	log.Printf("Starting background services...")
	go mpdinfo.Run(&mpdinfo.Config{
		Host:        "localhost",
		Port:        6600,
		DisplayHost: "localhost",
		DisplayPort: 7778,
	}, ctx, mpdCmdCh)

	go RunClock(lw, 20, ctx) // TODO: get width?
	go RunSysinfo(lw, 20, ctx)

	partyModeEntry := &Entry{
		Text:  "Party!",
		State: "✓",
	}

	partyModeEntry.ActionFunc = func() error {
		client, err := ambilight.NewClient(&ambilight.Config{
			Host: "localhost",
			Port: 4444,
		})

		if err != nil {
			return err
		}

		defer client.Close()

		enabled, err := client.Enabled()
		if err != nil {
			return err
		}

		if err := client.Enable(!enabled); err != nil {
			return err
		}

		if enabled {
			partyModeEntry.State = "×"
		} else {
			partyModeEntry.State = "✓"
		}

		return mgr.Display()
	}

	mainMenu := []MenuLine{
		&Separator{"MODES"},
		&Entry{
			Text:       "Music status",
			ActionFunc: switcher("mpd"),
		},
		&Entry{
			Text:       "Playlists",
			ActionFunc: switcher("playlists"),
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
		&Entry{
			Text:       "Switch Mono/Stereo",
			ActionFunc: nil, // TODO
		},
		&Entry{
			Text:       "Playback",
			ActionFunc: mpdCommand("stop", mpdCmdCh),
			State:      "⏹",
		},
		&Separator{"SYSTEM"},
		&Entry{
			Text:       "Power",
			ActionFunc: switcher("menu-power"),
		},
	}

	powerMenu := []MenuLine{
		&Entry{
			Text: "Poweroff",
			ActionFunc: func() error {
				drawShutdownscreen(lw)
				return sysCommand("systemctl", "poweroff")()
			},
		},
		&Entry{
			Text: "Reboot",
			ActionFunc: func() error {
				drawShutdownscreen(lw)
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

	mgr.AddTimedAction(400*time.Millisecond, func() error {
		ignoreRelease = true
		return mgr.SwitchTo("menu-main")
	})

	mgr.AddTimedAction(2*time.Second, func() error {
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
			mpdCmdCh <- "toggle"
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
			mpdCmdCh <- "next"
		case DirectionLeft:
			log.Printf("Play prev")
			mpdCmdCh <- "prev"
		}

		return nil
	})

	log.Printf("Waiting for a silent death...")
	<-ctx.Done()

	return mgr.Close()
}
