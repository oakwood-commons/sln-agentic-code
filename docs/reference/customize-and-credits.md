---
title: "Customize, Other Tools & Credits"
weight: 40
---

# Customize Models, Other AI Tools & Credits

How to adjust the starter configuration for different model providers, local
GPUs, or other coding harnesses -- plus open-source credits for everything
adapted in this repository.

---

## 1. Swapping Cloud Models or Using Local GPU Ollama

All model IDs in the starter kit live in **one file**
(`~/.config/opencode/opencode.jsonc`) rather than inside individual agent
Markdown files:

- **Different cloud provider (OpenAI, Google Gemini, GitHub Copilot):** Open
  `~/.config/opencode/opencode.jsonc` and swap the model strings under `"model"`
  and `"agents"` after authenticating with `opencode auth`.
- **Hybrid local GPU + cloud escalation (`ollama-hybrid`):** If you run local
  coding models via Ollama and want `$0` everyday coding loops while keeping
  strong cloud models on `planner`, `code-reviewer`, and `review-buster`,
  install the hybrid profile:
  ~~~bash
  scafctl run solution agentic-code -r profile=ollama-hybrid
  ~~~

---

## 2. Included User-Level Files Overview

Running `scafctl run solution agentic-code` places a minimal, agnostic set of
files into `~/.config/opencode/`:

~~~text
~/.config/opencode/
|-- opencode.jsonc                  # Model routing, permissions & scafctl MCP server
|-- dcp.jsonc                       # Dynamic Context Pruning schema
|-- THIRD-PARTY-NOTICE.txt          # MIT notices for adapted skills
|-- instructions/workflow-policy.md # 2-strike escalation, 3-strike stop, Copilot credit rules
|-- agents/
|   `-- {planner,code-reviewer,artifact-auditor,go-tester,docs-writer,pr-create,pr-reviewer,review-buster}.md
|-- commands/
|   `-- {plan,interview,wait-what,finish-issue,review,pr-create,pr-review,review-buster,handoff}.md
`-- skills/
    `-- {stepwise-interview,wait-what,grilling,codebase-design,tdd,diagnosing-bugs,
         finish-issue,git-github-workflow,go,go-debug,go-perf,scafctl,risk-scaled-effort,markdown}/SKILL.md
~~~

---

## 3. Mapping This Workflow to Other Coding Harnesses

Because **repo contracts** (`AGENTS.md`, `REVIEW.md`, `ci/go-ci/actions.yaml`)
are plain repository files, the exact same 5-step process maps cleanly onto
other harnesses:

| Layer | OpenCode (Included Starter) | Antigravity (`agy`) | Claude Code | Cursor / VS Code + Copilot |
|---|---|---|---|---|
| **Repo Contract** | `AGENTS.md` + `REVIEW.md` | `AGENTS.md` + `REVIEW.md` | `AGENTS.md` (or pointer in `CLAUDE.md`) + `REVIEW.md` | `AGENTS.md` + `.github/copilot-instructions.md` -> `REVIEW.md` |
| **Personal Skills** | `~/.config/opencode/skills/*/SKILL.md` | `~/.gemini/config/skills/*/SKILL.md` | `~/.claude/skills/*/SKILL.md` | `.cursor/rules/*.mdc` or user rules |
| **Isolated Subagents** | `~/.config/opencode/agents/*.md` | `invoke_subagent` (`Model: "inherit"`, parallel `code-reviewer` + `artifact-auditor`) | `.claude/agents/*.md` (`Task` tool) | Background / review agent modes |
| **Slash Commands** | `~/.config/opencode/commands/*.md` | `/plan`, `/grill-me`, `/boost`, `/swarm` + skill triggers | `~/.claude/commands/*.md` | `.github/prompts/*.prompt.md` |
| **`scafctl` MCP** | `opencode.jsonc` (`mcp.servers.scafctl`) | MCP config (`scafctl mcp serve`) | `claude_desktop_config.json` / `.mcp.json` | `.vscode/mcp.json` / `.cursor/mcp.json` |

> [!NOTE]
> **Parallel branches with Git worktrees:** Once comfortable on one branch, you
> can work multiple independent issues in parallel by creating a separate
> worktree per branch (`git worktree add ../<repo>-worktrees/<slug> -b
> feat/<issue>-<slug>`) inside `tmux`, `zellij`, or Herdr. Always start each
> agent session rooted inside its own worktree directory so `gopls` and git
> state never collide.

---

## 4. Open-Source Credits & Licenses

This repository is licensed under
**[Apache-2.0](https://github.com/oakwood-commons/sln-agentic-code/blob/main/LICENSE)**,
with full third-party license texts preserved in
[`THIRD-PARTY-NOTICE.txt`](https://github.com/oakwood-commons/sln-agentic-code/blob/main/THIRD-PARTY-NOTICE.txt).

- **[Matt Pocock -- `mattpocock/skills`](https://github.com/mattpocock/skills)**
  *(MIT License)*: Upstream source and inspiration adapted for `wait-what`,
  `grilling` (which inspired our companion one-question-at-a-time
  `stepwise-interview` skill), `codebase-design`, `tdd`, `diagnosing-bugs`, and
  `/handoff`.
- **[Samuel Berthe -- `samber/cc-skills-golang`](https://github.com/samber/cc-skills-golang)**
  *(MIT License)*: Upstream Go debugging and performance notes adapted in
  `go-debug` and `go-perf`.
- **[`oakwood-commons/scafctl`](https://github.com/oakwood-commons/scafctl)**
  *(Apache-2.0)* & **[`oakwood-commons/spiral-guard`](https://github.com/oakwood-commons/spiral-guard)**
  *(MIT)*: Portable workflow engine, Blueprint design tokens/assets, and loop
  circuit-breaker plugin.
- **[`@tarquinen/opencode-dcp`](https://github.com/Opencode-DCP/opencode-dynamic-context-pruning)**
  & **[`AxeForce/opencode-vibeguard-v2`](https://github.com/AxeForce/opencode-vibeguard-v2)**:
  Context pruning and local secret-placeholder redaction plugins for OpenCode.
- **[`alex-shpak/hugo-book`](https://github.com/alex-shpak/hugo-book)** *(MIT)* &
  **[Inter Font Project](https://github.com/rsms/inter)** *(SIL OFL 1.1)*:
  Documentation theme and `scafctl` logo wordmark vector outlines.
- **John Ousterhout** (*A Philosophy of Software Design*) & **Michael Feathers**
  (*Working Effectively with Legacy Code*): Foundational framing of deep modules
  and test seams.
