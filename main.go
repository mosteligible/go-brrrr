package main

import (
	"os"

	"github.com/mosteligible/go-brrrr/cmd"
)

func init() {
	outputDirectory := "./out"
	if _, err := os.Stat(outputDirectory); os.IsNotExist(err) {
		os.Mkdir(outputDirectory, os.ModePerm)
	}
}

func main() {
	cmd.Execute()
}
