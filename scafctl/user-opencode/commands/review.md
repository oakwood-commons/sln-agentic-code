---
description: "Run the read-only code-reviewer agent on the current diff, branch, or specified scope."
agent: code-reviewer
subagent: true
---

Review this code change: $ARGUMENTS

Default to uncommitted and staged changes if no scope is provided. Do not
modify files, commit, push, or open PRs. Return findings ordered by severity
and end with an explicit verdict.
