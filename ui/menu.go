package menu

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"

	"github.com/studentkittens/eulenfunk/display"
	"github.com/studentkittens/eulenfunk/ui/mpdinfo"
	"github.com/studentkittens/eulenfunk/util"
)

type Entry struct {
	Text   string
	Action Action
	State  string
}

type Action func() error

type Menu struct {
	Name    string
	Entries []*Entry
	Cursor  int

	lw *display.LineWriter
}

func NewMenu(name string, lw *display.LineWriter) (*Menu, error) {
	return &Menu{
		Name: name,
		lw:   lw,
	}, nil
}

func (mn *Menu) AddEntry(entry *Entry) {
	mn.Entries = append(mn.Entries, entry)
}

func (mn *Menu) ActiveName() string {
	if len(mn.Entries) == 0 {
		return ""
	}

	return mn.Entries[mn.Cursor].Text
}

func (mn *Menu) Scroll(move int) {
	mn.Cursor += move
	if mn.Cursor < 0 {
		mn.Cursor = 0
	}

	if mn.Cursor >= len(mn.Entries) {
		mn.Cursor = len(mn.Entries) - 1
	}
}

func (mn *Menu) Display() error {
	for pos, entry := range mn.Entries {
		prefix := "  "
		if pos == mn.Cursor {
			prefix = "> "
		}

		state := entry.State
		if state != "" {
			state = " [" + state + "]"
		}

		// TODO: get correct width.
		midLen := 20 - len(prefix) - len(state)
		line := fmt.Sprintf("%s%-*s%s", prefix, midLen, entry.Text, entry.State)

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

			if err := active.Click(); err != nil {
				name := active.ActiveName()
				log.Printf("Action for menu entry `%s` failed: %v", name, err)
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

func (mgr *MenuManager) AddMenu(name string, entries []*Entry) error {
	mgr.Lock()
	defer mgr.Unlock()

	menu, err := NewMenu(name, mgr.lw)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		menu.AddEntry(entry)
	}

	if mgr.Active == nil {
		mgr.Active = menu
		mgr.Active.Display()
	}

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
		Text:  "PartyMode",
		State: "On",
	}

	// partyModeEntry.Action = func() error {
	// 	switch partyModeEntry.State {
	// 	case "On":
	// 		partyModeEntry.State == "Off"
	// 	case
	// 	}
	// 	return nil
	// }

	mainMenu := []*Entry{
		{
			Text:   "Show status",
			Action: switcher("mpd"),
		}, {
			Text:   "Playlists",
			Action: switcher("playlists"),
		}, partyModeEntry,
		{
			Text:   "Clock",
			Action: switcher("clock"),
		}, {
			Text:   "System info",
			Action: switcher("sysinfo"),
		}, {
			Text:   "Statistics",
			Action: switcher("stats"),
		}, {
			Text:   "Switch Mono/Stereo",
			Action: nil, // TODO
		}, {
			Text:   "Toggle stop",
			Action: mpdCommand("stop", mpdCmdCh),
			State:  "play",
		}, {
			Text:   "Power",
			Action: switcher("menu-power"),
		},
	}

	powerMenu := []*Entry{
		{
			Text:   "Poweroff",
			Action: sysCommand("systemctl", "poweroff"),
		}, {
			Text:   "Reboot",
			Action: sysCommand("systemctl", "reboot"),
		}, {
			Text:   "Exit",
			Action: switcher("menu-main"),
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
