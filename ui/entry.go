package ui

type Action func() error

type Entry interface {
	// Selectable returns true when the user can select the entry
	Selectable() bool

	// Render draws the line as string (with len of w), possibly with an
	// "active" marker as indicated by `active`.
	Render(w int, active bool) string

	// Name returns the name of the entry
	Name() string

	// Action gets called when the user clicks the button.
	// Errors are logged but not handled otherwise.
	Action() error
}
