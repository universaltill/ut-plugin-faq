<!--
Sync Impact Report
- Version change: N/A → 1.0.0
- Modified principles: Initial definition (User-Value & I18N First; Marketplace Compliance by Default; Offline Resilience; Testable, Independent Increments; Security & Least Privilege)
- Added sections: Operational & Compliance Standards; Development Workflow & Quality Gates
- Removed sections: None
- Templates requiring updates: ✅ .specify/templates/spec-template.md (reviewed, no change) | ✅ .specify/templates/plan-template.md (reviewed, no change) | ✅ .specify/templates/tasks-template.md (reviewed, no change)
- Follow-up TODOs: None
-->

# UT Plugin FAQ Constitution

## Core Principles

### I. User-Value & I18N First
All work must prioritize merchant and cashier comprehension: multilingual content, RTL correctness, clear labels, and predictable navigation. Success is measured by reduced support needs and faster task completion across locales.

### II. Marketplace Compliance by Default
Every change must preserve marketplace readiness: valid manifest and taxonomy, declared capabilities/types, minimal permissions, required assets (icon, screenshots, docs), and zero MV-* validation errors before release.

### III. Offline Resilience
FAQ content and UI must function without connectivity after installation. Downloads and updates must resume safely, validate checksums, and never corrupt existing cached content.

### IV. Testable, Independent Increments
User stories and tasks must be independently deliverable and verifiable. Each increment needs acceptance scenarios, measurable outcomes, and clear scope boundaries to avoid coupling.

### V. Security & Least Privilege
Plugins request only the permissions necessary for UI display and local storage. Authentication context from POS must be respected; sensitive data avoidance is preferred over protection. Any telemetry must follow marketplace policy.

## Operational & Compliance Standards

- Support the published locale set (en-US, en-GB, fr-FR, ar-SA, fa-IR, tr-TR, es-ES, it-IT, pt-PT) with explicit fallbacks to English and correct RTL handling where applicable.
- Declare `canonical_type` suitable for UI/page plugins and align `plugin_entries` to Help/Support navigation with stable routes.
- Package FAQ content with version and checksum metadata; reject or retry corrupted downloads.
- Keep permissions constrained to UI rendering, cached storage, and only mandated telemetry; no external network calls for FAQ rendering.
- Maintain merchant-facing documentation and visible version/last-updated metadata inside the FAQ page.

## Development Workflow & Quality Gates

- Specifications must include prioritized, independently testable user stories, edge cases, functional requirements, entities, and measurable success criteria.
- Plans must document constraints for offline-first behavior, locale coverage, manifest validation, and permission minimization.
- Tasks must map to user stories and keep dependencies minimal; testing tasks are encouraged when explicitly requested by specs.
- Before delivery: ensure manifest validation passes, offline rendering verified, locale coverage confirmed (including RTL), and navigation entry registered/removed cleanly.
- Any complexity or deviation from principles requires written justification in plan.md (Complexity Tracking) and reviewer acknowledgment.

## Governance

- This constitution governs all specs, plans, and tasks for the UT Plugin FAQ project. Conflicts resolve in favor of this document.
- Amendments: propose changes via PR with rationale, version bump per semantic rules (MAJOR for principle changes/removals, MINOR for new/expanded guidance, PATCH for clarifications).
- Compliance: reviewers must check Constitution Check sections in plans; violations require documented justification.
- Source of truth: `.specify/memory/constitution.md`; templates must be updated or explicitly noted when guidance changes.

**Version**: 1.0.0 | **Ratified**: 2025-12-10 | **Last Amended**: 2025-12-10
