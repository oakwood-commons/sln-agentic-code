# Personal Agentic Workflow Policy

- **Repo contract vs. personal workflow**: Always read the target repository's
  `AGENTS.md` and `REVIEW.md` first for repo-specific build, test, lint, and
  architecture rules. Apply your personal planning, TDD, and `/finish-issue`
  discipline on top of those rules.
- **Plain-language alignment**: When requirements or design trade-offs are
  unclear, load `stepwise-interview` and ask one plain-English question at a
  time with a clear recommendation. If the user says `/wait-what` or signals
  confusion, drop jargon immediately and re-explain simply.
- **Two-strike escalation & three-strike circuit breaker**:
  - After **2 failed attempts** on the same failing check or bug with a faster
    coding model, escalate to a stronger reasoning model (with a minimal,
    secret-redacted summary of scope, files, attempts, and exact error output)
    instead of spinning in a retry loop.
  - After **3 failed fix attempts** (locally or in PR review rounds), **STOP**
    patching and ask the human -- repeated failure signals a design or boundary
    problem, not a line-level typo.
- **Pre-PR local gate (`$0`)**: Run `/finish-issue` (`code-reviewer` +
  `artifact-auditor` in parallel, plus repo tests/linters/security scan) before
  proposing a commit or opening a pull request so Round 1 reaches GitHub
  Copilot clean.
- **Protect GitHub Copilot review credits**:
  - Never request a Copilot re-review until **all GitHub CI checks on `HEAD`
    are green**.
  - Do not re-request Copilot on pure default-branch rebases or docs-only edits
    after a clean `No issues found` review.
  - Never treat a Copilot comment mentioning `"unable to review"` or
    `"insufficient remaining quota"` as a passing review.
- **Human-in-the-loop boundaries**: Commits, pushes, PR creation, and releases
  require explicit human approval. **Never merge pull requests (`gh pr merge`)**
  automatically -- the human always owns the merge button.
