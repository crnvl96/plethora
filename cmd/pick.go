package cmd

import (
	"fmt"

	"github.com/crnvl96/plethora/internal/ideas"
	"github.com/crnvl96/plethora/internal/ui"
	"github.com/spf13/cobra"
)

func init() {
	pickCmd := &cobra.Command{
		Use:   "pick [idea]",
		Short: "Run a specific idea",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				idea := args[0]
				if fn, ok := ideas.Ideas[idea]; ok {
					fn()
				} else {
					fmt.Printf("Unknown idea: %s\n", idea)
				}
			} else {
				ui.Picker(ui.PickerParams{Items: []string{"a", "b"}, Title: "Test"})
			}
		},
	}
	rootCmd.AddCommand(pickCmd)
}
