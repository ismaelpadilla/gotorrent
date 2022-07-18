package ui

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ismaelpadilla/gotorrent/interfaces"
	"github.com/pkg/browser"
)

var selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))

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

func InitialModel(torrents []interfaces.Torrent, config Config) Model {
	choices := make([]string, len(torrents))
	h := help.New()

	// call h.View to do some initialization that may cause problems if called later
	h.View(listKeys)

	searchInput := textinput.New()

	for i, t := range torrents {
		choices[i] = t.Title
	}
	return Model{
		client:           config.Client,
		torrents:         torrents,
		downloadLocation: config.DownloadFolder,
		mode:             List,
		keys:             listKeys,
		help:             h,
		persist:          config.Persist,
		searchInput:      searchInput,
		debug:            config.Debug,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) handleKeyPress(msg tea.KeyMsg) (bool, tea.Cmd) {
	var cmd tea.Cmd
	m.message = ""
	keyString := msg.String()
	switch m.mode {
	case List:
		// Which key was pressed?
		switch keyString {
		// These keys should exit the program.
		case "ctrl+c", "q", "esc":
			return true, nil

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			m.input = ""
			if m.cursorPosition > 0 {
				m.cursorPosition--
			} else {
				m.cursorPosition = 0
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			m.input = ""
			if m.cursorPosition < len(m.torrents)-1 {
				m.cursorPosition++
			} else {
				m.cursorPosition = len(m.torrents) - 1
			}

		// numbers and backspace change input number
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			m.input += msg.String()
			inputNumber, _ := strconv.Atoi(m.input)
			if inputNumber < len(m.torrents)-1 {
				m.cursorPosition = inputNumber
			}

		case "backspace":
			// we can do this safely because m.input contains numbers only
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
				m.cursorPosition, _ = strconv.Atoi(m.input)
			}

		// Search for a new torrent
		case "s":
			m.searchInput.SetValue("")
			m.searchInput.Focus()
			m.keys = searchKeys
			m.mode = Search

		// Show description
		case "d":
			t := m.getCurrentTorrent()
			if t.Description == "" {
				t.Description = t.FetchDescription()
			}
			m.keys = descriptionKeys
			m.mode = ShowDescription

		// Show files
		case "f":
			t := m.getCurrentTorrent()
			if t.Files == nil {
				t.Files = t.FetchFiles()
			}
			m.keys = filesKeys
			m.mode = ShowFiles

		// Copy magnet link to clipboard
		case "c":
			m.copyMagnetLinkToClipBoard(m.torrents[m.cursorPosition])

		// Enter navigates to magnet link
		case "enter":
			go visitMagnetLink(m.torrents[m.cursorPosition])
			if !m.persist {
				return true, nil
			}

		// Download torrent file
		case "t":
			m.message = "Downloading .torrent"
			cmd = cmdDownloadTorrentFile(*m)

		case "?":
			m.help.ShowAll = !m.help.ShowAll

			// adjust viewport, since toggling help changes footer size
			headerHeight := lipgloss.Height(m.headerView())
			footerHeight := lipgloss.Height(m.footerView())
			verticalMarginHeight := headerHeight + footerHeight
			m.viewport.Height = m.height - verticalMarginHeight
		}

		m.viewport.SetContent(m.GetContent())
	case ShowDescription, ShowFiles:
		switch keyString {
		case "ctrl+c":
			return true, nil

		case "q", "esc":
			m.keys = listKeys
			m.mode = List

		// Search for a new torrent
		case "s":
			m.searchInput.SetValue("")
			m.searchInput.Focus()
			m.keys = searchKeys
			m.mode = Search

		// Show description
		case "d":
			t := m.getCurrentTorrent()
			if t.Description == "" {
				t.Description = t.FetchDescription()
			}
			m.keys = descriptionKeys
			m.mode = ShowDescription

		// Show files
		case "f":
			t := m.getCurrentTorrent()
			if t.Files == nil {
				t.Files = t.FetchFiles()
			}
			m.keys = filesKeys
			m.mode = ShowFiles

		// Copy magnet link to clipboard
		case "c":
			m.copyMagnetLinkToClipBoard(m.torrents[m.cursorPosition])

		case "enter":
			go visitMagnetLink(m.torrents[m.cursorPosition])
			if !m.persist {
				return true, nil
			}

		// Download torrent file
		case "t":
			m.message = "Downloading .torrent"
			cmd = cmdDownloadTorrentFile(*m)

		case "?":
			m.help.ShowAll = !m.help.ShowAll

			// adjust viewport, since toggling help changes footer size
			headerHeight := lipgloss.Height(m.headerView())
			footerHeight := lipgloss.Height(m.footerView())
			verticalMarginHeight := headerHeight + footerHeight
			m.viewport.Height = m.height - verticalMarginHeight

		// The "up" and "k" keys scroll up
		case "up", "k":
			m.viewport.LineUp(1)

		// The "down" and "j" keys scroll down
		case "down", "j":
			m.viewport.LineDown(1)
		}
	case Search:
		switch keyString {
		case "ctrl+c":
			return true, nil

		case "esc":
			m.keys = listKeys
			m.mode = List

		case "enter":
			m.torrents = m.client.Search(m.searchInput.Value())
			m.mode = List
			m.keys = listKeys
		}
	}
	return false, cmd
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	useHighPerformanceRenderer := false
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	m.searchInput, cmd = m.searchInput.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case statusMsg:
		m.message = msg.message
	case errMsg:
		m.message = msg.err.Error()
	// Is it a key press?
	case tea.KeyMsg:
		shouldQuit, cmd := m.handleKeyPress(msg)
		if shouldQuit {
			return m, tea.Quit
		}
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

		m.viewport.SetContent(m.GetContent())

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.height = msg.Height
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			m.help.Width = msg.Width
			m.viewport.SetContent(m.GetContent())
			m.ready = true

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			// m.viewport.YPosition = headerHeight
		} else {
			m.height = msg.Height
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}
		// m.help.Width = msg.Width

		if useHighPerformanceRenderer {
			// Render (or re-render) the whole viewport. Necessary both to
			// initialize the viewport and when the window is resized.
			//
			// This is needed for high-performance rendering only.
			cmds = append(cmds, viewport.Sync(m.viewport))
		}
	}

	// adjust viewport if cursor position isn't visible
	// -1 because of the header line
	if m.cursorPosition < m.viewport.YOffset-1 {
		m.viewport.LineUp(m.viewport.YOffset - m.cursorPosition - 1)
	}
	if m.cursorPosition > m.viewport.Height+m.viewport.YOffset-2 {
		m.viewport.LineDown(m.cursorPosition - m.viewport.Height - m.viewport.YOffset + 2)
	}
	// vp, cmd := m.viewport.Update(msg)
	return m, tea.Batch(cmds...)
}

func (m Model) headerView() string {
	var title string
	switch m.mode {
	case List:
		title = "Select torrent to get, or input number and press enter\n"
	case ShowDescription:
		title = m.getCurrentTorrent().Title + "\n"
	case ShowFiles:
		title = m.getCurrentTorrent().Title + " files\n"
	case Search:
		title = "Enter query and press enter to search, or press esc to go back\n"
	}

	return title
}

func (m Model) footerView() string {
	info := "\nInput torrent number: "
	info += selectedStyle.Render(m.input) + "\n"
	info += m.message + "\n"

	helpView := m.help.View(m.keys)

	// debug info
	if m.debug {
		info += fmt.Sprintf("CursorPos: %d, Height: %d, Offset: %d\n", m.cursorPosition, m.viewport.Height, m.viewport.YOffset)
	}

	infoCentered := lipgloss.JoinHorizontal(lipgloss.Center, info)
	return infoCentered + helpView
}

func (m Model) GetContent() string {
	switch m.mode {
	case ShowDescription:
		return m.getCurrentTorrent().Description
	case ShowFiles:
		return m.GetTorrentFilesTable()
	case Search:
		return m.GetSearchContent()
	default:
		return m.GetTorrentsTable()
	}
}

func (m Model) GetSearchContent() string {
	return m.searchInput.View()
}

func (m Model) GetTorrentFilesTable() string {
	nameLength := getMaxFileNameLength(m.getCurrentTorrent().Files)
	nameLenghtAsString := strconv.Itoa(nameLength)
	// table header
	// the Name column with is variable, it is as wide as the longest name
	s := fmt.Sprintf("%3s %"+nameLenghtAsString+"s %9s\n", "No.", "Name", "Size")

	// Iterate over our choices
	for i, choice := range m.getCurrentTorrent().Files {
		s += fmt.Sprintf("%3d %"+nameLenghtAsString+"s %9s\n", i, choice.Name, choice.GetPrettySize())
	}
	return s
}

func (m Model) GetTorrentsTable() string {
	// table header
	s := fmt.Sprintf("%s %3s %64s %9s %4s %4s %s\n", " ", "No.", "Title", "Size", "S", "L", "Uploaded")

	// Iterate over our choices
	for i, choice := range m.torrents {
		dateInt, err := strconv.ParseInt(choice.Uploaded, 10, 64)
		var date string
		if err != nil {
			date = "err"
		} else {
			date = time.Unix(dateInt, 0).Format("2006-01-02")
		}

		// Is the cursor pointing at this choice?
		cursor := " "
		if m.cursorPosition == i {
			cursor = ">"
			s += selectedStyle.Render(fmt.Sprintf("%s %3d %64s %9s %4d %4d %s", cursor, i, choice.Title, choice.GetPrettySize(), choice.Seeders, choice.Leechers, date)) + "\n"
		} else {
			s += fmt.Sprintf("%s %3d %64s %9s %4d %4d %s\n", cursor, i, choice.Title, choice.GetPrettySize(), choice.Seeders, choice.Leechers, date)
		}
	}
	return s
}

func getMaxFileNameLength(torrentFiles []interfaces.TorrentFile) int {
	maxLength := 0
	for _, tf := range torrentFiles {
		length := utf8.RuneCountInString(tf.Name)
		if length > maxLength {
			maxLength = length
		}
	}
	return maxLength
}

func (m Model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func (m *Model) copyMagnetLinkToClipBoard(torrent interfaces.Torrent) {
	if err := clipboard.WriteAll(torrent.MagnetLink); err != nil {
		m.message = "Error while copying magnet link to clipboard"
	} else {
		m.message = "Magnet link copied to clipboard"
	}
}

func cmdDownloadTorrentFile(m Model) tea.Cmd {
	return func() tea.Msg {
		url := fmt.Sprintf("http://itorrents.org/torrent/%s.torrent", m.getCurrentTorrent().InfoHash)
		result, err := http.Get(url)
		if err != nil {
			return errMsg{err}
		}
		defer result.Body.Close()

		fileName := fmt.Sprintf("%s%s.torrent", m.downloadLocation, m.getCurrentTorrent().Title)
		file, err := os.Create(fileName)
		if err != nil {
			return errMsg{err}
		}
		defer file.Close()
		_, err = io.Copy(file, result.Body)
		if err != nil {
			return errMsg{err}
		}
		return statusMsg{"Downloaded file: " + fileName}
	}
}

func visitMagnetLink(torrent interfaces.Torrent) {
	err := browser.OpenURL(torrent.MagnetLink)
	if err != nil {
		fmt.Println("error")
	}
}

func (m *Model) getCurrentTorrent() *interfaces.Torrent {
	return &m.torrents[m.cursorPosition]
}
