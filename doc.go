// Package tuishell is a reusable 3-panel TUI framework for Bubble Tea applications.
//
// It handles the common infrastructure of panel-based terminal UIs:
// layout computation, panel routing, modal overlays, statusline,
// task management, and global keybindings — so you can focus on your
// domain logic.
//
// # Packages
//
//   - [shell] — core shell model that manages layout, routing, and state
//   - [style] — semantic color theme with 30 tokens
//   - [table] — themed table widget with loading/empty states
//   - [modal] — centered overlay modal with copy/submit actions
//   - [statusline] — mode indicator, project label, spinner, help keybindings
//   - [loader] — themed loading spinner
//   - [hub] — app launcher for hosting multiple tuishell apps
//
// # Usage
//
// Create panels that implement [tea.Model], then pass them to [shell.New]:
//
//	s := shell.New(shell.Config{
//	    Theme:     style.DefaultTheme(),
//	    LeftPanel: &myLeftPanel{},
//	    MainPanel: &myMainPanel{},
//	    Keybinds:  tuishell.GlobalKeys(false),
//	    LeftPanelStyle: lipgloss.NewStyle().Width(30),
//	})
package tuishell
