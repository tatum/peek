package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/term"
)

type runOpts struct {
	lang            string
	th              theme
	showLineNumbers bool
	usePager        bool
	noColor         bool
	isTTY           bool
	forceColor      bool
	forceKitty      bool
	imgMaxCols      int
	imgMaxRows      int
}

func main() {
	var (
		pager      bool
		lang       string
		themeName  string
		noLines    bool
		forceColor bool
		useFzf     bool
		forceKitty bool
		imgWidth   int
		imgHeight  int
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
	flag.BoolVar(&useFzf, "z", false, "pick a file with fzf, then peek it")
	flag.BoolVar(&useFzf, "fzf", false, "pick a file with fzf, then peek it")
	flag.BoolVar(&forceKitty, "force-kitty", false, "render images via kitty protocol even if terminal capability isn't detected")
	flag.IntVar(&imgWidth, "width", 0, "max image width in terminal columns (0 = full terminal width)")
	flag.IntVar(&imgHeight, "height", 0, "max image height in terminal rows (0 = unbounded)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: peek [flags] [file...]\n\nFlags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	noColor := os.Getenv("NO_COLOR") != ""
	isTTY := term.IsTerminal(int(os.Stdout.Fd()))

	opts := runOpts{
		lang:            lang,
		th:              resolveThemeFromEnv(themeName),
		showLineNumbers: !noLines,
		usePager:        pager,
		noColor:         noColor,
		isTTY:           isTTY,
		forceColor:      forceColor,
		forceKitty:      forceKitty,
		imgMaxCols:      imgWidth,
		imgMaxRows:      imgHeight,
	}

	files := flag.Args()

	if useFzf {
		picked, err := runFzf(lang, themeName, noLines)
		if err != nil {
			fmt.Fprintf(os.Stderr, "peek: %v\n", err)
			os.Exit(1)
		}
		if picked == "" {
			return
		}
		files = []string{picked}
	}

	if len(files) == 0 {
		if err := processReader(os.Stdin, "stdin", opts); err != nil {
			fmt.Fprintf(os.Stderr, "peek: %v\n", err)
			os.Exit(1)
		}
		return
	}

	for _, path := range files {
		if err := processFile(path, opts); err != nil {
			fmt.Fprintf(os.Stderr, "peek: %s: %v\n", path, err)
			os.Exit(1)
		}
	}
}

func processFile(path string, opts runOpts) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if format := detectImageFormat(path, data); format != imageFormatNone {
		return handleImage(data, format, path, opts)
	}

	if isBinary(data) {
		fmt.Fprintf(os.Stderr, "peek: %s: binary file, not rendering\n", path)
		return nil
	}

	return renderAndOutput(string(data), path, opts)
}

func processReader(r io.Reader, name string, opts runOpts) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	if format := detectImageFormat(name, data); format != imageFormatNone {
		return handleImage(data, format, name, opts)
	}

	if isBinary(data) {
		fmt.Fprintln(os.Stderr, "peek: binary input, not rendering")
		return nil
	}

	return renderAndOutput(string(data), name, opts)
}

// detectImageFormat resolves an image format from the filename's extension
// first, then falls back to magic-byte sniffing (covers stdin / extension-less
// paths). Returns imageFormatNone when nothing matches.
func detectImageFormat(name string, data []byte) imageFormat {
	if name != "" {
		if f := imageFormatFromExt(strings.ToLower(filepath.Ext(name))); f != imageFormatNone {
			return f
		}
	}
	return imageFormatFromMagic(data)
}

func handleImage(data []byte, format imageFormat, name string, opts runOpts) error {
	if !opts.isTTY {
		fmt.Fprintf(os.Stderr, "peek: %s: image, not rendering (stdout is not a terminal)\n", name)
		return nil
	}
	if opts.usePager {
		fmt.Fprintf(os.Stderr, "peek: %s: image, not rendering in pager mode\n", name)
		return nil
	}
	if !opts.forceKitty && !kittySupportEnv() {
		fmt.Fprintf(os.Stderr, "peek: %s: image, terminal does not appear to support kitty graphics (use --force-kitty to override)\n", name)
		return nil
	}
	return renderImageKitty(os.Stdout, data, format, imageRenderOpts{
		maxCols: opts.imgMaxCols,
		maxRows: opts.imgMaxRows,
	})
}

func renderAndOutput(content, filename string, opts runOpts) error {
	width := 80
	if opts.isTTY {
		if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && w > 0 {
			width = w
		}
	}

	out, err := render(content, filename, opts.lang, opts.th, width, opts.showLineNumbers)
	if err != nil {
		return err
	}

	if opts.noColor || (!opts.isTTY && !opts.forceColor) {
		out = stripANSI(out)
	}

	if opts.usePager && opts.isTTY {
		return outputWithPager(out)
	}

	fmt.Print(out)
	return nil
}
