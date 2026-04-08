// Package tuishell provides a reusable 3-panel TUI layout engine, shared context,
// global keybindings, and utility functions for bubbletea applications.
package tuishell

import "fmt"

// Max returns the larger of a or b.
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Min returns the smaller of a or b.
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Clamp restricts v to the range [low, high].
func Clamp(v, low, high int) int {
	return Min(Max(v, low), high)
}

// Truncate shortens s to limit characters, appending "..." if truncated.
func Truncate(s string, limit int) string {
	if len(s) >= Max(limit, 20) {
		return fmt.Sprintf("%v...", s[:limit])
	}
	return s
}
