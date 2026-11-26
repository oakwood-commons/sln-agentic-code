---
description: "Run the pre-PR AI self-review gate (parallel code-reviewer + artifact-auditor over the branch diff, merged fix pass, and final lint) before opening a pull request."
agent: build
---

Load the `finish-issue` skill and run its pre-PR self-review gate over the
current branch or uncommitted changes.

$ARGUMENTS
