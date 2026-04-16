package tuishell

import "charm.land/bubbles/v2/list"

// ConfigureList disables list keybindings that conflict with the shell's
// global keys (quit, force-quit). Call this after creating a list.Model
// that will be used inside a shell panel.
func ConfigureList(l *list.Model) {
	l.DisableQuitKeybindings()
}
