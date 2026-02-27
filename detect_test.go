package main

import "testing"

func TestDetectMarkdown(t *testing.T) {
	ft := detectFileType("README.md")
	if ft != fileTypeMarkdown {
		t.Errorf("expected markdown, got %v", ft)
	}
}

func TestDetectGoSource(t *testing.T) {
	ft := detectFileType("main.go")
	if ft != fileTypeCode {
		t.Errorf("expected code, got %v", ft)
	}
}

func TestDetectPython(t *testing.T) {
	ft := detectFileType("app.py")
	if ft != fileTypeCode {
		t.Errorf("expected code, got %v", ft)
	}
}

func TestDetectJSON(t *testing.T) {
	ft := detectFileType("config.json")
	if ft != fileTypeCode {
		t.Errorf("expected code, got %v", ft)
	}
}

func TestDetectUnknown(t *testing.T) {
	ft := detectFileType("notes.xyz")
	if ft != fileTypePlain {
		t.Errorf("expected plain, got %v", ft)
	}
}

func TestDetectNoExtension(t *testing.T) {
	ft := detectFileType("Makefile")
	if ft != fileTypeCode {
		t.Errorf("expected code for Makefile, got %v", ft)
	}
}

func TestDetectLangOverride(t *testing.T) {
	ft := detectFileTypeWithLang("data.txt", "json")
	if ft != fileTypeCode {
		t.Errorf("expected code when lang override is set, got %v", ft)
	}
}

func TestDetectLangOverrideMarkdown(t *testing.T) {
	ft := detectFileTypeWithLang("data.txt", "markdown")
	if ft != fileTypeMarkdown {
		t.Errorf("expected markdown when lang override is markdown, got %v", ft)
	}
}
