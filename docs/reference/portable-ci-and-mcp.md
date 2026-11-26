---
title: "Portable Go CI & scafctl MCP"
weight: 30
---

# Portable Go CI & `scafctl` MCP Reference

How [`oakwood-commons/scafctl`](https://github.com/oakwood-commons/scafctl)
powers both live AI schema lookups (`scafctl mcp serve`) and a single-source Go
CI workflow across local terminals, GitHub Actions, and Tekton.

---

## 1. One Go CI Definition, Three Environments

Instead of maintaining separate bash scripts inside `.github/workflows/`, Tekton
`Task` YAML, and a local `Makefile`,
[`ci/go-ci/actions.yaml`](https://github.com/oakwood-commons/sln-agentic-code/blob/main/ci/go-ci/actions.yaml)
defines the entire Go action [DAG](/docs/reference/glossary.md#dag) once behind
simple `-r key=value` parameters:

<div class="bp-flow" role="region" aria-label="One portable scafctl Go CI file used across three environments">
  <div class="bp-flow-row cols-3-plain">
    <div class="bp-flow-node">
      <span class="bp-flow-step">ENV 01 // LAPTOP</span>
      <span class="bp-flow-title">Local Terminal</span>
      <span class="bp-flow-desc"><code>scafctl run action ci -f ci/go-ci/actions.yaml</code></span>
    </div>
    <div class="bp-flow-node">
      <span class="bp-flow-step">ENV 02 // GITHUB</span>
      <span class="bp-flow-title">GitHub Actions</span>
      <span class="bp-flow-desc"><code>ci/github-actions/ci.yml</code></span>
    </div>
    <div class="bp-flow-node">
      <span class="bp-flow-step">ENV 03 // KUBERNETES</span>
      <span class="bp-flow-title">Tekton v1 Pipeline</span>
      <span class="bp-flow-desc"><code>ci/tekton/go-ci-{task,pipeline}.yaml</code></span>
    </div>
  </div>
  <div class="bp-flow-arrow-v" aria-hidden="true"><svg viewBox="0 0 24 24" fill="none" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 5v14M6 13l6 6 6-6"/></svg><span class="bp-flow-cond">All three invoke the same DAG</span><svg viewBox="0 0 24 24" fill="none" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 5v14M6 13l6 6 6-6"/></svg></div>
  <div class="bp-flow-node is-box">
    <span class="bp-flow-step" style="text-align:center;">SINGLE SOURCE // ci/go-ci/actions.yaml</span>
    <div class="bp-flow-node is-center">
      <span class="bp-flow-step">SETUP</span>
      <span class="bp-flow-title"><code>install-tools</code></span>
      <span class="bp-flow-desc">Detects OS/arch (Linux, macOS, Windows) and ensures Go 1.26+, <code>golangci-lint</code> &amp; <code>gosec</code>.</span>
    </div>
    <div class="bp-flow-arrow-v" aria-hidden="true"><svg viewBox="0 0 24 24" fill="none" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 5v14M6 13l6 6 6-6"/></svg><span class="bp-flow-cond">Parallel deterministic gates</span><svg viewBox="0 0 24 24" fill="none" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 5v14M6 13l6 6 6-6"/></svg></div>
    <div class="bp-flow-pills">
      <span>fmt-check (gofmt)</span>
      <span>vet (go vet)</span>
      <span>lint (golangci-lint)</span>
      <span>gosec (security)</span>
      <span>test (-race + coverage)</span>
    </div>
    <div class="bp-flow-arrow-v" aria-hidden="true"><svg viewBox="0 0 24 24" fill="none" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 5v14M6 13l6 6 6-6"/></svg></div>
    <div class="bp-flow-node is-center is-mint">
      <span class="bp-flow-step">FINAL ARTIFACT</span>
      <span class="bp-flow-title"><code>build</code> (<code>go build -trimpath</code>)</span>
      <span class="bp-flow-desc">Compiles only when every lint, security, race, and coverage check passes.</span>
    </div>
  </div>
</div>

### Supported Parameters (`-r key=value`)

| Parameter | Default | Description |
|---|---|---|
| `workdir` | `"."` | Go module directory to test and build |
| `tools_bin` | `"<cwd>/.tool/bin"` | Directory where missing CI tools (`go`, `golangci-lint`, `gosec`) are installed (relative overrides normalize to absolute paths) |
| `build_pkg` | `"./..."` | Package pattern passed to `go build` |
| `output_bin` | `""` | Optional binary output path (`go build -o`) |
| `min_coverage` | `70` | Minimum total statement coverage (`0` disables threshold) |
| `go_version` | `"1.26.11"` | Pinned Go version installed into `tools_bin/go-sdk` if Go 1.26+ is missing |
| `golangci_lint_version` | `"2.15.0"` | Pinned `golangci-lint` version installed into `tools_bin` if missing |
| `gosec_version` | `"2.30.0"` | Pinned `gosec` security scanner version installed into `tools_bin` if missing |
| `race` | `true` | Run `go test -race` when a C compiler and supported OS/arch are present |
| `run_lint` | `true` | Run pinned `golangci-lint` |
| `run_gosec` | `true` | Run pinned `gosec` security scanner |

### Copy-Paste Adapters

{{< tabs >}}
{{< tab "Local CLI" >}}

~~~shell
# Run full Go CI
scafctl run action ci -f ./ci/go-ci/actions.yaml -r workdir=examples/go-app

# Run only lint and test during inner-loop coding
scafctl run action lint test -f ./ci/go-ci/actions.yaml -r workdir=examples/go-app

# Run the CI file's own sandboxed functional tests
scafctl test functional -f ./ci/go-ci/actions.yaml
~~~

{{< /tab >}}
{{< tab "GitHub Actions" >}}

Copy [`ci/github-actions/ci.yml`](https://github.com/oakwood-commons/sln-agentic-code/blob/main/ci/github-actions/ci.yml)
to `.github/workflows/ci.yml` in your Go repository:

~~~yaml
- name: Run portable Go CI via scafctl
  run: |
    scafctl run action ci -f ./ci/go-ci/actions.yaml \
      -r workdir=. \
      -r build_pkg=./... \
      -r min_coverage=70
~~~

{{< /tab >}}
{{< tab "Tekton" >}}

Apply [`ci/tekton/go-ci-task.yaml`](https://github.com/oakwood-commons/sln-agentic-code/blob/main/ci/tekton/go-ci-task.yaml)
and [`ci/tekton/go-ci-pipeline.yaml`](https://github.com/oakwood-commons/sln-agentic-code/blob/main/ci/tekton/go-ci-pipeline.yaml):

~~~shell
kubectl apply -f ci/tekton/go-ci-task.yaml
kubectl apply -f ci/tekton/go-ci-pipeline.yaml
~~~

{{< /tab >}}
{{< /tabs >}}

> [!NOTE]
> **Splitting large pipelines with `compose:`:** When a `scafctl` solution
> outgrows a single file, add `compose: [resolvers.yaml, actions.yaml,
> tests.yaml]` so `scafctl` merges partial files cleanly at load time (`scafctl
> render solution --effective`).

---

## 2. Wiring `scafctl mcp serve` for AI Authoring

When AI agents write YAML from memory, people say they "hallucinate" fields --
when really they just lack schema [context](/docs/reference/glossary.md#context)
and a [deterministic](/docs/reference/glossary.md#deterministic) validator to
test their drafts against. `scafctl` solves both with its built-in Model Context
Protocol ([MCP](/docs/reference/glossary.md#mcp)) server (`scafctl mcp serve`),
already enabled in `scafctl/user-opencode/opencode.*.jsonc`:

~~~json
{
  "mcp": {
    "servers": {
      "scafctl": {
        "type": "local",
        "command": ["scafctl", "mcp", "serve"]
      }
    }
  }
}
~~~

*(Use the same `scafctl mcp serve` stdio command in `.vscode/mcp.json`,
`.cursor/mcp.json`, or `claude_desktop_config.json`.)*

### Most Useful `scafctl` MCP Tools

| When You Ask the AI To... | MCP Tool It Calls |
|---|---|
| Check valid fields on `solution`, `resolver`, or `workflow` | `explain_kind` / `get_solution_schema` |
| Inspect exact inputs/outputs of `exec`, `file`, `metadata`, or `github` | `get_provider_schema` |
| Test a CEL expression or Go template before saving | `evaluate_cel` / `validate_expression` |
| Lint a solution file for schema or dependency issues | `lint_solution` |
| Dry-run action order and resolved inputs safely | `preview_resolvers` / `preview_action` |
| Run `spec.testing.cases` in isolated temp sandboxes | `run_solution_tests` |

---

## 3. Connecting Jira User Stories via MCP {#jira-mcp}

If your team tracks work in **Jira** rather than GitHub Issues, connect
Atlassian's official OAuth 2.1 MCP server alongside `scafctl` in
`~/.config/opencode/opencode.jsonc`:

~~~jsonc
{
  "mcp": {
    "servers": {
      "scafctl": {
        "type": "local",
        "command": ["scafctl", "mcp", "serve"]
      },
      "atlassian": {
        "type": "remote",
        "url": "https://mcp.atlassian.com/v2/mcp"
      }
    }
  }
}
~~~

Once authenticated via your browser (`opencode mcp auth atlassian` if not
prompted automatically), the `planner` and `git-github-workflow` skills work
with Jira stories in plain English just like GitHub or local issues:
- **Drafting a story:** *"Let's create a Jira story in PROJ for optional
  expiring promo links..."*
- **Skeptical triage & planning:** *"Have a look at PROJ-123 and come up with a
  detailed plan to implement it."* (`planner` fetches `PROJ-123` over MCP,
  verifies every claim against live `file:line` code, and pairs naturally with
  `<type>/PROJ-123-<slug>` branches.)
