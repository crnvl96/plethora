// Package ui
package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
)

const (
	blueberry    = "99"
	lavenderRose = "212"
)

const (
	marginRight = 1
)

func List(items []string) *list.List {
	enumerator := list.Arabic
	enumeratorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(blueberry)).MarginRight(marginRight)
	itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(lavenderRose)).MarginRight(marginRight)

	return list.New(items).
		Enumerator(enumerator).
		EnumeratorStyle(enumeratorStyle).
		ItemStyle(itemStyle)
}
