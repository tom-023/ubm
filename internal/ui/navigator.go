package ui

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/pkg/browser"
	"github.com/tom-023/ubm/internal/bookmark"
	"github.com/tom-023/ubm/internal/category"
)

type NavigationItem struct {
	Type     string // "category", "bookmark", "back"
	Display  string
	Path     string
	Bookmark *bookmark.Bookmark
	Node     *category.Node
}

func NavigateBookmarks(categoryTree *category.Node, bookmarks []*bookmark.Bookmark) error {
	var navigateRecursive func(node *category.Node, path string) error
	navigateRecursive = func(node *category.Node, path string) error {
		items := []NavigationItem{}

		// Add back option if not at root
		if path != "" {
			items = append(items, NavigationItem{
				Type:    "back",
				Display: "⬅️  Back to parent",
			})
		}

		// Add subcategories
		for _, child := range node.Children {
			display := fmt.Sprintf("📁 %s", child.Name)
			if child.Count > 0 {
				display = fmt.Sprintf("📁 %s (%d)", child.Name, child.Count)
			}
			items = append(items, NavigationItem{
				Type:    "category",
				Display: display,
				Path:    child.Path,
				Node:    child,
			})
		}

		// Add bookmarks in current category
		for _, b := range bookmarks {
			if b.Category == path {
				items = append(items, NavigationItem{
					Type:     "bookmark",
					Display:  fmt.Sprintf("🔗 %s", b.Title),
					Bookmark: b,
				})
			}
		}

		if len(items) == 0 {
			fmt.Println("No bookmarks or categories found.")
			return nil
		}

		templates := &promptui.SelectTemplates{
			Label:    "{{ . | cyan | bold }}",
			Active:   "{{ \"▶\" | cyan | bold }} {{ .Display | cyan }}",
			Inactive: "  {{ .Display | faint }}",
			Selected: "{{ \"✔\" | green | bold }} {{ .Display | green }}",
			Details: `
{{ if eq .Type "bookmark" }}{{ if .Bookmark }}
{{ "--------- Bookmark Details ----------" | faint }}
{{ "Title:" | yellow }} {{ .Bookmark.Title | white }}
{{ "URL:" | yellow }}   {{ .Bookmark.URL | white }}
{{ if .Bookmark.Description }}{{ "Description:" | yellow }} {{ .Bookmark.Description | white }}{{ end }}
{{ end }}{{ end }}`,
		}

		searcher := func(input string, index int) bool {
			item := items[index]
			searchText := strings.ToLower(item.Display)
			if item.Bookmark != nil {
				searchText += " " + strings.ToLower(item.Bookmark.URL)
			}
			input = strings.ToLower(input)
			return strings.Contains(searchText, input)
		}

		prompt := promptui.Select{
			Label:     formatNavigationPath(path),
			Items:     items,
			Templates: templates,
			Searcher:  searcher,
			Size:      15,
		}

		i, _, err := prompt.Run()
		if err != nil {
			if err == promptui.ErrInterrupt {
				return nil
			}
			return err
		}

		selected := items[i]

		switch selected.Type {
		case "back":
			// Go back to parent
			parentPath := category.NewManager().GetParentPath(path)
			parentNode := findNode(categoryTree, parentPath)
			if parentNode != nil {
				return navigateRecursive(parentNode, parentPath)
			}
			return navigateRecursive(categoryTree, "")

		case "category":
			// Navigate into category
			return navigateRecursive(selected.Node, selected.Path)

		case "bookmark":
			// Open bookmark
			fmt.Printf("\nOpening: %s\n", selected.Bookmark.URL)
			if err := browser.OpenURL(selected.Bookmark.URL); err != nil {
				fmt.Printf("Error opening browser: %v\n", err)
				fmt.Printf("Please open manually: %s\n", selected.Bookmark.URL)
			}
			return nil // Exit after opening browser
		}

		return nil
	}

	return navigateRecursive(categoryTree, "")
}

func formatNavigationPath(path string) string {
	if path == "" {
		return "📚 Bookmarks"
	}
	return fmt.Sprintf("📚 Bookmarks > %s", strings.ReplaceAll(path, "/", " > "))
}