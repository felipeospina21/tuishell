# Mouse Support

tuishell provides the hooks needed to add mouse support to an app without
fragile geometry hacks. Bubble Tea v2 delivers mouse events natively once
`view.MouseMode = tea.MouseModeCellMotion` is set; tuishell handles routing
clicks to the correct panel and exposes the accessors a table needs for exact
hit-testing.

Three layers cooperate:

1. **Shell** — enables mouse mode, routes mouse events to the panel under the
   pointer, and focuses a panel on click (click-to-focus).
2. **Layout** — translates absolute coordinates to a panel via
   `Layout.PanelAt`.
3. **Table** — scrolls with the wheel, selects a row on click, and (optionally)
   marks rows for [bubblezone](https://github.com/lrstanley/bubblezone).

tuishell does **not** depend on bubblezone. If you want per-row hit-testing you
supply your own `*zone.Manager`, which satisfies the small `table.ZoneManager`
interface.

> **Runnable demo:** see [`shell/example-mouse`](../shell/example-mouse) for a
> complete app wiring `EnableMouse`, click-to-focus, wheel scroll, and
> bubblezone row selection. Run it with `cd shell/example-mouse && go run .`.

## Enabling mouse in the shell

Set `EnableMouse` in the shell config. This makes `RenderView()` enable
`tea.MouseModeCellMotion` (in addition to `AltScreen`) and makes `Update()`
route mouse messages by coordinate with click-to-focus.

```go
m := shell.New(shell.Config{
    Theme:       theme,
    LeftPanel:   &leftPanel{},
    MainPanel:   &mainPanel{},
    Keybinds:    tuishell.GlobalKeys(false),
    EnableMouse: true, // <- enable mouse
})
```

With `EnableMouse` set:

- A left click inside a visible panel sets `Ctx.FocusedPanel` to that panel.
- The original mouse message (with **absolute** coordinates) is forwarded to
  the panel under the pointer — not just the focused panel. Absolute
  coordinates are preserved so bubblezone (which scans at the root) keeps
  working; use `Layout.PanelAt` when you need panel-relative coordinates.
- While a modal is open, or when `EnableMouse` is false, mouse messages fall
  back to the focused panel (legacy behavior).

## Hit-testing panels: `Layout.PanelAt`

`Layout` exposes absolute rectangles (`LeftRect`, `MainRect`, `RightRect`) and a
helper that maps an absolute coordinate to a panel plus panel-relative
coordinates:

```go
panel, localX, localY, ok := m.Layout.PanelAt(mouse.X, mouse.Y)
if ok {
    // panel is LeftPanel / MainPanel / RightPanel
    // localX, localY are relative to that panel's top-left content cell
}
```

`PanelAt` accounts for the main frame's border/padding and the left-panel width,
so callers never duplicate layout math. It returns `ok == false` for clicks on
the frame border, the statusline, or empty space.

## Table mouse handling

`table.Model.Update` handles mouse events natively:

- **Wheel** — `MouseWheelUp`/`MouseWheelDown` move the cursor by
  `table.MouseWheelDelta` (1) row per tick, respecting the usual clamping.
- **Left click / release** — selects the row under the pointer when the Y
  coordinate maps to a visible row.

Key events remain gated on focus. Mouse events are gated on focus too by
default; set `AllowMouseWhenBlurred` (or use `table.WithAllowMouseWhenBlurred`)
to let an unfocused-but-visible table respond to the wheel and clicks:

```go
tbl := table.New(
    table.WithColumns(cols),
    table.WithRows(rows),
    table.WithAllowMouseWhenBlurred(true),
)
```

The table's `RowAt(y)` expects a Y coordinate **relative to the table's own
`View()` origin**. Because the shell forwards absolute coordinates, your panel
is responsible for translating them to table-local space (subtract the panel
origin from `PanelAt`, then any offset at which you render the table inside the
panel). Alternatively, use bubblezone (below) and let it handle translation.

### Geometry accessors

For manual geometric hit-testing the table exposes:

| Accessor          | Meaning                                                        |
| ----------------- | -------------------------------------------------------------- |
| `Start() int`     | Index of the first visible row (scroll offset).                |
| `VisibleRange()`  | `(start, end)` visible row indices; `end` is exclusive.        |
| `RowPitch() int`  | Terminal lines per row (`RowHeight + RowBottomMargin`).        |
| `HeaderHeight()`  | Lines above the first row within `View()` (header + separator).|
| `RowAt(y int)`    | `(index, ok)` for the row at a `View()`-relative Y coordinate. |

The exported constants `table.RowHeight` and `table.RowBottomMargin` document
the per-row layout so you never hardcode `2`.

## bubblezone: per-row click zones

The robust way to make rows clickable is [bubblezone](https://github.com/lrstanley/bubblezone),
which avoids brittle geometry entirely. tuishell stays bubblezone-free: you pass
your manager to the table, and it wraps each rendered row in a zone marker.

> `MaxWidth`/`MaxHeight` corrupt zone markers. tuishell marks the **outer**,
> already-trimmed row string (not the inner cells that use `MaxWidth`), so
> markers survive.

### 1. Create a manager and give it to the table

```go
import zone "github.com/lrstanley/bubblezone/v2"

func main() {
    zone.NewGlobal()
    // ...
}

tbl := table.New(
    table.WithColumns(cols),
    table.WithRows(rows),
    table.WithZoneManager(zone.DefaultManager),        // satisfies table.ZoneManager
    table.WithRowID(func(i int) string {               // unique id per row
        return fmt.Sprintf("row-%d", i)
    }),
)
```

`table.ZoneManager` is a one-method interface satisfied by a bubblezone
manager:

```go
type ZoneManager interface {
    Mark(id, v string) string
}
```

### 2. Scan at the root model's `View()`

bubblezone requires a single `Scan` at the root, wrapping the whole rendered
frame:

```go
func (a app) View() tea.View {
    v := a.shell.RenderView()          // AltScreen + MouseModeCellMotion (EnableMouse)
    v.SetContent(zone.Scan(v.Content)) // register + strip markers
    return v
}
```

### 3. Hit-test in your panel's `Update`

Because the shell forwards absolute-coordinate mouse messages to the panel under
the pointer, your panel can check zones directly:

```go
func (m *mainPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.MouseClickMsg:
        if msg.Button != tea.MouseLeft {
            break
        }
        for i := range m.table.Rows() {
            if zone.Get(fmt.Sprintf("row-%d", i)).InBounds(msg) {
                m.table.SetCursor(i)
                break
            }
        }
    }
    var cmd tea.Cmd
    m.table, cmd = m.table.Update(msg)
    return m, cmd
}
```

## Requirements recap

- `view.AltScreen = true` **and** `view.MouseMode = tea.MouseModeCellMotion`
  (both set by `RenderView()` when `EnableMouse` is true).
- bubblezone only works in alt-screen mode; call `zone.Scan` only at the root.
- Use `lipgloss.Width()` (not `len()`) around marked content.

## Fixed (always-open) three-column layout

By default the shell shows either **left + main** or **main + right** (opening
the right panel closes the left). Set `FixedPanels` to lay out all three
columns at once and keep them open for the shell's lifetime:

```go
m := shell.New(shell.Config{
    Theme:       theme,
    LeftPanel:   &leftPanel{},
    MainPanel:   &mainPanel{},
    RightPanel:  &rightPanel{},
    FixedPanels: true,
})
```

When `FixedPanels` is set:

- Both side panels start open and cannot be closed. `CloseLeftPanelMsg`,
  `CloseRightPanelMsg`, `ToggleFullscreenMsg`, and the global toggle/close
  keybinds (`ctrl+o`, `esc`) become no-ops.
- The layout allocates real width to all three columns (left fixed width, the
  remainder split between main and right).

`FixedPanels` is independent of `EnableMouse` — a fixed layout works for
keyboard-only apps too. Read it back with `m.IsFixedPanels()`.

## Toggle buttons (mouse only)

When `EnableMouse` is true **and** `FixedPanels` is false, the shell renders a
one-row button bar at the top of the content area with a left-panel and a
right-panel toggle. Clicking a toggle opens or closes that panel. The bar is
hidden when mouse support is off (nothing to click without a pointer) or when
panels are fixed (nothing to toggle).

The button labels are consumer-defined, since the panel content is
app-specific. Set `LeftButtonLabel` / `RightButtonLabel` in the config; the
shell wraps each label in brackets with a directional arrow (`«` closes an open
panel, `»` opens a closed one). They default to `"nav"` and `"details"`:

```go
m := shell.New(shell.Config{
    // ...
    EnableMouse:      true,
    LeftButtonLabel:  "menu",  // renders as "[ « menu ]"
    RightButtonLabel: "info",  // renders as "[ info » ]"
})
```

The bar's clickable regions are exposed on `Layout` for consumers who want to
hit-test them too (their widths track the configured labels automatically):

```go
if panel, ok := m.Layout.ButtonAt(mouse.X, mouse.Y); ok {
    // panel is LeftPanel or RightPanel — the panel that button toggles
}
```

The shell handles these clicks internally in `handleMouse` before panel
routing, so you get the buttons for free by setting `EnableMouse`.

## Panel-cycle keybind

A global keybind cycles focus through the currently visible panels
(left → main → right, skipping closed panels and wrapping around):

- `tab` — focus the next panel (`GlobalKeyMap.CyclePanel`)
- `shift+tab` — focus the previous panel (`GlobalKeyMap.CyclePanelBack`)

This works with or without mouse support and is the recommended way to move
focus in a fixed-panel layout.

## Active-panel indicator

The shell recolors the **focused** panel's border using `FocusedBorderColor`
(defaults to `Theme.Primary`) as an always-on active-panel indicator — no mouse
required. The left panel and the outer main frame are recolored automatically
when focused, as long as they have a border (a borderless style is left
untouched).

The **right** panel renders its own content (the shell does not wrap it in a
style), so recolor its border yourself using the shell's accessors:

```go
func (p *rightPanel) View() tea.View {
    color := lipgloss.Color("#585b70") // idle
    if p.focused {                      // set from Ctx.FocusedPanel
        color = p.focusColor            // = shell.FocusBorderColor()
    }
    body := lipgloss.NewStyle().
        Border(lipgloss.NormalBorder(), false, false, false, true).
        BorderForeground(color).
        Render(p.content)
    return tea.NewView(body)
}
```

`m.FocusedPanel()` and `m.FocusBorderColor()` expose the current focus and the
indicator color so consumer-rendered panels match the shell's built-in
indicator.
