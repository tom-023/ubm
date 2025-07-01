package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tom-023/ubm/internal/ui"
)

func moveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "move",
		Short: "Move bookmark to different category",
		Long:  `Interactively select a bookmark and move it to a different category.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load data
			data, categoryTree, err := loadDataAndBuildTree()
			if err != nil {
				return err
			}

			if len(data.Bookmarks) == 0 {
				fmt.Println("No bookmarks found.")
				return nil
			}

			// Navigate and select bookmark
			targetBookmark, err := ui.NavigateAndSelectBookmark(categoryTree, data.Bookmarks, "Select bookmark to move")
			if err != nil {
				return handleCancelError(err)
			}

			// Store original category
			originalCategory := targetBookmark.Category

			// Display current category
			fmt.Printf("\nMoving bookmark: %s\n", targetBookmark.Title)
			fmt.Printf("Current category: %s\n", ui.FormatCategory(originalCategory))

			// Select new category
			newCategory, err := ui.SelectCategory(categoryTree, originalCategory)
			if err != nil {
				return handleCancelError(err)
			}

			// Check if category actually changed
			if newCategory == originalCategory {
				fmt.Println("Bookmark remains in the same category.")
				return nil
			}

			// Update bookmark category
			targetBookmark.SetCategory(newCategory)

			// Add new category if it doesn't exist
			ensureCategoryExists(data, newCategory)

			// Update bookmark
			if err := store.UpdateBookmark(targetBookmark); err != nil {
				return fmt.Errorf("failed to update bookmark: %w", err)
			}

			// Save categories if new one was added
			if err := store.Save(data); err != nil {
				return fmt.Errorf("failed to save data: %w", err)
			}

			fmt.Printf("\n✅ Bookmark moved successfully!\n")
			fmt.Printf("Title: %s\n", targetBookmark.Title)
			fmt.Printf("From: %s\n", ui.FormatCategory(originalCategory))
			fmt.Printf("To: %s\n", ui.FormatCategory(newCategory))

			return nil
		},
	}
}