---
title: "Customize, Other Tools & Credits"
weight: 40
---

# Customize Models, Plugins, Other Tools & Credits

How to adjust the starter configuration for different model providers, local
GPUs, optional plugins, or other coding harnesses, plus open-source credits for
everything adapted in this repository.

---

## 1. The 5-Role Model Roster & Swapping Providers {#model-roster}

Using your most expensive model for every file search burns budget fast, while
using a bargain model for up-front architecture burns even more budget in failed
test retries. In practice, you only need **five model slots**, all configured in
one file (`~/.config/opencode/opencode.jsonc`):

| Slot | Where It Runs | Why It Earns Its Spot | Starter Default |
|---|---|---|---|
| **1. High-reasoning planner & local reviewer** | `plan`, `planner`, `code-reviewer`, `pr-reviewer`, `security-auditor` | Planning and security/diff review take 1 to 3 turns; catching one bad interface or false-positive boundary saves 20 turns of rewrite churn. | `anthropic/claude-opus-5-5` |
| **2. Code-writing workhorse** | `build`, `general`, `go-tester`, `artifact-auditor`, `docs-writer`, `pr-create`, `pr-autopilot` | Coding and testing take 15 to 40 tool turns; once the plan and `go test -race` are set, a mid-tier coder is 40% to 80% cheaper with near-flagship Go quality. | `anthropic/claude-sonnet-5-5` *(or local `qwen3.8:27b`)* |
| **3. Ultra-cheap utility subagents** | `explore`, `waiter`, `title`, `compaction` | Factual file search (`grep`/`glob`), bounded CI/PR polling (`waiter`), chat titles, and context pruning run in separate contexts at 3x to 10x lower cost and latency. | `anthropic/claude-haiku-5-5` *(or local `qwen3.5:9b`)* |
| **4. Second model brand (loop-breaker)** | `review-buster` & 2-strike escalations | Same-vendor models share training data and **correlated blind spots**. A second model family thinks differently when a loop gets stuck. | `openai/gpt-6.2-astra` *(swap to same vendor on Day 1 if using one API account)* |
| **5. Cold-room PR reviewer** | GitHub Copilot Code Review | Inspects only the final PR diff and `REVIEW.md` on GitHub with a separate review pipeline, reliably catching subtle edge cases local review missed. | Repo ruleset (`review_on_push: false`) |

> [!TIP]
> **Two cost rules most guides miss:**
> 1. **Do not go *too* cheap on Slot 2:** A weak coding model that takes 6 failed
>    test loops re-sends a 40k-token context window 6 times and costs *more* than
>    one clean pass on a mid-tier workhorse, which is why the **2-strike
>    escalation** rule steps up early.
> 2. **Keep Slot 3 in isolated subagents (Prompt Cache rule):** Never manually
>    flip models mid-conversation inside your main chat window just for one
>    "easy" turn. Changing model IDs mid-thread misses the provider's **~90%
>    prompt-cache discount** on existing chat history, whereas OpenCode's
>    separate `explore`, `title`, and `compaction` subagents run in their own
>    tiny contexts cleanly.

### Swapping Cloud Providers or Using Local GPU Ollama

- **Different cloud provider (OpenAI, Google Gemini, GitHub Copilot, OpenRouter):**
  Open `~/.config/opencode/opencode.jsonc` and swap the model strings under
  `"model"` and `"agents"` after authenticating with `opencode auth`. *(If you
  only connect one provider on Day 1, point `review-buster` at that provider's
  strongest reasoning model until you add a second brand.)*
- **Hybrid local GPU + cloud escalation (`ollama-hybrid`):** If you run local
  coding models via Ollama and want `$0` everyday coding loops while keeping
  strong cloud models on `planner`, `code-reviewer`, and `review-buster`,
  install the hybrid profile (add `-r overwrite=true` if switching an
  already-installed `cloud` starter config; `.bak` backups are created
  automatically):
  ~~~shell
  scafctl run solution agentic-code -r profile=ollama-hybrid -r overwrite=true
  ~~~

---

## 2. Recommended Optional Plugins {#recommended-plugins}

The starter `opencode.jsonc` intentionally ships with **zero active plugins** so
`opencode` starts immediately with no extra package downloads or corporate-proxy
friction, and nothing in the 5-step walkthrough depends on a plugin.

Once your base workflow is running, two open-source OpenCode plugins are worth
enabling depending on how you work:

| Plugin | What It Does | When to Enable It |
|---|---|---|
| **[`oakwood-commons/spiral-guard`](https://github.com/oakwood-commons/spiral-guard)** | Detects when an agent repeats the same failing command or reads dozens of files without progress, nudging it to pause, load `diagnosing-bugs`, and test one hypothesis at a time. | Long debugging sessions, local GPU models, or autonomous Herdr workers. |
| **[`AxeForce/opencode-vibeguard-v2`](https://github.com/AxeForce/opencode-vibeguard-v2)** | Replaces detected API keys and tokens with deterministic placeholders before prompts leave your machine for a cloud model, restoring them locally when shell tools run. | Repositories that use local `.env` files, test tokens, or cloud credentials. |

To enable either or both, add a `"plugins"` array to `~/.config/opencode/opencode.jsonc`:

~~~jsonc
{
  "plugins": [
    "github:oakwood-commons/spiral-guard#v2",
    "github:AxeForce/opencode-vibeguard-v2"
  ]
}
~~~

---

## 3. Included User-Level Files & Installer Flags

Running `scafctl run solution agentic-code` places a minimal, agnostic set of
files into `~/.config/opencode/` (generated directly from the solution bundle
via `scafctl run solution agentic-code --action list-files -o tree`; browse the
full text of each playbook in
[Skills, Agents & Commands](/docs/reference/skills-and-agents.md)):

{{< resource-code "opencode-tree.txt" "tree" >}}

Common user-setup and `--action install-go` overrides (`-r key=value`):
- **`harness`** (`"opencode"` default | `"agy"` | `"both"`), **`config_dir`**
  (`~/.config/opencode`), and **`agy_config_dir`** (`~/.gemini/config`).
- **`overwrite`** (`false` default | `true` to refresh bundled skills, agents,
  commands, and OpenCode profile with `.bak` backups while preserving unrelated
  `mcp_config.json` servers).
- **`install_gh`** (`true` default), **`gh_version`** (`"2.106.0"`), and
  **`user_bin_dir`** (`~/.local/bin`) for non-root CLI shims.
- **`go_version`** (`"1.26.11"`) and **`go_install_dir`** (`~/.local/go`) when
  running `--action install-go`.

---

## 4. Installing for Antigravity (`agy`) & Mapping to Other Harnesses

`agentic-code` installs natively into **both OpenCode and Antigravity (`agy`)**
from the same single-source bundle:

~~~shell
# OpenCode only (default -> ~/.config/opencode)
scafctl run solution agentic-code

# Antigravity (agy) only (-> ~/.gemini/config)
scafctl run solution agentic-code -r harness=agy

# Install and keep both harnesses configured on the same machine
scafctl run solution agentic-code -r harness=both
~~~

Because `agy` uses built-in slash commands rather than a user `commands/*.md`
folder and relies on `SKILL.md` playbooks when the main agent executes roles
inline, `-r harness=agy` uses `scafctl` CEL transforms to automatically export
role agents and standalone commands as companion `~/.gemini/config/skills/*/SKILL.md`
playbooks alongside `~/.gemini/config/agents/*/agent.md`,
`~/.gemini/config/rules/workflow-policy.md`,
`~/.gemini/config/THIRD-PARTY-NOTICE.txt`, and a non-destructive merge of
`scafctl mcp serve` into `~/.gemini/config/mcp_config.json`.

Because **repo contracts** (`AGENTS.md`, `REVIEW.md`, `ci/go-ci/actions.yaml`)
are plain repository files, the exact same 5-step process also maps cleanly onto
other harnesses:

| Layer | OpenCode (`-r harness=opencode`) | Antigravity (`-r harness=agy`) | Claude Code | Cursor / VS Code + Copilot |
|---|---|---|---|---|
| **Repo Contract** | `AGENTS.md` + `REVIEW.md` | `AGENTS.md` + `REVIEW.md` | `AGENTS.md` (or pointer in `CLAUDE.md`) + `REVIEW.md` | `AGENTS.md` + `.github/copilot-instructions.md` -> `REVIEW.md` |
| **Standing Rules** | `~/.config/opencode/instructions/workflow-policy.md` | `~/.gemini/config/rules/workflow-policy.md` | `~/.claude/CLAUDE.md` | User rules / `.cursor/rules/*.mdc` |
| **Personal Skills** | `~/.config/opencode/skills/*/SKILL.md` (26 domain skills) | `~/.gemini/config/skills/*/SKILL.md` (26 domain + 14 agent/command companion skills) | `~/.claude/skills/*/SKILL.md` | `.cursor/rules/*.mdc` or `.github/skills/*/SKILL.md` |
| **Isolated Subagents** | `~/.config/opencode/agents/*.md` | `~/.gemini/config/agents/*/agent.md` (`invoke_subagent`, `Model: "inherit"`) | `.claude/agents/*.md` (`Task` tool) | Background / review agent modes |
| **Slash Commands** | `~/.config/opencode/commands/*.md` | Built-in (`/plan`, `/boost`, `/swarm`) + companion skills | `~/.claude/commands/*.md` | `.github/prompts/*.prompt.md` |
| **`scafctl` MCP** | `opencode.jsonc` (`mcp.servers.scafctl`) | `~/.gemini/config/mcp_config.json` (`mcpServers.scafctl`) | `claude_desktop_config.json` / `.mcp.json` | `.vscode/mcp.json` / `.cursor/mcp.json` |

> [!NOTE]
> **Scaling to parallel branches:** Ready to work three or four independent
> issues at once across isolated Git worktrees without babysitting panes? See
> [Parallel Issues with Herdr](/docs/advanced/parallel-agents-with-herdr.md).

---

## 5. Open-Source Credits & Licenses

This repository is licensed under
**[Apache-2.0](https://github.com/oakwood-commons/sln-agentic-code/blob/main/LICENSE)**,
with third-party MIT license texts for adapted skills preserved in
[`THIRD-PARTY-NOTICE.txt`](https://github.com/oakwood-commons/sln-agentic-code/blob/main/THIRD-PARTY-NOTICE.txt).

- **[Matt Pocock: `mattpocock/skills`](https://github.com/mattpocock/skills)**
  *(MIT License)*: Upstream source and inspiration adapted for `wait-what`,
  `grilling` (which inspired the companion one-question-at-a-time
  `stepwise-interview` skill), `codebase-design`, `tdd`, `diagnosing-bugs`, and
  `/handoff`.
- **[Samuel Berthe: `samber/cc-skills-golang`](https://github.com/samber/cc-skills-golang)**
  *(MIT License)*: Upstream Go debugging and performance notes adapted in
  `go-debug` and `go-perf`.
- **[`oakwood-commons/scafctl`](https://github.com/oakwood-commons/scafctl)**
  *(Apache-2.0)*: Portable workflow engine, MCP server, and Blueprint design
  system.
- **[OpenCode](https://opencode.ai)** & **[`herdrdev/herdr`](https://github.com/herdrdev/herdr)**:
  Open-source terminal coding harness and agent-aware workspace manager covered
  across the walkthrough and Advanced guide.
