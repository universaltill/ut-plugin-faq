package main

import (
	"log"
	"os"
	"path/filepath"

	"ut-plugin-faq/src/storage"
	"ut-plugin-faq/src/ui"
)

func registerNavigation() ui.NavigationEntry {
	return ui.NewFAQNavigationEntry(
		"faq-page",
		"Help / FAQ",
		"assets/icon.png",
		"/plugin/faq",
		"help_support",
		10,
	)
}

// placeholder main until wired to the POS plugin SDK
func main() {
	cache := storage.Cache{BasePath: filepath.Join("data", "cache")}
	renderer := ui.PageRenderer{Cache: cache}
	bundle, err := renderer.LoadBundle("en-US")
	if err != nil {
		log.Printf("warning: failed to load bundle: %v", err)
	}
	if bundle != nil {
		log.Printf("loaded FAQ bundle locale=%s version=%s rtl=%v", bundle.Locale, bundle.Version, bundle.RTL)
	}

	entry := registerNavigation()
	if entry.Label == "" {
		os.Exit(1)
	}
	log.Printf("navigation registered: %s -> %s", entry.Label, entry.Route)
}
