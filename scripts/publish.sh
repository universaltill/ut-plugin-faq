#!/usr/bin/env bash
# Publish a packaged FAQ plugin release artifact to the marketplace.
#
# Uploads the artifact + manifest + checksum to the marketplace vendor release
# upload API. Reusable from CI and local/dev shells.
#
# Required environment:
#   MARKETPLACE_BASE_URL      e.g. http://localhost:8081
# Optional environment:
#   MARKETPLACE_UPLOAD_TOKEN  bearer token when the marketplace enforces one
#   MARKETPLACE_CHANNEL       stable (default) | beta | alpha
#   MARKETPLACE_LISTING_ID    existing listing UUID; omitted on first publish
#   RELEASE_NOTES             release notes text; defaults to the CHANGELOG
#                             section for the manifest version
#   ARTIFACT                  path to the .tar.gz; defaults to the artifact for
#                             the manifest version and host os/arch in dist/
#   TARGET_OS / TARGET_ARCH   used to locate the default artifact
set -euo pipefail

cd "$(dirname "$0")/.."

: "${MARKETPLACE_BASE_URL:?MARKETPLACE_BASE_URL is required (e.g. http://localhost:8081)}"

TARGET_OS=${TARGET_OS:-$(go env GOOS)}
TARGET_ARCH=${TARGET_ARCH:-$(go env GOARCH)}
CHANNEL=${MARKETPLACE_CHANNEL:-stable}

VERSION=$(go run ./tools/pkgtool version)
PLUGIN_ID=$(go run ./tools/pkgtool id)
ARTIFACT=${ARTIFACT:-dist/${PLUGIN_ID}_${VERSION}_${TARGET_OS}_${TARGET_ARCH}.tar.gz}

if [[ ! -f "$ARTIFACT" ]]; then
  echo "ERROR: artifact not found: $ARTIFACT (run scripts/package.sh first)" >&2
  exit 1
fi

if [[ -f "${ARTIFACT}.sha256" ]]; then
  CHECKSUM=$(cut -d' ' -f1 "${ARTIFACT}.sha256")
elif command -v sha256sum >/dev/null 2>&1; then
  CHECKSUM=$(sha256sum "$ARTIFACT" | cut -d' ' -f1)
else
  CHECKSUM=$(shasum -a 256 "$ARTIFACT" | cut -d' ' -f1)
fi

# Default release notes: the CHANGELOG section for this version.
if [[ -z "${RELEASE_NOTES:-}" ]]; then
  RELEASE_NOTES=$(awk -v ver="## ${VERSION}" '
    $0 == ver {found=1; next}
    /^## / {if (found) exit}
    found {print}
  ' CHANGELOG.md)
  RELEASE_NOTES=${RELEASE_NOTES:-"Release ${VERSION}"}
fi

UPLOAD_URL="${MARKETPLACE_BASE_URL%/}/ui/api/vendor/releases/upload"

CURL_ARGS=(
  --silent --show-error
  --write-out '\n%{http_code}'
  --form "artifact=@${ARTIFACT};type=application/gzip"
  --form "manifest=<src/manifest/manifest.json"
  --form "plugin_id=${PLUGIN_ID}"
  --form "version=${VERSION}"
  --form "channel=${CHANNEL}"
  --form "expected_hash=${CHECKSUM}"
  --form "release_notes=${RELEASE_NOTES}"
)
if [[ -n "${MARKETPLACE_LISTING_ID:-}" ]]; then
  CURL_ARGS+=(--form "listing_id=${MARKETPLACE_LISTING_ID}")
fi
if [[ -n "${MARKETPLACE_UPLOAD_TOKEN:-}" ]]; then
  CURL_ARGS+=(--header "Authorization: Bearer ${MARKETPLACE_UPLOAD_TOKEN}")
fi

echo "==> Publishing ${PLUGIN_ID} ${VERSION} (${CHANNEL}) to ${UPLOAD_URL}"
echo "    artifact: ${ARTIFACT}"
echo "    sha256:   ${CHECKSUM}"

RESPONSE=$(curl "${CURL_ARGS[@]}" "$UPLOAD_URL")
HTTP_CODE=${RESPONSE##*$'\n'}
BODY=${RESPONSE%$'\n'*}

echo "    response: HTTP ${HTTP_CODE}"
echo "$BODY"

if [[ "$HTTP_CODE" != 2* ]]; then
  echo "ERROR: upload failed with HTTP ${HTTP_CODE}" >&2
  case "$HTTP_CODE" in
    401) echo "hint: check MARKETPLACE_UPLOAD_TOKEN" >&2 ;;
    422) echo "hint: validation failed; see error.details above" >&2 ;;
    000) echo "hint: marketplace unreachable at ${MARKETPLACE_BASE_URL}" >&2 ;;
  esac
  exit 1
fi

# Surface release id/status for CI logs without requiring jq.
python3 - "$BODY" <<'EOF' 2>/dev/null || true
import json, sys
data = json.loads(sys.argv[1]).get("data", {})
print(f"release_id: {data.get('release_id')}")
print(f"status:     {data.get('status')} (scan: {data.get('scan_status')})")
EOF
