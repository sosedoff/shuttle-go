package main

import (
	"fmt"
)

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

func (conn *Connection) ReadFile(path string) (content string, err error) {
	if !conn.FileExists(path) {
		err = fmt.Errorf("File does not exist: %s", path)
		return
	}

	result := conn.Exec("cat " + path)

	if result.Success {
		content = result.Output
		err = nil
	} else {
		err = fmt.Errorf("Cant read file %s: %s", path, result.Output)
	}

	return content, err
}
