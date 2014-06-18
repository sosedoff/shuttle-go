package main

type App struct {
	name     string
	strategy string
	repo     string
	target   *Target
	conn     *Connection
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
