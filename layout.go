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
// These come from the styles of the consuming app's panels.
type LayoutConfig struct {
	MainFrameStyle  lipgloss.Style
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
		MainFrameStyle: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(t.Border),
		StatusBarStyle:  lipgloss.NewStyle().Margin(0, 0),
		LeftPanelStyle:  lipgloss.NewStyle(),
		RightPanelStyle: lipgloss.NewStyle(),
		LeftPanelWidth:  30,
		StatuslineLines: 1,
	}
}

// ComputeLayout calculates panel dimensions from the window size and panel state.
func ComputeLayout(win tea.WindowSizeMsg, cfg LayoutConfig, leftOpen, rightOpen, rightFullscreen bool) Layout {
	// Fixed-panels mode forces both side panels open and disables fullscreen.
	if cfg.FixedPanels {
		leftOpen = true
		rightOpen = true
		rightFullscreen = false
	}

	mainFrameX, mainFrameY := cfg.MainFrameStyle.GetFrameSize()

	innerW := win.Width - mainFrameX
	innerH := win.Height - mainFrameY

	slFrameY := cfg.StatusBarStyle.GetVerticalFrameSize()
	slH := cfg.StatuslineLines + slFrameY
	slFrameX := cfg.StatusBarStyle.GetHorizontalFrameSize()

	contentH := innerH - slH

	// Reserve one row at the top for the mouse-mode toggle button bar.
	barH := 0
	if cfg.ButtonBar {
		barH = 1
	}
	panelsH := contentH - barH

	leftW := 0
	if leftOpen && !rightFullscreen {
		leftW = cfg.LeftPanelWidth + cfg.LeftPanelStyle.GetHorizontalFrameSize()
	}

	mainW := innerW - leftW
	rightW := 0
	detailsFrameX := cfg.RightPanelStyle.GetHorizontalFrameSize()
	switch {
	case rightOpen && rightFullscreen:
		mainW = 0
		rightW = innerW - detailsFrameX
	case cfg.FixedPanels && rightOpen:
		// Three-column layout: split the space remaining after the left panel
		// between main and right.
		rightW = mainW/2 - detailsFrameX
		mainW = mainW - rightW - detailsFrameX
	case rightOpen && !leftOpen:
		rightW = mainW/2 - detailsFrameX
		mainW = mainW - rightW - detailsFrameX
	}

	leftH := panelsH - cfg.LeftPanelStyle.GetVerticalFrameSize()

	// Absolute origins for hit-testing. The main frame's left/top border and
	// padding push the body content inward; panels are laid out left→main→right.
	originX := cfg.MainFrameStyle.GetMarginLeft() +
		cfg.MainFrameStyle.GetBorderLeftSize() +
		cfg.MainFrameStyle.GetPaddingLeft()
	originY := cfg.MainFrameStyle.GetMarginTop() +
		cfg.MainFrameStyle.GetBorderTopSize() +
		cfg.MainFrameStyle.GetPaddingTop()

	// The button bar occupies the first content row; panels start below it.
	panelsY := originY + barH

	leftRect := PanelRect{X: originX, Y: panelsY, Width: leftW, Height: leftH}
	mainRect := PanelRect{X: originX + leftW, Y: panelsY, Width: mainW, Height: panelsH}
	rightRect := PanelRect{X: originX + leftW + mainW, Y: panelsY, Width: rightW, Height: panelsH}

	l := Layout{
		Window:     win,
		LeftPanel:  PanelSize{Width: leftW, Height: leftH},
		MainPanel:  PanelSize{Width: mainW, Height: panelsH},
		RightPanel: PanelSize{Width: rightW, Height: panelsH},
		Statusline: PanelSize{Width: innerW - slFrameX, Height: slH},
		ContentH:   panelsH,
		LeftRect:   leftRect,
		MainRect:   mainRect,
		RightRect:  rightRect,
	}

	if cfg.ButtonBar {
		l.ButtonBarRect = PanelRect{X: originX, Y: originY, Width: innerW, Height: 1}
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
		l.LeftToggleRect = PanelRect{X: originX, Y: originY, Width: leftBtnW, Height: 1}
		l.RightToggleRect = PanelRect{X: originX + innerW - rightBtnW, Y: originY, Width: rightBtnW, Height: 1}
	}

	return l
}
