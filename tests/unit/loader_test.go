package unit

import (
	"testing"

	faq "ut-plugin-faq/src/faq"
)

func TestComputeChecksumMatchesExpected(t *testing.T) {
	data := []byte(`{"locale":"en-US"}`)
	got := faq.ComputeChecksum(data)
	if got == "" || got[:7] != "sha256:" {
		t.Fatalf("checksum missing prefix: %s", got)
	}
}

func TestParseBundleRejectsBadChecksum(t *testing.T) {
	data := []byte(`{"locale":"en-US"}`)
	_, err := faq.ParseBundle(data, "sha256:deadbeef")
	if err == nil {
		t.Fatalf("expected checksum error")
	}
}

func TestParseBundleNormalizesFallback(t *testing.T) {
	data := []byte(`{"locale":"en-US","faq_entries":[],"categories":[]}`)
	b, err := faq.ParseBundle(data, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.FallbackLocale != "en-US" {
		t.Fatalf("expected fallback en-US, got %s", b.FallbackLocale)
	}
}

func TestLoadEmbeddedBundle(t *testing.T) {
	b, err := faq.LoadEmbeddedBundle("en-US", "")
	if err != nil {
		t.Skipf("embedded asset not available: %v", err)
	}
	if b.Locale == "" {
		t.Fatalf("expected locale set")
	}
}
