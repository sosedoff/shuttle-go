package main

import (
	"fmt"
	"github.com/mgutz/ansi"
	"os"
	"strings"
)

var (
	greenPrefix  = ansi.ColorCode("green") + "----->" + ansi.ColorCode("reset")
	redPrefix    = ansi.ColorCode("red") + "----->" + ansi.ColorCode("reset")
	yellowPrefix = ansi.ColorCode("yellow") + "----->" + ansi.ColorCode("reset")
)

func terminate(message string, status int) {
	fmt.Printf("%s ERROR: %s\n", redPrefix, strings.TrimSpace(message))
	os.Exit(1)
}

func logStep(message string) {
	fmt.Printf("%s %s\n", greenPrefix, strings.TrimSpace(message))
}

func logOutputLine(line string) {
	fmt.Println("      ", strings.TrimSpace(line))
}

func logOutput(output string) {
	trimmedOutput := strings.TrimSpace(output)

	if trimmedOutput == "" {
		return
	}

	logOutputLine("")

	lines := strings.Split(output, "\n")

	for _, line := range lines {
		logOutputLine(line)
	}
}

func exitWithError(err error) {
	fmt.Println(err)
	os.Exit(1)
}
