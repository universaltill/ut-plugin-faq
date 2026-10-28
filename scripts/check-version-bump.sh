#!/usr/bin/env bash
#
# Requires a manifest.json version bump on any PR that touches a shipped file
# (ut-docs#1940, rolled out here per ut-docs#4917).
#
# WHY THIS EXISTS: scripts/package.sh bundles manifest.json, content/,
# assets/, locales/, README.md, CHANGELOG.md and LICENSE into the marketplace
# release artifact. None of this repo's other checks (validate, package dry
# run, authors) look at the version field, so a PR that changes what ships
# but forgets to bump manifest.json's version lands on main, passes every
# check, and then SILENTLY NEVER SHIPS: auto-tag-release.yml reads the
# unchanged version, sees v<version> already tagged/released, and correctly
# does nothing. main and the marketplace diverge with no failing signal
# anywhere (ut-docs#1940).
#
# WHAT COUNTS AS "SHIPPED": deliberately keyed off the same set
# scripts/package.sh actually bundles, not "any file" -- a docs/code-reviews/
# record, a workflow edit, CLAUDE.md or other repo tooling must never force a
# version bump.
#
# Keep this list mirrored to package.sh's own entries (`entries=(...)` and
# every `entries+=(...)` line) -- nothing enforces that automatically, so a
# change to package.sh's bundle needs a matching edit here
# (check-version-bump.test.sh's "entries mirror" case re-reads package.sh's
# own source and fails on drift in either direction). The check logic itself
# lives in scripts/check-version-bump.body.sh, sourced at the bottom: one
# canonical body, byte-identical in every ut-plugin-* repo (canonical copy:
# ut-docs scripts/templates/check-version-bump.body.sh, ut-docs#4489). Never
# edit the body here; only SHIPPED_PATTERNS (and SHIPPED_NOTE) are per-repo.
set -euo pipefail
cd "$(dirname "$0")/.."

# SHIPPED_PATTERNS is read by the sourced body, which shellcheck does not
# follow without -x.
# shellcheck disable=SC2034
SHIPPED_PATTERNS=(
    'manifest.json'
    'content/*'
    'assets/*'
    'locales/*'
    'README.md'
    'CHANGELOG.md'
    'LICENSE'
)

# shellcheck source=scripts/check-version-bump.body.sh
. scripts/check-version-bump.body.sh
