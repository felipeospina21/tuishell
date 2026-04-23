# Shell Configuration

The `shell.Model` is the core of tuishell. It manages layout, panel routing, modals, statusline, and task state.

## Config Options

```go
type Config struct {
    Theme           style.Theme    // Required: color theme
    LeftPanel       tea.Model      // Required: navigation panel
    MainPanel       tea.Model      // Required: content panel
    RightPanel      tea.Model      // Optional: details panel
    AppIcon         string         // Optional: icon shown in statusline (e.g. "🎫")
    Keybinds        help.KeyMap    // Required: keybindings for help display
    DemoMode         bool           // Optional: enables demo-mode keys
    LeftPanelWidth  int            // Optional: default 30
    LeftPanelStyle  lipgloss.Style // Required: style for left panel border/padding
    RightPanelStyle lipgloss.Style // Optional: style for right panel
    MainFrameStyle  lipgloss.Style // Optional: outer border style
}
```

## Defaults

| Option | Default |
|--------|---------|
| `LeftPanelWidth` | 30 |
| `MainFrameStyle` | Normal border with `theme.Border` color |
| `RightPanel` | `nil` (no right panel) |
| `AppIcon` | `""` (empty statusline label) |
| `DemoMode` | `false` |

## Minimal Example

```go
s := shell.New(shell.Config{
    Theme:          myTheme,
    LeftPanel:      &myLeftPanel{},
    MainPanel:      &myMainPanel{},
    Keybinds:       tuishell.GlobalKeys(false),
    LeftPanelStyle: lipgloss.NewStyle().Width(30),
})
```

## Full Example

```go
theme := myTheme()

leftStyle := lipgloss.NewStyle().
    PaddingRight(2).
    Border(lipgloss.NormalBorder(), false, true, false, false).
    BorderForeground(theme.Border).
    Width(30)

rightStyle := lipgloss.NewStyle().
    PaddingLeft(2).
    Border(lipgloss.NormalBorder(), false, false, false, true).
    BorderForeground(theme.Border)

s := shell.New(shell.Config{
    Theme:           theme,
    LeftPanel:       &projectsPanel{},
    MainPanel:       &issuesPanel{},
    RightPanel:      &detailsPanel{},
    AppIcon:         "🦊",
    Keybinds:        tuishell.GlobalKeys(cfg.DemoMode),
    DemoMode:         cfg.DemoMode,
    LeftPanelWidth:  30,
    LeftPanelStyle:  leftStyle,
    RightPanelStyle: rightStyle,
})
```

## Accessing Shell State

The shell exposes several fields for reading state:

```go
m.shell.Ctx.FocusedPanel  // Current focus: LeftPanel, MainPanel, or RightPanel
m.shell.Ctx.Window        // Current window size
m.shell.Ctx.PanelHeight   // Computed content height
m.shell.Layout            // Computed layout dimensions
m.shell.Statusline        // Statusline model (can modify ProjectLabel, Content)
m.shell.Spinner           // Spinner model (use .View() for loading states)
```

## Panel Requirements

Panels must implement `tea.Model`:

```go
type tea.Model interface {
    Init() tea.Cmd
    Update(tea.Msg) (tea.Model, tea.Cmd)
    View() tea.View
}
```

Panels receive `tea.WindowSizeMsg` when their dimensions change. Use this to resize internal components.

## Wrapping Panels

If your panel's `Update` returns a concrete type, wrap it to preserve the type for interface assertions:

```go
type ProjectsPanel struct {
    *projects.Model
}

func (p ProjectsPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    m, cmd := p.Model.Update(msg)
    p.Model = &m
    return p, cmd  // Returns ProjectsPanel, not projects.Model
}

func (p ProjectsPanel) View() tea.View { return p.Model.View() }

// Implement SelectionProvider for automatic statusline updates
func (p ProjectsPanel) SelectedLabel() string {
    if item := p.Model.List.SelectedItem(); item != nil {
        return item.FilterValue()
    }
    return ""
}
```
