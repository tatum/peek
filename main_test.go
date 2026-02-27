package main

import (
	"os"
	"strings"
	"testing"
)

func TestIntegrationPythonFile(t *testing.T) {
	data, err := os.ReadFile("testdata/sample.py")
	if err != nil {
		t.Fatal(err)
	}
	out, err := render(string(data), "sample.py", "", resolveTheme(""), 80, true)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "greet") {
		t.Error("expected function name in output")
	}
}

func TestIntegrationMarkdownFile(t *testing.T) {
	data, err := os.ReadFile("testdata/sample.md")
	if err != nil {
		t.Fatal(err)
	}
	out, err := render(string(data), "sample.md", "", resolveTheme(""), 80, true)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "Sample") {
		t.Error("expected heading in output")
	}
}

func TestIntegrationJSONFile(t *testing.T) {
	data, err := os.ReadFile("testdata/sample.json")
	if err != nil {
		t.Fatal(err)
	}
	out, err := render(string(data), "sample.json", "", resolveTheme(""), 80, true)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "peek") {
		t.Error("expected JSON content in output")
	}
}

func TestStripANSI(t *testing.T) {
	input := "\033[38;5;208mhello\033[0m world"
	got := stripANSI(input)
	if got != "hello world" {
		t.Errorf("expected 'hello world', got %q", got)
	}
}
