package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type (
	DoneMsg   struct{ Err error }
	Component struct{ Title, Description string }
	model     struct{ list list.Model }

	Item struct {
		Component Component
		Command   tea.ExecCommand
		Callback  tea.ExecCallback
	}

	PickerParams struct {
		Items []list.Item
		Title string
	}
)

func (i *Item) Title() string       { return i.Component.Title }
func (i *Item) Description() string { return i.Component.Description }
func (i *Item) FilterValue() string { return i.Component.Description }

func (m *model) Init() tea.Cmd { return nil }
func (m *model) View() string  { return docStyle.Render(m.list.View()) }
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch key := msg.String(); key {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			i, ok := m.list.SelectedItem().(*Item)
			if ok {
				cmd := tea.Exec(i.Command, i.Callback)
				return m, cmd
			}
		}
	case DoneMsg:
		return m, tea.Quit
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func RenderPicker(p PickerParams) {
	m := model{list: list.New(p.Items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = p.Title

	prog := tea.NewProgram(&m, tea.WithAltScreen())
	if _, err := prog.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
