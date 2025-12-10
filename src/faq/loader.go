package faq

import (
	"crypto/sha256"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
)

//go:embed content/*.json
var embeddedContent embed.FS

type Category struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	SortOrder int    `json:"sort_order"`
}

type Entry struct {
	ID          string   `json:"id"`
	Category    string   `json:"category"`
	Question    string   `json:"question"`
	Answer      string   `json:"answer"`
	LastUpdated string   `json:"last_updated"`
	SortOrder   int      `json:"sort_order"`
	Keywords    []string `json:"keywords"`
}

type Bundle struct {
	Locale         string     `json:"locale"`
	Version        string     `json:"version"`
	ChecksumSHA256 string     `json:"checksum_sha256"`
	RTL            bool       `json:"rtl"`
	FallbackLocale string     `json:"fallback_locale"`
	Categories     []Category `json:"categories"`
	FAQEntries     []Entry    `json:"faq_entries"`
}

// ComputeChecksum returns the sha256 checksum with a "sha256:" prefix to align with manifest patterns.
func ComputeChecksum(data []byte) string {
	sum := sha256.Sum256(data)
	return fmt.Sprintf("sha256:%x", sum[:])
}

func validateChecksum(expected, actual string) error {
	if expected == "" || expected == "sha256:" {
		return nil
	}
	if expected != actual {
		return fmt.Errorf("checksum mismatch: expected %s got %s", expected, actual)
	}
	return nil
}

func normalizeFallback(b *Bundle) {
	if b.FallbackLocale == "" {
		b.FallbackLocale = "en-US"
	}
}

// ParseBundle decodes JSON bundle content and enforces checksum validation before returning a parsed bundle.
func ParseBundle(data []byte, expectedChecksum string) (*Bundle, error) {
	actual := ComputeChecksum(data)
	if err := validateChecksum(expectedChecksum, actual); err != nil {
		return nil, err
	}
	var bundle Bundle
	if err := json.Unmarshal(data, &bundle); err != nil {
		return nil, err
	}
	bundle.ChecksumSHA256 = actual
	normalizeFallback(&bundle)
	return &bundle, nil
}

// LoadBundleFromFile loads and validates a bundle from disk.
func LoadBundleFromFile(path, expectedChecksum string) (*Bundle, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseBundle(data, expectedChecksum)
}

// LoadEmbeddedBundle loads a locale bundle from the embedded content directory.
func LoadEmbeddedBundle(locale, expectedChecksum string) (*Bundle, error) {
	path := fmt.Sprintf("content/%s.json", locale)
	data, err := fs.ReadFile(embeddedContent, path)
	if err != nil {
		return nil, err
	}
	return ParseBundle(data, expectedChecksum)
}

// IsRTL reports whether the bundle content should render in RTL direction.
func (b *Bundle) IsRTL() bool {
	return b.RTL
}
