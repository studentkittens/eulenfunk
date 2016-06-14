package menu

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"

	"github.com/studentkittens/eulenfunk/display"
	"github.com/studentkittens/eulenfunk/mpd/mpdinfo"
	"github.com/studentkittens/eulenfunk/util"
)

type Action func() error

type Entry struct {
	Text   string
	Action Action
}

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
		line := entry.Text

		if pos == mn.Cursor {
			line = "> " + line
		} else {
			line = "  " + line
		}

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
	rotary               *util.Rotary
}

func NewMenuManager(lw *display.LineWriter) (*MenuManager, error) {
	rty, err := util.NewRotary()
	if err != nil {
		return nil, err
	}

	mgr := &MenuManager{
		Menus:        make(map[string]*Menu),
		TimedActions: make(map[time.Duration]Action),
		lw:           lw,
		rotary:       rty,
	}

	go func() {
		for state := range rty.Button {
			if mgr.Active == nil {
				continue
			}

			if !state {
				fmt.Println("Button released")
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

			for after, timedAction := range mgr.TimedActions {
				newDiff := duration - after
				if after <= duration && (action == nil || newDiff < diff) {
					diff = duration - after
					action = timedAction
				}
			}

			mgr.Unlock()

			if action != nil {
				action()
			}
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

func Run(ctx context.Context) error {
	// TODO: pass config
	cfg := &display.Config{
		Host: "localhost",
		Port: 7778,
	}

	lw, err := display.Connect(cfg)
	if err != nil {
		return err
	}

	defer lw.Close()

	// Switch to mpd initially:
	if _, err := lw.Formatf("switch mpd"); err != nil {
		return err
	}

	msg := util.Center("... booting ...", 20) // TODO
	if _, err := lw.Formatf("line mpd 2 %s", msg); err != nil {
		return err
	}

	mgr, err := NewMenuManager(lw)
	if err != nil {
		return err
	}

	// Some flags to coordinate actions:
	togglePlayback := false
	currentWindow := "mpd"

	switcher := func(mgr *MenuManager, lw *display.LineWriter, name string) func() error {
		return func() error {
			currentWindow = name
			return mgr.SwitchTo(name)
		}
	}

	// Start auxillary services:
	go mpdinfo.Run(&mpdinfo.Config{
		Host:        "localhost",
		Port:        6600,
		DisplayHost: "localhost",
		DisplayPort: 7778,
	}, ctx)

	go RunClock(lw, 20, ctx) // TODO: get width?
	go RunSysinfo(lw, 20, ctx)

	mainMenu := []*Entry{
		{
			"Show status", switcher(mgr, lw, "mpd"),
		}, {
			"Playlists", switcher(mgr, lw, "playlists"),
		}, {
			"Toggle PartyMode", nil, // TODO
		}, {
			"System info", switcher(mgr, lw, "sysinfo"),
		}, {
			"Clock", switcher(mgr, lw, "clock"),
		}, {
			"Switch Mono/Stereo", nil, // TODO
		}, {
			"Stop playback", nil, // TODO
		}, {
			"Power", switcher(mgr, lw, "menu-power"),
		},
	}

	powerMenu := []*Entry{
		{
			"Poweroff", nil, // TODO
		}, {
			"Reboot", nil, // TODO
		}, {
			"Exit", switcher(mgr, lw, "menu-main"),
		},
	}

	easterEggMenu := []*Entry{
		{
			"Schuhu?", nil, // TODO: Play actual shuhu.
		}, {
			"Exit", switcher(mgr, lw, "menu-main"),
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

	if err := mgr.AddMenu("menu-easteregg", easterEggMenu); err != nil {
		log.Printf("Add main-easteregg failed: %v", err)
		return err
	}

	mgr.AddTimedAction(10*time.Millisecond, func() error {
		togglePlayback = true
		return nil
	})

	mgr.AddTimedAction(500*time.Millisecond, func() error {
		togglePlayback = false
		return mgr.SwitchTo("menu-main")
	})

	mgr.AddTimedAction(2*time.Second, func() error {
		togglePlayback = false
		return mgr.SwitchTo("menu-power")
	})

	mgr.AddTimedAction(10*time.Second, func() error {
		togglePlayback = false
		return mgr.SwitchTo("menu-easteregg")
	})

	mgr.ReleaseAction(func() error {
		switch currentWindow {
		case "mpd":
			if togglePlayback {
				log.Printf("TOGGLE PLAYBACK!!!")
				togglePlayback = false
			}
		default:
			// This is a bit of a hack:
			// Enable "click to exit window" on most non-menu windows:
			if !strings.Contains(currentWindow, "menu") {
				return switcher(mgr, lw, "main-menu")()
			}
		}

		return nil
	})

	mgr.RotateAction(func() error {
		// TODO: check if in default
		log.Printf("rotate action")
		switch mgr.Direction() {
		case DirectionRight:
			log.Printf("Play next")
		case DirectionLeft:
			log.Printf("Play prev")
		}

		return nil
	})

	<-ctx.Done()

	return mgr.Close()
}
