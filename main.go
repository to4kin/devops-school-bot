package main

import (
	"fmt"
	"os"

	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %+v\n", err)
		os.Exit(1)
	}
}
