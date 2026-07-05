# Code Review — Release Packaging (story 6-2)

Date: 2026-07-05
Scope: LICENSE, manifest POS-contract alignment, version alignment (0.1.2), `tools/pkgtool`, `scripts/package.sh`, `scripts/validate.sh`, README/AGENTS/spec docs.
Method: two independent review passes (correctness: line-scan, removed-behavior audit, cross-repo consumer trace; cleanup: reuse/simplification/efficiency/altitude/conventions) before commit.

## Findings and dispositions

### Fixed

1. **validate.sh resolved relative manifest args against the repo root, not the caller's cwd** (correctness). Fixed: argument is absolutized before `cd`.
2. **Failed layout verification left a broken artifact + valid .sha256 in `dist/`** (correctness). Fixed: archive is built and verified in a temp dir and only moved into `dist/` after the layout check passes.
3. **Layout check used `grep -qx` (regex) so `.` matched any character** (correctness). Fixed: `grep -qxF`.
4. **`release-json` didn't require `-arch` and could emit `"target_arch": ""`** (correctness). Fixed: `-arch` required for both `stage-manifest` and `release-json`.
5. **`pkgtool` didn't validate `entrypoint`, which the POS host hard-requires** (correctness, cross-repo trace to `universal-till/internal/plugins/manifest.go:97`). Fixed: required field.
6. **Staged manifest target wasn't checked against `supported_architectures`** (correctness). Fixed: `stage-manifest` fails for unlisted os/arch pairs; `darwin/arm64` added to the manifest list (it was missing while being our dev platform).
7. **`.DS_Store`/junk files could ship inside the artifact; symlinked assets were staged as links** (correctness/determinism). Fixed: `cp -RL` + junk-file deletion before archiving.
8. **`FILELIST` temp file leaked on failure** (correctness). Fixed: single `$WORK` temp root covered by one EXIT trap.
9. **`go run ./tools/pkgtool` invoked 5× per package run (5 recompiles)** (efficiency). Fixed: built once into the temp dir and reused.
10. **Tests re-ran per target in a packaging matrix** (efficiency). Fixed: `SKIP_TESTS=1` escape hatch, documented.
11. **Unreachable `-os universal` half-support in pkgtool while package.sh could never produce universal artifacts** (dead code). Fixed: removed; `package.sh` now rejects `TARGET_OS=universal` with a clear message; README documents that this plugin ships per-arch artifacts only.
12. **`CheckVersionAlignment` hardcoded cwd-relative `package.json`, ignoring `-root`** (correctness). Fixed: joined with `-root`, skipped when root is empty.
13. **Spec docs still mandated `canonical_type: ui-extension`, contradicting the shipped manifest and the new validator** (consistency). Fixed: `specs/001-multilingual-faq-page/{contracts/manifest.md,spec.md,tasks.md,quickstart.md}` aligned to `page`.
14. **AGENTS.md Commands section was an empty placeholder while this change defined the commands** (conventions). Fixed: build/test/validate/package commands documented.
15. **Copied canonical-type taxonomy had no pointer to its source of truth** (drift risk). Fixed: comment referencing `universal-till/internal/plugins/manifest_verifier.go`.

### Accepted / deferred

16. **pkgtool fallback enforces the packaging + POS install contract, not the marketplace MV-\* rule set** (removed-behavior audit). Accepted: `uitill` does not exist yet, so the previous hard-fail made local validation impossible; marketplace-side validation runs authoritatively at upload (stories 6-3/6-4). Noted in validate.sh comments.
17. **`canonical_type: page` conflicts with the marketplace doc taxonomy (`ut-market-place/docs/manifest-validation-errors.md` lists `ui-extension`)**. Deferred to story 6-4, where that doc is being reworked: the POS host taxonomy (`page|button|payment|report|integration|background_job|device`) is the runtime authority and the marketplace doc will be aligned to it.
18. **`release.json` cannot contain `artifact_sha256` (self-hash impossibility)**. Accepted design decision: the external `.sha256` sidecar and the marketplace-computed checksum are canonical; documented in README. The release contract lists `artifact_sha256` as recommended, not required.

## Verification

- `go vet ./...`, `go test ./...` green (unit, integration, pkgtool).
- Packaged `darwin/arm64` and `linux/amd64`; archive layout verified; repeated runs produce identical structure.
- `TARGET_OS=linux TARGET_ARCH=arm` correctly rejected (not in `supported_architectures`).
- `scripts/validate.sh` works from repo root and from a subdirectory with a relative path.
