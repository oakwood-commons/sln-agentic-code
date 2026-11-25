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
- **Primary harness covered**: [OpenCode](https://opencode.ai) (with clear
  adaptation notes for Antigravity/`agy`, Claude Code, Cursor, and VS Code +
  GitHub Copilot)

---

## Zero-Clone Setup (`scafctl`)

Install [`scafctl`](https://github.com/oakwood-commons/scafctl) (`brew install
oakwood-commons/tap/scafctl` or from
[GitHub Releases](https://github.com/oakwood-commons/scafctl/releases)), then
initialize your environment in one command without cloning a repository:

~~~bash
# 1. Default workflow: install personal OpenCode skills, agents, commands & MCP into ~/.config/opencode
scafctl run solution agentic-code --dry-run
scafctl run solution agentic-code

# 2. Opt-in action: scaffold AGENTS.md, REVIEW.md, Copilot instructions & Go CI into a Go repo
scafctl run solution agentic-code --action init-repo \
  -r repo_dir=/path/to/your-repo \
  -r project_name=your-repo \
  -r module_name=github.com/your-org/your-repo

# 3. Opt-in action: configure GitHub repo settings & rulesets (squash-only PR_TITLE/PR_BODY, signed commits, Copilot review_on_push:false)
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
| `scafctl/solution.yaml` | Published `agentic-code` solution (default user OpenCode setup + opt-in `init-repo` and `configure-github` actions) |
| `scafctl/user-opencode/` | Agnostic **user-level** OpenCode configs (`cloud` default + `ollama-hybrid`), core subagents, slash commands, and 14 skills |
| `scafctl/repo/` & `scafctl/go-ci/` | Bundled repo-contract templates (`AGENTS.md`, `REVIEW.md`, `.github/copilot-instructions.md`) and portable Go CI |
| `docs/`, `home/`, `site/` | Hugo (`hugo-book`) 5-step Guided Walkthrough & Reference site styled with `scafctl`'s Blueprint theme |
| `AGENTS.md`, `REVIEW.md`, `.github/copilot-instructions.md` | Dogfooded **repo-level contracts** defining what this repo requires without forcing personal AI workflows |
| `ci/go-ci/`, `ci/github-actions/`, `ci/tekton/` | Single-source `scafctl` Go CI pipeline plus thin GitHub Actions and Tekton v1 wrappers |
| `examples/go-app/` | Stdlib-only Go URL shortener (`shorty`) used in the hands-on walkthrough and verified on every push |

---

## Local Repository Development

Root [`actions.yaml`](actions.yaml) is auto-discovered by `scafctl` across
Linux, macOS, and Windows:

~~~bash
# Install pinned local tools (hugo extended, golangci-lint, gosec) into .tool/bin
scafctl run action install-tools

# Run Go CI on examples/go-app, scafctl lint/bundle/functional tests, and Hugo build
scafctl run action ci

# Preview the Blueprint walkthrough site locally at http://localhost:1313
scafctl run action docs-serve
~~~

---

## License & Credits

Licensed under [Apache-2.0](LICENSE). Third-party MIT and SIL OFL 1.1 notices
for adapted skills and font outlines are preserved in
[`THIRD-PARTY-NOTICE.txt`](THIRD-PARTY-NOTICE.txt).
