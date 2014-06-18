package main

import (
	"fmt"
)

type Target struct {
	host     string
	user     string
	password string
	path     string
}

func (t *Target) getAddress() string {
	return fmt.Sprintf("%s:22", t.host)
}

func (t *Target) releasesPath() string {
	return t.path + "/releases"
}

func (t *Target) currentPath() string {
	return t.path + "/current"
}

func (t *Target) versionFilePath() string {
	return t.path + "/version"
}

func (t *Target) revisionFilePath() string {
	return t.path + "/REVISION"
}

func (t *Target) sharedPath() string {
	return t.path + "/shared"
}

func (t *Target) backupsPath() string {
	return t.path + "/backups"
}

func (t *Target) lockfilePath() string {
	return t.path + "/lock"
}

func (t *Target) repoPath() string {
	return t.path + "/repo"
}
