# ubm - URL Bookmark Manager

ubm (URL Bookmark Manager) is an interactive command-line tool that allows you to organize bookmarks in a tree structure.

## Features

- 📁 **Hierarchical Category Management**: Organize bookmarks in a tree structure
- 🔍 **Interactive Navigation**: Easy bookmark exploration with arrow keys
- 🌐 **Browser Integration**: Automatically open selected bookmarks in browser
- ✏️ **Edit Function**: Edit bookmark information later
- 📂 **Move Between Categories**: Move bookmarks to different categories

## Installation

### Install with Go

```bash
go install github.com/tom-023/ubm/cmd/ubm@latest
```

### Build from Source

```bash
git clone https://github.com/tom-023/ubm.git
cd ubm
go build -o ubm ./cmd/ubm
```

## Usage

### Add Bookmarks

```bash
# Add interactively (prompts for URL, title, and category)
ubm add
```

When adding a bookmark, you'll see:
1. URL prompt
2. Title prompt (with auto-suggested domain name)
3. Category selection screen where you can:
   - Navigate existing categories with arrow keys
   - Create new categories by selecting "➕ Create new category"
   - Select current category with "✅ Select this category"

### Browse Bookmarks

```bash
# Interactive navigation (opens selected bookmark in browser and exits)
ubm list

# Show all in tree format
ubm show

# Show in flat list
ubm show --flat
```

### Category Management

```bash
# Create category
ubm category create

# List categories
ubm category list

# Delete empty category
ubm category delete
```

### Edit Bookmarks

```bash
# Edit by title
ubm edit "bookmark title"

# Edit interactively
ubm edit

# Move by title
ubm move "bookmark title"

# Move interactively
ubm move

# Delete bookmark
ubm delete
```

## Keyboard Shortcuts

In interactive mode:

- `↑` `↓`: Select items
- `Enter`: Select/Confirm
- `Backspace` `Esc`: Go back to parent directory
- `/`: Toggle search mode
- `q` `Ctrl+C`: Quit

## Data Storage

Bookmarks are stored in:

- **Linux/macOS**: `~/.config/ubm/bookmarks.json`
- **Windows**: `%APPDATA%\ubm\bookmarks.json`

## Development

### Requirements

- Go 1.23+

### Dependencies

- [cobra](https://github.com/spf13/cobra) - CLI framework
- [promptui](https://github.com/manifoldco/promptui) - Interactive prompts
- [browser](https://github.com/pkg/browser) - Browser integration

### Build

```bash
go build -o ubm ./cmd/ubm
```

### Test

```bash
go test ./...
```

## License

MIT License