# Data Model: Multilingual FAQ Page Plugin

## Entities

### FAQEntry
- **Fields**: id, locale, category, question, answer, last_updated, sort_order, keywords
- **Rules**: locale must be in supported set; sort_order stable per category; questions/answers localized fully (no mixed language fragments).

### FAQCategory
- **Fields**: id, locale, name, sort_order
- **Rules**: categories localized; categories referenced by FAQEntry must exist for the same locale.

### LocalizationBundle
- **Fields**: locale, checksum_sha256, version, faq_entries[], categories[], rtl (bool), fallback_locale
- **Rules**: checksum validated on install/update; fallback_locale defaults to en-US for unsupported locales; rtl true for ar-SA and fa-IR.

### PluginNavigationEntry
- **Fields**: key, label, icon_path, route, parent_page_key, sort_order, is_active
- **Rules**: type `page`; parent_page_key targets Help/Support; localized label/icon per locale; disabled on uninstall.

## Relationships
- LocalizationBundle 1..* FAQEntry (by locale)
- LocalizationBundle 1..* FAQCategory (by locale)
- PluginNavigationEntry references LocalizationBundle for label/icon localization.

## State/Transitions
- **Installed**: bundles and manifest verified; navigation entry registered and active.
- **Updated**: new bundles downloaded; checksum validated; swap is atomic; on failure, rollback to previous bundle.
- **Disabled/Uninstalled**: navigation entry removed/hidden; cached bundles retained or purged per POS policy.
