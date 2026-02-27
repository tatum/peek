package main

import (
	"fmt"

	"github.com/charmbracelet/glamour"
)

func renderMarkdown(source string, th theme, width int) (string, error) {
	r, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle(th.glamourStyle),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return "", fmt.Errorf("create markdown renderer: %w", err)
	}

	out, err := r.Render(source)
	if err != nil {
		return "", fmt.Errorf("render markdown: %w", err)
	}

	return out, nil
}
