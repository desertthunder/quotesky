package main

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/desertthunder/quotesky/cmd"
)

func main() {
	if err := cmd.Execute(cmd.Port); err != nil {
		log.Errorf("failed to run cli: %s", err)
		os.Exit(1)
	}
}
