---
title: "3. Build in Small Tested Slices"
weight: 30
duration: "~5 MIN HANDS-ON"
goal: "Put a fakeable clock seam in place, then build expiration and click stats one failing test at a time."
---

# Step 3: Build in Small Tested Slices

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

## What We Do Instead

To make AI coding reliable, we pair the AI with **traditional
[deterministic](/docs/reference/glossary.md#deterministic) developer tools** --
TDD (`go test -race`), `golangci-lint`, and `gosec` -- so the AI has concrete
tools it can **test its ideas against and try again if they do not work, just
like a person would**.

Our `tdd` and `codebase-design` [skills](/docs/reference/glossary.md#skill)
keep that feedback loop fast and trustworthy through two simple rules:

- **Put a tiny [seam](/docs/reference/glossary.md#seam) at the clock:** Let
  `MemoryStore` (or `NewHandler`) accept a `now func() time.Time` function
  (defaulting to `time.Now`). Production code uses the real clock; unit tests
  advance a fake clock by 10 minutes in a microsecond with zero `time.Sleep`.
- **Work in vertical slices (`Red -> Green -> Refactor`):**
  1. **Red:** Write **one** failing test at the public HTTP handler and watch
     `go test` fail for the expected reason.
  2. **Green:** Write the smallest Go change that makes the whole suite pass.
  3. **Refactor:** Clean up naming or duplication while tests stay green, then
     move to the next slice.

{{< mermaid >}}
flowchart LR
    R["1. RED<br/>One failing httptest<br/>(e.g. expired link -> 410)"] --> G["2. GREEN<br/>Minimal Go code until<br/>go test -race passes"]
    G --> C["3. REFACTOR<br/>Tidy code while<br/>suite stays green"]
    C --> R
{{< /mermaid >}}

---

## Walk Through It

{{% steps %}}

1. **Build Slice 1: Expiring links (`410 Gone`)**

    Switch back to your `build` agent in OpenCode and prompt it to implement
    **only the first slice** using TDD:

    ~~~text
    Load the go, codebase-design, and tdd skills. Implement Slice 1 only
    (optional expires_in seconds on POST /links, returning 410 Gone on
    GET /{code} once expired, using an injectable clock seam instead of
    time.Sleep). Write the failing httptest first, confirm it fails, then make
    it pass.
    ~~~

    **What you should see:** OpenCode writes the failing test in
    `shortener_test.go` first, runs `go test` to confirm it fails (`got 302,
    want 410` or compile error on the new option), adds `ErrExpired` and the
    clock function, and re-runs `go test -race ./...` until green.

2. **Build Slice 2: Thread-safe click counter (`GET /links/{code}/stats`)**

    Now prompt OpenCode for the second slice -- including a concurrent test so
    Go's race detector proves click counting is safe under parallel traffic:

    ~~~text
    Now implement Slice 2 with TDD: increment the click count on each
    non-expired GET /{code} redirect, expose GET /links/{code}/stats returning
    code, url, clicks, and expired, and include a parallel goroutine test that
    passes under go test -race.
    ~~~

    **What you should see:** If `MemoryStore` tried to increment clicks under a
    read lock (`RLock`) or outside the mutex, `go test -race` immediately flags
    a `DATA RACE` and the agent fixes the lock before finishing the slice.

3. **Run the full portable Go CI check via `scafctl`**

    Verify formatting, `go vet`, `golangci-lint`, `gosec`, `-race` tests, the
    70%+ coverage threshold, and binary compilation in one command:

    ~~~bash
    scafctl run action go-ci
    ~~~

    **What you should see:** Every stage (`fmt-check`, `vet`, `lint`, `gosec`,
    `test`, `build`) finishes green in a few seconds.

{{% /steps %}}

---

## What If the AI Gets Stuck on a Bug?

Even good models occasionally hit tunnel vision on a tricky race condition or
compiler error. Two built-in safety nets prevent runaway token burn:

1. **Loop detection (`spiral-guard` [plugin](/docs/reference/glossary.md#plugin)):**
   Watches for repeated tool failures and nudges the agent to stop guessing,
   reproduce the bug cleanly (`diagnosing-bugs` skill), and test one hypothesis
   at a time.
2. **2-Strike Escalation & 3-Strike Stop:** Our personal policy tells fast coding
   models to escalate to a stronger reasoning [model](/docs/reference/glossary.md#model)
   after **2 failed tries** on the same check, and to **stop completely and ask
   you** after **3 failed tries**.

> [!NOTE]
> **Just reading along without OpenCode open?** You can inspect the starter code
> in [`examples/go-app/`](https://github.com/oakwood-commons/sln-agentic-code/tree/main/examples/go-app)
> anytime, and see how the same `scafctl` Go CI file runs in GitHub Actions and
> Tekton in [Portable Go CI & MCP Reference](/docs/reference/portable-ci-and-mcp.md).

Our tests pass -- next, let's see what happens when we run a free local code
review *before* committing in **Step 4**.
