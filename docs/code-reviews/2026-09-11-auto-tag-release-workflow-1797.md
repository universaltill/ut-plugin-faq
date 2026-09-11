# Code review: auto-tag-release.yml for ut-plugin-faq (ut-docs#1797)

**Date**: 2026-09-11
**Author**: Farshid Mirza (autonomous SDLC pipeline, `lane:cloud-54`)
**Reviewer**: independent fresh-context Sonnet subagent (read-only pass, no
prior context on the change — per `complexity:easy` model routing)
**Ticket**: universaltill/ut-docs#1797
**Branch**: `fix/1797-auto-tag-release-workflow`

## What changed

Added `.github/workflows/auto-tag-release.yml` to this repo — the workflow
already rolled out to `ut-plugin-tax-de`, `ut-plugin-language-de`,
`ut-plugin-language-es`, and four payment/tax plugin repos (ut-docs#1694,
#1700) that auto-tags a release the moment `manifest.json`'s version has no
matching git tag, then dispatches `release.yml` on the new tag.

This repo was left out of the #1700 mechanical rollout and split into its
own card because it isn't a safe byte-identical copy target: its real
default/integration branch is `001-multilingual-faq-page`
(`origin/HEAD`/`default_branch` both confirm this), not `main` — `main`
here is a stale, disjoint line of history with no `manifest.json` and no
common ancestor. Copying the canonical file unmodified (`push: branches:
[main]`) would silently never fire.

**Decision (recorded on the issue before building): option 1** — retarget
the trigger to `push: branches: [001-multilingual-faq-page]`, documented
inline, rather than option 2 (fix the repo's default branch back to `main`
and reconcile the disjoint history). Option 1 is the lower-risk default —
no history surgery on a live repo — and the issue's own text recommended it
absent a reason to think `main` should be the real default here; none was
found.

## Verification performed

- Fetched the canonical file from `ut-plugin-tax-de`'s `main` and diffed it
  byte-for-byte against this repo's copy: the only differences are the
  retargeted `branches:` trigger, the explanatory header comment, and (after
  the review's nit below) the tag-message branch name. Job logic,
  permissions, the semver guard, the race/retry handling, and the dispatch
  step are all untouched.
- Confirmed `origin/HEAD` → `001-multilingual-faq-page` directly via `git
  branch -a` (GitHub's own default-branch pointer, not an assumption
  inherited from the ticket text) and that this repo's own `ci.yml` already
  lists `branches: [main, 001-multilingual-faq-page]` as its two trigger
  branches — independent corroboration that `001-...` is the real
  integration branch, from a file this change didn't touch.
- Confirmed `release.yml`'s `workflow_dispatch` inputs are exactly `channel`
  (choice, includes `stable`) and `publish` (boolean) — matches the new
  workflow's dispatch call (`-f channel=stable -f publish=true`) with no
  adaptation needed.
- Confirmed `manifest.json`'s `version: "0.2.3"` parses cleanly against the
  workflow's semver guard, and that `v0.2.3` already exists as a tag — so
  the first real run on this branch will correctly no-op (tag exists) rather
  than false-fire, which is itself a clean proof of the tag-check logic
  without needing a live version bump to test against.
- YAML parses (`python3 -c "import yaml; yaml.safe_load(...)"`).
- `scripts/validate.sh` passes on the branch (unrelated to this change, but
  confirms the branch isn't otherwise broken).
- Permissions (`contents: write`, `actions: write`) match exactly what the
  job does — no broader grant than the canonical file already carries. Only
  third-party action is `actions/checkout@v4`, already used identically
  elsewhere in this repo. No secrets logged; uses `github.token` only. No
  self-retrigger risk: this workflow fires on branch push only, never on the
  tag push it performs.

## Review findings

1. **[nit, fixed]** The tag annotation message (`git tag -a ... -m
   "...on merge to main"`) was untouched canonical text and, left as-is,
   would be factually wrong for this repo. Fixed to name
   `001-multilingual-faq-page` instead — cosmetic (shows up in `git tag -n`
   / the GitHub tag view only, no functional effect), but worth getting
   right since the rest of the file is otherwise an exact, deliberate
   deviation from canonical.

## Verdict

**Approved.** The change is exactly the single, well-justified deviation
from the canonical rollout that the ticket describes, independently
corroborated by this repo's own `ci.yml` and `origin/HEAD` rather than
taken on the ticket's word. No blockers.
