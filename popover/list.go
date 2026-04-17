package popover

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/felipeospina21/tuishell"
	"github.com/felipeospina21/tuishell/style"
)

// ListModel is a compact list popover overlay with j/k navigation.
type ListModel struct {
	Header string
	Items  []tuishell.ListPopoverItem
	Cursor int
	theme  style.Theme
	width  int
	open   bool
}

// NewList creates a new list popover.
func NewList(t style.Theme) ListModel {
	return ListModel{theme: t, width: 40}
}

// Open configures the list popover with items.
func (m *ListModel) Open(header string, items []tuishell.ListPopoverItem) {
	m.Header = header
	m.Items = items
	m.Cursor = 0
	m.open = true
}

// Close hides the list popover.
func (m *ListModel) Close() { m.open = false }

// IsOpen reports whether the list popover is visible.
func (m ListModel) IsOpen() bool { return m.open }

// Update handles key events for the list popover.
func (m ListModel) Update(msg tea.Msg) (ListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg { return tuishell.CloseListPopoverMsg{} }
		case "enter":
			if len(m.Items) > 0 {
				v := m.Items[m.Cursor].Value
				return m, func() tea.Msg { return tuishell.SelectListPopoverMsg{Value: v} }
			}
		case "j", "down":
			if m.Cursor < len(m.Items)-1 {
				m.Cursor++
			}
		case "k", "up":
			if m.Cursor > 0 {
				m.Cursor--
			}
		}
	}
	return m, nil
}

// View renders the list popover over the given background.
func (m ListModel) View(bg string, screenW, screenH int) string {
	t := m.theme
	w := min(m.width, screenW-4)

	header := headerStyle(t).Width(w).Render(m.Header)

	selected := lipgloss.NewStyle().Foreground(t.Primary).Bold(true)
	normal := lipgloss.NewStyle().Foreground(t.Text)

	var rows []string
	for i, item := range m.Items {
		label := fmt.Sprintf("  %s", item.Label)
		if i == m.Cursor {
			label = fmt.Sprintf("▸ %s", item.Label)
			rows = append(rows, selected.Render(label))
		} else {
			rows = append(rows, normal.Render(label))
		}
	}

	content := lipgloss.JoinVertical(0, header, strings.Join(rows, "\n"))
	box := boxStyle(t).Width(w).Render(content)

	lines := strings.Split(bg, "\n")
	bgH := len(lines)
	if bgH == 0 {
		bgH = screenH
	}
	return overlay(t, box, bg, screenW, bgH)
}
