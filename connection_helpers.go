package main

func (conn *Connection) FileExists(path string) bool {
	return conn.Exec("test -f " + path).Success
}

func (conn *Connection) DirExists(path string) bool {
	return conn.Exec("test -d " + path).Success
}

func (conn *Connection) SymlinkExists(path string) bool {
	return conn.Exec("test -h " + path).Success
}

func (conn *Connection) ProcessExists(pid string) bool {
	return conn.Exec("ps -p " + pid).Success
}

func (conn *Connection) SvnInstalled() bool {
	return conn.Exec("which svn").Success
}
