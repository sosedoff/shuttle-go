package main

import (
	"fmt"
)

type App struct {
	target *Target
	conn   *Connection
	config *Config
}

func NewApp(target *Target, conn *Connection, config *Config) *App {
	return &App{
		target: target,
		conn:   conn,
		config: config,
	}
}

func (app *App) setupDirectoryStructure() {
	paths := []string{
		app.target.path,
		app.target.releasesPath(),
		app.target.backupsPath(),
		app.target.sharedPath(),
		app.target.sharedPath() + "/logs",
		app.target.sharedPath() + "/pids",
		app.target.sharedPath() + "/tmp",
	}

	for _, path := range paths {
		app.conn.Run("mkdir -p " + path)
	}
}

func (app *App) isLocked() bool {
	return app.conn.FileExists(app.target.lockfilePath())
}

func (app *App) writeLock() bool {
	return app.conn.Exec("touch " + app.target.lockfilePath()).Success
}

func (app *App) releaseLock() bool {
	return app.conn.Exec("rm " + app.target.lockfilePath()).Success
}

func (app *App) checkoutCode() bool {
	if !app.conn.GitInstalled() {
		fmt.Println("Git is not installed.")
		return false
	}

	if app.conn.DirExists(app.target.repoPath()) {
		// Check if repository remote has been changed
		if app.gitRemoteChanged() {
			fmt.Println("Git remote change detected.")

			// Just remote the reposity, its easier
			app.conn.Exec("rm -rf " + app.target.repoPath())
		} else {
			return app.updateCode()
		}
	}

	return app.cloneRepository()
}

func (app *App) cloneRepository() bool {
	fmt.Println("Cloning repository")

	branch := app.config.getBranch()
	cloneOpts := "--depth 25 --recursive --quiet"
	cloneCmd := fmt.Sprintf("git clone %s %s repo", cloneOpts, app.config.App["repo"])
	cmd := fmt.Sprintf("cd %s && %s", app.target.path, cloneCmd)
	result := app.conn.Exec(cmd)

	if !result.Success {
		fmt.Println("Failed to clone repository")
		fmt.Print(result.Output)
	}

	if branch != "master" {
		cmd = fmt.Sprintf("cd %s && git checkout %s", app.target.repoPath(), branch)
		result = app.conn.Exec(cmd)

		if !result.Success {
			fmt.Println("Failed to checkout branch")
			fmt.Print(result.Output)
		}
	}

	return result.Success
}

func (app *App) updateCode() bool {
	fmt.Println("Updating code")

	branch := app.config.getBranch()
	cmd := fmt.Sprintf("cd %s && git pull origin %s", app.target.repoPath(), branch)
	result := app.conn.Exec(cmd)

	if !result.Success {
		fmt.Println("Failed to update repository")
		fmt.Print(result.Output)
	}

	return result.Success
}

// Returns true if existing git repository remote has been changed
func (app *App) gitRemoteChanged() bool {
	oldRemote := app.conn.GitRemote(app.target.repoPath())
	newRemote := app.config.App["repo"]

	return oldRemote != newRemote
}