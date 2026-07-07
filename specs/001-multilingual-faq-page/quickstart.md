# Quickstart: Multilingual FAQ Page Plugin

## Prerequisites
- Go 1.21+
- Universal Till plugin SDK (Go)
- Access to Universal Till POS dev environment with marketplace mock or staging

## Setup
```bash
# from repo root
export GO111MODULE=on
# (optional) go env -w GOPRIVATE=... if using private modules
```

## Build & Validate
```bash
# format & lint
find . -name '*.go' -not -path './vendor/*' -print | xargs gofmt -w
go vet ./...

# or run helper script
./scripts/dev.sh

# tests
go test ./...

# manifest validation (illustrative command; replace with real tool)
uitill manifest validate manifest/manifest.json
```

## Run in POS Dev Mode
```bash
# build plugin artifact (example)
go build -o bin/ut-faq ./src

# install to POS (replace path/command with actual tooling)
ut-plugin-cli install ./bin/ut-faq

# start POS pointing to marketplace mock if needed
UT_DEV_MODE=true ./universal-till
```

## Locale Testing
- Set POS locale to each supported locale; open Help/FAQ page and verify full localization.
- For RTL (ar-SA, fa-IR): confirm layout direction, alignment, and icon mirroring.
- For unsupported locale (e.g., ja-JP): verify fallback to en-US with notice.

## Offline Testing
- Install plugin, then disconnect network; open FAQ page and ensure content loads from cache.
- Simulate interrupted update; confirm resume and checksum validation retain prior content.
- Auth: verify FAQ page is accessible only when signed-in staff are present.
- Disable/Uninstall: remove/disable plugin and ensure Help/FAQ entry disappears.
- Checksums: corrupt a bundle and ensure install/update refuses to promote it; verify rollback keeps prior content.

## Packaging for Marketplace
- Ensure manifest has `canonical_type` `page`, capability `page`, minimal permissions, required assets (icon 512x512, screenshots, README, license, changelog).
- Generate artifact with checksums/signatures per marketplace tooling.
- Submit to marketplace; verify zero MV-* errors.
