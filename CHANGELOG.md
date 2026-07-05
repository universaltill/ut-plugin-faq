# Changelog

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
