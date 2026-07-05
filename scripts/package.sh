#!/usr/bin/env bash
# Package the FAQ plugin into the canonical marketplace release artifact:
#   dist/<plugin-id>_<version>_<os>_<arch>.tar.gz
# per docs/plugins/release-artifact.md. The manifest is the version source of
# truth; package.json must agree or the build fails.
#
# Usage:
#   scripts/package.sh                    # package for the host os/arch
#   TARGET_OS=linux TARGET_ARCH=amd64 scripts/package.sh
#   SKIP_TESTS=1 scripts/package.sh       # when tests already ran (e.g. CI matrix)
set -euo pipefail

cd "$(dirname "$0")/.."

TARGET_OS=${TARGET_OS:-$(go env GOOS)}
TARGET_ARCH=${TARGET_ARCH:-$(go env GOARCH)}
if [[ "$TARGET_OS" == "universal" ]]; then
  echo "ERROR: this plugin ships a compiled binary; universal artifacts are not supported. Set TARGET_OS/TARGET_ARCH to a concrete pair." >&2
  exit 1
fi

WORK=$(mktemp -d)
trap 'rm -rf "$WORK"' EXIT
STAGE="$WORK/stage"
mkdir -p "$STAGE"

echo "==> Building pkgtool"
PKGTOOL="$WORK/pkgtool"
go build -o "$PKGTOOL" ./tools/pkgtool

echo "==> Validating manifest"
"$PKGTOOL" validate

VERSION=$("$PKGTOOL" version)
PLUGIN_ID=$("$PKGTOOL" id)
ARTIFACT="${PLUGIN_ID}_${VERSION}_${TARGET_OS}_${TARGET_ARCH}.tar.gz"

if [[ "${SKIP_TESTS:-}" != "1" ]]; then
  echo "==> Running tests"
  go test ./...
fi

echo "==> Building runtime payload (${TARGET_OS}/${TARGET_ARCH})"
mkdir -p "$STAGE/bin"
GOOS="$TARGET_OS" GOARCH="$TARGET_ARCH" CGO_ENABLED=0 \
  go build -trimpath -ldflags "-s -w" -o "$STAGE/bin/ut-faq" ./src

echo "==> Staging archive contents"
"$PKGTOOL" stage-manifest -out "$STAGE/manifest.json" -os "$TARGET_OS" -arch "$TARGET_ARCH"
cp -RL assets "$STAGE/assets"
mkdir -p "$STAGE/content"
cp src/faq/content/*.json "$STAGE/content/"
cp README.md CHANGELOG.md LICENSE "$STAGE/"
"$PKGTOOL" release-json \
  -out "$STAGE/release.json" -os "$TARGET_OS" -arch "$TARGET_ARCH" -artifact "$ARTIFACT"
# Strip OS junk so archives are identical across machines.
find "$STAGE" \( -name '.DS_Store' -o -name 'Thumbs.db' \) -delete

echo "==> Creating archive"
# Sorted file list + no macOS extended attrs keeps repeated runs structurally identical.
(cd "$STAGE" && find . -type f | sed 's|^\./||' | LC_ALL=C sort) > "$WORK/files"
TMP_ARCHIVE="$WORK/$ARTIFACT"
COPYFILE_DISABLE=1 tar -czf "$TMP_ARCHIVE" -C "$STAGE" -T "$WORK/files"

echo "==> Verifying archive layout"
LISTING=$(tar -tzf "$TMP_ARCHIVE")
for required in manifest.json bin/ut-faq release.json README.md CHANGELOG.md LICENSE \
  assets/icon.png assets/screenshot-en.png assets/screenshot-fr.png content/en-US.json; do
  if ! grep -qxF "$required" <<< "$LISTING"; then
    echo "ERROR: required file missing from archive: $required" >&2
    exit 1
  fi
done

echo "==> Publishing to dist/"
mkdir -p dist
rm -f "dist/${ARTIFACT}" "dist/${ARTIFACT}.sha256"
mv "$TMP_ARCHIVE" "dist/${ARTIFACT}"
if command -v sha256sum >/dev/null 2>&1; then
  (cd dist && sha256sum "$ARTIFACT" > "${ARTIFACT}.sha256")
else
  (cd dist && shasum -a 256 "$ARTIFACT" > "${ARTIFACT}.sha256")
fi
CHECKSUM=$(cut -d' ' -f1 "dist/${ARTIFACT}.sha256")

echo
echo "artifact:  dist/${ARTIFACT}"
echo "sha256:    ${CHECKSUM}"
echo "plugin:    ${PLUGIN_ID} ${VERSION} (${TARGET_OS}/${TARGET_ARCH})"
