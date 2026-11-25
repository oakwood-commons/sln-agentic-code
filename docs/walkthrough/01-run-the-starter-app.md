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

## What We Do Instead

First, we get the sample Go app running in two commands so you know everything
works. Second, we keep two kinds of AI files cleanly separated:

- **[Repo Contract](/docs/reference/glossary.md#repo-contract) (`AGENTS.md`, `REVIEW.md`)** --
  lives inside the Git repository and states *what* the project requires (test
  commands, coverage, Go rules) without forcing how anyone works.
- **Personal AI Workflow (`~/.config/opencode/`)** -- lives in your home folder
  and holds *how* you like to work with your [harness](/docs/reference/glossary.md#harness)
  ([OpenCode](https://opencode.ai)).

Both are installed safely by **[`scafctl`](https://github.com/oakwood-commons/scafctl)**,
a small CLI tool we use to run portable project actions and copy starter files
into place (creating `.bak` backups if you already have config files).

---

## Walk Through It

> [!NOTE]
> **Before you start:** Make sure you have **Go 1.26+** (`go version`),
> **`scafctl`** (`brew install oakwood-commons/tap/scafctl` or from
> [GitHub Releases](https://github.com/oakwood-commons/scafctl/releases)), and
> **[OpenCode](https://opencode.ai)** installed (`opencode auth` to connect your
> preferred cloud model provider).

{{% steps %}}

1. **Initialize your personal OpenCode setup with `agentic-code`**

    With `scafctl` installed, you can initialize `~/.config/opencode` from
    anywhere in one command (preview first with `--dry-run`):

    ~~~bash
    # Preview without writing any files
    scafctl run solution agentic-code --dry-run

    # Install starter commands, agents, skills & scafctl MCP into ~/.config/opencode
    # (backs up any existing modified files with .bak)
    scafctl run solution agentic-code
    ~~~

    **What you should see:** `agentic-code` completes, placing `opencode.jsonc`
    and the starter `/interview`, `/wait-what`, `/plan`, and `/finish-issue`
    commands in `~/.config/opencode/`.

2. **Clone this repository to get the `shorty` practice app & local tools**

    Now clone `sln-agentic-code` so you have the `shorty` starter Go app on your
    laptop. `scafctl` reads [`actions.yaml`](https://github.com/oakwood-commons/sln-agentic-code/blob/main/actions.yaml),
    detects whether you are on Linux, macOS, or Windows, and downloads the
    pinned Go linters into `.tool/bin` (no `sudo` required):

    ~~~bash
    git clone https://github.com/oakwood-commons/sln-agentic-code.git
    cd sln-agentic-code

    scafctl run action install-tools
    ~~~

    *(Working offline or testing uncommitted edits to the starter? Run
    `scafctl run solution -f ./scafctl/solution.yaml` from inside the clone.)*

    **What you should see:** `hugo`, `golangci-lint`, and `gosec` are verified
    or downloaded into `.tool/bin` in a few seconds.

3. **Run the starter Go tests & try `shorty` live**

    Verify that our sample Go URL shortener in `examples/go-app/` passes its
    starter checks:

    ~~~bash
    scafctl run action go-ci
    ~~~

    Now start `shorty` in one terminal tab and shorten a URL from a second tab:

    ~~~bash
    # Terminal 1: start the URL shortener on port 8080
    cd examples/go-app
    go run ./cmd/shorty -addr :8080
    ~~~

    ~~~bash
    # Terminal 2: create a short link, then follow the redirect headers
    curl -s -X POST http://localhost:8080/links \
      -H "Content-Type: application/json" \
      -d '{"url":"https://github.com/oakwood-commons/scafctl"}'

    # Replace <code> below with the 6-character code printed above:
    curl -i http://localhost:8080/<code>
    ~~~

    **What you should see:** `POST /links` returns `{"code":"...","url":"..."}`,
    and `GET /<code>` returns `HTTP/1.1 302 Found` with a `Location:` header
    pointing back to the original URL. Press `Ctrl+C` in Terminal 1 when done.

{{% /steps %}}

---

## Why We Need a New Feature

`shorty` works, and the whole server is ~150 lines of plain Go in
[`examples/go-app/internal/shortener/shortener.go`](https://github.com/oakwood-commons/sln-agentic-code/blob/main/examples/go-app/internal/shortener/shortener.go).
Right now it has two obvious gaps:

1. **Every link lives forever** -- even temporary promo or invite links.
2. **There is no click counter** -- you cannot check whether anyone ever used
   your link.

In **Step 2**, we will ask OpenCode to add **optional link expiration
(`expires_in`) and click stats (`GET /links/{code}/stats`)** -- starting with a
2-minute conversation instead of a blind code dump.

- *Using OpenAI, Gemini, Copilot, or local Ollama instead of Anthropic? See
  [Customizing Models & Plugins](/docs/reference/customize-and-credits.md).*
