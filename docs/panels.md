# Panels

Panels are the building blocks of a tuishell app. Each panel is a `tea.Model` that receives messages and renders its view.

## Panel Interface

```go
type tea.Model interface {
    Init() tea.Cmd
    Update(tea.Msg) (tea.Model, tea.Cmd)
    View() tea.View
}
```

## Window Size Messages

The shell sends `tea.WindowSizeMsg` to panels when their dimensions change. Use this to resize internal components:

```go
func (m *myPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        m.list.SetSize(msg.Width, msg.Height)
        return m, nil
    }
    // ...
}
```

## SelectionProvider Interface

Panels that provide a selected item label can implement `SelectionProvider`:

```go
type SelectionProvider interface {
    SelectedLabel() string
}
```

When the left panel closes (via `CloseLeftPanelMsg`), the shell checks if it implements this interface. If so, it updates the statusline label to:

```
AppIcon + " " + SelectedLabel()
```

Example: If `AppIcon` is `"🦊"` and `SelectedLabel()` returns `"my-project"`, the statusline shows `"🦊 my-project"`.

## Panel Wrappers

When your panel's `Update` method returns a concrete type, wrap it to preserve the type for interface assertions:

```go
// projects.Model has Update that returns (projects.Model, tea.Cmd)
// We need a wrapper that returns (tea.Model, tea.Cmd) but preserves the type

type ProjectsPanel struct {
    *projects.Model
}

func (p ProjectsPanel) Init() tea.Cmd { return p.Model.Init() }

func (p ProjectsPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    m, cmd := p.Model.Update(msg)
    p.Model = &m       // Update the pointer
    return p, cmd      // Return the wrapper
}

func (p ProjectsPanel) View() tea.View { return p.Model.View() }

// Implement SelectionProvider
func (p ProjectsPanel) SelectedLabel() string {
    if item := p.Model.List.SelectedItem(); item != nil {
        return item.FilterValue()
    }
    return ""
}
```

## Left Panel (Navigation)

The left panel typically contains a list for navigation:

```go
type leftPanel struct {
    list   list.Model
    width  int
    height int
}

func newLeftPanel(items []list.Item) *leftPanel {
    l := list.New(items, list.NewDefaultDelegate(), 0, 0)
    l.Title = "Projects"
    l.SetShowHelp(false)
    return &leftPanel{list: l}
}

func (m *leftPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.list.SetSize(msg.Width, msg.Height)
        return m, nil

    case tea.KeyPressMsg:
        if msg.String() == "enter" {
            if item := m.list.SelectedItem(); item != nil {
                return m, func() tea.Msg { return SelectItemMsg{item} }
            }
        }
    }

    var cmd tea.Cmd
    m.list, cmd = m.list.Update(msg)
    return m, cmd
}
```

## Main Panel (Content)

The main panel displays the primary content, often a table:

```go
type mainPanel struct {
    table       table.Model
    loading     bool
    spinnerView string
}

func (m *mainPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.table.SetSize(msg.Width, msg.Height)
        return m, nil

    case DataFetchedMsg:
        m.loading = false
        m.table.SetRows(msg.Rows)
        return m, nil
    }

    var cmd tea.Cmd
    m.table, cmd = m.table.Update(msg)
    return m, cmd
}

func (m *mainPanel) View() tea.View {
    return tea.NewView(table.RenderPanel(
        m.theme,
        &m.table,
        m.loading,
        m.spinnerView,
        "Items",
    ))
}
```

## Right Panel (Details)

The right panel shows details for the selected item:

```go
type detailsPanel struct {
    content string
    width   int
    height  int
}

func (m *detailsPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width, m.height = msg.Width, msg.Height
        return m, nil

    case DetailsFetchedMsg:
        m.content = msg.Content
        return m, nil
    }
    return m, nil
}

func (m *detailsPanel) View() tea.View {
    return tea.NewView(lipgloss.NewStyle().
        Width(m.width).
        Height(m.height).
        Render(m.content))
}
```

## Table Helper

tuishell provides `table.RenderPanel()` to handle loading/empty states:

```go
import "github.com/felipeospina21/tuishell/table"

func (m *mainPanel) View() tea.View {
    return tea.NewView(table.RenderPanel(
        m.theme,       // style.Theme
        &m.table,      // *table.Model
        m.loading,     // bool - show loading spinner
        m.spinnerView, // string - spinner view from shell
        "Header",      // string - table header text
    ))
}
```

This renders:
- Loading spinner when `loading` is true
- Empty message when table has no rows
- Header + table when table has rows
