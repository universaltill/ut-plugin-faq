# Changelog

## 0.2.3
- Populate the real `checksum_sha256` per locale bundle (was `""` in all 9
  files despite the data model requiring it — spec audit gap). New
  `scripts/checksum.py` computes it via byte-level substitution (the
  checksum field is zeroed to a same-length placeholder before hashing, so
  it doesn't need to hash itself) and `--check` mode is now wired into
  `validate.sh` so future content edits can't drift without updating it.

## 0.2.0
- Converted to an asset-only plugin: `runtime: "none"`, Go source, binary and
  `tools/pkgtool` removed (ADR-0001 — the POS renders `content/<locale>.json`
  natively, the shipped binary was never executed)
- One universal artifact per release instead of per-os/arch archives
- Content bundles moved from `src/faq/content/` to `content/`
- Release pipeline aligned with the other plugin repos (validate → package →
  publish → dev auto-approve)

## 0.1.2
- Add LICENSE (MIT) and align versions across manifest and package metadata
- Align manifest with the POS host install contract: `canonical_type: page`, `executable`, `entrypoint`, `device_arch`, `min_pos_version`, navigation `entries`
- Add release packaging (`scripts/package.sh`) producing the canonical marketplace artifact with `release.json` and SHA-256 checksum
- Add self-contained manifest/packaging validator (`tools/pkgtool`)

## 0.1.1
- Add locale-aware navigation labels and search rendering logic
- Add checksum-gated update flow with rollback safety
- Add version label for FAQ page metadata

## 0.1.0
- Initial scaffold for Universal Till FAQ plugin
- Embedded locale bundles and cache/loader
- Navigation entry and auth gating helpers
