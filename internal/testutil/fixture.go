package testutil

import (
	"time"
	
	"github.com/tom-023/ubm/internal/bookmark"
)

// CreateTestBookmark creates a bookmark for testing with predefined values
func CreateTestBookmark(title, url, category string) *bookmark.Bookmark {
	return &bookmark.Bookmark{
		ID:        "test-id-" + title,
		Title:     title,
		URL:       url,
		Category:  category,
		CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Tags:      []string{},
	}
}

// CreateTestBookmarks creates multiple test bookmarks
func CreateTestBookmarks() []*bookmark.Bookmark {
	return []*bookmark.Bookmark{
		CreateTestBookmark("Go Documentation", "https://golang.org/doc", "programming/go"),
		CreateTestBookmark("Python Tutorial", "https://python.org/tutorial", "programming/python"),
		CreateTestBookmark("React Guide", "https://react.dev", "programming/javascript"),
		CreateTestBookmark("Design Patterns", "https://refactoring.guru", "programming"),
		CreateTestBookmark("GitHub", "https://github.com", "tools"),
		CreateTestBookmark("Stack Overflow", "https://stackoverflow.com", ""),
	}
}

// SampleCategories returns a list of sample categories for testing
func SampleCategories() []string {
	return []string{
		"programming",
		"programming/go",
		"programming/python",
		"programming/javascript",
		"tools",
		"design",
		"design/ui",
		"design/ux",
	}
}

// SampleBookmarkCounts returns sample bookmark counts per category
func SampleBookmarkCounts() map[string]int {
	return map[string]int{
		"programming":            1,
		"programming/go":         1,
		"programming/python":     1,
		"programming/javascript": 1,
		"tools":                  1,
		"":                       1, // uncategorized
	}
}