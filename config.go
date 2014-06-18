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
	return &Target{
		host:     conf.Target["host"],
		user:     conf.Target["user"],
		password: conf.Target["password"],
		path:     conf.Target["deploy_to"],
	}
}
