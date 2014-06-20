package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
)

var VERSION = "0.1.0"

var options struct {
	Debug bool   `short:"d" long:"debug" description:"Enable debugging mode"`
	File  string `short:"f" long:"file" description:"Specify path to config"`
}

func main() {
	args, err := flags.ParseArgs(&options, os.Args)

	if err != nil {
		os.Exit(1)
	}

	if len(args) < 2 {
		fmt.Println("Command required")
		os.Exit(1)
	}

	cmd := args[1]

	// Check if config file path has been provided
	if options.File == "" {
		terminate("Please provide config file", 1)
	}

	config := ParseYamlConfig(options.File)
	if config == nil {
		terminate("Unable to parse config file", 1)
	}

	fmt.Printf("\nShuttle v%s\n\n", VERSION)

	target := config.NewTarget()
	conn, err := NewConnection(target)

	if err != nil {
		terminate("Unable to establish connection", 1)
	}
	conn.debug = options.Debug

	logStep("Connected to " + target.toString())

	app := NewApp(target, conn, config)

	if err = app.initialize(); err != nil {
		exitWithError(err)
	}

	if cmd == "deploy" {
		if app.isLocked() {
			terminate("Deployment is locked", 1)
		}

		// Create application deployment structure, directories, etc
		logStep("Preparing application structure")
		if err = app.setupDirectoryStructure(); err != nil {
			logStep("Failed to setup application structure")
			exitWithError(err)
		}

		if !app.writeLock() {
			terminate("Unable to write lock", 2)
		}

		// Clone repository or update codebase on specified deployment branch
		if err = app.checkoutCode(); err != nil {
			exitWithError(err)
		}

		if err = app.writeCurrentReleaseNumber(); err != nil {
			app.cleanupCurrentRelease()
		}

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
