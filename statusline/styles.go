package statusline

import (
	"charm.land/lipgloss/v2"
	"github.com/felipeospina21/tuishell/style"
)

// StatusBarStyle returns the outer status bar style.
func StatusBarStyle() lipgloss.Style {
	return lipgloss.NewStyle().Margin(0, 0)
}

// SpinnerStyle returns the spinner style for the given theme.
func SpinnerStyle(t style.Theme) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.PrimaryBright)
}

func statusText(t style.Theme) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.StatusText)
}

func statusStyle(t style.Theme) lipgloss.Style {
	return statusText(t).Padding(0, 1).MarginRight(1)
}

func encodingStyle(t style.Theme) lipgloss.Style {
	return statusText(t).Padding(0, 1).Background(t.StatusAccent1).Align(lipgloss.Right)
}

func projectStyle(t style.Theme) lipgloss.Style {
	return statusText(t).Padding(0, 1).Background(t.StatusAccent2)
}

func helpText() lipgloss.Style {
	return lipgloss.NewStyle().AlignHorizontal(lipgloss.Center)
}
