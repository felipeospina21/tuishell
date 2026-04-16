package table

import "testing"

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
		name                       string
		action                     func(m *Model)
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
