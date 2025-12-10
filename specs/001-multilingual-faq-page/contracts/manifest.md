# Contract: Plugin Manifest (FAQ Page)

## Required Fields
- `id`: Reverse-DNS unique plugin id.
- `name`: Human-readable name (max 50 chars).
- `version`: Semver (MAJOR.MINOR.PATCH).
- `description`: ≤200 chars summarizing FAQ purpose.
- `canonical_type`: `ui-extension` (page plugin).
- `capabilities`: include `page` (and any marketplace-required UI capability identifiers).
- `permissions`: minimal — UI render + local storage; optional telemetry if mandated.
- `locales`: ["en-US","en-GB","fr-FR","ar-SA","fa-IR","tr-TR","es-ES","it-IT","pt-PT"].
- `supported_architectures`: ["linux/amd64","linux/arm64","darwin/amd64"].
- `min_host_version`: Universal Till minimum compatible version.
- `resources`: icon (512x512), screenshots[], documentation (README/FAQ), license, changelog.
- `configuration_schema`: optional (no runtime config expected).

## Validation Expectations
- Zero MV-* errors on submission.
- Checksum/signature produced per marketplace packaging pipeline.
- Permissions audited for least privilege; network permissions not requested.

## Example (illustrative)
```json
{
  "id": "com.universaltill.ut-faq",
  "name": "Universal Till FAQ",
  "version": "1.0.0",
  "description": "Localized FAQ page for Universal Till POS",
  "canonical_type": "ui-extension",
  "capabilities": ["page"],
  "permissions": ["storage.local.10MB"],
  "locales": ["en-US","en-GB","fr-FR","ar-SA","fa-IR","tr-TR","es-ES","it-IT","pt-PT"],
  "supported_architectures": ["linux/amd64","linux/arm64","darwin/amd64"],
  "min_host_version": "1.0.0",
  "resources": {
    "icon": "assets/icon.png",
    "screenshots": ["assets/screenshot-en.png","assets/screenshot-fr.png"],
    "documentation": "README.md",
    "license": "LICENSE",
    "changelog": "CHANGELOG.md"
  }
}
```
