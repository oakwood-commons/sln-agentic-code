---
title: "Adopt in Your Own Repo"
weight: 20
---

# Adopt in Your Own Go Repository

Once you have tried the walkthrough on `shorty`, you can set up any of your own
Go repositories with the exact same repo contracts, portable Go CI, and GitHub
settings using the published `agentic-code` solution -- **no repository clone
required**.

---

## 1. Scaffold `AGENTS.md`, `REVIEW.md` & Go CI (`init-repo`)

Run `init-repo` against your Go project directory:

~~~bash
# Preview the files that will be created
scafctl run solution agentic-code --action init-repo \
  -r repo_dir=/path/to/your-service \
  -r project_name=your-service \
  -r module_name=github.com/your-org/your-service \
  -r min_coverage=70 \
  --dry-run

# Write the files (creates .bak backups if any target file already exists)
scafctl run solution agentic-code --action init-repo \
  -r repo_dir=/path/to/your-service \
  -r project_name=your-service \
  -r module_name=github.com/your-org/your-service \
  -r min_coverage=70
~~~

*(Running from a local clone of `sln-agentic-code`? Replace `agentic-code` with
`-f ./scafctl/solution.yaml`.)*

What gets added to `/path/to/your-service`:
- **`AGENTS.md`** -- build, test, lint, and Go conventions for contributors and
  coding agents.
- **`REVIEW.md`** -- single-pass PR review contract (focusing on bugs, callers,
  races, and `<last_reviewed_commit>..HEAD` delta re-reviews).
- **`.github/copilot-instructions.md`** -- points GitHub Copilot at `AGENTS.md`
  and `REVIEW.md` without duplicating rules.
- **`ci/go-ci/actions.yaml`** -- portable Go CI workflow (pass `-r
  include_go_ci=false` if your repo already has its own CI actions file).

---

## 2. Configure GitHub Repository Settings & Rulesets (`configure-github`)

Manually clicking through GitHub settings on every new repository is tedious and
easy to misconfigure (especially leaving `review_on_push` enabled by mistake).

Authenticate `scafctl` with GitHub once, preview with `--dry-run`, and apply:

~~~bash
# Log in once (uses scafctl's GitHub auth plugin)
scafctl auth login github

# Preview API operations
scafctl run solution agentic-code --action configure-github \
  -r owner=your-org \
  -r repo=your-service \
  --dry-run

# Apply idempotently (safe to re-run anytime -- updates existing rulesets in place)
scafctl run solution agentic-code --action configure-github \
  -r owner=your-org \
  -r repo=your-service
~~~

> [!NOTE]
> **Personal GitHub accounts vs. GitHub Organizations:** Organization repos
> default to `required_approvals=1` (with an `OrganizationAdmin` bypass). For a
> solo personal repository where you are the only contributor, pass `-r
> required_approvals=0` so you do not block yourself from merging your own PRs.

### What `configure-github` Sets Up

| Setting / Ruleset | Why It Matters for Agentic Workflows |
|---|---|
| **Squash-only merges (`PR_TITLE` + `PR_BODY`)** | Every PR lands on `main` as one clean Conventional Commit built from the PR title and body (which `/pr-create` keeps synced to `HEAD`), and deletes merged feature branches automatically. |
| **`Copilot review for default branch`** | Protects `main` against force-push/deletion and configures GitHub Copilot Code Review with **`review_on_push: false`** and **`review_draft_pull_requests: false`** so you only spend review quota on ready-for-review PRs. |
| **`all-branches-signed-commits`** | Requires cryptographic commit signatures (`git commit -s -S`) across all branches. Pass `-r require_signed_commits=false` if your team has not set up SSH/GPG signing yet. |
| **`default-branch` PR protection** | Requires pull requests, **`required_review_thread_resolution: true`** (every review comment must be replied to and resolved before merge), and stale-review dismissal. |
| **Security & optional GitHub Pages** | Enables Dependabot vulnerability alerts, and sets GitHub Pages to GitHub Actions `workflow` mode when `-r enable_pages=true` is passed. |
