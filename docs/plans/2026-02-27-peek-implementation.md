# peek Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build `peek`, a Go CLI that renders any text file with color and formatting in the terminal.

**Architecture:** Three layers — CLI (flag parsing, stdin/file reading, TTY detection), Detector (extension → language mapping via chroma), Renderer (glamour for markdown, chroma for code, plain text fallback). Output routes to stdout or pager.

**Tech Stack:** Go, github.com/alecthomas/chroma/v2, github.com/charmbracelet/glamour/v2, github.com/charmbracelet/lipgloss/v2

---

### Task 1: Project Scaffold

**Files:**
- Create: `go.mod`
- Create: `main.go`
- Create: `.gitignore`

**Step 1: Initialize Go module and add dependencies**

Run:
```bash
cd /Users/tatum/code/tries/2026-02-27-codecat
go mod init github.com/tatum/peek
go get github.com/alecthomas/chroma/v2@latest
go get github.com/charmbracelet/glamour/v2@latest
go get github.com/charmbracelet/lipgloss/v2@latest
```

**Step 2: Create .gitignore**

```
peek
*.exe
```

**Step 3: Create minimal main.go that compiles**

```go
package main

import "fmt"

func main() {
	fmt.Println("peek")
}
```

**Step 4: Verify it builds**

Run: `go build -o peek .`
Expected: binary created, no errors

**Step 5: Commit**

```bash
git add go.mod go.sum main.go .gitignore
git commit -m "feat: scaffold peek project with Go module"
```

---

### Task 2: File Type Detection

**Files:**
- Create: `detect.go`
- Create: `detect_test.go`

**Step 1: Write the failing tests**

```go
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
```

**Step 2: Run tests to verify they fail**

Run: `go test -run TestDetect -v ./...`
Expected: FAIL — functions not defined

**Step 3: Write implementation**

```go
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
```

**Step 4: Run tests to verify they pass**

Run: `go test -run TestDetect -v ./...`
Expected: all PASS

**Step 5: Commit**

```bash
git add detect.go detect_test.go
git commit -m "feat: add file type detection using chroma lexer registry"
```

---

### Task 3: Theme System

**Files:**
- Create: `theme.go`
- Create: `theme_test.go`

**Step 1: Write the failing tests**

```go
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
```

**Step 2: Run tests to verify they fail**

Run: `go test -run TestTheme -v ./... && go test -run TestDefault -v ./... && go test -run TestDracula -v ./... && go test -run TestMonokai -v ./...`
Expected: FAIL

**Step 3: Write implementation**

```go
package main

import "os"

type theme struct {
	chromaStyle  string
	glamourStyle string
}

// glamour only supports a few built-in style names; map common ones
var glamourStyles = map[string]string{
	"dracula":    "dracula",
	"dark":       "dark",
	"light":      "light",
	"tokyo-night": "tokyo-night",
	"pink":       "pink",
	"ascii":      "ascii",
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
```

**Step 4: Run tests to verify they pass**

Run: `go test -run "TestDefault|TestDracula|TestMonokai|TestTheme" -v ./...`
Expected: all PASS

**Step 5: Commit**

```bash
git add theme.go theme_test.go
git commit -m "feat: add theme resolution from env var and flag"
```

---

### Task 4: Code Renderer (chroma)

**Files:**
- Create: `render_code.go`
- Create: `render_code_test.go`

**Step 1: Write the failing tests**

```go
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
```

**Step 2: Run tests to verify they fail**

Run: `go test -run TestRenderCode -v ./...`
Expected: FAIL

**Step 3: Write implementation**

```go
package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/charmbracelet/lipgloss/v2"
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
```

**Step 4: Run tests to verify they pass**

Run: `go test -run TestRenderCode -v ./...`
Expected: all PASS

**Step 5: Commit**

```bash
git add render_code.go render_code_test.go
git commit -m "feat: add code renderer with chroma syntax highlighting"
```

---

### Task 5: Markdown Renderer (glamour)

**Files:**
- Create: `render_md.go`
- Create: `render_md_test.go`

**Step 1: Write the failing tests**

```go
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
	if !strings.Contains(out, "Hello World") {
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
```

**Step 2: Run tests to verify they fail**

Run: `go test -run TestRenderMarkdown -v ./...`
Expected: FAIL

**Step 3: Write implementation**

```go
package main

import (
	"fmt"

	"github.com/charmbracelet/glamour"
)

func renderMarkdown(source string, th theme, width int) (string, error) {
	r, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle(th.glamourStyle),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return "", fmt.Errorf("create markdown renderer: %w", err)
	}

	out, err := r.Render(source)
	if err != nil {
		return "", fmt.Errorf("render markdown: %w", err)
	}

	return out, nil
}
```

**Note:** Check the import path — glamour v2 may be `github.com/charmbracelet/glamour/v2` or just `github.com/charmbracelet/glamour` depending on what go get resolved. Adjust the import to match go.sum.

**Step 4: Run tests to verify they pass**

Run: `go test -run TestRenderMarkdown -v ./...`
Expected: all PASS

**Step 5: Commit**

```bash
git add render_md.go render_md_test.go
git commit -m "feat: add markdown renderer with glamour"
```

---

### Task 6: Renderer Dispatch

**Files:**
- Create: `render.go`
- Create: `render_test.go`

**Step 1: Write the failing tests**

```go
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
```

**Step 2: Run tests to verify they fail**

Run: `go test -run TestRenderDispatch -v ./...`
Expected: FAIL

**Step 3: Write implementation**

```go
package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
)

func render(content, filename, lang string, th theme, width int, showLineNumbers bool) (string, error) {
	ft := detectFileTypeWithLang(filename, lang)

	switch ft {
	case fileTypeMarkdown:
		return renderMarkdown(content, th, width)
	case fileTypeCode:
		return renderCode(content, filename, lang, th, showLineNumbers)
	default:
		return renderPlain(content, showLineNumbers), nil
	}
}

func renderPlain(content string, showLineNumbers bool) string {
	if !showLineNumbers {
		return content
	}

	lines := strings.Split(content, "\n")
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
```

**Step 4: Run tests to verify they pass**

Run: `go test -run TestRenderDispatch -v ./...`
Expected: all PASS

**Step 5: Commit**

```bash
git add render.go render_test.go
git commit -m "feat: add renderer dispatch with plain text fallback"
```

---

### Task 7: Pager Support

**Files:**
- Create: `pager.go`
- Create: `pager_test.go`

**Step 1: Write the failing tests**

```go
package main

import "testing"

func TestPagerCommand(t *testing.T) {
	cmd, args := pagerCommand()
	if cmd == "" {
		t.Error("expected non-empty pager command")
	}
	_ = args // just verify it doesn't panic
}

func TestPagerCommandFromEnv(t *testing.T) {
	t.Setenv("PAGER", "more")
	cmd, args := pagerCommand()
	if cmd != "more" {
		t.Errorf("expected 'more' from PAGER env, got %s", cmd)
	}
	if len(args) != 0 {
		t.Errorf("expected no args for 'more', got %v", args)
	}
}

func TestPagerCommandDefault(t *testing.T) {
	t.Setenv("PAGER", "")
	cmd, args := pagerCommand()
	if cmd != "less" {
		t.Errorf("expected 'less' as default, got %s", cmd)
	}
	found := false
	for _, a := range args {
		if a == "-R" {
			found = true
		}
	}
	if !found {
		t.Error("expected -R flag for less")
	}
}
```

**Step 2: Run tests to verify they fail**

Run: `go test -run TestPager -v ./...`
Expected: FAIL

**Step 3: Write implementation**

```go
package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func pagerCommand() (string, []string) {
	pager := os.Getenv("PAGER")
	if pager == "" {
		return "less", []string{"-R"}
	}

	parts := strings.Fields(pager)
	if len(parts) == 1 {
		return parts[0], nil
	}
	return parts[0], parts[1:]
}

func outputWithPager(content string) error {
	cmd, args := pagerCommand()
	pager := exec.Command(cmd, args...)
	pager.Stdout = os.Stdout
	pager.Stderr = os.Stderr

	stdin, err := pager.StdinPipe()
	if err != nil {
		return fmt.Errorf("pager stdin: %w", err)
	}

	if err := pager.Start(); err != nil {
		// Fallback: just print if pager fails
		fmt.Print(content)
		return nil
	}

	_, _ = io.WriteString(stdin, content)
	stdin.Close()

	return pager.Wait()
}
```

**Step 4: Run tests to verify they pass**

Run: `go test -run TestPager -v ./...`
Expected: all PASS

**Step 5: Commit**

```bash
git add pager.go pager_test.go
git commit -m "feat: add pager support with PAGER env fallback"
```

---

### Task 8: Binary File Detection

**Files:**
- Create: `binary.go`
- Create: `binary_test.go`

**Step 1: Write the failing tests**

```go
package main

import "testing"

func TestIsBinaryTrue(t *testing.T) {
	// Bytes with null characters indicate binary
	data := []byte{0x00, 0x01, 0x02, 0xFF}
	if !isBinary(data) {
		t.Error("expected binary detection for null bytes")
	}
}

func TestIsBinaryFalse(t *testing.T) {
	data := []byte("Hello, this is plain text.\nWith newlines.\n")
	if isBinary(data) {
		t.Error("expected text detection for plain ASCII")
	}
}

func TestIsBinaryUTF8(t *testing.T) {
	data := []byte("Hello 世界 🌍")
	if isBinary(data) {
		t.Error("expected text detection for valid UTF-8")
	}
}

func TestIsBinaryEmpty(t *testing.T) {
	if isBinary([]byte{}) {
		t.Error("empty input should not be binary")
	}
}
```

**Step 2: Run tests to verify they fail**

Run: `go test -run TestIsBinary -v ./...`
Expected: FAIL

**Step 3: Write implementation**

```go
package main

// isBinary checks if data looks like a binary file by scanning for null bytes
// in the first 8000 bytes (same heuristic git uses).
func isBinary(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	checkLen := 8000
	if len(data) < checkLen {
		checkLen = len(data)
	}
	for i := 0; i < checkLen; i++ {
		if data[i] == 0 {
			return true
		}
	}
	return false
}
```

**Step 4: Run tests to verify they pass**

Run: `go test -run TestIsBinary -v ./...`
Expected: all PASS

**Step 5: Commit**

```bash
git add binary.go binary_test.go
git commit -m "feat: add binary file detection"
```

---

### Task 9: CLI Entry Point

**Files:**
- Modify: `main.go`

**Step 1: Write the full CLI entry point**

Replace main.go with:

```go
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"golang.org/x/term"
)

func main() {
	var (
		pager      bool
		lang       string
		themeName  string
		noLines    bool
		forceColor bool
	)

	flag.BoolVar(&pager, "p", false, "open in pager mode")
	flag.BoolVar(&pager, "pager", false, "open in pager mode")
	flag.StringVar(&lang, "l", "", "force language for syntax highlighting")
	flag.StringVar(&lang, "lang", "", "force language for syntax highlighting")
	flag.StringVar(&themeName, "t", "", "color theme (overrides PEEK_THEME)")
	flag.StringVar(&themeName, "theme", "", "color theme (overrides PEEK_THEME)")
	flag.BoolVar(&noLines, "n", false, "hide line numbers")
	flag.BoolVar(&noLines, "no-lines", false, "hide line numbers")
	flag.BoolVar(&forceColor, "force-color", false, "force color output when piped")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: peek [flags] [file...]\n\nFlags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	// Respect NO_COLOR
	noColor := os.Getenv("NO_COLOR") != ""
	isTTY := term.IsTerminal(int(os.Stdout.Fd()))

	th := resolveThemeFromEnv(themeName)

	files := flag.Args()
	if len(files) == 0 {
		// Read from stdin
		if err := processReader(os.Stdin, "stdin", lang, th, !noLines, pager, noColor, isTTY, forceColor); err != nil {
			fmt.Fprintf(os.Stderr, "peek: %v\n", err)
			os.Exit(1)
		}
		return
	}

	for _, path := range files {
		if err := processFile(path, lang, th, !noLines, pager, noColor, isTTY, forceColor); err != nil {
			fmt.Fprintf(os.Stderr, "peek: %s: %v\n", path, err)
			os.Exit(1)
		}
	}
}

func processFile(path, lang string, th theme, showLineNumbers, usePager, noColor, isTTY, forceColor bool) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if isBinary(data) {
		fmt.Fprintf(os.Stderr, "peek: %s: binary file, not rendering\n", path)
		return nil
	}

	return renderAndOutput(string(data), path, lang, th, showLineNumbers, usePager, noColor, isTTY, forceColor)
}

func processReader(r io.Reader, name, lang string, th theme, showLineNumbers, usePager, noColor, isTTY, forceColor bool) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	if isBinary(data) {
		fmt.Fprintln(os.Stderr, "peek: binary input, not rendering")
		return nil
	}

	return renderAndOutput(string(data), name, lang, th, showLineNumbers, usePager, noColor, isTTY, forceColor)
}

func renderAndOutput(content, filename, lang string, th theme, showLineNumbers, usePager, noColor, isTTY, forceColor bool) error {
	width := 80
	if isTTY {
		if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && w > 0 {
			width = w
		}
	}

	out, err := render(content, filename, lang, th, width, showLineNumbers)
	if err != nil {
		return err
	}

	// Strip color if not TTY and not forced
	if noColor || (!isTTY && !forceColor) {
		out = stripANSI(out)
	}

	if usePager && isTTY {
		return outputWithPager(out)
	}

	fmt.Print(out)
	return nil
}
```

**Step 2: Add ANSI stripping utility — create `ansi.go`**

```go
package main

import "regexp"

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

func stripANSI(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}
```

**Step 3: Add golang.org/x/term dependency**

Run: `go get golang.org/x/term@latest`

**Step 4: Verify it builds**

Run: `go build -o peek .`
Expected: binary created, no errors

**Step 5: Manual smoke test**

Run: `./peek main.go`
Expected: syntax-highlighted Go code with line numbers

Run: `echo '# Hello' | ./peek -l markdown`
Expected: rendered markdown heading

Run: `./peek go.mod`
Expected: plain text or code-highlighted output

**Step 6: Commit**

```bash
git add main.go ansi.go
git commit -m "feat: wire CLI entry point with all flags and rendering pipeline"
```

---

### Task 10: Integration Tests & Polish

**Files:**
- Create: `testdata/sample.py`
- Create: `testdata/sample.md`
- Create: `testdata/sample.json`
- Create: `main_test.go`

**Step 1: Create test fixture files**

`testdata/sample.py`:
```python
def greet(name):
    return f"Hello, {name}!"

if __name__ == "__main__":
    print(greet("world"))
```

`testdata/sample.md`:
```markdown
# Sample Document

This is a **bold** statement with `inline code`.

## Code Block

```python
x = 42
```

- Item one
- Item two
```

`testdata/sample.json`:
```json
{
  "name": "peek",
  "version": "0.1.0",
  "features": ["syntax-highlighting", "markdown", "pager"]
}
```

**Step 2: Write integration tests**

```go
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
	if !strings.Contains(out, "Sample Document") {
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
```

**Step 3: Run all tests**

Run: `go test -v ./...`
Expected: all PASS

**Step 4: Full build and manual end-to-end test**

Run:
```bash
go build -o peek .
./peek testdata/sample.py
./peek testdata/sample.md
./peek testdata/sample.json
./peek -p testdata/sample.py   # should open in pager
echo '{"a":1}' | ./peek -l json
./peek nonexistent.txt         # should error gracefully
```

**Step 5: Commit**

```bash
git add testdata/ main_test.go ansi.go
git commit -m "test: add integration tests and test fixtures"
```

---

### Task 11: Final Build & Install

**Step 1: Clean build**

Run:
```bash
go build -o peek .
```

**Step 2: Install to GOPATH**

Run:
```bash
go install .
```
Expected: `peek` available on PATH

**Step 3: Verify from anywhere**

Run:
```bash
peek --help
```
Expected: usage output with all flags

**Step 4: Commit any remaining changes**

```bash
git add -A
git commit -m "chore: final build verification"
```
