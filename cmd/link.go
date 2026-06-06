package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/ro-56/taskr/internal/tickets"
	"github.com/spf13/cobra"
)

var linkCmd = &cobra.Command{
	Use:   "link <dependent-id> <depends-on-id>",
	Short: "Add a dependency: dependent waits on depends-on",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		err = tickets.Link(cwd, args[0], args[1])
		if err != nil {
			var cycleErr *tickets.ErrCycle
			switch {
			case errors.As(err, &cycleErr):
				fmt.Fprintf(os.Stderr, "error: %s\n", cycleErr.Error())
				os.Exit(1)
			case errors.Is(err, tickets.ErrSelfLink):
				fmt.Fprintln(os.Stderr, "error: a ticket cannot depend on itself")
				os.Exit(1)
			case errors.Is(err, tickets.ErrAlreadyLinked):
				fmt.Fprintf(os.Stderr, "already linked: %s depends on %s\n", args[0], args[1])
				return nil
			case errors.Is(err, tickets.ErrNotFound):
				fmt.Fprintf(os.Stderr, "ticket not found\n")
				os.Exit(1)
			default:
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
		}

		fmt.Printf("linked: %s depends on %s\n", args[0], args[1])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(linkCmd)
}
