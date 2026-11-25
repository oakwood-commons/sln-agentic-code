---
title: "2. Ask & Plan First"
weight: 20
duration: "~4 MIN HANDS-ON"
goal: "Settle three hidden product choices in plain English before the AI writes a single line of Go."
---

# Step 2: Ask & Plan First

## What Usually Goes Wrong

Suppose you open an AI tool and type: *"Add link expiration and click stats to
`shorty`."*

Sounds clear enough -- except you left three real product decisions unspoken:
- When someone visits an expired link, should they get `404 Not Found` or `410
  Gone`?
- If the caller omits `expires_in`, should the link live forever or expire after
  a default like 24 hours?
- After a link expires, can the creator still check `GET /links/{code}/stats` to
  see how many clicks it got?

Without guardrails, an AI silently guesses all three, writes 250 lines of Go
around those guesses, and forces you to untangle its code when you disagree.

People often call this *"AI hallucination"*, when in reality the AI is simply
giving the best answer it can based on its current
[context](/docs/reference/glossary.md#context). The trap is assuming the AI has
the same unspoken context a human teammate would have -- when often it does not.

## What We Do Instead

We mitigate missing context in two ways: we supply project rules in `AGENTS.md`,
and **we plan first so the AI explains what it is going to do** before writing a
single line of code. Three [slash commands](/docs/reference/glossary.md#slash-command)
in OpenCode make that effortless:

1. **`/interview`** -- loads a small playbook (a [skill](/docs/reference/glossary.md#skill))
   telling the AI to inspect the codebase silently first, then ask you **one
   plain-English question at a time** with its recommended answer.
2. **`/wait-what`** -- if the AI ever slips into dense technical jargon, typing
   `/wait-what` makes it pause and re-explain its last point in everyday words.
3. **`/plan`** -- dispatches a read-only planning helper (a
   [subagent](/docs/reference/glossary.md#subagent) that is blocked from editing
   files) to turn your answers into a short, step-by-step plan.

---

## Walk Through It

{{% steps %}}

1. **Launch OpenCode inside the repository**

    From the root of `sln-agentic-code`, start OpenCode:

    ~~~bash
    opencode
    ~~~

    OpenCode automatically reads [`AGENTS.md`](https://github.com/oakwood-commons/sln-agentic-code/blob/main/AGENTS.md)
    so it already knows how this repository builds and tests Go code.

2. **Start a one-question-at-a-time interview (`/interview`)**

    Paste this prompt into OpenCode:

    ~~~text
    /interview In examples/go-app, we want to add optional link expiration
    (expires_in seconds on POST /links) and a click-count endpoint
    (GET /links/{code}/stats).
    ~~~

    **What you should see:** Instead of editing files or dumping six questions
    at once, OpenCode reads `shortener.go` and asks **Question 1** in plain
    English, along with a clear **My recommendation** (for example, returning
    `410 Gone` for expired links while `404 Not Found` stays reserved for
    unknown codes).

    Answer each question in a word or two (`yes`, or state your preference). We
    recommend settling on:
    - **Expired redirect status:** `410 Gone` (with a new `ErrExpired` error).
    - **Missing `expires_in`:** Optional -- `0` or omitted means the link does
      not expire; negative numbers return `400 Bad Request`.
    - **Stats on expired links:** Still return `200 OK` from
      `GET /links/{code}/stats` (`{"code":"...","url":"...","clicks":N,"expired":true}`)
      so creators can see final click totals after a campaign ends.

    > [!TIP]
    > **Try `/wait-what` anytime:** If any question or explanation feels
    > unnecessarily wordy or full of buzzwords, reply `/wait-what` and watch the
    > agent re-pitch it in plain sentences.

3. **Generate a read-only blueprint (`/plan`)**

    Once the interview wraps up with zero open questions, run:

    ~~~text
    /plan Implement optional expires_in and GET /links/{code}/stats in
    examples/go-app based on the decisions we just agreed on.
    ~~~

    **What you should see:** The read-only `planner` subagent checks the
    existing `Store` interface in `shortener.go` and proposes a short blueprint
    broken into small slices -- **without modifying a single file yet**.

{{% /steps %}}

Notice what those 3 minutes bought you: you caught three product edge cases up
front, and you can read a 15-line plan instead of reviewing 250 lines of
surprise code.

In **Step 3**, we will build that plan in two small, tested slices.
