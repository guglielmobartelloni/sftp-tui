package main

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"os"

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
		// ftp.DialWithTimeout(5*time.Second),
		ftp.DialWithDisabledEPSV(true),
	)

	handleError(err)

	err = c.Login(os.Getenv("username"), os.Getenv("password"))

	handleError(err)

	r, err := c.Retr(".bashrc")
	if err != nil {
		panic(err)
	}
	defer r.Close()

	buf, err := ioutil.ReadAll(r)
	handleError(err)
	println(string(buf))

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
