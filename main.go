package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func main() {

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
