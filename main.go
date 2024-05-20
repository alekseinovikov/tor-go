package main

import (
	"github.com/anacrolix/torrent"
	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func main() {
	torrentClient, _ := torrent.NewClient(nil)

	fp := filepicker.New()
	fp.AllowedTypes = []string{".torrent"}
	fp.CurrentDirectory, _ = os.UserHomeDir()

	m := uiModel{
		filepicker:    fp,
		torrentClient: torrentClient,
	}
	_, _ = tea.NewProgram(&m).Run()
}
