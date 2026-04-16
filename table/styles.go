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

// EmptyMsg is the style applied when the table has no rows.
var EmptyMsg = lipgloss.NewStyle().Align(lipgloss.Center, lipgloss.Center)

// DocStyle is the outer border style for the table container.
func DocStyle(t style.Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(t.Border)
}

// RenderTable renders a table view wrapped in DocStyle.
func RenderTable(t style.Theme, tableView string) string {
	return DocStyle(t).Render(tableView)
}

// RenderPanel renders a table panel with loading state and optional header.
// Header is only shown when the table has rows.
func RenderPanel(t style.Theme, tbl *Model, loading bool, spinnerView, header string) string {
	if loading {
		return lipgloss.NewStyle().
			Width(tbl.W).
			Height(tbl.H).
			Align(lipgloss.Center, lipgloss.Center).
			Render(spinnerView + " Loading...")
	}
	tableView := DocStyle(t).Render(tbl.View())
	if len(tbl.Rows()) == 0 {
		return tableView
	}
	return lipgloss.JoinVertical(0, TitleStyle(t).Render(header), tableView)
}
