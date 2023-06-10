package ui

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ismaelpadilla/gotorrent/interfaces"
	"github.com/ismaelpadilla/gotorrent/ui/keys"
	"github.com/skratchdot/open-golang/open"
)

var selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))

func InitialModel(query string, config Config) Model {
	var mode Mode

	var torrents []interfaces.Torrent
	h := help.New()
	searchInput := textinput.New()

	if query == "" {
		mode = Search
		h.View(keys.SearchKeys)
		searchInput.Focus()
	} else {
		mode = List
		h.View(keys.ListKeys)
		torrents = config.Client.Search(query)
	}

	return Model{
		client:           config.Client,
		torrents:         torrents,
		downloadLocation: config.DownloadFolder,
		mode:             mode,
		keys:             keys.ListKeys,
		help:             h,
		persist:          config.Persist,
		searchInput:      searchInput,
		visitCommand:     config.VisitCommand,
		debug:            config.Debug,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		m.height = msg.Height
		m.help.Width = msg.Width

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.SetContent(m.GetContent())
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}
	default:
		m.viewport.SetContent(m.GetContent())
	}

	// adjust viewport if cursor position isn't visible
	// -1 because of the header line
	if m.cursorPosition < m.viewport.YOffset-1 {
		m.viewport.LineUp(m.viewport.YOffset - m.cursorPosition - 1)
	}
	if m.cursorPosition > m.viewport.Height+m.viewport.YOffset-2 {
		m.viewport.LineDown(m.cursorPosition - m.viewport.Height - m.viewport.YOffset + 2)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) handleKeyPress(msg tea.KeyMsg) (bool, tea.Cmd) {
	var cmd tea.Cmd
	m.message = ""
	keyString := msg.String()
	switch m.mode {
	case List:
		switch keyString {
		case "ctrl+c", "q", "esc":
			return true, nil

		case "up", "k":
			m.input = ""
			if m.cursorPosition > 0 {
				m.cursorPosition--
			} else {
				m.cursorPosition = 0
			}

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

		case "s":
			cmd = m.enterSearchMode()

		case "d":
			m.showDescription()

		case "f":
			m.showFiles()

		case "c":
			m.copyMagnetLinkToClipBoard()

		case "enter":
			go visitMagnetLink(m.visitCommand, m.torrents[m.cursorPosition])
			if !m.persist {
				return true, nil
			}

		case "t":
			cmd = m.downloadTorrent()

		case "g":
			go m.getCurrentTorrent().Client.NavigateTo(*m.getCurrentTorrent())

		case "?":
			m.toggleHelp()
		}

		m.viewport.SetContent(m.GetContent())
	case ShowDescription, ShowFiles:
		switch keyString {
		case "ctrl+c":
			return true, nil

		case "q", "esc":
			m.keys = keys.ListKeys
			m.mode = List

		case "s":
			cmd = m.enterSearchMode()

		case "d":
			m.showDescription()

		case "f":
			m.showFiles()

		case "c":
			m.copyMagnetLinkToClipBoard()

		case "enter":
			go visitMagnetLink(m.visitCommand, m.torrents[m.cursorPosition])
			if !m.persist {
				return true, nil
			}

		case "t":
			cmd = m.downloadTorrent()

		case "g":
			go m.getCurrentTorrent().Client.NavigateTo(*m.getCurrentTorrent())

		case "?":
			m.toggleHelp()

		case "up", "k":
			m.viewport.LineUp(1)

		case "down", "j":
			m.viewport.LineDown(1)
		}
	case Search:
		switch keyString {
		case "ctrl+c":
			return true, nil

		case "esc":
			if len(m.torrents) > 0 {
				m.keys = keys.ListKeys
				m.mode = List
			} else {
				return true, nil
			}

		case "enter":
			m.cursorPosition = 0
			m.torrents = m.client.Search(m.searchInput.Value())
			m.mode = List
			m.keys = keys.ListKeys
		}
	}
	return false, cmd
}

func (m *Model) headerView() string {
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

func (m *Model) footerView() string {
	info := "\nInput torrent number: "
	info += selectedStyle.Render(m.input) + "\n"
	info += m.message + "\n"

	helpView := m.help.View(m.keys)

	if m.debug {
		info += fmt.Sprintf("CursorPos: %d, Height: %d, Offset: %d\n", m.cursorPosition, m.viewport.Height, m.viewport.YOffset)
	}

	infoCentered := lipgloss.JoinHorizontal(lipgloss.Center, info)
	return infoCentered + helpView
}

func (m *Model) GetContent() string {
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

func (m *Model) GetSearchContent() string {
	return m.searchInput.View()
}

func (m *Model) GetTorrentFilesTable() string {
	nameLength := getMaxFileNameLength(m.getCurrentTorrent().Files)
	nameLenghtAsString := strconv.Itoa(nameLength)

	// table header
	// the Name column with is variable, it is as wide as the longest name
	s := fmt.Sprintf("%3s %"+nameLenghtAsString+"s %9s\n", "No.", "Name", "Size")

	// Iterate over our choices
	for i, file := range m.getCurrentTorrent().Files {
		s += fmt.Sprintf("%3d %"+nameLenghtAsString+"s %9s\n", i, file.Name, file.GetPrettySize())
	}
	return s
}

func (m *Model) GetTorrentsTable() string {
	// table header
	s := fmt.Sprintf("%s %3s %64s %9s %4s %4s %s\n", " ", "No.", "Title", "Size", "S", "L", "Uploaded")

	for i, torrent := range m.torrents {
		dateInt, err := strconv.ParseInt(torrent.Uploaded, 10, 64)
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
			s += selectedStyle.Render(fmt.Sprintf("%s %3d %64s %9s %4d %4d %s", cursor, i, torrent.Title, torrent.GetPrettySize(), torrent.Seeders, torrent.Leechers, date)) + "\n"
		} else {
			s += fmt.Sprintf("%s %3d %64s %9s %4d %4d %s\n", cursor, i, torrent.Title, torrent.GetPrettySize(), torrent.Seeders, torrent.Leechers, date)
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

func (m *Model) copyMagnetLinkToClipBoard() {
	torrent := m.torrents[m.cursorPosition]
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

func visitMagnetLink(visitCommand string, torrent interfaces.Torrent) {
	var err error
	if visitCommand == "" {
		err = open.Run(torrent.MagnetLink)
	} else {

		visitCmd := strings.Fields(visitCommand)

		if !strings.Contains(visitCommand, "%s") {
			visitCmd = append(visitCmd, "%s")
		}

		for i, word := range visitCmd {
			visitCmd[i] = strings.Replace(word, "%s", torrent.MagnetLink, -1)
		}

		cmd := exec.Command(visitCmd[0], visitCmd[1:]...)
		cmd.Stdout = os.Stdout
		err = cmd.Run()
	}
	if err != nil {
		fmt.Println("error")
	}
}

func (m *Model) getCurrentTorrent() *interfaces.Torrent {
	return &m.torrents[m.cursorPosition]
}

func (m *Model) enterSearchMode() tea.Cmd {
	m.searchInput.SetValue("")
	cmd := m.searchInput.Focus()
	m.keys = keys.SearchKeys
	m.mode = Search

	return cmd
}

func (m *Model) showDescription() {
	t := m.getCurrentTorrent()
	if t.Description == "" {
		t.Description = t.FetchDescription()
	}
	m.keys = keys.DescriptionKeys
	m.mode = ShowDescription
}

func (m *Model) showFiles() {
	t := m.getCurrentTorrent()
	if t.Files == nil {
		t.Files = t.FetchFiles()
	}
	m.keys = keys.FilesKeys
	m.mode = ShowFiles
}

func (m *Model) downloadTorrent() tea.Cmd {
	m.message = "Downloading .torrent"
	return cmdDownloadTorrentFile(*m)
}

func (m *Model) toggleHelp() {
	m.help.ShowAll = !m.help.ShowAll

	// adjust viewport, since toggling help changes footer size
	headerHeight := lipgloss.Height(m.headerView())
	footerHeight := lipgloss.Height(m.footerView())
	verticalMarginHeight := headerHeight + footerHeight
	m.viewport.Height = m.height - verticalMarginHeight
}
