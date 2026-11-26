---
name: go-perf
description: "Use when profiling Go performance, writing benchmarks, or comparing latency and memory allocations."
---

# Go Performance

Adapted from Go performance notes crediting `samber/cc-skills-golang` (MIT):
https://github.com/samber/cc-skills-golang (see `../../THIRD-PARTY-NOTICE.txt`).

- Define the target metric, profile the real workload, and record repeated
  baseline benchmarks before optimizing:
  `go test -run '^$' -bench BenchmarkName -benchmem -count=6 ./...`
- Change one variable at a time on the same hardware and Go version, and compare
  runs with `benchstat` when available; never claim wins from a single sample.
- Distinguish allocation rate (`B/op`, `allocs/op`) from retained heap. Remember
  CPU profiles miss off-CPU lock contention and I/O waits, so inspect block and
  mutex profiles when latency exceeds CPU time.
- Require measured workload evidence before tuning `GOGC` or `GOMEMLIMIT`.
