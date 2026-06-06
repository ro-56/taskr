package cmd

import (
	"fmt"
	"os"

	"github.com/ro-56/taskr/internal/tickets"
	"github.com/spf13/cobra"
)

var (
	addType     string
	addPriority int
	addMode     string
	addTags     []string
)

var addCmd = &cobra.Command{
	Use:   "add <title>",
	Short: "Create a new ticket",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		p := addPriority
		opts := tickets.AddOptions{
			Title:    args[0],
			Type:     addType,
			Priority: &p,
			Mode:     addMode,
			Tags:     addTags,
		}

		id, err := tickets.Add(cwd, opts)
		if err == tickets.ErrNotInitialized {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err != nil {
			return err
		}

		fmt.Println(id)
		return nil
	},
}

func init() {
	addCmd.Flags().StringVar(&addType, "type", "task", "ticket type (bug, feature, task, epic, chore)")
	addCmd.Flags().IntVar(&addPriority, "priority", 2, "priority 0=critical … 3=low")
	addCmd.Flags().StringVar(&addMode, "mode", "hitl", "mode (afk, hitl)")
	addCmd.Flags().StringSliceVar(&addTags, "tags", nil, "comma-separated tags")
	rootCmd.AddCommand(addCmd)
}
