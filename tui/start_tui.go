package tui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pkg/sftp"
	"github.com/guglielmobartelloni/sftp-tui/ssh"
)

//const (
//	username       = "samoorai"
//	password       = ""
//	privateKeyPath = "/Users/samurai/.ssh/id_rsa"
//	host           = "midas.usbx.me"
//	port           = "22"
//	knownHostsPath = "/Users/samurai/.ssh/known_hosts"
//)

func StartProgram(username, privateKeyPath, password, host, port, knownHostsPath string) {
	sshClient := ssh.ConnectSSH(
		username,
		privateKeyPath,
		password,
		host,
		port,
		knownHostsPath,
	)
	SftpClient, err := sftp.NewClient(sshClient)
	handleError(err)
	//Close open connnections
	defer SftpClient.Close()
	defer sshClient.Close()

	m := Model{
		List: list.New(
			CreateItemListModel(".", SftpClient),
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

