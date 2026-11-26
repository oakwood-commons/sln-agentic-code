---
title: "4. Let AI Self-Review Before the PR"
weight: 40
duration: "~4 MIN HANDS-ON"
goal: "Have fresh AI reviewers catch subtle bugs and missing docs on your branch automatically so you only inspect code once on the PR."
---

# Step 4: Let AI Self-Review Before the PR

## What Usually Goes Wrong

Committing to a feature branch is harmless, and you do **not** want to sit in an
editor manually reviewing code before every branch commit. What wastes everyone's
time is having the AI open a pull request the instant `go test` turns green:

1. **Self-review blindness:** If the *same* AI chat that just wrote the code
   checks its own work, its conversation history is already primed to defend the
   choices it just made.
2. **Forgotten supporting files:** Tests pass on the Go code, and then GitHub
   Copilot (or you, during PR review) has to point out that `expires_in` and
   `GET /links/{code}/stats` were never documented in `examples/go-app/README.md`,
   costing an avoidable CI and review round.

## What to Do Instead

Validating everything it possibly can before asking for your final review is **an
AI job, not yours**.

Your standing `workflow-policy.md` instructs your harness to automatically run
the **`finish-issue` pre-PR gate** whenever it wraps up an issue to open a pull
request (and you can also trigger it anytime in plain English or via
`/finish-issue`). It runs static security checks (`gosec`), then dispatches **two
fresh, read-only [subagents](/docs/reference/glossary.md#subagent) in parallel**
over the branch diff:

<div class="bp-flow" role="region" aria-label="AI-driven finish-issue pre-PR self-review flow">
  <div class="bp-flow-node is-center">
    <span class="bp-flow-step">01 // BRANCH CHANGES</span>
    <span class="bp-flow-title">AI Runs Deterministic Checks First</span>
    <span class="bp-flow-desc">Verifies <code>go test -race</code>, <code>golangci-lint</code>, and <code>gosec</code> on your laptop (<code>$0</code> CPU).</span>
  </div>
  <div class="bp-flow-arrow-v" aria-hidden="true"><svg viewBox="0 0 24 24" fill="none" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 5v14M6 13l6 6 6-6"/></svg><span class="bp-flow-cond">Dispatched in parallel (fresh read-only contexts)</span><svg viewBox="0 0 24 24" fill="none" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 5v14M6 13l6 6 6-6"/></svg></div>
  <div class="bp-flow-row cols-2">
    <div class="bp-flow-node">
      <span class="bp-flow-step">02A // CODE-REVIEWER</span>
      <span class="bp-flow-title">Adversarial Code Review</span>
      <span class="bp-flow-desc">Checks bugs, data races, callers outside the diff, and <code>REVIEW.md</code> rules.</span>
    </div>
    <div class="bp-flow-node">
      <span class="bp-flow-step">02B // ARTIFACT-AUDITOR</span>
      <span class="bp-flow-title">Supporting Files Audit</span>
      <span class="bp-flow-desc">Checks missing tests, CLI <code>--help</code> flags, examples, and <code>README.md</code> updates.</span>
    </div>
  </div>
  <div class="bp-flow-arrow-v" aria-hidden="true"><svg viewBox="0 0 24 24" fill="none" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 5v14M6 13l6 6 6-6"/></svg><span class="bp-flow-cond">AI fixes confirmed gaps automatically</span><svg viewBox="0 0 24 24" fill="none" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 5v14M6 13l6 6 6-6"/></svg></div>
  <div class="bp-flow-node is-center is-mint">
    <span class="bp-flow-step">03 // READY FOR PR</span>
    <span class="bp-flow-title">One Combined Fix Pass + Final CI</span>
    <span class="bp-flow-desc">Pauses only if an architectural trade-off needs your call, saving CI runs, Copilot quota, and your review time.</span>
  </div>
</div>

Because neither subagent wrote the code and both are read-only (`edit: deny`),
the AI catches its own blind spots before pushing. You still have not opened a
code editor or read a line of Go locally; the AI only interrupts you here if one
of its reviewers surfaces a genuine **design or API decision** that belongs to a
human.

---

## Walk Through It

{{< steps >}}

1. **Watch the AI audit its own work (`finish-issue`)**

    On your own repo or fork, telling your harness *"wrap this up on a feature
    branch and open a PR"* runs `finish-issue` automatically before pushing.
    Because this walkthrough practices locally on `shorty` without opening a PR
    yet, trigger that pre-PR self-review directly:

    ~~~text
    Run the pre-PR finish-issue check on the shorty changes and fix anything mechanical.
    ~~~

    **What you should see:** The AI runs Go build/test/`gosec` checks,
    dispatches `code-reviewer` and `artifact-auditor` in parallel, and merges
    their findings into one list:
    - **`artifact-auditor`** flags `examples/go-app/README.md` as **Missing
      update** because `POST /links` now accepts `expires_in` and
      `GET /links/{code}/stats` is undocumented.
    - **`code-reviewer`** checks edge cases against
      [`REVIEW.md`](https://github.com/oakwood-commons/sln-agentic-code/blob/main/REVIEW.md)
      (such as rejecting negative `expires_in` values or large integer overflows,
      and returning `404` on `/stats` for unknown codes).

2. **Let the AI fix confirmed gaps & verify full CI**

    The AI fixes the mechanical findings and missing `README.md` docs in **one
    pass** (pausing only if something requires a real product or API decision
    from you) and re-runs the repository verification suite:

    ~~~shell
    scafctl run action ci
    ~~~

    **What you should see:** Go CI, `scafctl` solution tests, and the docs build
    all finish green, with both implementation *and* docs ready for a clean pull
    request.

{{< /steps >}}

> [!TIP]
> **Branch commits stay frictionless; publishing stays gated:** Committing on a
> feature branch cannot hurt `main`, so `opencode.jsonc` lets the AI stage and
> commit on feature branches freely (while blocking `git commit --amend`),
> prompts before pushing or opening a PR (`git push`, `gh pr create`), and
> strictly blocks merging (`gh pr merge`: `deny`).

Next, **Step 5** opens the pull request, pairs with GitHub Copilot Code Review
without burning quota, and brings you in once at the end for **final architectural
and code sign-off**.
