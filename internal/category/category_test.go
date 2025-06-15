package category

import (
	"reflect"
	"sort"
	"testing"
)

func TestNewManager(t *testing.T) {
	m := NewManager()
	
	if m == nil {
		t.Fatal("NewManager() returned nil")
	}
	
	if m.root == nil {
		t.Fatal("Manager root is nil")
	}
	
	if !m.root.IsRoot {
		t.Error("Root node should have IsRoot = true")
	}
	
	if m.root.Name != "" {
		t.Errorf("Root node Name = %v, want empty string", m.root.Name)
	}
	
	if m.root.Path != "" {
		t.Errorf("Root node Path = %v, want empty string", m.root.Path)
	}
	
	if len(m.root.Children) != 0 {
		t.Errorf("Root node should have no children initially")
	}
}

func TestManager_BuildTree(t *testing.T) {
	tests := []struct {
		name           string
		categories     []string
		bookmarkCounts map[string]int
		wantNodes      []string // Expected node paths in tree
		wantCounts     map[string]int
	}{
		{
			name:       "empty categories",
			categories: []string{},
			bookmarkCounts: map[string]int{},
			wantNodes:  []string{},
			wantCounts: map[string]int{},
		},
		{
			name:       "single category",
			categories: []string{"programming"},
			bookmarkCounts: map[string]int{
				"programming": 5,
			},
			wantNodes: []string{"programming"},
			wantCounts: map[string]int{
				"programming": 5,
			},
		},
		{
			name:       "nested categories",
			categories: []string{"programming", "programming/go", "programming/go/tutorials"},
			bookmarkCounts: map[string]int{
				"programming":              2,
				"programming/go":           3,
				"programming/go/tutorials": 4,
			},
			wantNodes: []string{"programming", "programming/go", "programming/go/tutorials"},
			wantCounts: map[string]int{
				"programming":              2,
				"programming/go":           3,
				"programming/go/tutorials": 4,
			},
		},
		{
			name:       "multiple root categories",
			categories: []string{"programming", "design", "tools"},
			bookmarkCounts: map[string]int{
				"programming": 5,
				"design":      3,
				"tools":       2,
			},
			wantNodes: []string{"design", "programming", "tools"}, // Sorted
			wantCounts: map[string]int{
				"programming": 5,
				"design":      3,
				"tools":       2,
			},
		},
		{
			name:       "with uncategorized",
			categories: []string{"programming"},
			bookmarkCounts: map[string]int{
				"programming": 5,
				"":            3, // uncategorized
			},
			wantNodes: []string{"programming", ""}, // uncategorized has empty path
			wantCounts: map[string]int{
				"programming": 5,
				"":            3,
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewManager()
			root := m.BuildTree(tt.categories, tt.bookmarkCounts)
			
			// Collect all nodes from tree
			nodes := collectAllNodes(root)
			nodePaths := []string{}
			for _, node := range nodes {
				if !node.IsRoot {
					nodePaths = append(nodePaths, node.Path)
				}
			}
			
			// Sort for comparison
			sort.Strings(nodePaths)
			sort.Strings(tt.wantNodes)
			
			if !reflect.DeepEqual(nodePaths, tt.wantNodes) {
				t.Errorf("BuildTree() nodes = %v, want %v", nodePaths, tt.wantNodes)
			}
			
			// Check counts
			for _, node := range nodes {
				// Skip root node
				if node.IsRoot {
					continue
				}
				
				// Special handling for uncategorized node
				nodeKey := node.Path
				if node.Name == "uncategorized" && node.Path == "" {
					nodeKey = ""
				}
				
				if expectedCount, exists := tt.wantCounts[nodeKey]; exists {
					if node.Count != expectedCount {
						t.Errorf("Node path=%q name=%q count = %d, want %d", node.Path, node.Name, node.Count, expectedCount)
					}
				}
			}
		})
	}
}

func TestManager_GetCategories(t *testing.T) {
	m := NewManager()
	categories := []string{"programming", "programming/go", "design", "tools/cli"}
	counts := map[string]int{
		"programming":    1,
		"programming/go": 2,
		"design":         3,
		"tools/cli":      4,
	}
	
	m.BuildTree(categories, counts)
	got := m.GetCategories()
	
	// Sort for comparison
	sort.Strings(got)
	expected := []string{"design", "programming", "programming/go", "tools", "tools/cli"}
	
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("GetCategories() = %v, want %v", got, expected)
	}
}

func TestManager_CreateCategory(t *testing.T) {
	tests := []struct {
		name         string
		parentPath   string
		categoryName string
		existing     []string
		want         string
		wantErr      bool
		errMsg       string
	}{
		{
			name:         "create root category",
			parentPath:   "",
			categoryName: "programming",
			existing:     []string{},
			want:         "programming",
			wantErr:      false,
		},
		{
			name:         "create nested category",
			parentPath:   "programming",
			categoryName: "go",
			existing:     []string{"programming"},
			want:         "programming/go",
			wantErr:      false,
		},
		{
			name:         "empty category name",
			parentPath:   "",
			categoryName: "",
			existing:     []string{},
			want:         "",
			wantErr:      true,
			errMsg:       "category name cannot be empty",
		},
		{
			name:         "category name with slash",
			parentPath:   "",
			categoryName: "prog/ramming",
			existing:     []string{},
			want:         "",
			wantErr:      true,
			errMsg:       "category name cannot contain '/'",
		},
		{
			name:         "duplicate category",
			parentPath:   "",
			categoryName: "programming",
			existing:     []string{"programming"},
			want:         "",
			wantErr:      true,
			errMsg:       "category programming already exists",
		},
		{
			name:         "create deeply nested",
			parentPath:   "programming/go",
			categoryName: "tutorials",
			existing:     []string{"programming", "programming/go"},
			want:         "programming/go/tutorials",
			wantErr:      false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewManager()
			// Build tree with existing categories
			counts := make(map[string]int)
			for _, cat := range tt.existing {
				counts[cat] = 1
			}
			m.BuildTree(tt.existing, counts)
			
			got, err := m.CreateCategory(tt.parentPath, tt.categoryName)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.errMsg {
				t.Errorf("CreateCategory() error message = %v, want %v", err.Error(), tt.errMsg)
			}
			if got != tt.want {
				t.Errorf("CreateCategory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_ValidateCategory(t *testing.T) {
	tests := []struct {
		name     string
		category string
		wantErr  bool
	}{
		{
			name:     "valid root category",
			category: "programming",
			wantErr:  false,
		},
		{
			name:     "valid nested category",
			category: "programming/go/tutorials",
			wantErr:  false,
		},
		{
			name:     "empty category (uncategorized)",
			category: "",
			wantErr:  false,
		},
		{
			name:     "category with empty part",
			category: "programming//go",
			wantErr:  true,
		},
		{
			name:     "category ending with slash",
			category: "programming/",
			wantErr:  true,
		},
		{
			name:     "category starting with slash",
			category: "/programming",
			wantErr:  true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewManager()
			err := m.ValidateCategory(tt.category)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCategory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_GetParentPath(t *testing.T) {
	tests := []struct {
		name     string
		category string
		want     string
	}{
		{
			name:     "root category",
			category: "programming",
			want:     "",
		},
		{
			name:     "nested category",
			category: "programming/go",
			want:     "programming",
		},
		{
			name:     "deeply nested category",
			category: "programming/go/tutorials",
			want:     "programming/go",
		},
		{
			name:     "empty category",
			category: "",
			want:     "",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewManager()
			got := m.GetParentPath(tt.category)
			if got != tt.want {
				t.Errorf("GetParentPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_GetCategoryName(t *testing.T) {
	tests := []struct {
		name     string
		category string
		want     string
	}{
		{
			name:     "root category",
			category: "programming",
			want:     "programming",
		},
		{
			name:     "nested category",
			category: "programming/go",
			want:     "go",
		},
		{
			name:     "deeply nested category",
			category: "programming/go/tutorials",
			want:     "tutorials",
		},
		{
			name:     "empty category",
			category: "",
			want:     "uncategorized",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewManager()
			got := m.GetCategoryName(tt.category)
			if got != tt.want {
				t.Errorf("GetCategoryName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_CountBookmarksInCategory(t *testing.T) {
	bookmarks := []*struct{ Category string }{
		{Category: "programming"},
		{Category: "programming"},
		{Category: "programming/go"},
		{Category: "design"},
		{Category: ""},
		{Category: ""},
		{Category: ""},
	}
	
	m := NewManager()
	got := m.CountBookmarksInCategory(bookmarks)
	
	want := map[string]int{
		"programming":    2,
		"programming/go": 1,
		"design":         1,
		"":               3,
	}
	
	if !reflect.DeepEqual(got, want) {
		t.Errorf("CountBookmarksInCategory() = %v, want %v", got, want)
	}
}

// Helper function to collect all nodes in the tree
func collectAllNodes(node *Node) []*Node {
	nodes := []*Node{node}
	for _, child := range node.Children {
		nodes = append(nodes, collectAllNodes(child)...)
	}
	return nodes
}