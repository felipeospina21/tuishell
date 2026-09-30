# Mouse Support Demo

A runnable demo of tuishell's mouse-support hooks. It highlights everything
added for [issue #31](../../docs/mouse.md):

- **`EnableMouse`** — the shell turns on `tea.MouseModeCellMotion` and routes
  mouse events by coordinate.
- **Click-to-focus** — click the left (nav) or main (table) panel to focus it;
  the `● focused` / `○ click to focus` tag updates.
- **Click a nav item** — the left-panel items are bubblezone-marked; click one
  to select it.
- **Wheel scroll** — scroll the table with the mouse wheel (works even when the
  table isn't focused, via `WithAllowMouseWhenBlurred`).
- **Click-to-open details** — table rows are marked with
  [bubblezone](https://github.com/lrstanley/bubblezone) (`WithZoneManager` +
  `WithRowID`); clicking a row selects it **and opens the right details panel**
  for that row. (The shell lays out left+main *or* main+right, so opening the
  details panel closes the left nav — the same pattern downstream apps use.)
- **Toggle buttons** — because `EnableMouse` is on (and panels aren't fixed),
  the shell renders a top button bar: `[ « nav ]` and `[ details » ]`. Click a
  button to open/close that panel.
- **Panel-cycle key** — `tab` / `shift+tab` moves focus through the visible
  panels (works without a mouse).
- **Active-panel indicator** — the focused panel's border is recolored. The
  shell handles the left panel + outer frame; the details panel recolors its
  own border via `shell.FocusBorderColor()` (see `detailsPanel.View`).

## Run

This demo is a standalone Go module (it has its own `go.mod` and depends on
bubblezone), so run it from inside its directory:

```bash
cd shell/example-mouse && go run .
```

Then use your mouse: click between panels, spin the wheel over the table, and
click individual rows. Keyboard navigation (`j`/`k`, arrows) still works. Press
`q` or `ctrl+c` to quit, `?` for help.

## How it works

| Hook | Where |
| --- | --- |
| `shell.Config{ EnableMouse: true }` | `newApp()` in `main.go` |
| `zone.NewGlobal()` + `zone.Scan(view)` at the root | `main()` / `app.View()` |
| `table.WithZoneManager` + `table.WithRowID` | `newTablePanel()` |
| `zone.Get(id).InBounds(msg)` row hit-test | `tablePanel.Update` |
| click-to-focus (`Ctx.FocusedPanel` set by the shell) | broadcast in `app.Update` |

See [docs/mouse.md](../../docs/mouse.md) for the full reference.
