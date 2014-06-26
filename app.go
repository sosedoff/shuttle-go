package main

import (
	"fmt"
	"os/exec"
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

// Setup application directory structure
func (app *App) setup() error {
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

// Execute a command and return result
func (app *App) Run(command string) *Command {
	return app.conn.Exec(command)
}

// Returns true if remote server has a deployment lock file created by another
// deployment process
func (app *App) isLocked() bool {
	return app.conn.FileExists(app.target.lockfilePath)
}

// Writes deployment lock file to prevent simultaneous deployments
func (app *App) writeLock() bool {
	out, _ := exec.Command("hostname").Output()
	cmd := fmt.Sprintf("echo %s > %s", strings.TrimSpace(string(out)), app.target.lockfilePath)

	return app.conn.Exec(cmd).Success
}

// Get hostname of the last deployer stored in lockfile
func (app *App) lastDeployer() string {
	deployer := app.conn.Exec("cat " + app.target.lockfilePath).Output
	return strings.TrimSpace(deployer)
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

// Remove old releases
func (app *App) cleanupOldReleases() {
	keep := 10 // keep last 10
	cmd := fmt.Sprintf("cd %s && ls -1 | sort -rn", app.target.releasesPath)
	output := strings.TrimSpace(app.conn.Exec(cmd).Output)
	releases := strings.Split(output, "\r\n")

	if len(releases) < keep {
		return
	}

	logStep("Cleaning up old releases")

	for i, num := range releases {
		if i >= keep && num != "" {
			app.conn.Exec(fmt.Sprintf("rm -rf %s/%s", app.target.releasesPath, num))
		}
	}
}

func (app *App) symlinkCurrentRelease() error {
	logStep("Linking release")

	// If symlink already exists, unlink
	if app.conn.SymlinkExists(app.target.currentPath) {
		result := app.conn.Exec("unlink " + app.target.currentPath)

		if !result.Success {
			return fmt.Errorf(result.Output)
		}
	}

	// If current is a directory, remove it
	if app.conn.DirExists(app.target.currentPath) {
		result := app.conn.Exec("rm -rf " + app.target.currentPath)

		if !result.Success {
			return fmt.Errorf(result.Output)
		}
	}

	result := app.conn.Exec("ln -s " + app.currentReleasePath() + " " + app.target.currentPath)

	if !result.Success {
		return fmt.Errorf(result.Output)
	}

	// Write current version into RELEASE file
	cmd := fmt.Sprintf("echo %d > %s/RELEASE", app.currentRelease, app.currentReleasePath())
	if result = app.conn.Exec(cmd); !result.Success {
		return fmt.Errorf(result.Output)
	}

	// Write current version into "version" file
	if err := app.writeCurrentReleaseNumber(); err != nil {
		return err
	}

	// Cleanup old releases
	app.cleanupOldReleases()

	logStep(fmt.Sprintf("Release v%d has been deployed", app.currentRelease))
	fmt.Println("")

	return nil
}
