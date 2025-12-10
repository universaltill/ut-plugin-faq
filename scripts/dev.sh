#!/usr/bin/env bash
set -euo pipefail

# Format Go sources if any exist
if files=$(find . -name '*.go' -not -path './vendor/*' -print); then
  if [ -n "$files" ]; then
    gofmt -w $files
  fi
fi

# Run vet if module is present
if [ -f go.mod ]; then
  go vet ./...
fi
