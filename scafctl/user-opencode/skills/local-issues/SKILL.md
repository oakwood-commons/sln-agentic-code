---
name: local-issues
description: "Use when creating, listing, or reading local Markdown issue files in git-ignored temp/issues/<NNN>-<slug>.md for walkthrough practice, spikes, or offline workflows."
---

# Local Markdown Issues (`local-issues`)

Store lightweight issue specs as git-ignored Markdown files under
`temp/issues/<NNN>-<slug>.md` when practicing locally (such as `shorty` in
`examples/go-app/`), prototyping spikes, or working in a repository where you
do not have remote issue-tracker write access.

## 1. Creating a Local Issue (After `create-issue` Draft Approval)

1. Draft the issue using `create-issue` (Plain English at the top for human
   sign-off, Technical Context at the bottom for AI implementation) and confirm
   with the user first.
2. Ensure `temp/` is ignored by Git (check `.gitignore`) so scratch issues are
   never committed accidentally.
3. Pick the next zero-padded 3-digit number in `temp/issues/` (`001`, `002`,
   ...) and write `temp/issues/<NNN>-<slug>.md` (for example,
   `temp/issues/001-expiring-links-and-stats.md`).

## 2. Reading & Planning From Local Issues (`planner` and `work-herd`)

- **Single issue pickup:** When the user says *"look at issue 001"*, match
  `temp/issues/001-*.md`, read it, and run `planner`'s codebase validation
  exactly as you would for a GitHub or Jira ticket.
- **Herdr multi-issue waves (`work-herd`):** List open files in `temp/issues/`
  and map dependencies and hotspot files across `001`, `002`, `003` to build
  `temp/plans/<repo>-herd.md`. Because `git worktree` omits git-ignored
  directories, copy each assigned `temp/issues/<NNN>-<slug>.md` into
  `<worktree-path>/temp/issues/` before starting its worker.

## 3. Branch Naming

- **Branch name:** `<type>/<NNN>-<slug>` (for example,
  `feat/001-expiring-links`).
