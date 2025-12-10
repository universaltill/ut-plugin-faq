---

description: "Task list for feature implementation"
---

# Tasks: Multilingual FAQ Page Plugin

**Input**: Design documents from `/specs/001-multilingual-faq-page/`  
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Include targeted unit/integration checks where they lock in acceptance scenarios.  
**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: User story label (US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [X] T001 Create plugin skeleton directories per plan: `src/`, `src/manifest/`, `src/faq/content/`, `src/faq/`, `src/ui/`, `src/storage/`, `tests/unit/`, `tests/integration/`
- [X] T002 Initialize Go module and add dependencies in `go.mod` (Universal Till plugin SDK, JSON handling/encoding, checksum support)
- [X] T003 Configure formatting/lint helpers (`.golangci.yml` or simple `scripts/dev.sh` with `gofmt`, `go vet`) and update `README`/`quickstart.md` with commands

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure required before user stories

- [X] T004 Create manifest scaffold in `src/manifest/manifest.json` with required fields (id, name, version, canonical_type `ui-extension`, capability `page`, locales, architectures, permissions minimal)
- [X] T005 Implement selected localization approach using embedded JSON bundles; ensure plan.md reflects this decision
- [X] T006 Add localization bundle schema and sample content placeholders for each locale in `src/faq/content/*.json`, including version metadata and checksum fields
- [X] T007 Implement bundle loader and checksum validator in `src/faq/loader.go` to read embedded/cached JSON and verify integrity; reject/promote-block on checksum failure
- [X] T008 Implement storage cache utility in `src/storage/cache.go` using POS storage API for offline FAQ bundles (read/write, version swap)
- [X] T009 Define navigation entry contract and helpers in `src/ui/page.go` aligning to `plugin_entries` (type `page`, parent Help/Support, stable route)
- [X] T010 [P] Add unit tests scaffolding for loader/cache in `tests/unit/loader_test.go` and `tests/unit/cache_test.go`

**Checkpoint**: Loader, cache, manifest scaffold, and navigation helpers exist; unit test scaffolds ready.

---

## Phase 3: User Story 1 - Install and surface FAQ page (Priority: P1) 🎯 MVP

**Goal**: Install plugin, register Help/FAQ entry, render page offline from bundled content.

**Independent Test**: Install plugin; Help/FAQ entry appears with localized label/icon; offline load renders FAQ page from bundled content without network.

### Tests for User Story 1

- [X] T011 [P] [US1] Add integration test harness to simulate install/enable and verify navigation entry + route in `tests/integration/faq_page_test.go`
- [X] T012 [P] [US1] Add checksum failure test (corrupted bundle) on install/update in `tests/integration/faq_page_test.go`
- [X] T013 [P] [US1] Add auth/permission test to verify FAQ page requires signed-in staff and manifest permissions remain minimal in `tests/integration/faq_page_test.go`

### Implementation for User Story 1

- [X] T014 [US1] Complete manifest fields (resources, permissions minimal, locales, architectures) in `src/manifest/manifest.json`
- [X] T015 [US1] Implement POS auth context gating for FAQ route (signed-in staff) and align permissions scope in `src/ui/page.go` and `src/main.go`
- [X] T016 [US1] Embed default locale bundles with version metadata and register them with loader in `src/faq/loader.go`
- [X] T017 [US1] Implement navigation entry registration/removal with disable/uninstall cleanup in `src/ui/page.go` wiring route + `parent_page_key`
- [X] T018 [US1] Implement initial page render handler in `src/ui/page.go` that pulls bundled content via loader/cache and serves offline
- [X] T019 [US1] Add install/offline integration flow in `src/main.go` to wire plugin SDK lifecycle to page registration and cache init
- [X] T020 [US1] Update assets (icon 512x512, sample screenshots) paths referenced by manifest in `src/manifest/`
- [X] T021 [US1] Document install/offline, checksum validation, and auth/permission checks in `specs/001-multilingual-faq-page/quickstart.md` (append)

**Checkpoint**: FAQ page surfaces in Help/Support and works offline with bundled content.

---

## Phase 4: User Story 2 - View FAQs in preferred language (Priority: P2)

**Goal**: Render FAQ content in POS locale with fallback and RTL correctness; enable search/filter per locale.

**Independent Test**: Switch POS locale across supported set; FAQ page shows full localization/RTL. Unsupported locale falls back to en-US with notice. Search/filter returns locale-scoped results.

### Tests for User Story 2

- [ ] T022 [P] [US2] Add locale render tests (LTR + RTL + fallback) in `tests/integration/faq_page_test.go`
- [X] T023 [P] [US2] Add unit tests for search/filter per locale with relevance assertions (>95% expected hits for test data) in `tests/unit/render_test.go`

### Implementation for User Story 2

- [X] T024 [US2] Implement locale selection with fallback notice and RTL switch in `src/faq/loader.go` / `src/ui/page.go`
- [X] T025 [US2] Build rendering logic with category browse + keyword search scoped to locale in `src/faq/render.go`
- [X] T026 [US2] Localize navigation labels/icons using bundle data in `src/ui/page.go`
- [X] T027 [US2] Ensure unsupported locale path renders en-US content without layout breakage in `src/ui/page.go`
- [X] T028 [US2] Update locale bundles with translated UI strings and FAQ entries in `src/faq/content/*.json`

**Checkpoint**: Locale-aware rendering and search/filter verified across supported locales with correct RTL/fallback handling.

---

## Phase 5: User Story 3 - Keep FAQs current via marketplace updates (Priority: P3)

**Goal**: Deliver updated FAQ content via marketplace releases with checksum validation, resume support, and rollback safety.

**Independent Test**: Publish updated bundle; POS installs update and shows new content/version. Interrupted download resumes and preserves previous bundle until validation passes.

### Tests for User Story 3

- [X] T029 [P] [US3] Add update/rollback test to `tests/integration/faq_page_test.go` simulating interrupted download and checksum failure

### Implementation for User Story 3

- [X] T030 [US3] Implement content versioning metadata display in UI (`last_updated`/version badge) in `src/ui/page.go`
- [X] T031 [US3] Implement resumable download/apply flow using POS storage APIs in `src/storage/cache.go` with checksum gating
- [X] T032 [US3] Wire update hook to refresh bundles and navigation assets after successful validation in `src/main.go`
- [X] T033 [US3] Update manifest/changelog in `src/manifest/manifest.json` and `CHANGELOG.md` to reflect new content releases

**Checkpoint**: Updates apply safely with validation, resume, and visible version metadata.

---

## Phase N: Polish & Cross-Cutting Concerns

- [ ] T034 [P] Validate manifest with marketplace tooling (e.g., `uitill manifest validate src/manifest/manifest.json`) and fix MV-* issues (rerun after version bumps)
- [ ] T035 [P] Performance check: measure load times vs targets; optimize content parsing/render if needed (document in `quickstart.md`)
- [ ] T036 [P] Security/permissions audit to confirm least-privilege set in `src/manifest/manifest.json` and auth gating behavior
- [ ] T037 [P] Final documentation pass (README/FAQ content authoring guide) in `docs/` or `specs/001-multilingual-faq-page/`

---

## Dependencies & Execution Order

- Phase 1 Setup → Phase 2 Foundational → User Stories (US1 → US2 → US3) → Polish
- US2 depends on US1 foundational navigation/render scaffolding; US3 depends on cache/loader from US1/US2.

## Parallel Opportunities

- Setup tasks T001–T003 can run in parallel after repo checkout.
- Foundational T004–T008 mostly sequential; T010 test scaffolding can run in parallel.
- Within US1: T011–T016 largely sequential; tests T011–T012 can begin after foundational cache/loader/nav helpers exist.
- Within US2: T021/T022 tests in parallel with T023–T027 once US1 page render exists.
- Within US3: T028 test can parallel T029–T032 once cache/update hooks exist.

## MVP Scope

- Complete through User Story 1 (T001–T020) to deliver installable, offline FAQ page with navigation entry, auth gating, and bundled content.
