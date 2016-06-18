package ui

type ClickEntry struct {
	Text       string
	ActionFunc Action
}

func (en *ClickEntry) Render(w int, active bool) string {
	prefix := "  "
	if active {
		prefix = "❤ "
	}

	return prefix + en.Text
}

func (en *ClickEntry) Name() string {
	return en.Text
}

func (en *ClickEntry) Action() error {
	if en.ActionFunc != nil {
		return en.ActionFunc()
	}

	return nil
}

func (en *ClickEntry) Selectable() bool {
	return true
}
