package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/tom-023/ubm/internal/bookmark"
)

type Storage struct {
	filePath   string
	backupPath string
	mu         sync.RWMutex
}

type Data struct {
	Bookmarks  []*bookmark.Bookmark `json:"bookmarks"`
	Categories []string            `json:"categories"`
	UpdatedAt  time.Time           `json:"updated_at"`
}

func New(configDir string) (*Storage, error) {
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	return &Storage{
		filePath:   filepath.Join(configDir, "bookmarks.json"),
		backupPath: filepath.Join(configDir, "bookmarks.backup.json"),
	}, nil
}

func (s *Storage) Load() (*Data, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	file, err := os.Open(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Data{
				Bookmarks:  []*bookmark.Bookmark{},
				Categories: []string{},
				UpdatedAt:  time.Now(),
			}, nil
		}
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var data Data
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode data: %w", err)
	}

	return &data, nil
}

func (s *Storage) Save(data *Data) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data.UpdatedAt = time.Now()

	// Create backup if original file exists
	if _, err := os.Stat(s.filePath); err == nil {
		if err := s.createBackup(); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
	}

	// Marshal data
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	// Write to temporary file first
	tmpFile := s.filePath + ".tmp"
	if err := os.WriteFile(tmpFile, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write temporary file: %w", err)
	}

	// Rename temporary file to actual file
	if err := os.Rename(tmpFile, s.filePath); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}

func (s *Storage) createBackup() error {
	src, err := os.Open(s.filePath)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(s.backupPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = dst.ReadFrom(src)
	return err
}

func (s *Storage) AddBookmark(b *bookmark.Bookmark) error {
	data, err := s.Load()
	if err != nil {
		return err
	}

	// Check for duplicate URL in the same category
	for _, existing := range data.Bookmarks {
		if existing.URL == b.URL && existing.Category == b.Category {
			return fmt.Errorf("bookmark with URL %s already exists in category %s", b.URL, b.Category)
		}
	}

	data.Bookmarks = append(data.Bookmarks, b)
	
	// Add category if it doesn't exist
	categoryExists := false
	for _, cat := range data.Categories {
		if cat == b.Category {
			categoryExists = true
			break
		}
	}
	if !categoryExists && b.Category != "" {
		data.Categories = append(data.Categories, b.Category)
	}

	return s.Save(data)
}

func (s *Storage) UpdateBookmark(b *bookmark.Bookmark) error {
	data, err := s.Load()
	if err != nil {
		return err
	}

	found := false
	for i, existing := range data.Bookmarks {
		if existing.ID == b.ID {
			data.Bookmarks[i] = b
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("bookmark with ID %s not found", b.ID)
	}

	return s.Save(data)
}

func (s *Storage) DeleteBookmark(id string) error {
	data, err := s.Load()
	if err != nil {
		return err
	}

	bookmarks := []*bookmark.Bookmark{}
	found := false
	for _, b := range data.Bookmarks {
		if b.ID != id {
			bookmarks = append(bookmarks, b)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("bookmark with ID %s not found", id)
	}

	data.Bookmarks = bookmarks
	return s.Save(data)
}

func (s *Storage) GetBookmark(id string) (*bookmark.Bookmark, error) {
	data, err := s.Load()
	if err != nil {
		return nil, err
	}

	for _, b := range data.Bookmarks {
		if b.ID == id {
			return b, nil
		}
	}

	return nil, fmt.Errorf("bookmark with ID %s not found", id)
}

func (s *Storage) GetBookmarksByCategory(category string) ([]*bookmark.Bookmark, error) {
	data, err := s.Load()
	if err != nil {
		return nil, err
	}

	bookmarks := []*bookmark.Bookmark{}
	for _, b := range data.Bookmarks {
		if b.Category == category {
			bookmarks = append(bookmarks, b)
		}
	}

	return bookmarks, nil
}

func (s *Storage) SearchBookmarks(query string) ([]*bookmark.Bookmark, error) {
	data, err := s.Load()
	if err != nil {
		return nil, err
	}

	bookmarks := []*bookmark.Bookmark{}
	for _, b := range data.Bookmarks {
		if containsIgnoreCase(b.Title, query) || containsIgnoreCase(b.URL, query) || containsIgnoreCase(b.Description, query) {
			bookmarks = append(bookmarks, b)
		}
	}

	return bookmarks, nil
}

func containsIgnoreCase(s, substr string) bool {
	s = string([]rune(s))
	substr = string([]rune(substr))
	
	if len(substr) > len(s) {
		return false
	}
	
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] && s[i+j] != substr[j]+32 && s[i+j] != substr[j]-32 {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	
	return false
}