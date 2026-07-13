# Universal Till FAQ Plugin

[![CI](https://github.com/universaltill/ut-plugin-faq/actions/workflows/ci.yml/badge.svg)](https://github.com/universaltill/ut-plugin-faq/actions/workflows/ci.yml)

Multilingual, offline-capable FAQ page plugin for Universal Till POS.

A `page`-type plugin (ADR-0009 naming applies to newer repos): registers a
Help/Support menu entry; the POS renders the localized content bundle
(`content/<locale>.json`, RTL-aware) natively. **Asset-only, `runtime:
"none"`** — there is no plugin binary (ADR-0001; the Go runtime this repo
shipped up to 0.1.2 was never executed by the till).

## Features
- Marketplace-installable UI page entry under Help/Support
- Localized FAQ content (en-US, en-GB, fr-FR, ar-SA, fa-IR, tr-TR, es-ES, it-IT, pt-PT) with RTL support
- Fully offline: the content ships in the signed bundle

## Validate and Package a Release

One universal artifact per release:

```
dist/com.universaltill.ut-faq_<version>_universal.tar.gz
```

```bash
scripts/validate.sh   # manifest + a valid content bundle per declared locale
scripts/package.sh    # universal tar.gz + sha256 into dist/
```

Version bump checklist: update `version` in `manifest.json` and add a
`CHANGELOG.md` entry, then tag `v<version>`.

## Publish to the Marketplace

`scripts/publish.sh` uploads a packaged artifact to the marketplace vendor
release API. The same script serves local publishing and CI.

```bash
# Local/dev publish (marketplace running locally, package first)
scripts/package.sh
MARKETPLACE_BASE_URL=http://localhost:8081 scripts/publish.sh
```

Environment variables:

| Variable | Required | Purpose |
| --- | --- | --- |
| `MARKETPLACE_BASE_URL` | yes | Marketplace origin, e.g. `http://localhost:8081` |
| `MARKETPLACE_UPLOAD_TOKEN` | when enforced | Bearer token for the upload endpoint |
| `MARKETPLACE_CHANNEL` | no | `stable` (default), `beta`, or `alpha` |
| `MARKETPLACE_LISTING_ID` | no | Existing listing UUID; first publish auto-creates a draft listing |

### CI pipeline

`.github/workflows/release.yml` runs on `v*` tags (the tag must match the
manifest version) and on manual dispatch: validate → package the universal
artifact → publish via `scripts/publish.sh` → (dev) auto-approve via
`scripts/approve.sh` when the `AUTO_APPROVE` repo variable is `true`.
Secrets: `MARKETPLACE_BASE_URL`, `MARKETPLACE_UPLOAD_TOKEN`.
Vars: `AUTO_APPROVE`, `MARKETPLACE_LISTING_ID`.

## License
MIT — see [LICENSE](LICENSE)
