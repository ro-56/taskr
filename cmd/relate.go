package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/ro-56/taskr/internal/tickets"
	"github.com/spf13/cobra"
)

var relateCmd = &cobra.Command{
	Use:   "relate <id-1> <id-2> [id-3 ...]",
	Short: "Create bidirectional related links between all provided tickets",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		err = tickets.Relate(cwd, args...)
		if err != nil {
			switch {
			case errors.Is(err, tickets.ErrSelfRelate):
				fmt.Fprintln(os.Stderr, "error: a ticket cannot be related to itself")
				os.Exit(1)
			case errors.Is(err, tickets.ErrAlreadyRelated):
				fmt.Fprintln(os.Stderr, "already related: all pairs already linked")
				return nil
			case errors.Is(err, tickets.ErrNotFound):
				fmt.Fprintln(os.Stderr, "error: ticket not found")
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

		fmt.Printf("related: %s\n", joinArgs(args))
		return nil
	},
}

func joinArgs(args []string) string {
	out := args[0]
	for _, a := range args[1:] {
		out += " <-> " + a
	}
	return out
}

func init() {
	rootCmd.AddCommand(relateCmd)
}
