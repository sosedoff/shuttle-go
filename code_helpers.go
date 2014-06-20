package main

import (
	"fmt"
	"strings"
)

// Execute command within repository directory
func (app *App) repoExec(command string) *Command {
	// Prevent interactive mode when clone from HTTP/HTTPS
	opts := "export GIT_ASKPASS=echo"
	cmd := fmt.Sprintf("%s ; cd %s && %s", opts, app.target.repoPath, command)

	return app.conn.Exec(cmd)
}

// Clone repository or update an existing one from the upstream
func (app *App) checkoutCode() (err error) {
	if !app.conn.gitInstalled() {
		return fmt.Errorf("Git is not installed")
	}

	update := true

	// Clone git repository if directory does not exist
	if !app.conn.DirExists(app.target.repoPath) {
		update = false

		if err = app.cloneRepository(); err != nil {
			return
		}
	}

	// Run code update stuff unless its a new clone
	if update {
		if err = app.updateCode(); err != nil {
			return
		}
	}

	// Checkout branch specified in the config file
	if err = app.checkoutBranch(); err != nil {
		return
	}

	// Write index to the new release dir
	if err = app.checkoutIndex(); err != nil {
		return
	}

	return
}

// Checkout git branch specified by the configuration file
func (app *App) checkoutBranch() error {
	branch := app.config.getBranch()
	logStep("Using branch '" + branch + "'")

	result := app.repoExec("git checkout " + branch)
	if !result.Success {
		return fmt.Errorf(result.Output)
	}

	return nil
}

func (app *App) updateCode() (err error) {
	// Handle repository remote changes
	if app.remoteChanged() {
		logStep("Repository remote change detected: " + app.config.App["repo"])

		// If remote has been change, just delete repository and clone it again as
		// its an easier approach instead of switching remotes and updating codebase
		app.deleteRepository()

		if err = app.cloneRepository(); err != nil {
			return
		}
	}

	// Fetch repository changes from the remote
	return app.fetchCode()
}

func (app *App) fetchCode() error {
	logStep("Fetching latest code")

	// Make sure to hard reset before any updates
	if result := app.repoExec("git reset --hard"); !result.Success {
		return fmt.Errorf(result.Output)
	}

	// Fetch all new changes from the remote
	if result := app.repoExec("git fetch"); !result.Success {
		return fmt.Errorf(result.Output)
	}

	return nil
}

// Clone repository and checkout master branch
func (app *App) cloneRepository() error {
	repo := app.config.App["repo"]
	cmd := fmt.Sprintf(
		"export GIT_ASKPASS=echo ; cd %s && git clone --depth 50 --recursive --quiet %s repo",
		app.target.path, repo,
	)

	logStep("Cloning repository '" + repo + "'")
	result := app.conn.Exec(cmd)

	if !result.Success {
		return fmt.Errorf(result.Output)
	}

	return nil
}

// Removes git repository directory
func (app *App) deleteRepository() {
	app.conn.Exec("rm -rf " + app.target.repoPath)
}

// Returns true if existing git repository remote has been changed
func (app *App) remoteChanged() bool {
	oldRemote := app.conn.gitRemote(app.target.repoPath)
	newRemote := app.config.App["repo"]

	return oldRemote != newRemote
}

// Returns current git commit SHA
func (app *App) gitRevision() string {
	result := app.repoExec("git rev-parse HEAD")

	if result.Success {
		return strings.TrimSpace(result.Output)
	}

	return ""
}

func (app *App) checkoutIndex() error {
	result := app.repoExec("git checkout-index -a --prefix=" + app.currentReleasePath() + "/")

	if !result.Success {
		return fmt.Errorf(result.Output)
	}

	cmd := fmt.Sprintf("echo %s > %s/REVISION", app.gitRevision(), app.currentReleasePath())
	if result = app.conn.Exec(cmd); !result.Success {
		return fmt.Errorf(result.Output)
	}

	return nil
}
