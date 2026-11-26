---
title: "3. Build One Piece at a Time"
weight: 30
duration: "~5 MIN HANDS-ON"
goal: "Put a fakeable clock seam in place, then build link expiration and click stats one failing test at a time."
---

# Step 3: Build One Piece at a Time

## What Usually Goes Wrong

An AI model on its own is probabilistic -- it predicts code that *looks* right,
not code that is guaranteed to work. Left on autopilot without real feedback
tools, AI coding models make two classic mistakes in Go:

1. **They write everything in one giant batch** -- 200 lines of handler code and
   10 tests at once. When tests pass on the first try, you never saw them fail,
   so you do not know if they actually test anything.
2. **They hardcode `time.Now()` inside business logic** -- forcing expiration
   tests to call `time.Sleep(2 * time.Second)`, which makes test suites slow and
   flaky.

## What to Do Instead

Pair the AI with **good tests** and **[deterministic](/docs/reference/glossary.md#deterministic)
developer tools** (`go test -race`, `golangci-lint`, `gosec`). In an AI workflow,
good tests -- ones that check real observable behavior at public interfaces
rather than fragile private internals -- do **two essential jobs at once**:

1. **They help the AI write better new code:** Writing one failing test at the
   HTTP boundary first forces clean interfaces and gives the AI an unmistakable
   pass/fail target to iterate against instead of guessing.
2. **They give you confidence AI won't break existing code:** Every time the AI
   edits `shortener.go`, the *entire* existing suite runs alongside the new test
   -- proving non-expiring links, validation errors, and `404` responses still
   work untouched. Today's new tests become tomorrow's safety net for the next
   feature.

The starter `tdd` and `codebase-design`
[skills](/docs/reference/glossary.md#skill) keep that test suite fast and
trustworthy through two simple rules:

- **Put a tiny [seam](/docs/reference/glossary.md#seam) at the clock:** Let
  `MemoryStore` (or `NewHandler`) accept a `now func() time.Time` function
  (defaulting to `time.Now`). Production code uses the real clock; unit tests
  advance a fake clock by 10 minutes in a microsecond with zero `time.Sleep`.
- **Build one behavior at a time (`Red -> Green -> Refactor`):**
  1. **Red:** Write **one** failing test at the public HTTP handler and watch
     `go test` fail for the expected reason.
  2. **Green:** Write the smallest Go change that makes the whole suite pass.
  3. **Refactor:** Clean up naming or duplication while tests stay green, then
     move to the next behavior.

<div class="bp-flow" role="region" aria-label="Red-Green-Refactor test-first loop">
  <div class="bp-flow-row cols-3">
    <div class="bp-flow-node is-coral">
      <span class="bp-flow-step">01 // RED</span>
      <span class="bp-flow-title">One Failing Test</span>
      <span class="bp-flow-desc">Write one <code>httptest</code> at the public handler and watch it fail (e.g. <code>got 302, want 410</code>).</span>
    </div>
    <div class="bp-flow-arrow-h" aria-hidden="true"><svg viewBox="0 0 24 24" fill="none" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M5 12h14M13 6l6 6-6 6"/></svg></div>
    <div class="bp-flow-node is-mint">
      <span class="bp-flow-step">02 // GREEN</span>
      <span class="bp-flow-title">Minimal Go Code</span>
      <span class="bp-flow-desc">Write the smallest implementation change until <code>go test -race ./...</code> passes.</span>
    </div>
    <div class="bp-flow-arrow-h" aria-hidden="true"><svg viewBox="0 0 24 24" fill="none" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M5 12h14M13 6l6 6-6 6"/></svg></div>
    <div class="bp-flow-node">
      <span class="bp-flow-step">03 // REFACTOR</span>
      <span class="bp-flow-title">Tidy While Green</span>
      <span class="bp-flow-desc">Clean up naming or duplication while the full suite stays green.</span>
    </div>
  </div>
  <div class="bp-flow-loop"><svg viewBox="0 0 24 24" fill="none" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><polyline points="1 4 1 10 7 10"/><path d="M3.51 15a9 9 0 102.13-9.36L1 10"/></svg><span>Repeat from 01 // RED for the next piece of the plan</span></div>
</div>

---

## Walk Through It

{{< steps >}}

1. **Part 1: Add expiring links (`410 Gone`)**

    Tell OpenCode in plain English to implement **expiring links first** from
    the plan, test-first:

    ~~~text
    Let's implement expiring links first from the plan, writing the failing httptest before the code.
    ~~~

    **What you should see:** OpenCode's inner monologue loads the `go`,
    `codebase-design`, and `tdd` skills from your prose, writes the failing test
    in `shortener_test.go` first, runs `go test` to confirm it fails (`got 302,
    want 410` or compile error on the new option), adds `ErrExpired` and the
    clock function, and re-runs `go test -race ./...` until green.

2. **Part 2: Add the thread-safe click counter (`GET /links/{code}/stats`)**

    Now ask for the click counter -- including a concurrent test so Go's race
    detector proves click counting is safe under parallel traffic:

    ~~~text
    Now implement the click stats endpoint from the plan, again test-first, including a parallel goroutine test under go test -race.
    ~~~

    **What you should see:** If `MemoryStore` tried to increment clicks under a
    read lock (`RLock`) or outside the mutex, `go test -race` immediately flags
    a `DATA RACE` and the agent fixes the lock before finishing.

3. **Run the full portable Go CI check & smoke-test `shorty` live**

    Verify formatting, `go vet`, `golangci-lint`, `gosec`, `-race` tests, the
    70%+ coverage threshold, and binary compilation in one command:

    ~~~shell
    scafctl run action go-ci
    ~~~

    **What you should see:** Every stage (`fmt-check`, `vet`, `lint`, `gosec`,
    `test`, `build`) finishes green in a few seconds.

    You can also verify behavior in words without opening an editor -- ask
    OpenCode directly:

    ~~~text
    Start shorty briefly, create a link with expires_in=2, show me a working redirect and /stats count, then wait 2 seconds and confirm GET /{code} returns 410 Gone while /stats still reports clicks=1 and expired=true.
    ~~~

    Between deterministic CI checks and watching the live binary respond as
    expected, you know *what* you got matches *what* you asked for -- before
    reading a single line of Go yourself.

{{< /steps >}}

---

## What If the AI Gets Stuck on a Bug?

Even good models occasionally hit tunnel vision on a tricky race condition or
compiler error. Two built-in safety nets prevent runaway token burn:

1. **Structured bug triage (`diagnosing-bugs` [skill](/docs/reference/glossary.md#skill)):**
   Cues the agent to stop guessing, build a fast deterministic reproducer first,
   test one falsifiable hypothesis at a time, and lock in a regression test at
   the public seam.
2. **2-Strike Escalation & 3-Strike Stop:** Your `workflow-policy.md` tells
   faster coding models to escalate to a stronger reasoning
   [model](/docs/reference/glossary.md#model) after **2 failed tries** on the
   same check, and to **stop completely and ask you** after **3 failed tries**.
   *(Want automatic runtime loop detection too? See `spiral-guard` under
   [Recommended Optional Plugins](/docs/reference/customize-and-credits.md#recommended-plugins).)*

> [!NOTE]
> **Just reading along without OpenCode open?** You can inspect the starter code
> in [`examples/go-app/`](https://github.com/oakwood-commons/sln-agentic-code/tree/main/examples/go-app)
> anytime, and see how the same `scafctl` Go CI file runs in GitHub Actions and
> Tekton in [Portable Go CI & MCP Reference](/docs/reference/portable-ci-and-mcp.md).

Tests pass and the live server works -- next, **Step 4** has fresh AI subagents
audit the branch for subtle bugs and forgotten docs *before* opening a pull
request.
