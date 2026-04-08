// Package modal implements a centered overlay modal component.
package modal

import (
	"strings"

	"charm.land/bubbles/v2/help"
	"charm.land/lipgloss/v2"
	"github.com/felipeospina21/tuishell"
	"github.com/felipeospina21/tuishell/style"
)

// Model holds the state for the modal overlay.
type Model struct {
	Header    string
	Content   string
	Highlight bool
	IsError   bool
	ctx       *tuishell.AppContext
	theme     style.Theme
}

// New creates a new modal model.
func New(ctx *tuishell.AppContext, theme style.Theme) Model {
	return Model{ctx: ctx, theme: theme}
}

// View renders the modal box centered over the dimmed background.
func (m Model) View(background string) string {
	t := m.theme
	// Use actual background dimensions for proper centering
	bgLines := strings.Split(background, "\n")
	w := lipgloss.Width(background)
	h := len(bgLines)
	if w == 0 {
		w = m.ctx.Window.Width
	}
	if h == 0 {
		h = m.ctx.Window.Height
	}

	modalW := modalSize(w)
	modalH := modalSize(h)

	hdrStyle := headerStyle(t)
	header := hdrStyle.Width(modalW).Render(m.Header)
	if m.IsError {
		header = hdrStyle.Background(t.StatusError).Width(modalW).Render(m.Header)
	}
	ftr := helpStyle(t)
	footer := ftr.Render("esc close · ctrl+s submit · ctrl+y copy")
	box := boxStyle(t)
	contentH := max(modalH-lipgloss.Height(header)-lipgloss.Height(footer)-box.GetVerticalFrameSize(), 1)
	contentW := modalW - box.GetHorizontalFrameSize()

	content := m.Content
	if m.Highlight {
		content = lipgloss.NewStyle().
			Background(t.PrimaryBright).
			Foreground(t.TextInverse).
			Render(content)
	}

	body := lipgloss.NewStyle().
		Width(contentW).
		Height(contentH).
		MaxHeight(contentH).
		Render(content)

	rendered := box.Width(modalW).Height(modalH).Render(
		lipgloss.JoinVertical(0, header, body, footer),
	)

	// Dim the background and overlay the modal
	dim := dimStyle(t)
	dimmed := dimContent(dim, background, w, h)
	return placeOverlay(dim, w, h, rendered, dimmed)
}

// SetFocus sets the focused panel to the modal.
func (m *Model) SetFocus() {
	m.ctx.FocusedPanel = tuishell.ModalPanel
}

// RenderHelp renders a full help view with styles suited for the modal.
func (m Model) RenderHelp(km help.KeyMap) string {
	t := m.theme
	h := help.New()
	h.Styles.FullKey = lipgloss.NewStyle().Foreground(t.Text)
	h.Styles.FullDesc = lipgloss.NewStyle().Foreground(t.Muted)
	h.Styles.FullSeparator = lipgloss.NewStyle()
	return h.FullHelpView(km.FullHelp())
}

func modalSize(available int) int {
	switch {
	case available < 40:
		return available - 2
	case available < 80:
		return available * 85 / 100
	default:
		return available * 3 / 4
	}
}

// ContentWidth returns the usable content width inside the modal.
func ContentWidth(t style.Theme, windowW int) int {
	return modalSize(windowW) - boxStyle(t).GetHorizontalFrameSize()
}

// ContentHeight returns the usable content height inside the modal.
func ContentHeight(t style.Theme, windowH int) int {
	modalH := modalSize(windowH)
	header := headerStyle(t).Render("X")
	footer := helpStyle(t).Render("X")
	return max(modalH-lipgloss.Height(header)-lipgloss.Height(footer)-boxStyle(t).GetVerticalFrameSize(), 1)
}

func dimContent(dim lipgloss.Style, s string, w, h int) string {
	bg := lipgloss.NewStyle().Width(w).Height(h).Render(s)
	lines := strings.Split(bg, "\n")
	for i, line := range lines {
		lines[i] = dim.Render(stripAnsi(line))
	}
	return strings.Join(lines, "\n")
}

func placeOverlay(dim lipgloss.Style, w, h int, fg, bg string) string {
	fgLines := strings.Split(fg, "\n")
	bgLines := strings.Split(bg, "\n")

	fgW := lipgloss.Width(fg)
	fgH := len(fgLines)

	startY := (h - fgH) / 2
	startX := (w - fgW) / 2
	if startY < 0 {
		startY = 0
	}
	if startX < 0 {
		startX = 0
	}

	for i, fgLine := range fgLines {
		bgIdx := startY + i
		if bgIdx >= len(bgLines) {
			break
		}
		bgRunes := []rune(stripAnsi(bgLines[bgIdx]))

		var prefix string
		if startX > 0 && startX <= len(bgRunes) {
			prefix = dim.Render(string(bgRunes[:startX]))
		} else {
			prefix = strings.Repeat(" ", startX)
		}

		endX := startX + lipgloss.Width(fgLine)
		var suffix string
		if endX < len(bgRunes) {
			suffix = dim.Render(string(bgRunes[endX:]))
		}

		bgLines[bgIdx] = prefix + fgLine + suffix
	}

	return strings.Join(bgLines, "\n")
}

func stripAnsi(s string) string {
	var out strings.Builder
	inEsc := false
	for _, r := range s {
		if r == '\x1b' {
			inEsc = true
			continue
		}
		if inEsc {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '~' {
				inEsc = false
			}
			continue
		}
		out.WriteRune(r)
	}
	return out.String()
}
