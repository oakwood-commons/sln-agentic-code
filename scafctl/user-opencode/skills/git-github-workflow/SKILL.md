---
name: git-github-workflow
description: "Use when creating branches/worktrees, writing signed Conventional Commits, opening squash-merge PRs, syncing PR bodies, inspecting all four GitHub PR feedback surfaces, and running bounded GitHub Copilot review loops with cost guards and circuit breakers."
---

# Git, GitHub PR & Copilot Review Loop Workflow

## 1. Issues (Local, GitHub, or Jira MCP), Branching & Signed Commits

- **Issue tracking**: Track agreed scope in a GitHub Issue (`gh issue create` /
  `gh issue view <N>`), a Jira user story via MCP (`PROJ-123`), or a Markdown
  issue in git-ignored `temp/issues/<NNN>-<slug>.md` when practicing locally or
  working without upstream issue access.
- **Branch naming**: `<type>/<issue-id>-<slug>` (e.g., `feat/42-repo-readiness`,
  `feat/PROJ-123-expiring-links`, `feat/001-expiring-links`) or `<type>/<slug>`
  when no issue exists.
- **Isolated worktrees**: For parallel or disturbance-free work, create a
  sibling worktree (`git worktree add ../<repo>-worktrees/<slug> -b <branch>`)
  and start your coding session rooted inside that worktree so LSP and git state
  stay isolated.
- **Conventional commits + signing**: Format subjects as
  `<type>(<optional-scope>): <summary>` and sign commits with DCO sign-off when
  supported (`git commit -s -S`). Never run `git commit --amend` on published
  commits.

## 2. Squash-Merge PR Body Synchronization

Repositories configured with `squash_merge_commit_title=PR_TITLE` and
`squash_merge_commit_message=PR_BODY` build the permanent `main` commit from the
pull request title and body (discarding branch fixup commit messages):

1. Keep the PR title a valid Conventional Commit subject.
2. Whenever a new commit changes behavior, flags, or files on the branch, update
   the PR body in the same step (`gh pr edit <N> --title ... --body ...`) based
   on `git log <base>..HEAD` and `git diff --stat <base>...HEAD`.

## 3. Reading All Four GitHub PR Feedback Surfaces

Never assume `gh pr view` shows all review feedback. Always check all four:

1. **Inline review threads (GraphQL)** -- returns `id`, `isResolved`,
   `isOutdated`, file path, line, and comments:
   ~~~bash
   gh api graphql -f query='
     query($owner:String!, $repo:String!, $pr:Int!) {
       repository(owner:$owner, name:$repo) {
         pullRequest(number:$pr) {
           reviewThreads(first:100) {
             nodes {
               id isResolved isOutdated path line
               comments(first:20) { nodes { author { login } body createdAt } }
             }
           }
         }
       }
     }' -F owner="$OWNER" -F repo="$REPO" -F pr="$PR_NUM"
   ~~~
2. **Review bodies (`gh api repos/$OWNER/$REPO/pulls/$PR_NUM/reviews`)** --
   Copilot often places summary remarks or collapsed `<details>` findings in the
   top-level review body with zero inline threads.
3. **General PR comments (`gh api repos/$OWNER/$REPO/issues/$PR_NUM/comments`)**
   -- bot reports, coverage gates, and maintainer notes.
4. **CI check runs (`gh pr checks $PR_NUM`)** -- required build/test/security
   statuses on `HEAD`.

## 4. Reply & Resolve Are a Pair

When addressing or politely declining an inline review thread, reply inside that
thread AND resolve it via GraphQL (`resolveReviewThread`) in the same pass --
except when promising a code fix that has not been pushed yet. Re-query GraphQL
to confirm `isResolved: true`. Resolving a nit/decline thread does **not** push
a commit or spend another review round.

## 5. GitHub Copilot Review Cost Guardrails & Circuit Breakers

Repeatedly triggering AI PR reviews on every push burns GitHub Copilot premium
request quota fast and creates noisy review churn. Enforce these guardrails:

1. **Repo ruleset (`review_on_push: false`)**: Configure Copilot to review
   ready-for-review PRs once on open (`review_draft_pull_requests: false`),
   never automatically on every push.
2. **Local `/finish-issue` Before PR Open (0 PR quota)**: Catch bugs, missing
   tests, and `gosec` issues locally before Round 1 ever runs.
3. **CI-Green First**: **NEVER** request a Copilot re-review immediately after
   `git push`. Wait until 100% of GitHub CI checks on `HEAD` finish green so
   broken builds spend zero Copilot credits.
4. **Smart Post-Round-1 Re-Review Gate**:
   - **Pure rebase or docs-only delta after `No issues found`**: Do **not**
     re-request Copilot (0 Copilot review credits).
   - **Code fix addressing prior Copilot findings**: Batch all mechanical fixes
     into **one** locally reviewed commit + PR body sync, wait for green CI, and
     then request a Copilot re-review (`gh api -X POST
     repos/$OWNER/$REPO/pulls/$PR_NUM/requested_reviewers ...`) so Copilot
     verifies `<last_reviewed_sha>..HEAD` per `REVIEW.md` and clears `Changes
     recommended`.
5. **Quota-Error Guard**: If a Copilot review body contains `"unable to review"`
   or `"insufficient remaining quota"`, treat it as an unreviewed blocker and
   notify the human -- never treat it as `No issues found`.
6. **Three-Round Circuit Breaker & Human Stops**:
   - **Design questions (`Actionable-Design`)**: Stop immediately and ask the
     human whenever a comment raises an API, scope, or architecture choice.
   - **Max 3 fix rounds (`STALEMATE`)**: If a PR still has unresolved findings
     or failing checks after 3 fix rounds, **STOP** looping and surface it to
     the human (or run `/review-buster` for a root-cause reset).
   - **Never auto-merge**: Never run `gh pr merge` without explicit human
     instruction.
