package ui

import (
	"testing"
	"time"

	"github.com/tom-023/ubm/internal/bookmark"
	"github.com/tom-023/ubm/internal/storage"
)

func TestBuildCategoryTree(t *testing.T) {
	// Create test data
	bookmarks := []*bookmark.Bookmark{
		{Category: "programming"},
		{Category: "programming"},
		{Category: "programming/go"},
		{Category: "tools"},
		{Category: ""}, // uncategorized
	}
	
	categories := []string{
		"programming",
		"programming/go",
		"tools",
	}
	
	data := &storage.Data{
		Bookmarks:  bookmarks,
		Categories: categories,
		UpdatedAt:  time.Now(),
	}
	
	// Build tree
	tree := BuildCategoryTree(data)
	
	// Verify root node
	if tree == nil {
		t.Fatal("BuildCategoryTree returned nil")
	}
	
	if !tree.IsRoot {
		t.Error("Root node should have IsRoot = true")
	}
	
	// Verify children count (should have programming, tools, and uncategorized)
	expectedChildCount := 3
	if len(tree.Children) != expectedChildCount {
		t.Errorf("Root should have %d children, got %d", expectedChildCount, len(tree.Children))
	}
}

func TestCountBookmarksByCategory(t *testing.T) {
	bookmarks := []*bookmark.Bookmark{
		{Category: "programming"},
		{Category: "programming"},
		{Category: "programming/go"},
		{Category: "tools"},
		{Category: "tools"},
		{Category: "tools"},
		{Category: ""}, // uncategorized
		{Category: ""}, // uncategorized
	}
	
	counts := CountBookmarksByCategory(bookmarks)
	
	tests := []struct {
		category string
		want     int
	}{
		{"programming", 2},
		{"programming/go", 1},
		{"tools", 3},
		{"", 2},
		{"nonexistent", 0},
	}
	
	for _, tt := range tests {
		t.Run(tt.category, func(t *testing.T) {
			got := counts[tt.category]
			if got != tt.want {
				t.Errorf("CountBookmarksByCategory()[%q] = %d, want %d", tt.category, got, tt.want)
			}
		})
	}
}