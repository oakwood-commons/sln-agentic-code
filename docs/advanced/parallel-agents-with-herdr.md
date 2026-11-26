---
title: "Parallel Issues with Herdr"
weight: 10
---

# Scaling Up: Parallel Issues with Herdr, DAG Plans & Autonomous Copilot

Once Steps 1--5 click -- scoping in words, proving behavior with deterministic
CI, and reading code once at PR time -- watching **one** terminal work through
**one** issue at a time quickly becomes the bottleneck.

Running 3 or 4 issues in parallel without chaos comes down to a simple layout:
**you open one main Herdr window at the root of your repository**, and a
Coordinator session inside that window spins up **isolated worker worktrees
grouped together right underneath your repo in the left sidebar** -- each
covering a separate issue that can safely run in parallel.

Four bundled pieces in `agentic-code` (`herdr`, `work-herd`, `pr-autopilot`,
and `waiter`) power that loop:

1. **One Main Repo Window + Grouped Workers (`herdr` skill)** -- you stay in
   one Herdr window rooted in your main checkout; parallel worker worktrees nest
   directly underneath that parent repo in the sidebar (`working`, `blocked`,
   `done`, `idle`) with `--no-focus` so nothing steals your cursor or opens
   extra desktop windows.
2. **Herdr DAG Planning & Ground-Truth Coordinator (`work-herd` skill)** --
   inspects your shortlisted issues, builds an **issue dependency & hotspot
   [DAG](/docs/reference/glossary.md#dag)** (dependency order with no collisions)
   before spawning workers, double-checks workers' `"DONE"` claims, and refills
   open slots as you merge PRs.
3. **Autonomous Copilot Convergence (`pr-autopilot` + `waiter`)** -- each
   worker drives its own PR through GitHub Actions CI and GitHub Copilot
   reviews to merge-ready without human babysitting (and without ever running
   `gh pr merge`).

~~~tree
Herdr Window (Left Sidebar)
└── my-service                  [Coordinator -- root checkout on main]
    ├── #101 expiring-links     [working]   feat/101-expiring-links (worktree)
    ├── #102 portable-go-ci     [working]   feat/102-portable-go-ci (worktree)
    └── #103 powershell-tabs    [done]      feat/103-powershell-tabs (worktree)
~~~

<div class="bp-flow" role="region" aria-label="Herdr DAG planning, parallel worktrees, pr-autopilot, and 3-interrupt workflow">
  <div class="bp-flow-node is-center">
    <span class="bp-flow-step">01 // MAIN REPO WINDOW (COORDINATOR + work-herd SKILL)</span>
    <span class="bp-flow-title">Shortlist Sprint / Topic Issues -&gt; <code>temp/plans/&lt;repo&gt;-herd.md</code></span>
    <span class="bp-flow-desc">Filters backlog around your focus area, then maps dependencies + hotspot files into a conflict-free Issue DAG.</span>
  </div>
  <div class="bp-flow-arrow-v" aria-hidden="true"><svg viewBox="0 0 24 24" fill="none" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 5v14M6 13l6 6 6-6"/></svg><span class="bp-flow-cond">Spins up grouped child worktrees under the repo (--no-focus)</span><svg viewBox="0 0 24 24" fill="none" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 5v14M6 13l6 6 6-6"/></svg></div>
  <div class="bp-flow-row cols-3-plain">
    <div class="bp-flow-node">
      <span class="bp-flow-step">WORKER 01 // READY</span>
      <span class="bp-flow-title">Worktree <code>#101</code></span>
      <span class="bp-flow-desc">Plan -&gt; TDD -&gt; <code>finish-issue</code> -&gt; PR -&gt; <strong><code>pr-autopilot</code></strong></span>
    </div>
    <div class="bp-flow-node">
      <span class="bp-flow-step">WORKER 02 // READY</span>
      <span class="bp-flow-title">Worktree <code>#102</code></span>
      <span class="bp-flow-desc">Plan -&gt; TDD -&gt; <code>finish-issue</code> -&gt; PR -&gt; <strong><code>pr-autopilot</code></strong></span>
    </div>
    <div class="bp-flow-node">
      <span class="bp-flow-step">QUEUED IN DAG</span>
      <span class="bp-flow-title"><code>#104</code> (Blocked on <code>#101</code>)</span>
      <span class="bp-flow-desc">Auto-fills Worker 01's slot once you merge <code>#101</code>.</span>
    </div>
  </div>
  <div class="bp-flow-arrow-v" aria-hidden="true"><svg viewBox="0 0 24 24" fill="none" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 5v14M6 13l6 6 6-6"/></svg><span class="bp-flow-cond">Coordinator verifies live PR ground truth first</span><svg viewBox="0 0 24 24" fill="none" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 5v14M6 13l6 6 6-6"/></svg></div>
  <div class="bp-flow-row cols-3-plain">
    <div class="bp-flow-node is-mint">
      <span class="bp-flow-step">INTERRUPT 01 // DONE</span>
      <span class="bp-flow-title">PR Merge-Ready</span>
      <span class="bp-flow-desc">CI green + Copilot drained. Review diff once, click Merge, slot refills from DAG.</span>
    </div>
    <div class="bp-flow-node">
      <span class="bp-flow-step">INTERRUPT 02 // BLOCKED</span>
      <span class="bp-flow-title">Design or Approval Call</span>
      <span class="bp-flow-desc">Sidebar highlights <code>blocked</code> when a worker hits an API or scope choice.</span>
    </div>
    <div class="bp-flow-node is-amber">
      <span class="bp-flow-step">INTERRUPT 03 // STALEMATE</span>
      <span class="bp-flow-title">3-Strike Circuit Breaker</span>
      <span class="bp-flow-desc">Stops after 3 failed attempts instead of burning tokens or Copilot quota.</span>
    </div>
  </div>
</div>

---

## 1. One-Time Herdr Integration Setup

Herdr ships its own version-matched CLI skill inside the `herdr` binary (paired
with the `work-herd` orchestration skill in `~/.config/opencode/skills/work-herd/`):

~~~shell
# Install Herdr state-tracking hooks for OpenCode (also supports antigravity-cli, claude, codex, copilot)
herdr integration install opencode
herdr integration status

# Inspect Herdr's built-in CLI mechanics guide
herdr --skill
~~~

> [!NOTE]
> **Built-in safety gate (`HERDR_ENV=1`):** Every pane inside Herdr sets
> `HERDR_ENV=1`. Both `herdr --skill` and `work-herd` verify `HERDR_ENV=1`
> before running session commands so an external terminal never mutates your
> workspace by mistake.

---

## 2. Narrow the Candidate List, Then Build a Herdr DAG Plan

Real repositories often have hundreds or thousands of open issues, and four
green PRs that all touch the same shared file turn into three rebase conflicts
the moment the first PR merges.

Give your **Coordinator** session inside Herdr a **focus filter** -- next-sprint
priorities, a subsystem or problem area you care about, a label/milestone, or a
handful of issue numbers -- so it shortlists the right candidates first and then
builds the Herdr [DAG](/docs/reference/glossary.md#dag) plan:

~~~text
Give me a shortlist of issues to tackle in the next sprint (or issues related to CLI output and error handling), then come up with a Herdr plan once I confirm the list.
~~~

*(Already know the area, milestone, or issue IDs you want to run? Point at them
directly:)*

~~~text
Look at the open issues related to shortener storage and Go CI (or issues #101-#105) and come up with a Herdr plan.
~~~

Once the shortlist is set, the `work-herd` skill inspects those issues against
the codebase and writes a disposable, git-ignored plan doc
(`temp/plans/<repo>-herd.md`) mapping out:

1. **Logical dependencies** -- e.g., `#104` needs the new interface introduced
   by `#101`.
2. **Hotspot file edges** -- if both `#101` and `#105` edit the same root
   registration table or `actions.yaml`, `work-herd` adds a DAG edge (`#101 ->
   #105`) or assigns the shared table edit to one PR so concurrent branches
   never collide.
3. **Wave lineup & automatic slot refilling** -- Wave 1 fills up to your
   concurrency cap (default: 4 workers); as soon as you merge `#101`, the
   coordinator tears down `#101`'s worktree and immediately provisions `#104`
   into the open slot.

~~~text
Wave 1 (Ready Now -- non-overlapping packages):
  - #101 (internal/shortener)   --> unblocks #104
  - #102 (ci/go-ci)             --> independent
  - #103 (docs/reference)       --> independent

Wave 2 (Queued in DAG):
  - #104 (cmd/shorty flags)     [blocked on #101: uses new Store option]
  - #105 (internal/shortener)   [blocked on #101: hotspot shortener.go]
~~~

Once the DAG looks right, tell the coordinator: *"Looks good -- kick off Wave
1."*

---

## 3. Why Autonomous Copilot (`pr-autopilot` + `waiter`) Makes This Work

If four workers stopped the moment they ran `gh pr create`, you would still spend
the afternoon manually babysitting CI runs and GitHub Copilot comments across
four tabs. Instead, each worker's template (`work-herd/worker-prompt.md`)
automatically hands its new PR to **`pr-autopilot`**:

- **Zero context-window poll spam (`waiter` subagent):** Waiting 4 minutes for CI
  or Copilot runs inside a cheap 1-line `waiter` subagent (Slot 3 utility model)
  rather than dumping per-minute status JSON into the worker's context.
- **Reads all 4 PR surfaces:** Checks paginated GraphQL `reviewThreads`,
  collapsed `<details>` findings inside Copilot's review body, general PR
  comments, and `gh pr checks` -- guarding against Copilot `"insufficient
  remaining quota"` comments posing as clean reviews.
- **Protects Copilot quota (`review_on_push: false`):** Replies and resolves
  nits via GraphQL *without* pushing, batches mechanical code fixes into **one**
  verified push + PR body sync, and waits for **100% green CI** before
  requesting a Copilot delta re-review.
- **Confirms review drain before reporting `DONE`:** Checks the GitHub issue
  timeline to prove no trailing `review_requested` or `copilot_work_started`
  event is still running.

---

## 4. Coordinator Ground-Truth Check & Your Single Merge Step

Across real multi-agent runs, workers occasionally claim `"DONE"` prematurely
(for example, missing a page-2 comment or reporting before a check finishes).
Before notifying you, the `work-herd` coordinator **re-verifies live GitHub
ground truth** (`headRefOid`, `0` unresolved GraphQL threads, `0` open review
body findings, green CI, and drained Copilot timeline) and sends you the full
clickable PR URL.

Your job stays calm even with four workers running on a second workstation
(`herdr --remote <host>`):
1. **Ignore healthy workers** while they code, test, and address mechanical PR
   feedback.
2. **Answer `blocked` design questions** when Herdr's sidebar flags a real
   product choice.
3. **Read each `DONE` PR diff once on GitHub and click Merge** -- watching the
   coordinator pull `main`, clean up the worktree, and launch the next issue
   from your Herdr DAG.
