package modal

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/felipeospina21/tuishell"
)

// Re-export message types from the root package for backward compatibility.
type (
	CloseModalMsg     = tuishell.CloseModalMsg
	SubmitModalMsg    = tuishell.SubmitModalMsg
	CopyModalMsg      = tuishell.CopyModalMsg
	ResetHighlightMsg = tuishell.ResetHighlightMsg
)

// Update handles key presses for close, submit, and copy actions.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		match := tuishell.KeyMatcher(msg)
		switch {
		case match(Keybinds.Close):
			return m, func() tea.Msg { return tuishell.CloseModalMsg{} }
		case match(Keybinds.Submit):
			return m, func() tea.Msg { return tuishell.SubmitModalMsg{} }
		case match(Keybinds.Copy):
			m.Highlight = true
			return m, tea.Batch(
				func() tea.Msg { return tuishell.CopyModalMsg{} },
				tea.Tick(100*time.Millisecond, func(time.Time) tea.Msg { return tuishell.ResetHighlightMsg{} }),
			)
		}
	case tuishell.ResetHighlightMsg:
		m.Highlight = false
		return m, nil
	}
	return m, nil
}
