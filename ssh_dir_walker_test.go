package main

import (
	"fmt"
	"testing"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestDirList(t *testing.T) {
	sshClient := ConnectSSH("samoorai", "/Users/samurai/.ssh/id_rsa", "", "midas.usbx.me", "22", "/Users/samurai/.ssh/known_hosts")
	walker := &walker{
		sshClient:  sshClient,
		currentDir: "./",
	}
	fmt.Println(walker.LsDir())
	fmt.Println(walker.LsFiles())

}
