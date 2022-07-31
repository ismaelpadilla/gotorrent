package ui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/ismaelpadilla/gotorrent/interfaces"
)

type Mode int

const (
	List Mode = iota
	ShowDescription
	ShowFiles
	Search
)

type Model struct {
	client           interfaces.Client
	torrents         []interfaces.Torrent
	downloadLocation string
	cursorPosition   int
	input            string
	keys             help.KeyMap
	help             help.Model
	viewport         viewport.Model
	height           int
	ready            bool
	mode             Mode
	searchInput      textinput.Model
	message          string
	persist          bool
	debug            bool
}

type Config struct {
	Client         interfaces.Client
	Persist        bool
	DownloadFolder string
	Debug          bool
}

type errMsg struct{ err error }
type statusMsg struct{ message string }
