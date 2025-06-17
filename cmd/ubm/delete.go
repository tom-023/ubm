package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tom-023/ubm/internal/ui"
)

func deleteCmd() *cobra.Command {
	var skipConfirm bool

	cmd := &cobra.Command{
		Use:   "delete [ID]",
		Short: "Delete bookmark",
		Long:  `Delete a bookmark by ID or select interactively.`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var bookmarkID string

			if len(args) > 0 {
				bookmarkID = args[0]
			} else {
				// Interactive selection
				data, err := store.Load()
				if err != nil {
					return fmt.Errorf("failed to load data: %w", err)
				}

				if len(data.Bookmarks) == 0 {
					fmt.Println("No bookmarks found.")
					return nil
				}

				// Build category tree
				categoryTree := ui.BuildCategoryTree(data)

				// Navigate and select bookmark
				bookmark, err := ui.NavigateAndSelectBookmark(categoryTree, data.Bookmarks, "Select bookmark to delete")
				if err != nil {
					if ui.IsCancelError(err) {
						fmt.Println("Cancelled.")
						return nil
					}
					return err
				}
				bookmarkID = bookmark.ID
			}

			// Get bookmark details for confirmation
			bookmark, err := store.GetBookmark(bookmarkID)
			if err != nil {
				return fmt.Errorf("bookmark not found: %w", err)
			}

			// Confirm deletion
			if !skipConfirm {
				confirmMsg := fmt.Sprintf("Delete bookmark '%s' (%s)?", bookmark.Title, bookmark.URL)
				confirm, err := ui.Confirm(confirmMsg)
				if err != nil {
					if ui.IsCancelError(err) {
						fmt.Println("Cancelled.")
						return nil
					}
					return err
				}
				if !confirm {
					fmt.Println("Deletion cancelled.")
					return nil
				}
			}

			// Delete bookmark
			if err := store.DeleteBookmark(bookmarkID); err != nil {
				return fmt.Errorf("failed to delete bookmark: %w", err)
			}

			fmt.Printf("✅ Bookmark '%s' deleted successfully!\n", bookmark.Title)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&skipConfirm, "confirm", "y", false, "Skip confirmation prompt")

	return cmd
}