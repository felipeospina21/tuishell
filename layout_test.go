package tuishell

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func zeroFrameConfig(leftPanelWidth, statuslineLines int) LayoutConfig {
	return LayoutConfig{
		MainFrameStyle:  lipgloss.NewStyle(),
		StatusBarStyle:  lipgloss.NewStyle(),
		LeftPanelStyle:  lipgloss.NewStyle(),
		RightPanelStyle: lipgloss.NewStyle(),
		LeftPanelWidth:  leftPanelWidth,
		StatuslineLines: statuslineLines,
	}
}

func TestComputeLayout(t *testing.T) {
	win := tea.WindowSizeMsg{Width: 120, Height: 40}
	cfg := zeroFrameConfig(30, 1)

	tests := []struct {
		name            string
		leftOpen        bool
		rightOpen       bool
		rightFullscreen bool
		wantLeft        PanelSize
		wantMain        PanelSize
		wantRight       PanelSize
		wantStatusline  PanelSize
		wantContentH    int
	}{
		{
			name:           "all panels closed",
			wantLeft:       PanelSize{Width: 0, Height: 39},
			wantMain:       PanelSize{Width: 120, Height: 39},
			wantRight:      PanelSize{Width: 0, Height: 39},
			wantStatusline: PanelSize{Width: 120, Height: 1},
			wantContentH:   39,
		},
		{
			name:           "left open only",
			leftOpen:       true,
			wantLeft:       PanelSize{Width: 30, Height: 39},
			wantMain:       PanelSize{Width: 90, Height: 39},
			wantRight:      PanelSize{Width: 0, Height: 39},
			wantStatusline: PanelSize{Width: 120, Height: 1},
			wantContentH:   39,
		},
		{
			name:           "right open left closed",
			rightOpen:      true,
			wantLeft:       PanelSize{Width: 0, Height: 39},
			wantMain:       PanelSize{Width: 60, Height: 39},
			wantRight:      PanelSize{Width: 60, Height: 39},
			wantStatusline: PanelSize{Width: 120, Height: 1},
			wantContentH:   39,
		},
		{
			name:            "right fullscreen",
			rightOpen:       true,
			rightFullscreen: true,
			wantLeft:        PanelSize{Width: 0, Height: 39},
			wantMain:        PanelSize{Width: 0, Height: 39},
			wantRight:       PanelSize{Width: 120, Height: 39},
			wantStatusline:  PanelSize{Width: 120, Height: 1},
			wantContentH:    39,
		},
		{
			name:            "left open ignored when right fullscreen",
			leftOpen:        true,
			rightOpen:       true,
			rightFullscreen: true,
			wantLeft:        PanelSize{Width: 0, Height: 39},
			wantMain:        PanelSize{Width: 0, Height: 39},
			wantRight:       PanelSize{Width: 120, Height: 39},
			wantStatusline:  PanelSize{Width: 120, Height: 1},
			wantContentH:    39,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ComputeLayout(win, cfg, tt.leftOpen, tt.rightOpen, tt.rightFullscreen)

			if got.LeftPanel != tt.wantLeft {
				t.Errorf("LeftPanel = %+v, want %+v", got.LeftPanel, tt.wantLeft)
			}
			if got.MainPanel != tt.wantMain {
				t.Errorf("MainPanel = %+v, want %+v", got.MainPanel, tt.wantMain)
			}
			if got.RightPanel != tt.wantRight {
				t.Errorf("RightPanel = %+v, want %+v", got.RightPanel, tt.wantRight)
			}
			if got.Statusline != tt.wantStatusline {
				t.Errorf("Statusline = %+v, want %+v", got.Statusline, tt.wantStatusline)
			}
			if got.ContentH != tt.wantContentH {
				t.Errorf("ContentH = %d, want %d", got.ContentH, tt.wantContentH)
			}
			if got.Window != win {
				t.Errorf("Window = %+v, want %+v", got.Window, win)
			}
		})
	}
}

func TestComputeLayout_ZeroWindow(t *testing.T) {
	win := tea.WindowSizeMsg{Width: 0, Height: 0}
	cfg := zeroFrameConfig(30, 1)
	got := ComputeLayout(win, cfg, true, false, false)

	if got.ContentH != -1 {
		t.Errorf("ContentH = %d, want -1", got.ContentH)
	}
	if got.MainPanel.Width != -30 {
		t.Errorf("MainPanel.Width = %d, want -30", got.MainPanel.Width)
	}
}

func TestPanelRectContains(t *testing.T) {
	r := PanelRect{X: 5, Y: 2, Width: 10, Height: 4}
	tests := []struct {
		name string
		x, y int
		want bool
	}{
		{"top-left corner", 5, 2, true},
		{"inside", 10, 4, true},
		{"bottom-right edge exclusive x", 15, 4, false},
		{"bottom-right edge exclusive y", 10, 6, false},
		{"left of rect", 4, 3, false},
		{"above rect", 10, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.Contains(tt.x, tt.y); got != tt.want {
				t.Errorf("Contains(%d,%d) = %v, want %v", tt.x, tt.y, got, tt.want)
			}
		})
	}

	// Zero-size rect contains nothing.
	if (PanelRect{}).Contains(0, 0) {
		t.Error("zero PanelRect should contain nothing")
	}
}

func TestPanelAt_ZeroFrame(t *testing.T) {
	win := tea.WindowSizeMsg{Width: 120, Height: 40}
	cfg := zeroFrameConfig(30, 1)
	// left open only: left [0,30), main [30,120), height 39.
	l := ComputeLayout(win, cfg, true, false, false)

	tests := []struct {
		name      string
		x, y      int
		wantPanel FocusedPanel
		wantLX    int
		wantLY    int
		wantOK    bool
	}{
		{"left panel origin", 0, 0, LeftPanel, 0, 0, true},
		{"left panel inner", 10, 5, LeftPanel, 10, 5, true},
		{"main panel origin", 30, 0, MainPanel, 0, 0, true},
		{"main panel inner", 50, 10, MainPanel, 20, 10, true},
		{"statusline row (below content)", 40, 39, 0, 0, 0, false},
		{"beyond width", 200, 5, 0, 0, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			panel, lx, ly, ok := l.PanelAt(tt.x, tt.y)
			if ok != tt.wantOK {
				t.Fatalf("PanelAt(%d,%d) ok = %v, want %v", tt.x, tt.y, ok, tt.wantOK)
			}
			if !ok {
				return
			}
			if panel != tt.wantPanel || lx != tt.wantLX || ly != tt.wantLY {
				t.Errorf("PanelAt(%d,%d) = (panel %d, lx %d, ly %d), want (panel %d, lx %d, ly %d)",
					tt.x, tt.y, panel, lx, ly, tt.wantPanel, tt.wantLX, tt.wantLY)
			}
		})
	}
}

func TestPanelAt_BorderedFrame(t *testing.T) {
	win := tea.WindowSizeMsg{Width: 120, Height: 40}
	cfg := LayoutConfig{
		MainFrameStyle:  lipgloss.NewStyle().Border(lipgloss.NormalBorder()),
		StatusBarStyle:  lipgloss.NewStyle(),
		LeftPanelStyle:  lipgloss.NewStyle(),
		RightPanelStyle: lipgloss.NewStyle(),
		LeftPanelWidth:  30,
		StatuslineLines: 1,
	}
	l := ComputeLayout(win, cfg, true, false, false)

	// A 1-cell border pushes content origin to (1,1).
	if l.LeftRect.X != 1 || l.LeftRect.Y != 1 {
		t.Errorf("LeftRect origin = (%d,%d), want (1,1)", l.LeftRect.X, l.LeftRect.Y)
	}
	// Main starts after the left panel width.
	if l.MainRect.X != 1+l.LeftPanel.Width {
		t.Errorf("MainRect.X = %d, want %d", l.MainRect.X, 1+l.LeftPanel.Width)
	}

	// A click at the border origin (0,0) hits no panel.
	if _, _, _, ok := l.PanelAt(0, 0); ok {
		t.Error("PanelAt(0,0) on border should not hit a panel")
	}
	// A click just inside the border hits the left panel at local (0,0).
	if p, lx, ly, ok := l.PanelAt(1, 1); !ok || p != LeftPanel || lx != 0 || ly != 0 {
		t.Errorf("PanelAt(1,1) = (%d,%d,%d,%v), want (LeftPanel,0,0,true)", p, lx, ly, ok)
	}
}

func TestComputeLayout_FixedPanels(t *testing.T) {
	win := tea.WindowSizeMsg{Width: 120, Height: 40}
	cfg := zeroFrameConfig(30, 1)
	cfg.FixedPanels = true

	// Even though we pass leftOpen=false, rightOpen=false, fixed mode forces
	// all three columns open with real widths.
	l := ComputeLayout(win, cfg, false, false, false)

	if l.LeftPanel.Width != 30 {
		t.Errorf("fixed LeftPanel.Width = %d, want 30", l.LeftPanel.Width)
	}
	if l.RightPanel.Width <= 0 {
		t.Errorf("fixed RightPanel.Width = %d, want > 0", l.RightPanel.Width)
	}
	if l.MainPanel.Width <= 0 {
		t.Errorf("fixed MainPanel.Width = %d, want > 0", l.MainPanel.Width)
	}
	// Columns must tile without overlap: left+main+right == inner width.
	if sum := l.LeftPanel.Width + l.MainPanel.Width + l.RightPanel.Width; sum != win.Width {
		t.Errorf("column widths sum = %d, want %d", sum, win.Width)
	}
	// All three rects are hit-testable.
	if _, _, _, ok := l.PanelAt(l.RightRect.X, l.RightRect.Y); !ok {
		t.Error("right panel rect not hit-testable in fixed mode")
	}
}

func TestComputeLayout_ButtonBarReservesRow(t *testing.T) {
	win := tea.WindowSizeMsg{Width: 120, Height: 40}
	cfg := zeroFrameConfig(30, 1)

	base := ComputeLayout(win, cfg, true, false, false)
	cfg.ButtonBar = true
	withBar := ComputeLayout(win, cfg, true, false, false)

	// The bar consumes exactly one content row from the panels.
	if withBar.MainPanel.Height != base.MainPanel.Height-1 {
		t.Errorf("with button bar MainPanel.Height = %d, want %d",
			withBar.MainPanel.Height, base.MainPanel.Height-1)
	}
	// Panels are pushed down by one row.
	if withBar.MainRect.Y != base.MainRect.Y+1 {
		t.Errorf("with button bar MainRect.Y = %d, want %d",
			withBar.MainRect.Y, base.MainRect.Y+1)
	}
	// Bar rect occupies the top content row.
	if withBar.ButtonBarRect.Height != 1 || withBar.ButtonBarRect.Y != base.MainRect.Y {
		t.Errorf("ButtonBarRect = %+v, want 1-high row at y=%d", withBar.ButtonBarRect, base.MainRect.Y)
	}
}

func TestButtonAt(t *testing.T) {
	win := tea.WindowSizeMsg{Width: 120, Height: 40}
	cfg := zeroFrameConfig(30, 1)
	cfg.ButtonBar = true
	l := ComputeLayout(win, cfg, true, false, false)

	// Left toggle at x=0 on the bar row.
	if p, ok := l.ButtonAt(l.LeftToggleRect.X, l.LeftToggleRect.Y); !ok || p != LeftPanel {
		t.Errorf("ButtonAt(left toggle) = (%d,%v), want (LeftPanel,true)", p, ok)
	}
	// Right toggle near the far right of the bar row.
	if p, ok := l.ButtonAt(l.RightToggleRect.X+1, l.RightToggleRect.Y); !ok || p != RightPanel {
		t.Errorf("ButtonAt(right toggle) = (%d,%v), want (RightPanel,true)", p, ok)
	}
	// Middle of the bar hits no button.
	if _, ok := l.ButtonAt(60, l.ButtonBarRect.Y); ok {
		t.Error("ButtonAt(mid bar) should hit no button")
	}
}
