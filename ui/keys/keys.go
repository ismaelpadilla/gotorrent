package keys

import "github.com/charmbracelet/bubbles/key"

type keys struct {
	Up                key.Binding
	Down              key.Binding
	GetTorrent        key.Binding
	NavigateToTorrent key.Binding
	DownloadTorrent   key.Binding
	CopyMagnetLink    key.Binding
	ShowDescription   key.Binding
	ShowFiles         key.Binding
	GoBackEsc         key.Binding
	GoBackQEsc        key.Binding
	SearchS           key.Binding
	SearchEnter       key.Binding
	Help              key.Binding
	CtrlC             key.Binding
	QuitQEscCtrlC     key.Binding
}

// key definition used in other files
var allKeys = keys{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	GetTorrent: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "get torrent"),
	),
	NavigateToTorrent: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g", "go to torrent"),
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
	GoBackEsc: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "go back"),
	),
	GoBackQEsc: key.NewBinding(
		key.WithKeys("q", "esc"),
		key.WithHelp("q/esc", "go back"),
	),
	SearchS: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "search"),
	),
	SearchEnter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "search"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	CtrlC: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	QuitQEscCtrlC: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
