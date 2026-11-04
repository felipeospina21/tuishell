package shell

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/felipeospina21/tuishell"
	"github.com/felipeospina21/tuishell/style"
)

func newShell(t *testing.T, cfg Config) (Model, *fakePanel, *fakePanel, *fakePanel) {
	t.Helper()
	left := &fakePanel{name: "left"}
	main := &fakePanel{name: "main"}
	right := &fakePanel{name: "right"}
	if cfg.Theme == (style.Theme{}) {
		cfg.Theme = style.Theme{}
	}
	cfg.LeftPanel = left
	cfg.MainPanel = main
	cfg.RightPanel = right
	if cfg.Keybinds == nil {
		cfg.Keybinds = tuishell.GlobalKeys(false)
	}
	if cfg.LeftPanelWidth == 0 {
		cfg.LeftPanelWidth = 30
	}
	cfg.LeftPanelStyle = lipgloss.NewStyle()
	cfg.MainPanelStyle = lipgloss.NewStyle()
	m := New(cfg)
	m, _ = m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	return m, left, main, right
}

func keyPress(s string) tea.KeyPressMsg {
	return tea.KeyPressMsg{Code: tea.KeyTab, Text: s}
}

func TestFixedPanels_StartOpenAndCannotClose(t *testing.T) {
	m, _, _, _ := newShell(t, Config{FixedPanels: true})

	if !m.IsLeftOpen() || !m.IsRightOpen() {
		t.Fatalf("fixed panels should start with both open: left=%v right=%v", m.IsLeftOpen(), m.IsRightOpen())
	}

	// Close messages must be no-ops.
	m, _ = m.Update(tuishell.CloseLeftPanelMsg{})
	m, _ = m.Update(tuishell.CloseRightPanelMsg{})
	if !m.IsLeftOpen() || !m.IsRightOpen() {
		t.Errorf("fixed panels closed via message: left=%v right=%v", m.IsLeftOpen(), m.IsRightOpen())
	}

	// Global toggle/close keys must be no-ops too.
	gk := tuishell.GlobalKeys(false)
	m, _, _ = m.handleGlobalKeys(tea.KeyPressMsg{Code: 'o', Mod: tea.ModCtrl}) // ToggleLeftPanel (ctrl+o)
	_ = gk
	if !m.IsLeftOpen() {
		t.Error("fixed left panel closed via ctrl+o")
	}
}

func TestFixedPanels_ButtonBarHidden(t *testing.T) {
	// Fixed + mouse: no button bar (nothing to toggle).
	m, _, _, _ := newShell(t, Config{FixedPanels: true, EnableMouse: true})
	if m.showButtonBar() {
		t.Error("button bar should be hidden when panels are fixed")
	}
	if strings.Contains(m.View(), "nav ]") {
		t.Error("View should not render the toggle bar in fixed mode")
	}
}

func TestButtonBar_ShownOnlyWithMouseAndNotFixed(t *testing.T) {
	withMouse, _, _, _ := newShell(t, Config{EnableMouse: true})
	if !withMouse.showButtonBar() {
		t.Error("button bar should show when mouse enabled and not fixed")
	}
	if !strings.Contains(withMouse.View(), "nav ]") {
		t.Error("View should render the toggle bar when mouse enabled")
	}

	noMouse, _, _, _ := newShell(t, Config{EnableMouse: false})
	if noMouse.showButtonBar() {
		t.Error("button bar should be hidden when mouse disabled")
	}
}

func TestButtonLabels_Configurable(t *testing.T) {
	m, _, _, _ := newShell(t, Config{
		EnableMouse:      true,
		LeftButtonLabel:  "projects",
		RightButtonLabel: "info",
	})

	view := m.View()
	if !strings.Contains(view, "projects ]") {
		t.Error("custom left button label 'projects' not rendered")
	}
	if !strings.Contains(view, "info") {
		t.Error("custom right button label 'info' not rendered")
	}

	// The right toggle rect must track the (shorter) custom label so a click at
	// the rendered glyph still toggles the panel.
	rt := m.Layout.RightToggleRect
	// "[ info « ]" is 10 cells wide.
	if rt.Width != 10 {
		t.Errorf("RightToggleRect.Width = %d, want 10 for label 'info'", rt.Width)
	}
	m, _ = m.Update(tea.MouseClickMsg(tea.Mouse{Button: tea.MouseLeft, X: rt.X + 1, Y: rt.Y}))
	if !m.IsRightOpen() {
		t.Error("clicking the custom-labeled right toggle did not open the panel")
	}
}

func TestButtonLabels_DefaultWhenEmpty(t *testing.T) {
	m, _, _, _ := newShell(t, Config{EnableMouse: true})
	view := m.View()
	if !strings.Contains(view, "nav ]") || !strings.Contains(view, "details »") && !strings.Contains(view, "details «") {
		t.Errorf("default labels 'nav'/'details' not rendered:\n%s", view)
	}
}

func TestButtonClick_TogglesPanels(t *testing.T) {
	m, _, _, _ := newShell(t, Config{EnableMouse: true})
	// Left starts open; click the left toggle to close it.
	lt := m.Layout.LeftToggleRect
	m, _ = m.Update(tea.MouseClickMsg(tea.Mouse{Button: tea.MouseLeft, X: lt.X, Y: lt.Y}))
	if m.IsLeftOpen() {
		t.Error("clicking left toggle did not close the left panel")
	}

	// Right starts closed; click the right toggle to open it.
	rt := m.Layout.RightToggleRect
	m, _ = m.Update(tea.MouseClickMsg(tea.Mouse{Button: tea.MouseLeft, X: rt.X + 1, Y: rt.Y}))
	if !m.IsRightOpen() {
		t.Error("clicking right toggle did not open the right panel")
	}
}

func TestPanelCycle_ForwardAndBack(t *testing.T) {
	// Open all three so cycling has three stops.
	m, _, _, _ := newShell(t, Config{FixedPanels: true})
	// left -> main -> right -> left
	m, _, _ = m.handleGlobalKeys(cycleKey())
	if m.FocusedPanel() != tuishell.MainPanel {
		t.Fatalf("after 1 cycle focus = %d, want MainPanel", m.FocusedPanel())
	}
	m, _, _ = m.handleGlobalKeys(cycleKey())
	if m.FocusedPanel() != tuishell.RightPanel {
		t.Fatalf("after 2 cycles focus = %d, want RightPanel", m.FocusedPanel())
	}
	m, _, _ = m.handleGlobalKeys(cycleKey())
	if m.FocusedPanel() != tuishell.LeftPanel {
		t.Fatalf("after 3 cycles focus = %d, want LeftPanel (wrap)", m.FocusedPanel())
	}
	// Back one: left -> right.
	m, _, _ = m.handleGlobalKeys(cycleBackKey())
	if m.FocusedPanel() != tuishell.RightPanel {
		t.Fatalf("after cycle back focus = %d, want RightPanel", m.FocusedPanel())
	}
}

func TestPanelCycle_SkipsClosedPanels(t *testing.T) {
	// Default: left+main open, right closed. Cycle should be left <-> main.
	m, _, _, _ := newShell(t, Config{})
	m, _ = m.Update(tuishell.CloseRightPanelMsg{}) // ensure right closed, focus main
	m.Ctx.FocusedPanel = tuishell.LeftPanel
	m, _, _ = m.handleGlobalKeys(cycleKey())
	if m.FocusedPanel() != tuishell.MainPanel {
		t.Fatalf("cycle from left = %d, want MainPanel", m.FocusedPanel())
	}
	m, _, _ = m.handleGlobalKeys(cycleKey())
	if m.FocusedPanel() != tuishell.LeftPanel {
		t.Fatalf("cycle wrapped = %d, want LeftPanel (right is closed, skipped)", m.FocusedPanel())
	}
}

func TestFocusBorder_RecolorsFocusedPanelOnly(t *testing.T) {
	red := lipgloss.Color("#ff0000")
	m, _, _, _ := newShell(t, Config{FocusedBorderColor: red})

	bordered := lipgloss.NewStyle().Border(lipgloss.NormalBorder())

	// Focused panel gets the focus color.
	m.Ctx.FocusedPanel = tuishell.LeftPanel
	got := m.focusBorder(bordered, tuishell.LeftPanel)
	if got.GetBorderTopForeground() != red {
		t.Errorf("focused border color = %v, want red", got.GetBorderTopForeground())
	}

	// Non-focused panel is unchanged.
	unchanged := m.focusBorder(bordered, tuishell.MainPanel)
	if unchanged.GetBorderTopForeground() == red {
		t.Error("non-focused panel border should not be recolored")
	}

	// A borderless style is never recolored.
	plain := m.focusBorder(lipgloss.NewStyle(), tuishell.LeftPanel)
	if plain.GetBorderTopForeground() == red {
		t.Error("borderless style should not gain a focus color")
	}
}

func cycleKey() tea.KeyPressMsg     { return tea.KeyPressMsg{Code: tea.KeyTab} }
func cycleBackKey() tea.KeyPressMsg { return tea.KeyPressMsg{Code: tea.KeyTab, Mod: tea.ModShift} }

var _ = keyPress // keyPress kept for potential future text-key tests

// TestPanelCycle_ConfigurableKeys verifies that overriding the cycle keys moves
// the focus-cycle onto a custom key and that the default "tab" no longer cycles
// (so the consumer can reuse tab for in-panel tab switching).
func TestPanelCycle_ConfigurableKeys(t *testing.T) {
	m, _, _, _ := newShell(t, Config{
		FixedPanels:        true, // all three panels open
		CyclePanelKeys:     []string{"ctrl+l"},
		CyclePanelBackKeys: []string{"ctrl+h"},
	})
	m.Ctx.FocusedPanel = tuishell.LeftPanel

	// The custom forward key cycles: left -> main.
	ctrlL := tea.KeyPressMsg{Code: 'l', Mod: tea.ModCtrl}
	m, _, handled := m.handleGlobalKeys(ctrlL)
	if !handled {
		t.Fatal("ctrl+l should be handled as the cycle key")
	}
	if m.FocusedPanel() != tuishell.MainPanel {
		t.Fatalf("after ctrl+l focus = %d, want MainPanel", m.FocusedPanel())
	}

	// The custom back key cycles: main -> left.
	ctrlH := tea.KeyPressMsg{Code: 'h', Mod: tea.ModCtrl}
	m, _, handled = m.handleGlobalKeys(ctrlH)
	if !handled {
		t.Fatal("ctrl+h should be handled as the cycle-back key")
	}
	if m.FocusedPanel() != tuishell.LeftPanel {
		t.Fatalf("after ctrl+h focus = %d, want LeftPanel", m.FocusedPanel())
	}

	// The default tab must NOT cycle when overridden — it falls through
	// unhandled so the consumer can route it to in-panel tab switching.
	before := m.FocusedPanel()
	m, _, handled = m.handleGlobalKeys(cycleKey())
	if handled {
		t.Error("tab should not be handled as cycle when keys are overridden")
	}
	if m.FocusedPanel() != before {
		t.Errorf("tab changed focus (%d -> %d) despite override", before, m.FocusedPanel())
	}
}

func TestHelpKeyOpensScrollableSingleColumnModal(t *testing.T) {
	m, _, _, _ := newShell(t, Config{})

	m, cmd := m.Update(tea.KeyPressMsg{Code: '?', Text: "?"})
	if cmd == nil {
		t.Fatal("? should produce an OpenModalMsg command")
	}
	open, ok := cmd().(tuishell.OpenModalMsg)
	if !ok {
		t.Fatalf("expected OpenModalMsg, got %T", cmd())
	}
	m, _ = m.Update(open)
	if !m.IsModalOpen() {
		t.Fatal("help modal should be open")
	}
	if !strings.Contains(open.Content, "toggle help") {
		t.Errorf("help content missing global bindings:\n%s", open.Content)
	}
	for _, l := range strings.Split(open.Content, "\n") {
		if strings.Contains(l, "toggle help") && strings.Contains(l, "quit") {
			t.Errorf("bindings must be one per line, got %q", l)
		}
	}
}
