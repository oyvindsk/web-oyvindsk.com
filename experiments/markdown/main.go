package main

import (
	"fmt"
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
)

func main() {
	// Read the contents of input.md
	md, err := os.ReadFile("input.md")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input.md: %v\n", err)
		os.Exit(1)
	}

	// parsed := markdown.Parse(md)

	// https://blog.kowalczyk.info/article/cxn3/advanced-markdown-processing-in-go.html

	htmlFlags := html.CommonFlags | html.HrefTargetBlank | html.CompletePage
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	html := markdown.ToHTML(md, nil, renderer)
	fmt.Printf("\n%s\n", html)
}
