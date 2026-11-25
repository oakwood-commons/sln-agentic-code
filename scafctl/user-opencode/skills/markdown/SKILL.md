---
name: markdown
description: "Markdown authoring rules: tilde fences around nested backticks, direct prose, copy-pasteable commands, and ASCII punctuation."
---

# Markdown Formatting

## Fences

When a fenced code block contains backticks, do not wrap it in a plain
triple-backtick fence because inner backticks can close the block early:

- Prefer tilde fences (`~~~`) as the outer delimiter.
- Alternatively use 4 or more backticks when tildes are unsuitable.

## Style

- Keep prose direct, concrete, and easy to scan.
- Prefer copy-pasteable commands with realistic flags.
- Avoid emojis unless explicitly requested.
- Use ASCII punctuation (`--`, `->`, `"..."`) unless surrounding content
  requires Unicode.
