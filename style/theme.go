// Package style provides the Theme struct for tuishell components.
package style

import "image/color"

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
