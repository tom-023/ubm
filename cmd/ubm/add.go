package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tom-023/ubm/internal/bookmark"
	"github.com/tom-023/ubm/internal/ui"
	"github.com/tom-023/ubm/pkg/validator"
)

func addCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new URL bookmark interactively",
		Long:  `Add a new URL bookmark to your collection through interactive prompts.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			var url string
			var title string
			var err error

			// Get URL
			url, err = ui.PromptURL("")
			if err != nil {
				if ui.IsCancelError(err) {
					return ui.HandleCancelError(err)
				}
				return fmt.Errorf("failed to get URL: %w", err)
			}

			// Normalize and validate URL
			url, err = validator.NormalizeURL(url)
			if err != nil {
				return fmt.Errorf("invalid URL: %w", err)
			}

			// Get title
			// TODO: Auto-detect title from URL
			title, err = ui.PromptString("Title", extractDomainFromURL(url))
			if err != nil {
				if ui.IsCancelError(err) {
					return ui.HandleCancelError(err)
				}
				return fmt.Errorf("failed to get title: %w", err)
			}

			// Load existing data for category selection
			data, err := store.Load()
			if err != nil {
				return fmt.Errorf("failed to load data: %w", err)
			}

			// Build category tree
			categoryTree := ui.BuildCategoryTree(data)

			// Select category
			selectedCategory, err := ui.SelectCategory(categoryTree, "")
			if err != nil {
				if ui.IsCancelError(err) {
					return ui.HandleCancelError(err)
				}
				return fmt.Errorf("failed to select category: %w", err)
			}

			// Create new category if needed
			if selectedCategory != "" {
				categoryExists := false
				for _, cat := range data.Categories {
					if cat == selectedCategory {
						categoryExists = true
						break
					}
				}
				if !categoryExists {
					data.Categories = append(data.Categories, selectedCategory)
				}
			}

			// Create bookmark
			b := bookmark.New(title, url, selectedCategory)

			// Save bookmark
			if err := store.AddBookmark(b); err != nil {
				return fmt.Errorf("failed to save bookmark: %w", err)
			}

			fmt.Printf("✅ Bookmark added successfully!\n")
			fmt.Printf("Title: %s\n", b.Title)
			fmt.Printf("URL: %s\n", b.URL)
			fmt.Printf("Category: %s\n", ui.FormatCategory(b.Category))

			return nil
		},
	}

	return cmd
}

func extractDomainFromURL(url string) string {
	// Simple domain extraction
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "https://")
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		domain := parts[0]
		// Remove www. prefix
		domain = strings.TrimPrefix(domain, "www.")
		return domain
	}
	return url
}

