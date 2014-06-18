package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
)

var options struct {
	Debug       string `long:"debug" description:"Enable debugging mode"`
	File        string `long:"file" description:"Specify path to config"`
	Environment string `long:"environment" description:"Deployment environment"`
}

func main() {
	args, err := flags.ParseArgs(&options, os.Args)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(args) < 2 {
		fmt.Println("Command required")
		os.Exit(1)
	}

	target := Target{
		"192.168.33.10",
		"vagrant",
		"vagrant",
		"/home/vagrant/app",
	}

	conn, err := NewConnection(&target)

	if err != nil {
		panic("Unable to establish connection")
	}

	app := NewApp(&target, conn)
	app.setupDirectoryStructure()
}
