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
		Label:   label,
		Default: defaultValue,
	}
	return prompt.Run()
}

func PromptURL(defaultValue string) (string, error) {
	prompt := promptui.Prompt{
		Label:   "URL",
		Default: defaultValue,
		Validate: func(input string) error {
			if input == "" {
				return fmt.Errorf("URL cannot be empty")
			}
			return nil
		},
	}
	return prompt.Run()
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

		templates := &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   "▶ {{ .Display }}",
			Inactive: "  {{ .Display }}",
			Selected: "{{ .Display }}",
		}

		searcher := func(input string, index int) bool {
			item := items[index]
			name := strings.Replace(strings.ToLower(item.Display), " ", "", -1)
			input = strings.Replace(strings.ToLower(input), " ", "", -1)
			return strings.Contains(name, input)
		}

		prompt := promptui.Select{
			Label:     fmt.Sprintf("Select Category (current: %s)", formatPath(parentPath)),
			Items:     items,
			Templates: templates,
			Searcher:  searcher,
			Size:      10,
		}

		i, _, err := prompt.Run()
		if err != nil {
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

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "▶ {{ .Display }}",
		Inactive: "  {{ .Display }}",
		Selected: "{{ .Display }}",
	}

	searcher := func(input string, index int) bool {
		item := items[index]
		searchText := strings.ToLower(item.Bookmark.Title + " " + item.Bookmark.URL + " " + item.Bookmark.Category)
		input = strings.ToLower(input)
		return strings.Contains(searchText, input)
	}

	prompt := promptui.Select{
		Label:     label,
		Items:     items,
		Templates: templates,
		Searcher:  searcher,
		Size:      10,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	return items[i].Bookmark, nil
}

func Confirm(message string) (bool, error) {
	prompt := promptui.Prompt{
		Label:     message,
		IsConfirm: true,
	}
	_, err := prompt.Run()
	if err != nil {
		if err == promptui.ErrAbort {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func SelectEditField() (string, error) {
	fields := []string{
		"Title",
		"URL",
		"Category",
		"All fields",
	}

	prompt := promptui.Select{
		Label: "What would you like to edit?",
		Items: fields,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return fields[i], nil
}