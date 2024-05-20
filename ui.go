package main

import (
	"errors"
	"github.com/anacrolix/torrent"
	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
	"time"
)

type uiModel struct {
	filepicker    filepicker.Model
	torrentClient *torrent.Client
	quitting      bool
	selected      bool
	err           error
}

type clearErrorMsg struct{}
type allDownloaded struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func (m *uiModel) downloadTorrent(torrentFilePath string) tea.Cmd {
	return func() tea.Msg {
		t, err := m.torrentClient.AddTorrentFromFile(torrentFilePath)
		if err != nil {
			return err
		}

		<-t.GotInfo()
		t.DownloadAll()
		// processed := t.Info().NumPieces() / t.Stats().PiecesComplete
		m.torrentClient.WaitAll()
		return allDownloaded{}
	}
}

func (m uiModel) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m uiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			defer m.torrentClient.Close()
			return m, tea.Quit
		}
	case clearErrorMsg:
		m.err = nil
	case allDownloaded:
		defer m.torrentClient.Close()
		m.quitting = true
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)

	// Did the user select a file?
	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		m.selected = true
		return m, m.downloadTorrent(path)
	}

	// Did the user select a disabled file?
	// This is only necessary to display an error to the user.
	if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
		// Let's clear the selectedFile and display an error.
		m.selected = false
		m.err = errors.New(path + " is not valid.")
		return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}

	return m, cmd
}

func (m uiModel) View() string {
	if m.quitting {
		return "Finished... Quiting..."
	}

	var s strings.Builder
	s.WriteString("\n  ")

	if m.err != nil {
		s.WriteString(m.filepicker.Styles.DisabledFile.Render(m.err.Error()))
	} else if !m.selected {
		s.WriteString("Pick a file:")
	} else if m.selected {
		s.WriteString("Selected file is downloading... ")
	} else {
		s.WriteString("\n\n" + m.filepicker.View() + "\n")
	}

	return s.String()
}
