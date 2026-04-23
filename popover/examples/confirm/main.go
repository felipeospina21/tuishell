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
	confirm popover.ConfirmModel
	open    bool
	result  string
	width   int
	height  int
}

func newModel() model {
	t := defaultTheme()
	return model{theme: t, confirm: popover.NewConfirm(t)}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyPressMsg:
		if !m.open {
			switch msg.String() {
			case "c":
				m.open = true
				m.confirm.Open("Delete Item", "Are you sure you want to delete this item?", "Delete", "Cancel")
				return m, nil
			case "q", "ctrl+c":
				return m, tea.Quit
			}
			return m, nil
		}
	case tuishell.ConfirmPopoverYesMsg:
		m.open = false
		m.result = "Confirmed!"
		return m, nil
	case tuishell.ConfirmPopoverNoMsg:
		m.open = false
		m.result = "Cancelled"
		return m, nil
	}

	if m.open {
		var cmd tea.Cmd
		m.confirm, cmd = m.confirm.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m model) View() tea.View {
	bg := fmt.Sprintf("\n  Press 'c' to open confirm popover, 'q' to quit\n\n  Result: %s\n", m.result)
	bg = lipgloss.NewStyle().Width(m.width).Height(m.height).Render(bg)

	screen := bg
	if m.open {
		screen = m.confirm.View(bg, m.width, m.height)
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
