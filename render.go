package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func render(content, filename, lang string, th theme, width int, showLineNumbers bool) (string, error) {
	ft := detectFileTypeWithLang(filename, lang)

	switch ft {
	case fileTypeMarkdown:
		return renderMarkdown(content, th, width)
	case fileTypeCode:
		return renderCode(content, filename, lang, th, showLineNumbers)
	default:
		return renderPlain(content, showLineNumbers), nil
	}
}

func renderPlain(content string, showLineNumbers bool) string {
	if !showLineNumbers {
		return content
	}

	lines := strings.Split(content, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	width := len(fmt.Sprintf("%d", len(lines)))
	gutterStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("243"))

	var b strings.Builder
	for i, line := range lines {
		num := fmt.Sprintf("%*d", width, i+1)
		b.WriteString(gutterStyle.Render(num))
		b.WriteString("  ")
		b.WriteString(line)
		b.WriteString("\n")
	}
	return b.String()
}
