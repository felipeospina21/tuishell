// Package shell provides ShellModel, a reusable app shell that handles
// panel routing, modal lifecycle, task management, layout, and statusline.
package shell

import (
	"image/color"

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
	RightPanel      tea.Model // optional
	AppIcon         string    // e.g. "🎫" - shown in statusline, combined with selected item
	Keybinds        help.KeyMap
	DemoMode        bool
	LeftPanelWidth  int // default 30
	LeftPanelStyle  lipgloss.Style
	RightPanelStyle lipgloss.Style
	MainFrameStyle  lipgloss.Style

	// EnableMouse makes RenderView() enable tea.MouseModeCellMotion in
	// addition to AltScreen, and makes Update() route mouse messages to the
	// panel under the pointer (with click-to-focus). Defaults to false.
	EnableMouse bool

	// FixedPanels lays out the left, main, and right panels as three
	// always-visible columns. The panels start open and cannot be closed:
	// CloseLeftPanelMsg, CloseRightPanelMsg, ToggleFullscreenMsg, and the
	// global toggle/close keybinds become no-ops. Independent of EnableMouse
	// (a fixed layout is useful for keyboard-only apps too). Defaults to false.
	FixedPanels bool

	// FocusedBorderColor is the border color applied to the currently focused
	// panel as an active-panel indicator. When nil, the shell falls back to
	// Theme.Primary. Works whether or not mouse support is enabled.
	FocusedBorderColor color.Color

	// LeftButtonLabel and RightButtonLabel set the text shown inside the
	// mouse-mode toggle buttons (e.g. "nav", "details"). The shell wraps the
	// label with brackets and a directional arrow indicating open/close. When
	// empty they default to "nav" and "details" respectively.
	LeftButtonLabel  string
	RightButtonLabel string
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

	isLeftOpen         bool
	isRightOpen        bool
	isRightFullscreen  bool
	isModalOpen        bool
	enableMouse        bool
	fixedPanels        bool
	focusedBorderColor color.Color
	leftButtonLabel    string
	rightButtonLabel   string
	taskStatus         taskStatus
	taskErr            error
	prevFocus          tuishell.FocusedPanel
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
		cfg.MainFrameStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(t.Border)
	}

	ctx := tuishell.AppContext{FocusedPanel: tuishell.LeftPanel, DemoMode: cfg.DemoMode}
	sl := statusline.New(t, cfg.DemoMode, cfg.Keybinds)
	sl.ProjectLabel = cfg.AppIcon

	focusColor := cfg.FocusedBorderColor
	if focusColor == nil {
		focusColor = t.Primary
	}

	leftLabel := cfg.LeftButtonLabel
	if leftLabel == "" {
		leftLabel = "nav"
	}
	rightLabel := cfg.RightButtonLabel
	if rightLabel == "" {
		rightLabel = "details"
	}

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
		// Fixed layout starts with the right panel open too and keeps both
		// panels visible for the lifetime of the shell.
		isRightOpen:        cfg.FixedPanels,
		taskStatus:         taskIdle,
		enableMouse:        cfg.EnableMouse,
		fixedPanels:        cfg.FixedPanels,
		focusedBorderColor: focusColor,
		leftButtonLabel:    leftLabel,
		rightButtonLabel:   rightLabel,
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
