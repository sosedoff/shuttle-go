package main

import (
	"fmt"
)

type Target struct {
	host             string
	user             string
	password         string
	path             string
	releasesPath     string
	currentPath      string
	versionFilePath  string
	revisionFilePath string
	sharedPath       string
	backupsPath      string
	lockfilePath     string
	repoPath         string
}

func (t *Target) getAddress() string {
	return fmt.Sprintf("%s:22", t.host)
}

func (t *Target) toString() string {
	return fmt.Sprintf("%s@%s", t.user, t.host)
}
