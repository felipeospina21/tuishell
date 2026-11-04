package modal

import (
	"strings"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/lipgloss/v2"
)

const (
	helpMaxKeyWidth = 20
	helpSectionRule = "──"
)

// RenderHelpScrollable renders km.FullHelp() as a single vertical column
// suitable for a scrollable viewport (or the modal, which scrolls its content).
//
// Groups are stacked top to bottom and separated by a blank line. Within a
// group, bindings are rendered as follows:
//
//   - A binding with no description whose help text starts with "──" is a
//     section header, rendered bold in the primary color.
//
//   - Any other binding with no description is a note, rendered muted/italic.
//
//   - Otherwise it is a key/description row. Keys are left-aligned in a column
//     as wide as the widest key (capped), followed by the description.
//
//   - Disabled bindings are skipped.
//
//     ── Section Name ──
//     <key>   description
//
//     ── Another Section ──
//     <key>   description
func (m Model) RenderHelpScrollable(km help.KeyMap) string {
	t := m.theme
	groups := km.FullHelp()

	keyW := 0
	for _, g := range groups {
		for _, b := range g {
			if !b.Enabled() || b.Help().Desc == "" {
				continue
			}
			keyW = max(keyW, lipgloss.Width(b.Help().Key))
		}
	}
	keyW = min(keyW, helpMaxKeyWidth)

	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(t.Primary)
	noteStyle := lipgloss.NewStyle().Foreground(t.Muted).Italic(true)
	keyStyle := lipgloss.NewStyle().Foreground(t.Text).Bold(true).Width(keyW)
	descStyle := lipgloss.NewStyle().Foreground(t.Muted)

	var lines []string
	for gi, g := range groups {
		if gi > 0 {
			lines = append(lines, "")
		}
		for _, b := range g {
			lines = append(lines, renderHelpLine(b, headerStyle, noteStyle, keyStyle, descStyle)...)
		}
	}
	return strings.Join(lines, "\n")
}

func renderHelpLine(b key.Binding, header, note, keyStyle, desc lipgloss.Style) []string {
	if !b.Enabled() {
		return nil
	}
	h := b.Help()
	switch {
	case h.Desc == "" && h.Key == "":
		return nil
	case h.Desc == "" && strings.HasPrefix(h.Key, helpSectionRule):
		return []string{header.Render(h.Key)}
	case h.Desc == "":
		return []string{note.Render(h.Key)}
	default:
		return []string{keyStyle.Render(h.Key) + "  " + desc.Render(h.Desc)}
	}
}
