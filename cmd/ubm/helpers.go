package main

import (
	"fmt"

	"github.com/tom-023/ubm/internal/storage"
	"github.com/tom-023/ubm/internal/ui"
	"github.com/tom-023/ubm/internal/bookmark"
	"github.com/tom-023/ubm/internal/category"
)

// handleCancelError handles cancellation errors consistently across commands
func handleCancelError(err error) error {
	if err == nil {
		return nil
	}
	if ui.IsCancelError(err) {
		fmt.Println("Cancelled.")
		return nil
	}
	return err
}

// loadDataAndBuildTree loads storage data and builds category tree
func loadDataAndBuildTree() (*storage.Data, *category.Node, error) {
	data, err := store.Load()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load data: %w", err)
	}

	categoryTree := ui.BuildCategoryTree(data)
	return data, categoryTree, nil
}

// ensureCategoryExists adds a category to the data if it doesn't exist
func ensureCategoryExists(data *storage.Data, category string) {
	if category == "" {
		return
	}

	for _, cat := range data.Categories {
		if cat == category {
			return
		}
	}
	data.Categories = append(data.Categories, category)
}

// printBookmarkSuccess prints a success message for bookmark operations
func printBookmarkSuccess(operation string, b *bookmark.Bookmark) {
	fmt.Printf("✅ Bookmark %s successfully!\n", operation)
	fmt.Printf("Title: %s\n", b.Title)
	fmt.Printf("URL: %s\n", b.URL)
	fmt.Printf("Category: %s\n", ui.FormatCategory(b.Category))
}