---
name: publish-hygiene
description: "Use before posting or committing to any public repository (issues, PRs, commits, docs, or skills) to scrub secrets, local home paths (/home/<user>), private hostnames, internal tool names, and non-public organization details."
---

# Public Repository Publish Hygiene

Contributors work from many different employers, personal machines, and local
tooling setups. None of that private context belongs in anything published to a
public repository: GitHub issues, PR titles and bodies, PR or issue comments,
commit messages, or committed docs, examples, tests, and config files.

This checklist is intentionally **organization-agnostic** so you can apply it
regardless of which machine or employer environment you are working from.

## When to Run This Check

Run before any of these actions on a public repository:
- Creating or commenting on a GitHub issue (`create-issue`, `github-issues`)
- Opening a pull request or updating its title or body (`pr-create`, `pr-autopilot`)
- Writing a commit message or running `finish-issue`
- Adding or editing committed docs, examples, fixtures, skills, or configs

## What to Scrub

- **Credentials and secrets:** API keys, tokens, private keys, cloud project
  IDs tied to private accounts, or connection strings.
- **User-specific filesystem paths:** `/home/<username>`, `/Users/<username>`,
  `C:\Users\<username>`, or WSL mount paths baked into examples or logs.
- **Internal infrastructure:** private domains, internal hostnames, cluster or
  environment names, internal registries, VPN names, or private IP ranges.
- **Organization & identity details:** non-public company or team names,
  internal product codenames, corporate email domains, SSO tenant IDs, or
  employee IDs.
- **Screenshots, logs, and pasted output:** inspect visually as well as with
  grep, since command output frequently bakes in hostnames, usernames, or
  internal URLs.

## How to Scrub Without Losing Technical Meaning

Replace identifying details with neutral placeholders while keeping the
technical substance reproducible:
- Specific company or org -> `example-corp` or `Acme`
- Internal domain or registry -> `example.com`, `internal.example.com`, `registry.example.com`
- Local home path -> `~/my-repo` or `/path/to/repo`
- Internal tool codename -> a plain description of its role (for example,
  "the internal CI runner")

## Quick Self-Check Before Publishing

1. Would every line make sense to an external reader with zero knowledge of my
   machine or employer?
2. Are there any `/home/`, `/Users/`, internal domains, or non-public project
   IDs in diffs, examples, or command outputs?
3. Are there any real tokens or credentials anywhere in the text?

If any answer raises doubt, scrub before committing or posting, because git
history retains deleted lines unless rewritten.
