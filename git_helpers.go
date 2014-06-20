package main

import (
	"fmt"
	"regexp"
	"strings"
)

var regexpGitRemote = regexp.MustCompile("origin\t+(.*)\\s\\(fetch\\)")

func (conn *Connection) gitInstalled() bool {
	return conn.Exec("which git").Success
}

func (conn *Connection) gitRemote(path string) string {
	result := conn.Exec(fmt.Sprintf("cd %s && git remote -v", path))

	if !result.Success {
		return ""
	}

	matches := regexpGitRemote.FindAllString(result.Output, -1)

	if matches == nil {
		return ""
	}

	remoteStr := strings.Replace(matches[0], "\t", " ", -1)
	remote := strings.Split(remoteStr, " ")[1]

	return remote
}
