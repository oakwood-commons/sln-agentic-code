---
name: security-audit
description: "Use when reasoning about a security vulnerability -- auditing code for one, chasing a red-team or pentest report, judging whether a suspected leak is real, or fixing a credential/auth/multi-tenant flaw found during review. Covers proving a finding by execution, discarding non-findings, and the bug classes static review misses. The security-auditor agent loads this for whole-codebase audits; load it directly when fixing inline. Reconciling a whole external report against a repo is the security-sweep skill's step instead."
---

# Security audit

Static reading generates hypotheses. It does not generate findings. A finding
is a hypothesis you **reproduced**, and credibility dies the first time someone
disproves a reported vulnerability by running it.

This skill is the shared methodology. Two contexts use it:

- **Audit** -- a read-only sweep of a codebase producing a findings doc. The
  `security-auditor` agent owns that workflow; it loads this skill for the
  bug classes and the proof discipline below. The `security-sweep` skill is
  the end-to-end entry point that pairs the agent with a scanner battery.
- **Inline** -- you are reviewing or writing code and intend to fix what you
  find. Same proof discipline, different endgame: see "Fixing inline".

## Reproduce or retract

Build the code and exercise the path with a canary value before calling
anything a finding. A hypothesis you could not reproduce is reported as
*suspected, unconfirmed* -- never as a finding.

Scale the ceremony to the context: an audit reproduces every hypothesis
end-to-end; inline, a five-line direct call against the suspect function is
usually proof enough. What does not scale down is the rule that something
unproven is never stated as fact.

## Establish the trust model before rating

Frame by **trust boundary**, not severity. Ask which boundary a hypothesis
crosses if true -- typically: execute code on the host, reach *server*
credentials, reach *other users'* data or credentials. A hypothesis crossing no
boundary is not a finding.

Which boundaries exist is a property of the system's **purpose**, so read for
that purpose before rating anything. A capability is not a vulnerability: a
service designed to execute submitted code is no more "vulnerable to RCE" than
a CI runner is. Reporting a product's stated function as its worst finding
buries the real one.

Answer three questions from the repo itself -- design docs, command help,
config defaults, shipped deployment artifacts -- before any severity leaves
your hands:

- **What is the default exposure?** Bind address, auth default, and whether any
  Dockerfile/chart/manifest ships something more exposed than that default. A
  loopback-bound default is a real mitigating control; so is shipping no
  deployment artifact at all.
- **Is the dangerous capability the product's stated purpose?** If yes, the
  finding is the missing *boundary around* it (who may invoke it), not the
  capability.
- **What do the docs promise operators?** Every security claim in the
  documentation is a control you must verify exists.

Rate against a public scale -- OWASP API Top 10, CWE, CVSS vectors -- and name
the posture you rated. A finding that is severe once exposed and moot on
loopback must say both halves.

**Severity is never inherited.** Prior art -- another repo's issue tracker, a
past audit, a sibling codebase -- generates hypotheses, never ratings. That
system's defaults, exposure, and purpose differ from this one's. Re-derive
every rating here; importing one is how an industry-normal practice gets
reported as a vulnerability.

## Reading a canary

A canary in the output only proves a leak if the code under test put it there.

Before believing one, ask: **did my own harness plant this?** A test upstream
that echoes the caller's token, a fixture writing the secret into the response,
a template whose data context already holds the value -- each manufactures a
convincing false positive. Re-run with the harness's copy removed; if the
canary still appears, it is real.

Distinguish leaking a **credential** from leaking the **data that credential
unlocked**. They carry different severities and are trivially conflated.

## Bug classes to check explicitly

These survive diff review because no single diff contains them.

- **Guard with an escape hatch.** When a boolean gates a credential path,
  confirm every fallback sits *inside* the guard. A fallback after the `else`
  silently defeats it. Prove it with a direct call at both boolean values.
- **Guard whose unconfigured default protects nothing.** A merged fix that
  defaults to audit/log-only, or that each call site must opt into. Read the
  shipped default *and* the deployed configuration before accepting that a
  finding is closed -- "fixed in code" is not "fixed".
- **Shared cache keyed without identity.** In an authenticated service, a cache
  keyed by URL or resource ID alone serves user A's response to user B. Check
  the key derivation, whether authenticated responses are cached at all, and
  whether upstream `Cache-Control` / `Vary` is honored or overridden. Prove
  with two requests, differing credentials, one hit counter.
- **Process-global state holding per-caller data.** Package-level vars,
  singletons, long-lived maps that touch tokens. Per-request construction is
  the safe pattern -- confirm it rather than assume it.
- **Unsafe-by-default constructor.** A helper whose default return is the
  dangerous variant, so every caller must opt out. Inventory *all* call sites;
  the bug is the ones that forgot. Check `git log -S` for a prior fix that
  regressed -- an unguarded control drifts back.
- **Confused deputy.** Any endpoint where the server acts with its own identity
  on caller-supplied input. Ask what the caller can make the server *do*, not
  only what they can read.
- **Forwarded identity headers.** The confused-deputy class covers the server's
  own credential; this one is headers the caller makes the server *send*.
  Check every server-originated outbound request for inbound headers forwarded
  verbatim -- `Authorization`, `Cookie`, `X-Forwarded-*`,
  `X-Authenticated-*` -- where the upstream ends up attributing the request
  to an identity the caller chose.
- **Authentication without authorization.** Verify a real per-route scope
  check exists, not just a valid-token check. Read the middleware for
  commented-out scope logic and for one flat scope covering every route.
- **Documented control that does not exist.** A control promised in the docs
  but absent from the code is worse than an undocumented gap: operators read
  the promise and provision around it, so the gap ships with false assurance
  attached. Grep the documentation for every security claim -- auth
  requirements, role checks, loopback restrictions, redaction, allowlists --
  and verify each against the implementation. A doc table describing an
  access-control matrix is a checklist, not evidence.
- **Handler defined but never wired.** A function that looks like a live
  endpoint (`Register*Route`, a Huma operation, a Cobra command) is not a
  finding until you confirm it is actually registered in the startup
  sequence. Diff every handler/route function against what the server
  actually calls at boot *before* auditing it -- dead code treated as live
  attack surface is one of the most common audit false positives.
- **Interpreter capability inventory.** For every expression language,
  template engine, or action runner reachable from a request, enumerate
  *every* registered function or action and check each for server-mode
  gating -- not just the ones a prior finding named. Debug helpers and
  file/deploy actions regularly survive a narrow sandbox fix this way.
- **Availability budget.** What one request can make the server spend:
  uncapped sleep or loop primitives, context-unaware waits, unbounded body
  reads, missing server and per-request deadlines, a rate limiter
  constructed but never attached to any handler, or one keyed on a
  spoofable client IP.
- **Partially fixed finding.** The reported trigger is closed while the root
  cause lives on -- a length limit raised so a credential no longer fails
  validation, while validation errors still echo the submitted value.
  Re-derive the root cause before accepting a fix as closing a finding;
  re-running the original trigger proves only the trigger.
- **Arbitrary read: inventory everything mounted.** After proving a
  file-read primitive, enumerate everything it can reach -- mounted
  secrets, env vars, projected tokens -- and any token-exchange chain to
  cloud credentials (for example, a projected token exchanged via
  workload identity federation). Remediation then includes rotation and an
  audit-log review, not just the code fix.

## Non-findings

Reporting these burns the reader's trust in everything else:

- A caller supplying their own credential and receiving their own access. That
  is the passthrough model working.
- Secrets exposed only to whoever already supplied them.
- A "leak" of a file *path* rather than the file's contents.
- Theoretical reachability with no reproduction and no stated reason it could
  not be reproduced.
- **The industry norm.** Check the deviation against what comparable systems
  actually do before calling it a defect. Unauthenticated metrics endpoints,
  for instance, are standard practice -- controlled by network placement and a
  separate listener, not by auth -- so "metrics is unauthenticated" is a
  hardening observation at most, and the real question is whether the operator
  can bind it somewhere isolated.

Say explicitly that these were excluded, so the reader knows they were
considered rather than missed.

## Before testing a URL/host-shaped parameter

Trace where it actually goes before pointing a live payload at it. A
parameter that looks like a pass-through (`baseUrl`, `targetUrl`) may be
assigned to a local variable and never threaded into the client that does
the real fetch -- fire a live SSRF probe before confirming this and you hit
whatever the hardcoded default points at, which can be a real production
endpoint. One code-read is cheaper than one accidental request to prod.

## Sweep for siblings of every already-fixed pattern

A shipped fix protects the call site someone was looking at, not the
pattern. Once you know a pattern was fixed anywhere (a prior commit, a PR,
even a fix you just made this session), treat that as a lead, not a closed
case: **enumerate every other call site of the same underlying primitive**
and check each one for the same protection. Real vulnerabilities routinely
hide in sibling paths -- for example, a credential-fallback gate applied to
one function but not its near-identical sibling in another file, or a
URL-rejection check applied to one fetch path but not a second field fetched
through the same underlying function.

Concretely: grep every caller of the guarded primitive (not just the one
name-checked in the fix's commit message), and for each one ask "does this
call site have the same guard, or does it call the dangerous function
directly?" A clean sweep -- every sibling protected -- is a valid, useful
result; report it as such rather than manufacturing a finding to justify
the pass.

## Live-verifying against shared infrastructure

Reproducing locally proves the code has the flaw. It does not prove a
specific shared/staging/production system is still running that code,
still configured that way -- and closing that last gap is not yours to do
unilaterally. Never send a live exploit payload to shared infrastructure
using a credential you obtained yourself, and never try to obtain one to
force the point. That request needs the same sign-off a pentest would.

Instead, hand the human a self-contained script they run with their own
credential:

- Requires an explicit typed confirmation before it sends anything.
- Never prints a raw secret value -- compare in-process (`auth == expected`)
  and log only booleans/lengths, so the transcript itself cannot leak
  whatever the test retrieves.
- Reports three states, not two: **CONFIRMED**, **REFUTED**, and
  **INCONCLUSIVE**. A network path that never delivers a callback (a
  firewall, a workstation unreachable from the server's side) is not the
  same as "the bug is not there" -- collapsing that distinction turns an
  unrelated network gap into a false all-clear.
- If a preliminary check only proves the attack surface exists (endpoint
  live, accepts the input shape) rather than the dangerous outcome itself,
  label that result narrowly. A reachability check succeeding is not the
  same claim as the exploit succeeding, and calling both "CONFIRMED" invites
  legitimate pushback because the label overreached.

## Verify severity claims before they leave your hands

A CRITICAL/HIGH claim, whether your own or a subagent's, gets one more gate
before it goes in the doc: reproduce it yourself, not just re-read the trace
that produced it. Two independent traces agreeing is not the same as a
third-party execution. This is the same reproduce-or-retract rule above,
applied specifically to the highest-severity bucket where being wrong costs
the most credibility.

## Fixing inline

When the intent is to fix rather than report:

- **Keep the reproduction as a regression test.** This is the main divergence
  from audit mode, where exploit tests are scratch and deleted. A credential or
  multi-tenant bug with no test guarding it drifts back -- the
  unsafe-by-default class above exists because exactly that happened.
- **For a mode-gating fix, prefer a golden-list inventory test.** A CI test
  that lists every registered function, action, route, and read path, each
  classified server-safe or server-denied, fails when a new registration
  lands unclassified. This is a stronger control than any review-time
  prompt, and the one that keeps debug helpers and deploy actions from
  drifting back into reach.
- Fix the **default**, not the call site. If a constructor returns the
  dangerous variant, changing the two callers you found leaves the next caller
  exposed. Invert the default and make the dangerous path opt-in.
- Check whether the fix breaks existing callers, and say so. Closing a confused
  deputy turns a working tokenless request into a 400 for whoever relied on it.
- Sweep for siblings before declaring it fixed -- see "Sweep for siblings"
  above; the same discipline applies to a fix you just wrote, not only to
  a codebase-wide audit.

## Consolidating findings across multiple passes

A multi-pass or multi-agent audit tends to produce one narrative doc per
pass plus a running claim of "added to the plan." That claim is not true
until verified: grep the actual plan file for the fix number in the same
turn you claim to have added it. A single authoritative fix-tracking
document, merged into immediately rather than deferred, beats several
supplement docs each asserting a merge that never happened.
