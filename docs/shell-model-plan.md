# ShellModel — Reusable App Shell

## Problem

Every tuishell-based app (mrglab, mrjira) must manually wire the same behaviors: panel toggling, modal lifecycle, error display, spinner, statusline, resize handling, focus routing. This is ~200 lines of boilerplate per app that should live in tuishell.

## Goal

Provide a `ShellModel` that handles all common 3-panel TUI behavior. Apps only implement their panel content and communicate via messages.

## API Design

```go
shell := tuishell.NewShell(tuishell.ShellConfig{
    Theme:       myTheme,
    LeftPanel:   projectsList,    // tea.Model
    MainPanel:   issuesTable,     // tea.Model
    RightPanel:  nil,             // tea.Model, optional
    StatusLabel: "🎫 mrjira",
    Keybinds:    projects.Keybinds,
})
```

### ShellConfig

```go
type ShellConfig struct {
    Theme           style.Theme
    LeftPanel       tea.Model       // required
    MainPanel       tea.Model       // required
    RightPanel      tea.Model       // optional, nil = no right panel
    StatusLabel     string          // e.g. "🎫 mrjira" or " mrglab"
    Keybinds        help.KeyMap     // initial help keybinds
    LeftPanelWidth  int             // default 30
    LeftPanelStyle  lipgloss.Style  // optional override
    RightPanelStyle lipgloss.Style  // optional override
    MainFrameStyle  lipgloss.Style  // optional override
}
```

### ShellModel

```go
type ShellModel struct {
    // Panels
    Left       tea.Model
    Main       tea.Model
    Right      tea.Model

    // Built-in components
    Modal      modal.Model
    Statusline statusline.Model
    Spinner    spinner.Model

    // State
    Layout     Layout
    Ctx        AppContext

    // Internal
    isLeftOpen        bool
    isRightOpen       bool
    isRightFullscreen bool
    isModalOpen       bool
    taskStatus        taskStatus
    taskErr           error
}
```

## Built-in Behaviors

### 1. Panel Focus Routing

ShellModel routes `tea.KeyPressMsg` to the focused panel automatically:

```
FocusedPanel == LeftPanel  → msg goes to Left.Update()
FocusedPanel == MainPanel  → msg goes to Main.Update()
FocusedPanel == RightPanel → msg goes to Right.Update()
FocusedPanel == ModalPanel → msg goes to Modal.Update()
```

Global keys are checked first, before routing.

### 2. Global Keys (handled by shell, not panels)

| Key | Action |
|-----|--------|
| `ctrl+c` | Quit |
| `ctrl+o` | Toggle left panel, adjust focus |
| `?` | Open help modal with current keybinds |
| `@` | Open error modal (if task error exists) |

### 3. Modal Lifecycle

Shell handles these modal messages automatically:

| Message | Shell Action |
|---------|-------------|
| `modal.CloseModalMsg` | Close modal, restore focus to previous panel |
| `modal.CopyModalMsg` | Copy modal content to clipboard |
| `modal.ResetHighlightMsg` | Reset highlight flag |
| `modal.SubmitModalMsg` | Forward to app via `ShellSubmitMsg` |

Apps trigger modals by returning messages:

```go
// From any panel's Update():
return m, func() tea.Msg {
    return tuishell.OpenModalMsg{Header: "Error", Content: err.Error(), IsError: true}
}
```

### 4. Task Lifecycle

Shell manages loading/error states:

```go
// Panel returns this to start an async task
type StartTaskMsg struct {
    Cmd tea.Cmd  // the async command to run
}

// Shell automatically:
// 1. Sets statusline to Loading + spinner
// 2. Runs the command
// 3. On result, panel handles the typed msg
// 4. Panel returns FinishTaskMsg to clear loading state

type FinishTaskMsg struct {
    Err      error
    Keybinds help.KeyMap  // restore these keybinds
}
```

### 5. Layout & Resize

Shell handles `tea.WindowSizeMsg` automatically:
- Recomputes layout via `ComputeLayout()`
- Pushes dimensions to all panels via `tea.WindowSizeMsg` (panels receive their allocated size, not the full terminal)
- Updates statusline width

### 6. Spinner

Shell owns the spinner. Panels that need a spinner view receive it via:

```go
type SpinnerViewMsg struct {
    View string
}
```

### 7. Right Panel

Shell handles open/close/fullscreen:

```go
// Panel returns these to control right panel
type OpenRightPanelMsg struct{}
type CloseRightPanelMsg struct{}
type ToggleFullscreenMsg struct{}
```

### 8. View Composition

Shell composes the view automatically:

```
┌─────────────────────────────────────────┐
│ LeftPanel │ MainPanel    │ RightPanel   │
│           │              │ (optional)   │
│           │              │              │
├─────────────────────────────────────────┤
│ Statusline                              │
└─────────────────────────────────────────┘
│ Modal overlay (when open)               │
```

## Message Flow

```
tea.KeyPressMsg
    │
    ├─ Shell checks global keys (ctrl+c, ctrl+o, ?, @)
    │   └─ handled → no forwarding
    │
    ├─ Modal open? → forward to Modal.Update()
    │   └─ CloseModalMsg → shell closes modal
    │
    └─ Route to focused panel → panel.Update()
        └─ Panel returns cmd/msg
            ├─ StartTaskMsg → shell starts loading
            ├─ FinishTaskMsg → shell clears loading
            ├─ OpenModalMsg → shell opens modal
            ├─ OpenRightPanelMsg → shell opens right
            ├─ CloseRightPanelMsg → shell closes right
            ├─ SetKeybindsMsg → shell updates help
            ├─ SetStatusMsg → shell updates statusline
            └─ anything else → normal tea.Cmd flow
```

## What Apps Still Own

- Panel `tea.Model` implementations (content, keys, domain logic)
- Domain-specific messages (e.g. `MRListFetchedMsg`, `issuesFetchedMsg`)
- API clients and config
- Icon mappings and domain-specific styles
- Form handling (create MR, respond to comment)

## Migration Path

### Phase 1: Extract ShellModel
- Create `tuishell/shell.go` with `ShellModel`, `ShellConfig`
- Implement panel routing, global keys, resize, layout
- Implement modal lifecycle, task lifecycle, spinner

### Phase 2: Migrate mrjira
- Replace `mrjira/internal/tui/app/app.go` with ShellModel
- mrjira panels become simple `tea.Model`s
- Validate all behaviors work

### Phase 3: Migrate mrglab
- Replace `mrglab/internal/tui/app/` with ShellModel
- Keep domain logic in panel models
- Ensure backward compatibility (tests, visual parity)

### Phase 4: Cleanup
- Remove duplicated wiring code from both apps
- Remove shell-level logic from individual components
- Update docs

## Risks

1. **mrglab's Update is complex** — form handling, pending notes, multiple task types. ShellModel needs to be flexible enough to not block these patterns.
2. **Panel communication** — panels sometimes need to talk to each other (e.g. selecting a project triggers MR fetch). Shell needs a message bus or the app wraps ShellModel.
3. **Custom View overrides** — mrglab's View has special cases (initial screen, fetching state). Shell's View must support these or allow overrides.

## Mitigation

- ShellModel should be embeddable, not final. Apps can embed it and override `Update`/`View` for custom behavior.
- Use a `PreUpdate` hook pattern: shell calls `panel.Update()` and inspects returned messages for shell-level actions.
- Keep ShellModel's View composable — expose `RenderLeft()`, `RenderMain()`, `RenderRight()`, `RenderStatusline()` so apps can rearrange if needed.
