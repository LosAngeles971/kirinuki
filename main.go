package main

import (
	"log"
	"os"
	"os/user"
	"runtime"
	"it/losangeles971/kirinuki/cmd"
)

func setLogging() {
	var logfile string
	host := runtime.GOOS
	usr, err := user.Current()
	if err != nil {
		log.Panic("cannot detect home directory on ", host)
		os.Exit(1)
	}
	switch host {
	case "windows":
		logfile = usr.HomeDir + "/AppData/Local/Temp/kirinuki.log"
	case "darwin":
		logfile = usr.HomeDir + "/.kirinuki/kirinuki.log"
	case "linux":
		logfile = usr.HomeDir + "/.kirinuki/kirinuki.log"
	default:
		log.Panic("unsupported operating system ", host)
		os.Exit(1)
	}
	f, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("error opening file %s due to: %v", logfile, err)
	}
	defer f.Close()
	log.SetOutput(f)
}

func main() {
	//setLogging()
	cmd.Execute()
}
