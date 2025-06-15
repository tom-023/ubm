package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/tom-023/ubm/internal/bookmark"
	"github.com/tom-023/ubm/internal/testutil"
)

func TestNew(t *testing.T) {
	dir, cleanup := testutil.TempDir(t)
	defer cleanup()

	configDir := filepath.Join(dir, ".config", "ubm")
	s, err := New(configDir)

	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if s == nil {
		t.Fatal("New() returned nil")
	}

	expectedPath := filepath.Join(configDir, "bookmarks.json")
	if s.filePath != expectedPath {
		t.Errorf("filePath = %v, want %v", s.filePath, expectedPath)
	}

	expectedBackupPath := filepath.Join(configDir, "bookmarks.backup.json")
	if s.backupPath != expectedBackupPath {
		t.Errorf("backupPath = %v, want %v", s.backupPath, expectedBackupPath)
	}

	// Check that directory was created
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Error("Config directory was not created")
	}
}

func TestStorage_Load_Empty(t *testing.T) {
	dir, cleanup := testutil.TempDir(t)
	defer cleanup()

	s, err := New(dir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Load when file doesn't exist
	data, err := s.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if len(data.Bookmarks) != 0 {
		t.Errorf("Expected empty bookmarks, got %d", len(data.Bookmarks))
	}

	if len(data.Categories) != 0 {
		t.Errorf("Expected empty categories, got %d", len(data.Categories))
	}

	if data.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should not be zero for new data")
	}
}

func TestStorage_SaveAndLoad(t *testing.T) {
	dir, cleanup := testutil.TempDir(t)
	defer cleanup()

	s, err := New(dir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Create test data
	bookmarks := testutil.CreateTestBookmarks()
	data := &Data{
		Bookmarks:  bookmarks[:3], // Use first 3 bookmarks
		Categories: []string{"programming", "programming/go", "programming/python"},
		UpdatedAt:  time.Now(),
	}

	// Save data
	if err := s.Save(data); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Load data
	loaded, err := s.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Compare bookmarks
	if len(loaded.Bookmarks) != len(data.Bookmarks) {
		t.Errorf("Loaded %d bookmarks, want %d", len(loaded.Bookmarks), len(data.Bookmarks))
	}

	for i, b := range loaded.Bookmarks {
		if b.ID != data.Bookmarks[i].ID {
			t.Errorf("Bookmark[%d].ID = %v, want %v", i, b.ID, data.Bookmarks[i].ID)
		}
		if b.Title != data.Bookmarks[i].Title {
			t.Errorf("Bookmark[%d].Title = %v, want %v", i, b.Title, data.Bookmarks[i].Title)
		}
		if b.URL != data.Bookmarks[i].URL {
			t.Errorf("Bookmark[%d].URL = %v, want %v", i, b.URL, data.Bookmarks[i].URL)
		}
	}

	// Compare categories
	if !reflect.DeepEqual(loaded.Categories, data.Categories) {
		t.Errorf("Categories = %v, want %v", loaded.Categories, data.Categories)
	}

	// Check that UpdatedAt was updated
	if !loaded.UpdatedAt.After(data.UpdatedAt) && !loaded.UpdatedAt.Equal(data.UpdatedAt) {
		t.Error("UpdatedAt should be updated or equal")
	}
}

func TestStorage_Backup(t *testing.T) {
	dir, cleanup := testutil.TempDir(t)
	defer cleanup()

	s, err := New(dir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Save initial data
	data1 := &Data{
		Bookmarks:  []*bookmark.Bookmark{testutil.CreateTestBookmark("Test1", "https://test1.com", "test")},
		Categories: []string{"test"},
	}
	if err := s.Save(data1); err != nil {
		t.Fatalf("First save error = %v", err)
	}

	// Save new data (should create backup)
	data2 := &Data{
		Bookmarks:  []*bookmark.Bookmark{testutil.CreateTestBookmark("Test2", "https://test2.com", "test")},
		Categories: []string{"test"},
	}
	if err := s.Save(data2); err != nil {
		t.Fatalf("Second save error = %v", err)
	}

	// Check backup file exists and contains original data
	backupData, err := os.ReadFile(s.backupPath)
	if err != nil {
		t.Fatalf("Failed to read backup file: %v", err)
	}

	var backup Data
	if err := json.Unmarshal(backupData, &backup); err != nil {
		t.Fatalf("Failed to unmarshal backup: %v", err)
	}

	if len(backup.Bookmarks) != 1 || backup.Bookmarks[0].Title != "Test1" {
		t.Error("Backup should contain original data")
	}
}

func TestStorage_AddBookmark(t *testing.T) {
	dir, cleanup := testutil.TempDir(t)
	defer cleanup()

	s, err := New(dir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Add first bookmark
	b1 := testutil.CreateTestBookmark("Test1", "https://test1.com", "category1")
	if err := s.AddBookmark(b1); err != nil {
		t.Fatalf("AddBookmark() error = %v", err)
	}

	// Verify bookmark was added
	data, _ := s.Load()
	if len(data.Bookmarks) != 1 {
		t.Errorf("Expected 1 bookmark, got %d", len(data.Bookmarks))
	}
	if len(data.Categories) != 1 || data.Categories[0] != "category1" {
		t.Errorf("Expected category to be added")
	}

	// Try to add duplicate (same URL and category)
	b2 := testutil.CreateTestBookmark("Test1 Copy", "https://test1.com", "category1")
	err = s.AddBookmark(b2)
	if err == nil {
		t.Error("Expected error when adding duplicate bookmark")
	}

	// Add same URL in different category (should succeed)
	b3 := testutil.CreateTestBookmark("Test1 Different", "https://test1.com", "category2")
	if err := s.AddBookmark(b3); err != nil {
		t.Fatalf("AddBookmark() with different category error = %v", err)
	}

	// Add bookmark without category
	b4 := testutil.CreateTestBookmark("Uncategorized", "https://uncat.com", "")
	if err := s.AddBookmark(b4); err != nil {
		t.Fatalf("AddBookmark() uncategorized error = %v", err)
	}

	// Verify final state
	data, _ = s.Load()
	if len(data.Bookmarks) != 3 {
		t.Errorf("Expected 3 bookmarks, got %d", len(data.Bookmarks))
	}
	if len(data.Categories) != 2 {
		t.Errorf("Expected 2 categories, got %d", len(data.Categories))
	}
}

func TestStorage_UpdateBookmark(t *testing.T) {
	dir, cleanup := testutil.TempDir(t)
	defer cleanup()

	s, err := New(dir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Add initial bookmark
	b := testutil.CreateTestBookmark("Original", "https://original.com", "test")
	if err := s.AddBookmark(b); err != nil {
		t.Fatalf("AddBookmark() error = %v", err)
	}

	// Update bookmark
	b.Title = "Updated"
	b.URL = "https://updated.com"
	if err := s.UpdateBookmark(b); err != nil {
		t.Fatalf("UpdateBookmark() error = %v", err)
	}

	// Verify update
	loaded, _ := s.GetBookmark(b.ID)
	if loaded.Title != "Updated" {
		t.Errorf("Title = %v, want Updated", loaded.Title)
	}
	if loaded.URL != "https://updated.com" {
		t.Errorf("URL = %v, want https://updated.com", loaded.URL)
	}

	// Try to update non-existent bookmark
	nonExistent := testutil.CreateTestBookmark("NonExistent", "https://none.com", "test")
	err = s.UpdateBookmark(nonExistent)
	if err == nil {
		t.Error("Expected error when updating non-existent bookmark")
	}
}

func TestStorage_DeleteBookmark(t *testing.T) {
	dir, cleanup := testutil.TempDir(t)
	defer cleanup()

	s, err := New(dir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Add bookmarks
	b1 := testutil.CreateTestBookmark("Test1", "https://test1.com", "test")
	b2 := testutil.CreateTestBookmark("Test2", "https://test2.com", "test")
	s.AddBookmark(b1)
	s.AddBookmark(b2)

	// Delete first bookmark
	if err := s.DeleteBookmark(b1.ID); err != nil {
		t.Fatalf("DeleteBookmark() error = %v", err)
	}

	// Verify deletion
	data, _ := s.Load()
	if len(data.Bookmarks) != 1 {
		t.Errorf("Expected 1 bookmark after deletion, got %d", len(data.Bookmarks))
	}
	if data.Bookmarks[0].ID != b2.ID {
		t.Error("Wrong bookmark was deleted")
	}

	// Try to delete non-existent bookmark
	err = s.DeleteBookmark("non-existent-id")
	if err == nil {
		t.Error("Expected error when deleting non-existent bookmark")
	}
}

func TestStorage_GetBookmark(t *testing.T) {
	dir, cleanup := testutil.TempDir(t)
	defer cleanup()

	s, err := New(dir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Add bookmark
	b := testutil.CreateTestBookmark("Test", "https://test.com", "test")
	s.AddBookmark(b)

	// Get existing bookmark
	loaded, err := s.GetBookmark(b.ID)
	if err != nil {
		t.Fatalf("GetBookmark() error = %v", err)
	}
	if loaded.ID != b.ID {
		t.Errorf("Loaded wrong bookmark")
	}

	// Get non-existent bookmark
	_, err = s.GetBookmark("non-existent-id")
	if err == nil {
		t.Error("Expected error when getting non-existent bookmark")
	}
}

func TestStorage_GetBookmarksByCategory(t *testing.T) {
	dir, cleanup := testutil.TempDir(t)
	defer cleanup()

	s, err := New(dir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Add bookmarks in different categories
	bookmarks := []*bookmark.Bookmark{
		testutil.CreateTestBookmark("Go1", "https://go1.com", "programming/go"),
		testutil.CreateTestBookmark("Go2", "https://go2.com", "programming/go"),
		testutil.CreateTestBookmark("Python1", "https://py1.com", "programming/python"),
		testutil.CreateTestBookmark("Uncategorized", "https://uncat.com", ""),
	}

	for _, b := range bookmarks {
		s.AddBookmark(b)
	}

	// Test getting bookmarks by category
	tests := []struct {
		category string
		want     int
	}{
		{"programming/go", 2},
		{"programming/python", 1},
		{"", 1},
		{"non-existent", 0},
	}

	for _, tt := range tests {
		t.Run(tt.category, func(t *testing.T) {
			result, err := s.GetBookmarksByCategory(tt.category)
			if err != nil {
				t.Fatalf("GetBookmarksByCategory() error = %v", err)
			}
			if len(result) != tt.want {
				t.Errorf("Got %d bookmarks, want %d", len(result), tt.want)
			}
		})
	}
}

func TestStorage_SearchBookmarks(t *testing.T) {
	dir, cleanup := testutil.TempDir(t)
	defer cleanup()

	s, err := New(dir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Add bookmarks with various content
	bookmarks := []*bookmark.Bookmark{
		{
			ID:          "1",
			Title:       "Go Programming Language",
			URL:         "https://golang.org",
			Description: "Official Go website",
			Category:    "programming/go",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Tags:        []string{},
		},
		{
			ID:          "2",
			Title:       "Python Documentation",
			URL:         "https://python.org/doc",
			Description: "Learn Python programming",
			Category:    "programming/python",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Tags:        []string{},
		},
		{
			ID:          "3",
			Title:       "GitHub",
			URL:         "https://github.com",
			Description: "Code hosting platform",
			Category:    "tools",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Tags:        []string{},
		},
	}

	for _, b := range bookmarks {
		s.AddBookmark(b)
	}

	// Test search functionality
	tests := []struct {
		query string
		want  []string // Expected bookmark IDs
	}{
		{"go", []string{"1"}},                    // Title match
		{"python", []string{"2"}},                // Title and description match
		{"programming", []string{"1", "2"}},      // Title and description match
		{"github", []string{"3"}},                // URL match
		{"GITHUB", []string{"3"}},                // Case insensitive
		{"platform", []string{"3"}},              // Description match
		{"nonexistent", []string{}},              // No match
		{"https://", []string{"1", "2", "3"}},    // URL prefix match
	}

	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			results, err := s.SearchBookmarks(tt.query)
			if err != nil {
				t.Fatalf("SearchBookmarks() error = %v", err)
			}

			gotIDs := []string{}
			for _, b := range results {
				gotIDs = append(gotIDs, b.ID)
			}

			if len(gotIDs) != len(tt.want) {
				t.Errorf("SearchBookmarks(%q) returned %d results, want %d", tt.query, len(gotIDs), len(tt.want))
				return
			}

			// Check that all expected IDs are present
			for _, wantID := range tt.want {
				found := false
				for _, gotID := range gotIDs {
					if gotID == wantID {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("SearchBookmarks(%q) missing expected bookmark ID %s", tt.query, wantID)
				}
			}
		})
	}
}

func TestStorage_ConcurrentAccess(t *testing.T) {
	dir, cleanup := testutil.TempDir(t)
	defer cleanup()

	s, err := New(dir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Add initial bookmark
	b := testutil.CreateTestBookmark("Test", "https://test.com", "test")
	s.AddBookmark(b)

	// Simulate concurrent reads and writes
	done := make(chan bool)

	// Reader goroutine
	go func() {
		for i := 0; i < 10; i++ {
			_, err := s.Load()
			if err != nil {
				t.Errorf("Concurrent Load() error = %v", err)
			}
		}
		done <- true
	}()

	// Writer goroutine
	go func() {
		for i := 0; i < 10; i++ {
			newBookmark := testutil.CreateTestBookmark("Concurrent", "https://concurrent.com", "test")
			newBookmark.ID = newBookmark.ID + string(rune(i))
			err := s.AddBookmark(newBookmark)
			if err != nil && err.Error() != "bookmark with URL https://concurrent.com already exists in category test" {
				t.Errorf("Concurrent AddBookmark() error = %v", err)
			}
		}
		done <- true
	}()

	// Wait for both goroutines to complete
	<-done
	<-done

	// No assertions needed - test passes if no data races or panics occur
}

func TestContainsIgnoreCase(t *testing.T) {
	tests := []struct {
		s      string
		substr string
		want   bool
	}{
		{"Hello World", "hello", true},
		{"Hello World", "WORLD", true},
		{"Hello World", "ello", true},
		{"Hello World", "xyz", false},
		{"", "test", false},
		{"test", "", true},
		{"Test", "test", true},
		{"test", "TEST", true},
		{"Programming", "gram", true},
		{"日本語", "本", true},
	}

	for _, tt := range tests {
		t.Run(tt.s+"/"+tt.substr, func(t *testing.T) {
			got := containsIgnoreCase(tt.s, tt.substr)
			if got != tt.want {
				t.Errorf("containsIgnoreCase(%q, %q) = %v, want %v", tt.s, tt.substr, got, tt.want)
			}
		})
	}
}