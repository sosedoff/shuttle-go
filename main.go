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
	fmt.Println(message)
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

	cmd := args[1]

	config := ParseYamlConfig(options.File)
	if config == nil {
		terminate("Unable to parse config file", 1)
	}

	target := config.NewTarget()

	conn, err := NewConnection(target)

	if err != nil {
		terminate("Unable to establish connection", 1)
	}

	conn.debug = options.Debug

	app := NewApp(target, conn, config)

	if cmd == "deploy" {
		if app.isLocked() {
			terminate("Deployment is locked", 1)
		}

		app.setupDirectoryStructure()

		if !app.writeLock() {
			terminate("Unable to write lock", 2)
		}

		app.checkoutCode()

		if !app.releaseLock() {
			terminate("Unable to release lock", 2)
		}
	}

	if cmd == "unlock" {
		if !app.isLocked() {
			return
		}

		if !app.releaseLock() {
			terminate("Unable to release lock", 2)
		}
	}
}
