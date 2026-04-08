package style

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

// Theme defines the color tokens used across all tuishell components.
type Theme struct {
	Primary       color.Color
	PrimaryBright color.Color
	PrimaryFg     color.Color
	PrimaryDim    color.Color

	Info          color.Color
	InfoBright    color.Color
	Success       color.Color
	SuccessBright color.Color
	Danger        color.Color
	DangerBright  color.Color
	Warning       color.Color
	WarningBright color.Color
	Caution       color.Color

	Text            color.Color
	TextInverse     color.Color
	TextDimmed      color.Color
	Muted           color.Color
	Dim             color.Color
	Border          color.Color
	ModalBorder     color.Color
	SurfaceDim      color.Color
	SelectionBorder color.Color

	StatusText    color.Color
	StatusNormal  color.Color
	StatusLoading color.Color
	StatusError   color.Color
	StatusDev     color.Color
	StatusAccent1 color.Color
	StatusAccent2 color.Color
}

// DefaultTheme returns the default (mrglab) color scheme.
func DefaultTheme() Theme {
	return Theme{
		Primary:       lipgloss.Color(Violet[300]),
		PrimaryBright: lipgloss.Color(Violet[400]),
		PrimaryFg:     lipgloss.Color(Violet[50]),
		PrimaryDim:    lipgloss.Color(Violet[800]),

		Info:          lipgloss.Color(Blue[400]),
		InfoBright:    lipgloss.Color(Blue[500]),
		Success:       lipgloss.Color(Green[300]),
		SuccessBright: lipgloss.Color(Green[400]),
		Danger:        lipgloss.Color(Red[300]),
		DangerBright:  lipgloss.Color(Red[400]),
		Warning:       lipgloss.Color(Yellow[300]),
		WarningBright: lipgloss.Color(Yellow[400]),
		Caution:       lipgloss.Color(Orange[400]),

		Text:            lipgloss.Color("#C4C4C4"),
		TextInverse:     lipgloss.Color("#111"),
		TextDimmed:      lipgloss.Color("#777777"),
		Muted:           lipgloss.Color("#999999"),
		Dim:             lipgloss.Color("#444444"),
		Border:          lipgloss.Color("#3f4145"),
		ModalBorder:     lipgloss.Color("#666666"),
		SurfaceDim:      lipgloss.Color("#1e1e24"),
		SelectionBorder: lipgloss.Color("#AD58B4"),

		StatusText:    lipgloss.Color("#FFFDF5"),
		StatusNormal:  lipgloss.Color(Violet[600]),
		StatusLoading: lipgloss.Color("#1A7A94"),
		StatusError:   lipgloss.Color("#CE3060"),
		StatusDev:     lipgloss.Color("#4E8212"),
		StatusAccent1: lipgloss.Color("#A550DF"),
		StatusAccent2: lipgloss.Color("#6124DF"),
	}
}

// MainFrameStyle returns the outer border style for the main application frame.
func MainFrameStyle(t Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(t.Border)
}
