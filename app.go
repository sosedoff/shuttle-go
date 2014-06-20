package main

import (
	"fmt"
	"strconv"
	"strings"
)

type App struct {
	target         *Target
	conn           *Connection
	config         *Config
	lastRelease    int // Last deployed release number
	currentRelease int // Current (new) release number
}

func NewApp(target *Target, conn *Connection, config *Config) *App {
	return &App{
		target: target,
		conn:   conn,
		config: config,
	}
}

func (app *App) initialize() error {
	app.lastRelease = app.getLastReleaseNumber()
	app.currentRelease = app.lastRelease + 1

	return nil
}

func (app *App) currentReleasePath() string {
	return fmt.Sprintf("%s/%d", app.target.releasesPath, app.currentRelease)
}

// Creates directories necessary for other deployment steps
func (app *App) setupDirectoryStructure() error {
	paths := []string{
		app.target.path,
		app.target.releasesPath,
		app.target.backupsPath,
		app.target.sharedPath,
		app.target.sharedPath + "/logs",
		app.target.sharedPath + "/pids",
		app.target.sharedPath + "/tmp",
	}

	for _, path := range paths {
		if result := app.conn.Exec("mkdir -p " + path); !result.Success {
			return fmt.Errorf(result.Output)
		}
	}

	return nil
}

// Returns true if remote server has a deployment lock file created by another
// deployment process
func (app *App) isLocked() bool {
	return app.conn.FileExists(app.target.lockfilePath)
}

// Writes deployment lock file to prevent simultaneous deployments
func (app *App) writeLock() bool {
	return app.conn.Exec("touch " + app.target.lockfilePath).Success
}

// Removes deployment lock file after deployment sequience has been completed
func (app *App) releaseLock() bool {
	return app.conn.Exec("rm " + app.target.lockfilePath).Success
}

// Write a new release number to the release file
func (app *App) writeReleaseNumber(number int) error {
	cmd := fmt.Sprintf("echo %d > %s", number, app.target.versionFilePath)
	result := app.conn.Exec(cmd)

	if !result.Success {
		return fmt.Errorf(result.Output)
	}

	return nil
}

func (app *App) writeCurrentReleaseNumber() error {
	return app.writeReleaseNumber(app.currentRelease)
}

// Returns last deployed release number, stored in "version" file
// Version file could only contain a numeric value
func (app *App) getLastReleaseNumber() int {
	if !app.conn.FileExists(app.target.versionFilePath) {
		return 0
	}

	value, err := app.conn.ReadFile(app.target.versionFilePath)

	if err != nil {
		fmt.Println("Unable to read version file:", err)
		return 0
	}

	number, err := strconv.Atoi(strings.TrimSpace(value))

	// If contents of the version file is invalid also return 0
	if err != nil {
		fmt.Println("Invalid version format:", err)
		return 0
	}

	return number
}

// Removes current release if any of the deployment steps fails
func (app *App) cleanupCurrentRelease() {
	app.conn.Exec("rm -rf " + app.currentReleasePath())
}
