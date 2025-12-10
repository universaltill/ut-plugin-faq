#!/usr/bin/env bash
set -euo pipefail

MANIFEST=${1:-src/manifest/manifest.json}
if ! command -v uitill >/dev/null 2>&1; then
  echo "uitill not found; please install to validate manifest" >&2
  exit 1
fi
uitill manifest validate "$MANIFEST"
