---
name: security-sweep
description: "Use when asked to security-sweep, scan, or audit a whole repo end to end, or to reconcile an external red-team/pentest report against a repo. Runs a scanner battery (govulncheck, gosec, trivy) alongside the security-auditor agent, triages every scanner hit through reproduce-or-retract, produces one findings doc plus one issue draft per confirmed finding, and buckets an external report's findings as fixed/known/declined/untracked."
---

# Security sweep

The entry point for a whole-repo audit. Two engines run side by side:

- **Scanners** find the mechanical classes fast: known-vulnerable
  dependencies, risky API use, committed secrets, and misconfigured manifests.
- **The `security-auditor` agent** finds what no scanner models: trust
  boundaries, authz gaps, and multi-tenant leaks.

Neither engine's output is a finding. A scanner hit is a **hypothesis**,
exactly like a line of suspicious code: it becomes a finding only once
reproduced, per the `security-audit` skill. Load that skill now; this one
covers only the orchestration around it.

## 1. Scope

Confirm with the user, or state as assumptions:

- the repo or path, and which modules are in scope (`go.mod` locations,
  skipping `vendor/` and `testdata/`);
- the trust boundaries (default three: code execution on the host, server
  credentials, and other users' data or credentials);
- a scratch directory OUTSIDE the repo for scanner output, the findings
  doc, and issue drafts (unfixed-vulnerability writeups are never committed).

Done when: scope, boundaries, and the scratch path are stated in chat.

## 2. Check the tools

~~~bash
command -v govulncheck gosec trivy gitleaks
~~~

For each missing tool, give the user its install command and ask; never
install without asking. Proceed with what is present and record every
skipped scanner; the doc must say what was NOT scanned. Non-Go repos: `trivy
fs` still covers dependencies, secrets, and manifests; add the ecosystem's
own audit (`npm audit`, `pip-audit`) if installed.

Done when: each scanner is marked present, installed-with-consent, or
skipped.

## 3. Run the battery and the auditor in parallel

Dispatch the `security-auditor` agent with the scope and boundaries from
step 1, telling it the scanners run separately so it spends its pass on
logic, not on re-running them. While it works, run the battery (outputs to
the scratch dir):

~~~bash
# <module-id> = the module's path with "/" replaced by "-" (nested modules
# like services/api would otherwise break the report path). Run each from
# the module's directory.
(cd <module> && govulncheck -format json ./... > <scratch>/govulncheck-<module-id>.json)
(cd <module> && gosec -fmt=json -out <scratch>/gosec-<module-id>.json ./...)
# repo root
trivy fs --scanners vuln,misconfig,secret,license --format json -o <scratch>/trivy.json .
# optional, ask first: walks all of git history (gitleaks >= 8.19)
gitleaks git --report-path <scratch>/gitleaks.json .
~~~

Classify each exit code against the tool's documented set. An exit code inside
the set means the scan completed (findings or not): gosec and gitleaks 0/1,
govulncheck 0 in JSON mode (3 in text mode means findings), trivy 0. An
exit code outside the set is a failed scan even if a report file exists
(a scanner killed mid-run can leave truncated output); treat a missing,
empty, or unparseable report the same way, and say so in the coverage
section. If gitleaks is skipped, git history was not scanned for secrets
(`trivy fs` sees only the working tree) and the coverage section must say so.

A harness without subagents runs the battery first, then follows the
`security-auditor` agent's procedure inline.

Done when: every scanner has exited (record each exit code and tool
version) and the auditor has returned.

## 4. Triage every scanner hit

Dedupe across scanners, then put each hit in a bucket:

- **Dependency vuln.** govulncheck separates vulnerabilities your code
  *calls* from ones merely *imported* (in its JSON, a `finding` whose
  `trace` names a `function` is called; module- or package-level only is
  imported). Only called ones are finding candidates: confirm the affected
  version is the one in use and name the fixed version. Imported-only is
  dependency hygiene, listed separately.
- **Code-pattern hit (gosec).** Read the code. Most hits cross no trust
  boundary (unchecked errors, a file path from a CLI flag the operator
  already controls). Keep only those that would cross one, and hand that
  shortlist to the auditor (resume its session) for reproduce-or-retract;
  do not rate them from the scanner's severity.
- **Secret hit.** Verify it is a real credential, not a test fixture or a
  placeholder. A real one is a finding even if a later commit deleted it,
  because git history keeps it; the fix is rotation, not deletion. Never paste
  the value; record its name, file, and commit.
- **Misconfig hit** (Dockerfile, Kubernetes, Terraform). Check whether the
  deployed configuration actually carries it before rating.
- **License hit.** Not a security finding; a separate list for the owner.

Done when: every hit sits in exactly one bucket: confirmed (with
reproduction), suspected (with why it could not be run), disproved (with the
evidence), hygiene, or license.

## 5. Reconcile an external assessment

Run this step when a red-team or pentest report exists for the repo, since the
sweep is often driven by one and the reconciliation is where most of its value
lands. An external report's findings are hypotheses like any scanner hit: same
trust boundaries, same reproduce-or-retract.

1. Convert the report to markdown (e.g. pandoc) for grep-ability, into the
   scratch dir from step 1 with restrictive permissions. Assessment reports are
   usually classified: never commit one, and remove the converted copy when the
   reconciliation finishes.
2. Pin two deployment states separately: the assessed code revision and
   the versions deployed during testing (the assessment baseline, since a fix
   merged after that revision cannot be credited by the report), and the
   versions currently deployed in each environment (remediation status now).
3. Bucket every finding: **fixed** (merged, verified in code, and running in
   the currently deployed versions from step 2; a fix the deployments do not
   contain yet is tracked, not fixed), **known / tracked**, **declined by
   design** (requires a written risk acceptance, not just a closed issue
   state), **untracked / news**. Untracked findings go through step 4's triage
   like any other hypothesis before becoming findings.
4. Cross-reference each bucket against PRs, issues, and the tracker; flag
   state disagreements where a tracker claims fixed while the code says
   otherwise.
5. Re-verify every "fixed" claim against current code, because claims read from
   issue titles alone are often wrong.
6. Feed the misses back: every class the external team found that this
   sweep's own pass missed becomes a candidate addition to the `security-audit`
   skill's bug-class list. File that or say so in the report.

Done when: every finding sits in exactly one bucket, every "fixed" claim
has been re-verified in code, and the misses are recorded as feedback.

## 6. One findings doc

Merge into the auditor's doc (its format: per-boundary answers, confirmed /
suspected / disproved). Add a **Scanner coverage** section: tool, version,
exit code, hit counts per bucket, and every skipped scanner. Add the
reconciliation buckets when step 5 ran. Summarize dismissed noise by rule ID
and count, not hit by hit, so the reader knows it was considered without
re-reading noise.

Done when: grepping the doc finds every confirmed finding and the coverage
section (verify rather than assuming the merge landed).

## 7. Issue drafts

One draft per **confirmed** finding, in the scratch dir, in the
`create-issue` skill's draft shape (reproduction and actual output, no
secret values). Suspected findings stay in the doc unless the user asks for
drafts. Filing follows the `create-issue` skill's flow: show the drafts, and
file only on the human's go, asking whether any finding should be a private
security advisory instead of a public issue.

## 8. Report

Per-boundary verdicts, confirmed findings in priority order, the doc and
draft paths, and the coverage gaps (skipped scanners, unreachable
deployments). With a reconciliation: its buckets and the classes this sweep
missed. Say plainly when a boundary came back clean.
