---
name: go-debug
description: "Use when diagnosing Go panics, data races, deadlocks, hangs, resource leaks, or module/build errors."
---

# Go Debugging

Adapted from Go debugging notes crediting `samber/cc-skills-golang` (MIT):
https://github.com/samber/cc-skills-golang (see `../../THIRD-PARTY-NOTICE.txt`).

- Reproduce the symptom with a fast, deterministic command before editing code.
  Rank falsifiable hypotheses, change one variable at a time, and lock the fix
  in with a behavioral regression test.
- After 2 failed fix attempts on a fast model, escalate a minimal redacted
  evidence bundle to a stronger reasoning model; after 3 failed attempts, stop
  and ask the human.
- Check typed-`nil` interfaces, shadowed `err` variables, `defer` inside loops,
  unclosed files/HTTP bodies, slice backing-array aliasing, `nil`-map writes,
  and goroutine leaks.
- Use `go test -race` for shared-state code and `-count=N` for suspected flakes;
  inspect **both** stacks in a race report.
- Remember `SIGQUIT` dumps stacks AND terminates a Go process; prefer an
  internal `net/http/pprof` endpoint or `dlv` on live processes, and never
  expose `pprof` publicly.
- Use `go tool pprof` for CPU/heap/block/mutex profiles, redact secrets from
  captured output, and remove temporary debug probes before finishing.
