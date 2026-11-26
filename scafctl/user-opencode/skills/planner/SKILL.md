---
name: planner
description: "Use automatically before implementing, planning, or scoping any non-trivial feature, refactor, architecture change, or issue -- even when asked in plain English without a slash command. Inspects the real codebase first, skeptically questions whether the change is worth doing and fits project scope, challenges suboptimal 'how' suggestions, checks prior art (best over easiest), settles open edge cases, and produces a blueprint or issue spec before writing code."
---

# Skeptical Codebase Triage & Planning

Whenever a user proposes a non-trivial feature, asks to draft an issue or Jira
story, or asks to implement an open issue (`#123`, `PROJ-123`, or
`temp/issues/*.md`), **validate and plan before writing production code** --
even if they asked in plain English without a slash command.

## Step 0: Skeptical Triage Gate (Always First)

An idea, GitHub issue, Jira user story, or local issue file is a **hypothesis,
not a mandate** -- reporters often lack full codebase context, and even an
issue drafted earlier by an AI can miss callers or go stale.

1. **Inspect the real code first (`file:line`)**: Read `AGENTS.md`,
   `README.md`, and the relevant packages and tests before forming an opinion.
   Verify whether the gap or bug actually exists, is partially implemented, or
   was already solved.
2. **Question whether we should do it**: Confirm the request fits the project's
   purpose, scope, and maintenance bar. If it is redundant, out of scope, or
   better solved outside core, **push back with `file:line` evidence and stop**
   instead of producing a plan.
3. **Honor stated constraints, challenge suboptimal *how***:
   - Treat libraries, stack choices, and boundaries the user already stated as
     settled (do not grill them backwards on settled basics unless the codebase
     proves a direct conflict).
   - Never rubber-stamp a prescribed *"how"*: if the prompt or issue suggests a
     clunky, racy, or over-engineered mechanism when a simpler, safer, or more
     testable design exists, explain the better approach and why.
4. **Learn from prior art & favor *best* over *easiest***: Check existing
   codebase patterns (`go` skill for Go repos, `codebase-design` for deep
   modules and seams) and how mature, well-regarded tools solve the same
   problem. Explicitly call out any easy-but-inferior shortcut you rejected.
5. **Resolve genuinely open choices simply**: Look up codebase facts yourself.
   When real product or behavior edge cases remain unspoken, load
   `stepwise-interview` and settle them one plain-English question at a time
   with a clear recommendation.

## Output Modes

### Mode A -- Planning or Implementing a Change
Record a short **Validation** block first:
- Verdict (`implement` / `reshape` / `recommend-close`) with `file:line`
  evidence.
- Chosen design vs. rejected shortcuts (and why if diverging from a suggested
  fix), plus existing patterns/prior art reused.

If the verdict is `implement`, present a concise **Blueprint** and wait for the
user's go-ahead before editing production files:
1. **Summary** -- one plain paragraph on what will be built and why.
2. **Design & Seams** -- packages, public interfaces, and test seams.
3. **Ordered Implementation Steps** -- small, independently testable pieces
   (`S`/`M`/`L`).
4. **Error Handling** -- sentinel errors, `%w` wrapping, failure modes.
5. **Testing Strategy** -- unit, `-race`, and integration checks at public
   seams.
6. **Docs, Examples & Risks** -- supporting files and compatibility/security
   guardrails.

### Mode B -- Drafting an Issue or Story (Check Repo Tracker First)
After Step 0 and settling open product rules, load the `create-issue` skill and
draft the issue in **two halves** -- **Plain English at the top** (Summary &
Why, Agreed Behavior & Edge Cases, and Acceptance Criteria so a human can review
and approve in 20 seconds without reading code) followed by **Technical Context
at the bottom** (relevant packages/files, existing patterns, constraints, and
verification seams so an AI implementer does a clean first pass without brittle
line-by-line code prescriptions). Always show the draft in chat first and wait
for human confirmation before filing.

Determine **where** to create or read issues from `AGENTS.md` (so the user only
has to say *"create an issue"* or *"look at issue ..."*):
- **Declared in `AGENTS.md`:** Follow the repository's configured tracker --
  GitHub Issues (`gh issue`), Jira user stories (via MCP), local Markdown files
  (`temp/issues/<NNN>-<slug>.md`), or a declared combination (such as Jira for
  stories + GitHub/local for technical tasks).
- **Not yet declared in `AGENTS.md`:** Ask the user **once** whether this repo
  uses GitHub Issues, Jira, local `temp/issues/`, or a combination, and offer
  to add that one-line rule to `AGENTS.md` so they never have to specify it
  again.
- Treat local issue files, GitHub issues, and Jira stories identically when
  later asked to triage and plan them.
