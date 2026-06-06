package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "taskr",
	Short: "taskr — a file-based ticket tracker",
	Long:  `taskr manages tickets as markdown files in .tickets/ within your project directory.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func SetVersion(v string) {
	rootCmd.Version = v
	rootCmd.SetVersionTemplate("Version: {{.Version}}\n")
}
