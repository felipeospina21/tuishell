package table

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

// newTestTable creates a table with 10 rows, 1 column, viewport height 6 (visibleRows=3).
func newTestTable() Model {
	rows := make([]Row, 10)
	for i := range rows {
		rows[i] = Row{"row"}
	}
	m := New(
		WithColumns([]Column{{Title: "C", Width: 10}}),
		WithRows(rows),
	)
	m.viewport.SetHeight(6)
	m.UpdateViewport()
	return m
}

func TestScrolling(t *testing.T) {
	tests := []struct {
		name                           string
		action                         func(m *Model)
		wantCursor, wantStart, wantEnd int
	}{
		{
			name:       "initial state",
			action:     func(m *Model) {},
			wantCursor: 0, wantStart: 0, wantEnd: 3,
		},
		{
			name:       "MoveDown(1) stays in window",
			action:     func(m *Model) { m.MoveDown(1) },
			wantCursor: 1, wantStart: 0, wantEnd: 3,
		},
		{
			name:       "MoveDown(1)x3 slides window",
			action:     func(m *Model) { m.MoveDown(1); m.MoveDown(1); m.MoveDown(1) },
			wantCursor: 3, wantStart: 1, wantEnd: 4,
		},
		{
			name: "MoveUp from bottom slides window back",
			action: func(m *Model) {
				m.GotoBottom()
				m.MoveUp(4)
			},
			wantCursor: 5, wantStart: 5, wantEnd: 8,
		},
		{
			name:       "GotoBottom",
			action:     func(m *Model) { m.GotoBottom() },
			wantCursor: 9, wantStart: 7, wantEnd: 10,
		},
		{
			name: "GotoTop after GotoBottom",
			action: func(m *Model) {
				m.GotoBottom()
				m.GotoTop()
			},
			wantCursor: 0, wantStart: 0, wantEnd: 3,
		},
		{
			name:       "SetCursor to middle",
			action:     func(m *Model) { m.SetCursor(5) },
			wantCursor: 5, wantStart: 3, wantEnd: 6,
		},
		{
			name:       "MoveDown large n clamps to last",
			action:     func(m *Model) { m.MoveDown(100) },
			wantCursor: 9, wantStart: 7, wantEnd: 10,
		},
		{
			name: "MoveUp large n clamps to 0",
			action: func(m *Model) {
				m.SetCursor(5)
				m.MoveUp(100)
			},
			wantCursor: 0, wantStart: 0, wantEnd: 3,
		},
		{
			name:       "page down (MoveDown visibleRows)",
			action:     func(m *Model) { m.MoveDown(3) },
			wantCursor: 3, wantStart: 1, wantEnd: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newTestTable()
			tt.action(&m)

			if got := m.Cursor(); got != tt.wantCursor {
				t.Errorf("cursor = %d, want %d", got, tt.wantCursor)
			}
			if m.start != tt.wantStart {
				t.Errorf("start = %d, want %d", m.start, tt.wantStart)
			}
			if m.end != tt.wantEnd {
				t.Errorf("end = %d, want %d", m.end, tt.wantEnd)
			}
		})
	}
}

func TestGeometryAccessors(t *testing.T) {
	m := newTestTable() // visibleRows = 6 / (1+1) = 3

	if got := m.RowPitch(); got != RowHeight+RowBottomMargin {
		t.Errorf("RowPitch() = %d, want %d", got, RowHeight+RowBottomMargin)
	}
	if got := m.Start(); got != 0 {
		t.Errorf("Start() = %d, want 0", got)
	}
	start, end := m.VisibleRange()
	if start != 0 || end != 3 {
		t.Errorf("VisibleRange() = (%d,%d), want (0,3)", start, end)
	}
	if hh := m.HeaderHeight(); hh <= 0 {
		t.Errorf("HeaderHeight() = %d, want > 0", hh)
	}

	m.SetCursor(5) // window slides so start=3, end=6
	if s, e := m.VisibleRange(); s != 3 || e != 6 {
		t.Errorf("after SetCursor(5): VisibleRange() = (%d,%d), want (3,6)", s, e)
	}
}

func TestRowAt(t *testing.T) {
	m := newTestTable()
	hh := m.HeaderHeight()
	pitch := m.RowPitch()

	tests := []struct {
		name    string
		y       int
		wantRow int
		wantOK  bool
	}{
		{"header line", hh - 1, 0, false},
		{"first visible row", hh, 0, true},
		{"first row bottom margin", hh + RowHeight, 0, false},
		{"second visible row", hh + pitch, 1, true},
		{"third visible row", hh + 2*pitch, 2, true},
		{"below last visible row", hh + 3*pitch, 0, false},
		{"negative", -1, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRow, gotOK := m.RowAt(tt.y)
			if gotOK != tt.wantOK || (tt.wantOK && gotRow != tt.wantRow) {
				t.Errorf("RowAt(%d) = (%d,%v), want (%d,%v)", tt.y, gotRow, gotOK, tt.wantRow, tt.wantOK)
			}
		})
	}
}

func TestRowAt_ScrolledWindow(t *testing.T) {
	m := newTestTable()
	m.SetCursor(5) // start=3, end=6
	hh := m.HeaderHeight()
	pitch := m.RowPitch()

	// First visible line now maps to row index 3 (the current start).
	if row, ok := m.RowAt(hh); !ok || row != 3 {
		t.Errorf("RowAt(hh) after scroll = (%d,%v), want (3,true)", row, ok)
	}
	if row, ok := m.RowAt(hh + 2*pitch); !ok || row != 5 {
		t.Errorf("RowAt(hh+2*pitch) after scroll = (%d,%v), want (5,true)", row, ok)
	}
}

func TestMouseWheelScroll(t *testing.T) {
	m := newTestTable()
	m.Focus()

	// Wheel down moves cursor down by MouseWheelDelta.
	m, _ = m.Update(tea.MouseWheelMsg(tea.Mouse{Button: tea.MouseWheelDown}))
	if got := m.Cursor(); got != MouseWheelDelta {
		t.Errorf("after wheel down cursor = %d, want %d", got, MouseWheelDelta)
	}

	// Wheel up moves cursor back up.
	m, _ = m.Update(tea.MouseWheelMsg(tea.Mouse{Button: tea.MouseWheelUp}))
	if got := m.Cursor(); got != 0 {
		t.Errorf("after wheel up cursor = %d, want 0", got)
	}
}

func TestMouseClickSelectsRow(t *testing.T) {
	m := newTestTable()
	m.Focus()
	hh := m.HeaderHeight()
	pitch := m.RowPitch()

	// Click the third visible row.
	m, _ = m.Update(tea.MouseClickMsg(tea.Mouse{Button: tea.MouseLeft, Y: hh + 2*pitch}))
	if got := m.Cursor(); got != 2 {
		t.Errorf("after click on row 2, cursor = %d, want 2", got)
	}

	// Click on header does not change selection.
	m, _ = m.Update(tea.MouseClickMsg(tea.Mouse{Button: tea.MouseLeft, Y: 0}))
	if got := m.Cursor(); got != 2 {
		t.Errorf("after click on header, cursor = %d, want 2 (unchanged)", got)
	}

	// Non-left button is ignored.
	m, _ = m.Update(tea.MouseClickMsg(tea.Mouse{Button: tea.MouseRight, Y: hh}))
	if got := m.Cursor(); got != 2 {
		t.Errorf("after right click, cursor = %d, want 2 (unchanged)", got)
	}
}

func TestMouseReleaseSelectsRow(t *testing.T) {
	m := newTestTable()
	m.Focus()
	hh := m.HeaderHeight()

	m, _ = m.Update(tea.MouseReleaseMsg(tea.Mouse{Button: tea.MouseLeft, Y: hh}))
	if got := m.Cursor(); got != 0 {
		t.Errorf("after release on first row, cursor = %d, want 0", got)
	}
}

func TestMouseGatedByFocus(t *testing.T) {
	// Unfocused table ignores mouse by default.
	m := newTestTable() // not focused
	m, _ = m.Update(tea.MouseWheelMsg(tea.Mouse{Button: tea.MouseWheelDown}))
	if got := m.Cursor(); got != 0 {
		t.Errorf("unfocused table responded to wheel: cursor = %d, want 0", got)
	}

	// With AllowMouseWhenBlurred the unfocused table responds.
	m2 := newTestTable()
	m2.AllowMouseWhenBlurred = true
	m2, _ = m2.Update(tea.MouseWheelMsg(tea.Mouse{Button: tea.MouseWheelDown}))
	if got := m2.Cursor(); got != MouseWheelDelta {
		t.Errorf("blurred-allowed table cursor = %d, want %d", got, MouseWheelDelta)
	}
}

// fakeZone is a minimal ZoneManager for testing zone marking.
type fakeZone struct{ marked map[string]string }

func (f *fakeZone) Mark(id, v string) string {
	if f.marked == nil {
		f.marked = map[string]string{}
	}
	f.marked[id] = v
	return "[" + id + "]" + v + "[/" + id + "]"
}

func TestZoneMarkingWrapsRows(t *testing.T) {
	rows := []Row{{"a"}, {"b"}, {"c"}}
	fz := &fakeZone{}
	m := New(
		WithColumns([]Column{{Title: "C", Width: 10}}),
		WithRows(rows),
		WithZoneManager(fz),
		WithRowID(func(i int) string { return "row-" + string(rune('0'+i)) }),
	)
	m.viewport.SetHeight(10)
	m.UpdateViewport()

	if len(fz.marked) == 0 {
		t.Fatal("expected rows to be marked with zone ids, got none")
	}
	if _, ok := fz.marked["row-0"]; !ok {
		t.Errorf("expected zone id row-0 to be marked, marked=%v", fz.marked)
	}
}

func TestZoneMarkingDisabledWithoutRowID(t *testing.T) {
	fz := &fakeZone{}
	m := New(
		WithColumns([]Column{{Title: "C", Width: 10}}),
		WithRows([]Row{{"a"}}),
		WithZoneManager(fz),
		// no RowID hook
	)
	m.viewport.SetHeight(10)
	m.UpdateViewport()
	if len(fz.marked) != 0 {
		t.Errorf("expected no marking without RowID hook, got %v", fz.marked)
	}
}
