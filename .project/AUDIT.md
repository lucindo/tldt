# Audit Report — tldt (behavior-preserving full code audit)

## Outcome (applied)

All ✅ A groups applied + approved ⚠️ R1/R2/R4. Each group committed separately on `cleanup`, verified (build + vet + golangci-lint + `test -race` + golden I/O) before commit. Final state: lint 0 issues, all tests pass (366 test funcs), `go.mod`/`go.sum` unchanged, exported `pkg/tldt` API identical, golden I/O for all 15 baseline scenarios identical.

| Commit | Group |
|--------|-------|
| `d2b03c9` | A1 lint gate (`.golangci.yml` v2 defaults, `make lint`, gofmt, errcheck/staticcheck) |
| `3b88964` | A5 docs (threshold drift, missing flags, dup/stale cleanup) |
| `dd1aa2a` | A2/A3 dead-code + slop |
| `f65f1b1` | A4 test quality (2 no-op tests converted, 9 coverage tests) |
| `f2298ed` | **R1** reject non-positive `--sentences` (was a panic) |
| `829a816` | **R4** fetcher real metadata + sentinel errors (no fabricated values / string-matching) |
| `ff1790c` | **R2** dial-time SSRF validation (close DNS-rebinding TOCTOU) |

**Note on A2.2:** the `trRowNormalize`/`trSelectTopN` *function collapse* (CQ-6) was reduced to fixing the false comments only; the structural dedup belongs with the deferred 🔁 C refactors (test/call-site churn). The `tr*` functions remain.

### Remaining recommendations (not applied — your call)
Still open from section B: **R3** (SSRF blocklist gaps: `0.0.0.0`, IPv4-mapped, NAT64, CGNAT), **R5** (SanitizePII single-pass), **R6** (PII coverage: `sk-proj`/gh/AWS/PEM/SSN/Luhn), **R7** (`Detect` unused opts + nil error — API change), **R8** (UTF-8 byte-slice truncation in excerpts/`htmlmd` MaxLength), **R9** (installer `--target` installs claude anyway), **R10** (explain `Selected` keyed by string), **R11** (graph contract), **R12** (htmlmd swallows error), **R13** (legacy `src/` removal). Plus 🔁 C1–C3 refactors (deferred).

Discovered during apply:
- **R14 (BUG-5):** the `examples/openapi-client` happy path fetches a `swagger.json`, but `Fetch` only accepts `text/html` (now `ErrNonHTML`) — the documented example fails at runtime.
- **R15:** `examples/basic`, `examples/pipeline`, `examples/openapi-client` have stale `go.mod` (`go mod tidy` needed) and don't build — **pre-existing**, unrelated to audit changes (their modules were never touched). `examples/html-processor` builds.

---


Baseline: build/vet/`test -race` green · golangci-lint 19 issues (17 errcheck, 2 staticcheck) · 8 files not gofmt-clean.
Reviews run: bug, security, go-idiom, code-quality, test-quality, doc-quality, deslop. Findings deduped across reviews and classified.

**Legend:** ✅ = behavior-preserving (safe to apply) · ⚠️ = behavior-changing (recommend only, NOT applied) · 🔁 = large behavior-preserving refactor (opt-in).

---

## A. Safe to apply — behavior-preserving (proposed apply set)

### A1 · Lint gate adoption (D-7)
- **A1.1** Add `.golangci.yml` (`version: 2`, default linters), wire `make lint` → `golangci-lint run`.
- **A1.2** `gofmt -w` the 8 non-clean files (`cmd/tldt/main.go`, `internal/htmlmd/htmlmd.go`, `internal/installer/installer.go`, `internal/sanitizer/sanitizer.go` + test, `internal/summarizer/textrank.go`, `pkg/tldt/doc.go`, `pkg/tldt/tldt.go`). Fixes the misaligned flag block (CQ-9, SLOP-5). Whitespace-only.
- **A1.3** Fix the one production errcheck — `internal/fetcher/fetcher.go:111` (unchecked return).
- **A1.4** Fix 16 test-file errchecks (unchecked `os.Setenv`/`os.WriteFile`/etc. in config/fetcher/installer/main tests).
- **A1.5** Fix 2 staticcheck: `false_positive_test.go:88` (S1009 drop redundant `== nil` before `len()`), `installer/embed.go:7` (SA9009 reword the `// go:embed …` comment so it isn't read as a directive).
- *Net:* tree passes `golangci-lint run` + existing suite. *(GO-9, baseline)*

### A2 · Dead / speculative internal code
- **A2.1** Delete `detector.Report.QuarantinedIdxs` field + its misleading "quarantine=true" comment — never populated, references a nonexistent feature. *(CQ-5)*
- **A2.2** Delete duplicate `trRowNormalize`/`trSelectTopN` in `textrank.go`; call shared `rowNormalize`/`selectTopN`; drop the "avoid conflicts if LexRank runs in parallel" comment (false — same package, no parallelism). *Verify byte-identical before applying.* *(CQ-6)*

### A3 · Slop / micro-idiom (optional, trivial)
- **A3.1** Inline `count := len(findings)` at `main.go:247-252`. *(SLOP-4)*
- **A3.2** `base64` padding via `strings.Repeat` instead of `+=` loop, `detector.go:239-243`. *(GO-12)*
- **A3.3** Rename `applySentenceCap(text, cap)` param — `cap` shadows the builtin → `maxSentences`. *(GO-16)*

### A4 · Test quality (additive, behavior-preserving)
- **A4.1** Convert the two assertion-free tests to real assertions: `TestOutlierThresholdCalibration` (monotonicity) + `TestOutlierScoreDistribution` (bounds `0 ≤ score ≤ 1`). Currently always pass. *(TEST-6)*
- **A4.2** Add coverage for untested controls: fetcher `maxBytes` cap (TEST-1), CLI `--sanitize` (TEST-2), CLI `--detect-injection` incl. outlier hand-off (TEST-3), `SanitizePII` overlapping matches (TEST-4), Pipeline `Redactions` populated path (TEST-5), tokenizer abbrev/decimal boundaries (TEST-7), graph centrality unit (TEST-8).

### A5 · Documentation accuracy (doc-only)
- **A5.1** Fix threshold drift `0.85 → 0.99` in `README.md:76`, `pkg/tldt/tldt.go:39`, `pkg/tldt/doc.go:27` (code default is `0.99`). *(DOC-1, DOC-2)*
- **A5.2** Fix the wrong `--detect-pii` example output in `docs/security.md:35` to match real `pii-detect: N finding(s)` + per-finding format. *(DOC-3)*
- **A5.3** Add missing flags to README Flags table + short sections: `--detect-pii`, `--sanitize-pii` (DOC-6), `--from-html` (DOC-7), `--skill-dir`, `--target` (DOC-8).
- **A5.4** Remove the stale/duplicate package doc comment block in `pkg/tldt/tldt.go:1-9` (keep `doc.go` as single source; it's the fuller, correct one). *(DOC-4)*
- **A5.5** Drop the stale hardcoded "361 unit tests" count in `CHANGELOG.md:47`. *(DOC-5)*
- **A5.6** Cut the duplicated homoglyph paragraph `README.md:252-253` (already covered by the Confusable table row). *(DOC-9)*
- **A5.7** Align `examples/README.md` documented invocations with each example's actual flags (esp. html-processor `-url-mode`, pipeline lacking `-f`/`-url`). *(DOC-10)*

---

## B. Recommendations — behavior-changing (NOT applied; your call)

| ID | Sev | Finding | Source | Why behavior-changing |
|----|-----|---------|--------|------------------------|
| **R1** | high | `--sentences -1` (and other `<1`) **panics** (`makeslice: len out of range`) via `trSelectTopN`. Add `< 1` validation → clean error exit. | BUG-1 | panic → error exit |
| **R2** | high (sec) | **DNS-rebinding SSRF (TOCTOU):** fetcher validates the resolved IP, then the HTTP transport re-resolves independently at dial. Fix: custom `DialContext` that dials the *validated* IP. Legit public URLs unaffected. | SEC-1, GO-5 | new dial path; closes a real bypass |
| **R3** | med (sec) | SSRF blocklist gaps: `0.0.0.0`/`::` unspecified, IPv4-mapped IPv6, NAT64 `64:ff9b::/96`, CGNAT `100.64/10`, benchmark `198.18/15`. | SEC-2 | newly blocks some URLs |
| **R4** | med | **Fetcher API cluster** (interlocking): `pkg/tldt.Fetch` classifies errors by `strings.Contains(err.Error(),…)` instead of `errors.Is`; `FetchResult.{StatusCode,ContentType,FinalURL}` are **fabricated constants** (`200`/`text/html`/input URL — wrong after redirects). Fix: sentinel errors in `internal/fetcher` (`%w`), return real metadata or drop the fake fields. | GO-1, GO-2, CQ-1, SLOP-1, SLOP-2 | changes fetcher signature + public `FetchResult` shape |
| **R5** | med | `SanitizePII` redacts in a second whole-text pass independent of `DetectPII`; overlapping patterns can make reported findings ≠ actual redactions. Fix: single-pass span-based redaction. | GO-10, TEST-4 | redaction output can change on overlaps |
| **R6** | med (sec) | PII/secret coverage gaps: modern `sk-proj-…` keys (underscore/hyphen) unmatched; no GitHub `ghp_`/AWS/PEM/SSN patterns; credit-card has no Luhn check. | SEC-3 | more strings redacted → summary can shift |
| **R7** | low | `pkg/tldt.Detect(text, opts)` ignores `opts.OutlierThreshold` and never returns a non-nil `error`. Either honor `opts` or drop the param + error return. | GO-3, CQ-4, SLOP-3 | public API signature change |
| **R8** | low | UTF-8 **byte-slice truncation** splits multibyte runes → invalid UTF-8 in finding/`--explain` excerpts (`detector.go` ×6, `explain.go`) and in `htmlmd` `MaxLength` (public `ConvertHTML`). Fix: rune-aware `truncate` helper (also DRYs 6+ copies). | BUG-3, BUG-4, GO-11, CQ-8 | changes output bytes for multibyte input |
| **R9** | low | `installer --target opencode/cursor/agents` still installs the full Claude skill + hook + patches `~/.claude/settings.json` (claude entry unconditionally prepended); dead empty branch at `installer.go:90-96`. | BUG-2, GO-8 | changes install behavior |
| **R10** | low | `--explain` marks `Selected` by matching sentence **text**, so duplicate sentences are mis-flagged. Fix: mark by selected index (already computed). | GO-13 | changes explain output on dup sentences |
| **R11** | low | `graph` algorithm bypasses the `Summarizer` contract (no empty-input/`n`-cap handling; delegates raw to `didasy/tldr`), diverging from the documented "n > count ⇒ all sentences" guarantee. | GO-14 | could change graph output at edges |
| **R12** | low | `htmlmd` silently swallows `readability.FromReader` error and falls back to raw HTML with no signal. Decide: keep+comment, or surface a warning. | GO-15 | depends on chosen fix |
| **R13** | note | Legacy `src/` resumator tree (`//go:build ignore`) — per **D-2**, left untouched; flagged here for a future separate removal (with `assets/`, `SSL/`). | CQ-10 | repo cleanup, deliberate |

---

## C. Larger behavior-preserving refactors — opt-in (🔁)

These have a verification net (summarizer/CLI tests + golden I/O) and would meaningfully improve the code, but are big diffs. Apply only if you want them this pass.

- **C1** Collapse the **3–4× duplicated scoring pipeline**: `SummarizeExplain`/`SummarizeWithMatrix`/`ensemble` re-implement the same tokenize→IDF→matrix→power-iterate→select per algorithm. Extract shared `lexrankScores`/`textrankScores`. *(CQ-2)*
- **C2** Decompose the **~400-line `main`** into `resolveConfig` / `runSecurityStages` / `summarize` / `writeOutput` + a `usage()` for the 88-line help string. *(CQ-3, GO-6)*
- **C3** Collapse the triplicated summarizer-dispatch type-assertion / `err2` blocks (rides on C2). *(CQ-7)*

---

## Verified clean (no action)
SSRF redirect cap (≤5) + `LimitReader` wiring + body `Close` are correct · no data races (`confusableOnce` lazy-init sound) · hook script `tldt-hook.sh` has no command-injection path (uses `printf '%s'`, `json.dumps`) · no path traversal in installer · RE2 ⇒ no ReDoS · internal packages well-separated · JSON/markdown output examples in README match `formatter.go`.

---

## Proposed apply order (per D-8: one commit per group, verify before commit)
1. **Lint gate** (A1) — `.golangci.yml`, gofmt, errcheck/staticcheck → `make lint` green.
2. **Docs** (A5) — accuracy + missing flags.
3. **Dead-code / slop** (A2, A3).
4. **Test quality** (A4).
5. *(optional)* **Refactors** (C1–C3).

All ⚠️ (B) items left for your decision and tracked here.
