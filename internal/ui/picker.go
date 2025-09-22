// Package ui provides user interface components for the plethora application.
package ui

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Constants for styling.
const (
	moderateMagenta = "170"
	listHeight      = 14
	defaultWidth    = 20
)

// Styles for the list UI.
var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color(moderateMagenta))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

// item represents a selectable item in the list.
type item string

// FilterValue implements the list.Item interface.
func (i item) FilterValue() string { return "" }

// itemDelegate handles rendering of list items.
type itemDelegate struct{}

// Height returns the height of each item.
func (d itemDelegate) Height() int { return 1 }

// Spacing returns the spacing between items.
func (d itemDelegate) Spacing() int { return 0 }

// Update handles updates for the delegate (no-op).
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

// Render renders an individual list item.
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

// model holds the state of the picker UI.
type model struct {
	list     list.Model
	choice   string
	quitting bool
}

// Init initializes the model (no-op).
func (m model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i)
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders the current view of the model.
func (m model) View() string {
	if m.choice != "" {
		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice))
	}
	if m.quitting {
		return quitTextStyle.Render("Not hungry? That’s cool.")
	}
	return "\n" + m.list.View()
}

// PickerParams defines the arguments that must be sent to Picker
type PickerParams struct {
	Items []string
	Title string
}

// Picker runs the picker UI with hardcoded items.
func Picker(p PickerParams) {
	listItems := []list.Item{}
	for _, i := range p.Items {
		listItems = append(listItems, item(i))
	}

	l := list.New(listItems, itemDelegate{}, defaultWidth, listHeight)
	l.Title = p.Title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model{list: l}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
