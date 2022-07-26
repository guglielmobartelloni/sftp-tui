package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func main() {

	m := newModel()
	m.list.Title = "Files list"
	err := tea.NewProgram(m, tea.WithAltScreen()).Start()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

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

	for _, value := range fileList {
		item := &item{title: value}
		items = append(items, item)
	}

	return items

}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "choose"),
		),
		// remove: key.NewBinding(
		// 	key.WithKeys("x", "backspace"),
		// 	key.WithHelp("x", "delete"),
		// ),
	}
}

func newModel() model {
	return model{
		list:         list.New(createItemList(), newItemDelegate(newDelegateKeyMap()), 0, 0),
		delegateKeys: newDelegateKeyMap(),
	}
}
