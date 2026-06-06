package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/ro-56/taskr/internal/tickets"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start <id>",
	Short: "Transition a ticket to in_progress",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		err = tickets.Start(cwd, args[0])
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
			if errors.Is(err, tickets.ErrWrongStatus) {
				fmt.Fprintf(os.Stderr, "cannot start: ticket is already closed\n")
				os.Exit(1)
			}
			return err
		}

		fmt.Printf("started %s\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
