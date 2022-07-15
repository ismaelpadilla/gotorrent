package cmd

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ismaelpadilla/gotorrent/clients/thepiratebay"
	"github.com/ismaelpadilla/gotorrent/ui"
	"github.com/spf13/cobra"
)

var Debug bool

var rootCmd = &cobra.Command{
	Use:   "gotorrent <query>",
	Short: "gotorrent is a TUI for searching torrents in ThePirateBay",
	Args:  cobra.MinimumNArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		query := strings.Join(args, " ")

		client := thepiratebay.New()

		torrents := client.Search(query)

		p := tea.NewProgram(ui.InitialModel(torrents, Debug),
			tea.WithAltScreen(),
		)

		if err := p.Start(); err != nil {
			fmt.Printf("An error ocurred: %v", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "show debug information")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
