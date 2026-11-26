---
name: go
description: "Use when writing, editing, testing, or reviewing Go code, go.mod, go.work, scafctl Go actions, or Taskfiles."
---

# Go Development

- Read the repository's `AGENTS.md` and `REVIEW.md` first. Discover existing
  `scafctl` actions (`actions.yaml`), `Taskfile.yml` targets, or `Makefile`
  targets rather than assuming command names.
- Verify `command -v go`, `go version`, `GOTOOLCHAIN`, `GOWORK`, and `go.mod`
  at runtime. Do not rewrite `go` directives in `go.mod` casually to match your
  local machine.
- Use explicit `GOWORK=off` when testing standalone modules inside a
  multi-module repository so tests run against the intended module.
- Prefer project-defined check commands (`scafctl run action ci`, `task test`,
  `task lint`). Otherwise run `gofmt` on changed files, `go vet ./...`,
  `go test -race ./...`, and `go build ./...`. Never reformat untouched files
  across the entire repo during a scoped edit, and do not install global tools
  without asking.
- Wrap errors with context and `%w` (`fmt.Errorf("...: %w", err)`), preserve
  `errors.Is` / `errors.As` behavior, close files/bodies/goroutines cleanly, and
  propagate `context.Context`.
- Write table-driven subtests exercising observable behavior at public seams;
  avoid introducing external assertion libraries unless the repo already uses
  them.
- Break large features into small, independently testable steps. Do
  not run `go mod tidy` reflexively unless dependencies actually changed.
- Debugging: load `go-debug`. Benchmarks & profiling: load `go-perf`.
