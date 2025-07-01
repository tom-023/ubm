package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tom-023/ubm/internal/ui"
)

func listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Interactive navigation of bookmarked URLs",
		Long: `Navigate through your bookmarks interactively.
Use arrow keys to navigate, Enter to select, and q to quit.`,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load bookmarks
			data, categoryTree, err := loadDataAndBuildTree()
			if err != nil {
				return err
			}

			if len(data.Bookmarks) == 0 {
				fmt.Println("No bookmarks found. Use 'ubm add' to add your first bookmark.")
				return nil
			}

			// Start interactive navigation
			if err := ui.NavigateBookmarks(categoryTree, data.Bookmarks); err != nil {
				return handleCancelError(err)
			}

			return nil
		},
	}

	return cmd
}