---
name: github-issues
description: "Use when creating, viewing, listing, or filtering GitHub Issues via the gh CLI (for create-issue, planner, git-github-workflow, or work-herd Herdr plans)."
---

# GitHub Issues (`github-issues`)

Use the `gh` CLI to create, read, and shortlist GitHub Issues cleanly across
planning (`planner`), issue authoring (`create-issue`), branch naming
(`git-github-workflow`), and parallel Herdr orchestration (`work-herd`).

## 1. Creating an Issue (After `create-issue` Draft Approval)

Always draft the issue with `create-issue` first and wait for human sign-off
before running `gh issue create`. Pass the body securely to avoid shell quoting
corruption:

~~~bash
# Write body inside a repo-local ignored scratch file (works with sandboxed/Snap gh too)
mkdir -p temp
gh issue create --title "feat(shortener): add optional link expiration" --body-file temp/issue-draft.md
rm -f temp/issue-draft.md
~~~

Add `--label` or `--milestone` when the repository uses them (`gh label list`).

## 2. Reading & Shortlisting Issues (For `planner` and `work-herd`)

- **Single issue pickup:** `gh issue view <N> --comments`
- **Scoped shortlist for Herdr (`work-herd`):** Never pull hundreds of open
  issues at random. Filter by label, milestone, assignee, or search terms first:
  ~~~bash
  gh issue list --state open --label "sprint-next" --limit 20
  gh issue list --state open --search "shortener OR storage" --limit 20
  ~~~

## 3. Branch & PR Linking Conventions

- **Branch name:** `<type>/<N>-<slug>` (for example, `feat/42-expiring-links`).
- **Pull request footer:** Include `Closes #<N>` (or `Refs #<N>` for partial
  work) in the commit and PR body so merging the PR updates the issue
  automatically.
