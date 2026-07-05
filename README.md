# Universal Till FAQ Plugin

[![CI](https://github.com/universaltill/ut-plugin-faq/actions/workflows/ci.yml/badge.svg)](https://github.com/universaltill/ut-plugin-faq/actions/workflows/ci.yml)

Multilingual, offline-capable FAQ page plugin for Universal Till POS.

## Features
- Marketplace-installable UI page entry under Help/Support
- Localized FAQ content (en-US, en-GB, fr-FR, ar-SA, fa-IR, tr-TR, es-ES, it-IT, pt-PT) with RTL support
- Offline rendering with cache and checksum validation

## Usage
- Build: `go build -o bin/ut-faq ./src`
- Tests: `go test ./...`
- Manifest validate: `scripts/validate.sh` (uses `uitill` when installed, otherwise the repo-local `tools/pkgtool` validator)

## Build, Validate, and Package a Release

The plugin packages into the canonical marketplace release artifact defined in
`docs/plugins/release-artifact.md`:

```
dist/com.universaltill.ut-faq_<version>_<os>_<arch>.tar.gz
```

Steps:

```bash
# 1. Validate manifest + version alignment (manifest.json is the source of truth;
#    package.json must match)
scripts/validate.sh

# 2. Package for the host platform
scripts/package.sh

# 3. Or cross-package for a till target
TARGET_OS=linux TARGET_ARCH=amd64 scripts/package.sh
```

`package.sh` validates the manifest, runs the test suite (skip with
`SKIP_TESTS=1` when a CI matrix already ran them), builds the runtime payload
into `bin/ut-faq`, stages `manifest.json` (with `device_arch` rewritten for the
target and checked against `supported_architectures`), `assets/`, `content/`
(all locale bundles), docs, `LICENSE`, and a generated `release.json`, then
creates the versioned `.tar.gz` in `dist/` with a `.sha256` checksum file next
to it. The artifact only lands in `dist/` after the archive layout check
passes. The build fails if the manifest is invalid, versions are misaligned,
tests fail, or a required file is missing from the archive.

Notes:
- This plugin ships a compiled binary, so only concrete `os/arch` artifacts are
  produced; the `_universal` naming from the release contract does not apply.
- The archive checksum lives in the external `.sha256` sidecar (and is computed
  authoritatively by the marketplace at upload); `release.json` cannot contain
  its own archive hash.

Version bump checklist: update `version` in `src/manifest/manifest.json`, the
matching `version` in `package.json`, and add a `CHANGELOG.md` entry.

## License
MIT — see [LICENSE](LICENSE)
