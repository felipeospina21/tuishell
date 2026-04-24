package popover

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/felipeospina21/tuishell"
	"github.com/felipeospina21/tuishell/style"
)

// FilterModel is a compact multi-select filter popover overlay with sections and text inputs.
type FilterModel struct {
	Sections     []tuishell.FilterSection
	Inputs       []tuishell.FilterInput
	inputModels  []textinput.Model
	activeField  int // index across all fields (sections + inputs)
	activeCursor int // cursor within active checkbox section
	theme        style.Theme
	open         bool
}

// NewFilter creates a new filter popover.
func NewFilter(t style.Theme) FilterModel {
	return FilterModel{theme: t}
}

func (m FilterModel) totalFields() int {
	return len(m.Sections) + len(m.Inputs)
}

func (m FilterModel) isInputField() bool {
	return m.activeField >= len(m.Sections)
}

func (m FilterModel) inputIndex() int {
	return m.activeField - len(m.Sections)
}

// Open configures the filter popover with sections and input fields.
func (m *FilterModel) Open(sections []tuishell.FilterSection, inputs []tuishell.FilterInput) {
	m.Sections = sections
	m.Inputs = inputs
	m.inputModels = make([]textinput.Model, len(inputs))
	for i, inp := range inputs {
		ti := textinput.New()
		ti.Placeholder = inp.Placeholder
		ti.SetValue(inp.Value)
		ti.CharLimit = 64
		m.inputModels[i] = ti
	}
	m.activeField = 0
	m.activeCursor = 0
	m.open = true
}

// Close hides the filter popover.
func (m *FilterModel) Close() { m.open = false }

// IsOpen reports whether the filter popover is visible.
func (m FilterModel) IsOpen() bool { return m.open }

func (m FilterModel) helpText() string {
	// Always use the longest variant so the popover size stays constant.
	full := "enter apply · esc cancel · tab next · j/k navigate · space toggle"
	if m.isInputField() {
		short := "enter apply · esc cancel · tab next"
		// Pad to match full length
		return short + strings.Repeat(" ", max(0, len(full)-len(short)))
	}
	return full
}

// Update handles key events for the filter popover.
func (m FilterModel) Update(msg tea.Msg) (FilterModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		key := msg.String()

		// Always intercepted keys
		switch key {
		case "esc":
			return m, func() tea.Msg { return tuishell.CloseFilterPopoverMsg{} }
		case "enter":
			sel := make(map[string][]string)
			for _, s := range m.Sections {
				for _, o := range s.Options {
					if o.Selected {
						sel[s.Title] = append(sel[s.Title], o.Value)
					}
				}
			}
			inputs := make(map[string]string)
			for i, inp := range m.Inputs {
				inputs[inp.Title] = m.inputModels[i].Value()
			}
			return m, func() tea.Msg {
				return tuishell.ApplyFilterPopoverMsg{Selections: sel, Inputs: inputs}
			}
		case "tab":
			if m.totalFields() > 0 {
				if m.isInputField() {
					m.inputModels[m.inputIndex()].Blur()
				}
				m.activeField = (m.activeField + 1) % m.totalFields()
				if m.isInputField() {
					m.inputModels[m.inputIndex()].Focus()
				} else if m.activeField < len(m.Sections) {
					m.activeCursor = min(m.activeCursor, len(m.Sections[m.activeField].Options)-1)
				}
			}
			return m, nil
		case "shift+tab":
			if m.totalFields() > 0 {
				if m.isInputField() {
					m.inputModels[m.inputIndex()].Blur()
				}
				m.activeField = (m.activeField - 1 + m.totalFields()) % m.totalFields()
				if m.isInputField() {
					m.inputModels[m.inputIndex()].Focus()
				} else if m.activeField < len(m.Sections) {
					m.activeCursor = min(m.activeCursor, len(m.Sections[m.activeField].Options)-1)
				}
			}
			return m, nil
		}

		// Field-specific keys
		if m.isInputField() {
			// Forward to active textinput
			idx := m.inputIndex()
			var cmd tea.Cmd
			m.inputModels[idx], cmd = m.inputModels[idx].Update(msg)
			return m, cmd
		}

		// Checkbox section keys
		switch key {
		case "space":
			if m.activeField < len(m.Sections) {
				opts := m.Sections[m.activeField].Options
				if m.activeCursor < len(opts) {
					opts[m.activeCursor].Selected = !opts[m.activeCursor].Selected
					m.Sections[m.activeField].Options = opts
				}
			}
		case "j", "down":
			if m.activeField < len(m.Sections) {
				if m.activeCursor < len(m.Sections[m.activeField].Options)-1 {
					m.activeCursor++
				}
			}
		case "k", "up":
			if m.activeCursor > 0 {
				m.activeCursor--
			}
		}
	}
	return m, nil
}

// View renders the filter popover over the given background.
func (m FilterModel) View(bg string, screenW, screenH int) string {
	t := m.theme
	w := min(60, screenW-4)

	header := headerStyle(t).Render("Filters")

	titleStyle := lipgloss.NewStyle().Foreground(t.Primary).Bold(true)
	activeTitleStyle := lipgloss.NewStyle().Foreground(t.PrimaryBright).Bold(true).Italic(true)
	focused := lipgloss.NewStyle().Foreground(t.PrimaryBright).Italic(true)
	normal := lipgloss.NewStyle().Foreground(t.Text)

	// Render checkbox sections
	sectionBlocks := make([]string, len(m.Sections))
	colW := (w - 2) / 2
	for i, sec := range m.Sections {
		var lines []string
		if i == m.activeField {
			lines = append(lines, activeTitleStyle.Render(sec.Title))
		} else {
			lines = append(lines, titleStyle.Render(sec.Title))
		}
		for j, opt := range sec.Options {
			check := "[ ]"
			if opt.Selected {
				check = "[x]"
			}
			label := fmt.Sprintf("%s %s", check, opt.Label)
			if i == m.activeField && j == m.activeCursor {
				lines = append(lines, focused.Render(label))
			} else {
				lines = append(lines, normal.Render(label))
			}
		}
		sectionBlocks[i] = strings.Join(lines, "\n")
	}

	// Layout sections in 2-column grid
	var rows []string
	for i := 0; i < len(sectionBlocks); i += 2 {
		if i+1 < len(sectionBlocks) {
			left := lipgloss.NewStyle().Width(colW).Render(sectionBlocks[i])
			right := lipgloss.NewStyle().Width(colW).Render(sectionBlocks[i+1])
			rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, left, right))
		} else {
			rows = append(rows, sectionBlocks[i])
		}
	}

	sectionsView := strings.Join(rows, "\n\n")

	// Render input fields below sections
	var inputLines []string
	dimStyle := lipgloss.NewStyle().Foreground(t.Dim)
	for i, inp := range m.Inputs {
		fieldIdx := len(m.Sections) + i
		isActive := m.activeField == fieldIdx

		var label string
		if isActive {
			label = activeTitleStyle.Render(inp.Title + ": ")
		} else {
			label = titleStyle.Render(inp.Title + ": ")
		}

		var value string
		if isActive {
			value = m.inputModels[i].View()
		} else {
			v := m.inputModels[i].Value()
			if v == "" {
				value = dimStyle.Render(inp.Placeholder)
			} else {
				value = normal.Render(v)
			}
		}
		inputLines = append(inputLines, label+value)
	}

	help := lipgloss.NewStyle().Foreground(t.Dim).MarginTop(1).Render(m.helpText())

	var parts []string
	parts = append(parts, header)
	if len(sectionsView) > 0 {
		parts = append(parts, sectionsView)
	}
	if len(inputLines) > 0 {
		inputBlock := lipgloss.NewStyle().MarginTop(1).Render(strings.Join(inputLines, "\n"))
		parts = append(parts, inputBlock)
	}
	parts = append(parts, help)

	body := lipgloss.JoinVertical(0, parts...)
	box := boxStyle(t).Width(w).Render(body)

	lines := strings.Split(bg, "\n")
	bgH := len(lines)
	if bgH == 0 {
		bgH = screenH
	}
	return overlay(t, box, bg, screenW, bgH)
}
