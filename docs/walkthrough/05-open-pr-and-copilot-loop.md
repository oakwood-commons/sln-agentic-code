---
title: "5. Open a Clean PR & Bound Copilot"
weight: 50
duration: "~4 MIN READ + TRY"
goal: "Pair your local harness with GitHub Copilot Code Review while protecting your review quota and preventing endless loops."
---

# Step 5: Open a Clean PR & Bound GitHub Copilot

## What Usually Goes Wrong

Using an independent external AI reviewer like **GitHub Copilot Code Review** on
pull requests is a great second pair of eyes -- until three defaults bite you:

1. **Noisy style nitpicks:** Without repo instructions, Copilot leaves subjective
   `"Consider renaming..."` comments and re-audits untouched files on every push.
2. **Surprise quota burn:** Each Copilot PR review uses premium request quota on
   your GitHub plan. If a repo enables automatic review on every `git push` and
   an agent pushes 6 WIP commits while fixing CI, **one PR burns 6 paid
   reviews** -- mostly on broken code.
3. **Piecemeal review loops:** Fixing one comment triggers another round that
   spots one more comment in the same file, spinning for 5+ rounds.

## What We Do Instead

We pair our local OpenCode harness with GitHub Copilot using **repo review
files**, **credit-saving triggers**, and **hard circuit breakers**:

{{< mermaid >}}
flowchart LR
    L["Local /finish-issue<br/>($0 pre-PR gate)"] --> P["/pr-create<br/>Opens PR (Round 1 Copilot)"]
    P --> CI["Wait for 100% Green CI<br/>(Never review broken builds)"]
    CI --> R["/pr-review<br/>Reads all 4 PR surfaces"]
    R -->|"3+ rounds or design call"| H["STOP & Ask Human<br/>(or /review-buster)"]
    R -->|"Clean + CI green"| M["Human Clicks Merge"]
{{< /mermaid >}}

---

## How the Pieces Work Together

### 1. Give Copilot a Repo Rulebook (`REVIEW.md`)

Check two files at the root of this repository:
- **[`.github/copilot-instructions.md`](https://github.com/oakwood-commons/sln-agentic-code/blob/main/.github/copilot-instructions.md)** --
  GitHub Copilot reads this file automatically and follows its pointer straight
  to `REVIEW.md`.
- **[`REVIEW.md`](https://github.com/oakwood-commons/sln-agentic-code/blob/main/REVIEW.md)** --
  Tells Copilot to skip style nits already handled by `gofmt`/`golangci-lint`,
  check callers outside the diff, report an entire class of bug at once, and
  **on re-reviews, inspect only what changed since its last review
  (`<last_reviewed_commit>..HEAD`)**.

### 2. Protect Your GitHub Copilot Quota

> [!WARNING]
> **Turn off `review_on_push`!** In GitHub Repository Settings -> Rules ->
> Rulesets (or automatically via our `scafctl` `configure-github` action), set
> **`review_on_push: false`** and **`review_draft_pull_requests: false`**.

With those two settings in place, follow three simple habits:
1. **Let Round 1 run on a clean PR:** Because you ran `/finish-issue` locally in
   Step 4, Round 1 often comes back `No issues found` on the first pass.
2. **CI-Green First:** If you push a fix, **wait until every GitHub Actions CI
   check turns green** before requesting a Copilot re-review. Broken builds
   spend zero review credits.
3. **Skip re-reviews on pure rebases or docs tweaks:** Once Copilot reports `No
   issues found`, do not re-request it after a plain `main` rebase. Only request
   a re-review when you pushed code fixes to clear prior open findings.

### 3. Enforce Hard Circuit Breakers & Human Stops

> [!CAUTION]
> **Never let an AI loop on a pull request unattended without an exit condition.**

Our `/pr-create`, `/pr-review`, and `/review-buster` commands enforce four
non-negotiable stops:

| Safeguard | Rule |
|---|---|
| **Design trade-offs** | If a review comment questions an API shape, flag name, or behavior choice with two reasonable answers, **stop immediately and ask the human**. |
| **3-round circuit breaker** | Cap PR fix loops at **3 rounds maximum**. If Round 3 still raises new comments, stop -- or run `/review-buster`, which reads the *entire* review history with a different model family, diagnoses the shared root cause, and waits for your approval on one comprehensive plan. |
| **Quota-error check** | Never treat a Copilot comment saying `"unable to review"` or `"insufficient remaining quota"` as a clean pass. |
| **Human owns Merge** | `gh pr merge` is set to `"deny"` in `opencode.jsonc`. Agents get the PR green and resolved; you click **Merge**. |

---

## Try the Commands on Your Own Fork

When working in your own repository or fork:

~~~text
# 1. Create a <type>/<issue#>-<slug> branch, open a PR, or keep the PR body synced
/pr-create feat: add optional link expiration and click stats to shorty

# 2. Inspect all 4 GitHub PR feedback surfaces (inline threads, review bodies,
#    general comments, and CI checks) and choose what to do next
/pr-review

# 3. If an external reviewer ever nitpicks one line per round, diagnose the root cause
/review-buster
~~~

*(Need to pause mid-task and pick up tomorrow in a fresh window? Run `/handoff`
to write a secret-redacted summary note to `/tmp/opencode/`.)*
