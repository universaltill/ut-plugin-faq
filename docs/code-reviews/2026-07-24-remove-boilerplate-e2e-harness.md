# 2026-07-24 — Remove the boilerplate `tests/e2e/` harness

## Context
Spec-audit backlog item (`ut-docs/QUEUE.md`): "FAQ e2e tests are boilerplate,
not real coverage". `tests/e2e/example.spec.ts` was a generic-template
scaffold that never matched this plugin:

- Navigated `/faq`; the plugin's actual registered route is `/plugin/faq`
  (`manifest.json` entries[0].route).
- `faq-factory.ts` POSTed/DELETEd against `${API_URL}/faq`, a REST CRUD API
  that doesn't exist — this plugin is `runtime: "none"` (ADR-0001), asset-only,
  with no server of its own. The till renders `content/<locale>.json`
  natively; there is nothing here to POST to.
- Never wired into CI (`.github/workflows/ci.yml` runs `scripts/validate.sh`
  and `scripts/package.sh`, never `npm run test:e2e`) — the suite's breakage
  was invisible.

## Decision
Real browser-level coverage of the render path (locale fallback, RTL, the
client-side search JS) can only be exercised where the renderer lives —
`universal-till`'s `internal/pages/plugin_page.go` /
`web/ui/pages/plugin_content.html`. This repo ships content, not a page
renderer, so a Playwright suite here was structurally the wrong place for
that coverage regardless of whether it was pointed at the right route.

Added instead, in `universal-till` (see that repo's review doc,
`2026-07-24-faq-real-e2e-coverage.md`): a real e2e suite that installs this
plugin's actual `content/en-US.json` and `content/fa-IR.json` bundles into a
real booted till and drives Chromium against `/plugin/faq`, plus a Go unit
test closing the one gap that had literally zero coverage (every existing
fixture used `"rtl":false`; RTL rendering itself was never asserted).

## Changes (this repo)
- Removed `tests/` (the whole boilerplate suite: `e2e/example.spec.ts`,
  `support/fixtures/index.ts`, `support/fixtures/factories/faq-factory.ts`,
  `README.md`), `playwright.config.ts`, `.env.example`.
- `package.json`: dropped `test:e2e` script and the now-unused
  `@playwright/test` / `@faker-js/faker` / `typescript` devDependencies (no
  `.ts` files remain in the repo).

## Verification
`scripts/validate.sh` still green (manifest + all 9 content bundle
checksums) — this repo's real, CI-wired validation was untouched.

## Deliberately out of scope
- `package.json`'s own `version: 0.1.2` field is stale relative to
  `manifest.json`'s `0.2.3` — pre-existing, unrelated to the e2e harness,
  not touched here.
- `.nvmrc` left in place even though no `.ts`/node tooling remains in the
  repo post-cleanup — low-cost to keep for whatever npm-based tooling might
  land here later; removing it wasn't part of this task.
