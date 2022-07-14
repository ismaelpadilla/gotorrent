package ui

import (
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ismaelpadilla/gotorrent/interfaces"
	"github.com/pkg/browser"
)

var selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))

type model struct {
	torrents       []interfaces.Torrent
	cursorPosition int
	input          string
	keys           keyMap
	help           help.Model
	viewport       viewport.Model
	height         int
	ready          bool
}

func InitialModel(torrents []interfaces.Torrent) model {
	choices := make([]string, len(torrents))
	h := help.New()

	// call h.View to do some initialization that may cause problems if called later
	h.View(keys)

	for i, t := range torrents {
		choices[i] = t.Title
	}
	return model{
		torrents: torrents,
		keys:     keys,
		help:     h,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	useHighPerformanceRenderer := false
	var (
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Which key was pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

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
			m.cursorPosition, _ = strconv.Atoi(m.input)

		case "backspace":
			// we can do this safely because m.input contains numbers only
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
				m.cursorPosition, _ = strconv.Atoi(m.input)
			}

		// Enter navigates to magnet link
		case "enter":
			go browser.OpenURL(m.torrents[m.cursorPosition].MagnetLink)

		case m.keys.Help.Help().Key:
			m.help.ShowAll = !m.help.ShowAll

			// adjust viewport, since toggling help changes footer size
			headerHeight := lipgloss.Height(m.headerView())
			footerHeight := lipgloss.Height(m.footerView())
			verticalMarginHeight := headerHeight + footerHeight
			m.viewport.Height = m.height - verticalMarginHeight
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

func (m model) headerView() string {
	title := "Select torrent to get, or input number and press enter\n"
	return title
}

func (m model) footerView() string {
	info := "\nInput torrent number: "
	info += selectedStyle.Render(m.input) + "\n"

	helpView := m.help.View(m.keys)

	// debug info (TODO: only display if debug flag is used)
	info += fmt.Sprintf("\nCursorPos: %d, Height: %d, Offset: %d\n", m.cursorPosition, m.viewport.Height, m.viewport.YOffset)

	infoCentered := lipgloss.JoinHorizontal(lipgloss.Center, info)

	return infoCentered + helpView
}

func (m model) GetContent() string {
	// table header
	s := fmt.Sprintf("%s %3s %64s %9s %4s %4s %s\n", " ", "No.", "Title", "Size", "S", "L", "Uploaded")

	// Iterate over our choices
	for i, choice := range m.torrents {
		dateInt, _ := strconv.ParseInt(choice.Uploaded, 10, 64)
		date := time.Unix(dateInt, 0).Format("2006-01-02")

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

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}
