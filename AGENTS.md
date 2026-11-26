# AGENTS.md

Canonical repository guidance for any contributor or AI coding tool working in
`oakwood-commons/sln-agentic-code`.

This file states **what this repository requires** (commands, quality gates,
conventions, and publish hygiene). It intentionally does **not** dictate which
AI harness or personal workflow you use. Pull request review rules live
separately in [`REVIEW.md`](REVIEW.md).

## Overview

Published `scafctl` solution (`scafctl/solution.yaml` -> `agentic-code`),
Blueprint documentation site (Hugo `hugo-book`), portable `scafctl` Go CI
definitions (`ci/`), and runnable sample Go module (`examples/go-app/`) under
`oakwood-commons`.

## Issue Tracking & Git

- **Issue tracking**:
  - Repository changes to `oakwood-commons/sln-agentic-code`: GitHub Issues via
    `gh` CLI (`gh issue create` / `gh issue view`).
  - Walkthrough practice on `shorty` (`examples/go-app/`): local Markdown issues
    in git-ignored `temp/issues/<NNN>-<slug>.md` (unless the user is working on
    their own GitHub fork).
- **Conventional commits**: `feat:`, `fix:`, `docs:`, `chore:`, `ci:`, `test:`
  (`!` for breaking changes).
- **Signing & DCO**: every commit is signed and carries a `Signed-off-by:`
  trailer (`git commit -s -S`).
- **Squash-merge repository**: every PR becomes one commit on `main` built from
  the **PR title and body** (`PR_TITLE` / `PR_BODY`). Keep the PR title a valid
  conventional-commit subject and update the PR body whenever a new push
  changes what the PR does.
- **Copilot review discipline**: Copilot reviews ready-for-review PRs on open
  (`review_on_push: false`). Always wait for CI checks to pass before
  requesting a re-review, and only re-request Copilot after pushing code fixes
  for prior open findings.

## Build, Test & Lint (`scafctl` Actions)

Root [`actions.yaml`](actions.yaml) is auto-discovered by `scafctl`:

~~~shell
# Install or verify required tools (hugo, gh, go, golangci-lint, gosec) across Linux/macOS/Windows
scafctl run action install-tools

# Run full repository verification (Go CI on examples/go-app + scafctl lint & tests + Hugo build)
scafctl run action ci

# Lint and run functional tests across all repository scafctl solutions
scafctl run action lint-solutions
scafctl run action test-solutions

# Run Go CI against examples/go-app via the root shorthand
scafctl run action go-ci

# Build Hugo documentation site into ./public
scafctl run action docs-build

# Serve Hugo documentation locally at http://localhost:1313
scafctl run action docs-serve

# Configure GitHub repository settings & rulesets idempotently
scafctl run action configure-github -r owner=oakwood-commons -r repo=sln-agentic-code

# Run portable Go CI directly against examples/go-app
scafctl run action ci -f ci/go-ci/actions.yaml -r workdir=examples/go-app
~~~

## Go Conventions (`examples/go-app/`)

- Isolated Go module (`GOWORK=off`) so `examples/go-app/go.mod` and
  `site/go.mod` remain independent.
- Wrap errors with context and `%w` (`fmt.Errorf("...: %w", err)`), preserve
  `errors.Is` / `errors.As` behavior, and never panic in library packages.
- Maintain **70%+ test coverage** with table-driven tests (`go test -race ./...`)
  exercising public interfaces rather than private internals, and zero findings
  from `go vet`, `golangci-lint`, and `gosec`.

## `scafctl` Authoring & MCP Truth Source

- Consult `scafctl mcp serve` tools (`explain_kind`, `get_provider_schema`,
  `lint_solution`, `preview_action`, `validate_expression`) rather than
  inventing YAML fields.
- Every `scafctl` solution/action file (`actions.yaml`, `ci/go-ci/actions.yaml`,
  `scafctl/solution.yaml`) must pass `scafctl lint -f <file>`,
  `scafctl test functional -f <file>`, and (for `scafctl/solution.yaml`)
  `scafctl package solution -f ./scafctl/solution.yaml --dry-run --no-vendor`.

## Markdown & Publish Hygiene

- When a fenced code block contains backticks, wrap the outer fence in tildes
  (`~~~`) or 4+ backticks so inner fences never close the block early.
- **No dash punctuation in prose:** never use em-dashes, en-dashes, or spaced
  hyphens/double-hyphens (`" - "`, `" -- "`) between clauses. Rewrite sentences with
  colons, commas, parentheses, semicolons, or separate sentences, and avoid
  emojis.
- **Public repository**: never commit secrets, tokens, private hostnames,
  personal homelab skills, or user-specific filesystem paths (e.g.,
  `/home/<user>` or `/Users/<user>`). Scratch files belong in git-ignored
  `temp/`.
