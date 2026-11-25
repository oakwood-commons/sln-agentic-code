---
name: finish-issue
description: "Local 'Issue Definition of Done' quality gate: run static checks, dispatch code-reviewer and artifact-auditor IN PARALLEL over one frozen diff, fix whole defect classes in one pass, and run final lint before proposing any commit."
---

# Finish-Issue Local Quality Gate (`$0` Pre-PR Gate)

Run this gate locally as the final phase of any substantive change, before
proposing a commit or opening a pull request. Catching bugs locally costs zero
CI minutes and zero GitHub Copilot review credits.

## Phase 1 -- Build, Test & Static Security Scan

1. Inspect `git status --short` and `git diff --shortstat`. Stop if clean.
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
