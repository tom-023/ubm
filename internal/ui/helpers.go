package ui

import (
	"strings"

	"github.com/tom-023/ubm/internal/bookmark"
	"github.com/tom-023/ubm/internal/category"
	"github.com/tom-023/ubm/internal/storage"
)

// BuildCategoryTree builds a category tree from storage data
func BuildCategoryTree(data *storage.Data) *category.Node {
	catManager := category.NewManager()
	bookmarkCounts := CountBookmarksByCategory(data.Bookmarks)
	return catManager.BuildTree(data.Categories, bookmarkCounts)
}

// CountBookmarksByCategory counts bookmarks in each category
func CountBookmarksByCategory(bookmarks []*bookmark.Bookmark) map[string]int {
	counts := make(map[string]int)
	for _, b := range bookmarks {
		counts[b.Category]++
	}
	return counts
}

// FormatCategory formats a category string for display
func FormatCategory(cat string) string {
	if cat == "" {
		return "uncategorized"
	}
	return cat
}

// CreateSearcher creates a generic searcher function for promptui
func CreateSearcher(getSearchText func(int) string) func(string, int) bool {
	return func(input string, index int) bool {
		searchText := strings.ToLower(strings.Replace(getSearchText(index), " ", "", -1))
		input = strings.ToLower(strings.Replace(input, " ", "", -1))
		return strings.Contains(searchText, input)
	}
}