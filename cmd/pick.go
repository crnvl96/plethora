package cmd

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
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
			ideasList := []list.Item{}

			for _, idea := range ideas.Ideas {
				ideasList = append(ideasList, ui.Item{Header: idea.Header, Body: idea.Body, Command: idea.Command, Callback: idea.Callback})
			}

			if len(args) == 1 {
				arg := args[0]
				if idea, ok := ideas.Ideas[arg]; ok {
					tea.Exec(idea.Command, idea.Callback)
				} else {
					fmt.Printf("Unknown idea: %s\n", arg)
				}
			} else {
				ui.RenderPicker(ui.RenderParams{
					Items: ideasList,
					Title: "A Plethora of programming ideas",
				})
			}
		},
	}
	rootCmd.AddCommand(pickCmd)
}
