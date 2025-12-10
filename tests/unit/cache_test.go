package unit

import (
	"os"
	"testing"

	"ut-plugin-faq/src/storage"
)

func TestSaveRejectsBadChecksum(t *testing.T) {
	dir := t.TempDir()
	c := storage.Cache{BasePath: dir}
	data := []byte("hello")
	err := c.Save("en-US", data, "sha256:deadbeef")
	if err == nil {
		t.Fatalf("expected checksum mismatch error")
	}
	if _, err := os.Stat(c.BasePath + "/en-US.json"); !os.IsNotExist(err) {
		t.Fatalf("expected no file written on mismatch")
	}
}

func TestSaveAndLoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	c := storage.Cache{BasePath: dir}
	data := []byte("hello")
	// allow cache to compute checksum
	if err := c.Save("en-US", data, ""); err != nil {
		t.Fatalf("save failed: %v", err)
	}
	loaded, checksum, err := c.Load("en-US")
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if string(loaded) != string(data) {
		t.Fatalf("data mismatch")
	}
	if checksum == "" {
		t.Fatalf("expected checksum computed")
	}
}

func TestApplyUpdateChecksumMismatch(t *testing.T) {
	dir := t.TempDir()
	c := storage.Cache{BasePath: dir}
	res := c.ApplyUpdate("en-US", []byte("bad"), "sha256:deadbeef")
	if res.Applied || !res.PreviousKept || !res.ChecksumError {
		t.Fatalf("expected checksum mismatch handling, got %+v", res)
	}
}

func TestRemoveNonExistingIsNoop(t *testing.T) {
	dir := t.TempDir()
	c := storage.Cache{BasePath: dir}
	if err := c.Remove("en-US"); err != nil {
		t.Fatalf("remove should not error: %v", err)
	}
}
