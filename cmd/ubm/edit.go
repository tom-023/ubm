package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tom-023/ubm/internal/bookmark"
	"github.com/tom-023/ubm/internal/ui"
	"github.com/tom-023/ubm/pkg/validator"
)

func editCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "edit <title>",
		Short: "Edit existing bookmark by title",
		Long: `Edit a bookmark by specifying its title. You can edit the title, URL, or category.

Examples:
  ubm edit "GitHub"
  ubm edit "Google Search"`,
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

			// Select what to edit
			field, err := ui.SelectEditField()
			if err != nil {
				if err.Error() == "cancelled" {
					fmt.Println("Cancelled.")
					return nil
				}
				return err
			}

			// Store original values for comparison
			originalTitle := targetBookmark.Title
			originalURL := targetBookmark.URL
			originalCategory := targetBookmark.Category

			switch field {
			case "Title":
				fmt.Printf("\nOld Title: %s\n", originalTitle)
				fmt.Println("(Press Enter without typing to keep the current title)")
				newTitle, err := ui.PromptString("New title", "")
				if err != nil {
					if err.Error() == "cancelled" {
						fmt.Println("Cancelled.")
						return nil
					}
					return err
				}
				// If user didn't enter anything, keep the old title
				if newTitle == "" {
					newTitle = originalTitle
				}
				targetBookmark.SetTitle(newTitle)

			case "URL":
				fmt.Printf("\nOld URL: %s\n", originalURL)
				fmt.Println("(Press Enter without typing to keep the current URL)")
				newURL, err := ui.PromptString("New URL", "")
				if err != nil {
					if err.Error() == "cancelled" {
						fmt.Println("Cancelled.")
						return nil
					}
					return err
				}
				// If user didn't enter anything, keep the old URL
				if newURL == "" {
					newURL = originalURL
				}
				// Normalize and validate URL
				newURL, err = validator.NormalizeURL(newURL)
				if err != nil {
					return fmt.Errorf("invalid URL: %w", err)
				}
				targetBookmark.SetURL(newURL)
			}

			// Show changes summary
			fmt.Println("\n--- Changes Summary ---")
			if originalTitle != targetBookmark.Title {
				fmt.Printf("Title: %s → %s\n", originalTitle, targetBookmark.Title)
			}
			if originalURL != targetBookmark.URL {
				fmt.Printf("URL: %s → %s\n", originalURL, targetBookmark.URL)
			}
			if originalCategory != targetBookmark.Category {
				fmt.Printf("Category: %s → %s\n", formatCategory(originalCategory), formatCategory(targetBookmark.Category))
			}
			fmt.Println("---------------------")

			// Confirm changes
			confirm, err := ui.Confirm("Save changes?")
			if err != nil {
				if err.Error() == "cancelled" {
					fmt.Println("Cancelled.")
					return nil
				}
				return err
			}
			if !confirm {
				fmt.Println("Edit cancelled.")
				return nil
			}

			// Update bookmark
			if err := store.UpdateBookmark(targetBookmark); err != nil {
				return fmt.Errorf("failed to update bookmark: %w", err)
			}

			fmt.Printf("✅ Bookmark updated successfully!\n")
			fmt.Printf("Title: %s\n", targetBookmark.Title)
			fmt.Printf("URL: %s\n", targetBookmark.URL)
			fmt.Printf("Category: %s\n", formatCategory(targetBookmark.Category))

			return nil
		},
	}
}