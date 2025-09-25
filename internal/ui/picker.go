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
	Item struct {
		Header, Body string
		Command      tea.ExecCommand
		Callback     tea.ExecCallback
	}

	ExecDoneMsg struct {
		Err error
	}

	RenderParams struct {
		Items []list.Item
		Title string
	}

	model struct {
		list list.Model
	}
)

func (i Item) Title() string       { return i.Header }
func (i Item) Description() string { return i.Body }
func (i Item) FilterValue() string { return i.Body }

func (m model) Init() tea.Cmd { return nil }
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		if msg.String() == "enter" {
			i, ok := m.list.SelectedItem().(Item)

			if ok {
				cmd := tea.Exec(i.Command, i.Callback)
				return m, cmd
			}

		}

	case ExecDoneMsg:
		return m, tea.Quit

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func RenderPicker(p RenderParams) {
	m := model{list: list.New(p.Items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = p.Title

	pg := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := pg.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
