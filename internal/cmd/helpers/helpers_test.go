package helpers

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/tom-023/ubm/internal/bookmark"
	"github.com/tom-023/ubm/internal/storage"
	"github.com/tom-023/ubm/internal/ui"
)

func TestHandleCancelError(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantErr bool
	}{
		{
			name:    "nil error",
			err:     nil,
			wantErr: false,
		},
		{
			name:    "cancelled error",
			err:     ui.ErrCancelled,
			wantErr: false,
		},
		{
			name:    "regular error",
			err:     errors.New("some error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := HandleCancelError(tt.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleCancelError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnsureCategoryExists(t *testing.T) {
	tests := []struct {
		name     string
		data     *storage.Data
		category string
		want     []string
	}{
		{
			name: "add new category",
			data: &storage.Data{
				Categories: []string{"existing"},
			},
			category: "new",
			want:     []string{"existing", "new"},
		},
		{
			name: "category already exists",
			data: &storage.Data{
				Categories: []string{"existing", "duplicate"},
			},
			category: "duplicate",
			want:     []string{"existing", "duplicate"},
		},
		{
			name: "empty category",
			data: &storage.Data{
				Categories: []string{"existing"},
			},
			category: "",
			want:     []string{"existing"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			EnsureCategoryExists(tt.data, tt.category)
			if len(tt.data.Categories) != len(tt.want) {
				t.Errorf("EnsureCategoryExists() got %v categories, want %v", len(tt.data.Categories), len(tt.want))
			}
			for i, cat := range tt.want {
				if i >= len(tt.data.Categories) || tt.data.Categories[i] != cat {
					t.Errorf("EnsureCategoryExists() got categories %v, want %v", tt.data.Categories, tt.want)
					break
				}
			}
		})
	}
}

func TestPrintBookmarkSuccess(t *testing.T) {
	// This function only prints to stdout, so we're mainly testing that it doesn't panic
	b := &bookmark.Bookmark{
		ID:       "test-id",
		Title:    "Test Title",
		URL:      "https://example.com",
		Category: "test/category",
	}

	// Test that it doesn't panic with various operations
	operations := []string{"added", "updated", "deleted", "moved"}
	for _, op := range operations {
		t.Run(op, func(t *testing.T) {
			// If this panics, the test will fail
			PrintBookmarkSuccess(op, b)
		})
	}
}

func TestLoadDataAndBuildTree(t *testing.T) {
	// Create a temporary directory for test
	tmpDir, err := os.MkdirTemp("", "ubm-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test storage
	store, err := storage.New(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Add test data
	testData := &storage.Data{
		Bookmarks: []*bookmark.Bookmark{
			{
				ID:       "1",
				Title:    "Test 1",
				URL:      "https://test1.com",
				Category: "dev/go",
			},
			{
				ID:       "2", 
				Title:    "Test 2",
				URL:      "https://test2.com",
				Category: "dev/python",
			},
		},
		Categories: []string{"dev", "dev/go", "dev/python"},
	}

	// Save test data
	if err := store.Save(testData); err != nil {
		t.Fatalf("Failed to save test data: %v", err)
	}

	// Test LoadDataAndBuildTree
	data, tree, err := LoadDataAndBuildTree(store)
	if err != nil {
		t.Errorf("LoadDataAndBuildTree() error = %v", err)
	}

	if data == nil {
		t.Error("LoadDataAndBuildTree() returned nil data")
	}

	if tree == nil {
		t.Error("LoadDataAndBuildTree() returned nil tree")
	}

	if len(data.Bookmarks) != 2 {
		t.Errorf("LoadDataAndBuildTree() got %d bookmarks, want 2", len(data.Bookmarks))
	}

	// Test with non-existent storage
	nonExistentDir := filepath.Join(tmpDir, "non-existent")
	badStore, err := storage.New(nonExistentDir)
	if err != nil {
		t.Fatalf("Failed to create storage with non-existent dir: %v", err)
	}

	// This should still work (returns empty data)
	data2, tree2, err := LoadDataAndBuildTree(badStore)
	if err != nil {
		t.Errorf("LoadDataAndBuildTree() with new storage error = %v", err)
	}

	if data2 == nil || tree2 == nil {
		t.Error("LoadDataAndBuildTree() should return non-nil data and tree even for new storage")
	}
}