package ui

import (
	"strings"

	"github.com/studentkittens/eulenfunk/util"
)

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
