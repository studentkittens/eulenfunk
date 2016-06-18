package ui

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/studentkittens/eulenfunk/display"
	"github.com/studentkittens/eulenfunk/util"
)

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
				if err := active.Click(); err != nil {
					name := active.ActiveEntryName()
					log.Printf("Action for menu ClickEntry `%s` failed: %v", name, err)
				}

				mgr.Display()
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

func (mgr *MenuManager) AddMenu(name string, entries []Entry) error {
	mgr.Lock()
	defer mgr.Unlock()

	menu, err := NewMenu(name, mgr.lw)
	if err != nil {
		return err
	}

	for _, ClickEntry := range entries {
		menu.Entries = append(menu.Entries, ClickEntry)
	}

	if mgr.Active == nil {
		mgr.Active = menu
		mgr.Active.Display(mgr.Config.Width)
	}

	// Pre-select first selectable ClickEntry:
	menu.Scroll(0)

	mgr.Menus[name] = menu
	return nil
}

func (mgr *MenuManager) Close() error {
	return mgr.rotary.Close()
}
