# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v0.1.0] - 2026-04-16

### Added

- 3-panel shell layout with left (navigation list), main (content table), and right (details) panels
- Panel routing with focus management and keyboard navigation between panels
- Task lifecycle management with automatic loading states, spinner, and error handling
- Table widget with scrolling, styled rows, icon columns, and proportional column widths
- Modal overlay with dim background, copy/submit actions, and keybindings
- Statusline with mode indicator, project label, spinner, and help keybindings
- Theme system with 30 semantic color tokens and Tailwind-style color palettes
- Hub app picker for multi-app launcher
- Loader component for loading spinner views
- Global keybindings for help modal, quit, panel toggle, and dev-mode keys
- `SelectionProvider` interface for automatic statusline label updates
- `RenderPanel` helper for table loading/empty states
