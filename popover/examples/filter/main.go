package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/felipeospina21/tuishell"
	"github.com/felipeospina21/tuishell/popover"
	"github.com/felipeospina21/tuishell/style"
)

type model struct {
	theme  style.Theme
	filter popover.FilterModel
	open   bool
	result string
	width  int
	height int
}

func newModel() model {
	t := defaultTheme()
	return model{theme: t, filter: popover.NewFilter(t)}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyPressMsg:
		if !m.open {
			switch msg.String() {
			case "f":
				m.open = true
				m.filter.Open(
					[]tuishell.FilterSection{
						{Title: "Status", Options: []tuishell.FilterOption{
							{Label: "In Progress", Value: "in_progress", Selected: true},
							{Label: "To Do", Value: "todo"},
							{Label: "In Review", Value: "in_review"},
						}},
						{Title: "Priority", Options: []tuishell.FilterOption{
							{Label: "Critical", Value: "critical"},
							{Label: "High", Value: "high", Selected: true},
							{Label: "Medium", Value: "medium"},
							{Label: "Low", Value: "low"},
						}},
						{Title: "Type", Options: []tuishell.FilterOption{
							{Label: "Bug", Value: "bug", Selected: true},
							{Label: "Story", Value: "story"},
							{Label: "Task", Value: "task"},
						}},
					},
					[]tuishell.FilterInput{
						{Title: "Sprint", Placeholder: "e.g. 42"},
					},
				)
				return m, nil
			case "q", "ctrl+c":
				return m, tea.Quit
			}
			return m, nil
		}
	case tuishell.CloseFilterPopoverMsg:
		m.open = false
		return m, nil
	case tuishell.ApplyFilterPopoverMsg:
		m.open = false
		var parts []string
		for title, vals := range msg.Selections {
			parts = append(parts, fmt.Sprintf("%s: %s", title, strings.Join(vals, ", ")))
		}
		for title, val := range msg.Inputs {
			if val != "" {
				parts = append(parts, fmt.Sprintf("%s: %s", title, val))
			}
		}
		m.result = strings.Join(parts, " | ")
		return m, nil
	}

	if m.open {
		var cmd tea.Cmd
		m.filter, cmd = m.filter.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m model) View() tea.View {
	bg := fmt.Sprintf("\n  Press 'f' to open filter popover, 'q' to quit\n\n  Applied: %s\n", m.result)
	bg = lipgloss.NewStyle().Width(m.width).Height(m.height).Render(bg)

	screen := bg
	if m.open {
		screen = m.filter.View(bg, m.width, m.height)
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
