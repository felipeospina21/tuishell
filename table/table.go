// Package table provides a styled, scrollable table widget for bubbletea TUIs.
package table

import (
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/felipeospina21/tuishell"
)

const (
	rowBottomMargin = 1
	rowHeight       = 1
)

// Model defines a state for the table widget.
type Model struct {
	KeyMap       KeyMap
	EmptyMessage string
	W            int
	H            int

	cols      []Column
	rows      []Row
	cursor    int
	focus     bool
	styles    Styles
	styleFunc StyleFunc

	viewport viewport.Model
	start    int
	end      int
}

// Row represents one line in the table.
type Row []string

// Column defines the table structure.
type Column struct {
	Title    string
	Width    int
	Name     string
	Centered bool
}

// KeyMap defines keybindings.
type KeyMap struct {
	LineUp       key.Binding
	LineDown     key.Binding
	PageUp       key.Binding
	PageDown     key.Binding
	HalfPageUp   key.Binding
	HalfPageDown key.Binding
	GotoTop      key.Binding
	GotoBottom   key.Binding
}

// ShortHelp implements the KeyMap interface.
func (km KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{km.LineUp, km.LineDown}
}

// FullHelp implements the KeyMap interface.
func (km KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{km.LineUp, km.LineDown, km.GotoTop, km.GotoBottom},
		{km.PageUp, km.PageDown, km.HalfPageUp, km.HalfPageDown},
	}
}

// DefaultKeyMap returns a default set of keybindings.
func DefaultKeyMap() KeyMap {
	const spacebar = " "
	return KeyMap{
		LineUp: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		LineDown: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("b", "pgup"),
			key.WithHelp("b/pgup", "page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("f", "pgdown", spacebar),
			key.WithHelp("f/pgdn", "page down"),
		),
		HalfPageUp: key.NewBinding(
			key.WithKeys("u", "ctrl+u"),
			key.WithHelp("u", "½ page up"),
		),
		HalfPageDown: key.NewBinding(
			key.WithKeys("d", "ctrl+d"),
			key.WithHelp("d", "½ page down"),
		),
		GotoTop: key.NewBinding(
			key.WithKeys("home", "g"),
			key.WithHelp("g/home", "go to start"),
		),
		GotoBottom: key.NewBinding(
			key.WithKeys("end", "G"),
			key.WithHelp("G/end", "go to end"),
		),
	}
}

// Styles contains style definitions for this table component.
type Styles struct {
	Header   lipgloss.Style
	Cell     lipgloss.Style
	Selected lipgloss.Style
}

// DefaultStyles returns a set of default style definitions for this table.
func DefaultStyles() Styles {
	return Styles{
		Selected: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212")),
		Header:   lipgloss.NewStyle().Bold(true).Padding(0, 1),
		Cell:     lipgloss.NewStyle().Padding(0, 1),
	}
}

// SetStyles sets the table styles.
func (m *Model) SetStyles(s Styles) {
	m.styles = s
	m.UpdateViewport()
}

// StyleFunc customizes the style of a table cell based on the row, column, and value.
type StyleFunc func(row, col int, value string) lipgloss.Style

// Option is used to set options in New.
type Option func(*Model)

// New creates a new model for the table widget.
func New(opts ...Option) Model {
	m := Model{
		cursor:   0,
		viewport: viewport.New(viewport.WithWidth(0), viewport.WithHeight(20)),
		KeyMap:   DefaultKeyMap(),
		styles:   DefaultStyles(),
	}
	for _, opt := range opts {
		opt(&m)
	}
	m.UpdateViewport()
	return m
}

func WithColumns(cols []Column) Option  { return func(m *Model) { m.cols = cols } }
func WithRows(rows []Row) Option        { return func(m *Model) { m.rows = rows } }
func WithFocused(f bool) Option         { return func(m *Model) { m.focus = f } }
func WithStyles(s Styles) Option        { return func(m *Model) { m.styles = s } }
func WithStyleFunc(f StyleFunc) Option  { return func(m *Model) { m.styleFunc = f } }
func WithKeyMap(km KeyMap) Option       { return func(m *Model) { m.KeyMap = km } }

// WithHeight sets the height of the table.
func WithHeight(h int) Option {
	return func(m *Model) {
		m.viewport.SetHeight(h - lipgloss.Height(m.headersView()))
	}
}

// WithWidth sets the width of the table.
func WithWidth(w int) Option {
	return func(m *Model) { m.viewport.SetWidth(w) }
}

// Update is the Bubble Tea update loop.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.focus {
		return m, nil
	}
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.KeyMap.LineUp):
			m.MoveUp(1)
		case key.Matches(msg, m.KeyMap.LineDown):
			m.MoveDown(1)
		case key.Matches(msg, m.KeyMap.PageUp):
			m.MoveUp(m.viewport.Height())
		case key.Matches(msg, m.KeyMap.PageDown):
			m.MoveDown(m.viewport.Height())
		case key.Matches(msg, m.KeyMap.HalfPageUp):
			m.MoveUp(m.viewport.Height() / 2)
		case key.Matches(msg, m.KeyMap.HalfPageDown):
			m.MoveDown(m.viewport.Height() / 2)
		case key.Matches(msg, m.KeyMap.GotoTop):
			m.GotoTop()
		case key.Matches(msg, m.KeyMap.GotoBottom):
			m.GotoBottom()
		}
	}
	return m, nil
}

func (m Model) Focused() bool   { return m.focus }
func (m *Model) Focus()         { m.focus = true; m.UpdateViewport() }
func (m *Model) Blur()          { m.focus = false; m.UpdateViewport() }
func (m Model) SelectedRow() Row {
	if m.cursor < 0 || m.cursor >= len(m.rows) {
		return nil
	}
	return m.rows[m.cursor]
}
func (m Model) Rows() []Row       { return m.rows }
func (m Model) Columns() []Column { return m.cols }
func (m Model) Height() int       { return m.viewport.Height() }
func (m Model) Width() int        { return m.viewport.Width() }
func (m Model) Cursor() int       { return m.cursor }

func (m *Model) SetRows(r []Row)       { m.rows = r; m.UpdateViewport() }
func (m *Model) SetColumns(c []Column) { m.cols = c; m.UpdateViewport() }
func (m *Model) SetWidth(w int)        { m.viewport.SetWidth(w); m.UpdateViewport() }

func (m *Model) SetHeight(h int) {
	m.viewport.SetHeight(h - lipgloss.Height(m.headersView()))
	m.UpdateViewport()
}

func (m *Model) SetCursor(n int) {
	m.cursor = tuishell.Clamp(n, 0, len(m.rows)-1)
	m.UpdateViewport()
}

// View renders the component.
func (m Model) View() string {
	if len(m.viewport.View()) == 0 {
		return EmptyMsg.Width(m.W).Height(m.H).Render(m.EmptyMessage)
	}
	return m.headersView() + "\n" + m.viewport.View()
}

// UpdateViewport updates the list content based on the previously defined columns and rows.
func (m *Model) UpdateViewport() {
	renderedRows := make([]string, 0, len(m.rows))
	if m.cursor >= 0 {
		m.start = tuishell.Clamp(m.cursor-m.viewport.Height(), 0, m.cursor)
	} else {
		m.start = 0
	}
	m.end = tuishell.Clamp(m.cursor+m.viewport.Height(), m.cursor, len(m.rows))
	for i := m.start; i < m.end; i++ {
		renderedRows = append(renderedRows, m.renderRow(i))
	}
	m.viewport.SetContent(
		lipgloss.JoinVertical(lipgloss.Left, renderedRows...),
	)
}

func (m *Model) MoveUp(n int) {
	m.cursor = tuishell.Clamp(m.cursor-n, 0, len(m.rows)-1)
	switch {
	case m.start == 0:
		m.viewport.SetYOffset(tuishell.Clamp(m.viewport.YOffset(), 0, m.cursor))
	case m.start < m.viewport.Height():
		m.viewport.SetYOffset(tuishell.Clamp(tuishell.Clamp(m.viewport.YOffset()+n, 0, m.cursor), 0, m.viewport.Height()))
	case m.viewport.YOffset() >= 1:
		m.viewport.SetYOffset(tuishell.Clamp(m.viewport.YOffset()+n, 1, m.viewport.Height()))
	}
	m.UpdateViewport()
}

func (m *Model) MoveDown(n int) {
	m.cursor = tuishell.Clamp(m.cursor+n, 0, len(m.rows)-1)
	m.UpdateViewport()
	switch {
	case m.end == len(m.rows) && m.viewport.YOffset() > 0:
		m.viewport.SetYOffset(tuishell.Clamp(m.viewport.YOffset()-n, 1, m.viewport.Height()))
	case m.cursor > (m.end-m.start)/2 && m.viewport.YOffset() > 0:
		m.viewport.SetYOffset(tuishell.Clamp(m.viewport.YOffset()-n, 1, m.cursor))
	case m.viewport.YOffset() > 1:
	case m.cursor > m.viewport.YOffset()+m.viewport.Height()-1:
		m.viewport.SetYOffset(tuishell.Clamp(m.viewport.YOffset()+1, 0, 1))
	}
}

func (m *Model) GotoTop()    { m.MoveUp(m.cursor) }
func (m *Model) GotoBottom() { m.MoveDown(len(m.rows)) }

// FromValues creates table rows from a delimited string.
func (m *Model) FromValues(value, separator string) {
	rows := []Row{}
	for _, line := range strings.Split(value, "\n") {
		r := Row{}
		for _, field := range strings.Split(line, separator) {
			r = append(r, field)
		}
		rows = append(rows, r)
	}
	m.SetRows(rows)
}

func (m Model) headersView() string {
	s := make([]string, 0, len(m.cols))
	for _, col := range m.cols {
		if col.Width <= 0 {
			continue
		}
		st := lipgloss.NewStyle().Width(col.Width).MaxWidth(col.Width).Inline(true)
		if col.Centered {
			st = st.Align(lipgloss.Center, lipgloss.Center)
		}
		renderedCell := st.Render(ansi.Truncate(col.Title, col.Width, "…"))
		s = append(s, m.styles.Header.Render(renderedCell))
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, s...)
}

func (m *Model) renderRow(r int) string {
	s := make([]string, 0, len(m.cols))
	for i, value := range m.rows[r] {
		if m.cols[i].Width <= 0 {
			continue
		}
		var cellStyle lipgloss.Style
		if m.styleFunc != nil {
			cellStyle = m.styleFunc(r, i, value)
			if r == m.cursor {
				cellStyle = m.styles.Selected.Padding(0, 1)
			}
		} else {
			cellStyle = m.styles.Cell
		}

		st := lipgloss.NewStyle().
			Width(m.cols[i].Width).
			MaxWidth(m.cols[i].Width).
			Inline(true).
			AlignVertical(lipgloss.Center)
		if m.cols[i].Centered {
			st = st.AlignHorizontal(lipgloss.Center)
		}
		renderedCell := cellStyle.
			Height(rowHeight).
			Render(st.Render(ansi.Truncate(value, m.cols[i].Width, "…")))
		s = append(s, renderedCell)
	}

	row := lipgloss.JoinHorizontal(lipgloss.Left, s...)
	if r == m.cursor {
		return m.styles.Selected.MarginBottom(rowBottomMargin).Render(row)
	}
	return lipgloss.NewStyle().MarginBottom(rowBottomMargin).Render(row)
}
