---
name: stepwise-interview
description: "Use when resolving ambiguous requirements or complex architecture choices with a human: asks ONE plain-English question at a time with a clear recommendation instead of dumping a wall of technical questions. Pairs directly with wait-what."
---

# Stepwise Interview -- one plain question at a time

When solving a complex or ambiguous problem, dumping five or six dense
multi-part questions at once overloads the human and hides the real trade-offs.
Walk the design tree **one decision at a time in everyday language**.

## Rules

1. **Look up facts yourself first**: Never ask the user a question that reading
   the codebase, git history, CLI `--help`, or MCP tools can answer. Gather
   those facts silently before speaking.
2. **One question per turn**: Pick the single highest-leverage open decision
   whose prerequisites are already settled, ask only that question, and stop to
   wait for the user's answer.
3. **Plain, everyday English**:
   - Give 1-2 short sentences of context explaining *why* this choice matters.
   - Avoid internal jargon, framework buzzwords, or compressed shorthand unless
     explained inline in plain words.
   - If the user replies `/wait-what`, `"what?"`, or `"I don't follow"`, load
     the `wait-what` skill immediately and re-pitch the question more simply.
4. **Always include your recommendation**: Never hand the user a blank page.
   State your recommended choice and one clear sentence on why:

~~~text
### Question 1: <short plain title>

<1-2 short sentences of everyday context.>

- **My recommendation:** <concrete choice and why it is the safest/simplest default.>

Does that work for you, or would you prefer <clear alternative>?
~~~

5. **Adapt the tree after each answer**: Every answer settles a branch and may
   eliminate later questions altogether. Recompute what actually remains open
   before asking the next question.
6. **Stop before acting**: When no open assumptions remain, present a concise
   summary or blueprint and get the user's explicit go-ahead before writing
   production code.
