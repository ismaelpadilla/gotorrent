package keys

import "github.com/charmbracelet/bubbles/key"

type listKeyMap struct {
	Up                key.Binding
	Down              key.Binding
	Enter             key.Binding
	DownloadTorrent   key.Binding
	NavigateToTorrent key.Binding
	CopyMagnetLink    key.Binding
	ShowDescription   key.Binding
	ShowFiles         key.Binding
	Search            key.Binding
	Help              key.Binding
	Quit              key.Binding
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
		{k.Up, k.Down, k.Enter, k.NavigateToTorrent},                          // first column
		{k.DownloadTorrent, k.CopyMagnetLink, k.ShowDescription, k.ShowFiles}, // second column
		{k.Search, k.Help, k.Quit},                                            // third column
	}
}

var ListKeys = listKeyMap{
	Up:                allKeys.Up,
	Down:              allKeys.Down,
	Enter:             allKeys.GetTorrent,
	NavigateToTorrent: allKeys.NavigateToTorrent,
	DownloadTorrent:   allKeys.DownloadTorrent,
	CopyMagnetLink:    allKeys.CopyMagnetLink,
	ShowDescription:   allKeys.ShowDescription,
	ShowFiles:         allKeys.ShowFiles,
	Search:            allKeys.SearchS,
	Help:              allKeys.Help,
	Quit:              allKeys.QuitQEscCtrlC,
}
