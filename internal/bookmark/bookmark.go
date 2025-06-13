package bookmark

import (
	"time"

	"github.com/google/uuid"
)

type Bookmark struct {
	ID          string    `json:"id" yaml:"id"`
	Title       string    `json:"title" yaml:"title"`
	URL         string    `json:"url" yaml:"url"`
	Category    string    `json:"category" yaml:"category"`
	CreatedAt   time.Time `json:"created_at" yaml:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" yaml:"updated_at"`
	Tags        []string  `json:"tags,omitempty" yaml:"tags,omitempty"`
	Description string    `json:"description,omitempty" yaml:"description,omitempty"`
}

func New(title, url, category string) *Bookmark {
	now := time.Now()
	return &Bookmark{
		ID:        uuid.New().String(),
		Title:     title,
		URL:       url,
		Category:  category,
		CreatedAt: now,
		UpdatedAt: now,
		Tags:      []string{},
	}
}

func (b *Bookmark) Update() {
	b.UpdatedAt = time.Now()
}

func (b *Bookmark) AddTag(tag string) {
	for _, t := range b.Tags {
		if t == tag {
			return
		}
	}
	b.Tags = append(b.Tags, tag)
	b.Update()
}

func (b *Bookmark) RemoveTag(tag string) {
	tags := []string{}
	for _, t := range b.Tags {
		if t != tag {
			tags = append(tags, t)
		}
	}
	b.Tags = tags
	b.Update()
}

func (b *Bookmark) SetCategory(category string) {
	b.Category = category
	b.Update()
}

func (b *Bookmark) SetTitle(title string) {
	b.Title = title
	b.Update()
}

func (b *Bookmark) SetURL(url string) {
	b.URL = url
	b.Update()
}

func (b *Bookmark) SetDescription(description string) {
	b.Description = description
	b.Update()
}