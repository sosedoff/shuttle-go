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

func printVersion() {
	fmt.Printf("\nShuttle v%s\n\n", VERSION)
}

func realMain() {
	args, err := flags.ParseArgs(&options, os.Args)

	if err != nil {
		os.Exit(1)
	}

	if len(args) < 2 {
		fmt.Println("Command required")
		os.Exit(1)
	}

	cmd := args[1]

	if options.File == "" {
		terminate("Config file required", 1)
	}

	config := ParseYamlConfig(options.File)
	if config == nil {
		terminate("Unable to parse config file", 1)
	}

	if !config.isValidStrategy() {
		terminate("Invalid strategy: "+config.getStrategy(), 1)
	}

	if err = config.validateHooks(); err != nil {
		exitWithError(err)
	}

	printVersion()

	target := config.NewTarget()

	logStep("Establishing connection with remote server")
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

	if cmd == "setup" {
		if app.isLocked() {
			deployer := app.lastDeployer()
			message := fmt.Sprintf("Deployment is locked by %s", deployer)

			terminate(message, 1)
		}

		// Create application deployment structure, directories, etc
		logStep("Preparing application structure")
		if err = app.setup(); err != nil {
			logStep("Failed to setup application structure")
			exitWithError(err)
		}

		logStep("Application structure has been created")
	}

	if cmd == "deploy" {
		if app.isLocked() {
			deployer := app.lastDeployer()
			message := fmt.Sprintf("Deployment is locked by %s", deployer)

			terminate(message, 1)
		}

		// Create application deployment structure, directories, etc
		logStep("Preparing application structure")
		if err = app.setup(); err != nil {
			logStep("Failed to setup application structure")
			exitWithError(err)
		}

		if !app.writeLock() {
			terminate("Unable to write lock", 2)
		}

		// Make sure to release lock after established connection
		defer app.releaseLock()

		// Clone repository or update codebase on specified deployment branch
		if err = app.checkoutCode(); err != nil {
			exitWithError(err)
		}

		// Execute before_symlink hooks, no failures
		if err := app.runHook("before_link_release", false); err != nil {
			exitWithError(err)
		}

		// If current release cannot be symlinked, remove it
		if err = app.symlinkCurrentRelease(); err != nil {
			app.cleanupCurrentRelease()
			terminate("Unable to symlink current release", 3)
		}

		// Run user commands after release is linked, allow failures
		app.runHook("after_link_release", true)

		// Cleanup old releases
		app.cleanupOldReleases()

		// At this point deploy is considered completed
		logStep(fmt.Sprintf("Release v%d has been deployed", app.currentRelease))

		// Run user commands after release has been deployed
		app.runHook("after_deploy", true)

		fmt.Println("")
		return
	}

	if cmd == "lock" {
		if app.isLocked() {
			terminate("Deployment is already locked", 2)
		}

		if app.writeLock() {
			logStep("Deployment is successfully locked")
		} else {
			terminate("Unable to write lock", 2)
		}

		return
	}

	if cmd == "unlock" {
		if !app.isLocked() {
			terminate("Deployment is not locked", 2)
		}

		if app.releaseLock() {
			logStep("Deployment is successfully unlocked")
		} else {
			terminate("Unable to write lock", 2)
		}

		return
	}
}

func main() {
	realMain()
}
