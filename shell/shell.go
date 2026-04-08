// Package shell provides ShellModel, a reusable app shell that handles
// panel routing, modal lifecycle, task management, layout, and statusline.
package shell

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/felipeospina21/tuishell"
	"github.com/felipeospina21/tuishell/modal"
	"github.com/felipeospina21/tuishell/statusline"
	"github.com/felipeospina21/tuishell/style"
)

// Config configures a new Model.
type Config struct {
	Theme           style.Theme
	LeftPanel       tea.Model
	MainPanel       tea.Model
	RightPanel      tea.Model      // optional
	AppIcon         string         // e.g. "🎫" - shown in statusline, combined with selected item
	Keybinds        help.KeyMap
	DevMode         bool
	LeftPanelWidth  int            // default 30
	LeftPanelStyle  lipgloss.Style
	RightPanelStyle lipgloss.Style
	MainFrameStyle  lipgloss.Style
}

// Model handles all common 3-panel TUI behavior.
type Model struct {
	Left  tea.Model
	Main  tea.Model
	Right tea.Model

	Modal      modal.Model
	Statusline statusline.Model
	Spinner    spinner.Model

	Layout tuishell.Layout
	Ctx    tuishell.AppContext

	theme           style.Theme
	leftPanelStyle  lipgloss.Style
	rightPanelStyle lipgloss.Style
	mainFrameStyle  lipgloss.Style
	leftPanelWidth  int
	appIcon         string

	isLeftOpen        bool
	isRightOpen       bool
	isRightFullscreen bool
	isModalOpen       bool
	taskStatus        taskStatus
	taskErr           error
	prevFocus         tuishell.FocusedPanel
}

type taskStatus uint

const (
	taskIdle taskStatus = iota
	taskStarted
	taskFinished
)

// New creates a Model from the given config.
func New(cfg Config) Model {
	t := cfg.Theme
	if cfg.LeftPanelWidth == 0 {
		cfg.LeftPanelWidth = 30
	}
	if cfg.MainFrameStyle.GetWidth() == 0 {
		cfg.MainFrameStyle = style.MainFrameStyle(t)
	}

	ctx := tuishell.AppContext{FocusedPanel: tuishell.LeftPanel, DevMode: cfg.DevMode}
	sl := statusline.New(t, cfg.DevMode, cfg.Keybinds)
	sl.ProjectLabel = cfg.AppIcon

	return Model{
		Left:  cfg.LeftPanel,
		Main:  cfg.MainPanel,
		Right: cfg.RightPanel,
		Modal: modal.New(&ctx, t),

		Statusline: sl,
		Spinner: spinner.New(
			spinner.WithSpinner(spinner.Dot),
			spinner.WithStyle(statusline.SpinnerStyle(t)),
		),

		Ctx:             ctx,
		theme:           t,
		appIcon:         cfg.AppIcon,
		leftPanelStyle:  cfg.LeftPanelStyle,
		rightPanelStyle: cfg.RightPanelStyle,
		mainFrameStyle:  cfg.MainFrameStyle,
		leftPanelWidth:  cfg.LeftPanelWidth,
		isLeftOpen:      true,
		taskStatus:      taskIdle,
	}
}

// Init returns the initial commands for the shell and its children.
func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{m.Statusline.Init(), m.Spinner.Tick}
	if m.Left != nil {
		cmds = append(cmds, m.Left.Init())
	}
	if m.Main != nil {
		cmds = append(cmds, m.Main.Init())
	}
	if m.Right != nil {
		cmds = append(cmds, m.Right.Init())
	}
	return tea.Batch(cmds...)
}
