package table

import (
	"fmt"
	"math"
	"slices"
	"time"
)

// InitModelParams holds the parameters for creating a pre-configured table.
type InitModelParams struct {
	Rows      []Row
	Colums    []Column
	StyleFunc StyleFunc
	Height    int
	Width     int
	Styles    *Styles
}

// InitModel creates a focused table from the given params.
func InitModel(params InitModelParams) Model {
	opts := []Option{
		WithColumns(params.Colums),
		WithRows(params.Rows),
		WithFocused(true),
		WithWidth(params.Width),
		WithHeight(params.Height),
		WithStyleFunc(params.StyleFunc),
	}
	if params.Styles != nil {
		opts = append(opts, WithStyles(*params.Styles))
	}
	return New(opts...)
}

// ParseTimeString parses an RFC3339 time string.
func ParseTimeString(d string) time.Time {
	t, _ := time.Parse(time.RFC3339, d)
	return t
}

// FormatTime returns a human-readable relative time string.
func FormatTime(t time.Time) string {
	r := time.Since(t.Local())
	days := math.Floor(r.Hours()) / 24
	week := days / 7

	switch {
	case week > 4:
		return fmt.Sprintf("%.0fM", week/4)
	case days > 7:
		return fmt.Sprintf("%.0fw", week)
	case math.Floor(r.Hours()) > 24:
		return fmt.Sprintf("%.0fd", days)
	case math.Floor(r.Hours()) > 0:
		return fmt.Sprintf("%.0fh", r.Hours())
	case math.Floor(r.Minutes()) > 0:
		return fmt.Sprintf("%.0fm", r.Minutes())
	default:
		return fmt.Sprintf("%.0fs", r.Seconds())
	}
}

// FormatPercentage formats a float as a percentage string.
func FormatPercentage(v float32) string {
	if v == 0 {
		return ""
	}
	return fmt.Sprintf("%.2f %%", v)
}

// FormatDuration formats seconds as minutes.
func FormatDuration(d float32) string {
	seconds := d / 60.0
	x := time.Duration(d * float32(time.Second))
	switch {
	case seconds > 0:
		return fmt.Sprintf("%.0f m", x.Minutes())
	case seconds < 0:
		return fmt.Sprintf("%.0f m", x.Minutes())
	default:
		return ""
	}
}

// ColWidth computes a column width as a percentage of total width.
func ColWidth(w int, p int) int {
	return int(float32(w) * float32(p) / 100)
}

// RenderIcon returns the icon string if b is true, empty string otherwise.
func RenderIcon(b bool, icon string) string {
	if b {
		return icon
	}
	return ""
}

// GetColIndex returns the index of a column by name, or -1 if not found.
func GetColIndex(cols []Column, n string) int {
	return slices.IndexFunc(cols, func(c Column) bool {
		return c.Name == n
	})
}
