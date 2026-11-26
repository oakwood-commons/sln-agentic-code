---
name: grilling
description: "Batch-frontier structured interview for resolving multiple independent design questions in parallel rounds while honoring stated constraints and favoring best-practice designs over shortcuts. Prefer stepwise-interview when plain one-at-a-time Q&A is easier to digest."
---

<!-- Adapted from mattpocock/skills (MIT License): https://github.com/mattpocock/skills -->

# Grilling -- structured interview via frontier of rounds

Map open design questions as a **design tree**: every decision branches into
the decisions that depend on it.

> Note: For complex problems where a batch of questions feels overwhelming,
> prefer the `stepwise-interview` skill (one plain-English question at a time).
> Use `grilling` when you and the user want to clear all independent frontier
> questions in a single numbered batch.

## Honor stated context & challenge suboptimal "how"

1. **Never grill backwards**: Treat libraries, constraints, and goals the user
   already stated in their prompt as settled. Only put genuinely unresolved
   product or architecture decisions on the frontier.
2. **Prior art & best over easiest**: Ground every recommended answer (`->`) in
   how mature tools and idiomatic ecosystems solve the problem, preferring clean
   seams (`codebase-design`) over quick hacks.
3. **Push back on prescribed "how" when warranted**: If the user suggested an
   implementation mechanism that is more fragile or complex than a cleaner
   design, call out the better alternative explicitly instead of silently
   building around the weaker mechanism.

## Work the tree in rounds

The **frontier** is every open decision whose prerequisites are already settled.
Ask the current frontier with a recommended answer (`->`) on each item, then
wait for the user's answers before computing the next round:

~~~text
Q1 -- <question title>: <question body>

-> <your recommended answer grounded in best practice>
~~~

## Facts are your job, never the user's

When a question needs a fact from the codebase, comments, `README.md`,
`AGENTS.md`, docs, or an MCP tool rather than a human judgment call, look it up
yourself. Only put real design or scope choices to the user.

## Done condition

Stop interviewing once every branch of the design tree is settled, and wait for
explicit confirmation on the resulting plan before implementing.
