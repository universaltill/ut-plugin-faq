# 2026-07-21 — Populate `checksum_sha256` for every content bundle

## Context
Spec audit (`ut-docs/QUEUE.md`): `data-model.md` requires every
`LocalizationBundle` to carry a real `checksum_sha256`, "checksum validated
on install/update" — but all 9 `content/<locale>.json` files shipped `""`,
and nothing computed or verified it.

## Design
The field is self-referential (the checksum lives inside the file it
describes), so hashing the raw bytes as-is is circular. Fix: the checksum
covers the file with its own `checksum_sha256` **value** zeroed out to a
same-length placeholder (`sha256` hex digests are always exactly 64 chars,
so zeroing never changes the file's byte length or reflows anything else).
This is pure byte-level regex substitution — no JSON re-serialization or
cross-language canonicalization is involved, so a consumer in a different
language (the Go POS host) can verify it identically by doing the same
substitution and re-hashing, with zero risk of a Python-vs-Go JSON
formatting mismatch (key ordering, whitespace, number formatting) producing
false negatives.

## Changes
- `scripts/checksum.py` (new): `--write` recomputes and writes the checksum
  into every `content/<locale>.json`; `--check` recomputes and compares only,
  exits 1 listing any stale files.
- `scripts/validate.sh`: now runs `checksum.py --check` first, so CI catches
  "edited content, forgot to re-run --write" before it ships.
- All 9 `content/<locale>.json`: real checksums populated (verified
  idempotent — running `--write` twice in a row produces no further diff).
- `manifest.json`: 0.2.2 → 0.2.3. `CHANGELOG.md` updated.

## Verification
`scripts/validate.sh` and `scripts/package.sh` (the same two steps CI runs)
both green locally.

## Deliberately out of scope
- **Till-side verification** of the checksum is a separate change in
  `universal-till` (`loadContentBundle`) — tracked and implemented alongside
  this in the same session, see that repo's own review doc.
- The content files' internal `"version"` field still reads `0.2.1` while
  `manifest.json` is now `0.2.3` — pre-existing drift (it already didn't
  match `manifest.json`'s `0.2.2` before this change), tracked separately as
  "FAQ locale-fallback notice + version/last-updated metadata missing" in
  `ut-docs/QUEUE.md`. Not touched here to keep this change to checksum only.
