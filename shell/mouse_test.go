package shell

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/felipeospina21/tuishell"
	"github.com/felipeospina21/tuishell/style"
)

// fakePanel records the messages it receives so tests can assert routing.
type fakePanel struct {
	name     string
	received []tea.Msg
}

func (p *fakePanel) Init() tea.Cmd { return nil }

func (p *fakePanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Ignore the window-size messages the shell pushes during layout so tests
	// can focus on the mouse messages of interest.
	if _, ok := msg.(tea.WindowSizeMsg); !ok {
		p.received = append(p.received, msg)
	}
	return p, nil
}

func (p *fakePanel) View() tea.View { return tea.NewView("") }

func (p *fakePanel) gotMouse() bool {
	for _, m := range p.received {
		if _, ok := m.(tea.MouseMsg); ok {
			return true
		}
	}
	return false
}

func newTestShell(t *testing.T, enableMouse bool) (Model, *fakePanel, *fakePanel) {
	t.Helper()
	left := &fakePanel{name: "left"}
	main := &fakePanel{name: "main"}
	m := New(Config{
		Theme:          style.Theme{},
		LeftPanel:      left,
		MainPanel:      main,
		Keybinds:       tuishell.GlobalKeys(false),
		LeftPanelWidth: 30,
		LeftPanelStyle: lipgloss.NewStyle(),
		MainFrameStyle: lipgloss.NewStyle(), // zero frame => origin (0,0)
		EnableMouse:    enableMouse,
	})
	// Establish a known layout: left [0,30), main [30,120).
	m, _ = m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	return m, left, main
}

func TestClickToFocusMainPanel(t *testing.T) {
	m, _, main := newTestShell(t, true)

	if m.Ctx.FocusedPanel != tuishell.LeftPanel {
		t.Fatalf("initial focus = %d, want LeftPanel", m.Ctx.FocusedPanel)
	}

	// Click inside the main panel region (x >= 30).
	m, _ = m.Update(tea.MouseClickMsg(tea.Mouse{Button: tea.MouseLeft, X: 50, Y: 5}))

	if m.Ctx.FocusedPanel != tuishell.MainPanel {
		t.Errorf("after click on main, focus = %d, want MainPanel", m.Ctx.FocusedPanel)
	}
	if !main.gotMouse() {
		t.Error("main panel did not receive the mouse click")
	}
}

func TestClickForwardsToTargetNotFocused(t *testing.T) {
	m, left, main := newTestShell(t, true)
	// Focus starts on left. Click the main panel: message must go to main,
	// not the currently-focused left panel.
	_, _ = m.Update(tea.MouseClickMsg(tea.Mouse{Button: tea.MouseLeft, X: 50, Y: 5}))

	if left.gotMouse() {
		t.Error("left (previously focused) panel should not receive a click aimed at main")
	}
	if !main.gotMouse() {
		t.Error("main panel should receive the click at its coordinates")
	}
}

func TestClickOutsidePanelsIgnored(t *testing.T) {
	m, left, main := newTestShell(t, true)
	// Click beyond the window width: no panel, focus unchanged, nothing routed.
	m, _ = m.Update(tea.MouseClickMsg(tea.Mouse{Button: tea.MouseLeft, X: 500, Y: 5}))

	if m.Ctx.FocusedPanel != tuishell.LeftPanel {
		t.Errorf("focus changed on out-of-bounds click: %d", m.Ctx.FocusedPanel)
	}
	if left.gotMouse() || main.gotMouse() {
		t.Error("out-of-bounds click should not be routed to any panel")
	}
}

func TestMouseDisabledRoutesToFocused(t *testing.T) {
	m, left, main := newTestShell(t, false) // mouse disabled
	// With mouse disabled, a click on main's coordinates still goes to the
	// focused (left) panel via the fallback, and focus does not change.
	m, _ = m.Update(tea.MouseClickMsg(tea.Mouse{Button: tea.MouseLeft, X: 50, Y: 5}))

	if m.Ctx.FocusedPanel != tuishell.LeftPanel {
		t.Errorf("focus changed while mouse disabled: %d", m.Ctx.FocusedPanel)
	}
	if !left.gotMouse() {
		t.Error("with mouse disabled, click should fall back to focused (left) panel")
	}
	if main.gotMouse() {
		t.Error("with mouse disabled, non-focused main panel should not receive the click")
	}
}

func TestRenderViewMouseMode(t *testing.T) {
	withMouse, _, _ := newTestShell(t, true)
	if got := withMouse.RenderView().MouseMode; got != tea.MouseModeCellMotion {
		t.Errorf("RenderView().MouseMode = %v, want MouseModeCellMotion", got)
	}

	without, _, _ := newTestShell(t, false)
	if got := without.RenderView().MouseMode; got != tea.MouseModeNone {
		t.Errorf("RenderView().MouseMode = %v, want MouseModeNone", got)
	}
	if !without.RenderView().AltScreen {
		t.Error("RenderView().AltScreen should be true")
	}
}
