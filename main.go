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
