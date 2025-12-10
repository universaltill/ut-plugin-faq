package ui

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"ut-plugin-faq/src/faq"
	"ut-plugin-faq/src/storage"
)

// Simple auth context to gate FAQ page access.
type AuthContext struct {
	Authenticated bool
	Role          string
}

type NavigationEntry struct {
	Type          string
	Key           string
	Label         string
	IconPath      string
	Route         string
	ParentPageKey string
	SortOrder     int
	MenuGroup     string
	ConfigJSON    string
	IsActive      bool
}

// FAQPage encapsulates navigation and rendering behavior.
type FAQPage struct {
	Renderer         PageRenderer
	Navigation       NavigationEntry
	DefaultLocale    string
	SupportedLocales []string
	VersionLabel     string
}

// NewFAQNavigationEntry constructs a navigation entry aligned with plugin_entries contract.
func NewFAQNavigationEntry(key, label, iconPath, route, parent string, sortOrder int) NavigationEntry {
	return NavigationEntry{
		Type:          "page",
		Key:           key,
		Label:         label,
		IconPath:      iconPath,
		Route:         route,
		ParentPageKey: parent,
		SortOrder:     sortOrder,
		IsActive:      true,
	}
}

// Disable marks the entry inactive for uninstall/disable flows.
func (n *NavigationEntry) Disable() {
	n.IsActive = false
}

// RequireAuth ensures only authenticated staff can access the FAQ page.
func RequireAuth(ctx AuthContext) error {
	if !ctx.Authenticated {
		return errors.New("authentication required")
	}
	return nil
}

// RegisterFAQPage constructs the FAQ page with navigation metadata.
func RegisterFAQPage(renderer PageRenderer) FAQPage {
	entry := NewFAQNavigationEntry(
		"faq-page",
		"Help / FAQ",
		"assets/icon.png",
		"/plugin/faq",
		"help_support",
		10,
	)
	return FAQPage{
		Renderer:         renderer,
		Navigation:       entry,
		DefaultLocale:    "en-US",
		SupportedLocales: []string{"en-US", "en-GB", "fr-FR", "ar-SA", "fa-IR", "tr-TR", "es-ES", "it-IT", "pt-PT"},
		VersionLabel:     "0.1.0",
	}
}

// Unregister marks the navigation entry inactive (disable/uninstall).
func (p *FAQPage) Unregister() {
	p.Navigation.Disable()
}

// normalizeLocale enforces supported locale or falls back.
func (p FAQPage) normalizeLocale(locale string) string {
	if locale == "" {
		return p.DefaultLocale
	}
	for _, loc := range p.SupportedLocales {
		if strings.EqualFold(loc, locale) {
			return loc
		}
	}
	return p.DefaultLocale
}

// Handle returns the localized bundle after auth + locale resolution.
func (p FAQPage) Handle(ctx AuthContext, locale string) (*faq.Bundle, error) {
	if err := RequireAuth(ctx); err != nil {
		return nil, err
	}
	target := p.normalizeLocale(locale)
	bundle, err := p.Renderer.LoadBundle(target)
	if err != nil {
		if target != p.DefaultLocale {
			return p.Renderer.LoadBundle(p.DefaultLocale)
		}
		return nil, err
	}
	return bundle, nil
}

// UpdateBundles applies refreshed content with checksum validation and returns the outcome per locale.
func (p FAQPage) UpdateBundles(locale string, data []byte, expectedChecksum string) storage.UpdateResult {
	return p.Renderer.Cache.ApplyUpdate(locale, data, expectedChecksum)
}

// LocalizedLabel returns a locale-aware navigation label.
func (p FAQPage) LocalizedLabel(locale string) string {
	loc := p.normalizeLocale(locale)
	switch loc {
	case "fr-FR":
		return "Aide / FAQ"
	case "ar-SA", "fa-IR":
		return "مساعدة / الأسئلة"
	case "es-ES":
		return "Ayuda / FAQ"
	case "it-IT":
		return "Aiuto / FAQ"
	case "pt-PT":
		return "Ajuda / FAQ"
	case "tr-TR":
		return "Yardım / SSS"
	default:
		return "Help / FAQ"
	}
}

// PageRenderer loads localized FAQ content for the POS route.
type PageRenderer struct {
	Cache storage.Cache
}

// LoadBundle chooses cached bundle first, falling back to embedded content.
func (r PageRenderer) LoadBundle(locale string) (*faq.Bundle, error) {
	// Try cache
	if data, checksum, err := r.Cache.Load(locale); err == nil {
		if bundle, err := faq.ParseBundle(data, checksum); err == nil {
			return bundle, nil
		}
	}
	// Fallback to embedded
	bundle, err := faq.LoadEmbeddedBundle(locale, "")
	if err != nil {
		return nil, fmt.Errorf("load bundle: %w", err)
	}
	// Cache embedded for offline use
	if b, err := json.Marshal(bundle); err == nil {
		_ = r.Cache.Save(locale, b, bundle.ChecksumSHA256)
	}
	return bundle, nil
}
