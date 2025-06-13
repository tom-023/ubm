package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tom-023/ubm/internal/category"
	"github.com/tom-023/ubm/internal/ui"
)

func moveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "move",
		Short: "Move bookmark to different category",
		Long:  `Interactively select a bookmark and move it to a different category.`,
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

			// Select bookmark to move
			bookmark, err := ui.SelectBookmark(data.Bookmarks, "Select bookmark to move")
			if err != nil {
				return err
			}

			// Show current category
			currentCategory := bookmark.Category
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
			newCategory, err := ui.SelectCategory(categoryTree, bookmark.Category)
			if err != nil {
				return fmt.Errorf("failed to select category: %w", err)
			}

			// Check if category changed
			if newCategory == bookmark.Category {
				fmt.Println("Category unchanged.")
				return nil
			}

			// Confirm move
			confirmMsg := fmt.Sprintf("Move '%s' from '%s' to '%s'?", 
				bookmark.Title, 
				formatCategory(bookmark.Category), 
				formatCategory(newCategory))
			
			confirm, err := ui.Confirm(confirmMsg)
			if err != nil {
				return err
			}
			if !confirm {
				fmt.Println("Move cancelled.")
				return nil
			}

			// Update bookmark category
			bookmark.SetCategory(newCategory)

			// Update bookmark in storage
			if err := store.UpdateBookmark(bookmark); err != nil {
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
			fmt.Printf("Title: %s\n", bookmark.Title)
			fmt.Printf("New category: %s\n", formatCategory(newCategory))

			return nil
		},
	}
}