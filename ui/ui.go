package ui

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ismaelpadilla/gotorrent/interfaces"
	"github.com/pkg/browser"
)

type model struct {
	torrents       []interfaces.Torrent
	cursorPosition int
	input          string
	viewport       viewport.Model
	ready          bool
}

func InitialModel(torrents []interfaces.Torrent) model {
	choices := make([]string, len(torrents))
	for i, t := range torrents {
		choices[i] = t.Title
	}
	return model{
		torrents: torrents,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	useHighPerformanceRenderer := false
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Which key was pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
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
			browser.OpenURL(m.torrents[m.cursorPosition].MagnetLink)

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
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			m.viewport.SetContent(m.GetContent())
			m.ready = true

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			m.viewport.YPosition = headerHeight + 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}
		// println(m.viewport.Height)

		// if useHighPerformanceRenderer {
		// Render (or re-render) the whole viewport. Necessary both to
		// initialize the viewport and when the window is resized.
		//
		// This is needed for high-performance rendering only.
		// cmds = append(cmds, viewport.Sync(m.viewport))
		// }
	}

	// adjust viewport if cursor position isn't visible
	if m.cursorPosition < m.viewport.YOffset {
		m.viewport.LineUp(m.viewport.YOffset - m.cursorPosition)
	}
	if m.cursorPosition > m.viewport.Height+m.viewport.YOffset-1 {
		m.viewport.LineDown(m.cursorPosition - m.viewport.Height - m.viewport.YOffset + 1)
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) headerView() string {
	title := "Select torrent to get, or input number and press enter\n"
	return title
}

func (m model) footerView() string {
	info := fmt.Sprintf("\nInput torrent number: %s\n", m.input)

	// debug info (TODO: only display if debug flag is used)
	info += fmt.Sprintf("\nCursorPos: %d, Height: %d, Offset: %d\n", m.cursorPosition, m.viewport.Height, m.viewport.YOffset)

	return lipgloss.JoinHorizontal(lipgloss.Center, info)
}

func (m model) GetContent() string {
	s := ""

	// Iterate over our choices
	for i, choice := range m.torrents {

		// Is the cursor pointing at this choice?
		cursor := " "
		if m.cursorPosition == i {
			cursor = ">"
		}

		// Render the row
		s += fmt.Sprintf("%s %3d %s\n", cursor, i, choice.Title)
	}
	return s
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}
