package main

import (
	"strings"
	"testing"
)

func TestRenderMarkdownHeading(t *testing.T) {
	md := "# Hello World\n\nSome text here.\n"
	out, err := renderMarkdown(md, resolveTheme(""), 80)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "Hello") {
		t.Error("expected heading text in output")
	}
}

func TestRenderMarkdownCodeBlock(t *testing.T) {
	md := "```go\nfunc main() {}\n```\n"
	out, err := renderMarkdown(md, resolveTheme(""), 80)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "func") {
		t.Error("expected code block content in output")
	}
}

func TestRenderMarkdownWordWrap(t *testing.T) {
	md := "# Title\n\nThis is a paragraph.\n"
	out, err := renderMarkdown(md, resolveTheme(""), 40)
	if err != nil {
		t.Fatal(err)
	}
	if out == "" {
		t.Error("expected non-empty output with word wrap")
	}
}
