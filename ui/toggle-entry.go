package ui

import (
	"fmt"
	"sync"
	"unicode/utf8"
)

type ToggleEntry struct {
	// Text is the text displayed along the state
	Text string

	// State is the currently selected state as
	// index in Order.
	State int

	// Order in which the states are toggled through
	Order []string

	// Actions is a map of state to Actions.
	// When changing to a state the respective
	// Action is called.
	Actions map[string]Action

	// For safety:
	mu sync.Mutex
}

func (te *ToggleEntry) Render(w int, active bool) string {
	te.mu.Lock()
	defer te.mu.Unlock()

	prefix := "  "
	if active {
		prefix = "❤ "
	}

	stateText := "[" + te.Order[te.State] + "]"
	m := w - utf8.RuneCountInString(stateText) - utf8.RuneCountInString(prefix)
	return fmt.Sprintf("%s%-*s%s", prefix, m, te.Text, stateText)
}

func (te *ToggleEntry) Name() string {
	return te.Text
}

func (te *ToggleEntry) Action() error {
	te.mu.Lock()
	te.State = (te.State + 1) % len(te.Order)
	te.mu.Unlock()

	return te.SetState(te.Order[te.State], true)
}

func (te *ToggleEntry) Selectable() bool {
	return len(te.Actions) > 0
}

func (te *ToggleEntry) SetState(state string, click bool) error {
	fn, ok := te.Actions[state]
	if !ok {
		return fmt.Errorf("No such action `%s`", state)
	}

	if click {
		return fn()
	}

	return nil
}
