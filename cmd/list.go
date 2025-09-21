package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/spf13/cobra"
)

func init() {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all files under internal/",
		Run: func(cmd *cobra.Command, args []string) {
			paths := []string{}

			err := filepath.Walk("internal", func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					paths = append(paths, path)
				}
				return nil
			})
			if err != nil {
				fmt.Println(err)
			}

			enumeratorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("99")).MarginRight(1)
			itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("212")).MarginRight(1)

			l := list.New(paths).
				Enumerator(list.Roman).
				EnumeratorStyle(enumeratorStyle).
				ItemStyle(itemStyle)

			fmt.Println(l)
		},
	}

	rootCmd.AddCommand(listCmd)
}
