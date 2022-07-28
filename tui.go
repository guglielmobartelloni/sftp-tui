package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	docStyle           = lipgloss.NewStyle().Margin(1, 2)
	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

type item struct {
	title       string
	description string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

type model struct {
	list   list.Model
	walker *walker
}

func (m model) Init() tea.Cmd {
	sshClient := ConnectSSH("samoorai", "/Users/samurai/.ssh/id_rsa", "", "midas.usbx.me", "22", "/Users/samurai/.ssh/known_hosts")
	m.walker = &walker{
		sshClient:  sshClient,
		currentDir: "./",
	}
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			cmd := m.list.NewStatusMessage(statusMessageStyle(fmt.Sprintf("Donwloading %s", "banana")))
			m.walker.GetFile("banana", "/Users/samurai/Documents/progetti/ftp-tui/test")
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

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func createItemList() []list.Item {

	sshClient := ConnectSSH("samoorai", "/Users/samurai/.ssh/id_rsa", "", "midas.usbx.me", "22", "/Users/samurai/.ssh/known_hosts")
	walker := &walker{
		sshClient:  sshClient,
		currentDir: "./",
	}
	// walker.GetFile("banana", "/Users/samurai/Documents/progetti/ftp-tui/test")

	items := []list.Item{}
	fileList, err := walker.Ls()
	handleError(err)
	fmt.Println(fileList)

	for _, value := range fileList {
		item := &item{title: value}
		items = append(items, item)
	}

	return items[:len(items)-1]

}
