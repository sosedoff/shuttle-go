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
	return fmt.Sprintf("%s/releases", t.path)
}

func (t *Target) versionFilePath() string {
	return fmt.Sprintf("%s/version", t.path)
}

func (t *Target) sharedPath() string {
	return fmt.Sprintf("%s/shared", t.path)
}

func (t *Target) backupsPath() string {
	return fmt.Sprintf("%s/backups", t.path)
}

func (t *Target) lockfilePath() string {
	return fmt.Sprintf("%s/lock", t.path)
}

func (t *Target) repoPath() string {
	return fmt.Sprintf("%s/repo", t.path)
}
