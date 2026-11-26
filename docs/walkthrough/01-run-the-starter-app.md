---
title: "1. Run the Starter App"
weight: 10
duration: "~3 MIN HANDS-ON"
goal: "Install the local tools, load the starter commands into OpenCode, and watch the starter URL shortener run before changing a single line."
---

# Step 1: Run the Starter App

## What Usually Goes Wrong

Most AI coding guides make you read ten pages of setup options before you ever
run a program. Worse, teams often check one person's personal AI shortcuts into
the shared Git repo, forcing every teammate to work the exact same way.

## What to Do Instead

First, get the sample Go app running in two commands so you know everything
works. Second, keep two kinds of AI files cleanly separated:

- **[Repo Contract](/docs/reference/glossary.md#repo-contract) (`AGENTS.md`, `REVIEW.md`)** --
  lives inside the Git repository and states *what* the project requires (test
  commands, coverage, Go rules) without forcing how anyone works.
- **Personal AI Workflow (`~/.config/opencode/`)** -- lives in your home folder
  and holds *how* you like to work with your [harness](/docs/reference/glossary.md#harness)
  ([OpenCode](https://opencode.ai)).

Both are installed safely by **[`scafctl`](/docs/reference/glossary.md#scafctl)**,
an open-source workflow and scaffolding CLI used here to run portable project
actions and add missing starter files without stepping on an existing OpenCode
configuration.

---

## Walk Through It

> [!NOTE]
> **Works natively on Linux, macOS, and Windows:** Every tool in this guide runs
> across all three operating systems. Before Step 1, install two CLIs:
> - **[`scafctl`](/docs/reference/glossary.md#scafctl)** *(packaged for Linux,
>   macOS, and Windows on `amd64` + `arm64`)* -- follow the
>   **[`scafctl` Installation Guide](https://scafctl.oakwood-commons.org/docs/tutorials/getting-started/#installation)**
>   (`brew install oakwood-commons/tap/scafctl` on macOS/Linux, or `.tar.gz` /
>   Windows `.zip` binaries from
>   [GitHub Releases](https://github.com/oakwood-commons/scafctl/releases)).
> - **[OpenCode](https://opencode.ai)** *(Linux, macOS, Windows)* -- install
>   from [opencode.ai](https://opencode.ai) and run `opencode auth` once to
>   connect your cloud model provider.
> - **Go 1.26+** -- already installed? You are set. Missing or older Go? Step 1
>   includes a one-line `scafctl` action (`--action install-go`) that installs
>   Go 1.26+ for Linux, macOS, or Windows without `sudo` or admin rights.

{{< steps >}}

1. **Initialize your OpenCode setup (and Go 1.26+ if needed) with `agentic-code`**

    With `scafctl` installed on Linux, macOS, or Windows, initialize your
    OpenCode starter files and make sure Go 1.26+ is ready:

    ~~~shell
    # Preview without writing any files
    scafctl run solution agentic-code --dry-run

    # Add missing starter skills, agents, commands & CEL-merge MCP/policy into ~/.config/opencode
    # (preserves your existing models, providers & plugins, and creates a .bak backup before merging)
    scafctl run solution agentic-code

    # Using Antigravity (agy) instead of (or alongside) OpenCode?
    # Pass -r harness=agy (or -r harness=both) to generate ~/.gemini/config skills, agents, rules & MCP:
    # scafctl run solution agentic-code -r harness=agy

    # Ensure Go 1.26+ is installed (no-op if already 1.26+; otherwise installs into ~/.local/go)
    scafctl run solution agentic-code --action install-go
    ~~~

    > [!TIP]
    > **Already have your own OpenCode or Antigravity (`agy`) config?**
    > - **Smart CEL merge (no manual copy-pasting):** On OpenCode, `scafctl`
    >   parses any existing `opencode.json` or `opencode.jsonc` in place, keeps
    >   100% of your `"model"`, `"providers"`, `"agents"`, `"plugins"`, and
    >   custom skills untouched, and merges in only the missing `"instructions"`,
    >   `scafctl` MCP server, and safety permissions (`gh pr merge`: `deny`) with
    >   a `.bak` backup. With `-r harness=agy` (or `both`), it converts role
    >   agents and standalone commands into the extra `SKILL.md` playbooks `agy`
    >   needs and non-destructively merges `scafctl` into
    >   `~/.gemini/config/mcp_config.json`.
    > - **Want 100% isolation or a later refresh?** Install into a temporary
    >   folder (`scafctl run solution agentic-code -r config_dir=./temp/opencode-config`,
    >   or `-r agy_config_dir=./temp/agy-config` for `agy`) and launch OpenCode
    >   with `OPENCODE_CONFIG_DIR=$PWD/temp/opencode-config opencode`
    >   (`$env:OPENCODE_CONFIG_DIR="$PWD\temp\opencode-config"; opencode` in
    >   PowerShell) -- OpenCode still reuses your existing `opencode auth` logins
    >   while leaving `~/.config/opencode` untouched. To refresh bundled skills
    >   after a future release, pass `-r overwrite=true` (`.bak` backups are
    >   saved automatically).

    **What you should see:** `agentic-code` adds the starter skills (`planner`,
    `create-issue`, `stepwise-interview`, `tdd`, `go`, `security-audit`,
    `finish-issue`), subagents, and commands (browse all playbooks and see how
    prompts, rules, skills, subagents, and commands differ in
    [Skills, Agents & Commands](/docs/reference/skills-and-agents.md)), verifies
    `gh` CLI, and confirms Go 1.26+ is ready.

2. **Clone this repository to get the `shorty` practice app & local tools**

    Now clone `sln-agentic-code` so you have the `shorty` starter Go app on your
    machine. `scafctl` reads [`actions.yaml`](https://github.com/oakwood-commons/sln-agentic-code/blob/main/actions.yaml),
    detects whether you are on Linux, macOS, or Windows, and verifies or
    downloads the pinned repo tools (`go` 1.26+, `hugo`, `gh`, `golangci-lint`,
    and `gosec`) into `.tool/` with no `sudo` or admin rights required:

    ~~~shell
    git clone https://github.com/oakwood-commons/sln-agentic-code.git
    cd sln-agentic-code

    scafctl run action install-tools
    ~~~

    *(Working offline or testing uncommitted edits to the starter? Run
    `scafctl run solution -f ./scafctl/solution.yaml` from inside the clone.)*

    **What you should see:** `hugo`, `gh`, `golangci-lint`, and `gosec` are
    verified or downloaded into `.tool/bin` in a few seconds.

3. **Run the starter Go tests & try `shorty` live**

    Verify that the sample Go URL shortener in `examples/go-app/` passes its
    starter checks:

    ~~~shell
    scafctl run action go-ci
    ~~~

    Notice the `87%+` test coverage already guarding `POST /links` and
    `GET /{code}`. Those existing behavior tests are why you can let AI add new
    features in Step 3 with confidence -- every test run proves existing links
    and error cases still work without regressions. *(Adopting this workflow on
    an older codebase with weak test coverage? Make adding public-behavior tests
    your very first AI task before building new features.)*

    Now start `shorty` in one terminal tab and shorten a URL from a second tab
    (click **Bash** or **PowerShell** below if your shell differs -- your choice
    sticks across the rest of the guide):

    ~~~shell
    # Terminal 1: start the URL shortener on port 8080
    cd examples/go-app
    go run ./cmd/shorty -addr :8080
    ~~~

    {{< shell-toggle >}}
    {{< shell-pane "bash" >}}
    # Terminal 2: create a short link, then follow the redirect headers
    curl -s -X POST http://localhost:8080/links \
      -H "Content-Type: application/json" \
      -d '{"url":"https://github.com/oakwood-commons/scafctl"}'

    # Replace <code> below with the 6-character code printed above:
    curl -i http://localhost:8080/<code>
    {{< /shell-pane >}}
    {{< shell-pane "powershell" >}}
    # Terminal 2: create a short link, then inspect the 302 redirect headers
    $link = Invoke-RestMethod -Method Post -Uri http://localhost:8080/links `
      -ContentType "application/json" `
      -Body '{"url":"https://github.com/oakwood-commons/scafctl"}'
    $link

    # Inspect the 302 Found + Location header without following the redirect:
    curl.exe -i "http://localhost:8080/$($link.code)"
    {{< /shell-pane >}}
    {{< /shell-toggle >}}

    **What you should see:** `POST /links` returns `{"code":"...","url":"..."}`,
    and `GET /<code>` returns `HTTP/1.1 302 Found` with a `Location:` header
    pointing back to the original URL. Press `Ctrl+C` in Terminal 1 when done.

{{< /steps >}}

---

## Why `shorty` Needs a New Feature

`shorty` works, and the whole server is ~150 lines of plain Go in
[`examples/go-app/internal/shortener/shortener.go`](https://github.com/oakwood-commons/sln-agentic-code/blob/main/examples/go-app/internal/shortener/shortener.go).
Right now it has two obvious gaps:

1. **Every link lives forever** -- even temporary promo or invite links.
2. **There is no click counter** -- you cannot check whether anyone ever used
   your link.

In **Step 2**, you will have OpenCode scope **optional link expiration and click
counts** -- starting with a 2-minute conversation instead of a blind code dump.

- *Using OpenAI, Gemini, Copilot, or local Ollama instead of Anthropic? See
  [Customizing Models & Plugins](/docs/reference/customize-and-credits.md).*
