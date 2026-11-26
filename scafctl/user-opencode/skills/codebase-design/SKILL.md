---
name: codebase-design
description: "Shared vocabulary for designing deep modules. Use when designing or improving a module's interface, finding deepening opportunities, deciding where a seam goes, or making code more testable."
---

<!-- Adapted from mattpocock/skills (MIT License): https://github.com/mattpocock/skills -->

# Codebase Design

Design **deep modules**: substantial behaviour behind a small interface, placed
at a clean seam, testable through that interface. The aim is leverage for
callers, locality for maintainers, and testability for everyone.

## Glossary

- **Module:** anything with an interface and an implementation (a function,
  struct, package, or workflow).
- **Interface:** everything a caller must know to use the module correctly:
  type signatures, invariants, ordering constraints, error modes, and required
  configuration.
- **Implementation:** the body of code inside a module that hides complexity.
- **Depth:** leverage at the interface, meaning the amount of useful behaviour
  a caller or test exercises per unit of interface they must learn.
- **Seam:** a place where you can alter behaviour without editing in that
  place; the location where a module's interface lives.
- **Adapter:** a concrete implementation satisfying an interface at a seam
  (for example, a real directory reader vs. an in-memory map reader).
- **Leverage:** capability gained per unit of interface learned.
- **Locality:** changes, bugs, and verification stay concentrated in one place
  rather than spreading across callers.

## Principles

- **Depth is a property of the interface, not implementation line count.** Ask:
  can I reduce the number of methods, simplify parameters, and hide more
  details inside?
- **The deletion test.** Imagine deleting the module: if complexity disappears,
  it was a shallow pass-through; if complexity reappears across N callers, it
  was earning its keep.
- **The interface is the test surface.** Callers and tests cross the same seam.
- **One adapter means a hypothetical seam; two adapters mean a real one.** Do
  not invent seams unless something genuinely varies across them (such as
  production I/O vs. deterministic unit tests, or GitHub Actions vs. Tekton).

## Designing for testability

1. **Accept dependencies, don't construct them internally.**
2. **Return structured results, don't hide outputs in side effects.**
3. **Keep the surface area small.**
