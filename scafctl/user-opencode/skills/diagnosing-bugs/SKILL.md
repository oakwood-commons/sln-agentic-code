---
name: diagnosing-bugs
description: "Structured discipline for hard or flaky bugs: build a tight red-capable feedback loop first, minimize, rank falsifiable hypotheses, use tagged probes, and add a regression test."
---

<!-- Adapted from mattpocock/skills (MIT License): https://github.com/mattpocock/skills -->

# Diagnosing Bugs

Never guess at code changes before you can reliably reproduce the failure.
Redact any secret or token (`<REDACTED>`) in captured logs.

## Phase 1: Build a Tight Feedback Loop

Construct one fast, deterministic command (a targeted test, CLI invocation with
fixture input, or minimal reproduction script) that goes **red** on the exact
reported symptom and **green** once fixed. If you cannot build a reproduction
loop, stop and ask the user for the missing fixture or environment detail.

## Phase 2: Reproduce & Minimize

Confirm the failure matches the reported symptom and shrink the input and setup
until removing anything else makes the failure disappear.

## Phase 3: Rank 3 to 5 Falsifiable Hypotheses

Avoid anchoring on your first hunch. Write 3 to 5 ranked hypotheses with an
explicit prediction: *"If X is the cause, then Y will happen."*

## Phase 4: Instrument Cleanly

Test one variable at a time. Prefix temporary debug logs with a unique tag such
as `[DEBUG-a4f2]` so cleanup is a single `grep`.

## Phase 5: Regression Test, Fix & Circuit Breaker

1. Turn the minimized reproduction into a failing test at the public seam.
2. Apply the minimal root-cause fix and watch the test and full suite go green.
3. Remove all `[DEBUG-...]` probes.
4. **Circuit breakers**:
   - After **2 failed fix attempts** on a fast model, escalate the minimal
     evidence bundle to a stronger reasoning model.
   - After **3 failed fix attempts**, **STOP** patching and ask the human, as
     the surrounding abstraction or boundary needs redesign.
