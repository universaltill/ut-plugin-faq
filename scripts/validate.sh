#!/usr/bin/env bash
set -euo pipefail

# Resolve the optional manifest argument against the caller's cwd before
# changing directory, so relative paths keep working.
MANIFEST_ARG=${1:-}
if [[ -n "$MANIFEST_ARG" && "$MANIFEST_ARG" != /* ]]; then
  MANIFEST_ARG="$(pwd)/$MANIFEST_ARG"
fi

cd "$(dirname "$0")/.."
MANIFEST=${MANIFEST_ARG:-src/manifest/manifest.json}

# Prefer the platform CLI when installed; otherwise use the repo-local
# validator so validation never silently disappears. Note the fallback checks
# the packaging + POS install contract; marketplace-side validation runs again
# authoritatively at upload time.
if command -v uitill >/dev/null 2>&1; then
  uitill manifest validate "$MANIFEST"
else
  go run ./tools/pkgtool validate -manifest "$MANIFEST"
fi
