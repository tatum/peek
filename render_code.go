package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/charmbracelet/lipgloss"
)

func renderCode(code, filename, lang string, th theme, showLineNumbers bool) (string, error) {
	var lexer chroma.Lexer
	if lang != "" {
		lexer = lexers.Get(lang)
	}
	if lexer == nil {
		lexer = lexers.Match(filename)
	}
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	style := styles.Get(th.chromaStyle)
	if style == nil {
		style = styles.Fallback
	}

	formatter := formatters.Get("terminal256")
	if formatter == nil {
		formatter = formatters.Fallback
	}

	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return "", fmt.Errorf("tokenize: %w", err)
	}

	var buf bytes.Buffer
	if err := formatter.Format(&buf, style, iterator); err != nil {
		return "", fmt.Errorf("format: %w", err)
	}

	if !showLineNumbers {
		return buf.String(), nil
	}

	return addLineNumbers(buf.String()), nil
}

func addLineNumbers(highlighted string) string {
	lines := strings.Split(highlighted, "\n")
	// Remove trailing empty line from split
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
