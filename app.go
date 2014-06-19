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

func (app *App) setupDirectoryStructure() {
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
		app.conn.Run("mkdir -p " + path)
	}
}

func (app *App) isLocked() bool {
	return app.conn.FileExists(app.target.lockfilePath)
}

func (app *App) writeLock() bool {
	return app.conn.Exec("touch " + app.target.lockfilePath).Success
}

func (app *App) releaseLock() bool {
	return app.conn.Exec("rm " + app.target.lockfilePath).Success
}

// Clone repository or update an existing one from the upstream
func (app *App) checkoutCode() error {
	// Do not proceed if git is not installed, its the only hard requirement
	if !app.conn.GitInstalled() {
		return fmt.Errorf("Git executable is not installed")
	}

	if app.conn.DirExists(app.target.repoPath) {
		// Check if repository remote has changed.
		// When remote changes its not always easy to switch remotes.
		// In this case just remove the repo, its easier than updating it.
		if app.gitRemoteChanged() {
			app.conn.Exec("rm -rf " + app.target.repoPath)
		} else {
			return app.updateCode()
		}
	}

	return app.cloneRepository()
}

func (app *App) cloneRepository() error {
	branch := app.config.getBranch()
	cloneOpts := "--depth 25 --recursive --quiet"
	cloneCmd := fmt.Sprintf("git clone %s %s repo", cloneOpts, app.config.App["repo"])
	cmd := fmt.Sprintf("cd %s && %s", app.target.path, cloneCmd)
	result := app.conn.Exec(cmd)

	if !result.Success {
		return fmt.Errorf(result.Output)
	}

	if branch != "master" {
		if err := app.checkoutBranch(); err != nil {
			return err
		}
	}

	return nil
}

func (app *App) checkoutBranch() error {
	branch := app.config.getBranch()
	cmd := fmt.Sprintf("cd %s && git checkout %s", app.target.repoPath, branch)
	result := app.conn.Exec(cmd)

	if !result.Success {
		return fmt.Errorf(result.Output)
	}

	return nil
}

// Pulls new changes from the upstream
func (app *App) updateCode() error {
	branch := app.config.getBranch()
	cmd := fmt.Sprintf("cd %s && git pull origin %s", app.target.repoPath, branch)
	result := app.conn.Exec(cmd)

	if !result.Success {
		return fmt.Errorf(result.Output)
	}

	return nil
}

// Write a new release number to the release file
func (app *App) writeReleaseNumber(number string) error {
	cmd := fmt.Sprintf("echo %s > %s", number, app.target.versionFilePath)
	result := app.conn.Exec(cmd)

	if !result.Success {
		return fmt.Errorf(result.Output)
	}

	return nil
}

// Returns true if existing git repository remote has been changed
func (app *App) gitRemoteChanged() bool {
	oldRemote := app.conn.GitRemote(app.target.repoPath)
	newRemote := app.config.App["repo"]

	return oldRemote != newRemote
}

// Returns current git commit SHA
func (app *App) gitRevision() string {
	cmd := fmt.Sprintf("cd %s && git rev-parse HEAD", app.target.repoPath)
	sha := strings.TrimSpace(app.conn.Run(cmd))

	return sha
}

// Returns last deployed release number, stored in "version" file
// Version file could only contain a numeric value
func (app *App) getLastReleaseNumber() int {
	// Return 0 as a non-release if version file does not exist or deployment
	// directory/file structure has been broken
	if !app.conn.FileExists(app.target.versionFilePath) {
		return 0
	}

	value, err := app.conn.ReadFile(app.target.versionFilePath)

	// If file cant be read we also need to return 0
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
