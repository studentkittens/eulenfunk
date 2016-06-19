package ui

import (
	"fmt"
	"log"
	"sync"
	"unicode/utf8"
)

// ToggleEntry skips through an ordered list of states.
type ToggleEntry struct {
	// Text is the text displayed along the state
	Text string

	// Order in which the states are toggled through
	Order []string

	// Actions is a map of state to Actions.
	// When changing to a state the respective
	// Action is called.
	Actions map[string]Action

	// state is the currently selected state as
	// index in Order.
	state int

	// For safety:
	mu sync.Mutex
}

// Render draws an activity marker, the text and the current state.
func (te *ToggleEntry) Render(w int, active bool) string {
	te.mu.Lock()
	defer te.mu.Unlock()

	prefix := "  "
	if active {
		prefix = "❤ "
	}

	stateText := "[" + te.Order[te.state] + "]"
	m := w - utf8.RuneCountInString(stateText) - utf8.RuneCountInString(prefix)
	return fmt.Sprintf("%s%-*s%s", prefix, m, te.Text, stateText)
}

// Action goes to the next state and executes it's associated action
func (te *ToggleEntry) Action() error {
	te.mu.Lock()

	te.state = (te.state + 1) % len(te.Order)

	name := te.Order[te.state]
	fn, ok := te.Actions[name]
	if !ok {
		return fmt.Errorf("No such action: %s", name)
	}

	te.mu.Unlock()

	return fn()
}

// Selectable is true when there is a non-zero amount of actions
func (te *ToggleEntry) Selectable() bool {
	return len(te.Actions) > 0
}

// SetState sets the current state but does not call the associated action It's
// useful to set the state if it was changed by extern resources (i.e. a
// different mpd client)
func (te *ToggleEntry) SetState(state string) {
	if _, ok := te.Actions[state]; !ok {
		// Programmer error; a log line is enough:
		log.Printf("No such action `%s`", state)
		return
	}

	te.mu.Lock()
	for idx, toggle := range te.Order {
		if toggle == state {
			te.state = idx
			break
		}
	}
	te.mu.Unlock()

	return
}
