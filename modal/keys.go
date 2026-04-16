package modal

import (
	"charm.land/bubbles/v2/key"
	"github.com/felipeospina21/tuishell"
)

// KeyMap defines the keybindings for the modal overlay.
type KeyMap struct {
	Close  key.Binding
	Submit key.Binding
	Copy   key.Binding
	tuishell.GlobalKeyMap
}

// ShortHelp returns the keybindings for the short help view.
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Close, k.Submit, k.Copy, k.Quit}
}

// FullHelp returns the keybindings for the full help view.
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Close, k.Submit, k.Copy, k.Quit},
	}
}

// Keybinds is the default set of modal keybindings.
var Keybinds = KeyMap{
	Close: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "close modal"),
	),
	Submit: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "submit"),
	),
	Copy: key.NewBinding(
		key.WithKeys("ctrl+y"),
		key.WithHelp("ctrl+y", "copy"),
	),
	GlobalKeyMap: tuishell.GlobalKeys(false),
}
