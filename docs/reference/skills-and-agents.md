---
title: "Skills, Agents & Commands"
weight: 15
---

# Prompts vs. Rules vs. Skills vs. Subagents vs. Commands

Modern AI coding harnesses give you five different ways to steer an AI model.
Putting instructions in the wrong layer either bloats every prompt, biases code
reviews, or forces you to repeat boilerplate on every turn:

| Primitive | When It Loads | Context Window | Trigger Reliability | Best Used For |
|---|---|---|---|---|
| **1. Chat Prompt** | Current turn only | Main chat | `100%` manual (you type it) | Your plain-English goal, *why*, and specific feature constraints for this task. |
| **2. Standing Rules & Repo Contracts** (`AGENTS.md`, `REVIEW.md`, `workflow-policy.md`) | **Every turn** | Main chat | `100%` always-on | Short, universal invariants: build/test commands, safety limits (`gh pr merge`: `deny`). Keep concise so every turn stays lean. |
| **3. [Skills](/docs/reference/glossary.md#skill)** (`skills/*/SKILL.md`) | **Progressive disclosure:** only the 1-line `description` sits in context until matched | Main chat (or active subagent) | **Probabilistic** (model matches prose to `description`) | Deep domain playbooks (`tdd`, `create-issue`, `jira-issues`, `security-audit`, `go-perf`) needed on specific tasks without paying token cost on unrelated turns. |
| **4. [Subagents](/docs/reference/glossary.md#subagent)** (`agents/*.md`) | Only when delegated | **Fresh, isolated context** (with its own tools & model) | Delegated by skill, command, or main agent | Tasks needing **unbiased eyes** (`code-reviewer`, `security-auditor` with `edit: deny`), **separate models** (`review-buster`), or **noisy loops** (`waiter` polling CI without flooding main chat). |
| **5. [Slash Commands](/docs/reference/glossary.md#slash-command)** (`commands/*.md`) | Only when you type `/name` | Main chat or target subagent | **`100%` deterministic** | Explicit overrides (`/plan`, `/create-issue`, `/finish-issue`, `/wait-what`) when you want to force a workflow rather than rely on prose matching. |

### Why Different Harnesses Need Different Layouts (`OpenCode` vs. `Antigravity / agy`)

Not every coding harness exposes the same primitives:

- **OpenCode (`-r harness=opencode`, default):** Separates user customizations
  into `skills/`, `agents/`, and `commands/`. Workflows like `code-reviewer`,
  `artifact-auditor`, and `security-auditor` live in `agents/*.md`, while
  shortcuts like `/handoff` and `/review` live in `commands/*.md`.
- **Antigravity (`agy`) (`-r harness=agy` or `-r harness=both`):** Slash
  commands in `agy` are built-in rather than user-defined `commands/*.md` files,
  and skill-centric harnesses rely on `SKILL.md` playbooks whenever the main
  agent runs a review, handoff, or audit workflow inline. That means **`agy`
  needs more skills:** every role agent and standalone command workflow also
  gets a companion `SKILL.md` in `~/.gemini/config/skills/` (alongside
  `~/.gemini/config/rules/workflow-policy.md`, `~/.gemini/config/agents/*/agent.md`,
  `~/.gemini/config/THIRD-PARTY-NOTICE.txt`, and `~/.gemini/config/mcp_config.json`).

To keep both harnesses **100% in sync with zero duplicate files to maintain**,
`agentic-code` stores one canonical set of playbooks in `scafctl/user-opencode/`
and uses `scafctl` CEL transforms at install time to generate the extra
Antigravity `SKILL.md` and `agent.md` files automatically when you pass
`-r harness=agy` or `-r harness=both`.

---

## Browse Included Playbooks

Every drawer below renders directly at build time from `scafctl/user-opencode/`,
and `scafctl run action ci` verifies that every bundled Markdown playbook,
standing policy, agent, and command is listed here so this page can never drift
out of sync.

### Standing Workflow Policy (`instructions/*.md`)

Loaded on every turn (`~/.config/opencode/instructions/workflow-policy.md` in
OpenCode; `~/.gemini/config/rules/workflow-policy.md` in Antigravity):

<div class="bp-skill-list">
{{< starter-drawer path="instructions/workflow-policy.md" meta="always-on rule" summary="Standing personal workflow rules: read AGENTS.md/REVIEW.md first, auto-select skills from plain English, prefer best long-term design over shortcuts, enforce 2-strike escalation and 3-round circuit breakers, run finish-issue before PRs, and never auto-merge." >}}
</div>

---

### 1. Planning, Alignment & Multi-Tracker Issue Authoring

Used in **Step 2** (and Herdr planning) when scoping ideas, structuring
two-layer issues (Plain English at the top for humans, grounded codebase context
at the bottom for AI), and working across GitHub Issues, Jira tickets, or local
Markdown issues.

<div class="bp-skill-list">
{{< starter-drawer path="skills/planner/SKILL.md" meta="/plan | planner subagent" >}}
{{< starter-drawer path="skills/create-issue/SKILL.md" meta="/create-issue | human + AI sections" >}}
{{< starter-drawer path="skills/github-issues/SKILL.md" meta="gh issue CLI" >}}
{{< starter-drawer path="skills/jira-issues/SKILL.md" meta="Jira / Atlassian MCP" >}}
{{< starter-drawer path="skills/local-issues/SKILL.md" meta="temp/issues/*.md" >}}
{{< starter-drawer path="skills/stepwise-interview/SKILL.md" meta="/interview" >}}
{{< starter-drawer path="skills/wait-what/SKILL.md" meta="/wait-what" >}}
{{< starter-drawer path="skills/grilling/SKILL.md" >}}
{{< starter-drawer path="skills/risk-scaled-effort/SKILL.md" >}}
</div>

---

### 2. Design, TDD & Go Engineering

Used in **Step 3** to build clean interfaces, write behavior tests that protect
existing code from regressions, and diagnose Go concurrency or performance bugs.

<div class="bp-skill-list">
{{< starter-drawer path="skills/tdd/SKILL.md" >}}
{{< starter-drawer path="skills/codebase-design/SKILL.md" >}}
{{< starter-drawer path="skills/go/SKILL.md" meta="go-tester subagent" >}}
{{< starter-drawer path="skills/diagnosing-bugs/SKILL.md" >}}
{{< starter-drawer path="skills/go-debug/SKILL.md" >}}
{{< starter-drawer path="skills/go-perf/SKILL.md" >}}
</div>

---

### 3. Pre-PR Self-Review & Bounded GitHub PR Loop

Used in **Steps 4 and 5** so the AI audits its own branch locally before pushing
and drives GitHub Actions CI + GitHub Copilot reviews without burning quota.

<div class="bp-skill-list">
{{< starter-drawer path="skills/finish-issue/SKILL.md" meta="/finish-issue" >}}
{{< starter-drawer path="skills/git-github-workflow/SKILL.md" meta="/pr-create | /pr-review" >}}
{{< starter-drawer path="skills/pr-autopilot/SKILL.md" meta="/pr-autopilot | waiter subagent" >}}
</div>

---

### 4. Security Auditing & Public Repo Hygiene

Loaded when auditing a repository for vulnerabilities, reconciling scanner or
pentest findings via reproduce-or-retract, fixing auth/credential flaws inline,
or scrubbing private details before publishing to a public repo.

<div class="bp-skill-list">
{{< starter-drawer path="skills/security-audit/SKILL.md" meta="security-auditor subagent | reproduce-or-retract" >}}
{{< starter-drawer path="skills/security-sweep/SKILL.md" meta="govulncheck | gosec | trivy | gitleaks" >}}
{{< starter-drawer path="skills/publish-hygiene/SKILL.md" meta="public repo scrubber" >}}
</div>

---

### 5. Parallel Multi-Issue Orchestration with Herdr

Used in the **[Advanced Herdr Workflow](/docs/advanced/parallel-agents-with-herdr.md)**
to coordinate 3 to 4 parallel issues in isolated worktrees under one main repo
window.

<div class="bp-skill-list">
{{< starter-drawer path="skills/work-herd/SKILL.md" meta="worker-prompt.md" >}}
{{< starter-drawer path="skills/work-herd/worker-prompt.md" meta="worker handoff template" summary="Self-contained autonomous worker assignment template sent to each isolated Herdr worktree pane (scope, local CI + finish-issue gate, signed PR submission, pr-autopilot convergence, and human-only merge stop)." >}}
{{< starter-drawer path="skills/herdr/SKILL.md" >}}
</div>

---

### 6. `scafctl` Authoring & Documentation Hygiene

Used whenever editing portable CI actions, `scafctl` solutions, `README.md`
files, or Markdown docs.

<div class="bp-skill-list">
{{< starter-drawer path="skills/scafctl/SKILL.md" meta="scafctl MCP" >}}
{{< starter-drawer path="skills/markdown/SKILL.md" meta="docs-writer subagent" >}}
</div>

---

### 7. Included Subagents (`agents/*.md`)

Click any subagent to inspect its system prompt *(frontmatter `permissions` are
summarized in each badge and visible in full via `View raw on GitHub ->`; when
installing with `-r harness=agy` or `both`, `scafctl` installs these into
`~/.gemini/config/agents/` and generates companion `SKILL.md` playbooks in
`~/.gemini/config/skills/` automatically)*:

<div class="bp-skill-list">
{{< starter-drawer path="agents/planner.md" meta="read-only | Slot 1 (High reasoning)" >}}
{{< starter-drawer path="agents/code-reviewer.md" meta="read-only | Slot 1 (High reasoning)" >}}
{{< starter-drawer path="agents/security-auditor.md" meta="read-only | Slot 1 (High reasoning)" >}}
{{< starter-drawer path="agents/artifact-auditor.md" meta="read-only | Slot 2 (Workhorse)" >}}
{{< starter-drawer path="agents/go-tester.md" meta="Slot 2 (Workhorse)" >}}
{{< starter-drawer path="agents/docs-writer.md" meta="docs-only | Slot 2 (Workhorse)" >}}
{{< starter-drawer path="agents/pr-create.md" meta="Slot 2 (Workhorse)" >}}
{{< starter-drawer path="agents/pr-reviewer.md" meta="pause-before-write | Slot 1 (High reasoning)" >}}
{{< starter-drawer path="agents/pr-autopilot.md" meta="Slot 2 (Workhorse)" >}}
{{< starter-drawer path="agents/review-buster.md" meta="Slot 4 (2nd vendor)" >}}
{{< starter-drawer path="agents/waiter.md" meta="read-only | Slot 3 (Cheap utility)" >}}
</div>

---

### 8. Included Slash Commands (`commands/*.md`)

Click any command shortcut to inspect the exact prompt expansion it sends *(in
Antigravity / `agy`, standalone commands without a matching skill name are
exported as companion `SKILL.md` playbooks instead)*:

<div class="bp-skill-list">
{{< starter-drawer path="commands/plan.md" >}}
{{< starter-drawer path="commands/create-issue.md" >}}
{{< starter-drawer path="commands/interview.md" >}}
{{< starter-drawer path="commands/wait-what.md" >}}
{{< starter-drawer path="commands/finish-issue.md" >}}
{{< starter-drawer path="commands/review.md" >}}
{{< starter-drawer path="commands/pr-create.md" >}}
{{< starter-drawer path="commands/pr-review.md" >}}
{{< starter-drawer path="commands/pr-autopilot.md" >}}
{{< starter-drawer path="commands/review-buster.md" >}}
{{< starter-drawer path="commands/handoff.md" >}}
</div>
