package main

import (
	"strings"
	"testing"
)

func TestRenderCodeGo(t *testing.T) {
	code := `package main

func main() {
	println("hello")
}
`
	out, err := renderCode(code, "main.go", "", resolveTheme(""), true)
	if err != nil {
		t.Fatal(err)
	}
	if out == "" {
		t.Error("expected non-empty output")
	}
	// Should contain ANSI escape codes for color
	if !strings.Contains(out, "\033[") {
		t.Error("expected ANSI color codes in output")
	}
}

func TestRenderCodeWithLineNumbers(t *testing.T) {
	code := "line1\nline2\nline3\n"
	out, err := renderCode(code, "test.py", "", resolveTheme(""), true)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "1") {
		t.Error("expected line numbers in output")
	}
}

func TestRenderCodeWithLangOverride(t *testing.T) {
	code := `{"key": "value"}`
	out, err := renderCode(code, "data.txt", "json", resolveTheme(""), false)
	if err != nil {
		t.Fatal(err)
	}
	if out == "" {
		t.Error("expected non-empty output")
	}
}

func TestRenderCodeNoLineNumbers(t *testing.T) {
	code := "x = 1\n"
	out, err := renderCode(code, "test.py", "", resolveTheme(""), false)
	if err != nil {
		t.Fatal(err)
	}
	if out == "" {
		t.Error("expected non-empty output")
	}
}
