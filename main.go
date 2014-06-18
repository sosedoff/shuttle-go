package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
)

var options struct {
	Debug       bool   `short:"d" long:"debug" description:"Enable debugging mode"`
	File        string `short:"f" long:"file" description:"Specify path to config"`
	Environment string `short:"e" long:"environment" description:"Deployment environment"`
}

func terminate(message string, status int) {
	fmt.Println("Deployment it locked.")
	os.Exit(1)
}

func main() {
	args, err := flags.ParseArgs(&options, os.Args)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(args) < 2 {
		fmt.Println("Command required")
		os.Exit(1)
	}

	target := Target{
		"192.168.33.10",
		"vagrant",
		"vagrant",
		"/home/vagrant/app",
	}

	conn, err := NewConnection(&target)

	conn.debug = options.Debug

	if err != nil {
		panic("Unable to establish connection")
	}

	app := NewApp(&target, conn)

	if app.isLocked() {
		terminate("Deployment is locked", 1)
	}

	if !app.writeLock() {
		terminate("Unable to write lock", 2)
	}

	app.setupDirectoryStructure()

	if !app.releaseLock() {
		terminate("Unable to release lock", 2)
	}
}
