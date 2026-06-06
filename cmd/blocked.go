package cmd

import (
	"fmt"
	"os"

	"github.com/ro-56/taskr/internal/tickets"
	"github.com/spf13/cobra"
)

var blockedCmd = &cobra.Command{
	Use:   "blocked",
	Short: "List tickets that have at least one non-closed dependency",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		result, err := tickets.Blocked(cwd)
		if err != nil {
			return err
		}

		for _, t := range result {
			fmt.Printf("%-20s  %-40s  priority:%-2d  %s\n", t.ID, t.Title, t.Priority, t.Status)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(blockedCmd)
}
