---
title: "4. Check Locally Before Committing"
weight: 40
duration: "~4 MIN HANDS-ON"
goal: "Catch subtle bugs and forgotten documentation on your laptop for $0 before opening a pull request."
---

# Step 4: Check Locally Before Committing

## What Usually Goes Wrong

When `go test` turns green, it is tempting to immediately commit, push, and open
a pull request. Two annoying things usually happen next:

1. **Self-review blindness:** If you ask the *same* AI chat that just wrote your
   code *"Does this look good?"*, its conversation history is already primed to
   defend the choices it just made.
2. **Forgotten supporting files:** Your Go code works, so tests pass -- and then
   a teammate (or GitHub Copilot) points out that you forgot to document the new
   `expires_in` field and `GET /links/{code}/stats` endpoint in
   `examples/go-app/README.md`, costing another CI + review round.

## What We Do Instead

Before committing, we run one [slash command](/docs/reference/glossary.md#slash-command)
on our laptop: **`/finish-issue`**.

`/finish-issue` runs static security checks (`gosec`), then launches **two
fresh, read-only [subagents](/docs/reference/glossary.md#subagent) in parallel**
over your uncommitted `git diff`:

{{< mermaid >}}
flowchart LR
    D["Your Local git diff"] --> CR["code-reviewer subagent<br/>Checks bugs, races, security<br/>& REVIEW.md rules"]
    D --> AA["artifact-auditor subagent<br/>Checks missing tests, CLI help<br/>& README updates"]
    CR --> M["One Combined Fix List<br/>(Fixed locally for $0)"]
    AA --> M
{{< /mermaid >}}

Because neither subagent wrote the code -- and both are blocked from silently
editing files (`edit: deny`) -- you get an honest, unprimed checklist in seconds
at zero CI or GitHub Copilot cost.

---

## Walk Through It

{{% steps %}}

1. **Run `/finish-issue` in OpenCode**

    With your changes from Step 3 still uncommitted, run:

    ~~~text
    /finish-issue
    ~~~

    **What you should see:** OpenCode runs the Go build/test/`gosec` checks,
    dispatches `code-reviewer` and `artifact-auditor` at the same time, and
    merges their reports into a single list.

    Typical catches you will see right away on `shorty`:
    - **`artifact-auditor`** flags `examples/go-app/README.md` as **Missing
      update** because `POST /links` now accepts `expires_in` and the new
      `GET /links/{code}/stats` endpoint is not documented yet.
    - **`code-reviewer`** checks edge cases against the repo's
      [`REVIEW.md`](https://github.com/oakwood-commons/sln-agentic-code/blob/main/REVIEW.md)
      (for example, rejecting negative `expires_in` values or large integer
      overflows, and verifying `GET /links/{code}/stats` returns `404` on
      unknown codes).

2. **Watch one combined fix pass & final lint**

    `/finish-issue` fixes the confirmed items in **one pass** (stopping to ask
    you first if anything involves a real API or design choice) and re-runs the
    linter on the final diff.

    Verify the whole repository from your terminal one more time:

    ~~~bash
    scafctl run action ci
    ~~~

    **What you should see:** Go CI, `scafctl` solution tests, and the docs build
    all finish green, and `git diff` shows clean code *plus* updated docs.

{{% /steps %}}

> [!TIP]
> **Notice how `/finish-issue` stays silent on git commits:** Even after every
> check passes, `opencode.jsonc` keeps `git commit` and `git push` behind an
> approval prompt (`"ask"`), and blocks `git commit --amend` (`"deny"`). You
> always decide when code is ready to publish.

Now that our local diff is clean, **Step 5** shows how we open a pull request
and pair with GitHub Copilot Code Review without burning credits or getting
stuck in review loops.
