package bookmark

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	title := "Test Bookmark"
	url := "https://example.com"
	category := "test/category"

	b := New(title, url, category)

	if b == nil {
		t.Fatal("New() returned nil")
	}

	if b.Title != title {
		t.Errorf("Title = %v, want %v", b.Title, title)
	}

	if b.URL != url {
		t.Errorf("URL = %v, want %v", b.URL, url)
	}

	if b.Category != category {
		t.Errorf("Category = %v, want %v", b.Category, category)
	}

	if b.ID == "" {
		t.Error("ID should not be empty")
	}

	if b.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}

	if b.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should not be zero")
	}

	if !b.CreatedAt.Equal(b.UpdatedAt) {
		t.Error("CreatedAt and UpdatedAt should be equal for new bookmark")
	}

	if b.Tags == nil || len(b.Tags) != 0 {
		t.Error("Tags should be empty slice, not nil")
	}

	if b.Description != "" {
		t.Error("Description should be empty for new bookmark")
	}
}

func TestBookmark_Update(t *testing.T) {
	b := New("Test", "https://example.com", "test")
	originalUpdatedAt := b.UpdatedAt

	// Sleep to ensure time difference
	time.Sleep(10 * time.Millisecond)

	b.Update()

	if !b.UpdatedAt.After(originalUpdatedAt) {
		t.Error("UpdatedAt should be updated to a later time")
	}
}

func TestBookmark_AddTag(t *testing.T) {
	b := New("Test", "https://example.com", "test")
	originalUpdatedAt := b.UpdatedAt

	time.Sleep(10 * time.Millisecond)

	// Add first tag
	b.AddTag("tag1")
	if len(b.Tags) != 1 || b.Tags[0] != "tag1" {
		t.Errorf("Tags after first add = %v, want [tag1]", b.Tags)
	}
	if !b.UpdatedAt.After(originalUpdatedAt) {
		t.Error("UpdatedAt should be updated after adding tag")
	}

	// Add second tag
	b.AddTag("tag2")
	if len(b.Tags) != 2 || b.Tags[0] != "tag1" || b.Tags[1] != "tag2" {
		t.Errorf("Tags after second add = %v, want [tag1 tag2]", b.Tags)
	}

	// Try to add duplicate tag
	b.AddTag("tag1")
	if len(b.Tags) != 2 {
		t.Errorf("Tags should not contain duplicates, got %v", b.Tags)
	}
}

func TestBookmark_RemoveTag(t *testing.T) {
	b := New("Test", "https://example.com", "test")
	b.Tags = []string{"tag1", "tag2", "tag3"}
	originalUpdatedAt := b.UpdatedAt

	time.Sleep(10 * time.Millisecond)

	// Remove middle tag
	b.RemoveTag("tag2")
	if len(b.Tags) != 2 || b.Tags[0] != "tag1" || b.Tags[1] != "tag3" {
		t.Errorf("Tags after removing tag2 = %v, want [tag1 tag3]", b.Tags)
	}
	if !b.UpdatedAt.After(originalUpdatedAt) {
		t.Error("UpdatedAt should be updated after removing tag")
	}

	// Remove non-existent tag
	b.RemoveTag("tag4")
	if len(b.Tags) != 2 {
		t.Errorf("Tags should not change when removing non-existent tag, got %v", b.Tags)
	}

	// Remove all tags
	b.RemoveTag("tag1")
	b.RemoveTag("tag3")
	if len(b.Tags) != 0 {
		t.Errorf("Tags should be empty after removing all, got %v", b.Tags)
	}
}

func TestBookmark_SetCategory(t *testing.T) {
	b := New("Test", "https://example.com", "test")
	originalUpdatedAt := b.UpdatedAt

	time.Sleep(10 * time.Millisecond)

	newCategory := "new/category"
	b.SetCategory(newCategory)

	if b.Category != newCategory {
		t.Errorf("Category = %v, want %v", b.Category, newCategory)
	}
	if !b.UpdatedAt.After(originalUpdatedAt) {
		t.Error("UpdatedAt should be updated after setting category")
	}
}

func TestBookmark_SetTitle(t *testing.T) {
	b := New("Test", "https://example.com", "test")
	originalUpdatedAt := b.UpdatedAt

	time.Sleep(10 * time.Millisecond)

	newTitle := "New Title"
	b.SetTitle(newTitle)

	if b.Title != newTitle {
		t.Errorf("Title = %v, want %v", b.Title, newTitle)
	}
	if !b.UpdatedAt.After(originalUpdatedAt) {
		t.Error("UpdatedAt should be updated after setting title")
	}
}

func TestBookmark_SetURL(t *testing.T) {
	b := New("Test", "https://example.com", "test")
	originalUpdatedAt := b.UpdatedAt

	time.Sleep(10 * time.Millisecond)

	newURL := "https://newexample.com"
	b.SetURL(newURL)

	if b.URL != newURL {
		t.Errorf("URL = %v, want %v", b.URL, newURL)
	}
	if !b.UpdatedAt.After(originalUpdatedAt) {
		t.Error("UpdatedAt should be updated after setting URL")
	}
}

func TestBookmark_SetDescription(t *testing.T) {
	b := New("Test", "https://example.com", "test")
	originalUpdatedAt := b.UpdatedAt

	time.Sleep(10 * time.Millisecond)

	newDescription := "This is a test bookmark"
	b.SetDescription(newDescription)

	if b.Description != newDescription {
		t.Errorf("Description = %v, want %v", b.Description, newDescription)
	}
	if !b.UpdatedAt.After(originalUpdatedAt) {
		t.Error("UpdatedAt should be updated after setting description")
	}
}

func TestBookmark_MultipleUpdates(t *testing.T) {
	b := New("Test", "https://example.com", "test")

	// Perform multiple updates
	b.SetTitle("Updated Title")
	b.AddTag("important")
	b.SetCategory("updated/category")
	b.SetDescription("Updated description")

	// Verify all changes
	if b.Title != "Updated Title" {
		t.Errorf("Title = %v, want Updated Title", b.Title)
	}
	if len(b.Tags) != 1 || b.Tags[0] != "important" {
		t.Errorf("Tags = %v, want [important]", b.Tags)
	}
	if b.Category != "updated/category" {
		t.Errorf("Category = %v, want updated/category", b.Category)
	}
	if b.Description != "Updated description" {
		t.Errorf("Description = %v, want Updated description", b.Description)
	}

	// Verify UpdatedAt is after CreatedAt
	if !b.UpdatedAt.After(b.CreatedAt) {
		t.Error("UpdatedAt should be after CreatedAt after updates")
	}
}