package main

import (
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma/v2/lexers"
)

type fileType int

const (
	fileTypePlain    fileType = iota
	fileTypeMarkdown
	fileTypeCode
)

// detectFileType determines file type from filename.
func detectFileType(filename string) fileType {
	return detectFileTypeWithLang(filename, "")
}

// detectFileTypeWithLang determines file type, with optional language override.
func detectFileTypeWithLang(filename, lang string) fileType {
	if lang != "" {
		if lang == "markdown" || lang == "md" {
			return fileTypeMarkdown
		}
		return fileTypeCode
	}

	ext := strings.ToLower(filepath.Ext(filename))
	if ext == ".md" || ext == ".markdown" || ext == ".mdown" || ext == ".mkd" {
		return fileTypeMarkdown
	}

	// Use chroma's lexer registry to check if it's a known language
	lexer := lexers.Match(filename)
	if lexer != nil {
		return fileTypeCode
	}

	return fileTypePlain
}
