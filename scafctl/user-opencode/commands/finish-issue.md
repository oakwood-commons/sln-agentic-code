---
description: "Run the local 'Issue Definition of Done' quality gate (parallel code-reviewer + artifact-auditor over one frozen diff, merged fix pass, and final lint) before committing."
agent: build
---

Load the `finish-issue` skill and run its local quality gate over the current
changes.

$ARGUMENTS
