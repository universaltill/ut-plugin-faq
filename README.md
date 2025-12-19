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
- Manifest validate: `uitill manifest validate src/manifest/manifest.json`

## License
MIT
