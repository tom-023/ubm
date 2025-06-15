package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

// TempDir creates a temporary directory for tests and returns a cleanup function
func TempDir(t *testing.T) (string, func()) {
	t.Helper()
	
	dir, err := os.MkdirTemp("", "ubm-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	
	cleanup := func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Errorf("Failed to clean up temp dir %s: %v", dir, err)
		}
	}
	
	return dir, cleanup
}

// TempFile creates a temporary file with the given content
func TempFile(t *testing.T, dir, pattern string, content []byte) string {
	t.Helper()
	
	file, err := os.CreateTemp(dir, pattern)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer file.Close()
	
	if content != nil {
		if _, err := file.Write(content); err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
	}
	
	return file.Name()
}

// CreateTestConfig creates a test configuration directory with an empty bookmarks.json
func CreateTestConfig(t *testing.T) (string, func()) {
	t.Helper()
	
	dir, cleanup := TempDir(t)
	configDir := filepath.Join(dir, ".config", "ubm")
	
	if err := os.MkdirAll(configDir, 0755); err != nil {
		cleanup()
		t.Fatalf("Failed to create config dir: %v", err)
	}
	
	// Create empty bookmarks.json
	bookmarksPath := filepath.Join(configDir, "bookmarks.json")
	content := []byte(`{"bookmarks":[],"categories":[],"updated_at":"2024-01-01T00:00:00Z"}`)
	
	if err := os.WriteFile(bookmarksPath, content, 0644); err != nil {
		cleanup()
		t.Fatalf("Failed to create bookmarks.json: %v", err)
	}
	
	return dir, cleanup
}