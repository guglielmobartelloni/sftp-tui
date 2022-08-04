package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/guglielmobartelloni/sftp-tui/tui"
	"github.com/pkg/sftp"
)

const (
	username       = "samoorai"
	password       = ""
	privateKeyPath = "/Users/samurai/.ssh/id_rsa"
	host           = "midas.usbx.me"
	port           = "22"
	knownHostsPath = "/Users/samurai/.ssh/known_hosts"
)

var (
	sshClient = ConnectSSH(
		username,
		privateKeyPath,
		password,
		host,
		port,
		knownHostsPath,
	)
	SftpClient, err = sftp.NewClient(sshClient)
)

func main() {
	//Close open connnections
	defer SftpClient.Close()
	defer sshClient.Close()

	m := tui.Model{
		List: list.New(
			tui.CreateItemListModel(".", SftpClient),
			list.NewDefaultDelegate(), 0, 0),
		SftpClient: SftpClient,
	}
	m.List.Title = "File List"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

}

func HandleError(err error) {
	if err != nil {
		log.Fatal(err)
	}

}
