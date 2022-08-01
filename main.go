package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pkg/sftp"
)

var (
	username       = "samoorai"
	password       = ""
	privateKeyPath = "/Users/samurai/.ssh/id_rsa"
	host           = "midas.usbx.me"
	port           = "22"
	knownHostsPath = "/Users/samurai/.ssh/known_hosts"
)

func main() {
	sshClient := ConnectSSH(username, privateKeyPath, password, host, port, knownHostsPath)
	sftpClient, err := sftp.NewClient(sshClient)
	handleError(err)

	items := createItemList(sftpClient)

	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0), sftpClient: sftpClient}
	m.list.Title = "File List"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}

}
