package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

// kittySupportEnv returns true if the current environment looks like a terminal
// that speaks the kitty graphics protocol. This is a best-effort env sniff; the
// --force-kitty flag bypasses it.
func kittySupportEnv() bool {
	if os.Getenv("KITTY_WINDOW_ID") != "" {
		return true
	}
	if os.Getenv("GHOSTTY_RESOURCES_DIR") != "" {
		return true
	}
	if t := os.Getenv("TERM"); strings.Contains(t, "kitty") || strings.Contains(t, "ghostty") {
		return true
	}
	switch os.Getenv("TERM_PROGRAM") {
	case "ghostty", "WezTerm":
		return true
	}
	return false
}

// imageRenderOpts controls how an image is sized for the terminal.
type imageRenderOpts struct {
	maxCols  int // 0 = use terminal width
	maxRows  int // 0 = unlimited (let aspect ratio decide)
	cellW    int // pixel width of one cell; 0 to query
	cellH    int // pixel height of one cell; 0 to query
	termCols int // total terminal columns; 0 to query
}

// renderImageKitty writes a kitty-graphics-protocol escape sequence to w that
// displays the image in data. If the image is not PNG, it is decoded and
// re-sent as raw 32-bit RGBA pixels.
func renderImageKitty(w io.Writer, data []byte, format imageFormat, opts imageRenderOpts) error {
	cols, rows, err := imageCellSize(data, opts)
	if err != nil {
		return err
	}

	var (
		payload []byte
		fmtCode int
		width   int
		height  int
	)

	if format == imageFormatPNG {
		payload = data
		fmtCode = 100 // PNG passthrough
	} else {
		img, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			return fmt.Errorf("decode image: %w", err)
		}
		rgba, ok := img.(*image.RGBA)
		if !ok {
			b := img.Bounds()
			rgba = image.NewRGBA(b)
			draw.Draw(rgba, b, img, b.Min, draw.Src)
		}
		payload = rgba.Pix
		width = rgba.Bounds().Dx()
		height = rgba.Bounds().Dy()
		fmtCode = 32 // RGBA
	}

	return writeKittyChunks(w, payload, fmtCode, width, height, cols, rows)
}

// imageCellSize decides how many terminal cells the image should occupy.
func imageCellSize(data []byte, opts imageRenderOpts) (cols, rows int, err error) {
	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return 0, 0, fmt.Errorf("decode image config: %w", err)
	}

	cellW, cellH := opts.cellW, opts.cellH
	if cellW <= 0 || cellH <= 0 {
		qW, qH, ok := queryCellPixelSize(os.Stdout)
		if ok {
			cellW, cellH = qW, qH
		} else {
			// Reasonable fallback for typical terminals (e.g. 8x16).
			cellW, cellH = 8, 16
		}
	}

	termCols := opts.termCols
	if termCols <= 0 {
		if w, _, terr := term.GetSize(int(os.Stdout.Fd())); terr == nil && w > 0 {
			termCols = w
		} else {
			termCols = 80
		}
	}

	maxCols := opts.maxCols
	if maxCols <= 0 || maxCols > termCols {
		maxCols = termCols
	}

	cols, rows = fitImageCells(cfg.Width, cfg.Height, cellW, cellH, maxCols, opts.maxRows)
	return cols, rows, nil
}

// fitImageCells picks (cols, rows) for an image at the given cell size, fitting
// inside maxCols (>0). maxRows == 0 means unbounded vertically. Aspect ratio is
// preserved; the result is the largest size that fits.
func fitImageCells(imgW, imgH, cellW, cellH, maxCols, maxRows int) (cols, rows int) {
	if imgW <= 0 || imgH <= 0 || cellW <= 0 || cellH <= 0 || maxCols <= 0 {
		return 0, 0
	}
	// Pixel budget if we used the full width.
	pxW := maxCols * cellW
	pxH := imgH * pxW / imgW
	if maxRows > 0 && pxH > maxRows*cellH {
		pxH = maxRows * cellH
		pxW = imgW * pxH / imgH
	}
	cols = (pxW + cellW - 1) / cellW
	rows = (pxH + cellH - 1) / cellH
	if cols < 1 {
		cols = 1
	}
	if rows < 1 {
		rows = 1
	}
	if cols > maxCols {
		cols = maxCols
	}
	if maxRows > 0 && rows > maxRows {
		rows = maxRows
	}
	return cols, rows
}

// writeKittyChunks emits the kitty graphics control sequence(s) for payload.
// Format 100 = PNG (width/height inferred by terminal), 32 = RGBA (width/height
// must be provided). cols/rows controls display size in cells; 0 lets the
// terminal decide.
func writeKittyChunks(w io.Writer, payload []byte, fmtCode, width, height, cols, rows int) error {
	encoded := base64.StdEncoding.EncodeToString(payload)
	const chunkSize = 4096

	first := true
	for len(encoded) > 0 {
		var chunk string
		if len(encoded) > chunkSize {
			chunk = encoded[:chunkSize]
			encoded = encoded[chunkSize:]
		} else {
			chunk = encoded
			encoded = ""
		}
		more := 0
		if len(encoded) > 0 {
			more = 1
		}

		var ctrl strings.Builder
		if first {
			fmt.Fprintf(&ctrl, "a=T,f=%d", fmtCode)
			if fmtCode == 32 {
				fmt.Fprintf(&ctrl, ",s=%d,v=%d", width, height)
			}
			if cols > 0 {
				fmt.Fprintf(&ctrl, ",c=%d", cols)
			}
			if rows > 0 {
				fmt.Fprintf(&ctrl, ",r=%d", rows)
			}
			fmt.Fprintf(&ctrl, ",m=%d", more)
			first = false
		} else {
			fmt.Fprintf(&ctrl, "m=%d", more)
		}

		if _, err := fmt.Fprintf(w, "\x1b_G%s;%s\x1b\\", ctrl.String(), chunk); err != nil {
			return err
		}
	}
	// Cursor is left at the bottom-right of the image; advance to a fresh line.
	if _, err := fmt.Fprint(w, "\n"); err != nil {
		return err
	}
	return nil
}

// queryCellPixelSize asks the terminal for its cell pixel size via the
// "report text area in pixels" control (CSI 16 t -> CSI 6;H;W t). Returns
// (width, height, ok). Falls back to (0,0,false) if stdin isn't a TTY, the
// terminal doesn't respond promptly, or parsing fails.
func queryCellPixelSize(out *os.File) (int, int, bool) {
	in := os.Stdin
	fd := int(in.Fd())
	if !term.IsTerminal(fd) {
		return 0, 0, false
	}
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return 0, 0, false
	}
	defer term.Restore(fd, oldState)

	if _, err := out.WriteString("\x1b[16t"); err != nil {
		return 0, 0, false
	}

	// Read until we see the terminating 't' or run out of bytes. Terminals that
	// don't support this query simply won't reply; the read will block. To stay
	// safe we use a small buffered read with a tight non-blocking loop via the
	// raw fd, but the simplest portable approach is a single Read on the raw
	// stdin, which most kitty-class terminals satisfy near-instantly.
	buf := make([]byte, 32)
	n, err := in.Read(buf)
	if err != nil || n == 0 {
		return 0, 0, false
	}
	resp := string(buf[:n])
	// Expected: ESC [ 6 ; <H> ; <W> t
	i := strings.Index(resp, "\x1b[6;")
	if i < 0 {
		return 0, 0, false
	}
	rest := resp[i+4:]
	end := strings.IndexByte(rest, 't')
	if end < 0 {
		return 0, 0, false
	}
	parts := strings.Split(rest[:end], ";")
	if len(parts) != 2 {
		return 0, 0, false
	}
	var h, wd int
	if _, err := fmt.Sscanf(parts[0], "%d", &h); err != nil {
		return 0, 0, false
	}
	if _, err := fmt.Sscanf(parts[1], "%d", &wd); err != nil {
		return 0, 0, false
	}
	if h <= 0 || wd <= 0 {
		return 0, 0, false
	}
	// h and wd are pixel size of the *whole* text area. Convert to per-cell
	// size using current terminal dimensions.
	cols, rows, terr := term.GetSize(int(out.Fd()))
	if terr != nil || cols <= 0 || rows <= 0 {
		return 0, 0, false
	}
	return wd / cols, h / rows, true
}
