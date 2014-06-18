package main

type App struct {
	target *Target
	conn   *Connection
}

func NewApp(target *Target, conn *Connection) *App {
	return &App{
		target: target,
		conn:   conn,
	}
}

func (app *App) setupDirectoryStructure() {
	paths := []string{
		app.target.basePath(),
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
