package cmd

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/crnvl96/plethora/internal/ideas"
	"github.com/crnvl96/plethora/internal/ui"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "pick",
		Short: "Select an idea to run from a list",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			l := []list.Item{}
			for _, i := range ideas.Ideas {
				l = append(l, &i)
			}
			ui.RenderPicker(ui.PickerParams{Items: l, Title: "A Plethora of programming ideas"})
		},
	})
}
