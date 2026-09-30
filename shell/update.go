package shell

import (
	"fmt"
	"image/color"
	"time"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/felipeospina21/tuishell"
	"github.com/felipeospina21/tuishell/clipboard"
	"github.com/felipeospina21/tuishell/statusline"
	"github.com/felipeospina21/tuishell/style"
)

// Update handles all shell-level behavior and routes messages to panels.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Ctx.Window = msg
		m.recomputeLayout()
		cmds = append(cmds, m.pushSizeToPanels()...)
		return m, tea.Batch(cmds...)

	case tea.KeyPressMsg:
		updated, cmd, handled := m.handleGlobalKeys(msg)
		if handled {
			return updated, cmd
		}
		m = updated
		if m.isModalOpen {
			m.Modal, cmd = m.Modal.Update(msg)
			return m, cmd
		}
		cmd = m.routeToPanel(msg)
		cmds = append(cmds, cmd)

	case tuishell.StartTaskMsg:
		m.taskStatus = taskStarted
		m.Statusline.Status = statusline.ModesEnum.Loading
		cmds = append(cmds, msg.Cmd)

	case tuishell.FinishTaskMsg:
		m.taskStatus = taskFinished
		m.Statusline.SpinnerView = ""
		if msg.Err != nil {
			m.Statusline.Status = statusline.ModesEnum.Error
			m.Statusline.Content = msg.Err.Error()
			m.taskErr = msg.Err
		} else {
			mode := statusline.ModesEnum.Normal
			if m.Ctx.DemoMode {
				mode = statusline.ModesEnum.Demo
			}
			m.Statusline.Status = mode
			m.Statusline.Content = ""
			m.taskErr = nil
			if msg.Keybinds != nil {
				m.Statusline.Keybinds = msg.Keybinds
			}
		}

	case tuishell.OpenModalMsg:
		m.isModalOpen = true
		m.prevFocus = m.Ctx.FocusedPanel
		m.Modal.Header = msg.Header
		m.Modal.Content = msg.Content
		m.Modal.IsError = msg.IsError
		m.Modal.SetFocus()

	case tuishell.CloseModalMsg:
		m.isModalOpen = false
		m.Modal.IsError = false
		m.Ctx.FocusedPanel = m.prevFocus

	case tuishell.CopyModalMsg:
		clipboard.CopyToClipboard(m.Modal.Content)

	case tuishell.ResetHighlightMsg:
		m.Modal.Highlight = false

	case tuishell.SubmitModalMsg:
		m.isModalOpen = false
		m.Ctx.FocusedPanel = m.prevFocus
		return m, func() tea.Msg { return tuishell.ShellSubmitMsg{} }

	case tuishell.OpenRightPanelMsg:
		if !m.isRightOpen {
			m.isRightOpen = true
			m.Ctx.FocusedPanel = tuishell.RightPanel
			m.recomputeLayout()
			cmds = append(cmds, m.pushSizeToPanels()...)
		}

	case tuishell.CloseLeftPanelMsg:
		if m.fixedPanels {
			break
		}
		if m.isLeftOpen {
			m.isLeftOpen = false
			m.Ctx.FocusedPanel = tuishell.MainPanel
			m.updateProjectLabel()
			m.recomputeLayout()
			cmds = append(cmds, m.pushSizeToPanels()...)
		}

	case tuishell.CloseRightPanelMsg:
		if m.fixedPanels {
			break
		}
		m.isRightOpen = false
		m.isRightFullscreen = false
		m.Ctx.FocusedPanel = tuishell.MainPanel
		m.recomputeLayout()
		cmds = append(cmds, m.pushSizeToPanels()...)

	case tuishell.ToggleFullscreenMsg:
		if m.fixedPanels {
			break
		}
		m.isRightFullscreen = !m.isRightFullscreen
		m.recomputeLayout()
		cmds = append(cmds, m.pushSizeToPanels()...)

	case tuishell.SetKeybindsMsg:
		m.Statusline.Keybinds = msg.Keybinds

	case tuishell.SetStatusMsg:
		m.Statusline.Status = msg.Mode
		m.Statusline.Content = msg.Content

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		m.Statusline.SpinnerView = m.Spinner.View()
		cmds = append(cmds, cmd)

	case tea.MouseMsg:
		cmd := m.handleMouse(msg)
		cmds = append(cmds, cmd)

	default:
		cmd := m.routeToPanel(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// handleMouse routes a mouse message to the panel under the pointer. When
// mouse support is disabled it falls back to the focused-panel routing used
// for other messages. A left click inside a visible panel focuses that panel
// (click-to-focus) before the message is forwarded. Coordinates are left
// absolute so consumers using bubblezone (which scans at the root) can
// hit-test with the original coordinates; use Layout.PanelAt for manual
// panel-relative translation.
func (m *Model) handleMouse(msg tea.MouseMsg) tea.Cmd {
	if !m.enableMouse || m.isModalOpen {
		return m.routeToPanel(msg)
	}

	mouse := msg.Mouse()

	// Toggle-button clicks take priority over panel routing. Only left clicks
	// on the bar act; other mouse events over the bar are ignored.
	if m.showButtonBar() {
		if click, isClick := msg.(tea.MouseClickMsg); isClick && click.Button == tea.MouseLeft {
			if panel, ok := m.Layout.ButtonAt(mouse.X, mouse.Y); ok {
				return m.toggleButton(panel)
			}
		}
	}

	panel, _, _, ok := m.Layout.PanelAt(mouse.X, mouse.Y)
	if !ok {
		return nil
	}

	// Click-to-focus: a left click on a visible panel focuses it.
	if click, isClick := msg.(tea.MouseClickMsg); isClick && click.Button == tea.MouseLeft {
		m.Ctx.FocusedPanel = panel
	}

	return m.routeToPanelTarget(panel, msg)
}

// toggleButton opens or closes the panel behind a clicked toggle button,
// running the same layout recompute + size push that the keybinds do. State
// changes are applied directly (not emitted as messages) so a single click is
// immediately reflected.
func (m *Model) toggleButton(panel tuishell.FocusedPanel) tea.Cmd {
	switch panel {
	case tuishell.LeftPanel:
		m.isLeftOpen = !m.isLeftOpen
		if m.isLeftOpen {
			m.Ctx.FocusedPanel = tuishell.LeftPanel
		} else {
			m.Ctx.FocusedPanel = tuishell.MainPanel
			m.updateProjectLabel()
		}
		m.recomputeLayout()
		return tea.Batch(m.pushSizeToPanels()...)
	case tuishell.RightPanel:
		m.isRightOpen = !m.isRightOpen
		if m.isRightOpen {
			m.Ctx.FocusedPanel = tuishell.RightPanel
		} else {
			m.isRightFullscreen = false
			m.Ctx.FocusedPanel = tuishell.MainPanel
		}
		m.recomputeLayout()
		return tea.Batch(m.pushSizeToPanels()...)
	}
	return nil
}

// View composes the full shell view.
func (m Model) View() string {
	var left, main, right string
	l := m.Layout

	// Active-panel indicator: recolor the focused panel's border. The left
	// panel carries its own (consumer-supplied) border; the main panel's
	// indicator is the outer frame border. The right panel renders its own
	// content (and border) as before, so its focus indicator is left to the
	// consumer's own view via Ctx.FocusedPanel. Recoloring never changes
	// width, so layout math is unaffected. Works with or without mouse support.
	leftStyle := m.focusBorder(m.leftPanelStyle, tuishell.LeftPanel)
	frameStyle := m.focusBorder(m.mainFrameStyle, tuishell.MainPanel)

	if m.isLeftOpen && !m.isRightFullscreen && m.Left != nil {
		left = leftStyle.Height(l.LeftPanel.Height).Render(m.Left.View().Content)
	}
	if !m.isRightFullscreen && m.Main != nil {
		main = lipgloss.NewStyle().
			Width(l.MainPanel.Width).
			Height(l.MainPanel.Height).
			Render(m.Main.View().Content)
	}
	if m.isRightOpen && m.Right != nil {
		right = m.Right.View().Content
	}

	body := lipgloss.JoinHorizontal(lipgloss.Top, left, main, right)

	// Prepend the mouse-mode toggle button bar when enabled.
	if m.showButtonBar() {
		body = lipgloss.JoinVertical(lipgloss.Left, m.buttonBarView(), body)
	}

	sl := m.Statusline.View()
	screen := frameStyle.Render(lipgloss.JoinVertical(lipgloss.Left, body, sl))

	if m.isModalOpen {
		screen = m.Modal.View(screen)
	}
	return screen
}

// focusBorder returns style with its border foreground recolored to the
// configured focus color when the given panel is the focused one; otherwise it
// returns style unchanged. It only recolors when style actually has a border,
// so panels without a border are unaffected.
func (m Model) focusBorder(style lipgloss.Style, panel tuishell.FocusedPanel) lipgloss.Style {
	if m.Ctx.FocusedPanel != panel {
		return style
	}
	if style.GetBorderStyle() == (lipgloss.Border{}) {
		return style
	}
	return style.BorderForeground(m.focusedBorderColor)
}

// leftButtonGlyph / rightButtonGlyph build the toggle-button strings from the
// consumer-configured labels. The arrow indicates the action: « closes an open
// panel, » opens a closed one. Both the rendered bar and the layout hit-rects
// derive their widths from these, so clicks always land on the glyph.
func (m Model) leftButtonGlyph() string {
	arrow := "»"
	if m.isLeftOpen {
		arrow = "«"
	}
	return fmt.Sprintf("[ %s %s ]", arrow, m.leftButtonLabel)
}

func (m Model) rightButtonGlyph() string {
	arrow := "»"
	if m.isRightOpen {
		arrow = "«"
	}
	return fmt.Sprintf("[ %s %s ]", m.rightButtonLabel, arrow)
}

// buttonBarView renders the one-row mouse-mode toggle bar. The left toggle
// opens/closes the left panel; the right toggle opens/closes the right panel.
// The glyph widths are kept in sync with the layout hit-rects via
// recomputeLayout, which measures the same glyphs.
func (m Model) buttonBarView() string {
	btn := lipgloss.NewStyle().Foreground(m.theme.Primary)
	leftBtn := btn.Render(m.leftButtonGlyph())
	rightBtn := btn.Render(m.rightButtonGlyph())

	width := m.Layout.ButtonBarRect.Width
	gap := width - lipgloss.Width(leftBtn) - lipgloss.Width(rightBtn)
	if gap < 1 {
		gap = 1
	}
	return leftBtn + lipgloss.NewStyle().Width(gap).Render("") + rightBtn
}

// RenderView returns a tea.View with AltScreen enabled. When the shell was
// configured with EnableMouse, it also enables tea.MouseModeCellMotion so
// mouse events (clicks, wheel, motion) are delivered to the program.
func (m Model) RenderView() tea.View {
	v := tea.NewView(m.View())
	v.AltScreen = true
	if m.enableMouse {
		v.MouseMode = tea.MouseModeCellMotion
	}
	return v
}

// --- Accessors ---

// FocusedPanel reports which panel currently has focus. Consumers can use this
// to render their own active-panel indicator (e.g. recoloring a right-panel
// border, which the shell renders as raw content).
func (m Model) FocusedPanel() tuishell.FocusedPanel { return m.Ctx.FocusedPanel }

// FocusBorderColor returns the color used for the active-panel border
// indicator, so consumer-rendered panels can match the shell's indicator.
func (m Model) FocusBorderColor() color.Color { return m.focusedBorderColor }

// IsFixedPanels reports whether the shell is in fixed three-column mode.
func (m Model) IsFixedPanels() bool { return m.fixedPanels }

// IsLeftOpen reports whether the left panel is visible.
func (m Model) IsLeftOpen() bool { return m.isLeftOpen }

// IsRightOpen reports whether the right panel is visible.
func (m Model) IsRightOpen() bool { return m.isRightOpen }

// IsRightFullscreen reports whether the right panel is in fullscreen mode.
func (m Model) IsRightFullscreen() bool { return m.isRightFullscreen }

// IsModalOpen reports whether the modal overlay is visible.
func (m Model) IsModalOpen() bool { return m.isModalOpen }

// TaskErr returns the error from the last completed task, or nil.
func (m Model) TaskErr() error { return m.taskErr }

// Theme returns the theme used by this shell.
func (m Model) Theme() style.Theme { return m.theme }

// --- Internal ---

// handleGlobalKeys returns (updated model, cmd, handled).
// When handled is true, the caller should return immediately.
func (m Model) handleGlobalKeys(msg tea.KeyPressMsg) (Model, tea.Cmd, bool) {
	match := tuishell.KeyMatcher(msg)
	gk := tuishell.GlobalKeys(m.Ctx.DemoMode)

	switch {
	case match(gk.Quit):
		return m, tea.Quit, true

	case match(gk.CyclePanel):
		m.Ctx.FocusedPanel = m.nextVisiblePanel(m.Ctx.FocusedPanel, true)
		return m, nil, true

	case match(gk.CyclePanelBack):
		m.Ctx.FocusedPanel = m.nextVisiblePanel(m.Ctx.FocusedPanel, false)
		return m, nil, true

	case !m.fixedPanels && m.isRightOpen && match(gk.CloseRightPanel):
		m.isRightOpen = false
		m.isRightFullscreen = false
		m.Ctx.FocusedPanel = tuishell.MainPanel
		m.recomputeLayout()
		cmds := m.pushSizeToPanels()
		return m, tea.Batch(cmds...), true

	case !m.fixedPanels && match(gk.ToggleLeftPanel):
		m.isLeftOpen = !m.isLeftOpen
		if m.isRightOpen {
			m.isRightOpen = false
			m.isRightFullscreen = false
		}
		if m.isLeftOpen {
			m.Ctx.FocusedPanel = tuishell.LeftPanel
		} else {
			m.Ctx.FocusedPanel = tuishell.MainPanel
		}
		m.recomputeLayout()
		cmds := m.pushSizeToPanels()
		return m, tea.Batch(cmds...), true

	case match(gk.Help):
		content := m.Modal.RenderHelp(m.Statusline.Keybinds)
		return m, func() tea.Msg {
			return tuishell.OpenModalMsg{Header: "Keybindings", Content: content}
		}, true

	case match(gk.OpenModal):
		if m.taskErr != nil {
			err := m.taskErr
			return m, func() tea.Msg {
				return tuishell.OpenModalMsg{Header: "Error", Content: err.Error(), IsError: true}
			}, true
		}

	case m.Ctx.DemoMode && match(gk.ThrowError):
		return m, func() tea.Msg {
			return tuishell.FinishTaskMsg{Err: fmt.Errorf("simulated error for testing")}
		}, true

	case m.Ctx.DemoMode && match(gk.MockFetch):
		return m, func() tea.Msg {
			return tuishell.StartTaskMsg{Cmd: func() tea.Msg {
				time.Sleep(2 * time.Second)
				return tuishell.FinishTaskMsg{}
			}}
		}, true
	}

	return m, nil, false
}

func (m *Model) recomputeLayout() {
	cfg := tuishell.LayoutConfig{
		MainFrameStyle:  m.mainFrameStyle,
		StatusBarStyle:  statusline.StatusBarStyle(),
		LeftPanelStyle:  m.leftPanelStyle,
		RightPanelStyle: m.rightPanelStyle,
		LeftPanelWidth:  m.leftPanelWidth,
		StatuslineLines: 1,
		FixedPanels:     m.fixedPanels,
		ButtonBar:       m.showButtonBar(),
	}
	if cfg.ButtonBar {
		// Measure the actual glyphs so the hit-rects track the configured
		// labels. The arrow is a single rune either way, so width is stable
		// across open/close state.
		cfg.LeftButtonWidth = lipgloss.Width(m.leftButtonGlyph())
		cfg.RightButtonWidth = lipgloss.Width(m.rightButtonGlyph())
	}
	m.Layout = tuishell.ComputeLayout(m.Ctx.Window, cfg, m.isLeftOpen, m.isRightOpen, m.isRightFullscreen)
	m.Statusline.Width = m.Layout.Statusline.Width
	m.Ctx.PanelHeight = m.Layout.ContentH
}

// nextVisiblePanel returns the next (or previous, when forward is false)
// focusable panel that is currently visible, cycling in left→main→right order.
// The main panel is always visible; left/right are included only when open.
func (m Model) nextVisiblePanel(current tuishell.FocusedPanel, forward bool) tuishell.FocusedPanel {
	var order []tuishell.FocusedPanel
	if m.isLeftOpen {
		order = append(order, tuishell.LeftPanel)
	}
	if !m.isRightFullscreen {
		order = append(order, tuishell.MainPanel)
	}
	if m.isRightOpen {
		order = append(order, tuishell.RightPanel)
	}
	if len(order) == 0 {
		return current
	}

	idx := 0
	for i, p := range order {
		if p == current {
			idx = i
			break
		}
	}
	if forward {
		idx = (idx + 1) % len(order)
	} else {
		idx = (idx - 1 + len(order)) % len(order)
	}
	return order[idx]
}

// showButtonBar reports whether the mouse-mode toggle button bar should be
// rendered: only when mouse support is on and the panels are not fixed
// (fixed panels cannot be toggled, so there is nothing to click).
func (m Model) showButtonBar() bool {
	return m.enableMouse && !m.fixedPanels
}

func (m *Model) routeToPanel(msg tea.Msg) tea.Cmd {
	return m.routeToPanelTarget(m.Ctx.FocusedPanel, msg)
}

// routeToPanelTarget forwards msg to the given panel, updating that panel's model.
func (m *Model) routeToPanelTarget(target tuishell.FocusedPanel, msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch target {
	case tuishell.LeftPanel:
		if m.Left != nil {
			m.Left, cmd = m.Left.Update(msg)
		}
	case tuishell.MainPanel:
		if m.Main != nil {
			m.Main, cmd = m.Main.Update(msg)
		}
	case tuishell.RightPanel:
		if m.Right != nil {
			m.Right, cmd = m.Right.Update(msg)
		}
	}
	return cmd
}

func (m *Model) updateProjectLabel() {
	if sp, ok := m.Left.(tuishell.SelectionProvider); ok {
		if label := sp.SelectedLabel(); label != "" {
			m.Statusline.ProjectLabel = m.appIcon + " " + label
		}
	}
}

func (m *Model) pushSizeToPanels() []tea.Cmd {
	var cmds []tea.Cmd
	l := m.Layout
	if m.Left != nil {
		var cmd tea.Cmd
		m.Left, cmd = m.Left.Update(tea.WindowSizeMsg{Width: l.LeftPanel.Width, Height: l.LeftPanel.Height})
		cmds = append(cmds, cmd)
	}
	if m.Main != nil {
		var cmd tea.Cmd
		m.Main, cmd = m.Main.Update(tea.WindowSizeMsg{Width: l.MainPanel.Width, Height: l.MainPanel.Height})
		cmds = append(cmds, cmd)
	}
	if m.Right != nil {
		var cmd tea.Cmd
		m.Right, cmd = m.Right.Update(tea.WindowSizeMsg{Width: l.RightPanel.Width, Height: l.RightPanel.Height})
		cmds = append(cmds, cmd)
	}
	return cmds
}
