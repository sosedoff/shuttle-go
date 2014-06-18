package main

import(
  "testing"
  "github.com/stretchr/testify/assert"
)

var target = Target{
  host: "localhost",
  user: "user",
  password: "password",
  path: "/var/www/app",
}

func Test_getAddress(t *testing.T) {
  assert.Equal(t, target.getAddress(), "localhost:22")
}

func Test_releasesPath(t *testing.T) {
  assert.Equal(t, target.releasesPath(), "/var/www/app/releases")
}

func Test_sharedPath(t *testing.T) {
  assert.Equal(t, target.sharedPath(), "/var/www/app/shared")
}

func Test_backupsPath(t *testing.T) {
  assert.Equal(t, target.backupsPath(), "/var/www/app/backups")
}

func Test_lockfilePath(t *testing.T) {
  assert.Equal(t, target.lockfilePath(), "/var/www/app/lock")
}

func Test_repoPath(t *testing.T) {
  assert.Equal(t, target.repoPath(), "/var/www/app/repo")
}