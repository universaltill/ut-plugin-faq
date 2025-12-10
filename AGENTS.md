# ut-plugin-faq Development Guidelines

Auto-generated from all feature plans. Last updated: 2025-12-10

## Active Technologies
- Go 1.21 + Universal Till plugin SDK (Go), embedded JSON locale bundles (no go-i18n), stdlib embed/checksum (001-multilingual-faq-page)
- Local plugin cache via POS storage; bundled locale JSON content (001-multilingual-faq-page)

## Project Structure

```text
src/
├── main.go
├── manifest/
│   └── manifest.json
├── faq/
│   ├── content/
│   ├── loader.go
│   └── render.go
├── ui/
│   └── page.go
└── storage/
    └── cache.go

tests/
├── integration/
│   └── faq_page_test.go
└── unit/
    ├── loader_test.go
    └── render_test.go
```

## Commands

# Add commands for build/test/manifest validation as they are defined

## Code Style

- Go fmt/vet; keep permissions minimal and manifest compliant

## Recent Changes
- 001-multilingual-faq-page: Added Go SDK offline/i18n FAQ page structure and tech stack

<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
