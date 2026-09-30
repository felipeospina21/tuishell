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

// Row layout constants. Each rendered row occupies RowHeight content lines
// plus RowBottomMargin blank lines below it. Consumers doing geometric
// hit-testing should use RowPitch() rather than hardcoding these values.
const (
	// RowBottomMargin is the number of blank lines rendered below each row.
	RowBottomMargin = 1
	// RowHeight is the number of content lines occupied by a single row.
	RowHeight = 1

	// Deprecated: unexported aliases kept for internal readability.
	rowBottomMargin = RowBottomMargin
	rowHeight       = RowHeight
)

// ZoneManager is the minimal subset of a bubblezone *zone.Manager that the
// table needs to mark rows as clickable. It is satisfied by
// *github.com/lrstanley/bubblezone/v2.Manager without tuishell taking a hard
// dependency on bubblezone. Mark wraps v in zero-width markers identified by id.
type ZoneManager interface {
	Mark(id, v string) string
}

// Model defines a state for the table widget.
type Model struct {
	KeyMap       KeyMap
	EmptyMessage string
	W            int
	H            int

	// AllowMouseWhenBlurred, when true, lets the table respond to mouse
	// wheel and click events even while it does not have focus. Key events
	// are always gated on focus. Defaults to false.
	AllowMouseWhenBlurred bool

	// RowID, when non-nil together with a ZoneManager set via SetZoneManager,
	// returns a unique zone id for the row at the given index. Rows are then
	// wrapped in zone markers so consumers can hit-test them with bubblezone.
	RowID func(index int) string

	cols      []Column
	rows      []Row
	cursor    int
	focus     bool
	styles    Styles
	styleFunc StyleFunc
	zone      ZoneManager

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
			key.WithKeys("pgdown", spacebar),
			key.WithHelp("pgdn", "page down"),
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

// WithColumns sets the table columns.
func WithColumns(cols []Column) Option { return func(m *Model) { m.cols = cols } }

// WithRows sets the table rows.
func WithRows(rows []Row) Option { return func(m *Model) { m.rows = rows } }

// WithFocused sets the initial focus state.
func WithFocused(f bool) Option { return func(m *Model) { m.focus = f } }

// WithStyles sets the table styles.
func WithStyles(s Styles) Option { return func(m *Model) { m.styles = s } }

// WithStyleFunc sets a custom style function for cell rendering.
func WithStyleFunc(f StyleFunc) Option { return func(m *Model) { m.styleFunc = f } }

// WithKeyMap sets the table keybindings.
func WithKeyMap(km KeyMap) Option { return func(m *Model) { m.KeyMap = km } }

// WithZoneManager sets a bubblezone manager used to mark rows as clickable.
// Use together with WithRowID.
func WithZoneManager(z ZoneManager) Option { return func(m *Model) { m.zone = z } }

// WithRowID sets a hook returning a unique zone id for the row at index.
// When set along with a ZoneManager, each row is wrapped in a zone marker.
func WithRowID(fn func(index int) string) Option { return func(m *Model) { m.RowID = fn } }

// WithAllowMouseWhenBlurred lets the table respond to mouse wheel and click
// events even while unfocused. Key events remain gated on focus.
func WithAllowMouseWhenBlurred(b bool) Option {
	return func(m *Model) { m.AllowMouseWhenBlurred = b }
}

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

// MouseWheelDelta is the number of rows the cursor moves per mouse wheel tick.
const MouseWheelDelta = 1

// Update is the Bubble Tea update loop.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if !m.focus {
			return m, nil
		}
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

	case tea.MouseWheelMsg:
		if !m.mouseEnabled() {
			return m, nil
		}
		switch msg.Button {
		case tea.MouseWheelUp:
			m.MoveUp(MouseWheelDelta)
		case tea.MouseWheelDown:
			m.MoveDown(MouseWheelDelta)
		}

	case tea.MouseClickMsg:
		if !m.mouseEnabled() || msg.Button != tea.MouseLeft {
			return m, nil
		}
		if row, ok := m.RowAt(msg.Y); ok {
			m.SetCursor(row)
		}

	case tea.MouseReleaseMsg:
		if !m.mouseEnabled() || msg.Button != tea.MouseLeft {
			return m, nil
		}
		if row, ok := m.RowAt(msg.Y); ok {
			m.SetCursor(row)
		}
	}
	return m, nil
}

// mouseEnabled reports whether the table should process mouse events given its
// current focus state and the AllowMouseWhenBlurred option.
func (m Model) mouseEnabled() bool { return m.focus || m.AllowMouseWhenBlurred }

// Focused reports whether the table has focus.
func (m Model) Focused() bool { return m.focus }

// Focus sets the table to focused state and updates the viewport.
func (m *Model) Focus() { m.focus = true; m.UpdateViewport() }

// Blur removes focus from the table and updates the viewport.
func (m *Model) Blur() { m.focus = false; m.UpdateViewport() }

// SelectedRow returns the currently selected row, or nil if no rows exist.
func (m Model) SelectedRow() Row {
	if m.cursor < 0 || m.cursor >= len(m.rows) {
		return nil
	}
	return m.rows[m.cursor]
}

// Rows returns all table rows.
func (m Model) Rows() []Row { return m.rows }

// Columns returns all table columns.
func (m Model) Columns() []Column { return m.cols }

// Height returns the viewport height.
func (m Model) Height() int { return m.viewport.Height() }

// Width returns the viewport width.
func (m Model) Width() int { return m.viewport.Width() }

// Cursor returns the index of the currently selected row.
func (m Model) Cursor() int { return m.cursor }

// Start returns the index of the first visible row (the current scroll offset).
func (m Model) Start() int { return m.start }

// VisibleRange returns the [start, end) range of currently visible row indices.
// end is exclusive.
func (m Model) VisibleRange() (start, end int) { return m.start, m.end }

// RowPitch returns the number of terminal lines each row occupies, including
// its bottom margin (RowHeight + RowBottomMargin). Use this instead of
// hardcoding the value when doing geometric hit-testing.
func (m Model) RowPitch() int { return RowHeight + RowBottomMargin }

// HeaderHeight returns the number of terminal lines occupied by the header
// area above the first row within View(), including the separator line.
// A click at a Y coordinate (relative to the table View origin) less than
// HeaderHeight() falls on the header, not a row.
func (m Model) HeaderHeight() int {
	if len(m.rows) == 0 {
		return 0
	}
	// View() renders headersView() + "\n" + viewport, so the separator adds one line.
	return lipgloss.Height(m.headersView()) + 1
}

// RowAt maps a Y coordinate, relative to the top-left origin of this table's
// View() output, to a row index. It returns (index, true) when the coordinate
// falls on a currently visible row, or (0, false) otherwise (e.g. the header,
// a between-row margin, or empty space below the last row).
func (m Model) RowAt(y int) (int, bool) {
	if len(m.rows) == 0 {
		return 0, false
	}
	rel := y - m.HeaderHeight()
	if rel < 0 {
		return 0, false
	}
	pitch := m.RowPitch()
	// Reject clicks that land on the bottom margin between rows.
	if rel%pitch >= RowHeight {
		return 0, false
	}
	idx := m.start + rel/pitch
	if idx < m.start || idx >= m.end {
		return 0, false
	}
	return idx, true
}

// SetZoneManager sets the bubblezone manager used to mark rows as clickable.
// Combined with RowID (see WithRowID), each rendered row is wrapped in a zone
// marker. Passing nil disables zone marking.
func (m *Model) SetZoneManager(z ZoneManager) { m.zone = z; m.UpdateViewport() }

// SetRows replaces the table rows and updates the viewport.
func (m *Model) SetRows(r []Row) { m.rows = r; m.UpdateViewport() }

// SetColumns replaces the table columns and updates the viewport.
func (m *Model) SetColumns(c []Column) { m.cols = c; m.UpdateViewport() }

// SetWidth sets the viewport width and updates the viewport.
func (m *Model) SetWidth(w int) { m.viewport.SetWidth(w); m.UpdateViewport() }

// SetHeight sets the viewport height (minus header) and updates the viewport.
func (m *Model) SetHeight(h int) {
	m.viewport.SetHeight(h - lipgloss.Height(m.headersView()))
	m.UpdateViewport()
}

// SetCursor moves the cursor to row n, clamped to valid bounds.
func (m *Model) SetCursor(n int) {
	m.cursor = tuishell.Clamp(n, 0, len(m.rows)-1)
	m.UpdateViewport()
}

// View renders the component.
func (m Model) View() string {
	if len(m.rows) == 0 {
		w, h := m.W, m.H
		if w == 0 {
			w = m.viewport.Width()
		}
		if h == 0 {
			h = m.viewport.Height()
		}
		return EmptyMsg.Width(w).Height(h).Render(m.EmptyMessage)
	}
	return m.headersView() + "\n" + m.viewport.View()
}

// UpdateViewport updates the list content based on the previously defined columns and rows.
func (m *Model) UpdateViewport() {
	visibleRows := m.viewport.Height() / (rowHeight + rowBottomMargin)
	if visibleRows < 1 {
		visibleRows = 1
	}

	// Keep cursor visible within the window
	if m.cursor < m.start {
		m.start = m.cursor
	} else if m.cursor >= m.start+visibleRows {
		m.start = m.cursor - visibleRows + 1
	}
	m.start = tuishell.Clamp(m.start, 0, max(0, len(m.rows)-visibleRows))
	m.end = tuishell.Clamp(m.start+visibleRows, 0, len(m.rows))

	renderedRows := make([]string, 0, m.end-m.start)
	for i := m.start; i < m.end; i++ {
		renderedRows = append(renderedRows, m.renderRow(i))
	}
	m.viewport.SetContent(
		lipgloss.JoinVertical(lipgloss.Left, renderedRows...),
	)
}

// MoveUp moves the cursor up by n rows.
func (m *Model) MoveUp(n int) {
	m.cursor = tuishell.Clamp(m.cursor-n, 0, len(m.rows)-1)
	m.UpdateViewport()
}

// MoveDown moves the cursor down by n rows.
func (m *Model) MoveDown(n int) {
	m.cursor = tuishell.Clamp(m.cursor+n, 0, len(m.rows)-1)
	m.UpdateViewport()
}

// GotoTop moves the cursor to the first row.
func (m *Model) GotoTop() { m.MoveUp(m.cursor) }

// GotoBottom moves the cursor to the last row.
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

	var rendered string
	if r == m.cursor {
		rendered = m.styles.Selected.MarginBottom(rowBottomMargin).Render(row)
	} else {
		rendered = lipgloss.NewStyle().MarginBottom(rowBottomMargin).Render(row)
	}

	// Wrap the fully-rendered outer row in a zone marker so bubblezone can
	// hit-test it. Marking the outer string (rather than inner cells) keeps
	// the markers away from the per-cell MaxWidth hard-trims that would
	// otherwise corrupt them.
	if m.zone != nil && m.RowID != nil {
		if id := m.RowID(r); id != "" {
			rendered = m.zone.Mark(id, rendered)
		}
	}
	return rendered
}
