package main

import (
	"fmt"
	"sort"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/tom-023/ubm/internal/category"
	"github.com/tom-023/ubm/internal/ui"
)

func categoryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "category",
		Short: "Manage categories",
		Long:  `Manage bookmark categories including creating, listing, and deleting categories.`,
	}

	cmd.AddCommand(
		categoryCreateCmd(),
		categoryListCmd(),
		categoryDeleteCmd(),
	)

	return cmd
}

func categoryCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "Create new category interactively",
		Long:  `Create a new category by selecting a parent category and providing a name.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load existing data
			data, err := store.Load()
			if err != nil {
				return fmt.Errorf("failed to load data: %w", err)
			}

			// Build category tree
			catManager := category.NewManager()
			bookmarkCounts := make(map[string]int)
			for _, b := range data.Bookmarks {
				bookmarkCounts[b.Category]++
			}
			categoryTree := catManager.BuildTree(data.Categories, bookmarkCounts)

			// Select parent category
			fmt.Println("Select parent category (or press Enter for root level):")
			parentPath, err := ui.SelectCategory(categoryTree, "")
			if err != nil {
				return fmt.Errorf("failed to select parent category: %w", err)
			}

			// If user selected "Create new category" at root, parentPath will be the new category name
			var newCategoryPath string
			if parentPath != "" && !categoryExists(parentPath, data.Categories) {
				// User created a root-level category
				newCategoryPath = parentPath
			} else {
				// Get new category name
				name, err := ui.PromptString("New category name", "")
				if err != nil {
					return fmt.Errorf("failed to get category name: %w", err)
				}

				// Create full path
				newCategoryPath, err = catManager.CreateCategory(parentPath, name)
				if err != nil {
					return fmt.Errorf("failed to create category: %w", err)
				}
			}

			// Check if category already exists
			for _, cat := range data.Categories {
				if cat == newCategoryPath {
					return fmt.Errorf("category '%s' already exists", newCategoryPath)
				}
			}

			// Add to categories
			data.Categories = append(data.Categories, newCategoryPath)
			sort.Strings(data.Categories)

			// Save data
			if err := store.Save(data); err != nil {
				return fmt.Errorf("failed to save data: %w", err)
			}

			fmt.Printf("✅ Category '%s' created successfully!\n", newCategoryPath)
			return nil
		},
	}
}

func categoryListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all categories in tree format",
		Long:  `Display all categories in a hierarchical tree structure.`,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load data
			data, err := store.Load()
			if err != nil {
				return fmt.Errorf("failed to load data: %w", err)
			}

			if len(data.Categories) == 0 {
				fmt.Println("No categories found. Use 'ubm category create' to create your first category.")
				return nil
			}

			// Build and display category tree
			catManager := category.NewManager()
			bookmarkCounts := make(map[string]int)
			for _, b := range data.Bookmarks {
				bookmarkCounts[b.Category]++
			}
			tree := catManager.BuildTree(data.Categories, bookmarkCounts)

			fmt.Println("📁 Categories:")
			printCategoryNode(tree, "", true)

			return nil
		},
	}
}

func categoryDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete",
		Short: "Delete empty category",
		Long:  `Delete a category. Only empty categories (with no bookmarks) can be deleted.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load data
			data, err := store.Load()
			if err != nil {
				return fmt.Errorf("failed to load data: %w", err)
			}

			if len(data.Categories) == 0 {
				fmt.Println("No categories found.")
				return nil
			}

			// Build category tree
			catManager := category.NewManager()
			bookmarkCounts := make(map[string]int)
			for _, b := range data.Bookmarks {
				bookmarkCounts[b.Category]++
			}
			_ = catManager.BuildTree(data.Categories, bookmarkCounts)

			// Filter out non-empty categories
			emptyCategories := []string{}
			for _, cat := range data.Categories {
				if bookmarkCounts[cat] == 0 {
					// Also check if it has subcategories
					hasSubcategories := false
					for _, otherCat := range data.Categories {
						if otherCat != cat && len(otherCat) > len(cat) && 
							otherCat[:len(cat)] == cat && otherCat[len(cat)] == '/' {
							hasSubcategories = true
							break
						}
					}
					if !hasSubcategories {
						emptyCategories = append(emptyCategories, cat)
					}
				}
			}

			if len(emptyCategories) == 0 {
				fmt.Println("No empty categories found. Only empty categories can be deleted.")
				return nil
			}

			// Select category to delete
			fmt.Println("Select category to delete (only empty categories are shown):")
			selectedCategory, err := selectEmptyCategory(emptyCategories)
			if err != nil {
				return fmt.Errorf("failed to select category: %w", err)
			}

			// Confirm deletion
			confirm, err := ui.Confirm(fmt.Sprintf("Delete category '%s'?", selectedCategory))
			if err != nil {
				return err
			}
			if !confirm {
				fmt.Println("Deletion cancelled.")
				return nil
			}

			// Remove category
			newCategories := []string{}
			for _, cat := range data.Categories {
				if cat != selectedCategory {
					newCategories = append(newCategories, cat)
				}
			}
			data.Categories = newCategories

			// Save data
			if err := store.Save(data); err != nil {
				return fmt.Errorf("failed to save data: %w", err)
			}

			fmt.Printf("✅ Category '%s' deleted successfully!\n", selectedCategory)
			return nil
		},
	}
}

func printCategoryNode(node *category.Node, prefix string, isLast bool) {
	if !node.IsRoot {
		connector := "├── "
		if isLast {
			connector = "└── "
		}
		fmt.Printf("%s%s%s", prefix, connector, node.Name)
		if node.Count > 0 {
			fmt.Printf(" (%d bookmarks)", node.Count)
		}
		fmt.Println()
	}

	childPrefix := prefix
	if !node.IsRoot {
		if isLast {
			childPrefix += "    "
		} else {
			childPrefix += "│   "
		}
	}

	for i, child := range node.Children {
		isLastChild := i == len(node.Children)-1
		printCategoryNode(child, childPrefix, isLastChild)
	}
}

func categoryExists(path string, categories []string) bool {
	for _, cat := range categories {
		if cat == path {
			return true
		}
	}
	return false
}

func selectEmptyCategory(categories []string) (string, error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "▶ {{ . }}",
		Inactive: "  {{ . }}",
		Selected: "{{ . }}",
	}

	prompt := promptui.Select{
		Label:     "Select category to delete",
		Items:     categories,
		Templates: templates,
		Size:      10,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return categories[i], nil
}