package tuishell

import (
	"charm.land/bubbles/v2/help"
	tea "charm.land/bubbletea/v2"
)

// Shell-level messages that panels return to control the shell.

// StartTaskMsg tells the shell to enter loading state and run Cmd.
type StartTaskMsg struct {
	Cmd tea.Cmd
}

// FinishTaskMsg tells the shell the async task completed.
type FinishTaskMsg struct {
	Err      error
	Keybinds help.KeyMap
}

// OpenModalMsg tells the shell to open the modal overlay.
type OpenModalMsg struct {
	Header  string
	Content string
	IsError bool
}

// OpenRightPanelMsg tells the shell to show the right panel.
type OpenRightPanelMsg struct{}

// CloseRightPanelMsg tells the shell to hide the right panel.
type CloseRightPanelMsg struct{}

// CloseLeftPanelMsg tells the shell to hide the left panel.
type CloseLeftPanelMsg struct{}

// ToggleFullscreenMsg tells the shell to toggle right panel fullscreen.
type ToggleFullscreenMsg struct{}

// SetKeybindsMsg tells the shell to update the statusline help keybinds.
type SetKeybindsMsg struct {
	Keybinds help.KeyMap
}

// SetStatusMsg tells the shell to update the statusline content.
type SetStatusMsg struct {
	Mode    string
	Content string
}

// ShellSubmitMsg is forwarded to the app when the modal submit key is pressed.
type ShellSubmitMsg struct{}

// Modal message types — defined here to avoid import cycles between
// tuishell and tuishell/modal. The modal package re-exports these.

// CloseModalMsg signals that the modal should be closed.
type CloseModalMsg struct{}

// SubmitModalMsg signals that the modal content should be submitted.
type SubmitModalMsg struct{}

// CopyModalMsg signals that the modal content should be copied.
type CopyModalMsg struct{}

// ResetHighlightMsg signals that the modal highlight should be cleared.
type ResetHighlightMsg struct{}

// SelectionProvider is implemented by panels that can provide a selected item label.
type SelectionProvider interface {
	SelectedLabel() string
}
