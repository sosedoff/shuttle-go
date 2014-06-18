package main

import (
	"fmt"
	"io/ioutil"
	"launchpad.net/goyaml"
)

type Config struct {
	App   map[string]string
	Hooks map[string]interface{}
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
