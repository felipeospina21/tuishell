package tuishell

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/felipeospina21/tuishell/style"
)

// PanelSize holds the computed width and height for a UI region.
type PanelSize struct {
	Width  int
	Height int
}

// Layout holds all computed dimensions for the current window size and panel state.
type Layout struct {
	Window     tea.WindowSizeMsg
	LeftPanel  PanelSize
	MainPanel  PanelSize
	RightPanel PanelSize
	Statusline PanelSize
	ContentH   int
}

// LayoutConfig provides the frame sizes needed to compute the layout.
// These come from the styles of the consuming app's panels.
type LayoutConfig struct {
	MainFrameStyle    lipgloss.Style
	StatusBarStyle    lipgloss.Style
	LeftPanelStyle    lipgloss.Style
	RightPanelStyle   lipgloss.Style
	LeftPanelWidth    int
	StatuslineLines   int
}

// DefaultLayoutConfig returns a config using the default theme.
func DefaultLayoutConfig(t style.Theme) LayoutConfig {
	return LayoutConfig{
		MainFrameStyle:  style.MainFrameStyle(t),
		StatusBarStyle:  lipgloss.NewStyle().Margin(0, 0),
		LeftPanelStyle:  lipgloss.NewStyle(),
		RightPanelStyle: lipgloss.NewStyle(),
		LeftPanelWidth:  30,
		StatuslineLines: 1,
	}
}

// ComputeLayout calculates panel dimensions from the window size and panel state.
func ComputeLayout(win tea.WindowSizeMsg, cfg LayoutConfig, leftOpen, rightOpen, rightFullscreen bool) Layout {
	mainFrameX, mainFrameY := cfg.MainFrameStyle.GetFrameSize()

	innerW := win.Width - mainFrameX
	innerH := win.Height - mainFrameY

	slFrameY := cfg.StatusBarStyle.GetVerticalFrameSize()
	slH := cfg.StatuslineLines + slFrameY
	slFrameX := cfg.StatusBarStyle.GetHorizontalFrameSize()

	contentH := innerH - slH

	leftW := 0
	if leftOpen && !rightFullscreen {
		leftW = cfg.LeftPanelWidth + cfg.LeftPanelStyle.GetHorizontalFrameSize()
	}

	mainW := innerW - leftW
	rightW := 0
	detailsFrameX := cfg.RightPanelStyle.GetHorizontalFrameSize()
	if rightOpen && rightFullscreen {
		mainW = 0
		rightW = innerW - detailsFrameX
	} else if rightOpen && !leftOpen {
		rightW = mainW/2 - detailsFrameX
		mainW = mainW - rightW - detailsFrameX
	}

	leftH := contentH - cfg.LeftPanelStyle.GetVerticalFrameSize()

	return Layout{
		Window:     win,
		LeftPanel:  PanelSize{Width: leftW, Height: leftH},
		MainPanel:  PanelSize{Width: mainW, Height: contentH},
		RightPanel: PanelSize{Width: rightW, Height: contentH},
		Statusline: PanelSize{Width: innerW - slFrameX, Height: slH},
		ContentH:   contentH,
	}
}
