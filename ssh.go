package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type SSHFileManager struct {
	sshClient  *SshClient
	currentDir []string
}

func StartConnection() *SSHFileManager {
	ssh, err := NewSshClient(
		os.Getenv("username"),
		os.Getenv("host"),
		22,
		"/Users/samurai/.ssh/id_rsa",
		"")

	if err != nil {
		log.Printf("SSH init error %v", err)
	}

	return &SSHFileManager{currentDir: []string{"."}, sshClient: ssh}
}

func (fm *SSHFileManager) Ls() []string {
	lsOutput, err := fm.sshClient.RunCommand("ls")
	fileList := strings.Split(lsOutput, "\n")
	if err != nil {
		log.Printf("SSH run command error %v", err)
	}
	return fileList
}

func (fm *SSHFileManager) Cd(dir string) []string {
	command := fmt.Sprintf("cd %s; ls", strings.Join(fm.currentDir, "/")+"/"+dir)
	output, err := fm.sshClient.RunCommand(command)
	fileList := strings.Split(output, "\n")
	if err != nil {
		log.Printf("SSH run command error %v", err)
	}
	fm.currentDir = append(fm.currentDir, dir)
	return fileList
}
