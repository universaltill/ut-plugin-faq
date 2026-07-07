package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func validManifest() *Manifest {
	return &Manifest{
		ID:            "com.universaltill.ut-faq",
		Name:          "Universal Till FAQ",
		Version:       "0.1.2",
		CanonicalType: "page",
		Executable:    "bin/ut-faq",
		Entrypoint:    "./bin/ut-faq",
		DeviceArch:    "any",
		Permissions:   []string{"storage.local.10MB"},
		Locales:       []string{"en-US"},
		Entries: []ManifestEntry{
			{Type: "page", Key: "faq-page", Label: "Help / FAQ", Route: "/plugin/faq"},
		},
	}
}

func TestValidateManifestHappyPath(t *testing.T) {
	if errs := ValidateManifest(validManifest(), ""); len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestValidateManifestRequiredFields(t *testing.T) {
	m := validManifest()
	m.ID = ""
	m.Executable = ""
	m.Entrypoint = ""
	errs := ValidateManifest(m, "")
	if len(errs) != 3 {
		t.Fatalf("expected 3 errors, got %v", errs)
	}
}

func TestSupportsTarget(t *testing.T) {
	m := validManifest()
	if !m.SupportsTarget("linux", "amd64") {
		t.Fatal("empty supported_architectures should allow any target")
	}
	m.SupportedArch = []string{"linux/amd64", "darwin/arm64"}
	if !m.SupportsTarget("linux", "amd64") {
		t.Fatal("expected listed target to be supported")
	}
	if m.SupportsTarget("linux", "arm") {
		t.Fatal("expected unlisted target to be rejected")
	}
}

func TestValidateManifestRejectsBadSemver(t *testing.T) {
	m := validManifest()
	m.Version = "v1.0"
	errs := ValidateManifest(m, "")
	if len(errs) != 1 || !strings.Contains(errs[0], "semver") {
		t.Fatalf("expected semver error, got %v", errs)
	}
}

func TestValidateManifestRejectsUnknownCanonicalType(t *testing.T) {
	m := validManifest()
	m.CanonicalType = "ui-extension"
	errs := ValidateManifest(m, "")
	if len(errs) != 1 || !strings.Contains(errs[0], "canonical_type") {
		t.Fatalf("expected canonical_type error, got %v", errs)
	}
}

func TestValidateManifestPageNeedsEntries(t *testing.T) {
	m := validManifest()
	m.Entries = nil
	errs := ValidateManifest(m, "")
	if len(errs) != 1 || !strings.Contains(errs[0], "navigation entry") {
		t.Fatalf("expected navigation entry error, got %v", errs)
	}
}

func TestValidateManifestChecksResourceAndLocaleFiles(t *testing.T) {
	root := t.TempDir()
	m := validManifest()
	m.Resources.Icon = "assets/icon.png"

	errs := ValidateManifest(m, root)
	if len(errs) != 2 {
		t.Fatalf("expected missing icon + missing locale content, got %v", errs)
	}

	if err := os.MkdirAll(filepath.Join(root, "assets"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "assets", "icon.png"), []byte("png"), 0o644); err != nil {
		t.Fatal(err)
	}
	contentDir := filepath.Join(root, "src", "faq", "content")
	if err := os.MkdirAll(contentDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(contentDir, "en-US.json"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}

	if errs := ValidateManifest(m, root); len(errs) != 0 {
		t.Fatalf("expected no errors after creating files, got %v", errs)
	}
}

func TestCheckVersionAlignment(t *testing.T) {
	dir := t.TempDir()
	pkg := filepath.Join(dir, "package.json")
	if err := os.WriteFile(pkg, []byte(`{"version":"0.1.2"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	m := validManifest()
	if err := CheckVersionAlignment(m, pkg); err != nil {
		t.Fatalf("expected aligned versions, got %v", err)
	}

	m.Version = "0.9.9"
	if err := CheckVersionAlignment(m, pkg); err == nil {
		t.Fatal("expected mismatch error")
	}
}

func TestRealManifestIsValid(t *testing.T) {
	root := filepath.Join("..", "..")
	m, err := LoadManifest(filepath.Join(root, "src", "manifest", "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	if errs := ValidateManifest(m, root); len(errs) != 0 {
		t.Fatalf("repo manifest invalid: %v", errs)
	}
	if err := CheckVersionAlignment(m, filepath.Join(root, "package.json")); err != nil {
		t.Fatal(err)
	}
}
