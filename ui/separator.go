package ui

import (
	"strings"

	"github.com/studentkittens/eulenfunk/util"
)

// Separator is a visual separator between several other menu entries.
// Other than that it has no function.
type Separator struct {
	// Title of the separator (drawn centered)
	Title string
}

// Render draws the separator centered in `w`.
// `active` is ignored.
func (sp *Separator) Render(w int, active bool) string {
	return util.Center(strings.ToUpper(" "+sp.Title+" "), w, '━')
}

// Action is a no-op.
func (sp *Separator) Action() error {
	return nil
}

// Selectable is always false.
func (sp *Separator) Selectable() bool {
	return false
}
