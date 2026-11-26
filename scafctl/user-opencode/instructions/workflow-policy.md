# Personal Agentic Workflow Policy

- **Repo contract vs. personal workflow**: Always read the target repository's
  `AGENTS.md` and `REVIEW.md` first for repo-specific build, test, lint,
  architecture, and **issue-tracking** rules (`github`, `jira`, `local`
  `temp/issues/`, or a combination). Never make the user repeat where issues
  live in every prompt -- if a repo has not declared an issue mechanism yet, ask
  once and offer to record it in `AGENTS.md`.
- **Auto-select skills from plain-English prose**: Users should not have to
  memorize slash commands or repeat boilerplate instructions. Match their
  natural prose to your installed skills automatically:
  - Load `planner` before starting any non-trivial feature, refactor, issue
    draft, or issue pickup -- inspect the codebase first, skeptically verify it
    is worth doing, and align on a plan before writing code.
  - Load `create-issue` when drafting or filing issues/stories so every issue
    has a plain-English top half for human sign-off and grounded repo context at
    the bottom for AI implementation (shown in chat before filing).
  - Load `stepwise-interview` when product edge cases are open, and `wait-what`
    whenever the user says `"what?"` or `"I don't follow"`.
  - Load `tdd` and `codebase-design` (plus `go` in Go repos) when implementing
    incremental steps, `security-audit` (or `security-sweep` for whole-repo
    audits) when investigating or fixing security flaws, `publish-hygiene`
    before publishing to public repos, `finish-issue` automatically before
    opening a pull request, `pr-autopilot` when driving an open PR through CI
    and Copilot reviews, and `work-herd` when asked for a Herdr plan/DAG or
    multi-worker wave.
- **Engineering judgment (best over easiest, never rubber-stamp *how*)**:
  - Honor libraries, stack choices, and constraints the user already stated
    without grilling them backwards on settled basics.
  - Check how mature tools and idiomatic packages solve the problem, and
    recommend the **best long-term design** (`codebase-design`) rather than the
    easiest shortcut.
  - Do not just build what you think the user wants to hear: if a prompt
    prescribes a clunky or fragile *"how"* when a simpler, safer, or more
    testable design exists, explain the better approach up front before coding.
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
- **Pre-PR AI self-review gate (0 CI / 0 PR quota)**: Automatically run
  `finish-issue` (`code-reviewer` + `artifact-auditor` in parallel, plus free
  local tests/linters/security scan) on your feature branch before pushing and
  opening a pull request so Round 1 reaches GitHub Copilot and human PR review
  clean.
- **Protect GitHub Copilot review credits**:
  - Never request a Copilot re-review until **all GitHub CI checks on `HEAD`
    are green**.
  - Do not re-request Copilot on pure default-branch rebases or docs-only edits
    after a clean `No issues found` review.
  - Never treat a Copilot comment mentioning `"unable to review"` or
    `"insufficient remaining quota"` as a passing review.
- **Human-in-the-loop boundaries**: Never commit directly to `main`/`master`;
  work on feature branches, require human approval before pushing to remotes,
  opening PRs, or cutting releases, and **never merge pull requests (`gh pr
  merge`)** automatically -- the human owns architecture, guidance, and final PR
  acceptance.
