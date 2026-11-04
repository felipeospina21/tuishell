package modal

import (
	"strings"

	"charm.land/lipgloss/v2"
)

// ResetScroll scrolls the modal content back to the top. The shell calls it
// whenever new content is opened.
func (m *Model) ResetScroll() { m.scroll = 0 }

// viewport returns the content width/height available inside the modal and the
// number of wrapped content lines. ok is false when window size is unknown.
func (m Model) viewport() (contentW, contentH, total int, ok bool) {
	if m.ctx == nil || m.ctx.Window.Width == 0 || m.ctx.Window.Height == 0 {
		return 0, 0, 0, false
	}
	contentW = ContentWidth(m.theme, m.ctx.Window.Width)
	contentH = ContentHeight(m.theme, m.ctx.Window.Height)
	wrapped := lipgloss.NewStyle().Width(max(contentW, 1)).Render(m.Content)
	return contentW, contentH, len(strings.Split(wrapped, "\n")), true
}

// scrollKey applies a scrolling key to m and reports whether it was one.
func (m *Model) scrollKey(k string) bool {
	_, h, total, ok := m.viewport()
	if !ok {
		return false
	}
	maxScroll := max(total-h, 0)
	switch k {
	case "j", "down":
		m.scroll++
	case "k", "up":
		m.scroll--
	case "pgdown", "ctrl+d":
		m.scroll += max(h/2, 1)
	case "pgup", "ctrl+u":
		m.scroll -= max(h/2, 1)
	case "g", "home":
		m.scroll = 0
	case "G", "end":
		m.scroll = maxScroll
	default:
		return false
	}
	m.scroll = min(max(m.scroll, 0), maxScroll)
	return true
}
