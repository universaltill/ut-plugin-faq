#!/usr/bin/env bash
# Validates the FAQ page-plugin manifest: marketplace-required fields
# (id/name/semver/permissions/locales), asset-only runtime (ADR-0001),
# exactly one type="page" entry with a route, and a content bundle for
# every declared locale.
set -euo pipefail
cd "$(dirname "$0")/.."

echo "Checking content bundle checksums..."
python3 scripts/checksum.py --check

python3 - <<'PY'
import json, os, re, sys
m = json.load(open("manifest.json"))
errs = []
if not re.match(r'^[a-z0-9]+([.-][a-z0-9]+)*$', m.get("id","")): errs.append("bad id")
if not m.get("name"): errs.append("missing name")
if not re.match(r'^\d+\.\d+\.\d+', m.get("version","")): errs.append("bad version")
if not m.get("permissions"): errs.append("missing permissions")
if not m.get("locales"): errs.append("missing locales")
if m.get("runtime") != "none": errs.append("runtime must be 'none' — the till renders the content bundle (ADR-0001)")
if m.get("device_arch") != "any": errs.append("device_arch must be 'any'")
if m.get("canonical_type") != "page": errs.append("canonical_type must be 'page'")
pages = [e for e in m.get("entries", []) if e.get("type") == "page"]
if len(pages) != 1:
    errs.append(f"expected exactly 1 page entry, got {len(pages)}")
else:
    if not pages[0].get("key"): errs.append("page entry missing key")
    if not pages[0].get("label"): errs.append("page entry missing label")
    if not pages[0].get("route"): errs.append("page entry missing route")
for loc in m.get("locales", []):
    if not os.path.isfile(f"content/{loc}.json"):
        errs.append(f"missing content bundle content/{loc}.json")
    else:
        try:
            json.load(open(f"content/{loc}.json"))
        except Exception as e:
            errs.append(f"content/{loc}.json invalid JSON: {e}")
if errs:
    print("FAIL: " + "; ".join(errs)); sys.exit(1)
print(f"ok {m['id']} v{m['version']} ({len(m['locales'])} locales)")
PY
