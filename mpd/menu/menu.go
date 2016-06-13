package menu

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/studentkittens/eulenfunk/display"
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
				// We don't do anything yet...
			fmt.Println("Button released")
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

func switcher(mgr *MenuManager, lw *display.LineWriter, name string) func() error {
	return func() error {
		return mgr.SwitchTo(name)
	}
}

func Run() error {
	cfg := &display.Config{
		Host: "localhost",
		Port: 7778,
	}

	lw, err := display.Connect(cfg)
	if err != nil {
		return err
	}

	defer lw.Close()

	mgr, err := NewMenuManager(lw)
	if err != nil {
		return err
	}

	// Start clock and sysinfo screen:
	killClock, killSysinfo := make(chan bool), make(chan bool)
	go RunClock(lw, 20, killClock) // TODO: get width?
	go RunSysinfo(lw, 20, killSysinfo)

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
			"Schuhu?", nil,
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
		log.Printf("TODO: Toggle playback")
		return nil
	})

	mgr.AddTimedAction(500*time.Millisecond, func() error {
		return mgr.SwitchTo("menu-main")
	})

	mgr.AddTimedAction(2*time.Second, func() error {
		return mgr.SwitchTo("menu-power")
	})

	mgr.AddTimedAction(10*time.Second, func() error {
		return mgr.SwitchTo("menu-easteregg")
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

	log.Printf("Press CTRL-C to shut down")
	ctrlCh := make(chan os.Signal, 1)
	signal.Notify(ctrlCh, os.Interrupt)
	<-ctrlCh

	killClock <- true
	killSysinfo <- true

	return mgr.Close()
}
