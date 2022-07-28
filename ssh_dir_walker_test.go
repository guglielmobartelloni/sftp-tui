package main

import (
	"testing"

	"golang.org/x/exp/slices"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestFilesList(t *testing.T) {
	sshClient := ConnectSSH(username, privateKeyPath, password, host, port, knownHostsPath)
	walker := &walker{
		sshClient:  sshClient,
		currentDir: "./",
	}
	files, err := walker.LsFiles()
	handleError(err)
	if !sliceContains(files, []string{"banana", "fragola"}) {
		t.Error("Error mock files not included")
	}

}

func sliceContains(list, other []string) bool {

	for _, v := range other {
		if slices.Contains(list, v) {
			return true
		}
	}

	return false

}
