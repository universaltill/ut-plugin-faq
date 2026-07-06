# Code Review — Story 6-5: Publish Pipeline

Date: 2026-07-06
Reviewer: Claude (autonomous review before commit)
Scope: `scripts/publish.sh`, `.github/workflows/release.yml`, `README.md`
(publish + CI docs).

## What the change delivers

`scripts/publish.sh` uploads a packaged artifact to the marketplace vendor
release API (`/ui/api/vendor/releases/upload`) — the same script for local and
CI use. `.github/workflows/release.yml` runs on `v*` tags and manual dispatch:
test → validate → tag/version check → package (linux/amd64 + arm64) → attach
artifacts → publish. No release logic lives in YAML; CI orchestrates the repo
scripts.

## Contract alignment with the marketplace (verified)

The upload form fields (`artifact`, `manifest` as an inline `<file` value,
`plugin_id`, `version`, `channel`, `expected_hash`, `release_notes`, optional
`listing_id`, optional `Authorization: Bearer`) match exactly what
`handleVendorReleaseUploadAPI` / `IngestArtifact` accept. The response parser
reads `data.release_id` / `data.status` / `data.scan_status`, matching
`IngestResult`. Because validation now runs synchronously (story 6-4),
`scan_status` reported by the script is the settled `passed`/`failed`, not
`pending`.

## Correctness findings

| # | Finding | Disposition |
|---|---------|-------------|
| C1 | **Multi-arch matrix would break CI.** Publishing both `linux/amd64` and `linux/arm64` to the same `(version, channel)` hits the marketplace's global `(version, channel)` unique index (flaw flagged in 6-3), so the second leg fails with "already exists". | **Fixed.** The publish step is guarded to `matrix.target_arch == 'amd64'`; both arches are still attached to the workflow run. Limitation documented in the README and in a workflow comment, tied to the `(listing, version, channel)` uniqueness follow-up. |
| C2 | `expected_hash` is taken from the `.sha256` sidecar (`cut -d' ' -f1`), falling back to `sha256sum`/`shasum`. Matches package.sh's sidecar format. | Verified correct. |
| C3 | HTTP code / body split via `curl --write-out '\n%{http_code}'` and bash parameter expansion; non-2xx exits non-zero with actionable hints (401/422/000). `set -euo pipefail` throughout. | Verified correct. |
| C4 | `manifest=<src/manifest/manifest.json` sends the manifest as a form **value**, which `resolveManifest` reads first (before file/artifact fallback). | Verified correct. |
| C5 | Tag/version guard (`go run ./tools/pkgtool version` vs `${GITHUB_REF_NAME#v}`) prevents tag/manifest drift; runs only on tag pushes. | Verified correct. |

## Cleanup findings

| # | Finding | Disposition |
|---|---------|-------------|
| K1 | The tag/version guard runs once per matrix leg (redundant across arches). | **Accepted** — negligible cost, keeps the job self-contained. |
| K2 | `release_notes` defaults to the CHANGELOG section for the version via `awk`, falling back to `Release <version>`. | Verified correct; good ergonomics. |

## Verification

- `bash -n scripts/publish.sh` — clean.
- `go run ./tools/pkgtool version` → `0.1.2`, `... id` → `com.universaltill.ut-faq`
  (the values publish.sh derives).
- Full end-to-end publish against a running marketplace is exercised in story 6-6.
