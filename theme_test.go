package main

import "testing"

func TestDefaultTheme(t *testing.T) {
	th := resolveTheme("")
	if th.chromaStyle == "" {
		t.Error("expected non-empty default chroma style")
	}
	if th.glamourStyle == "" {
		t.Error("expected non-empty default glamour style")
	}
}

func TestDraculaTheme(t *testing.T) {
	th := resolveTheme("dracula")
	if th.chromaStyle != "dracula" {
		t.Errorf("expected chroma style dracula, got %s", th.chromaStyle)
	}
	if th.glamourStyle != "dracula" {
		t.Errorf("expected glamour style dracula, got %s", th.glamourStyle)
	}
}

func TestMonokaiTheme(t *testing.T) {
	th := resolveTheme("monokai")
	if th.chromaStyle != "monokai" {
		t.Errorf("expected chroma style monokai, got %s", th.chromaStyle)
	}
}

func TestThemeFromEnv(t *testing.T) {
	t.Setenv("PEEK_THEME", "github-dark")
	th := resolveThemeFromEnv("")
	if th.chromaStyle != "github-dark" {
		t.Errorf("expected github-dark from env, got %s", th.chromaStyle)
	}
}

func TestThemeFlagOverridesEnv(t *testing.T) {
	t.Setenv("PEEK_THEME", "github-dark")
	th := resolveThemeFromEnv("dracula")
	if th.chromaStyle != "dracula" {
		t.Errorf("expected dracula from flag override, got %s", th.chromaStyle)
	}
}
