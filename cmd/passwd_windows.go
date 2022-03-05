package cmd

/*
But the second reason is…. if you are using Cygwin/mintty/git-bash on Windows, those Windows shells are unable to reach down to the OS API, and will throw the exact same error of the “handle is invalid”.

This issue is not directly fixable and really not an issue with Go. If you switch to Powershell or even CMD then executing ReadPassword will work as expected. You may then switch back to your shell for all other commands that don’t invoke ReadPassword. If you must stay in Cygwin/minty/git-bash then take a look at https://github.com/rprichard/winpty project, it might solve your issue
*/

import (
    "os"
    "syscall"
    "golang.org/x/crypto/ssh/terminal"
	"log"
)

func askPassword() string {
	log.Println("Enter your Kirinuki passphrase")
	data, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Panic(err)
		os.Exit(1)
	}
	return string(data)
}