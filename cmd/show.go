package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ro-56/taskr/internal/tickets"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Display full ticket info",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		result, err := tickets.Show(cwd, args[0])
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

		printTicket(result)
		return nil
	},
}

func printTicket(r *tickets.ShowResult) {
	// frontmatter fields in a readable order
	fields := []string{"id", "title", "status", "type", "priority", "mode", "created", "updated", "tags", "dependencies", "links"}
	for _, f := range fields {
		v, ok := r.Frontmatter[f]
		if !ok {
			continue
		}
		fmt.Printf("%-14s %v\n", f+":", formatField(v))
	}

	if r.Body != "" {
		fmt.Println()
		fmt.Print(r.Body)
	}

	fmt.Println()
	fmt.Printf("depends on:  %s\n", formatDepRefs(r.DependsOn))
	fmt.Printf("required by: %s\n", formatDepRefs(r.RequiredBy))
}

func formatField(v any) string {
	switch val := v.(type) {
	case []any:
		if len(val) == 0 {
			return "none"
		}
		parts := make([]string, len(val))
		for i, item := range val {
			parts[i] = fmt.Sprint(item)
		}
		return strings.Join(parts, ", ")
	default:
		return fmt.Sprint(v)
	}
}

func formatDepRefs(refs []tickets.DepRef) string {
	if len(refs) == 0 {
		return "none"
	}
	parts := make([]string, len(refs))
	for i, r := range refs {
		parts[i] = fmt.Sprintf("%s (%s)", r.ID, r.Status)
	}
	return strings.Join(parts, ", ")
}

func init() {
	rootCmd.AddCommand(showCmd)
}
