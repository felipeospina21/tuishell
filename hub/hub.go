package hub

import (
	"fmt"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/felipeospina21/tuishell/style"
)

type appID int

const (
	pickerView appID = iota
	appView
)

// AppEntry defines a launchable app.
type AppEntry struct {
	Name        string
	Desc        string
	NewModel    func() tea.Model
}

func (a AppEntry) Title() string       { return a.Name }
func (a AppEntry) Description() string { return a.Desc }
func (a AppEntry) FilterValue() string { return a.Name }

// Model is the hub launcher that switches between a picker and child apps.
type Model struct {
	active   appID
	picker   list.Model
	apps     []AppEntry
	child    tea.Model
	theme    style.Theme
	winSize  tea.WindowSizeMsg
}

var backKey = key.NewBinding(
	key.WithKeys("ctrl+h"),
	key.WithHelp("ctrl+h", "back to launcher"),
)

// New creates a hub model from the given app entries.
func New(apps []AppEntry) Model {
	theme := style.DefaultTheme()

	items := make([]list.Item, len(apps))
	for i, a := range apps {
		items[i] = a
	}

	delegate := list.NewDefaultDelegate()
	l := list.New(items, delegate, 50, 14)
	l.Title = "  tuishell"
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(theme.PrimaryBright).
		Bold(true).
		MarginBottom(1)

	return Model{
		active: pickerView,
		picker: l,
		apps:   apps,
		theme:  theme,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.winSize = msg
		m.picker.SetSize(msg.Width-4, msg.Height-4)
		if m.child != nil {
			var cmd tea.Cmd
			m.child, cmd = m.child.Update(msg)
			return m, cmd
		}
		return m, nil

	case tea.KeyPressMsg:
		if m.active == appView && key.Matches(msg, backKey) {
			m.active = pickerView
			return m, nil
		}

		if m.active == pickerView {
			switch msg.String() {
			case "enter":
				idx := m.picker.Index()
				if idx >= 0 && idx < len(m.apps) {
					m.child = m.apps[idx].NewModel()
					m.active = appView
					cmds := []tea.Cmd{m.child.Init()}
					if m.winSize.Width > 0 {
						cmds = append(cmds, func() tea.Msg { return m.winSize })
					}
					return m, tea.Batch(cmds...)
				}
			case "q", "ctrl+c":
				return m, tea.Quit
			}

			var cmd tea.Cmd
			m.picker, cmd = m.picker.Update(msg)
			return m, cmd
		}
	}

	if m.active == appView && m.child != nil {
		var cmd tea.Cmd
		m.child, cmd = m.child.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) View() tea.View {
	if m.active == appView && m.child != nil {
		return m.child.View()
	}

	t := m.theme
	hint := lipgloss.NewStyle().
		Foreground(t.Muted).
		MarginTop(1).
		Render(fmt.Sprintf("enter select · q quit · %s back", backKey.Help().Key))

	content := m.picker.View() + "\n" + hint

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}
