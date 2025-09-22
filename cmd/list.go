package cmd

import (
	"fmt"

	"github.com/crnvl96/plethora/internal/ideas"
	"github.com/crnvl96/plethora/internal/ui"
	"github.com/spf13/cobra"
)

func init() {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all ideas",
		Run: func(cmd *cobra.Command, args []string) {
			ideasList := []string{}
			for name := range ideas.Ideas {
				ideasList = append(ideasList, name)
			}
			fmt.Print(ui.List(ideasList))
		},
	}

	rootCmd.AddCommand(listCmd)
}
