package main

import (
	"crypto/tls"
	"fmt"
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
	)

	logErr(err)

	err = c.Login(os.Getenv("username"), os.Getenv("password"))

	logErr(err)

	fmt.Println(c.List("/"))

	if err := c.Quit(); err != nil {
		log.Fatal(err)
	}
}

func logErr(err error) {
	if err != nil {
		log.Fatal(err)
	}

}
