package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tom-023/ubm/internal/cmd/helpers"
	"github.com/tom-023/ubm/internal/ui"
	"github.com/tom-023/ubm/pkg/validator"
)

func editCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "edit",
		Short: "Edit existing bookmark",
		Long:  `Interactively select a bookmark and edit its title or URL.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load data
			data, categoryTree, err := helpers.LoadDataAndBuildTree(store)
			if err != nil {
				return err
			}

			if len(data.Bookmarks) == 0 {
				fmt.Println("No bookmarks found.")
				return nil
			}

			// Navigate and select bookmark
			targetBookmark, err := ui.NavigateAndSelectBookmark(categoryTree, data.Bookmarks, "Select bookmark to edit")
			if err != nil {
				return helpers.HandleCancelError(err)
			}

			// Select what to edit
			field, err := ui.SelectEditField()
			if err != nil {
				return helpers.HandleCancelError(err)
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
					return helpers.HandleCancelError(err)
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
					return helpers.HandleCancelError(err)
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
				fmt.Printf("Category: %s → %s\n", ui.FormatCategory(originalCategory), ui.FormatCategory(targetBookmark.Category))
			}
			fmt.Println("---------------------")

			// Confirm changes
			confirm, err := ui.Confirm("Save changes?")
			if err != nil {
				return helpers.HandleCancelError(err)
			}
			if !confirm {
				fmt.Println("Edit cancelled.")
				return nil
			}

			// Update bookmark
			if err := store.UpdateBookmark(targetBookmark); err != nil {
				return fmt.Errorf("failed to update bookmark: %w", err)
			}

			helpers.PrintBookmarkSuccess("updated", targetBookmark)

			return nil
		},
	}
}