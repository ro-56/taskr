package cmd

import (
	"fmt"
	"os"

	"github.com/ro-56/taskr/internal/tickets"
	"github.com/spf13/cobra"
)

var initPrefix string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a taskr project in the current directory",
	Long: `Creates .tickets/, .tickets/archive/, and .tickets/config.json.
The PREFIX defaults to an uppercased slug of the current directory name.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		err = tickets.Init(cwd, initPrefix)
		if err == tickets.ErrAlreadyInitialized {
			fmt.Println("already initialized")
			return nil
		}
		if err != nil {
			return err
		}

		fmt.Println("initialized")
		return nil
	},
}

func init() {
	initCmd.Flags().StringVar(&initPrefix, "prefix", "", "project prefix (default: uppercased slug of directory name)")
	rootCmd.AddCommand(initCmd)
}
