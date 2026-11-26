---
name: create-issue
description: "Use when drafting an issue, bug report, feature request, or user story: explores live code first, writes Plain English at the top for human review plus Technical Context at the bottom for AI execution, previews in chat, and files via github-issues, jira-issues, or local-issues."
---

# Two-Layer Issue Authoring (`create-issue`)

Every well-written issue serves two different readers in sequence:
1. **The human decision-maker right now**, who wants to verify what will be
   built, why it matters, and how edge cases behave in plain English without
   wading through file paths or Go types.
2. **The AI (or teammate) implementing it next**, who starts in a fresh context
   window and produces much stronger code when given grounded package pointers,
   existing patterns, and test expectations up front.

Never write implementation code while running this skill. Explore, align,
preview the draft, and file only after the human confirms.

## 1. Explore & Clarify Before Drafting

1. **Inspect real code first:** Read `AGENTS.md`, `README.md`, and the relevant
   packages and tests (`path/to/file`) to confirm the gap or bug exists and note
   existing interfaces, error conventions, and test helpers.
2. **Clarify open product rules simply:** Do not guess user-facing behavior.
   When 1 to 3 product or edge-case choices are genuinely open (status codes,
   defaults, error handling, non-goals), load `stepwise-interview` and settle
   them one plain-English question at a time with a recommended default.

## 2. Draft With the Two-Half Template

Keep the top half jargon-free for a 20-second human review, and place grounded
codebase context in the bottom half without locking implementers into brittle
line-by-line code:

~~~markdown
# <Clear action title, e.g. feat(shortener): optional link expiration and click stats>

## Summary & Why
<2 to 4 plain-English sentences stating the problem or use case, who benefits, and why it matters.>

## Agreed Behavior & Edge Cases
- <Plain-English rule 1: defaults, inputs, and happy-path behavior>
- <Plain-English rule 2: edge cases, error responses, or expiration states>
- **Out of scope:** <Anything explicitly excluded so scope stays tight>

## Acceptance Criteria
- [ ] <Observable outcome 1 a human can verify without reading source code>
- [ ] <Observable outcome 2>
- [ ] <Existing behavior remains unbroken and covered by automated tests>

---

## Technical Context (For Implementation)
- **Relevant packages & files:** <Package paths, entrypoints, and existing types/interfaces found during exploration>
- **Existing patterns to reuse:** <Test helpers, sentinel errors, and seams/constructors already in the repo>
- **Constraints & guardrails:** <Settled library boundaries, security/trust checks, compatibility rules, and tempting shortcuts to avoid without prescribing line-by-line code>
- **Verification & scope:** <Public test seams to add (`go test -race`, table-driven cases, docs/README updates) plus rough size: XS / S / M / L>
~~~

### Formatting & Hygiene Rules
- Load `markdown` for formatting rules (tilde fences `~~~` around nested
  backticks, and clean prose with zero em-dashes or `--` dash clauses).
- On public repositories, apply `publish-hygiene` before previewing or filing so
  secrets, `/home/<user>` paths, private hostnames, and internal org names never
  appear in an issue.

## 3. Preview First, Then File via the Matching Tracker Skill

1. **Show the draft in chat first** and wait for human approval before creating
   any external ticket.
2. **Load the matching tracker skill** based on what the user asked for (or the
   repository default in `AGENTS.md` when no tracker was named):
   - **GitHub Issues** (`"create a GitHub issue..."` or repo default `github`):
     load `github-issues`.
   - **Jira Tickets** (`"create a Jira ticket / story..."` or repo default
     `jira`): load `jira-issues`.
   - **Local Markdown** (`"save a local issue..."` or repo default `local`):
     load `local-issues`.
