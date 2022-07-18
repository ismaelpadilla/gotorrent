package ui

import "github.com/charmbracelet/bubbles/key"

type searchKeyMap struct {
	Enter  key.Binding
	GoBack key.Binding
	Help   key.Binding
	Quit   key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k searchKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Enter, k.GoBack, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k searchKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Enter},  // first column
		{k.GoBack}, // second column
		{k.Quit},   // third column
	}
}

var searchKeys = searchKeyMap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "search"),
	),
	GoBack: key.NewBinding(
		key.WithKeys("q", "esc"),
		key.WithHelp("q", "go back"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}
