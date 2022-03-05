package cmd

import (
    "os"
    "golang.org/x/crypto/ssh/terminal"
	"log"
)

func askPassword() string {
    log.Println("Enter your Kirinuki passphrase")
    data, err := terminal.ReadPassword(int(os.Stdin.Fd()))
    if err != nil {
        log.Panic(err)
        os.Exit(1)
    }
    return string(data)
}