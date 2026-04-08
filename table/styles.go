package table

import (
	"charm.land/lipgloss/v2"
	"github.com/felipeospina21/tuishell/style"
)

// ThemedStyles returns table styles using the given theme.
func ThemedStyles(t style.Theme) Styles {
	return Styles{
		Header: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(t.PrimaryBright).
			BorderBottom(true).
			Bold(false).
			Padding(0, 1),
		Selected: lipgloss.NewStyle().
			Foreground(t.PrimaryFg).
			Background(t.PrimaryDim).
			Bold(false),
		Cell: lipgloss.NewStyle().Padding(0, 1),
	}
}

// TitleStyle returns the table title style for the given theme.
func TitleStyle(t style.Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		Margin(0, 0, 0, 1).
		Foreground(t.Primary).
		Bold(true)
}

// EmptyMsg is the style for the empty table message.
var EmptyMsg = lipgloss.NewStyle().Align(lipgloss.Center, lipgloss.Center)

// DocStyle is the outer border style for the table container.
func DocStyle(t style.Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(t.Border)
}
