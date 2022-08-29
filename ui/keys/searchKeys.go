package keys

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

var SearchKeys = searchKeyMap{
	Enter:  allKeys.SearchEnter,
	GoBack: allKeys.GoBackEsc,
	Help:   allKeys.Help,
	Quit:   allKeys.CtrlC,
}
