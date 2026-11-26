---
name: work-herd
description: "Use when asked to 'come up with a Herdr plan', build an issue dependency DAG from GitHub Issues, Jira tickets, or local issues, or orchestrate multiple autonomous workers in parallel across Herdr worktrees."
---

# Parallel Multi-Issue Orchestration with Herdr (`work-herd`)

Orchestrate a shortlist of issues across isolated Git worktrees in
[Herdr](https://herdr.dev) (`HERDR_ENV=1`). Composes with `herdr --skill` (live
CLI syntax), `github-issues`, `jira-issues`, `local-issues`, `planner`,
`finish-issue`, `git-github-workflow`, `pr-autopilot`, and `waiter`.

## Golden Rules

1. **You never merge (`gh pr merge` is forbidden).** The human reviews the
   finished PR diff once and clicks Merge; you process post-merge cleanup and
   refill the freed worker slot from the DAG.
2. **You coordinate; workers write code.** One worker per issue, one issue per
   worker, capped at the agreed concurrency limit (default: **4** active
   workers). If work has no issue yet, ask the human before provisioning (or
   draft one with `create-issue`).
3. **Never touch a human-driven pane, and always pass `--no-focus`** when
   creating or opening worktrees so you never steal the human's cursor.
4. **Verify ground truth rather than trusting a worker's `"DONE"` claim on
   faith.**
5. **Surface `blocked` design questions and `STALEMATE` stops to the human;
   never guess or auto-approve permission prompts.**

---

## Phase 1: Scope Candidate Issues, Then Build the Herdr Plan & DAG

Real repositories can have hundreds or thousands of open tickets. Never pull an
unbounded backlog and pick tickets at random:

- **Load the matching tracker skill:**
  - **GitHub Issues:** load `github-issues` (`gh issue list` filtered by label,
    milestone, or search query).
  - **Jira Tickets:** when the user says *"look at my Jira tickets to make a
    Herdr plan"*, load `jira-issues` and query open sprint or project tickets
    via JQL (`PROJ-101`, `PROJ-102`).
  - **Local Markdown Issues:** load `local-issues` and inspect `temp/issues/*.md`.
- **If the user has not narrowed the scope yet** (no explicit IDs, sprint,
  milestone, label, or topic area), query recent or sprint-scoped tickets around
  their stated goal, propose a focused shortlist (typically **4 to 8 tickets**),
  and confirm that list with the human first.
- **Once the candidate list is clear**, inspect each ticket against the codebase
  to map which packages and files it will touch, and write a local, git-ignored
  plan file at `temp/plans/<repo>-herd.md` containing:

1. **Issue Dependency DAG (`Ready Now` vs. `Blocked On`):**
   - Add edges for **logical dependencies** (Ticket B needs types or behavior
     introduced by Ticket A).
   - Add edges for **hotspot files**: any shared file more than one ticket will
     edit (a central registry, CLI root command, shared table, or config
     schema). Concurrent edits to the same lines cause rebase conflicts and
     extra review rounds, so either **serialize** them (`#101 -> #104` or
     `PROJ-101 -> PROJ-104`) or assign the shared edit to one ticket.
2. **Wave Schedule:** Group independent, conflict-free tickets into **Wave 1**
   (up to the concurrency cap), **Wave 2** (unblocked as Wave 1 PRs merge), and
   so on.
3. **Active Worker Registry:** Track `ws | pane | issue | branch | PR | status`
   for every provisioned worker.

If the user asked for the Herdr plan first, present the concise DAG, hotspot
notes, and Wave 1 lineup and wait for their go-ahead before spawning workers.

---

## Phase 2: Provisioning Each Ready Worker

First verify `test "${HERDR_ENV:-}" = 1` and consult `herdr --skill` / `herdr
worktree` / `herdr agent` for exact CLI syntax. For each ready issue in the
current wave:

1. **Create an isolated worktree off latest `origin/main` with `--no-focus`:**
   Open each branch worktree nested under the parent repo workspace (never
   cramming workers into narrow split panes, as TUIs need ~90+ columns).
2. **Start the worker agent** (`--kind opencode`, or `agy` / `claude` matching
   the user's harness) inside that worktree's pane.
3. **Deliver `worker-prompt.md` (in this skill directory)** with `__ISSUE__`,
   `__REPO__`, `__BRANCH__`, and `__AGENTNAME__` substituted (ensuring zero
   unreplaced `__` placeholders remain).
4. **Verify delivery (`agent_prompted` is NOT delivery):** Wait for the TUI
   ready state before sending, then confirm the worker transitions from `idle`
   -> `working` within ~30s (a prompt sent while a TUI is still booting can drop
   onto the welcome screen, leaving the pane `idle` looking finished). Re-send
   once if still on the welcome screen.
5. **Record the worker** in `temp/plans/<repo>-herd.md`.

---

## Phase 3: Event-Driven Monitoring (Never Sleep-Poll)

Never run `for ...; do sleep 60; herdr agent get ...; done` in your main
context, because per-minute status dumps bloat your context window.

1. Run `herdr api snapshot` to inspect all workers at once and handle any
   already `done` / `blocked` / `idle` worker immediately.
2. Block on state transitions with `herdr agent wait <worker> --until done
   --until blocked --until idle --timeout 600000` (10 min bounded timeout;
   returns in ~0s as soon as state changes).
3. When polling GitHub externally where no local event exists, delegate to the
   `waiter` subagent (one bounded shell loop returning a single summary line).

---

## Phase 4: Independent Ground-Truth Verification

Workers routinely misreport `"DONE"` (checking page 1 of unpaginated reviews,
missing unresolved GraphQL threads, or reporting before Copilot finishes).
Before notifying the human that a PR is **Merge-Ready**, independently verify:

1. **CI checks on current `headRefOid`:** All required checks finished green
   (`gh pr view <PR> --json headRefOid,reviewDecision,statusCheckRollup`).
2. **All inline threads resolved:** Query GraphQL `reviewThreads` (paginating
   if needed) and confirm `0` unresolved threads.
3. **Copilot review drained & covering `HEAD`:**
   - No pending Copilot request in `reviewRequests`, and the latest Copilot
     event in `repos/<owner>/<repo>/issues/<PR>/timeline --paginate` is
     `reviewed` (not an in-flight `review_requested` or `copilot_work_started`).
   - Latest valid (non-quota-error) `copilot-pull-request-reviewer` review
     covers `headRefOid` (or covers the prior round with only a locally
     reviewed surgical or rebase delta per `pr-autopilot`), with `0` unaddressed
     `<details>` review-body findings.

Always include **full clickable URLs** (`https://github.com/<owner>/<repo>/pull/<N>`)
when reporting **Merge-Ready (`DONE`)**, **Blocked / Design Question**, or
**Stalemate** to the human.

---

## Phase 5: Post-Merge Cleanup & DAG Slot Refill

Whenever the human merges a PR:

1. Fetch and pull `main`, then tear down that worker's worktree and pane cleanly
   (checking `git status` first so uncommitted work is never deleted).
2. Check other open worker PRs for merge conflicts created by the merge and
   have affected workers rebase if needed.
3. Mark the ticket merged in `temp/plans/<repo>-herd.md`, recompute newly
   unblocked tickets in the **DAG**, and provision the next ready ticket into
   the freed concurrency slot.
