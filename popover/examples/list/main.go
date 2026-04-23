package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/felipeospina21/tuishell"
	"github.com/felipeospina21/tuishell/popover"
	"github.com/felipeospina21/tuishell/style"
)

type model struct {
	theme  style.Theme
	list   popover.ListModel
	open   bool
	result string
	width  int
	height int
}

func newModel() model {
	t := defaultTheme()
	return model{theme: t, list: popover.NewList(t)}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyPressMsg:
		if !m.open {
			switch msg.String() {
			case "l":
				m.open = true
				m.list.Open("Select a color", []tuishell.ListPopoverItem{
					{Label: "Red", Value: "red"},
					{Label: "Green", Value: "green"},
					{Label: "Blue", Value: "blue"},
					{Label: "Yellow", Value: "yellow"},
					{Label: "Violet", Value: "violet"},
					{Label: "Orange", Value: "orange"},
				})
				return m, nil
			case "q", "ctrl+c":
				return m, tea.Quit
			}
			return m, nil
		}
	case tuishell.CloseListPopoverMsg:
		m.open = false
		return m, nil
	case tuishell.SelectListPopoverMsg:
		m.open = false
		m.result = msg.Value
		return m, nil
	}

	if m.open {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m model) View() tea.View {
	bg := fmt.Sprintf("\n  Press 'l' to open list popover, 'q' to quit\n\n  Selected: %s\n", m.result)
	bg = lipgloss.NewStyle().Width(m.width).Height(m.height).Render(bg)

	screen := bg
	if m.open {
		screen = m.list.View(bg, m.width, m.height)
	}

	v := tea.NewView(screen)
	v.AltScreen = true
	return v
}

func main() {
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
