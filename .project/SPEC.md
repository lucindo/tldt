# Specification: Full Code Audit (Behavior-Preserving)

## Problem
The `tldt` codebase has just migrated off the GSD planning workflow onto devskills and carries history from its `resumator` origin (legacy `src/`, mixed assets). Before further feature work, the maintainer needs a systematic quality pass across the active code — correctness, security, Go idioms, tests, and docs — applying only behavior-preserving improvements so the code is left in a measurably better state without any change to what the tool does for users.

## Scope
**In scope**
- Quality audit of active code: `cmd/tldt/`, `pkg/tldt/`, `internal/*`, plus `examples/`, `Makefile`, `docs/`, `README.md`.
- Running the relevant `/ds-*` review commands and consolidating their findings.
- Applying fixes that are strictly behavior-preserving: dead-code removal, idiom/style alignment, comment/doc accuracy, test quality, internal refactors (≤ same public behavior), AI-slop cleanup.
- Reporting findings that *would* change behavior (real bugs, security gaps) as recommendations — not auto-applied.

**Out of scope**
- Any change to observable functionality: CLI flags/output, `pkg/tldt` public API signatures and semantics, exit codes, file formats.
- New features, new dependencies, or version bumps not required by a fix.
- Reviving or modifying the legacy `src/` resumator service beyond flagging it (it is `//go:build ignore`).
- Performance optimization that alters output.

## Users
- **Maintainer (primary):** wants a trustworthy, idiomatic, well-tested codebase and a prioritized list of anything that needs a behavior change, decided explicitly rather than silently applied.
- **Library consumers of `pkg/tldt`:** must see zero API or behavior drift across the audit.

## Functional Requirements
- **FR-1:** The audit SHALL run each applicable `/ds-*` review command over the active code and capture its findings: `/ds-code-quality-review`, `/ds-bug-review`, `/ds-security-review`, `/ds-go-review`, `/ds-test-quality-review`, `/ds-doc-quality-review`, `/ds-deslop`.
- **FR-2:** The audit SHALL classify every finding as either **behavior-preserving** (safe to apply) or **behavior-changing** (recommend only).
- **FR-3:** The audit SHALL apply behavior-preserving fixes to the working tree.
- **FR-4:** The audit SHALL NOT apply any fix classified as behavior-changing; each such finding SHALL be recorded with location, severity, and proposed remediation for maintainer decision.
- **FR-5:** The audit SHALL preserve the public surface of `pkg/tldt` (exported identifiers, signatures, documented semantics) unchanged.
- **FR-6:** The audit SHALL preserve all CLI-observable behavior: flag set, defaults, stdout content, stderr stats/warnings routing, exit codes, and output formats (text/json/markdown).
- **FR-7:** The audit SHALL verify behavior preservation by running the full test suite before and after changes and confirming the same pass result.
- **FR-8:** The audit SHALL leave the repository building (`go build ./...`), vetting (`go vet ./...`), and passing `go test -race ./...`.
- **FR-9:** The audit SHALL produce a consolidated report summarizing applied fixes and outstanding recommendations, grouped by `/ds-*` source and severity.

## Non-Functional Requirements
- **NFR-1:** Behavior parity: 100% of pre-existing tests that passed before the audit SHALL pass after, with no test removed solely to make the suite green.
- **NFR-2:** Diff hygiene: every changed line SHALL trace to a specific audit finding; no opportunistic edits to untouched code (per AGENTS.md §3).
- **NFR-3:** Reversibility: behavior-preserving changes SHALL be committed in logically grouped commits so any group can be reverted independently.
- **NFR-4:** No new third-party dependency SHALL be added; `go.mod`/`go.sum` change only via `go mod tidy` if an import becomes unused.

## Interfaces
- **Inputs:** the working tree on the current branch; the `/ds-*` review skills listed in FR-1.
- **Outputs:**
  - Modified source files (behavior-preserving only).
  - A consolidated audit report (inline, and optionally `.project/AUDIT.md`).
  - Git commits grouping the applied fixes.
- **No external systems**; the audit runs entirely locally with no network calls.

## Constraints
- Language/runtime: Go 1.26.2, module `github.com/gleicon/tldt` (per `.project/PROJECT.md` and Go language profile).
- Toolchain: `go build`, `go vet`, `go test -race`, `golangci-lint`, `make` targets.
- Forbidden approaches:
  - Changing observable behavior under the banner of "cleanup".
  - Deleting pre-existing dead code that the maintainer hasn't approved — flag it (AGENTS.md §3); the legacy `src/` tree is flagged, not modified.
  - Weakening or deleting tests to pass the suite.
  - Adding speculative abstractions or configurability (AGENTS.md §2).

## Technical Profile
- Primary language: Go 1.26.2
- Runtime target: native binary (darwin/arm64 and goreleaser targets)
- Build toolchain: `go build` / `go install`, `Makefile`, `.goreleaser.yaml`
- Testing framework: standard `testing` (table-driven, `httptest`), `go test -race`, coverage via `make test-cover`

## Acceptance Criteria
- **AC-1:** Given the active code, when each `/ds-*` command in FR-1 is run, then a findings list exists for each command (or an explicit "no findings"). *(verifies FR-1)*
- **AC-2:** Given the collected findings, when the report is produced, then every finding carries a behavior-preserving vs behavior-changing classification. *(verifies FR-2, FR-9)*
- **AC-3:** Given behavior-preserving findings, when fixes are applied, then the working tree contains those changes and only those. *(verifies FR-3, NFR-2)*
- **AC-4:** Given a behavior-changing finding, when the audit completes, then no code implementing that change is present in the diff and the finding appears in the recommendations list. *(verifies FR-4)*
- **AC-5:** Given `git diff` of the audit, when the public API of `pkg/tldt` is compared (e.g. `go doc ./pkg/tldt` before/after), then exported signatures and doc semantics are identical. *(verifies FR-5)*
- **AC-6:** Given a fixed set of representative inputs (stdin pipe, `-f`, `--url` via httptest, each `--format`, each `--algorithm`), when run before and after the audit, then stdout, stderr classification, and exit codes match byte-for-byte. *(verifies FR-6)*
- **AC-7:** Given the suite, when `go test -race ./...` runs before and after, then the same set of tests passes and no previously-passing test fails. *(verifies FR-7, NFR-1)*
- **AC-8:** Given the final tree, when `go build ./...` and `go vet ./...` run, then both exit 0. *(verifies FR-8)*
- **AC-9:** Given the applied changes, when commits are inspected, then they are grouped by finding category and individually revertible. *(verifies NFR-3)*
- **AC-10:** Given the final tree, when `go.mod`/`go.sum` are diffed, then changes (if any) result only from `go mod tidy` removing newly-unused imports. *(verifies NFR-4)*

### Coverage
- FR-1 → AC-1
- FR-2 → AC-2
- FR-3 → AC-3
- FR-4 → AC-4
- FR-5 → AC-5
- FR-6 → AC-6
- FR-7 → AC-7
- FR-8 → AC-8
- FR-9 → AC-2
- NFR-1 → AC-7
- NFR-2 → AC-3
- NFR-3 → AC-9
- NFR-4 → AC-10

## Resolved Decisions (via /ds-grill-me)
- **D-1 (cadence):** Checkpoint workflow — stop after the review pass and present classified findings in `.project/AUDIT.md` before applying any fix; then apply hands-off.
- **D-2 (legacy `src/`):** Leave untouched; flag as a one-line removal recommendation only. Same for legacy `assets/`, `SSL/`.
- **D-3 (report):** Persist to `.project/AUDIT.md`; append "applied" status as fixes land.
- **D-4 (borderline fixes):** Apply-and-verify when a test or the golden I/O snapshot covers the path; downgrade to recommend-only when no behavioral net exists.
- **D-5 (`--url` snapshot):** Excluded from the golden I/O snapshot (network-nondeterministic); parity rides on `internal/fetcher` httptest tests. Updates AC-6.
- **D-6 (`examples/`):** Light scope — doc-accuracy pass against `pkg/tldt` + one `go build ./...` per example module at verify; no deep idiom/test audit.
- **D-7 (lint gate):** Adopt `golangci-lint` as a gate — add `.golangci.yml` (`version: 2`, default linters: `errcheck govet ineffassign staticcheck unused`), wire `make lint` to `golangci-lint run`, make the tree pass. This adds a behavior-preserving requirement, not a behavior change.
- **D-8 (commits):** One commit per fix category on the `cleanup` branch, committed only after that group passes build + vet + lint + `test -race` + golden I/O.

## Open Questions
- None outstanding — all resolved above.
