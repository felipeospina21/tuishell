# Messages

Panels communicate with the shell by returning commands that produce these messages. The shell handles them automatically.

## Task Management

### StartTaskMsg

Tells the shell to enter loading state and run an async command.

```go
return m, func() tea.Msg {
    return tuishell.StartTaskMsg{Cmd: m.fetchData}
}
```

The shell will:
1. Set statusline to loading mode with spinner
2. Execute the provided command
3. Wait for `FinishTaskMsg`

### FinishTaskMsg

Signals that an async task completed.

```go
return m, func() tea.Msg {
    return tuishell.FinishTaskMsg{
        Err:      err,           // nil on success
        Keybinds: panel.Keybinds, // Update statusline help
    }
}
```

The shell will:
1. Set statusline to normal (or error if `Err != nil`)
2. Update help keybindings if provided

## Panel Control

### CloseLeftPanelMsg

Hides the left panel and moves focus to main panel.

```go
return m, func() tea.Msg { return tuishell.CloseLeftPanelMsg{} }
```

If the left panel implements `SelectionProvider`, the shell automatically updates the statusline label to `AppIcon + " " + SelectedLabel()`.

### OpenRightPanelMsg / CloseRightPanelMsg

Shows or hides the right panel.

```go
return m, func() tea.Msg { return tuishell.OpenRightPanelMsg{} }
return m, func() tea.Msg { return tuishell.CloseRightPanelMsg{} }
```

### ToggleFullscreenMsg

Toggles the right panel between normal and fullscreen mode.

```go
return m, func() tea.Msg { return tuishell.ToggleFullscreenMsg{} }
```

## Modal

### OpenModalMsg

Opens the modal overlay with content.

```go
return m, func() tea.Msg {
    return tuishell.OpenModalMsg{
        Header:  "Error",
        Content: err.Error(),
        IsError: true,
    }
}
```

### CloseModalMsg

Closes the modal overlay.

```go
return m, func() tea.Msg { return tuishell.CloseModalMsg{} }
```

### ShellSubmitMsg

Forwarded to your app when the user presses the modal submit key. Handle this in your app's `Update`:

```go
case tuishell.ShellSubmitMsg:
    // User submitted the modal form
    return m, m.processSubmission()
```

## Statusline

### SetKeybindsMsg

Updates the help keybindings shown in the statusline.

```go
return m, func() tea.Msg {
    return tuishell.SetKeybindsMsg{Keybinds: panel.Keybinds}
}
```

### SetStatusMsg

Updates the statusline mode and content.

```go
return m, func() tea.Msg {
    return tuishell.SetStatusMsg{
        Mode:    "INSERT",
        Content: "Editing...",
    }
}
```

## Common Patterns

### Selection with panel close

```go
case tea.KeyPressMsg:
    if msg.String() == "enter" {
        item := m.list.SelectedItem()
        return m, tea.Batch(
            func() tea.Msg { return SelectItemMsg{item} },
            func() tea.Msg { return tuishell.CloseLeftPanelMsg{} },
        )
    }
```

### Fetch with loading state

```go
case SelectProjectMsg:
    return m, func() tea.Msg {
        return tuishell.StartTaskMsg{Cmd: m.fetchProjectData(msg.ID)}
    }

case ProjectDataMsg:
    return m, func() tea.Msg {
        return tuishell.FinishTaskMsg{Err: msg.Err, Keybinds: projectKeybinds}
    }
```

### Open details panel

```go
case tea.KeyPressMsg:
    if msg.String() == "enter" {
        return m, tea.Batch(
            func() tea.Msg { return tuishell.OpenRightPanelMsg{} },
            func() tea.Msg { return tuishell.StartTaskMsg{Cmd: m.fetchDetails} },
        )
    }
```
