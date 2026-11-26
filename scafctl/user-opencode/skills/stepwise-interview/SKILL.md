---
name: stepwise-interview
description: "Use when resolving ambiguous requirements or architecture choices with a human: honors already-stated constraints, checks prior art for the best (not easiest) design, challenges weak 'how' suggestions, and asks ONE plain-English question at a time. Pairs directly with wait-what."
---

# Stepwise Interview: One Plain Question at a Time

When solving a non-trivial or ambiguous problem, dumping five or six dense
multi-part questions at once overloads the human, while grilling the human
backwards on basics they already stated wastes their time. Walk only the
**genuinely open** design tree **one decision at a time in everyday language**.

## Rules

1. **Read repo context silently first**: Never ask the user a question that
   reading the codebase, comments, `README.md`, `AGENTS.md`, git history, CLI
   `--help`, or MCP tools can answer. Gather those facts before speaking.
2. **Honor stated request context (never grill backwards)**: Treat libraries,
   stack choices, boundaries, and goals the user already supplied in their
   prompt as settled inputs. Do not ask them to re-justify settled preferences
   unless a hard technical conflict exists in the codebase.
3. **Bring prior art, favor *best* over *easiest*, and challenge weak "how"**:
   - Check how mature tools and idiomatic libraries solve this problem. Anchor
     recommendations in the **cleanest long-term design** (`codebase-design`),
     calling out when the *best* approach differs from the *easiest* shortcut.
   - Do not rubber-stamp a prescribed *"how"*: if the user's prompt suggests a
     specific implementation mechanism when a simpler, safer, or more testable
     design exists, briefly share that better option and why before locking in
     the plan.
4. **One question per turn**: Pick the single highest-leverage open decision
   whose prerequisites are already settled, ask only that question, and stop to
   wait for the user's answer.
5. **Plain, everyday English**:
   - Give 1 to 2 short sentences of context explaining *why* this choice matters.
   - Avoid internal jargon, framework buzzwords, or compressed shorthand unless
     explained inline in plain words.
   - If the user replies `/wait-what`, `"what?"`, or `"I don't follow"`, load
     the `wait-what` skill immediately and re-pitch the question more simply.
6. **Always include your recommendation**: Never hand the user a blank page. On
   Turn 1, if helpful, share 1 to 2 sentences on best-practice prior art for the
   overall approach, then present Question 1 with your recommended choice:

~~~text
### Question 1: <short plain title>

<1 to 2 short sentences of everyday context.>

- **My recommendation:** <concrete choice grounded in best practice, and why it beats the shortcut.>

Does that work for you, or would you prefer <clear alternative>?
~~~

7. **Adapt the tree after each answer**: Every answer settles a branch and may
   eliminate later questions altogether. Recompute what actually remains open
   before asking the next question.
8. **Stop before acting**: When no open assumptions remain, present a concise
   summary or blueprint and get the user's explicit go-ahead before writing
   production code.
