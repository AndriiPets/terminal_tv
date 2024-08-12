package ui

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Quit  key.Binding
	Pause key.Binding
	Enter key.Binding
	Up    key.Binding
	Down  key.Binding
}
