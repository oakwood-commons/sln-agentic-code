---
title: "Glossary"
weight: 10
---

# Plain-English Glossary

Quick definitions for the terms linked from the walkthrough.

---

## Context (and "Hallucination") {#context}

Everything the AI can currently see when answering you, split across two halves:
- **Repo context (automatic):** Existing code structure, comments, `README.md`,
  `AGENTS.md`, `REVIEW.md`, and recent tool output.
- **Request context (yours to supply):** What you want built, libraries or
  constraints you already prefer (so you do not get grilled backwards on
  settled basics), and open design questions -- **without** prescribing *how* to
  implement it, so the AI can suggest best practices instead of rubber-stamping
  a suboptimal hunch.

What people call *"AI hallucination"* in coding is usually just the AI giving
the best answer it can from incomplete context. Supplying both halves of context
and letting the `planner` [skill](#skill) validate and plan before coding
surfaces missing decisions before code is written.

---

## Deterministic Tools {#deterministic}

Traditional developer tools -- such as `go test -race` (TDD), `gofmt`,
`golangci-lint`, `gosec`, and `scafctl lint` -- that always produce the exact
same pass/fail result for the same input. Because AI models are probabilistic
(they predict plausible text), giving them deterministic tools lets the AI test
its ideas against reality, see exact compiler or test failures, and try again
just like a human developer does.

---

## DAG (Directed Acyclic Graph) {#dag}

A plain-English name for a **dependency map with one-way arrows and no circular
loops** (*"Step B depends on Step A, while Steps A and C have nothing in common
and can run at the same time"*).

DAGs show up in two important places in this workflow:
- **`scafctl` action pipelines (`ci/go-ci/actions.yaml`):** Independent checks
  (`fmt-check`, `vet`, `lint`, `gosec`, `test`) automatically run **in
  parallel**, while `build` waits until all five succeed -- keeping local and CI
  runs fast without fragile shell ordering.
- **Parallel issue waves in Herdr (`work-herd` [skill](#skill)):** Before
  spawning multiple AI workers at once, mapping candidate issues into a DAG
  separates **safe parallel issues** (Wave 1) from issues that depend on another
  PR's new code or touch the same shared *"hotspot"* file (queued for Wave 2).
  Without a DAG, four parallel PRs editing the same file turn into three rebase
  conflicts and wasted review rounds the moment the first PR merges.

---

## `scafctl` {#scafctl}

An open-source declarative scaffolding and workflow CLI
(**[Installation Guide](https://scafctl.oakwood-commons.org/docs/tutorials/getting-started/#installation)**
| **[Documentation](https://scafctl.oakwood-commons.org/)**
| **[`oakwood-commons/scafctl` on GitHub](https://github.com/oakwood-commons/scafctl)**),
packaged natively for **Linux, macOS, and Windows** (`amd64` and `arm64`) with
an embedded cross-platform POSIX shell so every action runs identically on all
three OSes without Git Bash or WSL. Across this guide it handles three jobs:
1. **Safe environment, toolchain & repo setup (`scafctl run solution agentic-code`):**
   Installs starter skills, subagents, and commands into OpenCode
   (`~/.config/opencode/`), Antigravity (`~/.gemini/config/` via
   `-r harness=agy`, automatically generating the companion skills `agy` needs
   for agents and commands), or both (`-r harness=both`); uses multi-step
   resolvers + CEL (`map.merge`, `json.unmarshal`, `json.marshalPretty`) to
   non-destructively merge `workflow-policy.md`, `scafctl` MCP, and safety
   permissions into existing config files (with `.bak` backups); verifies or
   installs Go 1.26+ without admin rights (`--action install-go`); scaffolds
   repo contracts (`--action init-repo`); and configures GitHub rulesets
   idempotently (`--action configure-github`).
2. **Portable [deterministic](#deterministic) CI (`scafctl run action ci`):**
   Executes one cross-platform action [DAG](#dag) (`gofmt`, `go vet`,
   `golangci-lint`, `gosec`, `go test -race`, `go build`) identically across
   Linux, macOS, Windows, GitHub Actions, and Tekton.
3. **Live schema validation via [MCP](#mcp) (`scafctl mcp serve`):** Gives AI
   assistants live schema lookup, CEL/Go-template evaluation, and dry-run
   previews so they never guess YAML fields from memory.

---

## MCP (Model Context Protocol) {#mcp}

A standard local pipe (such as `scafctl mcp serve`) that lets an AI assistant
look up live schemas and run deterministic validators directly instead of
guessing YAML fields from memory.

---

## Model {#model}

The AI brain generating responses (for example, Claude Opus 5.5, Claude Sonnet
5.5, Claude Haiku 5.5, GPT-6.2, Gemini 4, or a local Qwen 3.8 model via
Ollama). Rather than using one model for everything, the starter setup splits
work across **five distinct roles** (see the full breakdown in
[The 5-Role Model Roster](/docs/reference/customize-and-credits.md#model-roster)):
1. **High-reasoning model** (`planner`, `code-reviewer`) for 1--3 turn planning
   and local adversarial review, where catching one bad seam saves 20 turns.
2. **Mid-tier coding workhorse** (`build`, `go-tester`) for 15--40 turn TDD
   implementation loops at 40%--80% lower cost.
3. **Ultra-cheap utility model** (`explore`, `title`, `compaction`) running in
   isolated subagents so file searches never invalidate your main chat's prompt
   cache.
4. **Second model brand** (`review-buster`) from a different vendor to avoid
   correlated blind spots when breaking out of a stuck loop.
5. **Independent PR reviewer** (**GitHub Copilot Code Review**) inspecting only
   the final PR diff in a cold-room context on GitHub.

---

## Harness {#harness}

The coding application on your computer (**OpenCode**, **Antigravity (`agy`)**,
**Claude Code**, **Cursor**, or **VS Code + GitHub Copilot**) that connects a
[model](#model) to your files, terminal, permissions, and tools. This guide uses
**OpenCode** because it is open-source and model-agnostic, while the workflow
itself works in any harness.

---

## Agent & Subagent {#subagent}

An **agent** combines a role prompt, a model choice, and a set of allowed tools:
- **Primary agents** (`plan`, `build`) run in your main chat window.
- **Subagents** (such as `planner`, `code-reviewer`, `security-auditor`,
  `artifact-auditor`, `go-tester`, `docs-writer`, `pr-create`, `pr-reviewer`,
  `pr-autopilot`, `waiter`, and `review-buster`) run a focused job in
  a separate, temporary context and return a concise report. Making `planner`,
  `code-reviewer`, `security-auditor`, `artifact-auditor`, and `waiter`
  read-only (`edit: deny`) prevents them from
  silently changing your code or rubber-stamping mistakes from your main chat
  history.

---

## Skill {#skill}

A short `SKILL.md` playbook (such as `planner`, `create-issue`,
`stepwise-interview`, `wait-what`, `tdd`, `codebase-design`, `go`,
`security-audit`, `diagnosing-bugs`, or `finish-issue`) whose `description:`
header is exposed to the model on every turn alongside a `skill` tool. When you
write in plain English, the model decides in its inner monologue which
playbooks to load -- keeping prompts short without stuffing 2,000 lines of rules
into every turn. Because prose matching is **probabilistic** (strong cloud
models invoke skills reliably, while smaller models can occasionally skip a tool
call), pair skills with always-on `workflow-policy.md` rules and optional
[slash commands](#slash-command). Browse all 22 bundled skills -- plus a
side-by-side comparison of prompts, rules, skills, subagents, and commands -- in
[Skills, Agents & Commands](/docs/reference/skills-and-agents.md).

---

## Slash Command {#slash-command}

A **deterministic** shortcut you type in chat (`/plan`, `/create-issue`,
`/finish-issue`, `/review`, `/pr-create`, `/pr-review`, `/pr-autopilot`,
`/review-buster`, `/handoff`, `/interview`, `/wait-what`) when you want to force
a specific subagent or prompt expansion 100% of the time -- especially useful on
faster or smaller local models where prose skill-selection can be less
consistent.

---

## Plugin {#plugin}

Optional extensions running inside the [harness](#harness) around model turns
and tool calls. The starter `opencode.jsonc` ships with **zero plugins enabled
by default** so `opencode` starts immediately without external plugin downloads;
see [Recommended Optional Plugins](/docs/reference/customize-and-credits.md#recommended-plugins)
for suggested opt-in additions.

---

## Seam & Deep Module {#seam}

Terms used inside the adapted `codebase-design` and `tdd` skills:
- **Deep module** -- a package or function that packs genuinely useful behavior
  behind a small, easy-to-learn interface (such as `shortener.NewHandler(store)`).
- **Seam** -- a clean plug-in point where callers and tests can swap behavior
  without editing internals (for example, passing an in-memory `Store` or fake
  clock in unit tests vs. a real clock/database in production). Avoid inventing
  seams unless at least **two** real callers/tests use them.

---

## Repo Contract {#repo-contract}

Files checked into a Git repository (`AGENTS.md`, `REVIEW.md`,
`.github/copilot-instructions.md`, and `ci/go-ci/actions.yaml`) that state
**what the project expects from every contributor** -- build/test commands,
coverage thresholds, Go error-handling rules, and the PR code review bar --
without forcing anyone's personal AI editor or slash commands.
