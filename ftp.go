package main

import (
	"crypto/tls"
	"log"
	"os"
	"time"

	ftp "github.com/jlaffaye/ftp"
)

func connectToServer() {
	c, err := ftp.Dial(
		os.Getenv("ftpserver"),
		ftp.DialWithExplicitTLS(
			&tls.Config{
				InsecureSkipVerify: true,
			},
		),
		ftp.DialWithDebugOutput(os.Stdout),
		ftp.DialWithTimeout(5*time.Second),
		ftp.DialWithDisabledEPSV(true),
	)

	handleError(err)

	err = c.Login(os.Getenv("username"), os.Getenv("password"))

	handleError(err)

	c.ChangeDir("/media")

	if err := c.Quit(); err != nil {
		log.Fatal(err)
	}
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}

}
