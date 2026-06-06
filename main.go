package main

import (
	"fmt"
	"os"

	"github.com/ro-56/taskr/cmd"
)

var version = "dev"

func main() {
	cmd.SetVersion(version)
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
