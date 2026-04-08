// Package statusline implements the bottom status bar component.
package statusline

import (
	"image/color"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/spinner"
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
	Spinner      spinner.Model
	Help         help.Model
	Keybinds     help.KeyMap
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
		Spinner: spinner.New(
			spinner.WithSpinner(spinner.Dot),
			spinner.WithStyle(SpinnerStyle(theme)),
		),
		Help:  help.New(),
		theme: theme,
	}
}

func (m Model) Init() tea.Cmd {
	return m.Spinner.Tick
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case spinner.TickMsg:
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		return m, nil
	default:
		return m, cmd
	}
}

func (m Model) View() string {
	t := m.theme
	width := m.Width
	w := lipgloss.Width

	modeColor := m.modeBackground()
	statusKey := statusStyle(t).Background(modeColor).Render(m.Status)
	statusVal := statusText(t).Render(tuishell.Truncate(m.Content, width/4))
	encoding := encodingStyle(t).Render("UTF-8")
	projectName := projectStyle(t).Render(m.ProjectLabel)

	helpWidth := width - w(statusKey) - w(statusVal) - w(encoding) - w(projectName) - 2
	if helpWidth < 0 {
		helpWidth = 0
	}
	m.Help.SetWidth(helpWidth)
	helpView := helpText().
		Width(helpWidth + 2).
		Render(" " + m.Help.View(m.Keybinds) + " ")

	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		statusKey,
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
