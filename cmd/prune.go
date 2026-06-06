package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/ro-56/taskr/internal/tickets"
	"github.com/spf13/cobra"
)

var pruneCmd = &cobra.Command{
	Use:   "prune <id>",
	Short: "Remove ticket from all dependency lists and clear its own",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		err = tickets.Prune(cwd, args[0])
		if err != nil {
			if errors.Is(err, tickets.ErrNotFound) {
				fmt.Fprintf(os.Stderr, "ticket not found: %s\n", args[0])
				os.Exit(1)
			}
			var ambErr *tickets.ErrAmbiguous
			if errors.As(err, &ambErr) {
				fmt.Fprintln(os.Stderr, "ambiguous ID — did you mean one of:")
				for _, m := range ambErr.Matches {
					fmt.Fprintf(os.Stderr, "  %s\n", m)
				}
				os.Exit(1)
			}
			return err
		}

		fmt.Printf("pruned: %s removed from all dependency lists\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pruneCmd)
}
