---
title: "Start Here"
type: docs
---

<div class="bp-hero" role="region" aria-labelledby="bp-hero-heading">
  <div class="bp-hero-header">
    <img class="bp-lockup" src="images/scafctl-logo.svg" alt="scafctl Configuration Flow">
    <span class="bp-edition">OAKWOOD COMMONS / HANDS-ON GUIDE</span>
  </div>
  <div class="bp-hero-body">
    <div>
      <div class="bp-eyebrow"><span class="bp-dot" aria-hidden="true"></span><span>AGENTIC GO WALKTHROUGH</span></div>
      <h1 class="bp-hero-title" id="bp-hero-heading">You steer the design.<br>AI writes and proves the code.</h1>
      <p class="bp-hero-desc">Let AI do what it is best at -- writing Go code, tests, and exhaustive local validation -- while you focus on architecture, plain-English guidance, and one final review on the Pull Request.</p>
      <div class="bp-cta-row">
        <a class="bp-btn-primary" href="docs/walkthrough/01-run-the-starter-app/">Start Step 01 -&gt;</a>
        <a class="bp-btn-secondary" href="docs/reference/glossary/">Open Glossary</a>
      </div>
    </div>
    <svg class="bp-diagram" viewBox="0 0 360 300" aria-hidden="true">
      <g class="construction">
        <path d="M42 25V272M318 25V272M24 50H336M24 239H336M180 12V282" stroke-dasharray="3 5"/>
        <path d="M85 17H275M85 12V22M275 12V22M336 65V221M331 65H341M331 221H341M30 38H54M42 26V50M306 251H330M318 239V263"/>
        <circle cx="180" cy="142" r="112" stroke-dasharray="2 6"/>
      </g>
      <g class="flow">
        <path d="M160 48H190L204 62V100H156V52Q156 48 160 48ZM190 48V62H204M180 101V159M90 178V164Q90 159 96 159H264Q270 159 270 165V178M180 159V178"/>
        <rect x="66" y="178" width="48" height="48" rx="7"/>
        <rect x="156" y="178" width="48" height="48" rx="7"/>
        <rect x="246" y="178" width="48" height="48" rx="7"/>
        <ellipse cx="90" cy="195" rx="10" ry="4"/>
        <path d="M80 195V209C80 215 100 215 100 209V195M80 202C80 208 100 208 100 202M172 195L179 202L172 209M182 195L189 202L182 209"/>
        <path d="M265 193L279 202L265 211Z"/>
      </g>
      <circle class="status" cx="209" cy="100" r="10"/>
      <path class="status-mark" d="M205 100L208 103L213 97"/>
      <text x="90" y="248" text-anchor="middle">CLARIFY</text>
      <text x="180" y="248" text-anchor="middle">TEST + BUILD</text>
      <text x="270" y="248" text-anchor="middle">REVIEW + CI</text>
      <text x="180" y="35" text-anchor="middle">AGENTIC WORKFLOW</text>
    </svg>
  </div>
  <div class="bp-hero-footer">
    <span>DECLARATIVE BY DESIGN.</span>
    <span class="bp-state">5-STEP WALKTHROUGH // ~20 MINUTES</span>
  </div>
</div>

## The Problem: Why "Prompt and Pray" Goes Sideways

When most people try coding with an AI assistant, the first five minutes feel
like magic -- and the next hour feels like cleaning up a mess. Three things
cause almost all of that frustration:

1. **What people call "hallucination" is really missing context.** An AI does
   not invent wrong code out of spite -- it gives the best answer it can based
   on whatever [context](/docs/reference/glossary.md#context) is in front of it.
   The mistake people make is either assuming the AI knows unspoken context a
   human teammate would have, or telling the AI *how* to code something when it
   actually knows a cleaner design. Fix that by **combining repo context
   (`AGENTS.md`, `REVIEW.md`) with clear request context** (your goals and
   preferred libraries, without prescribing *how*) and **planning first** before
   writing code.
2. **AI is probabilistic, so it needs good tests and deterministic tools.** Left
   to itself, an AI only predicts code that *looks* plausible -- and adding a new
   feature can quietly break existing behavior you already shipped. Pairing the
   AI with **strong behavior tests (`go test -race`) and linters (`golangci-lint`,
   `gosec`)** gives it [deterministic](/docs/reference/glossary.md#deterministic)
   feedback to write better new code *and* gives you the regression safety net to
   let AI extend your codebase with confidence.
3. **Unchecked pull requests burn time and review credits.** Opening a PR before
   running local checks triggers failing CI builds, noisy review threads, and
   wasted GitHub Copilot quota on broken commits.

---

## Why You Only Need a Terminal and a Browser (No IDE Required)

Writing boilerplate and running validation loops is now an AI job; **making
architectural decisions, providing product guidance, and deciding whether to
accept the finished PR is the developer's job**. With a CLI
[harness](/docs/reference/glossary.md#harness) and good guardrails, **you do not
need a code editor open at all while building**:

1. **Steer architecture & intent in words (Step 2):** Talk through what you want
   and why in plain English, make the design calls that matter, and capture the
   agreement in an issue (GitHub Issues, Jira via MCP, or local `temp/issues/`).
2. **Let AI write code & prove behavior (Steps 3--4):** The AI writes tests and
   code, verifies them with deterministic CI (`go test -race`, `golangci-lint`,
   `gosec`), and runs its own fresh multi-agent self-review (`finish-issue`)
   before opening a PR -- pausing only when an architectural trade-off needs your
   call.
3. **Make the final call once -- on the Pull Request (Step 5):** Because the AI
   already validated everything it could on your branch, **the only time you
   inspect code line-by-line is once on the finished GitHub PR** before clicking
   Merge.

---

## The Fix: Four Small Habits

You do not need a complicated framework to fix this. You only need four small
habits in order:

<div class="bp-pillars">
  <div class="bp-card">
    <div class="bp-card-head">
      <svg viewBox="0 0 24 24" fill="none" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><circle cx="12" cy="12" r="10"/><path d="M9.09 9a3 3 0 015.83 1c0 2-3 3-3 3"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
      <h3>1. Supply Context &amp; Plan First</h3>
    </div>
    <p>Give the AI your repo context, let it ask plain-English questions one at a time, and have it explain its plan before touching code.</p>
  </div>
  <div class="bp-card">
    <div class="bp-card-head">
      <svg viewBox="0 0 24 24" fill="none" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M14.7 6.3a1 1 0 000 1.4l1.6 1.6a1 1 0 001.4 0l3.77-3.77a6 6 0 01-7.94 7.94l-6.91 6.91a2.12 2.12 0 01-3-3l6.91-6.91a6 6 0 017.94-7.94l-3.76 3.76z"/></svg>
      <h3>2. Rely on Good Tests &amp; Deterministic CI</h3>
    </div>
    <p>Pair the AI with strong public-behavior tests, race checks, and linters via one portable <a href="docs/reference/glossary/#scafctl"><code>scafctl</code></a> command -- helping AI write better code while proving existing features never break.</p>
  </div>
  <div class="bp-card">
    <div class="bp-card-head">
      <svg viewBox="0 0 24 24" fill="none" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M22 11.08V12a10 10 0 11-5.93-9.14"/><path d="M22 4L12 14.01l-3-3"/></svg>
      <h3>3. Let AI Self-Review Before the PR</h3>
    </div>
    <p>Have fresh, read-only AI subagents audit the branch for bugs and forgotten docs automatically before opening a PR -- saving CI runs, Copilot quota, and your attention.</p>
  </div>
  <div class="bp-card">
    <div class="bp-card-head">
      <svg viewBox="0 0 24 24" fill="none" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
      <h3>4. Bound the PR Loop</h3>
    </div>
    <p>Give GitHub Copilot a clear repo rulebook (<code>REVIEW.md</code>), wait for green CI before re-reviewing, stop at 3 rounds, and keep a human on Merge.</p>
  </div>
</div>

---

## What You Will Build in This Walkthrough

Inside [`examples/go-app/`](https://github.com/oakwood-commons/sln-agentic-code/tree/main/examples/go-app)
is **`shorty`**, a tiny ~150-line Go URL shortener with two endpoints:
- `POST /links` -- takes a long URL and returns a 6-character short code.
- `GET /{code}` -- redirects your browser to the original URL.

Right now, links in `shorty` live forever and you cannot tell whether anyone
ever clicked them. Across five short steps, you will add **optional link
expiration (`expires_in`) and a click counter (`GET /links/{code}/stats`)**:

<div class="bp-roadmap">
  <a class="bp-roadmap-item" href="docs/walkthrough/01-run-the-starter-app/">
    <span class="bp-rm-num">01</span>
    <span>
      <span class="bp-rm-title">Run the Starter App</span>
      <span class="bp-rm-desc">Install scafctl and OpenCode, merge the starter workflow into your config, and shorten your first URL.</span>
    </span>
    <span class="bp-rm-time">~3 MIN</span>
  </a>
  <a class="bp-roadmap-item" href="docs/walkthrough/02-ask-and-plan-first/">
    <span class="bp-rm-num">02</span>
    <span>
      <span class="bp-rm-title">Ask &amp; Plan First</span>
      <span class="bp-rm-desc">Speak in plain English while auto-loaded skills inspect the code, question assumptions, and plan first.</span>
    </span>
    <span class="bp-rm-time">~4 MIN</span>
  </a>
  <a class="bp-roadmap-item" href="docs/walkthrough/03-build-in-small-slices/">
    <span class="bp-rm-num">03</span>
    <span>
      <span class="bp-rm-title">Build One Piece at a Time</span>
      <span class="bp-rm-desc">Add link expiration and a thread-safe click counter one failing test at a time.</span>
    </span>
    <span class="bp-rm-time">~5 MIN</span>
  </a>
  <a class="bp-roadmap-item" href="docs/walkthrough/04-check-locally-first/">
    <span class="bp-rm-num">04</span>
    <span>
      <span class="bp-rm-title">Let AI Self-Review Before the PR</span>
      <span class="bp-rm-desc">Have fresh AI subagents catch edge cases and missing docs on your branch so you only inspect code once on the PR.</span>
    </span>
    <span class="bp-rm-time">~4 MIN</span>
  </a>
  <a class="bp-roadmap-item" href="docs/walkthrough/05-open-pr-and-copilot-loop/">
    <span class="bp-rm-num">05</span>
    <span>
      <span class="bp-rm-title">Open a Clean PR &amp; Bound Copilot</span>
      <span class="bp-rm-desc">Pair your local checks with GitHub Copilot Code Review without burning credits or looping forever.</span>
    </span>
    <span class="bp-rm-time">~4 MIN</span>
  </a>
</div>

New to `scafctl`? Grab it for Linux, macOS, or Windows via the
**[`scafctl` Installation Guide](https://scafctl.oakwood-commons.org/docs/tutorials/getting-started/#installation)**
before starting Step 01. Need to look up a term, set up your own repository, or
wire Tekton/GitHub Actions CI later? Everything extra lives in the **Reference &
Lookup** section at the bottom of the left sidebar so it stays out of your way
while you walk through the steps.
