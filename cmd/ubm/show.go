package main

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/tom-023/ubm/internal/bookmark"
	"github.com/tom-023/ubm/internal/category"
)

func showCmd() *cobra.Command {
	var flat bool

	cmd := &cobra.Command{
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

			if flat {
				// Display flat list
				displayFlatList(data.Bookmarks)
			} else {
				// Display tree structure
				displayTree(data.Bookmarks, data.Categories)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&flat, "flat", "f", false, "Show flat list instead of tree")

	return cmd
}

func displayFlatList(bookmarks []*bookmark.Bookmark) {
	// Sort by category, then by title
	sort.Slice(bookmarks, func(i, j int) bool {
		if bookmarks[i].Category == bookmarks[j].Category {
			return bookmarks[i].Title < bookmarks[j].Title
		}
		return bookmarks[i].Category < bookmarks[j].Category
	})

	currentCategory := ""
	for _, b := range bookmarks {
		if b.Category != currentCategory {
			currentCategory = b.Category
			if currentCategory == "" {
				fmt.Println("\n📁 uncategorized:")
			} else {
				fmt.Printf("\n📁 %s:\n", currentCategory)
			}
		}
		fmt.Printf("  🔗 %s - %s\n", b.Title, b.URL)
	}
}

func displayTree(bookmarks []*bookmark.Bookmark, categories []string) {
	// Build category tree
	catManager := category.NewManager()
	bookmarkCounts := make(map[string]int)
	for _, b := range bookmarks {
		bookmarkCounts[b.Category]++
	}
	tree := catManager.BuildTree(categories, bookmarkCounts)

	// Group bookmarks by category
	bookmarksByCategory := make(map[string][]*bookmark.Bookmark)
	for _, b := range bookmarks {
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