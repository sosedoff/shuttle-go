package main

import (
	"fmt"
	"io/ioutil"
	"launchpad.net/goyaml"
)

type Config struct {
	App    map[string]string
	Hooks  map[string]interface{}
	Target map[string]string
}

func ParseYamlConfig(path string) *Config {
	buff, err := ioutil.ReadFile(path)

	if err != nil {
		return nil
	}

	config := Config{}

	if err = goyaml.Unmarshal(buff, &config); err != nil {
		fmt.Println(err)
		return nil
	}

	return &config
}

func (conf *Config) NewTarget() *Target {
	deployPath := conf.Target["deploy_to"]

	if deployPath == "" {
		deployPath = conf.Target["path"]
	}

	return &Target{
		host:            conf.Target["host"],
		user:            conf.Target["user"],
		password:        conf.Target["password"],
		path:            deployPath,
		releasesPath:    deployPath + "/releases",
		currentPath:     deployPath + "/current",
		versionFilePath: deployPath + "/version",
		sharedPath:      deployPath + "/shared",
		backupsPath:     deployPath + "/backups",
		lockfilePath:    deployPath + "/lock",
		repoPath:        deployPath + "/repo",
	}
}

func (conf *Config) getBranch() string {
	branch := conf.App["branch"]

	if branch == "" {
		branch = "master"
	}

	return branch
}

func (conf *Config) getStrategy() string {
	strategy := conf.App["strategy"]

	if strategy == "" {
		strategy = "static"
	}

	return strategy
}

func (conf *Config) isValidStrategy() bool {
	if conf.getStrategy() == "static" {
		return true
	}

	return false
}
