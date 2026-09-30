// Example app showcasing tuishell's mouse-support hooks (issue #31):
//
//   - shell.Config.EnableMouse enables MouseModeCellMotion and coordinate
//     routing / click-to-focus.
//   - Click a panel to focus it (click-to-focus).
//   - Click a nav item (left panel) to select it — bubblezone-marked rows.
//   - Scroll the table with the mouse wheel.
//   - Click a table row (main panel) to select it AND open the details panel,
//     via bubblezone (zone.Mark on rows + zone.Scan at the root View).
//
// Run from this directory with: go run .
package main

import (
	"fmt"
	"os"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"

	"github.com/felipeospina21/tuishell"
	"github.com/felipeospina21/tuishell/shell"
	"github.com/felipeospina21/tuishell/table"
)

var theme = defaultTheme()

func main() {
	// bubblezone: initialize the global manager. zone.Scan (in View) strips the
	// markers; zone.Get(id).InBounds(mouse) does the hit-testing.
	zone.NewGlobal()

	p := tea.NewProgram(newApp())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

// ── App ─────────────────────────────────────────────────────────────

type app struct {
	shell shell.Model
}

func newApp() tea.Model {
	leftStyle := lipgloss.NewStyle().
		PaddingRight(2).
		Border(lipgloss.NormalBorder(), false, true, false, false).
		BorderForeground(theme.Border).
		Width(24)

	rightStyle := lipgloss.NewStyle().
		PaddingLeft(2).
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(theme.Border)

	s := shell.New(shell.Config{
		Theme:           theme,
		LeftPanel:       newNavPanel(),
		MainPanel:       newTablePanel(),
		RightPanel:      newDetailsPanel(),
		AppIcon:         "🖱️",
		Keybinds:        demoKeys,
		LeftPanelWidth:  24,
		LeftPanelStyle:  leftStyle,
		RightPanelStyle: rightStyle,
		EnableMouse:     true, // ← the feature: mouse mode + routing + click-to-focus
		// Button labels are consumer-defined (the panel content is app-specific).
		LeftButtonLabel:  "menu",
		RightButtonLabel: "info",
	})

	return app{shell: s}
}

func (m app) Init() tea.Cmd { return m.shell.Init() }

func (m app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	prevFocus := m.shell.Ctx.FocusedPanel

	// App-level messages emitted by panels on mouse interaction.
	switch msg := msg.(type) {
	case openDetailsMsg:
		if dp, ok := m.shell.Right.(*detailsPanel); ok {
			dp.row = msg.row
		}
		// The shell lays out left+main OR main+right (not all three at full
		// width), so close the left panel before opening the right one — the
		// same pattern tron uses. Each shell message recomputes layout and
		// pushes new sizes to the panels.
		var c1, c2 tea.Cmd
		if m.shell.IsLeftOpen() {
			m.shell, c1 = m.shell.Update(tuishell.CloseLeftPanelMsg{})
		}
		m.shell, c2 = m.shell.Update(tuishell.OpenRightPanelMsg{})
		cmds := []tea.Cmd{c1, c2}
		cmds = append(cmds, m.broadcastFocus()...)
		return m, tea.Batch(cmds...)
	}

	var cmd tea.Cmd
	m.shell, cmd = m.shell.Update(msg)
	cmds := []tea.Cmd{cmd}

	// The shell updates Ctx.FocusedPanel on click-to-focus. When it changes,
	// tell every panel so they can render a focus indicator.
	if m.shell.Ctx.FocusedPanel != prevFocus {
		cmds = append(cmds, m.broadcastFocus()...)
	}

	return m, tea.Batch(cmds...)
}

func (m *app) broadcastFocus() []tea.Cmd {
	fc := focusChangedMsg{panel: m.shell.Ctx.FocusedPanel}
	var cmds []tea.Cmd
	var c tea.Cmd
	m.shell.Left, c = m.shell.Left.Update(fc)
	cmds = append(cmds, c)
	m.shell.Main, c = m.shell.Main.Update(fc)
	cmds = append(cmds, c)
	if m.shell.Right != nil {
		m.shell.Right, c = m.shell.Right.Update(fc)
		cmds = append(cmds, c)
	}
	return cmds
}

// View wraps the shell output in zone.Scan so bubblezone can register the
// markers. Scan must be called only at the root model.
func (m app) View() tea.View {
	v := m.shell.RenderView() // AltScreen + MouseModeCellMotion (EnableMouse)
	v.SetContent(zone.Scan(v.Content))
	return v
}

// ── Messages ────────────────────────────────────────────────────────

type (
	focusChangedMsg struct{ panel tuishell.FocusedPanel }
	openDetailsMsg  struct{ row int }
)

// ── Keybinds (shown in the statusline / help) ───────────────────────

var demoKeys = newKeyMap(
	key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab/click", "focus panel")),
	key.NewBinding(key.WithKeys("j"), key.WithHelp("wheel", "scroll table")),
	key.NewBinding(key.WithKeys("enter"), key.WithHelp("click row", "open details")),
	key.NewBinding(key.WithKeys(" "), key.WithHelp("[«]/[»]", "toggle panels")),
)

type keyMap struct{ bindings []key.Binding }

func newKeyMap(b ...key.Binding) keyMap {
	return keyMap{bindings: append(b, tuishell.CommonKeys...)}
}
func (k keyMap) ShortHelp() []key.Binding  { return k.bindings }
func (k keyMap) FullHelp() [][]key.Binding { return [][]key.Binding{k.bindings} }

// ── Left panel (bubblezone-marked navigation) ───────────────────────

const navZonePrefix = "nav-item-"

var navItems = []string{"Dashboard", "Issues", "Pipelines", "Settings"}

type navPanel struct {
	cursor  int
	focused bool
}

func newNavPanel() *navPanel {
	return &navPanel{focused: true} // shell focuses LeftPanel by default
}

func (m *navPanel) Init() tea.Cmd { return nil }

func (m *navPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case focusChangedMsg:
		m.focused = msg.panel == tuishell.LeftPanel
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(navItems)-1 {
				m.cursor++
			}
		}
		return m, nil

	case tea.MouseClickMsg:
		if msg.Button == tea.MouseLeft {
			for i := range navItems {
				if zone.Get(fmt.Sprintf("%s%d", navZonePrefix, i)).InBounds(msg) {
					m.cursor = i
					return m, nil
				}
			}
		}
	}
	return m, nil
}

func (m *navPanel) View() tea.View {
	title := lipgloss.NewStyle().Foreground(theme.Primary).Bold(true).Render("Navigate")

	lines := make([]string, len(navItems))
	for i, label := range navItems {
		style := lipgloss.NewStyle().Padding(0, 1)
		text := "  " + label
		if i == m.cursor {
			style = style.Foreground(theme.PrimaryFg).Background(theme.PrimaryDim).Bold(true)
			text = "› " + label
		} else {
			style = style.Foreground(theme.Text)
		}
		// Mark each item so clicks map to it (bubblezone).
		lines[i] = zone.Mark(fmt.Sprintf("%s%d", navZonePrefix, i), style.Render(text))
	}

	hint := lipgloss.NewStyle().Foreground(theme.TextDimmed).
		Render("\nClick an item, or\nclick a panel to\nfocus it.")

	body := lipgloss.JoinVertical(lipgloss.Left,
		focusTag(m.focused),
		title,
		strings.Join(lines, "\n"),
		hint,
	)
	return tea.NewView(body)
}

func (m *navPanel) SelectedLabel() string { return navItems[m.cursor] }

// ── Main panel (bubblezone table) ───────────────────────────────────

const rowZonePrefix = "table-row-"

type tablePanel struct {
	tbl        table.Model
	focused    bool
	lastAction string
	width      int
	height     int
}

func newTablePanel() *tablePanel {
	rows := make([]table.Row, len(mockRows))
	for i, r := range mockRows {
		rows[i] = table.Row{r.id, r.name, r.status}
	}
	s := table.ThemedStyles(theme)
	tbl := table.New(
		table.WithColumns(tableCols(60)),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithStyles(s),
		// Mouse hooks from issue #31:
		table.WithZoneManager(zone.DefaultManager), // consumer supplies the manager
		table.WithRowID(func(i int) string { // unique zone id per row
			return fmt.Sprintf("%s%d", rowZonePrefix, i)
		}),
		table.WithAllowMouseWhenBlurred(true), // wheel-scroll even when unfocused
	)
	return &tablePanel{tbl: tbl}
}

func (m *tablePanel) Init() tea.Cmd { return nil }

func (m *tablePanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.resizeTable()
		return m, nil

	case focusChangedMsg:
		m.focused = msg.panel == tuishell.MainPanel
		return m, nil

	case tea.MouseWheelMsg:
		// Let the table's built-in wheel handler scroll, then report.
		var cmd tea.Cmd
		m.tbl, cmd = m.tbl.Update(msg)
		m.lastAction = "wheel scroll → row " + fmt.Sprint(m.tbl.Cursor())
		return m, cmd

	case tea.MouseClickMsg:
		return m, m.handleClick(msg)

	case tea.MouseReleaseMsg:
		// A left click emits both a click and a release. We hit-test on the
		// click; swallow the release so it does not reach the table's built-in
		// geometric handler (which would re-select using absolute coordinates).
		return m, nil
	}

	var cmd tea.Cmd
	m.tbl, cmd = m.tbl.Update(msg)
	return m, cmd
}

// handleClick uses bubblezone to find the clicked row, selects it, and opens
// the details panel for it.
func (m *tablePanel) handleClick(msg tea.MouseClickMsg) tea.Cmd {
	if msg.Button != tea.MouseLeft {
		return nil
	}
	for i := range m.tbl.Rows() {
		if zone.Get(fmt.Sprintf("%s%d", rowZonePrefix, i)).InBounds(msg) {
			m.tbl.SetCursor(i)
			m.lastAction = "clicked row " + fmt.Sprint(i) + " (" + mockRows[i].name + ")"
			row := i
			return func() tea.Msg { return openDetailsMsg{row: row} }
		}
	}
	return nil
}

func (m *tablePanel) View() tea.View {
	// Chrome above the table: focus tag (1) + title (1). Status line below (1).
	title := lipgloss.NewStyle().Foreground(theme.Primary).Bold(true).
		Render("Items — wheel to scroll, click a row to open details")
	action := m.lastAction
	if action == "" {
		action = "…interact with the mouse"
	}
	status := lipgloss.NewStyle().Foreground(theme.TextDimmed).Render("Last: " + action)

	// Render the table plainly. We do NOT wrap it in an extra bordered style:
	// the rows carry zero-width bubblezone markers, and wrapping marked content
	// in another Border()/Width() style shifts the marker offsets so hit-testing
	// misfires. The table's header already draws its own bottom border.
	body := lipgloss.JoinVertical(lipgloss.Left,
		focusTag(m.focused),
		title,
		m.tbl.View(),
		status,
	)
	return tea.NewView(body)
}

// mainChromeY = focus tag (1) + title (1) + status (1) + header row (1) +
// header bottom border (1).
const mainChromeY = 5

func (m *tablePanel) resizeTable() {
	w := m.width - 2
	if w < 20 {
		w = 20
	}
	h := m.height - mainChromeY
	if h < 3 {
		h = 3
	}
	m.tbl.SetColumns(tableCols(w))
	m.tbl.SetWidth(w)
	m.tbl.SetHeight(h)
	m.tbl.W = w
	m.tbl.H = h
}

func tableCols(w int) []table.Column {
	nameW := w - 6 - 14
	if nameW < 8 {
		nameW = 8
	}
	return []table.Column{
		{Title: "ID", Width: 6},
		{Title: "Name", Width: nameW},
		{Title: "Status", Width: 14},
	}
}

// ── Right panel (details, opened on row click) ──────────────────────

type detailsPanel struct {
	row     int
	focused bool
	width   int
	height  int
}

func newDetailsPanel() *detailsPanel { return &detailsPanel{row: -1} }

func (m *detailsPanel) Init() tea.Cmd { return nil }

func (m *detailsPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil
	case focusChangedMsg:
		m.focused = msg.panel == tuishell.RightPanel
		return m, nil
	case openDetailsMsg:
		m.row = msg.row
		return m, nil
	}
	return m, nil
}

// vpWidth is the content width inside the panel's left-border separator.
func (m *detailsPanel) vpWidth() int {
	w := m.width - 1
	if w < 12 {
		return 12
	}
	return w
}

func (m *detailsPanel) View() tea.View {
	if m.height <= 0 {
		return tea.NewView("")
	}
	w := m.vpWidth()

	// Active-panel indicator: recolor our own left-border when focused, to
	// match the shell's built-in indicator (theme.Primary == the shell's
	// default FocusBorderColor). The shell recolors the left panel + main
	// frame automatically; the right panel renders its own content, so we
	// self-indicate here.
	borderColor := theme.Border
	if m.focused {
		borderColor = theme.Primary
	}

	// Left-border separator between main and right panel (tron pattern).
	container := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(borderColor).
		Width(w).
		Height(m.height)

	title := lipgloss.NewStyle().Foreground(theme.Primary).Bold(true).
		MaxWidth(w).Inline(true).Render("Details")

	if m.row < 0 || m.row >= len(mockRows) {
		hint := lipgloss.NewStyle().Foreground(theme.TextDimmed).Width(w).
			Render("Click a table row\nto open its details.")
		return tea.NewView(container.Render(lipgloss.JoinVertical(lipgloss.Left,
			focusTag(m.focused), title, "", hint,
		)))
	}

	r := mockRows[m.row]
	label := lipgloss.NewStyle().Foreground(theme.TextDimmed)
	value := lipgloss.NewStyle().Foreground(theme.Text)
	line := func(k, v string) string {
		return lipgloss.NewStyle().Width(w).Inline(true).MaxWidth(w).
			Render(label.Render(k) + value.Render(v))
	}

	body := lipgloss.JoinVertical(lipgloss.Left,
		focusTag(m.focused),
		title,
		"",
		line("ID:     ", r.id),
		line("Name:   ", r.name),
		line("Status: ", r.status),
		"",
		lipgloss.NewStyle().Foreground(theme.TextDimmed).Width(w).
			Render("Opened from a mouse click on the table row."),
	)
	return tea.NewView(container.Render(body))
}

// ── Focus indicator ─────────────────────────────────────────────────

func focusTag(focused bool) string {
	if focused {
		return lipgloss.NewStyle().Foreground(theme.SuccessBright).Bold(true).Render("● focused")
	}
	return lipgloss.NewStyle().Foreground(theme.Dim).Render("○ click to focus")
}

// ── Mock data ───────────────────────────────────────────────────────

type mockRow struct{ id, name, status string }

var mockRows = []mockRow{
	{"001", "Fix login timeout", "Open"},
	{"002", "Add pagination", "In Progress"},
	{"003", "Update API docs", "Done"},
	{"004", "Migrate DB schema", "In Review"},
	{"005", "Refactor auth module", "Open"},
	{"006", "Add rate limiting", "Open"},
	{"007", "Fix memory leak", "In Progress"},
	{"008", "Update dependencies", "Done"},
	{"009", "Add health check", "Open"},
	{"010", "Improve logging", "In Review"},
	{"011", "Cache invalidation", "Open"},
	{"012", "Add metrics export", "In Progress"},
	{"013", "Fix flaky test", "Done"},
	{"014", "Upgrade runtime", "In Review"},
	{"015", "Add dark mode", "Open"},
}
