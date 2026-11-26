---
title: "2. Ask & Plan First"
weight: 20
duration: "~4 MIN HANDS-ON"
goal: "Describe what you want and why in plain English, capture a two-layer issue, and skeptically verify a plan against the code."
---

# Step 2: Ask & Plan First

## What Usually Goes Wrong

People often call bad AI code *"hallucination"*, when in reality the AI is simply
giving the best answer it can from incomplete
[context](/docs/reference/glossary.md#context). Everyday prompts usually fail in
two opposite ways:

1. **Jumping straight to code with zero intent (*"Add expiring links and
   stats"*):** With no *why* and no plan step, the AI silently guesses product
   rules (`404` vs. `410 Gone`, default TTL, expired stats) and dumps 250 lines
   of Go you did not ask for.
2. **Prescribing *how* instead of *what* (*"Run a 1-second `time.Ticker`
   goroutine to delete expired map entries..."*):** AI models try hard to give
   you what they think you asked for. Over-specifying the mechanism locks the AI
   into clunky concurrency, and deleting expired entries from the map silently
   wipes out final promo click counts.

## Plain English + Auto-Loaded Skills

You do not need a code editor open or a paragraph of meta-instructions on every
prompt. Each installed [skill](/docs/reference/glossary.md#skill) has a short
`description` header exposed to the model on every turn alongside your standing
`workflow-policy.md`:

- **Skip repeating repo facts:** Your harness already reads `AGENTS.md`,
  `go.mod`, and `shortener.go`, so it already knows `shorty` is in-memory and
  stdlib-only. Only mention stack choices when introducing something **not** in
  the repo yet (for example, *"and let's use `sqlite`"*).
- **Supply *why*, leave the API & *how* open:** State what you want and the use
  case (*"for temporary promo links"*). Let `planner`, `create-issue`, and
  `stepwise-interview` inspect existing routes, propose the cleanest API and
  mechanism, and ask one plain-English question at a time for genuine edge
  cases.
- **Probabilistic prose vs. deterministic slash commands:** Matching skills from
  plain English is a model judgment call. Strong reasoning models load
  `planner` and `create-issue` reliably from natural prose, while smaller or
  faster models can occasionally skip a tool call. Use
  [slash commands](/docs/reference/glossary.md#slash-command) (`/plan`,
  `/create-issue`, `/wait-what`, `/finish-issue`) whenever you want a **100%
  deterministic override**.

<div class="bp-compare">
  <div class="bp-compare-card is-bad">
    <div class="bp-compare-head">
      <svg viewBox="0 0 24 24" fill="none" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/></svg>
      <span class="bp-compare-tag">Avoid</span>
      <span class="bp-compare-title">No intent, straight to code (forces silent guesses)</span>
    </div>
    <div class="bp-compare-prompt">"Add expiring links and stats."</div>
    <p class="bp-compare-outcome"><strong>Result:</strong> No use case and no plan gate. The AI guesses three product rules and writes 250 lines of Go immediately.</p>
  </div>
  <div class="bp-compare-card is-warn">
    <div class="bp-compare-head">
      <svg viewBox="0 0 24 24" fill="none" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
      <span class="bp-compare-tag">Avoid</span>
      <span class="bp-compare-title">Prescribing how (locks in a flawed mechanism)</span>
    </div>
    <div class="bp-compare-prompt">"Add expires_in int64, run a 1s time.Ticker goroutine to delete expired map entries, and count clicks in a sync.Map."</div>
    <p class="bp-compare-outcome"><strong>Result:</strong> AI builds your ticker goroutine instead of a simpler lazy check, and deleting expired entries wipes out final click totals.</p>
  </div>
  <div class="bp-compare-card is-good">
    <div class="bp-compare-head">
      <svg viewBox="0 0 24 24" fill="none" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M22 11.08V12a10 10 0 11-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
      <span class="bp-compare-tag">Do This</span>
      <span class="bp-compare-title">Goal + why in plain English (skills handle the rest)</span>
    </div>
    <div class="bp-compare-prompt">"Let's create an issue for shorty to support optional expiring links (for temporary promo links) and a way to check click counts."</div>
    <p class="bp-compare-outcome"><strong>Result:</strong> <code>create-issue</code> and <code>planner</code> inspect <code>shortener.go</code>, clarify edge cases one at a time, and draft a two-layer issue for your review.</p>
  </div>
</div>

---

## Walk Through It

{{< steps >}}

1. **Launch your coding harness inside the repository**

    From the root of `sln-agentic-code`, start OpenCode (or `agy`):

    ~~~shell
    opencode
    ~~~

2. **Describe what you want & review the two-layer issue draft**

    Paste this plain-English prompt:

    ~~~text
    Let's create an issue for shorty to support optional expiring links (for temporary promo links) and a way to check click counts.
    ~~~

    **What happens under the hood:**
    - Your harness loads `create-issue`, `stepwise-interview`, and
      `local-issues`, inspects `examples/go-app/`, and sees `MemoryStore` + Go
      stdlib-only on its own.
    - Looking at existing `POST /links` and `GET /{code}`, it **proposes** a
      clean API (`expires_in` seconds on `POST /links` and
      `GET /links/{code}/stats` so stats never collide with 6-character codes)
      and asks **one plain-English question at a time** for the edge cases:
      - **Expired redirect status:** `410 Gone` (`ErrExpired`), keeping `404 Not
        Found` for unknown codes.
      - **Omitted `expires_in`:** Optional (`0` or omitted lives forever;
        negative numbers return `400 Bad Request`).
      - **Stats after expiry:** Because promo creators care about final campaign
        results, `GET /links/{code}/stats` still returns `200 OK` with
        `{"code":"...","url":"...","clicks":N,"expired":true}`.
    - Before saving, `create-issue` previews a **two-half draft** in chat:
      1. **Top half in Plain English (for you):** *Summary & Why*, *Agreed
         Behavior & Edge Cases*, and *Acceptance Criteria*, so you can verify
         intent in 20 seconds without reading Go code.
      2. **Bottom half with Technical Context (for the AI):** relevant packages,
         existing test patterns, constraints, and verification seams, giving
         Step 3's coding session grounded context without locking it into brittle
         line-by-line code prescriptions.
    - Once you confirm, `local-issues` saves
      `temp/issues/001-expiring-links-and-stats.md` so you can practice locally
      without needing write access to the upstream GitHub repo.

    > [!NOTE]
    > **Same issue format across Local files, GitHub Issues, and Jira:** The
    > starter bundle pairs `create-issue` with three tracker skills
    > (`local-issues`, `github-issues`, and `jira-issues` via
    > [Atlassian MCP](/docs/reference/portable-ci-and-mcp.md#jira-mcp)). On real
    > projects you can say *"create a GitHub issue for..."* or *"create a Jira
    > ticket in PROJ for..."* (or configure your repo's default in `AGENTS.md`
    > via `init-repo`).

    > [!TIP]
    > **Confused by a term?** Reply `"what do you mean?"` or `"I don't follow"`.
    > Your standing policy and the `wait-what` skill cue the AI to drop jargon
    > and re-explain in plain sentences (or type `/wait-what` to force it
    > explicitly).

3. **Start a fresh session & skeptically plan from the issue**

    Start a clean session (`/new`, or exit and launch `opencode` again) so
    brainstorming chatter is cleared, then point `planner` at the issue:

    ~~~text
    Have a look at issue 001 and come up with a detailed plan to implement it.
    ~~~

    *(On a GitHub or Jira project, use the same prompt with `issue #42` or
    `PROJ-123`.)*

    **What you should see:** Even when an AI helped draft the issue, `planner`
    treats it as a **hypothesis**: verifying every assumption against
    `shortener.go` (`file:line`), recommending a lazy expiration check via an
    injectable clock seam instead of a background goroutine, and breaking the
    work into two small testable steps **without modifying a single file yet**.

{{< /steps >}}

In **Step 3**, you will build those two steps and verify them with tests,
linters, and a live local run, still without opening a code editor.
