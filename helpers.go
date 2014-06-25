package main

import (
	"fmt"
	"os"
	"strings"
)

func terminate(message string, status int) {
	fmt.Printf("-----> ERROR: %s\n", strings.TrimSpace(message))
	os.Exit(1)
}

func exitWithError(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func logStep(str string) {
	fmt.Printf("-----> %s\n", str)
}
