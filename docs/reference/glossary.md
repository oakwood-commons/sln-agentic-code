---
title: "Glossary"
weight: 10
---

# Plain-English Glossary

Quick definitions for the terms linked from the walkthrough.

---

## Context (and "Hallucination") {#context}

Everything the AI can currently see when answering you: your prompt, files it
has read, `AGENTS.md`, `REVIEW.md`, and recent tool output. What people call
*"AI hallucination"* in coding is usually just the AI giving the best answer it
can from incomplete context because we assumed it knew unspoken project rules a
human teammate would know. Supplying repo files and having the AI explain its
plan first (`/interview`, `/plan`) surfaces missing context before code is
written.

---

## Deterministic Tools {#deterministic}

Traditional developer tools -- such as `go test -race` (TDD), `gofmt`,
`golangci-lint`, `gosec`, and `scafctl lint` -- that always produce the exact
same pass/fail result for the same input. Because AI models are probabilistic
(they predict plausible text), giving them deterministic tools lets the AI test
its ideas against reality, see exact compiler or test failures, and try again
just like a human developer does.

---

## MCP (Model Context Protocol) {#mcp}

A standard local pipe (such as `scafctl mcp serve`) that lets an AI assistant
look up live schemas and run deterministic validators directly instead of
guessing YAML fields from memory.

---

## Model {#model}

The AI brain generating responses (for example, Claude Opus, Claude Sonnet,
GPT-5, Gemini, or a local Qwen model running via Ollama). Different steps reward
different models:
- **Stronger reasoning models** for planning (`/plan`) and code review
  (`code-reviewer`), where catching one subtle design flaw saves an hour.
- **Faster workhorse models** for writing bounded Go tests and code once the
  plan is settled.
- **Tiny utility models** for quick file searches, chat titles, and summarizing
  long sessions.

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
- **Subagents** (`planner`, `code-reviewer`, `artifact-auditor`, `go-tester`,
  `pr-reviewer`, `review-buster`) run a focused job in a separate, temporary
  context and return a concise report. Making `planner` and `code-reviewer`
  read-only (`edit: deny`) prevents them from silently changing your code or
  rubber-stamping mistakes from your main chat history.

---

## Skill {#skill}

A short `SKILL.md` instruction file loaded **only when needed** (for example,
`stepwise-interview`, `wait-what`, `tdd`, `codebase-design`, `go`,
`diagnosing-bugs`, or `finish-issue`). Loading focused playbooks on demand keeps
everyday prompts small and accurate compared to stuffing 2,000 lines of rules
into every single turn.

---

## Slash Command {#slash-command}

A short command you type in chat (`/interview`, `/wait-what`, `/plan`,
`/review`, `/finish-issue`, `/pr-create`, `/pr-review`, `/review-buster`,
`/handoff`) that selects the right agent and expands a repeatable prompt.

---

## Plugin {#plugin}

Small extensions running inside the [harness](#harness) around model turns and
tool calls. Our starter `opencode.jsonc` enables three open-source plugins:
- **`opencode-vibeguard-v2`** -- replaces detected secrets/tokens with safe
  placeholders before sending prompts to a cloud model, restoring them locally
  when tools run.
- **`@tarquinen/opencode-dcp`** -- prunes outdated file reads and verbose tool
  outputs from long conversations so context stays sharp and affordable.
- **`oakwood-commons/spiral-guard`** -- detects when an agent fails the same
  command repeatedly or reads dozens of files without progress, nudging it to
  pause and reassess.

---

## Seam & Deep Module {#seam}

Terms from software design used in our `codebase-design` and `tdd` skills:
- **Deep module** *(adapted from John Ousterhout)* -- a package or function that
  packs genuinely useful behavior behind a small, easy-to-learn interface
  (such as `shortener.NewHandler(store)`).
- **Seam** *(from Michael Feathers)* -- a clean plug-in point where callers and
  tests can swap behavior without editing internals (for example, passing an
  in-memory `Store` or fake clock in unit tests vs. a real clock/database in
  production). Avoid inventing seams unless at least **two** real callers/tests
  use them.

---

## Repo Contract {#repo-contract}

Files checked into a Git repository (`AGENTS.md`, `REVIEW.md`,
`.github/copilot-instructions.md`, and `ci/go-ci/actions.yaml`) that state
**what the project expects from every contributor** -- build/test commands,
coverage thresholds, Go error-handling rules, and the PR code review bar --
without forcing anyone's personal AI editor or slash commands.
