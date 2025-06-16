package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tom-023/ubm/internal/bookmark"
	"github.com/tom-023/ubm/internal/category"
	"github.com/tom-023/ubm/internal/ui"
)

func moveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "move <title>",
		Short: "Move bookmark to different category by title",
		Long: `Move a bookmark to a different category by specifying its title.

Examples:
  ubm move "GitHub"
  ubm move "Google Search"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load data
			data, err := store.Load()
			if err != nil {
				return fmt.Errorf("failed to load data: %w", err)
			}

			if len(data.Bookmarks) == 0 {
				fmt.Println("No bookmarks found.")
				return nil
			}

			// Find bookmark by title (exact match first, then partial match)
			targetTitle := args[0]
			var targetBookmark *bookmark.Bookmark
			var candidates []*bookmark.Bookmark

			// First, try exact match
			for _, b := range data.Bookmarks {
				if b.Title == targetTitle {
					targetBookmark = b
					break
				}
			}

			// If no exact match, find partial matches
			if targetBookmark == nil {
				for _, b := range data.Bookmarks {
					if strings.Contains(strings.ToLower(b.Title), strings.ToLower(targetTitle)) {
						candidates = append(candidates, b)
					}
				}

				if len(candidates) == 0 {
					return fmt.Errorf("no bookmark found matching '%s'", targetTitle)
				} else if len(candidates) == 1 {
					targetBookmark = candidates[0]
				} else {
					// Multiple matches found, let user select
					targetBookmark, err = ui.SelectBookmark(candidates, fmt.Sprintf("Multiple bookmarks found for '%s'. Select one:", targetTitle))
					if err != nil {
						if err.Error() == "cancelled" {
							fmt.Println("Cancelled.")
							return nil
						}
						return err
					}
				}
			}

			// Show current category
			currentCategory := targetBookmark.Category
			if currentCategory == "" {
				currentCategory = "uncategorized"
			}
			fmt.Printf("\nCurrent category: %s\n", currentCategory)

			// Build category tree
			catManager := category.NewManager()
			bookmarkCounts := make(map[string]int)
			for _, b := range data.Bookmarks {
				bookmarkCounts[b.Category]++
			}
			categoryTree := catManager.BuildTree(data.Categories, bookmarkCounts)

			// Select new category
			newCategory, err := ui.SelectCategory(categoryTree, targetBookmark.Category)
			if err != nil {
				if err.Error() == "cancelled" {
					fmt.Println("Cancelled.")
					return nil
				}
				return fmt.Errorf("failed to select category: %w", err)
			}

			// Check if category changed
			if newCategory == targetBookmark.Category {
				fmt.Println("Category unchanged.")
				return nil
			}

			// Confirm move
			confirmMsg := fmt.Sprintf("Move '%s' from '%s' to '%s'?", 
				targetBookmark.Title, 
				formatCategory(targetBookmark.Category), 
				formatCategory(newCategory))
			
			confirm, err := ui.Confirm(confirmMsg)
			if err != nil {
				if err.Error() == "cancelled" {
					fmt.Println("Cancelled.")
					return nil
				}
				return err
			}
			if !confirm {
				fmt.Println("Move cancelled.")
				return nil
			}

			// Update bookmark category
			targetBookmark.SetCategory(newCategory)

			// Update bookmark in storage
			if err := store.UpdateBookmark(targetBookmark); err != nil {
				return fmt.Errorf("failed to update bookmark: %w", err)
			}

			// Add new category if it doesn't exist
			if newCategory != "" {
				categoryExists := false
				for _, cat := range data.Categories {
					if cat == newCategory {
						categoryExists = true
						break
					}
				}
				if !categoryExists {
					data.Categories = append(data.Categories, newCategory)
					if err := store.Save(data); err != nil {
						return fmt.Errorf("failed to save categories: %w", err)
					}
				}
			}

			fmt.Printf("✅ Bookmark moved successfully!\n")
			fmt.Printf("Title: %s\n", targetBookmark.Title)
			fmt.Printf("New category: %s\n", formatCategory(newCategory))

			return nil
		},
	}
}