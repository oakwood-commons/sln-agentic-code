---
name: herdr
description: "Use when setting up Herdr or controlling its workspaces, grouped git worktrees, panes, tabs, or coding agents from inside a Herdr session (HERDR_ENV=1)."
---

# Herdr Workspace & Worktree Control

Before issuing any session-control command, verify you are running inside a
Herdr-managed pane:

~~~shell
test "${HERDR_ENV:-}" = 1
~~~

If unset, stop and tell the user to start or attach this session inside Herdr
(`herdr`). Do not set `HERDR_ENV=1` manually to bypass this check. (Setup and
`herdr integration` commands may run outside a managed pane.)

## Authoritative CLI Discovery

Run `herdr --skill` and consult the installed binary's command-group help
(`herdr workspace`, `herdr worktree`, `herdr agent`, `herdr pane`, `herdr api`)
rather than guessing flags. Never run bare `herdr` for flag discovery (it opens
the TUI).

## Single-Window Coordinator & Grouped Worktree Layout

1. **One main Herdr window rooted at the repository:**
   - The human opens one Herdr window in the root of the repository (`main`
     checkout), where the **Coordinator** agent runs.
   - Always pass **`--no-focus`** when creating/opening worktrees or tabs so you
     never steal the human's cursor, and never close or hijack panes driven by
     the human.
2. **Group parallel issue workers underneath the parent repo:**
   - Open each issue's isolated Git worktree attached to the parent repository
     workspace (`herdr worktree create --workspace <parent-ws-id> --branch
     <type>/<issue>-<slug> --label "#<issue> <slug>" --no-focus` or `herdr
     worktree open --workspace <parent-ws-id> ... --no-focus`) so all parallel
     workers appear **grouped together directly underneath the main repo** in
     Herdr's left sidebar.
   - If a task also touches external sibling repositories, group those into one
     dedicated `<topic>-upstream` workspace rather than cluttering the main
     repo's tab bar.
   - **Full-width panes, never narrow splits:** Terminal coding UIs need ~90+
     columns for readable diffs and status output; give each worker its own
     worktree entry rather than splitting one pane four ways.
3. **Zero-prompt startup & verified delivery:**
   - Start workers matching the user's harness (`--kind opencode`, `--kind agy`,
     or `--kind claude`). When launching `agy` on a brand-new worktree path,
     pass `-- --settings '{"trustedWorkspaces":["<abs-worktree-path>"]}'` so it
     starts without blocking on a workspace-trust prompt.
   - Never assume `agent_prompted` means the worker started working: read/wait
     until the worker transitions from `idle` -> `working`, and inspect `herdr
     pane read <pane> --source visible` before resending a timed-out prompt.
