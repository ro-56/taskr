package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/ro-56/taskr/internal/tickets"
	"github.com/spf13/cobra"
)

var depTreeFull bool

var depTreeCmd = &cobra.Command{
	Use:   "dep-tree <id>",
	Short: "Render ASCII tree of a ticket's dependencies",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		root, err := tickets.DepTree(cwd, args[0], depTreeFull)
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
			return err
		}

		fmt.Print(tickets.RenderDepTree(root))
		return nil
	},
}

func init() {
	depTreeCmd.Flags().BoolVar(&depTreeFull, "full", false, "recursively expand all dependencies")
	rootCmd.AddCommand(depTreeCmd)
}
