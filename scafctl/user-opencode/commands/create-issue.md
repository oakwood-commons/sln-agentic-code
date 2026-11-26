---
description: "Explore the codebase, clarify open edge cases in plain English, and draft a two-layer issue (Plain English for humans + Technical Context for AI) for GitHub, Jira, or local temp/issues/"
---

Load the `create-issue` skill and draft a well-structured issue for:

$ARGUMENTS

1. Read `AGENTS.md` and explore the relevant packages and tests first (`path/to/file`) without modifying any source code.
2. If genuine product or edge-case choices are open, ask one plain-English question at a time (`stepwise-interview`) with a clear recommendation.
3. Present the draft in chat using the two-half format (**Plain English Summary, Behavior & Acceptance Criteria** at the top for fast human review; **Technical Context, Constraints & Verification** at the bottom for AI accuracy without brittle code prescriptions). Wait for confirmation, then file it using `github-issues`, `jira-issues`, or `local-issues`.
