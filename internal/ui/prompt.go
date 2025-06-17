package ui

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/tom-023/ubm/internal/bookmark"
	"github.com/tom-023/ubm/internal/category"
)

func PromptString(label string, defaultValue string) (string, error) {
	prompt := promptui.Prompt{
		Label:     label,
		Default:   defaultValue,
		Templates: StandardPromptTemplates,
	}
	result, err := prompt.Run()
	if err != nil {
		return "", WrapCancelError(err)
	}
	return result, nil
}

func PromptURL(defaultValue string) (string, error) {
	return PromptURLWithLabel("URL", defaultValue)
}

func PromptURLWithLabel(label string, defaultValue string) (string, error) {
	prompt := promptui.Prompt{
		Label:   label,
		Default: defaultValue,
		Validate: func(input string) error {
			if input == "" {
				return fmt.Errorf("URL cannot be empty")
			}
			return nil
		},
		Templates: StandardPromptTemplates,
	}
	result, err := prompt.Run()
	if err != nil {
		return "", WrapCancelError(err)
	}
	return result, nil
}

func SelectCategory(categoryTree *category.Node, currentPath string) (string, error) {
	type categoryItem struct {
		Display string
		Path    string
		IsBack  bool
		IsNew   bool
	}

	var selectCategoryRecursive func(node *category.Node, parentPath string) (string, error)
	selectCategoryRecursive = func(node *category.Node, parentPath string) (string, error) {
		items := []categoryItem{}

		// Add back option if not at root
		if parentPath != "" {
			items = append(items, categoryItem{
				Display: "⬅️  Back to parent",
				Path:    "",
				IsBack:  true,
			})
		}

		// Add child categories
		for _, child := range node.Children {
			display := fmt.Sprintf("📁 %s", child.Name)
			if child.Count > 0 {
				display = fmt.Sprintf("📁 %s (%d bookmarks)", child.Name, child.Count)
			}
			items = append(items, categoryItem{
				Display: display,
				Path:    child.Path,
			})
		}

		// Add option to create new category
		items = append(items, categoryItem{
			Display: "➕ Create new category",
			Path:    "",
			IsNew:   true,
		})

		// Add option to select current directory if not at root
		if parentPath != "" {
			items = append(items, categoryItem{
				Display: fmt.Sprintf("✅ Select this category (%s)", parentPath),
				Path:    parentPath,
			})
		}

		searcher := CreateSearcher(func(index int) string {
			return items[index].Display
		})


		prompt := promptui.Select{
			Label:     fmt.Sprintf("Select Category (current: %s)", formatPath(parentPath)),
			Items:     items,
			Templates: ItemSelectTemplates,
			Searcher:  searcher,
			Size:      GetSelectSize(len(items)),
			HideHelp:  true,
		}

		i, _, err := prompt.Run()
		if err != nil {
			if IsCancelError(err) {
				return "", ErrCancelled
			}
			return "", err
		}

		selected := items[i]

		if selected.IsBack {
			// Go back to parent
			parent := findParentNode(categoryTree, parentPath)
			if parent != nil {
				return selectCategoryRecursive(parent, category.NewManager().GetParentPath(parentPath))
			}
			return selectCategoryRecursive(categoryTree, "")
		}

		if selected.IsNew {
			// Create new category
			name, err := PromptString("New category name", "")
			if err != nil {
				return "", err
			}
			if parentPath == "" {
				return name, nil
			}
			return fmt.Sprintf("%s/%s", parentPath, name), nil
		}

		if strings.HasPrefix(selected.Display, "✅") {
			// Selected current directory
			return selected.Path, nil
		}

		// Navigate into subdirectory
		childNode := findNode(categoryTree, selected.Path)
		if childNode != nil && len(childNode.Children) > 0 {
			return selectCategoryRecursive(childNode, selected.Path)
		}

		// Leaf node selected
		return selected.Path, nil
	}

	return selectCategoryRecursive(categoryTree, currentPath)
}

func findNode(root *category.Node, path string) *category.Node {
	if root.Path == path {
		return root
	}
	for _, child := range root.Children {
		if node := findNode(child, path); node != nil {
			return node
		}
	}
	return nil
}

func findParentNode(root *category.Node, childPath string) *category.Node {
	for _, child := range root.Children {
		if child.Path == childPath {
			return root
		}
		if parent := findParentNode(child, childPath); parent != nil {
			return parent
		}
	}
	return nil
}

func formatPath(path string) string {
	if path == "" {
		return "/"
	}
	return path
}

func SelectBookmark(bookmarks []*bookmark.Bookmark, label string) (*bookmark.Bookmark, error) {
	if len(bookmarks) == 0 {
		return nil, fmt.Errorf("no bookmarks found")
	}

	type bookmarkItem struct {
		Display  string
		Bookmark *bookmark.Bookmark
	}

	items := []bookmarkItem{}
	for _, b := range bookmarks {
		display := fmt.Sprintf("%s - %s [%s]", b.Title, b.URL, b.Category)
		if b.Category == "" {
			display = fmt.Sprintf("%s - %s [uncategorized]", b.Title, b.URL)
		}
		items = append(items, bookmarkItem{
			Display:  display,
			Bookmark: b,
		})
	}


	searcher := CreateSearcher(func(index int) string {
		item := items[index]
		return item.Bookmark.Title + " " + item.Bookmark.URL + " " + item.Bookmark.Category
	})

	prompt := promptui.Select{
		Label:     label,
		Items:     items,
		Templates: ItemSelectTemplates,
		Searcher:  searcher,
		Size:      GetSelectSize(len(items)),
		HideHelp:  true,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	return items[i].Bookmark, nil
}

func Confirm(message string) (bool, error) {
	// Use Select instead of Prompt with IsConfirm for better color control
	items := []string{"Yes", "No"}
	prompt := promptui.Select{
		Label:     message,
		Items:     items,
		Templates: StandardSelectTemplates,
		Size:      2,
		HideHelp:  true,
	}

	i, _, err := prompt.Run()
	if err != nil {
		if IsCancelError(err) {
			return false, ErrCancelled
		}
		return false, err
	}

	return i == 0, nil
}

func SelectEditField() (string, error) {
	fields := []string{
		"Title",
		"URL",
	}

	prompt := promptui.Select{
		Label:     "What would you like to edit?",
		Items:     fields,
		Templates: StandardSelectTemplates,
		Size:      2,
		HideHelp:  true,
	}

	i, _, err := prompt.Run()
	if err != nil {
		if err == promptui.ErrInterrupt || err == promptui.ErrEOF {
			return "", fmt.Errorf("cancelled")
		}
		return "", err
	}

	return fields[i], nil
}