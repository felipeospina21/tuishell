package shell

import (
	"fmt"
	"time"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/felipeospina21/tuishell"
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
			if m.Ctx.DevMode {
				mode = statusline.ModesEnum.Dev
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
		// Apps handle clipboard — forwarded via returned cmd

	case tuishell.ResetHighlightMsg:
		m.Modal.Highlight = false

	case tuishell.SubmitModalMsg:
		m.isModalOpen = false
		m.Ctx.FocusedPanel = m.prevFocus
		return m, func() tea.Msg { return tuishell.ShellSubmitMsg{} }

	case tuishell.OpenRightPanelMsg:
		if !m.isRightOpen {
			m.isRightOpen = true
			m.recomputeLayout()
			cmds = append(cmds, m.pushSizeToPanels()...)
		}

	case tuishell.CloseLeftPanelMsg:
		if m.isLeftOpen {
			m.isLeftOpen = false
			m.Ctx.FocusedPanel = tuishell.MainPanel
			m.updateProjectLabel()
			m.recomputeLayout()
			cmds = append(cmds, m.pushSizeToPanels()...)
		}

	case tuishell.CloseRightPanelMsg:
		m.isRightOpen = false
		m.isRightFullscreen = false
		m.recomputeLayout()
		cmds = append(cmds, m.pushSizeToPanels()...)

	case tuishell.ToggleFullscreenMsg:
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

	default:
		cmd := m.routeToPanel(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View composes the full shell view.
func (m Model) View() string {
	var left, main, right string
	l := m.Layout

	if m.isLeftOpen && !m.isRightFullscreen && m.Left != nil {
		left = m.leftPanelStyle.Height(l.LeftPanel.Height).Render(m.Left.View().Content)
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
	sl := m.Statusline.View()
	screen := m.mainFrameStyle.Render(lipgloss.JoinVertical(lipgloss.Left, body, sl))

	if m.isModalOpen {
		screen = m.Modal.View(screen)
	}
	return screen
}

// RenderView returns a tea.View with AltScreen enabled.
func (m Model) RenderView() tea.View {
	v := tea.NewView(m.View())
	v.AltScreen = true
	return v
}

// --- Accessors ---

func (m Model) IsLeftOpen() bool        { return m.isLeftOpen }
func (m Model) IsRightOpen() bool       { return m.isRightOpen }
func (m Model) IsRightFullscreen() bool { return m.isRightFullscreen }
func (m Model) IsModalOpen() bool       { return m.isModalOpen }
func (m Model) TaskErr() error          { return m.taskErr }
func (m Model) Theme() style.Theme      { return m.theme }

// --- Internal ---

// handleGlobalKeys returns (updated model, cmd, handled).
// When handled is true, the caller should return immediately.
func (m Model) handleGlobalKeys(msg tea.KeyPressMsg) (Model, tea.Cmd, bool) {
	match := tuishell.KeyMatcher(msg)
	gk := tuishell.GlobalKeys(m.Ctx.DevMode)

	switch {
	case match(gk.Quit):
		return m, tea.Quit, true

	case match(gk.ToggleLeftPanel):
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

	case m.Ctx.DevMode && match(gk.ThrowError):
		return m, func() tea.Msg {
			return tuishell.FinishTaskMsg{Err: fmt.Errorf("simulated error for testing")}
		}, true

	case m.Ctx.DevMode && match(gk.MockFetch):
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
	}
	m.Layout = tuishell.ComputeLayout(m.Ctx.Window, cfg, m.isLeftOpen, m.isRightOpen, m.isRightFullscreen)
	m.Statusline.Width = m.Layout.Statusline.Width
	m.Ctx.PanelHeight = m.Layout.ContentH
}

func (m *Model) routeToPanel(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch m.Ctx.FocusedPanel {
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
