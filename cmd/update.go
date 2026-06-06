package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ro-56/taskr/internal/tickets"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update ticket fields non-interactively",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		flags := cmd.Flags()
		opts := tickets.UpdateOptions{}

		if flags.Changed("title") {
			v, _ := flags.GetString("title")
			opts.Title = &v
		}
		if flags.Changed("priority") {
			v, _ := flags.GetInt("priority")
			opts.Priority = &v
		}
		if flags.Changed("mode") {
			v, _ := flags.GetString("mode")
			opts.Mode = &v
		}
		if flags.Changed("tags") {
			raw, _ := flags.GetString("tags")
			opts.Tags = strings.Split(raw, ",")
		}

		err = tickets.Update(cwd, args[0], opts)
		if err != nil {
			var ambErr *tickets.ErrAmbiguous
			if errors.As(err, &ambErr) {
				fmt.Fprintln(os.Stderr, "ambiguous ID — did you mean one of:")
				for _, m := range ambErr.Matches {
					fmt.Fprintf(os.Stderr, "  %s\n", m)
				}
				os.Exit(1)
			}
			if errors.Is(err, tickets.ErrNotFound) {
				fmt.Fprintf(os.Stderr, "ticket not found: %s\n", args[0])
				os.Exit(1)
			}
			if errors.Is(err, tickets.ErrInvalidPriority) {
				fmt.Fprintln(os.Stderr, "invalid priority: must be 0–3")
				os.Exit(1)
			}
			if errors.Is(err, tickets.ErrInvalidMode) {
				fmt.Fprintln(os.Stderr, "invalid mode: must be afk or hitl")
				os.Exit(1)
			}
			return err
		}

		fmt.Printf("updated %s\n", args[0])
		return nil
	},
}

func init() {
	updateCmd.Flags().String("title", "", "New title")
	updateCmd.Flags().Int("priority", 0, "New priority (0–3)")
	updateCmd.Flags().String("mode", "", "New mode (afk or hitl)")
	updateCmd.Flags().String("tags", "", "Comma-separated tags (replaces existing)")
	rootCmd.AddCommand(updateCmd)
}
