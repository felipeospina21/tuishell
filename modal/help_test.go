package modal

import (
	"fmt"
	"strings"
	"testing"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/felipeospina21/tuishell"
	"github.com/felipeospina21/tuishell/style"
)

type fakeKeyMap struct{ groups [][]key.Binding }

func (k fakeKeyMap) ShortHelp() []key.Binding  { return nil }
func (k fakeKeyMap) FullHelp() [][]key.Binding { return k.groups }

var _ help.KeyMap = fakeKeyMap{}

func header(s string) key.Binding {
	return key.NewBinding(key.WithKeys(""), key.WithHelp(s, ""))
}

func bind(k, d string) key.Binding {
	return key.NewBinding(key.WithKeys(k), key.WithHelp(k, d))
}

func newModal(w, h int) Model {
	ctx := &tuishell.AppContext{Window: tea.WindowSizeMsg{Width: w, Height: h}}
	return New(ctx, style.Theme{})
}

func TestRenderHelpScrollable_SingleColumnWithHeaders(t *testing.T) {
	off := bind("x", "hidden")
	off.SetEnabled(false)
	km := fakeKeyMap{groups: [][]key.Binding{
		{header("── Global ──"), header("Modifiers: note"), bind("?", "help"), off},
		{header("── Actions ──"), bind("ctrl+r", "restart daemon")},
	}}

	out := newModal(100, 30).RenderHelpScrollable(km)
	lines := strings.Split(out, "\n")

	find := func(sub string) int {
		for i, l := range lines {
			if strings.Contains(l, sub) {
				return i
			}
		}
		return -1
	}
	g, n, h, a, r := find("── Global ──"), find("Modifiers: note"), find("help"), find("── Actions ──"), find("restart daemon")
	if g < 0 || n < 0 || h < 0 || a < 0 || r < 0 {
		t.Fatalf("missing lines in output:\n%s", out)
	}
	if !(g < n && n < h && h < a && a < r) {
		t.Errorf("groups not stacked in order: %d %d %d %d %d", g, n, h, a, r)
	}
	if strings.TrimSpace(lines[a-1]) != "" {
		t.Errorf("expected blank line between groups, got %q", lines[a-1])
	}
	if strings.Contains(out, "hidden") {
		t.Error("disabled binding must be skipped")
	}
	// Descriptions line up in one column: key column is padded to widest key.
	if lipgloss.Width(lines[h][:strings.Index(lines[h], "help")]) != lipgloss.Width(lines[r][:strings.Index(lines[r], "restart")]) {
		t.Error("descriptions are not aligned across rows")
	}
}

func TestModal_ScrollsLongContent(t *testing.T) {
	m := newModal(100, 20)
	var rows []string
	for i := 0; i < 80; i++ {
		rows = append(rows, fmt.Sprintf("row %d", i))
	}
	m.Content = strings.Join(rows, "\n")

	view := func() string { return m.View(strings.Repeat(strings.Repeat(" ", 100)+"\n", 20)) }

	if !strings.Contains(view(), "row 0") || strings.Contains(view(), "row 79") {
		t.Fatal("expected to start at the top of the content")
	}
	if !strings.Contains(view(), "j/k scroll") {
		t.Error("footer should advertise scrolling when content overflows")
	}

	m, _ = m.Update(tea.KeyPressMsg{Code: 'G', Text: "G"})
	if !strings.Contains(view(), "row 79") || strings.Contains(view(), "row 0\n") {
		t.Error("G should jump to the end")
	}
	m, _ = m.Update(tea.KeyPressMsg{Code: 'g', Text: "g"})
	if !strings.Contains(view(), "row 0") {
		t.Error("g should jump to the top")
	}
	m, _ = m.Update(tea.KeyPressMsg{Code: 'j', Text: "j"})
	if m.scroll != 1 {
		t.Errorf("j should scroll one line, scroll=%d", m.scroll)
	}
	m, _ = m.Update(tea.KeyPressMsg{Code: 'k', Text: "k"})
	m, _ = m.Update(tea.KeyPressMsg{Code: 'k', Text: "k"})
	if m.scroll != 0 {
		t.Errorf("scroll must not go negative, scroll=%d", m.scroll)
	}

	m.scroll = 5
	m.ResetScroll()
	if m.scroll != 0 {
		t.Error("ResetScroll should return to top")
	}
}

func TestModal_ShortContentHasNoScrollHint(t *testing.T) {
	m := newModal(100, 30)
	m.Content = "short"
	v := m.View(strings.Repeat(strings.Repeat(" ", 100)+"\n", 30))
	if strings.Contains(v, "j/k scroll") {
		t.Error("no scroll hint expected for short content")
	}
}
