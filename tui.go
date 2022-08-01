package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"

	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pkg/sftp"
)

var (
	docStyle           = lipgloss.NewStyle().Margin(1, 2)
	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render

	fileItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"}).
			Render
	dirItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#0477b5", Dark: "#0477b5"}).
			Render
)

type item struct {
	title       string
	description string
	rawValue    fs.FileInfo
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

type model struct {
	list       list.Model
	sftpClient *sftp.Client
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			var cmd tea.Cmd
			selectedItem := m.list.SelectedItem().(*item).rawValue
			selectedItemName := selectedItem.Name()
			if selectedItem.IsDir() {
				cmd = m.list.NewStatusMessage(statusMessageStyle(fmt.Sprintf("This is a dir %s", selectedItemName)))
			} else {
				cmd = m.list.NewStatusMessage(statusMessageStyle(fmt.Sprintf("Downloading %s", selectedItemName)))
				err := m.downloadFile(selectedItemName)
				handleError(err)
			}

			return m, cmd
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) downloadFile(fileToDownload string) error {
	srcFile, err := m.sftpClient.Open(fileToDownload)
	handleError(err)
	defer srcFile.Close()
	destFile, err := os.Create(filepath.Join(".", fileToDownload))
	defer destFile.Close()
	handleError(err)
	_, err = io.Copy(destFile, srcFile)
	return err
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func createItemList(sftpClient *sftp.Client) []list.Item {
	fileList, err := sftpClient.ReadDir(".")
	handleError(err)

	items := []list.Item{}

	for _, value := range fileList {
		var decoratedItem string
		if value.IsDir() {
			decoratedItem = dirItemStyle(value.Name())
		} else {
			decoratedItem = fileItemStyle(value.Name())
		}

		item := &item{title: decoratedItem, rawValue: value}
		items = append(items, item)
	}

	return items

}
