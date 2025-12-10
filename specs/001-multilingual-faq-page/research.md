# Research: Multilingual FAQ Page Plugin

## Decisions

- **Language/SDK**: Go 1.21 with Universal Till plugin SDK
  - **Rationale**: SDK supports Go; aligns with POS runtime targets and keeps packaging simple.
  - **Alternatives**: JS/Python SDKs — rejected to avoid bundling interpreters and to keep startup time low.

- **Content format**: Bundled JSON per locale with checksum metadata
  - **Rationale**: Simple to embed, diff, and validate; plays well with checksum verification and offline caching.
  - **Alternatives**: Markdown + render at runtime — adds parsing overhead; binary blobs — harder to diff/verify.

- **Locale handling**: POS locale preference with explicit fallback to en-US; RTL switch for ar-SA/fa-IR
  - **Rationale**: Matches marketplace i18n guidance and constitution; minimizes UX surprises.
  - **Alternatives**: Browser Accept-Language — not applicable; dynamic translation — requires network and risks inconsistency.

- **Search/filter**: In-memory keyword search scoped to locale with simple tokenization
  - **Rationale**: Dataset size is small (~50-200 entries); avoids extra storage/indices; works offline.
  - **Alternatives**: External search service — needs network; SQLite — unnecessary complexity for static content.

- **Permissions**: UI + local storage only (plus optional minimal telemetry if mandated)
  - **Rationale**: Satisfies least-privilege principle; FAQ rendering needs no external access.
  - **Alternatives**: Network permissions — rejected; not required for static FAQ.

- **Update strategy**: Marketplace-driven updates with checksum validation and resume support
  - **Rationale**: Aligns with marketplace download/resume semantics; prevents corrupted cache.
  - **Alternatives**: Custom updater — redundant; manual imports — out of scope for marketplace listing.

- **Navigation placement**: Help/Support area via `plugin_entries` page entry with stable route
  - **Rationale**: Matches POS data model; keeps entry discoverable and removable cleanly.
  - **Alternatives**: Custom menu injection — brittle and may violate UI expectations.

## Open Points

None outstanding; all clarifications resolved within scope.
