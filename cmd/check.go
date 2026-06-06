package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/ro-56/taskr/internal/tickets"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Validate all tickets and report invariant violations",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		violations, err := tickets.Check(cwd)
		if err != nil {
			if errors.Is(err, tickets.ErrNotInitialized) {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return err
		}

		for _, v := range violations {
			fmt.Printf("%s: %s\n", v.TicketID, v.Message)
		}

		if len(violations) > 0 {
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
