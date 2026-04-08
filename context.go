package tuishell

import tea "charm.land/bubbletea/v2"

// FocusedPanel identifies which panel currently has focus.
type FocusedPanel uint

// Panel focus constants.
const (
	LeftPanel FocusedPanel = iota
	MainPanel
	RightPanel
	ModalPanel
)

// AppContext holds shared state passed to all tuishell components.
// Embed this in your app-specific context to add domain fields.
type AppContext struct {
	Window       tea.WindowSizeMsg
	DevMode      bool
	FocusedPanel FocusedPanel
	PanelHeight  int
}
