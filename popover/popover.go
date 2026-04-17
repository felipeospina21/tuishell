// Package popover provides lightweight overlay components — compact alternatives to the full modal.
package popover

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/felipeospina21/tuishell/style"
)

func boxStyle(t style.Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.ModalBorder).
		Padding(0, 1)
}

func headerStyle(t style.Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(t.StatusText).
		Background(t.StatusNormal).
		Padding(0, 1).
		MarginBottom(1)
}

func dimStyle(t style.Theme) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.Dim)
}

// overlay renders fg centered over a dimmed bg.
func overlay(t style.Theme, fg, bg string, w, h int) string {
	dim := dimStyle(t)
	dimmed := dimContent(dim, bg, w, h)
	return placeOverlay(dim, w, h, fg, dimmed)
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

	startY := max((h-fgH)/2, 0)
	startX := max((w-fgW)/2, 0)

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
