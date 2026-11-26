---
name: pr-autopilot
description: "Autonomous GitHub PR review and CI convergence loop: polls all four PR feedback surfaces (GraphQL reviewThreads, review bodies including collapsed Copilot details, general comments, and CI checks), batches mechanical fixes into one push + PR body sync, replies and resolves nits without pushing, enforces the Copilot credit gate, and stops ONLY when merge-ready, a design call is needed, or the 3-strike cap is hit. Never merges."
---

# PR Autopilot -- Autonomous PR & Copilot Convergence Loop

Drive an open GitHub pull request to a clean, merge-ready state without wasting
GitHub Copilot review quota or spinning in endless fix loops. Composes with
`git-github-workflow` and `waiter`.

**Hard rule: NEVER run `gh pr merge` or dismiss reviews.** Your only terminal
outcomes are **Merge-Ready (`DONE`)**, **Pause for Human Decision**, or
**Stalemate (3-strike / safety cap)**.

## Phase 0 -- Target & Reviewer Probe

1. Resolve `<owner>/<repo>` and `<PR>` (`gh repo view --json owner,name` and
   `gh pr view --json number,url,isDraft,headRefName,headRefOid,reviewDecision,reviewRequests,reviews`).
2. Verify GraphQL `reviewThreads` enumeration and `resolveReviewThread`
   capability are available before starting.
3. Probe whether GitHub Copilot reviews this repo (`copilot-pull-request-reviewer`
   in `reviews`, `reviewRequests`, or recent repo PRs). If the PR is a draft
   with no bot review, converge on CI and human threads only and report draft
   status at exit (never run `gh pr ready` unilaterally). If Copilot is active
   on the repo but unrequested on a ready PR, request once (`gh pr edit <PR>
   --add-reviewer @copilot`).
4. Default bounds: `base_interval = 60s`, `max_iterations = 60`, `max_runtime =
   2h` (or caller-supplied bounds such as `30` / `45m` inside a Herdr worker).

## Phase 1 -- Read All Four PR Surfaces (Every Cycle)

1. Check `gh api rate_limit` first; back off cleanly if `< 15%` remaining.
2. Fetch all four surfaces (always `--paginate` REST list endpoints):
   - **Inline review threads (GraphQL `reviewThreads`)** -- `id`, `isResolved`,
     `isOutdated`, `path`, `line`, and comments.
   - **Review bodies (`gh pr view --json reviews,headRefOid,reviewDecision`)**
     -- parse collapsed `<details>` findings (`Suppressed comments`, `Previously
     missed`, `needs a closer look`) and `file:line` items in the latest valid
     `copilot-pull-request-reviewer` review. **Quota-error guard:** NEVER treat
     a review body containing `"unable to review"` or `"insufficient remaining
     quota"` as a clean review; latch `rerequest_blocked = true` and surface it
     if Round 1 never ran.
   - **General PR comments (`repos/$OWNER/$REPO/issues/$PR/comments --paginate`)**
     -- CI bots, policy checks, maintainer notes.
   - **CI check runs (`gh pr checks <PR>`)** -- track unchanged pending checks
     across cycles to detect stuck non-moving gates.

## Phase 2 -- Classify Every Item Before Acting

| Category | Meaning | Action |
|---|---|---|
| **Nit / Question / Already addressed / Disagree / Outdated** | Non-blocking wording, answered question, or stale note | **Auto:** reply with clear reasoning + resolve via GraphQL (**no push**) |
| **Actionable-Mechanical** | Bounded bug, missing nil/error check, missing test, security scanner fix | **Auto:** fix locally, batch into **one** verified commit + push + PR body sync |
| **Actionable-Design / Ambiguous** | Public API shape, breaking change, repo/policy setting, or 2+ valid designs | **STOP & Ask Human** (after handling routine nits/mechanicals in this fetch) |

**3-strike circuit breaker:** Stop and escalate to the human (or suggest
`/review-buster`) if the same finding/check fails after **3** fix attempts or a
fix introduces more errors than it resolves. Security scanner findings are
always actionable (never dismissed as nits).

## Phase 3 -- Batch Push & Copilot Credit Gate (`review_on_push: false`)

1. Apply all `Actionable-Mechanical` fixes in the cycle first, run the repo's
   build/race-test/lint/`gosec` checks + local `code-reviewer` check on the
   fix, sign-commit **once** (`git commit -s -S`), `git fetch && git rebase`
   onto default branch, push **once**, sync the PR title/body (`git-github-workflow`),
   and reply + resolve fixed threads.
2. **Never request Copilot immediately on `git push`:**
   - Wait (via the `waiter` subagent or one bounded shell poll -- never
     per-minute context spam) until **100% of required CI checks on `HEAD`
     finish green** so broken builds spend zero Copilot credits.
   - **Pure rebase / lockfile / docs-only delta after `No issues found`:** Do
     **not** re-request Copilot; record `minor_delta_verified_head = HEAD`.
   - **Code fix addressing prior findings:** Once CI is 100% green, request one
     Copilot re-review (`POST repos/$OWNER/$REPO/pulls/$PR/requested_reviewers`
     with `reviewers[]=copilot-pull-request-reviewer[bot]`) so Copilot verifies
     `<last_reviewed_oid>..HEAD` per `REVIEW.md` and clears `Changes
     recommended`.

## Phase 4 -- Drain Check & Terminal Report

Before declaring **Merge-Ready (`DONE`)**, refresh live PR state and confirm:
- Required CI checks on `HEAD` are finished and green.
- `0` unresolved GraphQL review threads and `0` open Copilot review-body
  findings.
- Copilot's timeline (`repos/$OWNER/$REPO/issues/$PR/timeline --paginate`) shows
  no in-flight `review_requested` or `copilot_work_started` after its latest
  `reviewed` event, and the latest valid review covers `HEAD` (or covers the
  prior round + locally verified rebase/minor delta).

Emit a concise terminal summary with the full PR URL, per-finding dispositions
(`<id> -> fixed <sha> | declined | answered`), CI status, and bot coverage `[review <oid> | HEAD <headOid>]` -- and let the human review and click Merge.
