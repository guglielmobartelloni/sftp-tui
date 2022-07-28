package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/bramvdbogaerde/go-scp"
	"golang.org/x/crypto/ssh"
)

type walker struct {
	sshClient  *ssh.Client
	currentDir string
}

func (w *walker) LsFiles() ([]string, error) {
	output, err := RunCommand("ls -p | grep -v /", w.sshClient)
	// output, err := RunCommand("ls --color", w.sshClient)
	fileList := strings.Split(output, "\n")
	return fileList[:len(fileList)-1], err
}

func (w *walker) LsDir() ([]string, error) {
	output, err := RunCommand("find . -maxdepth 1 -type d -print", w.sshClient)
	dirList := strings.Split(output, "\n")
	//Strip the first two char from the file list
	dirList = dirList[1 : len(dirList)-1]
	for i, v := range dirList {
		dirList[i] = v[2:]
	}
	return dirList, err
}

func (w *walker) Cd(dir string) ([]string, error) {
	output, err := RunCommand(fmt.Sprintf("cd %s && ls", dir), w.sshClient)
	if err == nil {
		w.currentDir = fmt.Sprintf("%s%s/", w.currentDir, dir)
	}
	return strings.Split(output, "\n"), err
}

func (w *walker) GetFile(fileName string, destinationPath string) *os.File {

	client, err := scp.NewClientBySSH(w.sshClient)
	if err != nil {
		fmt.Println("Error creating new SSH session from existing connection", err)
	}

	err = client.Connect()
	if err != nil {
		fmt.Println("Couldn't establish a connection to the remote server ", err)
	}

	buff := new(bytes.Buffer)

	path := w.currentDir + fileName

	client.CopyFromRemotePassThru(context.Background(), buff, path, nil)
	f, err := os.Create(destinationPath)
	handleError(err)
	defer f.Close()
	f.Write(buff.Bytes())
	return f
}
