// Package loader provides a reusable loading indicator component.
package loader

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/felipeospina21/tuishell/style"
)

// View renders a spinner frame with a "Loading..." label.
func View(theme style.Theme, spinnerView string) string {
	textStyle := lipgloss.NewStyle().Foreground(theme.Primary)
	return fmt.Sprintf("\n %s %s\n\n", spinnerView, textStyle.Render("Loading..."))
}
