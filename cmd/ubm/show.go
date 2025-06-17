package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tom-023/ubm/internal/bookmark"
	"github.com/tom-023/ubm/internal/category"
	"github.com/tom-023/ubm/internal/storage"
	"github.com/tom-023/ubm/internal/ui"
)

func showCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Display all bookmarks in tree format",
		Long:  `Display all bookmarks organized by their categories in a tree structure.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load bookmarks
			data, err := store.Load()
			if err != nil {
				return fmt.Errorf("failed to load bookmarks: %w", err)
			}

			if len(data.Bookmarks) == 0 {
				fmt.Println("No bookmarks found. Use 'ubm add' to add your first bookmark.")
				return nil
			}

			// Display tree structure
			displayTree(data)

			return nil
		},
	}
}

func displayTree(data *storage.Data) {
	// Build category tree
	tree := ui.BuildCategoryTree(data)

	// Group bookmarks by category
	bookmarksByCategory := make(map[string][]*bookmark.Bookmark)
	for _, b := range data.Bookmarks {
		bookmarksByCategory[b.Category] = append(bookmarksByCategory[b.Category], b)
	}

	fmt.Println("📚 Bookmarks:")
	printNode(tree, "", true, bookmarksByCategory)
}

func printNode(node *category.Node, prefix string, isLast bool, bookmarksByCategory map[string][]*bookmark.Bookmark) {
	if !node.IsRoot {
		connector := "├── "
		if isLast {
			connector = "└── "
		}
		fmt.Printf("%s%s📁 %s", prefix, connector, node.Name)
		if node.Count > 0 {
			fmt.Printf(" (%d)", node.Count)
		}
		fmt.Println()

		// Print bookmarks in this category
		if bookmarks, exists := bookmarksByCategory[node.Path]; exists {
			childPrefix := prefix
			if isLast {
				childPrefix += "    "
			} else {
				childPrefix += "│   "
			}
			for i, b := range bookmarks {
				bookmarkConnector := "├── "
				if i == len(bookmarks)-1 && len(node.Children) == 0 {
					bookmarkConnector = "└── "
				}
				fmt.Printf("%s%s🔗 %s\n", childPrefix, bookmarkConnector, b.Title)
			}
		}
	}

	// Update prefix for children
	childPrefix := prefix
	if !node.IsRoot {
		if isLast {
			childPrefix += "    "
		} else {
			childPrefix += "│   "
		}
	}

	// Print children
	for i, child := range node.Children {
		isLastChild := i == len(node.Children)-1
		printNode(child, childPrefix, isLastChild, bookmarksByCategory)
	}

	// Print uncategorized bookmarks at root level
	if node.IsRoot {
		if bookmarks, exists := bookmarksByCategory[""]; exists && len(bookmarks) > 0 {
			fmt.Printf("%s└── 📁 uncategorized (%d)\n", prefix, len(bookmarks))
			for i, b := range bookmarks {
				bookmarkConnector := "├── "
				if i == len(bookmarks)-1 {
					bookmarkConnector = "└── "
				}
				fmt.Printf("%s    %s🔗 %s\n", prefix, bookmarkConnector, b.Title)
			}
		}
	}
}