---
name: last-look
description: "Use when the user sits down to review or merge ready PRs, asks 'walk me through the PRs to merge' or 'are these ready to merge', or runs /last-look. Verifies live CI, review history, and 3-way mergeability, then briefs strictly one PR per turn in plain English."
---

# Last Look: One-at-a-Time Human Merge Walkthrough

Guide the human owner through ready pull requests **one at a time** so they can
make a confident merge, trim, or hold call in ~60 seconds without scrolling
dozens of test files. Complements `pr-review` (line-level defect hunting) by
serving the final human ownership decision after `pr-autopilot` or local review
converges.

## Rules of Engagement

- **One PR per turn.** Even when 5 PRs are ready, present a 1-line queue header
  (`### 1 of 5: ...`), brief **only the current PR**, and stop for the human's
  response (`merged`, `next`, `trim <X>`, `retrigger`, `skip`, `hold`).
- **Plain, everyday English.** Write in clear, everyday English rather than
  compressed shorthand: explain outcomes from a user's, operator's, or caller's
  perspective instead of reciting internal struct, helper, or variable names.
- **Never merge.** `gh pr merge` is forbidden; the human clicks Merge in the
  GitHub UI.

## Step 1: Verify Live Ground Truth (Never Trust Memory)

Before presenting a PR, verify every fact from live command output on current
`HEAD`:

1. **PR & review state**: Fetch live CI checks, unresolved GraphQL thread count,
   review coverage on `headRefOid`, `PR_TITLE` / `PR_BODY`, `baseRefName`, and
   `headRefName` (`gh pr view`, `gh pr checks`, and `gh api graphql`), then
   fetch both refs (`git fetch origin <base> <headRefName>`) and set
   `<pr_head>` to `origin/<headRefName>` (or `<headRefOid>` for fork PRs via
   `git fetch origin pull/<num>/head`) so commands work whether run inside a
   feature branch or from a `main` Coordinator checkout.
2. **Prod vs. test size & review history**: Run `git log --stat --oneline
   origin/<base>..<pr_head>` and `git diff --numstat origin/<base>...<pr_head>`.
   Separate **production LOC** from **test/doc LOC**, and separate the initial
   implementation commit(s) from post-review fix commits.
3. **Fresh base-branch conflict check**: Run a read-only 3-way merge check
   (`git merge-tree --write-tree <pr_head> origin/<base>`) to catch physical
   conflicts and semantic collisions with recent `<base>` commits before the
   human clicks Merge. Note any file overlap with sibling open PRs in the queue.

Completion criterion: every number, check status, and conflict claim in the
brief comes from live command output on current `HEAD`.

## Step 2: Present the Single-PR Brief

Render one scannable card (omit `Deferred / shortcuts` when empty):

~~~markdown
### <k> of <N>: [<owner>/<repo>#<num>](<url>) | <plain-English summary>

- **What it does**:
  Plain-English explanation of what changes for users, operators, or callers.
- **Still needed & size**:
  Why the problem still matters on current `<base>` + size split (`+P/-M` prod,
  `+T/-U` tests/docs across `F` files).
- **Breaking changes & new surface**:
  Any `!` breaking behavior shifts, new CLI flags/commands, config/schema keys,
  public API symbols, dependency additions, or `.github/workflows` edits (`None,
  internal only` when none). Flag missing `!` on breaking conventional-commit
  titles.
- **Review drift & bloat**:
  Commit 1 vs. `HEAD`: did review rounds fix real bugs or add defensive
  over-engineering? Call out removable bloat and any **declined review threads**
  where human judgment could differ.
- **Blast radius & spot-check files**:
  `Low` / `Medium` / `High` (what breaks if a subtle bug slips through) plus
  clickable links to the **1 to 2 core files** where the real logic lives.
- **Deferred / shortcuts**:
  *(Include only when non-empty)* Scope punted from the issue, `TODO`s left
  behind, or follow-up work assumed.
- **Merge gate**:
  CI status (distinguishing real failures from hung external webhooks), review
  coverage on `HEAD`, **squash-commit check** (PR title and body accurately
  match `HEAD`), and 3-way conflict result against `origin/<base>`.

Ready to merge `#<num>`, trim anything, or skip to next?
~~~

Stop calling tools and end the turn.

## Step 3: Between Merges (`"merged"` / `"next"`)

When the human reports `"merged"`:

1. Confirm `state == MERGED` on GitHub (`gh pr view <num> --json state`).
2. If a local or Herdr linked worktree exists for the merged branch, verify it
   has no uncommitted changes, remove it cleanly, and fast-forward the local
   `<base>` checkout (`git fetch origin <base> && git merge --ff-only
   origin/<base>`).
3. Re-run Step 1 (especially `git merge-tree --write-tree <next_pr_head>
   origin/<base>`) on the **next** PR against the newly updated `<base>`, then
   present the next brief.
