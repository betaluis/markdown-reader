package main

import (
	"bytes"
	"os"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"
)

// RenderMarkdown reads a markdown file and converts it to HTML
func RenderMarkdown(filepath string) (string, error) {
	// Read the markdown file
	content, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	// Create Goldmark markdown parser with extensions
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,        // GitHub Flavored Markdown
			extension.Typographer, // Smart quotes, dashes
			highlighting.NewHighlighting(
				highlighting.WithStyle("github"),
				highlighting.WithFormatOptions(
					html.WithClasses(false), // Use inline styles
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(), // Auto-generate heading IDs
		),
		goldmark.WithRendererOptions(
			goldmarkhtml.WithHardWraps(), // Line breaks create <br>
			goldmarkhtml.WithUnsafe(),    // Allow raw HTML in markdown
		),
	)

	// Convert markdown to HTML
	var buf bytes.Buffer
	if err := md.Convert(content, &buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}
