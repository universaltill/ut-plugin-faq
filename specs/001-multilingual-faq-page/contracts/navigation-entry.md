# Contract: Plugin Navigation Entry (FAQ Page)

## Fields (aligns with POS plugin_entries)
- `type`: `page`
- `key`: unique key for FAQ page
- `label`: localized label per locale
- `icon_path`: path to icon asset packaged with plugin
- `route`: stable route handled by plugin
- `parent_page_key`: Help/Support menu key
- `sort_order`: integer ordering within Help/Support
- `is_active`: enable/disable flag
- `menu_group`: optional grouping under Help/Support
- `config_json`: optional for layout metadata (theme, categories ordering)

## Behavioral Expectations
- Entry is registered on install/enable; removed/hidden on disable/uninstall.
- Uses POS locale to select label/icon; falls back to en-US if missing.
- Route resolves offline using cached/bundled content; no network dependency.
