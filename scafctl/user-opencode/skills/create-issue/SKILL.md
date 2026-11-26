---
name: create-issue
description: "Use when drafting or filing an issue, feature request, bug report, or user story (GitHub Issues, Jira via MCP, or local temp/issues/*.md). Explores the codebase first, structures every issue into two halves -- Plain English at the top so a human can review and approve intent in 20 seconds without reading code, followed by grounded Technical Context at the bottom so an AI implementer succeeds on the first pass without brittle code prescriptions -- and shows the draft for human confirmation before filing."
---

# Two-Layer Issue Authoring (Plain English for Humans + Context for AI)

Every well-written issue serves **two different readers** in sequence:
1. **The human decision-maker right now** -- who wants to verify *what* will be
   built, *why*, and how edge cases behave in plain English without wading
   through file paths or Go types.
2. **The AI (or teammate) implementing it next** -- who starts in a fresh
   context window and does dramatically better work when given grounded package
   pointers, existing patterns, and test expectations up front.

Never implement code while running this skill. Explore, align, draft, and file
only after the human approves.

## 1. Explore & Settle Open Product Questions First

1. **Read repo rules first:** Check `AGENTS.md` and `README.md` to see the
   project's scope, conventions, and configured issue tracker (`github`, `jira`,
   `local` `temp/issues/`, or a combination).
2. **Verify against real code (`path/to/file`):** Inspect the relevant packages
   and tests to confirm the feature gap or bug actually exists and see what
   interfaces, patterns, and test helpers already exist.
3. **Clarify unspoken product rules simply:** Do not guess user-facing behavior.
   If 1--3 product or edge-case choices are genuinely open (status codes,
   defaults, error handling, non-goals), load `stepwise-interview` and settle
   them one plain-English question at a time with a recommended default.

## 2. Draft Using the Two-Half Structure

Keep the **top half** jargon-free for fast human sign-off, and put **repo
grounding** in the bottom half without prescribing brittle line-by-line code:

~~~markdown
# <Clear action title -- e.g. feat(shortener): optional link expiration and click stats>

## Summary & Why
<2--4 plain-English sentences stating the problem or use case, who it helps, and why it matters.>

## Agreed Behavior & Edge Cases
- <Plain-English rule 1 -- defaults, inputs, and happy-path behavior>
- <Plain-English rule 2 -- edge cases, error responses, or expiration/empty states>
- **Out of scope:** <Anything explicitly excluded so scope does not creep>

## Acceptance Criteria
- [ ] <Observable outcome 1 a human can verify without reading source code>
- [ ] <Observable outcome 2>
- [ ] <Existing behavior remains unbroken and covered by automated tests>

---

## Technical Context (For Implementation)
- **Relevant packages & files:** <Package paths, entrypoints, and existing types/interfaces discovered during exploration>
- **Existing patterns to reuse:** <Test helpers, error sentinel conventions, seams/constructors already in the repo>
- **Constraints & guardrails:** <Settled library/stack boundaries, trust/security checks, compatibility rules, and any tempting shortcut to avoid -- state boundaries, NOT line-by-line code to copy>
- **Verification & scope:** <Public test seams to add/extend (`go test -race`, table-driven cases, docs/README updates) + rough size: XS / S / M / L>
~~~

### Formatting & Hygiene Rules
- **Do not lock in brittle *how*:** Point to seams, interfaces, and constraints,
  not 80 lines of pseudo-implementation that may go stale or block a cleaner
  design during planning.
- **Markdown safety:** Use ASCII punctuation (`--`, `->`, `...`) and wrap outer
  code blocks in tildes (`~~~`) whenever inner examples contain backticks.
- **Public repo scrub:** On public repositories, apply `publish-hygiene` before
  showing or posting -- never include real secrets, local `/home/<user>` paths,
  private hostnames, or non-public organization names.

## 3. Show the Draft First, Then File

1. **Always preview in chat first:** Show the drafted title and body (or the
   saved draft path when preparing multiple security/backlog drafts) and wait
   for the human's go-ahead before creating external tracker items.
2. **Route per `AGENTS.md`:**
   - **GitHub Issues (`github`):** Write the body to a temporary file and run
     `gh issue create --title "..." --body-file <tmpfile>` to avoid shell
     escaping bugs.
   - **Jira (`jira`):** Create the story via the configured Jira MCP tool using
     the repository's project key.
   - **Local Markdown (`local`):** Write to git-ignored
     `temp/issues/<NNN>-<slug>.md`.
   - **Not declared yet in `AGENTS.md`:** Ask the user **once** which tracker
     this repo uses and offer to record that one-line rule in `AGENTS.md`.
