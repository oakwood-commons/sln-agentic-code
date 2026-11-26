---
name: finish-issue
description: "Use automatically before pushing or opening a pull request (or via /finish-issue): runs static checks, dispatches code-reviewer and artifact-auditor IN PARALLEL over the branch diff, fixes defect classes in one pass, and runs final lint."
---

# Finish-Issue Pre-PR Self-Review Gate

Run this gate automatically on your feature branch before pushing and opening a
pull request (or against uncommitted changes when asked directly). Deterministic
checks (`go test -race`, linters, `gosec`) run free on your CPU, and catching
defects before opening the PR saves CI minutes, GitHub Copilot PR review
credits, and human review time.

## Phase 1 -- Build, Test & Static Security Scan

1. Inspect `git status --short` and the diff (`git diff` if uncommitted, or
   `git diff origin/main...HEAD` on a committed feature branch). Stop if empty.
2. Run the repository's build and test commands (`scafctl run action ci`,
   `task test`, or `go test -race ./...`). Fix any red test first.
3. Run local static security checks (`gosec ./...` for Go repos) and smoke-test
   CLI `--help` + one valid and one invalid invocation when CLI flags changed.

## Phase 2 -- Parallel Review Over One Frozen Diff

Dispatch `code-reviewer` and `artifact-auditor` **in parallel in a single turn**
over the same diff:

- **Standard diff (`<= 800` changed lines)**: 1 `code-reviewer` + 1
  `artifact-auditor`.
- **Large or security-sensitive diff (`> 800` lines or auth/concurrency/secrets)**:
  fan `code-reviewer` out into 3 focused parallel lenses alongside
  `artifact-auditor`:
  1. *Callers, Flags & Upstream Contracts* (`$GOMODCACHE` / cross-repo grep)
  2. *Security, Trust Boundaries & File Modes*
  3. *Concurrency, Resource Lifecycle & Test Reality*

## Phase 3 -- Merge & Fix Defect Classes

Merge all findings into one list before editing:

- Fix every `BLOCKER` and `MAJOR` issue and every real missing artifact.
- **Fix the defect class, not just the single line**: whenever a finding has
  sibling occurrences elsewhere in the diff, fix all of them in the same pass.
- **Pause and ask the human** whenever a finding involves a real design or API
  trade-off rather than a clear bug.

## Phase 4 -- Final Lint & Summary

Re-run the repository's tests and linter on the final diff, then summarize
findings resolved, artifact status, and verification output. Never commit or
push automatically inside this gate.
