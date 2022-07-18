package ui

import "github.com/charmbracelet/bubbles/key"

type filesKeyMap struct {
	Up              key.Binding
	Down            key.Binding
	Enter           key.Binding
	CopyMagnetLink  key.Binding
	ShowDescription key.Binding
	GoBack          key.Binding
	Help            key.Binding
	Quit            key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k filesKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k filesKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Enter},                         // first column
		{k.CopyMagnetLink, k.ShowDescription, k.GoBack}, // second column
		{k.GoBack, k.Quit},                              // third column
	}
}

var filesKeys = filesKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "get torrent"),
	),
	CopyMagnetLink: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "copy magnet link"),
	),
	ShowDescription: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "show description"),
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
