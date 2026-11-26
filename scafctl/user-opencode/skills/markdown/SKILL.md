---
name: markdown
description: "Markdown authoring rules: tilde fences around nested backticks, direct prose without em-dashes or dash clauses, and copy-pasteable commands."
---

# Markdown Formatting

## Fences

When a fenced code block contains backticks, do not wrap it in a plain
triple-backtick fence because inner backticks can close the block early:

- Prefer tilde fences (`~~~`) as the outer delimiter.
- Alternatively use 4 or more backticks when tildes are unsuitable.

## Punctuation & Style (Strict No-Dash Rule)

- **Never use em-dashes, en-dashes, or spaced hyphens/double-hyphens (`" - "`, `" -- "`) as sentence punctuation.** Do not substitute `-` or `--` where an em-dash would go.
- Rewrite sentences so no dash punctuation is needed at all: use colons (`:`), commas (`,`), semicolons (`;`), parentheses `(...)`, or separate sentences.
- Write numeric or step ranges with words (`1 to 3 turns`, `Steps 1 through 5`) rather than hyphen or dash ranges.
- Hyphens remain normal inside compound modifiers (`read-only`, `open-source`), CLI flags (`--dry-run`), Markdown list bullets (`- `), YAML `---` delimiters, and table separators (`|---|---|`).
- Keep prose direct, concrete, and easy to scan.
- Prefer copy-pasteable commands with realistic flags, and avoid emojis unless explicitly requested.
