package main

import (
	"fmt"
	"os"
)

func terminate(message string, status int) {
	fmt.Println(message)
	os.Exit(1)
}

func exitWithError(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func logStep(str string) {
	fmt.Printf("-----> %s\n", str)
}
