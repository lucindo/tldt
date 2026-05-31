# Plan

Source: `.project/SPEC.md` — Full Code Audit (Behavior-Preserving)

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
