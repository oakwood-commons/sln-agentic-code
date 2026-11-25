---
name: risk-scaled-effort
description: "Use when starting work to match investigation, model tier, testing, and review depth to actual task risk."
---

# Risk-Scaled Effort

Match effort to risk rather than turning trivial edits into heavyweight quests
or rushing high-risk changes.

## Classify First

- **Low risk** (typos, small doc clarifications, isolated renames): make the
  minimal accurate edit, run the targeted check, and finish.
- **Medium risk** (standard features, bug fixes, CI updates): inspect callers,
  agree on seams when non-obvious, write table-driven tests, and run `/finish-issue`.
- **High risk** (auth, secrets/redaction, concurrency, destructive state, public
  repo/ruleset mutations): gather concrete evidence first, use a stronger
  reasoning model for `/plan` and review, verify `--dry-run` / rollback paths,
  and get explicit user confirmation before state-changing actions.

## Avoid Waste

- Do not re-read unchanged files repeatedly in the same session.
- Do not let a minor review nit trigger an unrequested cross-cutting refactor.
- Escalate to a stronger model after 2 failed attempts on the same check rather
  than spinning in retry loops.
