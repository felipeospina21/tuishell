// Example app showcasing tuishell's 3-panel layout with async loading.
// Run with: go run ./example
package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/felipeospina21/tuishell"
	"github.com/felipeospina21/tuishell/shell"
	"github.com/felipeospina21/tuishell/style"
	"github.com/felipeospina21/tuishell/table"
)

var theme = style.DefaultTheme()

func main() {
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
		Width(25)

	s := shell.New(shell.Config{
		Theme:          theme,
		LeftPanel:      newProjectsPanel(),
		MainPanel:      &itemsPanel{},
		RightPanel:     &detailsPanel{},
		AppIcon:        "📦",
		Keybinds:       projectsKeys,
		LeftPanelWidth: 25,
		LeftPanelStyle: leftStyle,
	})

	return app{shell: s}
}

func (m app) Init() tea.Cmd { return m.shell.Init() }

func (m app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	// Left panel: user picked a project → start loading items
	case selectProjectMsg:
		if mp, ok := m.shell.Main.(*itemsPanel); ok {
			mp.project = msg.name
			mp.loading = true
		}
		cmds = append(cmds,
			func() tea.Msg { return tuishell.CloseLeftPanelMsg{} },
			func() tea.Msg {
				return tuishell.StartTaskMsg{Cmd: fetchItems(msg.name)}
			},
		)

	// Async result: items fetched
	case itemsFetchedMsg:
		if mp, ok := m.shell.Main.(*itemsPanel); ok {
			mp.loading = false
			mp.items = msg.items
		}
		m.shell.Main, _ = m.shell.Main.Update(msg)
		cmds = append(cmds,
			func() tea.Msg {
				return tuishell.FinishTaskMsg{Keybinds: itemsKeys}
			},
		)

	// Main panel: user pressed enter on a row → open details with loading
	case viewDetailMsg:
		if dp, ok := m.shell.Right.(*detailsPanel); ok {
			dp.title = msg.name
			dp.body = ""
			dp.ready = false
		}
		m.shell.Ctx.FocusedPanel = tuishell.RightPanel
		cmds = append(cmds,
			func() tea.Msg { return tuishell.OpenRightPanelMsg{} },
			func() tea.Msg {
				return tuishell.StartTaskMsg{Cmd: fetchDetail(msg.name)}
			},
		)

	// Async result: detail fetched
	case detailFetchedMsg:
		if dp, ok := m.shell.Right.(*detailsPanel); ok {
			dp.body = msg.body
			dp.ready = true
			dp.viewport.SetContent(msg.body)
		}
		cmds = append(cmds,
			func() tea.Msg {
				return tuishell.FinishTaskMsg{Keybinds: detailsKeys}
			},
		)

	// Details panel: close
	case closeDetailMsg:
		m.shell.Ctx.FocusedPanel = tuishell.MainPanel
		cmds = append(cmds,
			func() tea.Msg { return tuishell.CloseRightPanelMsg{} },
			func() tea.Msg { return tuishell.SetKeybindsMsg{Keybinds: itemsKeys} },
		)
	}

	var cmd tea.Cmd
	m.shell, cmd = m.shell.Update(msg)
	cmds = append(cmds, cmd)

	// Sync spinner to items panel
	if mp, ok := m.shell.Main.(*itemsPanel); ok {
		mp.spinnerView = m.shell.Spinner.View()
	}

	return m, tea.Batch(cmds...)
}

func (m app) View() tea.View { return m.shell.RenderView() }

// ── Messages ────────────────────────────────────────────────────────

type (
	selectProjectMsg struct{ name string }
	itemsFetchedMsg  struct{ items []mockItem }
	viewDetailMsg    struct{ name string }
	detailFetchedMsg struct{ body string }
	closeDetailMsg   struct{}
)

// ── Async commands ──────────────────────────────────────────────────

func fetchItems(project string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(800 * time.Millisecond)
		return itemsFetchedMsg{items: mockItems(project)}
	}
}

func fetchDetail(name string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(600 * time.Millisecond)
		return detailFetchedMsg{body: mockDetail(name)}
	}
}

// ── Keybinds ────────────────────────────────────────────────────────

var projectsKeys = newKeyMap(
	key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
)

var itemsKeys = newKeyMap(
	key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "details")),
	key.NewBinding(key.WithKeys("j"), key.WithHelp("j/k", "navigate")),
)

var detailsKeys = newKeyMap(
	key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "close")),
	key.NewBinding(key.WithKeys("j"), key.WithHelp("j/k", "scroll")),
)

type keyMap struct{ bindings []key.Binding }

func newKeyMap(b ...key.Binding) keyMap {
	return keyMap{bindings: append(b, tuishell.CommonKeys...)}
}
func (k keyMap) ShortHelp() []key.Binding        { return k.bindings }
func (k keyMap) FullHelp() [][]key.Binding        { return [][]key.Binding{k.bindings} }

// ── Left panel (projects) ───────────────────────────────────────────

type projectItem string

func (i projectItem) Title() string       { return string(i) }
func (i projectItem) Description() string { return "" }
func (i projectItem) FilterValue() string { return string(i) }

type projectsPanel struct{ list list.Model }

func newProjectsPanel() *projectsPanel {
	items := []list.Item{
		projectItem("acme-api"),
		projectItem("acme-web"),
		projectItem("acme-mobile"),
		projectItem("acme-infra"),
	}
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Projects"
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.Styles.Title = lipgloss.NewStyle().Foreground(theme.Primary).Bold(true)
	return &projectsPanel{list: l}
}

func (m *projectsPanel) Init() tea.Cmd { return nil }

func (m *projectsPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
		return m, nil
	case tea.KeyPressMsg:
		if msg.Key().Code == tea.KeyEnter {
			if i, ok := m.list.SelectedItem().(projectItem); ok {
				name := string(i)
				return m, func() tea.Msg { return selectProjectMsg{name: name} }
			}
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *projectsPanel) View() tea.View { return tea.NewView(m.list.View()) }

func (m *projectsPanel) SelectedLabel() string {
	if i := m.list.SelectedItem(); i != nil {
		return i.FilterValue()
	}
	return ""
}

// ── Main panel (items table) ────────────────────────────────────────

const (
	tableBorderX  = 2
	tableBorderY  = 2
	headerLines   = 1
	tableOverhead = 1
)

type itemsPanel struct {
	tbl         table.Model
	project     string
	items       []mockItem
	loading     bool
	spinnerView string
	width       int
	height      int
}

func (m *itemsPanel) Init() tea.Cmd { return nil }

func (m *itemsPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.tbl.EmptyMessage = "Select a project"
		tw := m.tableWidth()
		m.tbl.W = tw
		m.tbl.H = m.height - tableBorderY
		if len(m.tbl.Rows()) > 0 {
			h := m.height - headerLines - tableBorderY - tableOverhead
			if h < 3 {
				h = 3
			}
			m.tbl.SetColumns(itemCols(tw))
			m.tbl.SetWidth(tw)
			m.tbl.SetHeight(h)
		}
		return m, nil

	case itemsFetchedMsg:
		rows := make([]table.Row, len(msg.items))
		for i, it := range msg.items {
			rows[i] = table.Row{it.id, it.name, it.status, it.priority}
		}
		tw := m.tableWidth()
		h := m.height - headerLines - tableBorderY - tableOverhead
		if h < 3 {
			h = 3
		}
		cols := itemCols(tw)
		s := table.ThemedStyles(theme)
		m.tbl = table.InitModel(table.InitModelParams{
			Rows: rows, Colums: cols, Styles: &s, Width: tw, Height: h,
		})
		return m, nil

	case tea.KeyPressMsg:
		if msg.String() == "enter" && len(m.items) > 0 {
			idx := m.tbl.Cursor()
			if idx >= 0 && idx < len(m.items) {
				name := m.items[idx].name
				return m, func() tea.Msg { return viewDetailMsg{name: name} }
			}
		}
		var cmd tea.Cmd
		m.tbl, cmd = m.tbl.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *itemsPanel) View() tea.View {
	header := fmt.Sprintf("%s — Items", m.project)
	return tea.NewView(table.RenderPanel(theme, &m.tbl, m.loading, m.spinnerView, header))
}

func (m *itemsPanel) tableWidth() int {
	return m.width - table.DocStyle(theme).GetHorizontalFrameSize() - tableBorderX
}

func itemCols(w int) []table.Column {
	return []table.Column{
		{Title: "ID", Width: 6},
		{Title: "Name", Width: w - 6 - 10 - 8 - 6},
		{Title: "Status", Width: 10},
		{Title: "Priority", Width: 8},
	}
}

// ── Right panel (details) ───────────────────────────────────────────

var detailsPanelStyle = lipgloss.NewStyle().
	MarginTop(1).
	Border(lipgloss.NormalBorder(), true, false, true, true).
	BorderForeground(theme.Border)

var (
	hdrBorder = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b)
	}()
	ftrBorder = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return hdrBorder.BorderStyle(b)
	}()
)

type detailsPanel struct {
	viewport viewport.Model
	title    string
	body     string
	ready    bool
}

func (m *detailsPanel) Init() tea.Cmd { return nil }

func (m *detailsPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		frameY := detailsPanelStyle.GetVerticalFrameSize()
		h := msg.Height - frameY - 3 - 3 - 2 // header(3) + footer(3) + 2 newlines
		if h < 1 {
			h = 1
		}
		m.viewport = viewport.New(viewport.WithWidth(msg.Width), viewport.WithHeight(h))
		if m.body != "" {
			m.viewport.SetContent(m.body)
		}
		return m, nil
	case tea.KeyPressMsg:
		if msg.String() == "esc" || msg.String() == "q" {
			return m, func() tea.Msg { return closeDetailMsg{} }
		}
	}
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m *detailsPanel) View() tea.View {
	if !m.ready {
		return tea.NewView(detailsPanelStyle.Render("Loading..."))
	}

	t := m.title
	if t == "" {
		t = "Details"
	}
	hdr := hdrBorder.Render(t)
	hline := strings.Repeat("─", max(0, m.viewport.Width()-lipgloss.Width(hdr)))
	header := lipgloss.JoinHorizontal(lipgloss.Center, hdr, hline)

	info := ftrBorder.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	fline := strings.Repeat("─", max(0, m.viewport.Width()-lipgloss.Width(info)))
	footer := lipgloss.JoinHorizontal(lipgloss.Center, fline, info)

	return tea.NewView(detailsPanelStyle.Render(
		fmt.Sprintf("%s\n%s\n%s", header, m.viewport.View(), footer),
	))
}

// ── Mock data ───────────────────────────────────────────────────────

type mockItem struct {
	id, name, status, priority string
}

func mockItems(project string) []mockItem {
	return []mockItem{
		{"001", "Fix login timeout", "Open", "High"},
		{"002", "Add pagination", "In Progress", "Medium"},
		{"003", "Update API docs", "Done", "Low"},
		{"004", "Migrate DB schema", "In Review", "High"},
		{"005", "Refactor auth module", "Open", "Medium"},
		{"006", "Add rate limiting", "Open", "High"},
		{"007", "Fix memory leak", "In Progress", "Critical"},
		{"008", "Update dependencies", "Done", "Low"},
		{"009", "Add health check", "Open", "Medium"},
		{"010", "Improve logging", "In Review", "Low"},
	}
}

func mockDetail(name string) string {
	label := lipgloss.NewStyle().Foreground(theme.TextDimmed)
	value := lipgloss.NewStyle().Foreground(theme.Text)
	section := lipgloss.NewStyle().Bold(true).Foreground(theme.Primary).MarginTop(1)

	var b strings.Builder
	b.WriteString(section.Render("Summary"))
	b.WriteString("\n")
	b.WriteString(value.Render(name))
	b.WriteString("\n\n")
	b.WriteString(label.Render("Status: ") + value.Render("Open"))
	b.WriteString("\n")
	b.WriteString(label.Render("Priority: ") + value.Render("High"))
	b.WriteString("\n")
	b.WriteString(label.Render("Assignee: ") + value.Render("Jane Doe"))
	b.WriteString("\n")
	b.WriteString(label.Render("Created: ") + value.Render("3 days ago"))
	b.WriteString("\n\n")
	b.WriteString(section.Render("Description"))
	b.WriteString("\n")
	b.WriteString(value.Render("This is a detailed view for the selected item. "))
	b.WriteString(value.Render("It demonstrates the right panel with a viewport, "))
	b.WriteString(value.Render("header/footer, and the panel-owned style pattern."))
	b.WriteString("\n\n")
	b.WriteString(section.Render("Activity"))
	b.WriteString("\n")
	for i := 1; i <= 15; i++ {
		b.WriteString(value.Render(fmt.Sprintf("  • Comment #%d — sample activity for scroll test\n", i)))
	}
	return b.String()
}
