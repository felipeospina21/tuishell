package modal

import (
	"charm.land/bubbles/v2/key"
	"github.com/felipeospina21/tuishell"
)

type KeyMap struct {
	Close  key.Binding
	Submit key.Binding
	Copy   key.Binding
	tuishell.GlobalKeyMap
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Close, k.Submit, k.Copy, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Close, k.Submit, k.Copy, k.Quit},
	}
}

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
