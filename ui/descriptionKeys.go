package ui

import "github.com/charmbracelet/bubbles/key"

type descriptionKeyMap struct {
	Up        key.Binding
	Down      key.Binding
	Enter     key.Binding
	ShowFiles key.Binding
	GoBack    key.Binding
	Help      key.Binding
	Quit      key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k descriptionKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k descriptionKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Enter}, // first column
		{k.ShowFiles, k.GoBack}, // second column
		{k.GoBack, k.Quit},      // third column
	}
}

var descriptionKeys = descriptionKeyMap{
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
	ShowFiles: key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "show files"),
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
