package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tom-023/ubm/internal/category"
	"github.com/tom-023/ubm/internal/ui"
	"github.com/tom-023/ubm/pkg/validator"
)

func editCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "edit",
		Short: "Edit existing bookmark",
		Long:  `Interactively select a bookmark and edit its title, URL, or category.`,
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

			// Select bookmark to edit
			bookmark, err := ui.SelectBookmark(data.Bookmarks, "Select bookmark to edit")
			if err != nil {
				return err
			}

			// Select what to edit
			field, err := ui.SelectEditField()
			if err != nil {
				return err
			}

			switch field {
			case "Title":
				newTitle, err := ui.PromptString("New title", bookmark.Title)
				if err != nil {
					return err
				}
				bookmark.SetTitle(newTitle)

			case "URL":
				newURL, err := ui.PromptURL(bookmark.URL)
				if err != nil {
					return err
				}
				// Normalize and validate URL
				newURL, err = validator.NormalizeURL(newURL)
				if err != nil {
					return fmt.Errorf("invalid URL: %w", err)
				}
				bookmark.SetURL(newURL)

			case "Category":
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
					return err
				}
				bookmark.SetCategory(newCategory)

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

			case "All fields":
				// Edit title
				newTitle, err := ui.PromptString("New title", bookmark.Title)
				if err != nil {
					return err
				}
				bookmark.SetTitle(newTitle)

				// Edit URL
				newURL, err := ui.PromptURL(bookmark.URL)
				if err != nil {
					return err
				}
				// Normalize and validate URL
				newURL, err = validator.NormalizeURL(newURL)
				if err != nil {
					return fmt.Errorf("invalid URL: %w", err)
				}
				bookmark.SetURL(newURL)

				// Edit category
				catManager := category.NewManager()
				bookmarkCounts := make(map[string]int)
				for _, b := range data.Bookmarks {
					bookmarkCounts[b.Category]++
				}
				categoryTree := catManager.BuildTree(data.Categories, bookmarkCounts)

				newCategory, err := ui.SelectCategory(categoryTree, bookmark.Category)
				if err != nil {
					return err
				}
				bookmark.SetCategory(newCategory)

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
			}

			// Confirm changes
			confirm, err := ui.Confirm("Save changes?")
			if err != nil {
				return err
			}
			if !confirm {
				fmt.Println("Edit cancelled.")
				return nil
			}

			// Update bookmark
			if err := store.UpdateBookmark(bookmark); err != nil {
				return fmt.Errorf("failed to update bookmark: %w", err)
			}

			fmt.Printf("✅ Bookmark updated successfully!\n")
			fmt.Printf("Title: %s\n", bookmark.Title)
			fmt.Printf("URL: %s\n", bookmark.URL)
			fmt.Printf("Category: %s\n", formatCategory(bookmark.Category))

			return nil
		},
	}
}