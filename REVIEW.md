# Copilot code review: how to review this repo

You are reviewing as an **independent final reviewer**, not a re-run of this
repository's authoring workflow. Your job is to find substantive bugs, broken
contracts, security or publish-hygiene leaks, and missing tests that local
checks missed -- not to nitpick style.

## Sources of truth

Do not invent parallel repo conventions -- read the canonical contract directly:

- [`AGENTS.md`](AGENTS.md) -- build/test/lint commands, Go conventions,
  `scafctl` validation gates, Markdown rules, and public publish hygiene.

Where this file and `AGENTS.md` conflict on a repository convention,
`AGENTS.md` wins; this file governs *how* to conduct the review.

## Method: review the whole change before writing any findings

1. **Read the full diff and PR description together** to understand what
   changed and why before judging individual hunks.
2. **Inventory every changed interface**: Go functions/types in
   `examples/go-app/`, `scafctl` resolver parameters and action names in
   `actions.yaml` / `ci/` / `scafctl/`, CI wrapper inputs in `.github/` and
   `ci/tekton/`, and documented commands across `docs/`. Verify callers,
   templates, and docs stay in sync with those interfaces.
3. **Check test coverage against the actual change**: confirm new or modified
   Go behavior has real edge-case and failure-path assertions (`go test -race
   ./...`, 70%+ coverage) and changed `scafctl` solutions keep valid
   `spec.testing.cases` assertions.

## What to weigh heavily

- **Regressions and edge cases**: empty/nil inputs, missing files, non-zero
  exit codes, cross-platform (`linux`, `darwin`, `windows`) shell/path bugs in
  `scafctl` `exec` actions, and error wrapping (`%w` / `errors.Is`).
- **Interface drift across docs, CI wrappers, and starter templates**: commands,
  resolver flags (`-r key=value`), or file paths in `docs/` that do not match
  `actions.yaml`, `ci/`, or `scafctl/`.
- **Trust boundaries and shell safety**: unquoted shell interpolation, path
  traversal outside target directories, overly permissive file modes, or
  destructive GitHub API calls without opt-in flags.
- **Publish hygiene**: this repository is public -- flag any hardcoded user
  home paths, internal hostnames, private package references, or credentials in
  docs, starter templates, or tests.
- **License & attribution integrity**: changes touching adapted third-party
  skills must keep [`THIRD-PARTY-NOTICE.txt`](THIRD-PARTY-NOTICE.txt) accurate.
- **PR title and body accuracy**: this repo squash-merges using `PR_TITLE` and
  `PR_BODY`. Flag a non-conventional title or a PR body that no longer matches
  `HEAD`.

## What not to spend the review on

- Style nits already enforced by `gofmt`, `golangci-lint`, or `scafctl lint`.
- Speculative `"Consider ..."` abstractions or subjective prose rewrites with
  no accuracy or clarity bug.
- Defensive guards on internal paths that have no external caller.

## Follow the implication

If one finding reveals a class of mistake (for example, an unquoted shell path,
a stale CLI flag in docs, or a missing fence escape around nested backticks),
search the rest of the diff for sibling occurrences and report the defect class
once covering all locations rather than drip-feeding one comment per round.

## Re-reviews (after a fix is pushed)

Scope re-reviews strictly to `<last_reviewed_commit>..HEAD`:

1. **Verify prior findings and their fixes**, checking whether the fix missed a
   sibling occurrence or introduced a local regression.
2. **Inspect only files modified since `<last_reviewed_commit>`** plus direct
   consumers of any interface changed in that delta. Do not re-audit untouched
   files that already passed an earlier round.
3. **Confirm PR title and body sync** with `HEAD`.
