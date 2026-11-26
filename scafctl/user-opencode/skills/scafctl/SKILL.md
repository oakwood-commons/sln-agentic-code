---
name: scafctl
description: "Use when authoring, editing, testing, packaging, or running scafctl Solution and actions.yaml files."
---

# `scafctl` Solution & Action Authoring

Consult `scafctl mcp serve` tools rather than guessing schema fields, provider
inputs, or CEL functions from memory.

## Two Solution Archetypes

1. **Repository task file (`actions.yaml`):**
   - Checked into a repo root and auto-discovered by `scafctl run action <name>`
     (no `-f` needed inside that directory).
2. **Published catalog solution (`solution.yaml` / `scafctl/solution.yaml`):**
   - Packaged via `scafctl package solution` and published to an OCI catalog
     (e.g., `ghcr.io/<org>/solutions/<name>` following the `sln-<name>` repo
     convention).
   - **Design a real default workflow (`explicit: false`):** `scafctl run
     solution <name>` executes all non-explicit actions out of the box with zero
     flags.
   - **Mark secondary/admin actions `explicit: true`:** opt-in actions (like
     `init-repo`, `configure-github`, `docs-serve`, or `cleanup`) set
     `explicit: true` so users invoke them selectively via
     `scafctl run solution <name> --action <action>` (or `-a <action>`).
   - **Keep bundled files under the solution directory:** static bundle
     discovery rejects `../` parent-directory paths (`path escapes bundle
     root`). Declare non-obvious globs under `bundle.include`.

## Authoring & Verification Loop

1. **Schema & providers first:** Call `explain_kind` and `get_provider_schema`
   before writing resolvers or actions. Omit `type:` on resolvers whose provider
   returns an object/map.
2. **Cross-platform `exec` actions (`mvdan/sh`):**
   - Check OS/arch via the `metadata` provider (`_.platform.os`,
     `_.platform.arch`).
   - On Windows, use `;` as `PATH` separator, append `.exe` to binary names, and
     strip trailing `\r` (`CR="$(printf '\r')"`) when parsing command output.
   - Note that when a `continueOnError: true` action fails, downstream
     `dependsOn` actions are skipped (`dependency-failed`).
3. **Deterministic validation before finishing:**
   - `scafctl lint -f <file>` (or `lint_solution` MCP tool) with zero errors or
     stale dependency warnings.
   - `scafctl test functional -f <file>` (or `run_solution_tests` MCP tool),
     where every test case in `spec.testing.cases` carries `description` and
     `assertions` (or `snapshot`).
   - For catalog solutions (`bundle:` present), run
     `scafctl package solution -f <file> --dry-run --no-vendor` (omit
     `--no-vendor` when `bundle.plugins` are vendored locally) to prove bundle
     completeness before committing.
