// Package statusline implements the bottom status bar component.
package statusline

import (
	"image/color"

	"charm.land/bubbles/v2/help"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/felipeospina21/tuishell"
	"github.com/felipeospina21/tuishell/style"
)

// Modes defines the possible status bar mode labels.
type Modes struct {
	Normal  string
	Loading string
	Error   string
	Dev     string
}

// ModesEnum contains the available status bar mode values.
var ModesEnum = Modes{
	Normal:  "NORMAL",
	Loading: "LOADING",
	Error:   "ERROR",
	Dev:     "DEVELOP",
}

// Model holds the state for the status bar.
type Model struct {
	Status       string
	Content      string
	Width        int
	ProjectLabel string
	SpinnerView  string
	Help         help.Model
	Keybinds     help.KeyMap
	devMode      bool
	theme        style.Theme
}

// New creates a new status bar model.
func New(theme style.Theme, devMode bool, keybinds help.KeyMap) Model {
	status := ModesEnum.Normal
	if devMode {
		status = ModesEnum.Dev
	}
	return Model{
		Status:   status,
		Keybinds: keybinds,
		Help:     help.New(),
		devMode:  devMode,
		theme:    theme,
	}
}

// Init returns nil; the statusline has no initial command.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles window resize messages to adjust the statusline width.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m.Width = msg.Width
	}
	return m, nil
}

// View renders the status bar with mode indicator, spinner, help, and project label.
func (m Model) View() string {
	t := m.theme
	width := m.Width
	w := lipgloss.Width

	modeColor := m.modeBackground()
	statusKey := statusStyle(t).Background(modeColor).Render(m.Status)
	
	// Show spinner next to mode label when loading
	spinnerView := ""
	if m.Status == ModesEnum.Loading && m.SpinnerView != "" {
		spinnerView = " " + m.SpinnerView
	}
	
	statusVal := statusText(t).Render(tuishell.Truncate(m.Content, width/4))
	encoding := encodingStyle(t).Render("UTF-8")
	projectName := projectStyle(t).Render(m.ProjectLabel)

	helpWidth := width - w(statusKey) - w(spinnerView) - w(statusVal) - w(encoding) - w(projectName) - 2
	if helpWidth < 0 {
		helpWidth = 0
	}
	m.Help.SetWidth(helpWidth)
	helpView := helpText().
		Width(helpWidth + 2).
		Render(" " + m.Help.View(m.Keybinds) + " ")

	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		statusKey,
		spinnerView,
		statusVal,
		helpView,
		encoding,
		projectName,
	)

	return StatusBarStyle().Render(bar)
}

func (m Model) modeBackground() color.Color {
	t := m.theme
	switch m.Status {
	case ModesEnum.Loading:
		return t.StatusLoading
	case ModesEnum.Error:
		return t.StatusError
	case ModesEnum.Dev:
		return t.StatusDev
	default:
		return t.StatusNormal
	}
}

// GetFrameSize returns the total frame size of the status bar.
func GetFrameSize() (int, int) {
	return StatusBarStyle().GetFrameSize()
}

// Theme returns the theme used by this statusline.
func (m Model) Theme() style.Theme {
	return m.theme
}
