package main

import "os"

type theme struct {
	chromaStyle  string
	glamourStyle string
}

// glamour only supports a few built-in style names; map common ones
var glamourStyles = map[string]string{
	"dracula":     "dracula",
	"dark":        "dark",
	"light":       "light",
	"tokyo-night": "tokyo-night",
	"pink":        "pink",
	"ascii":       "ascii",
}

func resolveThemeFromEnv(flagValue string) theme {
	name := flagValue
	if name == "" {
		name = os.Getenv("PEEK_THEME")
	}
	return resolveTheme(name)
}

func resolveTheme(name string) theme {
	if name == "" {
		return theme{
			chromaStyle:  "monokai",
			glamourStyle: "dark",
		}
	}

	glamour := "dark"
	if gs, ok := glamourStyles[name]; ok {
		glamour = gs
	}

	return theme{
		chromaStyle:  name,
		glamourStyle: glamour,
	}
}
