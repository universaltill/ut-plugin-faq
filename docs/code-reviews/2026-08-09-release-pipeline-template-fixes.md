# 2026-08-09 — Release-pipeline template fixes (ut-docs#166)

## Context
ut-docs#166 flagged two of three release-pipeline template defects found in
`ut-plugin-payment-sumup`/`ut-plugin-integration-webhook` as also present
here — significant because `ut-docs/for-plugin-developers.md` names this
repo ("`ut-plugin-faq` is the reference implementation — copy its layout")
as the template new plugin authors are told to copy, so fixing it matters
beyond this one repo. The Go-specific `go vet`/`gofmt` CI gate does not
apply — this is a content-only plugin with no Go toolchain.

1. `latest.tar.gz.sha256` — and, found during this fix's own review, the
   versioned artifact's own `.sha256` sidecar too — recorded a
   `dist/`-prefixed path instead of the bare filename, so
   `sha256sum -c` fails for any self-hoster who downloads the artifact +
   checksum pair into one directory.
2. `scripts/approve.sh` reused `MARKETPLACE_UPLOAD_TOKEN` (a vendor-upload
   credential) as the bearer token for the marketplace's admin
   review/approve endpoints. Fails closed today, not currently exploitable.

## Changes
- **`scripts/package.sh`**: the `sha256sum "$OUT"` line now `cd`s into
  `dist/` and hashes the bare filename, fixing the versioned artifact's
  sidecar (found during review — the original ticket only named the
  `latest.tar.gz` copy).
- **`.github/workflows/release.yml`**: the "Create GitHub Release" step
  regenerates `latest.tar.gz.sha256` against the renamed
  `dist/latest.tar.gz` directly instead of copying the versioned `.sha256`
  file verbatim.
- **`scripts/approve.sh`**: now prefers `MARKETPLACE_ADMIN_TOKEN`, falling
  back to `MARKETPLACE_UPLOAD_TOKEN` (unchanged current behavior) when
  unset. **Client-side prep only** — ut-cloud's `authorizeStaff` accepts
  only the upload-token value today, so a genuinely distinct
  `MARKETPLACE_ADMIN_TOKEN` would 401 every call. Filed
  **universaltill/ut-docs#496** for the real fix; comments in both files
  say explicitly not to provision a distinct secret until #496 ships.

## Independent review
Fresh-context Opus review (different model from the implementer), run once
across all three repos together (identical diffs where applicable). No
blockers; the should-fix items relevant to this repo (checksum bug also in
the versioned sidecar; the admin-token warning wording) were applied before
commit — see `ut-plugin-payment-sumup`'s twin record for the full findings
detail and verification methodology.

## Verification
- `scripts/validate.sh && scripts/package.sh` — green (9 locale bundles,
  checksums current).
- Reproduced the self-hoster checksum failure pre-fix and confirmed both
  the versioned sidecar and the `latest.tar.gz` pair pass `sha256sum -c`
  post-fix.
- Token-selection logic verified in isolation across all 4
  admin/upload-token combinations.
- `bash -n` and YAML parse on every edited file.

## Deliberately out of scope
- ut-cloud-side distinct admin credential — universaltill/ut-docs#496.
- Sweeping the other `ut-plugin-*` repos — universaltill/ut-docs#497.
- `README.md`'s secrets table still lists only `MARKETPLACE_UPLOAD_TOKEN` —
  left as-is deliberately: documenting `MARKETPLACE_ADMIN_TOKEN` as
  something an operator should set would contradict the "don't provision it
  yet" guidance in the code comments until #496 actually ships.
