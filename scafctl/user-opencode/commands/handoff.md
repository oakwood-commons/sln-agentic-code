---
description: "Compact the current session into a secret-redacted Markdown handoff file under /tmp/opencode/ so a fresh session can resume cleanly."
agent: build
---

<!-- Adapted from mattpocock/skills (MIT License): https://github.com/mattpocock/skills -->

Write a concise handoff document under `/tmp/opencode/handoff-<topic>-<date>.md`
(never inside the repository working tree) summarizing:

$ARGUMENTS

Include:
1. **Context:** what this session is building and why.
2. **State:** done, in progress, and remaining (link files/PRs by path or URL
   rather than duplicating content).
3. **Key Decisions:** non-obvious choices settled with the user.
4. **Open Questions:** unresolved items for the next session.
5. **Next Steps & Skills to Load:** exact commands and skills to load first.

Scrub any credentials, tokens, or personal data (`<REDACTED>`). Do not commit
or push.
