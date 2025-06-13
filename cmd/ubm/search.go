package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/pkg/browser"
	"github.com/tom-023/ubm/internal/ui"
)

func searchCmd() *cobra.Command {
	var categoryFilter string

	cmd := &cobra.Command{
		Use:   "search [QUERY]",
		Short: "Search bookmarks by title or URL",
		Long:  `Search for bookmarks that match the given query in their title, URL, or description.`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := strings.Join(args, " ")

			// Search bookmarks
			results, err := store.SearchBookmarks(query)
			if err != nil {
				return fmt.Errorf("search failed: %w", err)
			}

			// Filter by category if specified
			if categoryFilter != "" {
				filtered := results[:0]
				for _, b := range results {
					if b.Category == categoryFilter {
						filtered = append(filtered, b)
					}
				}
				results = filtered
			}

			if len(results) == 0 {
				fmt.Printf("No bookmarks found matching '%s'", query)
				if categoryFilter != "" {
					fmt.Printf(" in category '%s'", categoryFilter)
				}
				fmt.Println()
				return nil
			}

			fmt.Printf("Found %d bookmark(s) matching '%s':\n\n", len(results), query)

			// Select bookmark to open
			selected, err := ui.SelectBookmark(results, "Select bookmark to open")
			if err != nil {
				if err.Error() == "^C" {
					return nil
				}
				return err
			}

			// Open selected bookmark
			fmt.Printf("\nOpening: %s\n", selected.URL)
			if err := browser.OpenURL(selected.URL); err != nil {
				fmt.Printf("Error opening browser: %v\n", err)
				fmt.Printf("Please open manually: %s\n", selected.URL)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&categoryFilter, "category", "c", "", "Search within specific category")

	return cmd
}