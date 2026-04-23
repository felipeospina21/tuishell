package tuishell

import tea "charm.land/bubbletea/v2"

// FocusedPanel identifies which panel currently has focus.
type FocusedPanel uint

// Panel focus constants.
const (
// LeftPanel indicates the left (navigation) panel has focus.
	LeftPanel FocusedPanel = iota
	// MainPanel indicates the main (content) panel has focus.
	MainPanel
	// RightPanel indicates the right (details) panel has focus.
	RightPanel
	// ModalPanel indicates the modal overlay has focus.
	ModalPanel
)

// AppContext holds shared state passed to all tuishell components.
// Embed this in your app-specific context to add domain fields.
type AppContext struct {
	Window       tea.WindowSizeMsg
	DemoMode      bool
	FocusedPanel FocusedPanel
	PanelHeight  int
}
