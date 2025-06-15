package category

import (
	"fmt"
	"path"
	"sort"
	"strings"
)

type Node struct {
	Name      string
	Path      string
	Children  []*Node
	Count     int
	IsRoot    bool
}

type Manager struct {
	root *Node
}

func NewManager() *Manager {
	return &Manager{
		root: &Node{
			Name:     "",
			Path:     "",
			Children: []*Node{},
			Count:    0,
			IsRoot:   true,
		},
	}
}

func (m *Manager) BuildTree(categories []string, bookmarkCounts map[string]int) *Node {
	m.root = &Node{
		Name:     "",
		Path:     "",
		Children: []*Node{},
		Count:    0,
		IsRoot:   true,
	}

	// Sort categories for consistent ordering
	sort.Strings(categories)

	for _, category := range categories {
		m.addCategory(category, bookmarkCounts[category])
	}

	// Add uncategorized if it has bookmarks
	if count, exists := bookmarkCounts[""]; exists && count > 0 {
		m.root.Children = append(m.root.Children, &Node{
			Name:     "uncategorized",
			Path:     "",
			Children: []*Node{},
			Count:    count,
		})
	}

	return m.root
}

func (m *Manager) addCategory(categoryPath string, count int) {
	if categoryPath == "" {
		return
	}

	parts := strings.Split(categoryPath, "/")
	currentNode := m.root
	currentPath := ""

	for i, part := range parts {
		if part == "" {
			continue
		}

		if currentPath == "" {
			currentPath = part
		} else {
			currentPath = currentPath + "/" + part
		}

		// Check if this part already exists
		found := false
		for _, child := range currentNode.Children {
			if child.Name == part {
				currentNode = child
				found = true
				// Update count for leaf nodes
				if i == len(parts)-1 {
					child.Count = count
				}
				break
			}
		}

		if !found {
			newNode := &Node{
				Name:     part,
				Path:     currentPath,
				Children: []*Node{},
				Count:    0,
			}
			// Set count only for leaf nodes
			if i == len(parts)-1 {
				newNode.Count = count
			}
			currentNode.Children = append(currentNode.Children, newNode)
			currentNode = newNode
		}
	}
}

func (m *Manager) GetCategories() []string {
	categories := []string{}
	m.collectCategories(m.root, &categories)
	return categories
}

func (m *Manager) collectCategories(node *Node, categories *[]string) {
	if !node.IsRoot && node.Path != "" {
		*categories = append(*categories, node.Path)
	}
	for _, child := range node.Children {
		m.collectCategories(child, categories)
	}
}

func (m *Manager) CreateCategory(parentPath, name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("category name cannot be empty")
	}

	// Validate name
	if strings.Contains(name, "/") {
		return "", fmt.Errorf("category name cannot contain '/'")
	}

	newPath := name
	if parentPath != "" {
		newPath = path.Join(parentPath, name)
	}

	// Check if category already exists
	categories := m.GetCategories()
	for _, cat := range categories {
		if cat == newPath {
			return "", fmt.Errorf("category %s already exists", newPath)
		}
	}

	return newPath, nil
}

func (m *Manager) ValidateCategory(categoryPath string) error {
	if categoryPath == "" {
		return nil // Empty category is valid (uncategorized)
	}

	parts := strings.Split(categoryPath, "/")
	for _, part := range parts {
		if part == "" {
			return fmt.Errorf("invalid category path: %s", categoryPath)
		}
	}

	return nil
}

func (m *Manager) GetParentPath(categoryPath string) string {
	if categoryPath == "" {
		return ""
	}

	parts := strings.Split(categoryPath, "/")
	if len(parts) <= 1 {
		return ""
	}

	return strings.Join(parts[:len(parts)-1], "/")
}

func (m *Manager) GetCategoryName(categoryPath string) string {
	if categoryPath == "" {
		return "uncategorized"
	}

	parts := strings.Split(categoryPath, "/")
	return parts[len(parts)-1]
}

func (m *Manager) CountBookmarksInCategory(bookmarks []*struct {
	Category string
}) map[string]int {
	counts := make(map[string]int)
	
	for _, b := range bookmarks {
		counts[b.Category]++
	}
	
	return counts
}