package storage

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Cache provides simple offline storage for localized FAQ bundles.
type Cache struct {
	BasePath string
}

// path resolves the on-disk location for a locale bundle.
func (c Cache) path(locale string) string {
	return filepath.Join(c.BasePath, fmt.Sprintf("%s.json", locale))
}

func computeChecksum(data []byte) string {
	sum := sha256.Sum256(data)
	return fmt.Sprintf("sha256:%x", sum[:])
}

// Save writes a bundle to cache after verifying checksum. Corrupted bundles are rejected.
func (c Cache) Save(locale string, data []byte, expectedChecksum string) error {
	actual := computeChecksum(data)
	if expectedChecksum != "" && expectedChecksum != "sha256:" && expectedChecksum != actual {
		return fmt.Errorf("checksum mismatch for %s: expected %s got %s", locale, expectedChecksum, actual)
	}

	if err := os.MkdirAll(c.BasePath, 0o755); err != nil {
		return err
	}

	tmp := c.path(locale) + ".tmp"
	target := c.path(locale)

	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, target)
}

// Load returns cached bundle bytes and the computed checksum.
func (c Cache) Load(locale string) ([]byte, string, error) {
	data, err := os.ReadFile(c.path(locale))
	if err != nil {
		return nil, "", err
	}
	return data, computeChecksum(data), nil
}

// Remove deletes a cached bundle for the given locale.
func (c Cache) Remove(locale string) error {
	if locale == "" {
		return errors.New("locale required")
	}
	err := os.Remove(c.path(locale))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
