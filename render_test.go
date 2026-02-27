package main

import (
	"strings"
	"testing"
)

func TestRenderDispatchMarkdown(t *testing.T) {
	md := "# Hello\n\nWorld.\n"
	out, err := render(md, "README.md", "", resolveTheme(""), 80, true)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "Hello") {
		t.Error("expected markdown rendering")
	}
}

func TestRenderDispatchCode(t *testing.T) {
	code := "package main\n"
	out, err := render(code, "main.go", "", resolveTheme(""), 80, true)
	if err != nil {
		t.Fatal(err)
	}
	if out == "" {
		t.Error("expected code rendering")
	}
}

func TestRenderDispatchPlainText(t *testing.T) {
	text := "just some text\nwith lines\n"
	out, err := render(text, "notes.xyz", "", resolveTheme(""), 80, true)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "just some text") {
		t.Error("expected plain text output")
	}
}

func TestRenderDispatchLangOverride(t *testing.T) {
	code := `{"key": "value"}`
	out, err := render(code, "data.txt", "json", resolveTheme(""), 80, true)
	if err != nil {
		t.Fatal(err)
	}
	// Should have ANSI codes because it was forced to render as JSON
	if !strings.Contains(out, "\033[") {
		t.Error("expected syntax highlighted output with lang override")
	}
}
