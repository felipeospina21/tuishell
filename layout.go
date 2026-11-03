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

// clampNonNeg returns n, or 0 when n is negative.
func clampNonNeg(n int) int {
	if n < 0 {
		return 0
	}
	return n
}

// PanelRect holds the absolute on-screen position and size of a panel region,
// in terminal cell coordinates. X/Y are the zero-based coordinates of the
// top-left cell of the panel's content area (inside any frame border).
type PanelRect struct {
	X      int
	Y      int
	Width  int
	Height int
}

// Contains reports whether the absolute coordinate (x, y) falls inside the rect.
func (r PanelRect) Contains(x, y int) bool {
	return r.Width > 0 && r.Height > 0 &&
		x >= r.X && x < r.X+r.Width &&
		y >= r.Y && y < r.Y+r.Height
}

// Layout holds all computed dimensions for the current window size and panel state.
type Layout struct {
	Window     tea.WindowSizeMsg
	LeftPanel  PanelSize
	MainPanel  PanelSize
	RightPanel PanelSize
	Statusline PanelSize
	ContentH   int

	// Absolute on-screen rectangles for hit-testing, in terminal cells.
	LeftRect  PanelRect
	MainRect  PanelRect
	RightRect PanelRect

	// Button-bar rectangles for mouse hit-testing. These are non-empty only
	// when the shell renders the mouse-mode toggle button bar (mouse enabled
	// and panels not fixed). ButtonBarRect is the full bar row; LeftToggleRect
	// and RightToggleRect are the individual clickable toggle regions.
	ButtonBarRect   PanelRect
	LeftToggleRect  PanelRect
	RightToggleRect PanelRect
}

// ButtonAt reports which toggle button (if any) is under the absolute
// coordinate (x, y). The returned panel is LeftPanel or RightPanel to indicate
// which panel that button toggles; ok is false when no button is hit.
func (l Layout) ButtonAt(x, y int) (panel FocusedPanel, ok bool) {
	if l.LeftToggleRect.Contains(x, y) {
		return LeftPanel, true
	}
	if l.RightToggleRect.Contains(x, y) {
		return RightPanel, true
	}
	return 0, false
}

// PanelAt maps an absolute mouse coordinate to the panel under it, returning
// the panel-relative coordinates (localX, localY) and ok=true when the point
// lands inside a visible panel. When the point is outside every panel (a
// frame border, the statusline, or empty space) it returns ok=false.
//
// This is the recommended entry point for consumers doing mouse hit-testing:
// it accounts for the frame border, left-panel width/border, and right-panel
// placement so callers do not have to duplicate layout math.
func (l Layout) PanelAt(x, y int) (panel FocusedPanel, localX, localY int, ok bool) {
	if l.LeftRect.Contains(x, y) {
		return LeftPanel, x - l.LeftRect.X, y - l.LeftRect.Y, true
	}
	if l.MainRect.Contains(x, y) {
		return MainPanel, x - l.MainRect.X, y - l.MainRect.Y, true
	}
	if l.RightRect.Contains(x, y) {
		return RightPanel, x - l.RightRect.X, y - l.RightRect.Y, true
	}
	return 0, 0, 0, false
}

// LayoutConfig provides the frame sizes needed to compute the layout.
// These come from the styles of the consuming app's panels. Each panel is laid
// out as an independent box using its own style; tuishell imposes no outer
// frame. The consumer decides each panel's border (or lack of one).
type LayoutConfig struct {
	MainPanelStyle  lipgloss.Style
	StatusBarStyle  lipgloss.Style
	LeftPanelStyle  lipgloss.Style
	RightPanelStyle lipgloss.Style
	LeftPanelWidth  int
	StatuslineLines int

	// FixedPanels lays out left+main+right as three always-visible columns.
	// When set, the leftOpen/rightOpen flags passed to ComputeLayout are
	// treated as true and the right panel is allocated real width alongside
	// the left panel (which the default two-column model does not do).
	FixedPanels bool

	// ButtonBar reserves one content row at the top for the mouse-mode toggle
	// button bar and populates the button rects on the returned Layout.
	ButtonBar bool

	// LeftButtonWidth and RightButtonWidth are the rendered cell widths of the
	// two toggle buttons, used to size their hit-rects so clicks land on the
	// glyphs regardless of the (consumer-configurable) label text. Used only
	// when ButtonBar is set.
	LeftButtonWidth  int
	RightButtonWidth int
}

// DefaultLayoutConfig returns a config using the given theme.
func DefaultLayoutConfig(t style.Theme) LayoutConfig {
	return LayoutConfig{
		MainPanelStyle: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(t.Border),
		StatusBarStyle:  lipgloss.NewStyle().Margin(0, 0),
		LeftPanelStyle:  lipgloss.NewStyle(),
		RightPanelStyle: lipgloss.NewStyle(),
		LeftPanelWidth:  30,
		StatuslineLines: 1,
	}
}

// ComputeLayout calculates panel dimensions from the window size and panel
// state. Each panel is an independent box laid out left→main→right with its own
// style/border; there is no outer frame. The returned PanelSize widths/heights
// are CONTENT sizes (what each panel should render into, before its own border
// is added). The statusline spans the full width below the panel row.
func ComputeLayout(win tea.WindowSizeMsg, cfg LayoutConfig, leftOpen, rightOpen, rightFullscreen bool) Layout {
	// Fixed-panels mode forces both side panels open and disables fullscreen.
	if cfg.FixedPanels {
		leftOpen = true
		rightOpen = true
		rightFullscreen = false
	}

	// Per-panel frame sizes (border + padding + margin) from each style.
	leftFX, leftFY := cfg.LeftPanelStyle.GetFrameSize()
	mainFX, mainFY := cfg.MainPanelStyle.GetFrameSize()
	rightFX, rightFY := cfg.RightPanelStyle.GetFrameSize()

	slFrameY := cfg.StatusBarStyle.GetVerticalFrameSize()
	slH := cfg.StatuslineLines + slFrameY
	slFrameX := cfg.StatusBarStyle.GetHorizontalFrameSize()

	// Reserve one row at the top for the mouse-mode toggle button bar.
	barH := 0
	if cfg.ButtonBar {
		barH = 1
	}
	// Total height available to the panel row (on-screen, including borders).
	panelsTotalH := win.Height - slH - barH

	// Total on-screen widths (content + own frame) for each panel.
	leftTotal := 0
	if leftOpen && !rightFullscreen {
		leftTotal = cfg.LeftPanelWidth + leftFX
	}

	rightTotal := 0
	mainTotal := win.Width - leftTotal
	switch {
	case rightOpen && rightFullscreen:
		mainTotal = 0
		rightTotal = win.Width
	case rightOpen:
		// Split the space remaining after the left panel between main and
		// right (leftTotal is 0 when the left panel is closed). Covers fixed
		// mode, left-closed, and left+right both open.
		rightTotal = mainTotal / 2
		mainTotal = mainTotal - rightTotal
	}

	// Convert on-screen totals to content sizes by subtracting each panel's
	// own frame. Clamp at 0 so a panel narrower than its frame never goes
	// negative.
	leftW := clampNonNeg(leftTotal - leftFX)
	if leftTotal == 0 {
		leftW = 0
	}
	mainW := clampNonNeg(mainTotal - mainFX)
	if mainTotal == 0 {
		mainW = 0
	}
	rightW := clampNonNeg(rightTotal - rightFX)
	if rightTotal == 0 {
		rightW = 0
	}

	leftH := clampNonNeg(panelsTotalH - leftFY)
	mainH := clampNonNeg(panelsTotalH - mainFY)
	rightH := clampNonNeg(panelsTotalH - rightFY)

	// Absolute origins for hit-testing. Panels sit flush left→right with no
	// outer frame; the button bar (when shown) occupies the first row. Each
	// panel's content area is offset inward by that panel's own left border +
	// left padding.
	panelsY := barH
	leftContentY := panelsY + cfg.LeftPanelStyle.GetBorderTopSize() + cfg.LeftPanelStyle.GetPaddingTop()
	mainContentY := panelsY + cfg.MainPanelStyle.GetBorderTopSize() + cfg.MainPanelStyle.GetPaddingTop()
	rightContentY := panelsY + cfg.RightPanelStyle.GetBorderTopSize() + cfg.RightPanelStyle.GetPaddingTop()

	leftContentX := cfg.LeftPanelStyle.GetBorderLeftSize() + cfg.LeftPanelStyle.GetPaddingLeft()
	mainContentX := leftTotal + cfg.MainPanelStyle.GetBorderLeftSize() + cfg.MainPanelStyle.GetPaddingLeft()
	rightContentX := leftTotal + mainTotal + cfg.RightPanelStyle.GetBorderLeftSize() + cfg.RightPanelStyle.GetPaddingLeft()

	leftRect := PanelRect{X: leftContentX, Y: leftContentY, Width: leftW, Height: leftH}
	mainRect := PanelRect{X: mainContentX, Y: mainContentY, Width: mainW, Height: mainH}
	rightRect := PanelRect{X: rightContentX, Y: rightContentY, Width: rightW, Height: rightH}

	l := Layout{
		Window:     win,
		LeftPanel:  PanelSize{Width: leftW, Height: leftH},
		MainPanel:  PanelSize{Width: mainW, Height: mainH},
		RightPanel: PanelSize{Width: rightW, Height: rightH},
		Statusline: PanelSize{Width: win.Width - slFrameX, Height: slH},
		ContentH:   panelsTotalH,
		LeftRect:   leftRect,
		MainRect:   mainRect,
		RightRect:  rightRect,
	}

	if cfg.ButtonBar {
		l.ButtonBarRect = PanelRect{X: 0, Y: 0, Width: win.Width, Height: 1}
		// Widths come from the shell, which measures the rendered glyphs so the
		// hit-rects track the consumer-configured labels. Fall back to sane
		// defaults if unset.
		leftBtnW := cfg.LeftButtonWidth
		if leftBtnW <= 0 {
			leftBtnW = 9
		}
		rightBtnW := cfg.RightButtonWidth
		if rightBtnW <= 0 {
			rightBtnW = 13
		}
		l.LeftToggleRect = PanelRect{X: 0, Y: 0, Width: leftBtnW, Height: 1}
		l.RightToggleRect = PanelRect{X: win.Width - rightBtnW, Y: 0, Width: rightBtnW, Height: 1}
	}

	return l
}
