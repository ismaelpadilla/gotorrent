package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ismaelpadilla/gotorrent/clients/thepiratebay"
	"github.com/ismaelpadilla/gotorrent/ui"
)

func main() {
	args := os.Args
	var query string
	if len(args) < 2 {
		fmt.Println("Usage: gotorrent <query>")
		return
	}

	query = strings.Join(args[1:], " ")

	client := thepiratebay.New()

	torrents := client.Search(query)

	p := tea.NewProgram(ui.InitialModel(torrents))

	if err := p.Start(); err != nil {
		fmt.Printf("An error ocurred: %v", err)
		os.Exit(1)
	}
}
