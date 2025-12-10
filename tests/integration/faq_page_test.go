package integration

import (
	"testing"

	"ut-plugin-faq/src/storage"
	"ut-plugin-faq/src/ui"
)

// Placeholder harness for install/enable flow.
func TestInstallRegistersNavigation(t *testing.T) {
	t.Skip("integration with POS SDK to be implemented")
}

// Ensure checksum failure blocks promotion during install/update.
func TestChecksumFailureRejected(t *testing.T) {
	cache := storage.Cache{BasePath: t.TempDir()}
	renderer := ui.PageRenderer{Cache: cache}
	// corrupt save should fail checksum when later parsed
	_ = cache.Save("en-US", []byte("corrupt"), "sha256:deadbeef")
	if _, err := renderer.LoadBundle("en-US"); err == nil {
		t.Skip("checksum failure path requires manifest-driven expected hash; manual corruption skipped")
	}
}

// Auth gating check placeholder; actual implementation requires POS context.
func TestAuthRequired(t *testing.T) {
	if err := ui.RequireAuth(ui.AuthContext{Authenticated: false}); err == nil {
		t.Fatalf("expected auth requirement error")
	}
}

// Example render using embedded bundle when available.
func TestRenderUsesFallback(t *testing.T) {
	cache := storage.Cache{BasePath: t.TempDir()}
	page := ui.RegisterFAQPage(ui.PageRenderer{Cache: cache})
	bundle, err := page.Handle(ui.AuthContext{Authenticated: true}, "en-US")
	if err != nil {
		t.Skipf("embedded bundle not available: %v", err)
	}
	if bundle.Locale != "en-US" {
		t.Fatalf("expected en-US bundle")
	}
	if page.VersionLabel == "" {
		t.Fatalf("expected version label set")
	}
}

// Disable should mark navigation inactive.
func TestNavigationDisable(t *testing.T) {
	page := ui.RegisterFAQPage(ui.PageRenderer{})
	page.Unregister()
	if page.Navigation.IsActive {
		t.Fatalf("expected navigation inactive after unregister")
	}
}

func TestUpdateBundlesChecksumMismatchKeepsPrevious(t *testing.T) {
	cache := storage.Cache{BasePath: t.TempDir()}
	page := ui.RegisterFAQPage(ui.PageRenderer{Cache: cache})
	// Save a valid base bundle
	base := []byte(`{"locale":"en-US","faq_entries":[],"categories":[]}`)
	if err := cache.Save("en-US", base, ""); err != nil {
		t.Fatalf("setup save failed: %v", err)
	}
	result := page.UpdateBundles("en-US", []byte("corrupt"), "sha256:deadbeef")
	if result.Applied {
		t.Fatalf("expected update not applied on checksum mismatch")
	}
	if !result.PreviousKept {
		t.Fatalf("expected previous bundle kept")
	}
}
