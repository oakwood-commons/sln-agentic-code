---
title: "Portable Go CI & scafctl MCP"
weight: 30
---

# Portable Go CI & `scafctl` MCP Reference

How [`oakwood-commons/scafctl`](https://github.com/oakwood-commons/scafctl)
powers both live AI schema lookups (`scafctl mcp serve`) and our single-source
Go CI workflow across local terminals, GitHub Actions, and Tekton.

---

## 1. One Go CI Definition, Three Environments

Instead of maintaining separate bash scripts inside `.github/workflows/`, Tekton
`Task` YAML, and a local `Makefile`,
[`ci/go-ci/actions.yaml`](https://github.com/oakwood-commons/sln-agentic-code/blob/main/ci/go-ci/actions.yaml)
defines the entire Go pipeline once behind simple `-r key=value` parameters:

{{< mermaid >}}
flowchart TD
    A1["Local Terminal<br/>scafctl run action ci -f ci/go-ci/actions.yaml"] --> S
    A2["GitHub Actions<br/>ci/github-actions/ci.yml"] --> S
    A3["Tekton v1 Task & Pipeline<br/>ci/tekton/go-ci-{task,pipeline}.yaml"] --> S

    subgraph S["ci/go-ci/actions.yaml"]
        IT["install-tools<br/>(OS/arch -> golangci-lint & gosec)"]
        FC["fmt-check (gofmt)"]
        VT["vet (go vet)"]
        LN["lint (golangci-lint)"]
        GS["gosec (security scan)"]
        TS["test (go test -race + coverage)"]
        BD["build (go build -trimpath)"]

        IT --> LN
        IT --> GS
        FC --> BD
        VT --> BD
        LN --> BD
        GS --> BD
        TS --> BD
    end
{{< /mermaid >}}

### Supported Parameters (`-r key=value`)

| Parameter | Default | Description |
|---|---|---|
| `workdir` | `"."` | Go module directory to test and build |
| `build_pkg` | `"./..."` | Package pattern passed to `go build` |
| `output_bin` | `""` | Optional binary output path (`go build -o`) |
| `min_coverage` | `70` | Minimum total statement coverage (`0` disables threshold) |
| `race` | `true` | Run `go test -race` when a C compiler and supported OS/arch are present |
| `run_lint` | `true` | Run pinned `golangci-lint` |
| `run_gosec` | `true` | Run pinned `gosec` security scanner |

### Copy-Paste Adapters

{{< tabs "ci-environment-adapters" >}}
{{% tab "Local CLI" %}}

~~~bash
# Run full Go CI
scafctl run action ci -f ./ci/go-ci/actions.yaml -r workdir=examples/go-app

# Run only lint and test during inner-loop coding
scafctl run action lint test -f ./ci/go-ci/actions.yaml -r workdir=examples/go-app

# Run the CI file's own sandboxed functional tests
scafctl test functional -f ./ci/go-ci/actions.yaml
~~~

{{% /tab %}}
{{% tab "GitHub Actions" %}}

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

{{% /tab %}}
{{% tab "Tekton" %}}

Apply [`ci/tekton/go-ci-task.yaml`](https://github.com/oakwood-commons/sln-agentic-code/blob/main/ci/tekton/go-ci-task.yaml)
and [`ci/tekton/go-ci-pipeline.yaml`](https://github.com/oakwood-commons/sln-agentic-code/blob/main/ci/tekton/go-ci-pipeline.yaml):

~~~bash
kubectl apply -f ci/tekton/go-ci-task.yaml
kubectl apply -f ci/tekton/go-ci-pipeline.yaml
~~~

{{% /tab %}}
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
