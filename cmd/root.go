package cmd

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ismaelpadilla/gotorrent/clients/thepiratebay"
	"github.com/ismaelpadilla/gotorrent/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Debug bool
var Persist bool
var DownloadFolder string

var rootCmd = &cobra.Command{
	Use:   "gotorrent <query>",
	Short: "gotorrent is a TUI for searching torrents in ThePirateBay",
	Args:  cobra.MinimumNArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		DownloadFolder = viper.GetString("download-folder")

		query := strings.Join(args, " ")

		client := thepiratebay.New()

		torrents := client.Search(query)

		// DownloadLocation represents a folder, it should end with "/"
		if DownloadFolder != "" && !strings.HasSuffix(DownloadFolder, "/") {
			DownloadFolder = DownloadFolder + "/"
		}

		config := ui.Config{
			Client:         client,
			Persist:        Persist,
			DownloadFolder: DownloadFolder,
			Debug:          Debug,
		}

		p := tea.NewProgram(ui.InitialModel(torrents, config),
			tea.WithAltScreen(),
		)

		if err := p.Start(); err != nil {
			fmt.Printf("An error ocurred: %v", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	setFlags()
	loadConfig()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func setFlags() {
	rootCmd.Flags().BoolVarP(&Debug, "debug", "d", false, "show debug information")
	rootCmd.Flags().BoolVarP(&Persist, "persist", "p", false, "keep gotorrent open after selecting torrent")
	rootCmd.Flags().StringVarP(&DownloadFolder, "download-folder", "f", "", "folder where files are downloaded")
}

func loadConfig() {
	err := viper.BindPFlag("download-folder", rootCmd.Flags().Lookup("download-folder"))
	if err != nil {
		panic(err)
	}

	viper.AddConfigPath("$HOME/.config/gotorrent/")
	viper.AddConfigPath(".")

	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file was found but some other error was produced
			panic(err)
		}
	}
}
