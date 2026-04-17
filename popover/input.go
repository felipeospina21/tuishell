package popover

import (
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/felipeospina21/tuishell"
	"github.com/felipeospina21/tuishell/style"
)

// InputModel is a compact input popover overlay.
type InputModel struct {
	Header string
	Input  textinput.Model
	theme  style.Theme
	width  int
}

// NewInput creates a new input popover.
func NewInput(t style.Theme) InputModel {
	ti := textinput.New()
	ti.CharLimit = 256
	return InputModel{theme: t, Input: ti, width: 40}
}

// Open configures and focuses the input popover.
func (m *InputModel) Open(header, placeholder string) {
	m.Header = header
	m.Input.Placeholder = placeholder
	m.Input.SetValue("")
	m.Input.Focus()
}

// Update handles key events for the input popover.
func (m InputModel) Update(msg tea.Msg) (InputModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg { return tuishell.CloseInputPopoverMsg{} }
		case "enter":
			v := m.Input.Value()
			return m, func() tea.Msg { return tuishell.SubmitInputPopoverMsg{Value: v} }
		}
	}
	var cmd tea.Cmd
	m.Input, cmd = m.Input.Update(msg)
	return m, cmd
}

// View renders the input popover over the given background.
func (m InputModel) View(bg string, screenW, screenH int) string {
	t := m.theme
	w := min(m.width, screenW-4)

	header := headerStyle(t).Width(w).Render(m.Header)
	m.Input.SetWidth(w - boxStyle(t).GetHorizontalFrameSize())
	input := m.Input.View()

	content := lipgloss.JoinVertical(0, header, input)
	box := boxStyle(t).Width(w).Render(content)

	lines := strings.Split(bg, "\n")
	bgH := len(lines)
	if bgH == 0 {
		bgH = screenH
	}
	return overlay(t, box, bg, screenW, bgH)
}
