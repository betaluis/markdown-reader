# mdview

Local markdown file viewer with live reload and GitHub-style rendering.

## Features

- **GitHub-style rendering** - Clean, familiar markdown styling
- **Live reload** - Auto-refresh browser when file changes
- **Syntax highlighting** - Code blocks with 200+ languages supported
- **Auto browser launch** - Opens your default browser automatically
- **Simple CLI** - Easy to use command-line interface

## Installation

### Prerequisites

- Go 1.21 or higher

### Build from Source

```bash
# Clone the repository
git clone https://github.com/betaluis/markdown-reader.git
cd markdown-reader

# Build the binary
go build -o mdview

# Install to ~/.local/bin
cp mdview ~/.local/bin/

# Verify installation (ensure ~/.local/bin is in your PATH)
mdview --version
```

## Usage

### Basic Usage

```bash
mdview <file.md>
```

This will:
1. Start a local server on port 3000
2. Open your default browser
3. Display the rendered markdown
4. Watch for file changes and auto-reload

### Examples

```bash
# View a markdown file
mdview README.md

# Use a custom port
mdview --port 8080 docs/guide.md

# Don't open browser automatically
mdview --no-browser file.md

# View help
mdview --help

# Show version
mdview --version
```

### Options

- `--port <port>` - Custom port (default: 3000)
- `--no-browser` - Don't auto-open browser
- `--help, -h` - Show help message
- `--version, -v` - Show version

## How It Works

1. **Markdown Parsing** - Uses [Goldmark](https://github.com/yuin/goldmark) with GitHub Flavored Markdown (GFM) support
2. **Syntax Highlighting** - Powered by [Chroma](https://github.com/alecthomas/chroma)
3. **File Watching** - [fsnotify](https://github.com/fsnotify/fsnotify) monitors file changes
4. **Live Reload** - WebSocket connection broadcasts reload signals
5. **Browser Launch** - [pkg/browser](https://github.com/pkg/browser) for cross-platform support

## Supported Markdown Features

- Headers (h1-h6)
- Text formatting (bold, italic, strikethrough)
- Lists (ordered, unordered, nested)
- Links and auto-linking
- Code blocks with syntax highlighting
- Inline code
- Blockquotes
- Tables (GitHub Flavored Markdown)
- Task lists
- Horizontal rules
- Images

## Development

### Project Structure

```
local-md-reader/
├── main.go         # CLI parsing and orchestration
├── server.go       # HTTP server and WebSocket handler
├── renderer.go     # Markdown to HTML conversion
├── watcher.go      # File change monitoring
├── templates.go    # HTML template and CSS
├── examples/       # Test markdown files
│   └── test.md
├── go.mod          # Go module definition
└── README.md       # This file
```

### Testing

```bash
# Run with test file
go run . examples/test.md

# Or use the binary
mdview examples/test.md
```

## Troubleshooting

**Browser doesn't open automatically:**
- Check if `xdg-open` is available on Linux
- Use the `--no-browser` flag and open `http://localhost:3000` manually

**Port already in use:**
- Use a different port: `mdview --port 8080 file.md`

**File changes not detected:**
- Some editors use atomic writes - save the file directly instead

## License

MIT

## Author

Built with Go, Goldmark, and Chroma