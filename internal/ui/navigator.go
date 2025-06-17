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

// BookmarkAction defines what to do when a bookmark is selected
type BookmarkAction func(*bookmark.Bookmark) error

// navigateWithAction is the common navigation function
func navigateWithAction(categoryTree *category.Node, bookmarks []*bookmark.Bookmark, label string, action BookmarkAction) error {
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

		// Use DetailedSelectTemplates but override Details for bookmark display
		templates := &promptui.SelectTemplates{
			Label:    DetailedSelectTemplates.Label,
			Active:   DetailedSelectTemplates.Active,
			Inactive: DetailedSelectTemplates.Inactive,
			Selected: DetailedSelectTemplates.Selected,
			Details: `
{{ if eq .Type "bookmark" }}{{ if .Bookmark }}
{{ "--------- Bookmark Details ----------" | faint }}
{{ "Title:" | yellow }} {{ .Bookmark.Title | white }}
{{ "URL:" | yellow }}   {{ .Bookmark.URL | white }}
{{ if .Bookmark.Description }}{{ "Description:" | yellow }} {{ .Bookmark.Description | white }}{{ end }}
{{ end }}{{ end }}`,
		}

		searcher := CreateSearcher(func(index int) string {
			item := items[index]
			searchText := item.Display
			if item.Bookmark != nil {
				searchText += " " + item.Bookmark.URL
			}
			return searchText
		})

		promptLabel := label
		if promptLabel == "" {
			promptLabel = formatNavigationPath(path)
		} else {
			promptLabel = fmt.Sprintf("%s - %s", label, formatNavigationPath(path))
		}

		prompt := promptui.Select{
			Label:     promptLabel,
			Items:     items,
			Templates: templates,
			Searcher:  searcher,
			Size:      15,
		}

		i, _, err := prompt.Run()
		if err != nil {
			return WrapCancelError(err)
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
			// Execute the action on the selected bookmark
			if action != nil {
				return action(selected.Bookmark)
			}
			return nil
		}

		return nil
	}

	return navigateRecursive(categoryTree, "")
}

// NavigateBookmarks opens the selected bookmark in browser
func NavigateBookmarks(categoryTree *category.Node, bookmarks []*bookmark.Bookmark) error {
	return navigateWithAction(categoryTree, bookmarks, "", func(b *bookmark.Bookmark) error {
		fmt.Printf("\nOpening: %s\n", b.URL)
		if err := browser.OpenURL(b.URL); err != nil {
			fmt.Printf("Error opening browser: %v\n", err)
			fmt.Printf("Please open manually: %s\n", b.URL)
		}
		return nil
	})
}

func formatNavigationPath(path string) string {
	if path == "" {
		return "📚 Bookmarks"
	}
	return fmt.Sprintf("📚 Bookmarks > %s", strings.ReplaceAll(path, "/", " > "))
}

// NavigateAndSelectBookmark allows navigating through categories to select a bookmark
func NavigateAndSelectBookmark(categoryTree *category.Node, bookmarks []*bookmark.Bookmark, prompt string) (*bookmark.Bookmark, error) {
	var selectedBookmark *bookmark.Bookmark
	err := navigateWithAction(categoryTree, bookmarks, prompt, func(b *bookmark.Bookmark) error {
		selectedBookmark = b
		return nil
	})
	if err != nil {
		return nil, err
	}
	return selectedBookmark, nil
}
