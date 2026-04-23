# Theming

tuishell uses a semantic color token system. All components receive a `Theme` at construction and use its tokens for consistent styling.

## Theme Struct

```go
type Theme struct {
    // Primary accent colors
    Primary       color.Color  // Main accent (titles, selections)
    PrimaryBright color.Color  // Hover/active states
    PrimaryFg     color.Color  // Text on primary background
    PrimaryDim    color.Color  // Subtle primary backgrounds

    // Semantic colors
    Info          color.Color
    InfoBright    color.Color
    Success       color.Color
    SuccessBright color.Color
    Danger        color.Color
    DangerBright  color.Color
    Warning       color.Color
    WarningBright color.Color
    Caution       color.Color

    // Neutral colors
    Text            color.Color  // Primary text
    TextInverse     color.Color  // Text on light backgrounds
    TextDimmed      color.Color  // Secondary text
    Muted           color.Color  // Tertiary text
    Dim             color.Color  // Very subtle text
    Border          color.Color  // Panel borders
    ModalBorder     color.Color  // Modal overlay border
    SurfaceDim      color.Color  // Subtle backgrounds
    SelectionBorder color.Color  // Selected item border

    // Statusline colors
    StatusText    color.Color  // Statusline text
    StatusNormal  color.Color  // Normal mode background
    StatusLoading color.Color  // Loading mode background
    StatusError   color.Color  // Error mode background
    StatusDemo     color.Color  // Demo mode indicator
    StatusAccent1 color.Color  // Accent segment 1
    StatusAccent2 color.Color  // Accent segment 2
}
```

## Creating a Theme

Define all 30 tokens to create your theme:

```go
import "charm.land/lipgloss/v2"

myTheme := style.Theme{
    Primary:       lipgloss.Color("#0052CC"),
    PrimaryBright: lipgloss.Color("#2684FF"),
    PrimaryFg:     lipgloss.Color("#FFFFFF"),
    PrimaryDim:    lipgloss.Color("#0747A6"),

    Info:          lipgloss.Color("#0065FF"),
    InfoBright:    lipgloss.Color("#4C9AFF"),
    Success:       lipgloss.Color("#36B37E"),
    SuccessBright: lipgloss.Color("#57D9A3"),
    Danger:        lipgloss.Color("#FF5630"),
    DangerBright:  lipgloss.Color("#FF7452"),
    Warning:       lipgloss.Color("#FFAB00"),
    WarningBright: lipgloss.Color("#FFC400"),
    Caution:       lipgloss.Color("#FF8B00"),

    Text:            lipgloss.Color("#C4C4C4"),
    TextInverse:     lipgloss.Color("#111"),
    TextDimmed:      lipgloss.Color("#777777"),
    Muted:           lipgloss.Color("#999999"),
    Dim:             lipgloss.Color("#444444"),
    Border:          lipgloss.Color("#3f4145"),
    ModalBorder:     lipgloss.Color("#666666"),
    SurfaceDim:      lipgloss.Color("#1e1e24"),
    SelectionBorder: lipgloss.Color("#0052CC"),

    StatusText:    lipgloss.Color("#FFFFFF"),
    StatusNormal:  lipgloss.Color("#0052CC"),
    StatusLoading: lipgloss.Color("#0065FF"),
    StatusError:   lipgloss.Color("#FF5630"),
    StatusDemo:     lipgloss.Color("#36B37E"),
    StatusAccent1: lipgloss.Color("#6554C0"),
    StatusAccent2: lipgloss.Color("#403294"),
}
```

## Using Theme in Components

Pass the theme when creating the shell:

```go
s := shell.New(shell.Config{
    Theme: jiraTheme,
    // ...
})
```

For custom components, access theme colors:

```go
func (m *myPanel) View() tea.View {
    titleStyle := lipgloss.NewStyle().
        Foreground(m.theme.Primary).
        Bold(true)
    
    return tea.NewView(titleStyle.Render("Title"))
}
```

## Component Styles

tuishell components expose themed style functions:

```go
// Table styles
table.ThemedStyles(theme)     // Returns bubbles table.Styles
table.TitleStyle(theme)       // Table header style
table.DocStyle(theme)         // Table container style

// Modal styles
modal.OverlayStyle(theme)     // Dim background
modal.ContentStyle(theme)     // Modal content box

// Statusline styles
statusline.ModeStyle(theme, mode)  // Mode indicator style
```
