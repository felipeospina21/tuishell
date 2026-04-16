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
