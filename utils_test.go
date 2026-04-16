package tuishell

import "testing"

func TestMax(t *testing.T) {
	tests := []struct {
		name   string
		a, b   int
		expect int
	}{
		{"a greater", 5, 3, 5},
		{"b greater", 3, 5, 5},
		{"equal", 4, 4, 4},
		{"negative", -1, -5, -1},
		{"zero", 0, 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Max(tt.a, tt.b); got != tt.expect {
				t.Errorf("Max(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.expect)
			}
		})
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		name   string
		a, b   int
		expect int
	}{
		{"a smaller", 3, 5, 3},
		{"b smaller", 5, 3, 3},
		{"equal", 4, 4, 4},
		{"negative", -5, -1, -5},
		{"zero", 0, 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Min(tt.a, tt.b); got != tt.expect {
				t.Errorf("Min(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.expect)
			}
		})
	}
}

func TestClamp(t *testing.T) {
	tests := []struct {
		name          string
		v, low, high  int
		expect        int
	}{
		{"within range", 5, 0, 10, 5},
		{"below low", -1, 0, 10, 0},
		{"above high", 15, 0, 10, 10},
		{"at low", 0, 0, 10, 0},
		{"at high", 10, 0, 10, 10},
		{"zero range", 5, 3, 3, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Clamp(tt.v, tt.low, tt.high); got != tt.expect {
				t.Errorf("Clamp(%d, %d, %d) = %d, want %d", tt.v, tt.low, tt.high, got, tt.expect)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		limit  int
		expect string
	}{
		{"short string ignored", "hello", 3, "hello"},
		{"under 20 chars never truncated", "nineteen chars!!", 5, "nineteen chars!!"},
		{"exactly 20 chars with limit 10", "12345678901234567890", 10, "1234567890..."},
		{"long string truncated", "this is a very long string that exceeds twenty", 10, "this is a ..."},
		{"limit above 20 no truncation needed", "short", 25, "short"},
		{"limit 0 short string", "hi", 0, "hi"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Truncate(tt.s, tt.limit); got != tt.expect {
				t.Errorf("Truncate(%q, %d) = %q, want %q", tt.s, tt.limit, got, tt.expect)
			}
		})
	}
}
