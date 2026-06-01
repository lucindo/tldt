# Plan

Source: `.project/SPEC.md` — Full Code Audit (Behavior-Preserving)

## Now

**State** — ALL audit recommendations applied. 16 follow-up commits on `cleanup` (`163cd34`→`47712ec`): comment-ID strip, R15, R3, R6, R5, R7(A), R8, R9, R10, R11, R12, R13, R14, C1, C2/C3, plus doc updates. Combined with the earlier 9 commits, the whole audit (A-groups, R1–R15, C1–C3) is shipped. Tree green on every commit: lint 0, `test -race` all pass, 15/15 golden byte-identical, `pkg/tldt` API stable. Nothing from sections B/C remains open — see AUDIT.md "Section B + C — ALL APPLIED".

**Next** — Run the gopls/`go fix` pre-pass: `go fix ./...` across all modules, then resolve every gopls hint (`gopls check -severity=hint $(find . -name "*.go")`) — accumulated hints are rangeint modernization + `interface{}`→`any` (installer `PatchSettingsJSON`, openapi-client struct). Apply behavior-preserving fixes, verify build+lint+`test -race`+golden+API, commit as its own group. THEN the final `/ds-*` review pass — sequential, one command at a time, **deslop last** (order in Roadmap).

**Open questions** — None blocking. Maintainer chose: R7 option A, R6 incl. SSN+Luhn, R13 full legacy removal, C1–C3 done now. The gopls pre-pass + final reviews are the only remaining roadmap items.

## Roadmap

### Baseline
- [x] Behavior baseline captured: build/vet/`test -race` all green; golangci-lint baseline = 19 issues (17 errcheck, 2 staticcheck)
- [x] Public API snapshot captured: `go doc -all ./pkg/tldt` saved (`.audit-baseline/pkg-tldt-api.txt`)
- [x] Golden I/O snapshot captured: 15 scenarios (stdin, `-f`, 4 algorithms × 3 formats, verbose, sanitize/detect) — stdout SHAs + stderr recorded; `--url` excluded per D-5

### Review pass (collect findings, no edits)
- [x] `/ds-bug-review` findings collected over active code
- [x] `/ds-security-review` findings collected over active code
- [x] `/ds-go-review` findings collected over active code
- [x] `/ds-code-quality-review` findings collected over active code
- [x] `/ds-test-quality-review` findings collected
- [x] `/ds-doc-quality-review` findings collected (README, docs/, package docs)
- [x] `/ds-deslop` findings collected
- [x] All findings consolidated and classified → `.project/AUDIT.md` (A=apply, B=recommend, C=opt-in refactor)
- [~] **CHECKPOINT: maintainer reviews AUDIT.md before any edit**

### Apply safe fixes (behavior-preserving only) — DONE
- [x] golangci-lint gate adopted (commit d2b03c9): `.golangci.yml` v2 defaults, `make lint` wired, tree green
- [x] Doc accuracy fixes (commit 3b88964): threshold drift, missing flags, dup/stale doc cleanup
- [x] Deslop + dead-code fixes (commit dd1aa2a): QuarantinedIdxs, false comments, micro-idioms
- [x] Test-quality fixes (commit f65f1b1): assertion-free tests converted; 9 coverage tests added
- [x] `go.mod`/`go.sum` unchanged (no imports became unused)

### Apply approved behavior-changing fixes (R1, R2, R4 — user-approved) — DONE
- [x] R1 (commit f2298ed): reject `--sentences < 1`; summarizer select clamps negative n
- [x] R4 (commit 829a816): fetcher returns real metadata + sentinel errors; pkg/tldt.Fetch uses errors.Is
- [x] R2 (commit ff1790c): SSRF dial-time IP validation (closes DNS-rebinding TOCTOU)

### Verify — DONE
- [x] API snapshot unchanged: exported `pkg/tldt` func/type set identical (FetchResult shape unchanged; values now truthful per R4)
- [x] Golden I/O unchanged: all 15 baseline scenarios match byte-for-byte
- [x] Suite parity: `go test -race ./...` passes (366 test funcs)
- [x] Build + vet + lint green on final tree
- [x] `go.mod`/`go.sum` unchanged after `go mod tidy` (NFR-4)
- [~] examples: `html-processor` builds; basic/pipeline/openapi-client have pre-existing stale go.mod (flagged R15)

### Report & ship — DONE
- [x] Consolidated audit report produced → `.project/AUDIT.md` (applied + remaining recommendations)
- [x] Fixes committed in 7 independently-revertible groups on `cleanup`
- [x] Remaining behavior-changing recommendations (R3, R5–R15, C1–C3) handed to maintainer in AUDIT.md

### Ending cleanup
- [x] Strip leaked planning IDs/lingo from code comments (commit 163cd34) — GSD IDs + audit R1 IDs removed across 13 files; comments only, golden/API unchanged.
- [x] R15 (commit 20d1c23): `go mod tidy` examples basic/pipeline/openapi-client (pre-existing stale go.mod). main module untouched.
- [x] R3 (commit 9228002): SSRF blocklist — unspecified/CGNAT/NAT64/benchmark; IPv4-mapped covered by regression tests.
- [x] R6 (commit 27ec822): PII coverage — modern sk-/GitHub/PEM/SSN + Luhn credit-card; detect/redact kept consistent.
- [x] Reviewed and applied ALL remaining B/C items one-by-one with maintainer: R3,R5,R6,R7(A),R8,R9,R10,R11,R12,R13,R14 + C1,C2,C3 (commits 9228002→68c1f01). See AUDIT.md "Section B + C — ALL APPLIED". Tree green throughout; golden + API parity held on every commit.
- [ ] **Pre-pass hygiene (before reviews):** run `go fix ./...` across all modules, then resolve every gopls hint: `gopls check -severity=hint $(find . -name "*.go")`. Apply behavior-preserving fixes; verify build + lint + `test -race` + golden I/O + API unchanged. Commit as its own group.
- [ ] Final `/ds-*` pass over the full audit diff — **run sequentially, one command at a time** (may delegate each to an agent, but never in parallel). **`/ds-deslop` runs LAST.** Order:
  1. [ ] `/ds-code-quality-review` — maintainability of the changes
  2. [ ] `/ds-go-review` — Go idioms (fetcher transport/SSRF, error wrapping, clamps, new PII paths)
  3. [ ] `/ds-bug-review` — correctness of the behavior-changing commits (R1/R2/R3/R4/R5/R6)
  4. [ ] `/ds-security-review` — re-verify SSRF dial path + blocklist (R2/R3) and PII detect/redact surface (R4/R6)
  5. [ ] `/ds-test-quality-review` — coverage of the new behavior
  6. [ ] `/ds-deslop` — slop in new/edited code (LAST)
  - [ ] Consolidate residual findings into AUDIT.md; apply behavior-preserving fixes, recommend the rest
