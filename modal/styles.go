package modal

import (
	"charm.land/lipgloss/v2"
	"github.com/felipeospina21/tuishell/style"
)

func headerStyle(t style.Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		Margin(0, 0, 0, 1).
		Bold(true).
		Foreground(t.StatusText).
		Background(t.StatusNormal).
		Padding(0, 1).
		MarginBottom(1)
}

func helpStyle(t style.Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.Border).
		MarginTop(1)
}

func boxStyle(t style.Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.ModalBorder).
		Padding(1)
}

func dimStyle(t style.Theme) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.Dim)
}
