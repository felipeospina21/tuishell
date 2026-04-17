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
	theme   style.Theme
	input   popover.InputModel
	open    bool
	result  string
	width   int
	height  int
}

func newModel() model {
	t := style.DefaultTheme()
	return model{theme: t, input: popover.NewInput(t)}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyPressMsg:
		if !m.open {
			switch msg.String() {
			case "i":
				m.open = true
				m.input.Open("Enter your name", "Type here...")
				return m, m.input.Input.Focus()
			case "q", "ctrl+c":
				return m, tea.Quit
			}
			return m, nil
		}
	case tuishell.CloseInputPopoverMsg:
		m.open = false
		return m, nil
	case tuishell.SubmitInputPopoverMsg:
		m.open = false
		m.result = msg.Value
		return m, nil
	}

	if m.open {
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m model) View() tea.View {
	bg := fmt.Sprintf("\n  Press 'i' to open input popover, 'q' to quit\n\n  Last input: %s\n", m.result)
	bg = lipgloss.NewStyle().Width(m.width).Height(m.height).Render(bg)

	screen := bg
	if m.open {
		screen = m.input.View(bg, m.width, m.height)
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
