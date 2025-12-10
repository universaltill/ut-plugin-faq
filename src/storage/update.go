package storage

import "fmt"

// UpdateResult captures the outcome of an update attempt.
type UpdateResult struct {
	Applied       bool
	PreviousKept  bool
	ChecksumError bool
	Message       string
}

// ApplyUpdate attempts to save new bundle data with checksum gating, preserving prior data on failure.
func (c Cache) ApplyUpdate(locale string, data []byte, expectedChecksum string) UpdateResult {
	if expectedChecksum == "" {
		return UpdateResult{Applied: false, PreviousKept: true, ChecksumError: true, Message: "missing expected checksum"}
	}
	actual := computeChecksum(data)
	if expectedChecksum != actual {
		return UpdateResult{Applied: false, PreviousKept: true, ChecksumError: true, Message: fmt.Sprintf("checksum mismatch expected %s got %s", expectedChecksum, actual)}
	}
	if err := c.Save(locale, data, expectedChecksum); err != nil {
		return UpdateResult{Applied: false, PreviousKept: true, Message: err.Error()}
	}
	return UpdateResult{Applied: true, PreviousKept: false, Message: "update applied"}
}
