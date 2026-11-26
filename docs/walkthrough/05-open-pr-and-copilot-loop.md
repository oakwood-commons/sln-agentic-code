---
title: "5. Open a Clean PR & Bound Copilot"
weight: 50
duration: "~4 MIN READ + TRY"
goal: "Pair your local harness with GitHub Copilot Code Review while protecting your review quota and preventing endless loops."
---

# Step 5: Open a Clean PR & Bound GitHub Copilot

## What Usually Goes Wrong

Using an independent external AI reviewer like **GitHub Copilot Code Review** on
pull requests is a great second pair of eyes until three defaults bite you:

1. **Noisy style nitpicks:** Without repo instructions, Copilot leaves subjective
   `"Consider renaming..."` comments and re-audits untouched files on every push.
2. **Surprise quota burn:** Each Copilot PR review uses premium request quota on
   your GitHub plan. If a repo enables automatic review on every `git push` and
   an agent pushes 6 WIP commits while fixing CI, **one PR burns 6 paid
   reviews**, mostly on broken code.
3. **Piecemeal review loops:** Fixing one comment triggers another round that
   spots one more comment in the same file, spinning for 5+ rounds.

## What to Do Instead

Even after Step 4's local review passes clean, **GitHub Copilot Code Review
almost always spots new edge cases worth knowing about** because it runs a
separate review pipeline in a true **cold-room context** (seeing only your final
PR diff and `REVIEW.md` on GitHub, with zero shared blind spots from your local
session or authoring model family; see
[The 5-Role Model Roster](/docs/reference/customize-and-credits.md#model-roster)).

Local `finish-issue` clears the noisy first ~80% of bugs and missing docs before
spending any CI runs or PR review quota, so Copilot can spend 100% of its Round
1 attention on the subtle remaining 20%, paired with **repo review files**,
**credit-saving triggers**, and **hard circuit breakers**:

<div class="bp-flow" role="region" aria-label="Bounded GitHub PR and Copilot review workflow">
  <div class="bp-flow-row cols-3">
    <div class="bp-flow-node">
      <span class="bp-flow-step">01 // LOCAL GATE</span>
      <span class="bp-flow-title"><code>finish-issue</code> Passes</span>
      <span class="bp-flow-desc">Clears ~80% of bugs &amp; missing docs locally (0 PR review quota used).</span>
    </div>
    <div class="bp-flow-arrow-h" aria-hidden="true"><svg viewBox="0 0 24 24" fill="none" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M5 12h14M13 6l6 6-6 6"/></svg></div>
    <div class="bp-flow-node">
      <span class="bp-flow-step">02 // OPEN PR</span>
      <span class="bp-flow-title"><code>/pr-create</code></span>
      <span class="bp-flow-desc">Opens ready-for-review PR and triggers Round 1 Copilot review once.</span>
    </div>
    <div class="bp-flow-arrow-h" aria-hidden="true"><svg viewBox="0 0 24 24" fill="none" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M5 12h14M13 6l6 6-6 6"/></svg></div>
    <div class="bp-flow-node">
      <span class="bp-flow-step">03 // CI FIRST</span>
      <span class="bp-flow-title">Wait for 100% Green CI</span>
      <span class="bp-flow-desc">Never re-request Copilot on broken builds.</span>
    </div>
  </div>
  <div class="bp-flow-arrow-v" aria-hidden="true"><svg viewBox="0 0 24 24" fill="none" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 5v14M6 13l6 6 6-6"/></svg></div>
  <div class="bp-flow-node is-center">
    <span class="bp-flow-step">04 // INSPECT ALL 4 PR SURFACES</span>
    <span class="bp-flow-title"><code>/pr-review</code></span>
    <span class="bp-flow-desc">Reads inline review threads, top-level review bodies, general PR comments, and CI checks.</span>
  </div>
  <div class="bp-flow-arrow-v" aria-hidden="true"><svg viewBox="0 0 24 24" fill="none" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 5v14M6 13l6 6 6-6"/></svg><span class="bp-flow-cond">Branch on review state</span><svg viewBox="0 0 24 24" fill="none" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 5v14M6 13l6 6 6-6"/></svg></div>
  <div class="bp-flow-row cols-2">
    <div class="bp-flow-node is-amber">
      <span class="bp-flow-step">CIRCUIT BREAKER // 3+ ROUNDS OR DESIGN CALL</span>
      <span class="bp-flow-title">STOP &amp; Ask Human</span>
      <span class="bp-flow-desc">Pause immediately on API/design trade-offs, or run <code>/review-buster</code> on a second model brand.</span>
    </div>
    <div class="bp-flow-node is-mint">
      <span class="bp-flow-step">MERGE READY // CLEAN + CI GREEN</span>
      <span class="bp-flow-title">Human Reviews Diff &amp; Clicks Merge</span>
      <span class="bp-flow-desc">Your single line-by-line code review happens here on the finished PR.</span>
    </div>
  </div>
</div>

---

## How the Pieces Work Together

### 1. Give Copilot a Repo Rulebook (`REVIEW.md`)

Check two files at the root of this repository:
- **[`.github/copilot-instructions.md`](https://github.com/oakwood-commons/sln-agentic-code/blob/main/.github/copilot-instructions.md):**
  GitHub Copilot reads this file automatically and follows its pointer straight
  to `REVIEW.md`.
- **[`REVIEW.md`](https://github.com/oakwood-commons/sln-agentic-code/blob/main/REVIEW.md):**
  Tells Copilot to skip style nits already handled by `gofmt`/`golangci-lint`,
  check callers outside the diff, report an entire class of bug at once, and
  **on re-reviews, inspect only what changed since its last review
  (`<last_reviewed_commit>..HEAD`)**.

### 2. Protect Your GitHub Copilot Quota

> [!WARNING]
> **Turn off `review_on_push`!** In GitHub Repository Settings -> Rules ->
> Rulesets (or automatically via the `scafctl` `configure-github` action), set
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

The `/pr-create`, `/pr-review`, `/pr-autopilot`, and `/review-buster` commands
enforce four non-negotiable stops:

| Safeguard | Rule |
|---|---|
| **Design trade-offs** | If a review comment questions an API shape, flag name, or behavior choice with two reasonable answers, **stop immediately and ask the human**. |
| **3-round circuit breaker** | Cap PR fix loops at **3 rounds maximum**. If Round 3 still raises new comments, stop or run `/review-buster`, which reads the *entire* review history with a different model family, diagnoses the shared root cause, and waits for your approval on one comprehensive plan. |
| **Quota-error check** | Never treat a Copilot comment saying `"unable to review"` or `"insufficient remaining quota"` as a clean pass. |
| **Human owns Merge (your one code review)** | `gh pr merge` is set to `"deny"` in `opencode.jsonc`. Because words, deterministic tests, live smoke runs, and `finish-issue` already did the heavy lifting in your terminal, **PR time is the one place you sit down to read the finished diff**, weigh Copilot's comments, and click **Merge** without needing a code editor. |

---

## Try the Commands on Your Own Fork

When working in your own repository or fork:

~~~text
# 1. Create a <type>/<issue#>-<slug> branch, open a PR, or keep the PR body synced
/pr-create feat: add optional link expiration and click stats to shorty

# 2a. One-shot interactive review: inspect all 4 PR feedback surfaces and choose what to do
/pr-review

# 2b. Hands-off convergence loop: wait for green CI, batch mechanical fixes, resolve nits,
#     and stop automatically when merge-ready, on a design call, or after 3 strikes
/pr-autopilot

# 3. If an external reviewer ever nitpicks one line per round, diagnose the root cause
/review-buster
~~~

*(Need to pause mid-task and pick up tomorrow in a fresh window? Run `/handoff`
to write a secret-redacted summary note to `/tmp/opencode/`.)*
