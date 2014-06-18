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

func (target *Target) getAddress() string {
	return fmt.Sprintf("%s:22", target.host)
}

func (target *Target) basePath() string {
	return target.path
}

func (target *Target) releasesPath() string {
	return fmt.Sprintf("%s/releases", target.basePath())
}

func (target *Target) sharedPath() string {
	return fmt.Sprintf("%s/shared", target.basePath())
}

func (target *Target) backupsPath() string {
	return fmt.Sprintf("%s/backups", target.basePath())
}

func (target *Target) lockfilePath() string {
	return fmt.Sprintf("%s/lock", target.basePath())
}
