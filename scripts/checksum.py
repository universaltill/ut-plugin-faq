#!/usr/bin/env python3
"""Computes/verifies each content/<locale>.json's checksum_sha256 field.

The field is self-referential (the file's own checksum lives inside the
file), so it can't hash the file bytes as-is. Instead: the checksum covers
the file with its own checksum_sha256 VALUE zeroed out to a same-length
placeholder (sha256 hex digests are always exactly 64 chars, so zeroing
never changes the file's byte length or reflows anything). This is pure
byte-level substitution — no JSON re-serialization/canonicalization is
involved, so a consumer in any language (the Go POS host, here) verifies it
identically by doing the same substitution and re-hashing.

Modes:
  --write  recompute and write the checksum into every content/<locale>.json
  --check  recompute and compare only; exit 1 (listing the stale files) if
           any file's stored checksum doesn't match its actual content —
           catches "edited content, forgot to re-run --write".
"""
import hashlib
import re
import sys
from pathlib import Path

CONTENT_DIR = Path(__file__).resolve().parent.parent / "content"
FIELD_PATTERN = re.compile(r'("checksum_sha256":\s*")([0-9a-f]{0,64})(")')
PLACEHOLDER = "0" * 64


def compute(raw: str) -> str:
    zeroed, n = FIELD_PATTERN.subn(r"\g<1>" + PLACEHOLDER + r"\g<3>", raw, count=1)
    if n != 1:
        raise ValueError("expected exactly one checksum_sha256 field")
    return hashlib.sha256(zeroed.encode("utf-8")).hexdigest()


def main() -> int:
    mode = sys.argv[1] if len(sys.argv) > 1 else "--check"
    if mode not in ("--write", "--check"):
        print(f"usage: {sys.argv[0]} [--write|--check]", file=sys.stderr)
        return 2

    files = sorted(CONTENT_DIR.glob("*.json"))
    if not files:
        print(f"no content files found under {CONTENT_DIR}", file=sys.stderr)
        return 2

    stale = []
    for path in files:
        raw = path.read_text(encoding="utf-8")
        digest = compute(raw)
        m = FIELD_PATTERN.search(raw)
        current = m.group(2) if m else ""
        if current == digest:
            continue
        if mode == "--write":
            updated, n = FIELD_PATTERN.subn(r"\g<1>" + digest + r"\g<3>", raw, count=1)
            assert n == 1
            path.write_text(updated, encoding="utf-8")
            print(f"wrote {path.name}: {digest}")
        else:
            stale.append(path.name)

    if mode == "--check" and stale:
        print("FAIL: checksum_sha256 out of date (run scripts/checksum.py --write): " + ", ".join(stale), file=sys.stderr)
        return 1

    if mode == "--check":
        print(f"ok — checksum_sha256 current for all {len(files)} locale bundles")
    return 0


if __name__ == "__main__":
    sys.exit(main())
