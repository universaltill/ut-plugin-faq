# ut-plugin-faq — rules for working in this repo

The first sample **Universal Till plugin** (multilingual FAQ page) and the
template other plugins copy. Full guide: `docs` repo →
`for-plugin-developers.md`; standards → `reference/coding-standards.md`.

## The manifest is the contract (enforced by scripts/validate.sh)
- `src/manifest/manifest.json` is the source of truth. `version` must match
  `package.json` and the `CHANGELOG.md` entry; a `vX.Y.Z` tag must match too.
- Required fields: `id`, `name`, `version`, `canonical_type` (page/button/payment/
  report/integration/background_job/device), `executable`, `entrypoint`,
  `device_arch`, `min_pos_version`, plus navigation `entries`.
- Run `scripts/validate.sh` before packaging; it fails on drift.
- The signed manifest must stay compatible with the POS `plugins.Manifest` struct
  and the marketplace `signing.CanonicalManifest` — don't add fields unilaterally.

## Conventions
- **No hardcoded user-facing strings** — content lives in `content/<locale>.json`
  (en-US, fr-FR, ar-SA). RTL must render correctly.
- Keep requested `permissions` minimal.
- Offline-first: the plugin bundle is self-contained; no network needed to render.
- `gofmt`/`go vet` clean; `tools/pkgtool` is stdlib-only (no third-party deps).

## Build / package / publish
- Build: `go build -o bin/ut-faq ./src`
- Validate: `scripts/validate.sh`
- Package: `scripts/package.sh` (env: `TARGET_OS`, `TARGET_ARCH`, `SKIP_TESTS=1`)
- Publish: `scripts/publish.sh` (`MARKETPLACE_BASE_URL`, `MARKETPLACE_UPLOAD_TOKEN`)

## Before committing
- `go test ./...` and `scripts/validate.sh`. Feature branch; record a code review
  in `docs/code-reviews/<date>-<topic>.md`; then merge to `main`.
