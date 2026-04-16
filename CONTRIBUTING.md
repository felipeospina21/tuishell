# Contributing to tuishell

## Getting Started

1. Clone the repo:
   ```bash
   git clone https://github.com/felipeospina21/tuishell.git
   cd tuishell
   ```
2. Requires **Go 1.26+**.
3. Build and test:
   ```bash
   go build ./...
   go test ./...
   ```

## Development

tuishell is typically used inside a `go.work` workspace alongside consumer modules (mrglab, jiraf). When working on tuishell standalone, prefix Go commands with `GOWORK=off` to avoid workspace interference:

```bash
GOWORK=off go build ./...
GOWORK=off go test ./...
```

Run `go vet ./...` before submitting changes.

## Pull Requests

- Branch from `main` — never push directly to it.
- Keep PRs focused on a single change.
- Include tests for new features or bug fixes.
- Run checks locally before opening a PR:
  ```bash
  go build ./...
  go test ./...
  go vet ./...
  ```

## Code Style

- Follow standard Go idioms and conventions.
- Format code with `gofmt` (or `goimports`).
- Run [`golangci-lint`](https://golangci-lint.run/) if available:
  ```bash
  golangci-lint run
  ```

## Reporting Issues

Use [GitHub Issues](https://github.com/felipeospina21/tuishell/issues). Include:

- What you expected vs. what happened
- Go version (`go version`)
- OS and terminal emulator
