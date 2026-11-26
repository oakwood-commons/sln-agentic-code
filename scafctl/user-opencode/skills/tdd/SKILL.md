---
name: tdd
description: "Use automatically when implementing a feature or fixing a bug -- writes one failing test at a public seam first, makes it pass, and refactors (Red -> Green -> Refactor)."
---

<!-- Adapted from mattpocock/skills (MIT License): https://github.com/mattpocock/skills -->

# Test-Driven Development

Tests verify observable behavior through public interfaces, not private
implementation details. Code internals can change completely during a refactor;
good tests still pass.

## Test at Pre-Agreed Seams

A **seam** is the public interface where callers and tests observe behavior
without reaching inside. Confirm the public interface and test seams up front
before writing tests.

## Three Anti-Patterns to Avoid

1. **Implementation-coupled** -- mocking private helpers or asserting internal
   fields so tests break on harmless refactors.
2. **Tautological** -- recomputing the expected value using the exact formula
   under test instead of asserting against an independent literal or spec.
3. **All-tests-up-front batching** -- writing 20 speculative tests up front and
   then all the code. Work **one behavior at a time** instead: one failing test
   -> minimal implementation -> green suite -> refactor -> repeat.

## Rules of the Loop

1. **Red before green**: Write one failing test first and watch it fail.
2. **Red must fail for the expected reason**: Confirm it fails because the
   behavior is missing, not because of a syntax or fixture error.
3. **Green means the whole suite**: Run the full package/module test and lint
   suite before moving on.
4. **Refactor after green**: Clean up naming and duplication in the touched
   code while tests stay green before starting the next failing test.
