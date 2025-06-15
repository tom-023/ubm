# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

ubm (URL Bookmark Manager) is a CLI tool for managing bookmarks with hierarchical category organization and interactive navigation.

## Build and Development Commands

```bash
# Build the application
go build -o ubm ./cmd/ubm

# Install locally
go install ./cmd/ubm

# Run tests (when implemented)
go test ./...

# Update dependencies
go mod tidy
```

## Architecture and Component Interaction

### Core Architecture Pattern
The project follows a layered architecture with clear separation of concerns:

1. **Command Layer** (`cmd/ubm/`): CLI commands using Cobra framework. Each command file (add.go, list.go, etc.) implements a specific user action.

2. **Storage Layer** (`internal/storage/`): Handles data persistence with JSON files. Key features:
   - Thread-safe operations using sync.RWMutex
   - Automatic backup creation before writes
   - Atomic writes using temporary files

3. **Domain Models** (`internal/bookmark/`, `internal/category/`): Core business logic
   - Bookmarks are stored with UUID identifiers
   - Categories use slash-separated paths (e.g., "programming/go/tutorials")

4. **UI Layer** (`internal/ui/`): Interactive prompts and navigation
   - Uses promptui for interactive selections
   - Custom color templates for dark terminal visibility

### Key Design Decisions

1. **Global Storage Instance**: A single storage instance is created in main.go and shared across all commands via a package-level variable.

2. **Category Tree Building**: Categories are stored as flat paths in JSON but converted to tree structures at runtime for navigation.

3. **Interactive UI Flow**: All user interactions go through the ui package, which handles cancellation (Ctrl+C) gracefully.

### Data Flow Example (Add Command)
```
User Input → ui.PromptURL() → validator.NormalizeURL() → ui.PromptString() (title) 
→ ui.SelectCategory() → bookmark.New() → store.AddBookmark() → JSON file
```

## Important Implementation Details

### Storage Format
Data is stored in `~/.config/ubm/bookmarks.json` with this structure:
```json
{
  "bookmarks": [...],
  "categories": ["programming", "programming/go", ...],
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### UI Color Scheme
- Prompts: Cyan bold
- Active items: Cyan
- Inactive items: Faint
- Selected items: Green with checkmark

### Cancel Handling
All interactive prompts check for `promptui.ErrInterrupt` or `promptui.ErrEOF` and return a "cancelled" error that commands handle by displaying "Cancelled." and exiting gracefully.