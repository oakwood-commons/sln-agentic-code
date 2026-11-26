<p>
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="static/images/scafctl-logo.svg">
    <img src="static/images/scafctl-logo-light.svg" alt="scafctl Configuration Flow" height="44">
  </picture>
</p>

# `sln-agentic-code` -- Agentic Coding Solution & Guided Walkthrough

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)
[![CI](https://github.com/oakwood-commons/sln-agentic-code/actions/workflows/ci.yml/badge.svg)](https://github.com/oakwood-commons/sln-agentic-code/actions/workflows/ci.yml)
[![Docs](https://img.shields.io/badge/docs-github.io-blue)](https://oakwood-commons.github.io/sln-agentic-code/)

Official [`scafctl`](https://github.com/oakwood-commons/scafctl) solution
(`agentic-code`), step-by-step Blueprint walkthrough, repository contracts, and
portable Go CI pipeline from **Oakwood Commons**.

- **Documentation & 5-Step Walkthrough**: <https://oakwood-commons.github.io/sln-agentic-code/>
- **Published OCI Solution**: `ghcr.io/oakwood-commons/solutions/agentic-code`
- **Harnesses supported out of the box**: [OpenCode](https://opencode.ai)
  (`default`) and **Antigravity (`agy`)** (`-r harness=agy` or
  `-r harness=both`), with mapping notes for Claude Code, Cursor, and VS Code +
  GitHub Copilot

---

## Zero-Clone Setup (`scafctl` -- Linux, macOS & Windows)

[`scafctl`](https://github.com/oakwood-commons/scafctl) ships native binaries
for **Linux, macOS, and Windows** (`amd64` and `arm64`) with an embedded
cross-platform shell -- follow the
**[`scafctl` Installation Guide](https://scafctl.oakwood-commons.org/docs/tutorials/getting-started/#installation)**
(`brew install oakwood-commons/tap/scafctl` on macOS/Linux, or `.tar.gz` /
Windows `.zip` binaries from
[GitHub Releases](https://github.com/oakwood-commons/scafctl/releases)). Then
initialize your environment without cloning a repository:

~~~bash
# 1. Default workflow: add starter skills, agents & commands into OpenCode (~/.config/opencode),
#    Antigravity (~/.gemini/config via -r harness=agy), or both (-r harness=both) with non-destructive CEL config merge
scafctl run solution agentic-code --dry-run
scafctl run solution agentic-code
# scafctl run solution agentic-code -r harness=agy

# 2. Opt-in action: verify Go 1.26+ or install it into ~/.local/go across Linux, macOS, and Windows (no sudo/admin)
scafctl run solution agentic-code --action install-go

# 3. Opt-in action: scaffold AGENTS.md, REVIEW.md, Copilot instructions & Go CI into a Go repo
scafctl run solution agentic-code --action init-repo \
  -r repo_dir=/path/to/your-repo \
  -r project_name=your-repo \
  -r module_name=github.com/your-org/your-repo

# 4. Opt-in action: configure GitHub repo settings & rulesets (squash-only PR_TITLE/PR_BODY, signed commits, Copilot review_on_push:false)
scafctl run solution agentic-code --action configure-github \
  -r owner=your-org \
  -r repo=your-repo
~~~

*(Running from a local clone of this repository before a release tag is pushed?
Replace `agentic-code` above with `-f ./scafctl/solution.yaml`.)*

---

## What Is Inside This Repository

| Directory / File | Role |
|---|---|
| `scafctl/solution.yaml` | Published `agentic-code` solution (user OpenCode and/or Antigravity `agy` setup + opt-in `install-go`, `init-repo`, and `configure-github` actions) |
| `scafctl/user-opencode/` | Single-source user-level starter configs (`cloud` default + `ollama-hybrid`), 11 subagents, 11 slash commands, and 22 skills (including `planner`, `create-issue`, `security-audit`, `security-sweep`, `publish-hygiene`, `herdr`, `work-herd`, and `pr-autopilot` -- compiled automatically for `agy` when `-r harness=agy` or `both` is passed) |
| `scafctl/repo/` & `scafctl/go-ci/` | Bundled repo-contract templates (`AGENTS.md`, `REVIEW.md`, `.github/copilot-instructions.md`) and portable Go CI |
| `docs/`, `home/`, `site/` | Hugo (`hugo-book`) 5-step Guided Walkthrough, Advanced Herdr Workflows & Reference site styled with `scafctl`'s Blueprint theme |
| `AGENTS.md`, `REVIEW.md`, `.github/copilot-instructions.md` | Dogfooded **repo-level contracts** defining what this repo requires without forcing personal AI workflows |
| `ci/go-ci/`, `ci/github-actions/`, `ci/tekton/` | Single-source `scafctl` Go CI pipeline plus thin GitHub Actions and Tekton v1 wrappers |
| `examples/go-app/` | Stdlib-only Go URL shortener (`shorty`) used in the hands-on walkthrough and verified on every push |

---

## Local Repository Development

Root [`actions.yaml`](actions.yaml) is auto-discovered by `scafctl` across
Linux, macOS, and Windows:

~~~bash
# Install pinned local tools (hugo extended, gh, go 1.26+, golangci-lint, gosec) into .tool/bin
scafctl run action install-tools

# Run Go CI on examples/go-app, scafctl lint/bundle/functional tests, and Hugo build
scafctl run action ci

# Lint and run functional tests across actions.yaml, ci/go-ci/actions.yaml, and scafctl/solution.yaml
scafctl run action test-solutions

# Preview the Blueprint walkthrough site locally at http://localhost:1313
scafctl run action docs-serve
~~~

---

## License & Credits

Licensed under [Apache-2.0](LICENSE). Third-party MIT notices for adapted
skills are preserved in [`THIRD-PARTY-NOTICE.txt`](THIRD-PARTY-NOTICE.txt).
