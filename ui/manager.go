package ui

import (
	"log"
	"sync"
	"time"

	"github.com/studentkittens/eulenfunk/display"
	"github.com/studentkittens/eulenfunk/util"
)

const (
	// DirectionNone means no movement (initial state)
	DirectionNone = iota

	// DirectionRight means a clockwise rotation.
	DirectionRight

	// DirectionLeft means a counter-clockwise rotation.
	DirectionLeft
)

// MenuManager handles the switching and drawing of several Menus
// and possibly "normal" windows.
type MenuManager struct {
	sync.Mutex
	Config *Config

	// TODO: cleanup
	active               *Menu
	menus                map[string]*Menu
	timedActions         map[time.Duration]Action
	timedActionExecd     map[time.Duration]bool
	lw                   *display.LineWriter
	rotateActions        []Action
	releaseActions       []Action
	currValue, lastValue int
	activeWindow         string
	rotary               *util.Rotary
	ignoreNextRelease    bool
}

func (mgr *MenuManager) callAction(action Action, typ string) {
	// Order is reversed - This is intentional!
	// callAction should be called locked, but the action
	// should be called unlocked.
	mgr.Unlock()
	defer mgr.Lock()

	if err := action(); err != nil {
		log.Printf("%s action failed: %v", typ, err)
	}
}

func (mgr *MenuManager) handleButtonEvent(state bool) {
	mgr.Lock()
	defer mgr.Unlock()

	if mgr.active == nil {
		return
	}

	switch state {
	case false:
		log.Printf("Button released (ignore: %v)", mgr.ignoreNextRelease)
		mgr.timedActionExecd = make(map[time.Duration]bool)

		if !mgr.ignoreNextRelease {
			for _, action := range mgr.releaseActions {
				mgr.callAction(action, "release")
			}

			mgr.display()
		}

		mgr.ignoreNextRelease = false
	case true:
		log.Printf("Button pressed")

		// Check if we're actually in a menu:
		if mgr.ActiveWindow() == mgr.active.Name {
			mgr.callAction(mgr.active.Click, "pressed")
			mgr.display()
		}
	}
}

func (mgr *MenuManager) handlePressedEvent(duration time.Duration) {
	mgr.Lock()
	defer mgr.Unlock()

	log.Printf("Pressed for %s", duration)

	// Find the action with smallest non-negative diff:
	var diff time.Duration
	var action Action
	var actionTime time.Duration

	for after, timedAction := range mgr.timedActions {
		newDiff := duration - after
		if after <= duration && (action == nil || newDiff < diff) {
			diff = duration - after
			action = timedAction
			actionTime = after
		}
	}

	if action != nil && !mgr.timedActionExecd[actionTime] {
		mgr.callAction(action, "timed")
		mgr.timedActionExecd[actionTime] = true
	}
}

func (mgr *MenuManager) handleValueEvent(value int) {
	mgr.Lock()
	defer mgr.Unlock()

	if mgr.active == nil {
		return
	}

	mgr.lastValue = mgr.currValue
	mgr.currValue = value
	name := mgr.active.Name
	diff := mgr.currValue - mgr.lastValue

	mgr.active.Scroll(diff)
	if _, err := mgr.lw.Printf("move %s %d", name, diff); err != nil {
		log.Printf("move failed: %v", err)
	}

	for _, action := range mgr.rotateActions {
		mgr.callAction(action, "rotate")
	}

	mgr.display()
}

// NewMenuManager returns a new MenuManager that sends it's data to `lw` and switches to `initialWin`.
func NewMenuManager(cfg *Config, lw *display.LineWriter, initialWin string) (*MenuManager, error) {
	rty, err := util.NewRotary()
	if err != nil {
		return nil, err
	}

	// Switch to mpd initially:
	if _, err := lw.Printf("switch %s", initialWin); err != nil {
		return nil, err
	}

	mgr := &MenuManager{
		menus:            make(map[string]*Menu),
		timedActions:     make(map[time.Duration]Action),
		timedActionExecd: make(map[time.Duration]bool),
		activeWindow:     initialWin,
		Config:           cfg,
		lw:               lw,
		rotary:           rty,
	}

	go func() {
		for state := range rty.Button {
			mgr.handleButtonEvent(state)
		}
	}()

	go func() {
		for duration := range rty.Pressed {
			mgr.handlePressedEvent(duration)
		}
	}()

	go func() {
		for value := range rty.Value {
			mgr.handleValueEvent(value)
		}
	}()

	return mgr, nil
}

// Display sends the current active menu to the display server.
// (but it does not switch to it!)
func (mgr *MenuManager) Display() {
	mgr.Lock()
	defer mgr.Unlock()

	mgr.display()
}

func (mgr *MenuManager) display() {
	// Just log the error, since the user can't do anything about it:
	if err := mgr.active.Display(mgr.Config.Width); err != nil {
		log.Printf("Failed to display current state: %v", err)
	}
}

// ActiveWindow returns the currently active window shown
// by the display server.
func (mgr *MenuManager) ActiveWindow() string {
	mgr.Lock()
	defer mgr.Unlock()

	return mgr.activeWindow
}

// Direction returns the direction in which the
// rotary button was rotated last. If it was not
// rotated yet, the duration is DirectionNone.
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

// Value returns the current value of the rotary button.
func (mgr *MenuManager) Value() int {
	mgr.Lock()
	defer mgr.Unlock()

	return mgr.currValue
}

// RotateAction registers an action to be called
// when the user rotates the knob.
func (mgr *MenuManager) RotateAction(a Action) {
	mgr.Lock()
	defer mgr.Unlock()

	mgr.rotateActions = append(mgr.rotateActions, a)
}

// ReleaseAction registers an action to be called when
// the rotary button is released.
func (mgr *MenuManager) ReleaseAction(a Action) {
	mgr.Lock()
	defer mgr.Unlock()

	mgr.releaseActions = append(mgr.releaseActions, a)
}

// SwitchTo switches to the menu named `menu`.
// NOTE: use this instead of directly talking to the display
// server since it also switches the input focus.
func (mgr *MenuManager) SwitchTo(name string) error {
	mgr.Lock()
	defer mgr.Unlock()

	if menu, ok := mgr.menus[name]; ok {
		mgr.active = menu
		mgr.display()
	}

	if _, err := mgr.lw.Printf("switch %s", name); err != nil {
		log.Printf("switch failed: %v", err)
		return err
	}

	mgr.ignoreNextRelease = true
	mgr.activeWindow = name
	return nil
}

// AddTimedAction register an action to be called after pressing the rotary
// button with a duration of `after`.
func (mgr *MenuManager) AddTimedAction(after time.Duration, action Action) {
	mgr.Lock()
	defer mgr.Unlock()

	mgr.timedActions[after] = action
}

// AddMenu adds a (potentially new) menu to the manager with `name` and
// `entries`. If the menu already exists it will be subsituted with the new one.
func (mgr *MenuManager) AddMenu(name string, entries []Entry) error {
	mgr.Lock()
	defer mgr.Unlock()

	menu, err := NewMenu(name, mgr.lw)
	if err != nil {
		return err
	}

	menu.Entries = append(menu.Entries, entries...)

	// Why? Because AddMenu may be called more than once with different entries.
	// If first a long menu is given and then a short, the diff will still
	// contain the lines of the longer one:
	if _, err := mgr.lw.Printf("truncate %s %d", name, len(entries)); err != nil {
		log.Printf("Failed to truncate menu %s: %v", name, err)
	}

	if mgr.active == nil {
		mgr.active = menu
		mgr.display()
	}

	// Pre-select first selectable entry:
	menu.Scroll(0)

	mgr.menus[name] = menu
	return nil
}

// Close closes all input resources
func (mgr *MenuManager) Close() error {
	return mgr.rotary.Close()
}
