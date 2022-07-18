package ui

import "github.com/charmbracelet/bubbles/key"

type listKeyMap struct {
	Up              key.Binding
	Down            key.Binding
	Enter           key.Binding
	DownloadTorrent key.Binding
	CopyMagnetLink  key.Binding
	ShowDescription key.Binding
	ShowFiles       key.Binding
	Search          key.Binding
	Help            key.Binding
	Quit            key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k listKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k listKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Enter, k.Search},                                     // first column
		{k.DownloadTorrent, k.CopyMagnetLink, k.ShowDescription, k.ShowFiles}, // second column
		{k.Help, k.Quit}, // third column
	}
}

var listKeys = listKeyMap{
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
	DownloadTorrent: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "download .torrent"),
	),
	CopyMagnetLink: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "copy magnet link"),
	),
	ShowDescription: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "show description"),
	),
	ShowFiles: key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "show files"),
	),
	Search: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "search"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
