package popover

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/felipeospina21/tuishell"
	"github.com/felipeospina21/tuishell/style"
)

// ConfirmModel is a compact confirmation popover overlay with two buttons.
type ConfirmModel struct {
	Header  string
	Message string
	Confirm string
	Cancel  string
	Focused int // 0 = confirm (left), 1 = cancel (right)
	theme   style.Theme
	width   int
	open    bool
}

// NewConfirm creates a new confirmation popover.
func NewConfirm(t style.Theme) ConfirmModel {
	return ConfirmModel{theme: t, width: 40}
}

// Open configures the confirmation popover with header, message, and button labels.
func (m *ConfirmModel) Open(header, message, confirm, cancel string) {
	m.Header = header
	m.Message = message
	m.Confirm = confirm
	m.Cancel = cancel
	m.Focused = 1 // start on cancel to prevent accidental confirmations
	m.open = true
}

// Close hides the confirmation popover.
func (m *ConfirmModel) Close() { m.open = false }

// IsOpen reports whether the confirmation popover is visible.
func (m ConfirmModel) IsOpen() bool { return m.open }

// Update handles key events for the confirmation popover.
func (m ConfirmModel) Update(msg tea.Msg) (ConfirmModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc", "n":
			return m, func() tea.Msg { return tuishell.ConfirmPopoverNoMsg{} }
		case "y":
			return m, func() tea.Msg { return tuishell.ConfirmPopoverYesMsg{} }
		case "enter":
			if m.Focused == 0 {
				return m, func() tea.Msg { return tuishell.ConfirmPopoverYesMsg{} }
			}
			return m, func() tea.Msg { return tuishell.ConfirmPopoverNoMsg{} }
		case "tab", "l", "right":
			if m.Focused < 1 {
				m.Focused = 1
			}
		case "shift+tab", "h", "left":
			if m.Focused > 0 {
				m.Focused = 0
			}
		}
	}
	return m, nil
}

// View renders the confirmation popover over the given background.
func (m ConfirmModel) View(bg string, screenW, screenH int) string {
	t := m.theme
	w := min(m.width, screenW-4)

	header := headerStyle(t).Width(w).Render(m.Header)
	body := lipgloss.NewStyle().Foreground(t.Text).Width(w).Render(m.Message)

	dimBtn := lipgloss.NewStyle().Foreground(t.Dim)
	focusedConfirm := lipgloss.NewStyle().Foreground(t.Danger).Bold(true)
	focusedCancel := lipgloss.NewStyle().Foreground(t.Muted).Bold(true)

	var confirmBtn, cancelBtn string
	if m.Focused == 0 {
		confirmBtn = focusedConfirm.Render("[ " + m.Confirm + " ]")
	} else {
		confirmBtn = dimBtn.Render("[ " + m.Confirm + " ]")
	}
	if m.Focused == 1 {
		cancelBtn = focusedCancel.Render("[ " + m.Cancel + " ]")
	} else {
		cancelBtn = dimBtn.Render("[ " + m.Cancel + " ]")
	}

	buttons := lipgloss.NewStyle().MarginTop(1).Render(
		lipgloss.JoinHorizontal(lipgloss.Top, confirmBtn, " ", cancelBtn),
	)

	content := lipgloss.JoinVertical(0, header, body, buttons)
	box := boxStyle(t).Width(w).Render(content)

	lines := strings.Split(bg, "\n")
	bgH := len(lines)
	if bgH == 0 {
		bgH = screenH
	}
	return overlay(t, box, bg, screenW, bgH)
}
