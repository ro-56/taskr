package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/ro-56/taskr/internal/tickets"
	"github.com/spf13/cobra"
)

var unlinkCmd = &cobra.Command{
	Use:   "unlink <dependent-id> <depends-on-id>",
	Short: "Remove a dependency link",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		err = tickets.Unlink(cwd, args[0], args[1])
		if err != nil {
			if errors.Is(err, tickets.ErrNotFound) {
				fmt.Fprintln(os.Stderr, "ticket not found")
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

		fmt.Printf("unlinked: %s no longer depends on %s\n", args[0], args[1])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(unlinkCmd)
}
