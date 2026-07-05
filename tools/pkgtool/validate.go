package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Manifest mirrors the fields pkgtool needs from src/manifest/manifest.json.
// It intentionally covers both the marketplace contract (capabilities, locales,
// resources) and the POS host install contract (canonical_type, executable,
// device_arch, entries).
type Manifest struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Version       string            `json:"version"`
	CanonicalType string            `json:"canonical_type"`
	Executable    string            `json:"executable"`
	Entrypoint    string            `json:"entrypoint"`
	DeviceArch    string            `json:"device_arch"`
	SupportedArch []string          `json:"supported_architectures"`
	Permissions   []string          `json:"permissions"`
	Locales       []string          `json:"locales"`
	Entries       []ManifestEntry   `json:"entries"`
	Resources     ManifestResources `json:"resources"`
}

type ManifestEntry struct {
	Type  string `json:"type"`
	Key   string `json:"key"`
	Label string `json:"label"`
	Route string `json:"route"`
}

type ManifestResources struct {
	Icon          string   `json:"icon"`
	Screenshots   []string `json:"screenshots"`
	Documentation string   `json:"documentation"`
	License       string   `json:"license"`
	Changelog     string   `json:"changelog"`
}

var semverRe = regexp.MustCompile(`^\d+\.\d+\.\d+(-[0-9A-Za-z.-]+)?(\+[0-9A-Za-z.-]+)?$`)

// validCanonicalTypes mirrors the POS host taxonomy in
// universal-till/internal/plugins/manifest_verifier.go (isValidCanonicalType).
// Keep the two lists in sync.
var validCanonicalTypes = map[string]bool{
	"page":           true,
	"button":         true,
	"payment":        true,
	"report":         true,
	"integration":    true,
	"background_job": true,
	"device":         true,
}

// LoadManifest reads and decodes a manifest file.
func LoadManifest(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read manifest: %w", err)
	}
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parse manifest: %w", err)
	}
	return &m, nil
}

// ValidateManifest checks required fields and value formats. root is the repo
// root used to resolve resource paths; pass "" to skip file-existence checks.
func ValidateManifest(m *Manifest, root string) []string {
	var errs []string

	if m.ID == "" {
		errs = append(errs, "id is required")
	}
	if m.Name == "" {
		errs = append(errs, "name is required")
	}
	if m.Version == "" {
		errs = append(errs, "version is required")
	} else if !semverRe.MatchString(m.Version) {
		errs = append(errs, fmt.Sprintf("version %q is not valid semver", m.Version))
	}
	if !validCanonicalTypes[m.CanonicalType] {
		errs = append(errs, fmt.Sprintf("canonical_type %q is not in the POS taxonomy", m.CanonicalType))
	}
	if m.Executable == "" {
		errs = append(errs, "executable is required")
	}
	if m.Entrypoint == "" {
		errs = append(errs, "entrypoint is required (POS host ParseManifest rejects manifests without it)")
	}
	if m.DeviceArch == "" {
		errs = append(errs, "device_arch is required")
	}
	if len(m.Permissions) == 0 {
		errs = append(errs, "permissions must declare at least one scope")
	}
	if len(m.Locales) == 0 {
		errs = append(errs, "locales must list at least one locale")
	}
	if m.CanonicalType == "page" && len(m.Entries) == 0 {
		errs = append(errs, "page plugins must declare at least one navigation entry")
	}
	for i, e := range m.Entries {
		if e.Type == "" || e.Key == "" || e.Label == "" {
			errs = append(errs, fmt.Sprintf("entries[%d] must set type, key, and label", i))
		}
		if e.Type == "page" && e.Route == "" {
			errs = append(errs, fmt.Sprintf("entries[%d] page entry must set route", i))
		}
	}

	if root != "" {
		for _, rel := range resourcePaths(m) {
			if _, err := os.Stat(filepath.Join(root, rel)); err != nil {
				errs = append(errs, fmt.Sprintf("resource %q not found in repo", rel))
			}
		}
		for _, locale := range m.Locales {
			content := filepath.Join(root, "src", "faq", "content", locale+".json")
			if _, err := os.Stat(content); err != nil {
				errs = append(errs, fmt.Sprintf("locale %q has no content file at src/faq/content/%s.json", locale, locale))
			}
		}
	}

	return errs
}

// resourcePaths returns every file path the manifest references, relative to
// the repo root (and, after packaging, to the archive root).
func resourcePaths(m *Manifest) []string {
	var paths []string
	add := func(p string) {
		if strings.TrimSpace(p) != "" {
			paths = append(paths, p)
		}
	}
	add(m.Resources.Icon)
	for _, s := range m.Resources.Screenshots {
		add(s)
	}
	add(m.Resources.Documentation)
	add(m.Resources.License)
	add(m.Resources.Changelog)
	return paths
}

// SupportsTarget reports whether the manifest's supported_architectures list
// covers the given os/arch target. An empty list means unrestricted.
func (m *Manifest) SupportsTarget(targetOS, targetArch string) bool {
	if len(m.SupportedArch) == 0 {
		return true
	}
	want := targetOS + "/" + targetArch
	for _, pair := range m.SupportedArch {
		if strings.TrimSpace(pair) == want {
			return true
		}
	}
	return false
}

// CheckVersionAlignment verifies package.json agrees with the manifest, which
// is the version source of truth.
func CheckVersionAlignment(m *Manifest, packageJSONPath string) error {
	data, err := os.ReadFile(packageJSONPath)
	if err != nil {
		return fmt.Errorf("read package.json: %w", err)
	}
	var pkg struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return fmt.Errorf("parse package.json: %w", err)
	}
	if pkg.Version != m.Version {
		return fmt.Errorf("version mismatch: manifest.json=%s package.json=%s (manifest is the source of truth)", m.Version, pkg.Version)
	}
	return nil
}
