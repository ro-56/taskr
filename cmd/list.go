package cmd

import (
	"fmt"
	"os"

	"github.com/ro-56/taskr/internal/tickets"
	"github.com/spf13/cobra"
)

var (
	listFlagTags   bool
	listFlagCount  bool
	listFlagStatus string
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List tickets",
	RunE:    runList,
}

func runList(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	switch {
	case listFlagTags:
		tags, err := tickets.ListTags(cwd)
		if err != nil {
			return err
		}
		for _, t := range tags {
			fmt.Println(t)
		}

	case listFlagCount:
		counts, err := tickets.ListCounts(cwd)
		if err != nil {
			return err
		}
		for _, status := range []string{"open", "in_progress", "closed"} {
			fmt.Printf("%-12s %d\n", status+":", counts[status])
		}

	case listFlagStatus != "":
		results, err := tickets.ListByStatus(cwd, listFlagStatus)
		if err != nil {
			return err
		}
		for _, r := range results {
			fmt.Printf("%-20s %s\n", r.ID, r.Title)
		}

	default:
		results, err := tickets.List(cwd)
		if err != nil {
			return err
		}
		for _, r := range results {
			fmt.Printf("%-20s %-12s %s\n", r.ID, r.Status, r.Title)
		}
	}

	return nil
}

func init() {
	listCmd.Flags().BoolVar(&listFlagTags, "tags", false, "Print all unique tags")
	listCmd.Flags().BoolVar(&listFlagCount, "count", false, "Print per-status counts")
	listCmd.Flags().StringVar(&listFlagStatus, "status", "", "Filter by status (open, in_progress, closed)")
	rootCmd.AddCommand(listCmd)
}
