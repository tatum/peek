package main

import (
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma/v2/lexers"
)

type fileType int

const (
	fileTypePlain fileType = iota
	fileTypeMarkdown
	fileTypeCode
	fileTypeImage
)

type imageFormat int

const (
	imageFormatNone imageFormat = iota
	imageFormatPNG
	imageFormatJPEG
	imageFormatGIF
	imageFormatWebP
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
	if imageFormatFromExt(ext) != imageFormatNone {
		return fileTypeImage
	}

	// Use chroma's lexer registry to check if it's a known language
	lexer := lexers.Match(filename)
	if lexer != nil {
		return fileTypeCode
	}

	return fileTypePlain
}

// imageFormatFromExt maps a lowercased extension (including leading dot) to an
// imageFormat, or imageFormatNone if it isn't a supported image extension.
func imageFormatFromExt(ext string) imageFormat {
	switch ext {
	case ".png":
		return imageFormatPNG
	case ".jpg", ".jpeg":
		return imageFormatJPEG
	case ".gif":
		return imageFormatGIF
	case ".webp":
		return imageFormatWebP
	}
	return imageFormatNone
}

// imageFormatFromMagic sniffs the leading bytes for a supported image format.
// Returns imageFormatNone if no signature matches.
func imageFormatFromMagic(data []byte) imageFormat {
	if len(data) >= 8 && string(data[:8]) == "\x89PNG\r\n\x1a\n" {
		return imageFormatPNG
	}
	if len(data) >= 3 && data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		return imageFormatJPEG
	}
	if len(data) >= 6 && (string(data[:6]) == "GIF87a" || string(data[:6]) == "GIF89a") {
		return imageFormatGIF
	}
	if len(data) >= 12 && string(data[:4]) == "RIFF" && string(data[8:12]) == "WEBP" {
		return imageFormatWebP
	}
	return imageFormatNone
}
