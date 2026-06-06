package cmd

import (
	"fmt"
	"os"

	"github.com/ro-56/taskr/internal/tickets"
	"github.com/spf13/cobra"
)

var readyMode string

var readyCmd = &cobra.Command{
	Use:   "ready",
	Short: "List actionable tickets (no blocking deps)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		result, err := tickets.Ready(cwd, readyMode)
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
	readyCmd.Flags().StringVar(&readyMode, "mode", "", "filter by mode (e.g. afk)")
	rootCmd.AddCommand(readyCmd)
}
