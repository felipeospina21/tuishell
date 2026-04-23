package style

import "charm.land/lipgloss/v2"

// Presets contains all built-in theme presets keyed by name.
var Presets = map[string]Theme{
	// ── Catppuccin Mocha (dark) ──
	"catppuccin-mocha": {
		Primary:       lipgloss.Color("#89b4fa"),  // Blue
		PrimaryBright: lipgloss.Color("#b4befe"),  // Lavender
		PrimaryFg:     lipgloss.Color("#cdd6f4"),  // Text
		PrimaryDim:    lipgloss.Color("#45475a"),  // Surface1

		Info:          lipgloss.Color("#74c7ec"),  // Sapphire
		InfoBright:    lipgloss.Color("#89dceb"),  // Sky
		Success:       lipgloss.Color("#a6e3a1"),  // Green
		SuccessBright: lipgloss.Color("#94e2d5"),  // Teal
		Danger:        lipgloss.Color("#f38ba8"),  // Red
		DangerBright:  lipgloss.Color("#eba0ac"),  // Maroon
		Warning:       lipgloss.Color("#f9e2af"),  // Yellow
		WarningBright: lipgloss.Color("#fab387"),  // Peach
		Caution:       lipgloss.Color("#fab387"),  // Peach

		Text:            lipgloss.Color("#cdd6f4"),  // Text
		TextInverse:     lipgloss.Color("#11111b"),  // Crust
		TextDimmed:      lipgloss.Color("#a6adc8"),  // Subtext0
		Muted:           lipgloss.Color("#9399b2"),  // Overlay2
		Dim:             lipgloss.Color("#6c7086"),  // Overlay0
		Border:          lipgloss.Color("#585b70"),  // Surface2
		ModalBorder:     lipgloss.Color("#7f849c"),  // Overlay1
		SurfaceDim:      lipgloss.Color("#1e1e2e"),  // Base
		SelectionBorder: lipgloss.Color("#cba6f7"),  // Mauve

		StatusText:    lipgloss.Color("#cdd6f4"),  // Text
		StatusNormal:  lipgloss.Color("#45475a"),  // Surface1
		StatusLoading: lipgloss.Color("#94e2d5"),  // Teal
		StatusError:   lipgloss.Color("#f38ba8"),  // Red
		StatusDev:     lipgloss.Color("#a6e3a1"),  // Green
		StatusAccent1: lipgloss.Color("#cba6f7"),  // Mauve
		StatusAccent2: lipgloss.Color("#b4befe"),  // Lavender
	},

	// ── Catppuccin Macchiato ──
	"catppuccin-macchiato": {
		Primary:       lipgloss.Color("#8aadf4"),  // Blue
		PrimaryBright: lipgloss.Color("#b7bdf8"),  // Lavender
		PrimaryFg:     lipgloss.Color("#cad3f5"),  // Text
		PrimaryDim:    lipgloss.Color("#494d64"),  // Surface1

		Info:          lipgloss.Color("#7dc4e4"),  // Sapphire
		InfoBright:    lipgloss.Color("#91d7e3"),  // Sky
		Success:       lipgloss.Color("#a6da95"),  // Green
		SuccessBright: lipgloss.Color("#8bd5ca"),  // Teal
		Danger:        lipgloss.Color("#ed8796"),  // Red
		DangerBright:  lipgloss.Color("#ee99a0"),  // Maroon
		Warning:       lipgloss.Color("#eed49f"),  // Yellow
		WarningBright: lipgloss.Color("#f5a97f"),  // Peach
		Caution:       lipgloss.Color("#f5a97f"),  // Peach

		Text:            lipgloss.Color("#cad3f5"),  // Text
		TextInverse:     lipgloss.Color("#181926"),  // Crust
		TextDimmed:      lipgloss.Color("#a5adcb"),  // Subtext0
		Muted:           lipgloss.Color("#939ab7"),  // Overlay2
		Dim:             lipgloss.Color("#6e738d"),  // Overlay0
		Border:          lipgloss.Color("#5b6078"),  // Surface2
		ModalBorder:     lipgloss.Color("#8087a2"),  // Overlay1
		SurfaceDim:      lipgloss.Color("#24273a"),  // Base
		SelectionBorder: lipgloss.Color("#c6a0f6"),  // Mauve

		StatusText:    lipgloss.Color("#cad3f5"),  // Text
		StatusNormal:  lipgloss.Color("#494d64"),  // Surface1
		StatusLoading: lipgloss.Color("#8bd5ca"),  // Teal
		StatusError:   lipgloss.Color("#ed8796"),  // Red
		StatusDev:     lipgloss.Color("#a6da95"),  // Green
		StatusAccent1: lipgloss.Color("#c6a0f6"),  // Mauve
		StatusAccent2: lipgloss.Color("#b7bdf8"),  // Lavender
	},

	// ── Catppuccin Frappé ──
	"catppuccin-frappe": {
		Primary:       lipgloss.Color("#8caaee"),  // Blue
		PrimaryBright: lipgloss.Color("#babbf1"),  // Lavender
		PrimaryFg:     lipgloss.Color("#c6d0f5"),  // Text
		PrimaryDim:    lipgloss.Color("#51576d"),  // Surface1

		Info:          lipgloss.Color("#85c1dc"),  // Sapphire
		InfoBright:    lipgloss.Color("#99d1db"),  // Sky
		Success:       lipgloss.Color("#a6d189"),  // Green
		SuccessBright: lipgloss.Color("#81c8be"),  // Teal
		Danger:        lipgloss.Color("#e78284"),  // Red
		DangerBright:  lipgloss.Color("#ea999c"),  // Maroon
		Warning:       lipgloss.Color("#e5c890"),  // Yellow
		WarningBright: lipgloss.Color("#ef9f76"),  // Peach
		Caution:       lipgloss.Color("#ef9f76"),  // Peach

		Text:            lipgloss.Color("#c6d0f5"),  // Text
		TextInverse:     lipgloss.Color("#232634"),  // Crust
		TextDimmed:      lipgloss.Color("#a5adce"),  // Subtext0
		Muted:           lipgloss.Color("#949cbb"),  // Overlay2
		Dim:             lipgloss.Color("#737994"),  // Overlay0
		Border:          lipgloss.Color("#626880"),  // Surface2
		ModalBorder:     lipgloss.Color("#838ba7"),  // Overlay1
		SurfaceDim:      lipgloss.Color("#303446"),  // Base
		SelectionBorder: lipgloss.Color("#ca9ee6"),  // Mauve

		StatusText:    lipgloss.Color("#c6d0f5"),  // Text
		StatusNormal:  lipgloss.Color("#51576d"),  // Surface1
		StatusLoading: lipgloss.Color("#81c8be"),  // Teal
		StatusError:   lipgloss.Color("#e78284"),  // Red
		StatusDev:     lipgloss.Color("#a6d189"),  // Green
		StatusAccent1: lipgloss.Color("#ca9ee6"),  // Mauve
		StatusAccent2: lipgloss.Color("#babbf1"),  // Lavender
	},

	// ── Catppuccin Latte (light) ──
	"catppuccin-latte": {
		Primary:       lipgloss.Color("#1e66f5"),  // Blue
		PrimaryBright: lipgloss.Color("#7287fd"),  // Lavender
		PrimaryFg:     lipgloss.Color("#4c4f69"),  // Text
		PrimaryDim:    lipgloss.Color("#bcc0cc"),  // Surface1

		Info:          lipgloss.Color("#209fb5"),  // Sapphire
		InfoBright:    lipgloss.Color("#04a5e5"),  // Sky
		Success:       lipgloss.Color("#40a02b"),  // Green
		SuccessBright: lipgloss.Color("#179299"),  // Teal
		Danger:        lipgloss.Color("#d20f39"),  // Red
		DangerBright:  lipgloss.Color("#e64553"),  // Maroon
		Warning:       lipgloss.Color("#df8e1d"),  // Yellow
		WarningBright: lipgloss.Color("#fe640b"),  // Peach
		Caution:       lipgloss.Color("#fe640b"),  // Peach

		Text:            lipgloss.Color("#4c4f69"),  // Text
		TextInverse:     lipgloss.Color("#eff1f5"),  // Base
		TextDimmed:      lipgloss.Color("#6c6f85"),  // Subtext0
		Muted:           lipgloss.Color("#7c7f93"),  // Overlay2
		Dim:             lipgloss.Color("#9ca0b0"),  // Overlay0
		Border:          lipgloss.Color("#acb0be"),  // Surface2
		ModalBorder:     lipgloss.Color("#8c8fa1"),  // Overlay1
		SurfaceDim:      lipgloss.Color("#eff1f5"),  // Base
		SelectionBorder: lipgloss.Color("#8839ef"),  // Mauve

		StatusText:    lipgloss.Color("#eff1f5"),  // Base
		StatusNormal:  lipgloss.Color("#bcc0cc"),  // Surface1
		StatusLoading: lipgloss.Color("#179299"),  // Teal
		StatusError:   lipgloss.Color("#d20f39"),  // Red
		StatusDev:     lipgloss.Color("#40a02b"),  // Green
		StatusAccent1: lipgloss.Color("#8839ef"),  // Mauve
		StatusAccent2: lipgloss.Color("#7287fd"),  // Lavender
	},

	// ── Rosé Pine ──
	"rose-pine": {
		Primary:       lipgloss.Color("#c4a7e7"),  // Iris
		PrimaryBright: lipgloss.Color("#ebbcba"),  // Rose
		PrimaryFg:     lipgloss.Color("#e0def4"),  // Text
		PrimaryDim:    lipgloss.Color("#26233a"),  // Highlight Med

		Info:          lipgloss.Color("#9ccfd8"),  // Foam
		InfoBright:    lipgloss.Color("#c4a7e7"),  // Iris
		Success:       lipgloss.Color("#9ccfd8"),  // Foam
		SuccessBright: lipgloss.Color("#3e8fb0"),  // Pine
		Danger:        lipgloss.Color("#eb6f92"),  // Love
		DangerBright:  lipgloss.Color("#eb6f92"),  // Love
		Warning:       lipgloss.Color("#f6c177"),  // Gold
		WarningBright: lipgloss.Color("#f6c177"),  // Gold
		Caution:       lipgloss.Color("#ebbcba"),  // Rose

		Text:            lipgloss.Color("#e0def4"),  // Text
		TextInverse:     lipgloss.Color("#191724"),  // Base
		TextDimmed:      lipgloss.Color("#908caa"),  // Subtle
		Muted:           lipgloss.Color("#6e6a86"),  // Muted
		Dim:             lipgloss.Color("#524f67"),  // Highlight High
		Border:          lipgloss.Color("#403d52"),  // Highlight Med
		ModalBorder:     lipgloss.Color("#524f67"),  // Highlight High
		SurfaceDim:      lipgloss.Color("#191724"),  // Base
		SelectionBorder: lipgloss.Color("#c4a7e7"),  // Iris

		StatusText:    lipgloss.Color("#e0def4"),  // Text
		StatusNormal:  lipgloss.Color("#21202e"),  // Surface
		StatusLoading: lipgloss.Color("#9ccfd8"),  // Foam
		StatusError:   lipgloss.Color("#eb6f92"),  // Love
		StatusDev:     lipgloss.Color("#3e8fb0"),  // Pine
		StatusAccent1: lipgloss.Color("#c4a7e7"),  // Iris
		StatusAccent2: lipgloss.Color("#ebbcba"),  // Rose
	},

	// ── Tokyo Night ──
	"tokyo-night": {
		Primary:       lipgloss.Color("#7aa2f7"),  // Blue
		PrimaryBright: lipgloss.Color("#2ac3de"),  // Cyan
		PrimaryFg:     lipgloss.Color("#c0caf5"),  // Foreground
		PrimaryDim:    lipgloss.Color("#3b4261"),  // Comment

		Info:          lipgloss.Color("#7dcfff"),  // Light Blue
		InfoBright:    lipgloss.Color("#2ac3de"),  // Cyan
		Success:       lipgloss.Color("#9ece6a"),  // Green
		SuccessBright: lipgloss.Color("#73daca"),  // Teal
		Danger:        lipgloss.Color("#f7768e"),  // Red
		DangerBright:  lipgloss.Color("#ff9e64"),  // Orange
		Warning:       lipgloss.Color("#e0af68"),  // Yellow
		WarningBright: lipgloss.Color("#ff9e64"),  // Orange
		Caution:       lipgloss.Color("#ff9e64"),  // Orange

		Text:            lipgloss.Color("#c0caf5"),  // Foreground
		TextInverse:     lipgloss.Color("#1a1b26"),  // Background
		TextDimmed:      lipgloss.Color("#a9b1d6"),  // Fg Dark
		Muted:           lipgloss.Color("#565f89"),  // Comment
		Dim:             lipgloss.Color("#3b4261"),  // Line
		Border:          lipgloss.Color("#292e42"),  // Bg Highlight
		ModalBorder:     lipgloss.Color("#565f89"),  // Comment
		SurfaceDim:      lipgloss.Color("#1a1b26"),  // Background
		SelectionBorder: lipgloss.Color("#bb9af7"),  // Purple

		StatusText:    lipgloss.Color("#c0caf5"),  // Foreground
		StatusNormal:  lipgloss.Color("#292e42"),  // Bg Highlight
		StatusLoading: lipgloss.Color("#73daca"),  // Teal
		StatusError:   lipgloss.Color("#f7768e"),  // Red
		StatusDev:     lipgloss.Color("#9ece6a"),  // Green
		StatusAccent1: lipgloss.Color("#bb9af7"),  // Purple
		StatusAccent2: lipgloss.Color("#7aa2f7"),  // Blue
	},

	// ── Dracula ──
	"dracula": {
		Primary:       lipgloss.Color("#bd93f9"),  // Purple
		PrimaryBright: lipgloss.Color("#ff79c6"),  // Pink
		PrimaryFg:     lipgloss.Color("#f8f8f2"),  // Foreground
		PrimaryDim:    lipgloss.Color("#44475a"),  // Current Line

		Info:          lipgloss.Color("#8be9fd"),  // Cyan
		InfoBright:    lipgloss.Color("#bd93f9"),  // Purple
		Success:       lipgloss.Color("#50fa7b"),  // Green
		SuccessBright: lipgloss.Color("#50fa7b"),  // Green
		Danger:        lipgloss.Color("#ff5555"),  // Red
		DangerBright:  lipgloss.Color("#ff6e6e"),  // Red bright
		Warning:       lipgloss.Color("#f1fa8c"),  // Yellow
		WarningBright: lipgloss.Color("#ffb86c"),  // Orange
		Caution:       lipgloss.Color("#ffb86c"),  // Orange

		Text:            lipgloss.Color("#f8f8f2"),  // Foreground
		TextInverse:     lipgloss.Color("#282a36"),  // Background
		TextDimmed:      lipgloss.Color("#bfbfbf"),  // Fg dimmed
		Muted:           lipgloss.Color("#6272a4"),  // Comment
		Dim:             lipgloss.Color("#44475a"),  // Current Line
		Border:          lipgloss.Color("#44475a"),  // Current Line
		ModalBorder:     lipgloss.Color("#6272a4"),  // Comment
		SurfaceDim:      lipgloss.Color("#282a36"),  // Background
		SelectionBorder: lipgloss.Color("#ff79c6"),  // Pink

		StatusText:    lipgloss.Color("#f8f8f2"),  // Foreground
		StatusNormal:  lipgloss.Color("#44475a"),  // Current Line
		StatusLoading: lipgloss.Color("#8be9fd"),  // Cyan
		StatusError:   lipgloss.Color("#ff5555"),  // Red
		StatusDev:     lipgloss.Color("#50fa7b"),  // Green
		StatusAccent1: lipgloss.Color("#bd93f9"),  // Purple
		StatusAccent2: lipgloss.Color("#ff79c6"),  // Pink
	},
}
