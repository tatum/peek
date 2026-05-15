package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"
	"testing"
)

func TestImageFormatFromExt(t *testing.T) {
	cases := map[string]imageFormat{
		".png":  imageFormatPNG,
		".jpg":  imageFormatJPEG,
		".jpeg": imageFormatJPEG,
		".gif":  imageFormatGIF,
		".webp": imageFormatWebP,
		".txt":  imageFormatNone,
		"":      imageFormatNone,
	}
	for ext, want := range cases {
		if got := imageFormatFromExt(ext); got != want {
			t.Errorf("imageFormatFromExt(%q) = %v, want %v", ext, got, want)
		}
	}
}

func TestImageFormatFromMagicPNG(t *testing.T) {
	data := []byte("\x89PNG\r\n\x1a\n" + strings.Repeat("x", 16))
	if got := imageFormatFromMagic(data); got != imageFormatPNG {
		t.Errorf("PNG magic: got %v", got)
	}
}

func TestImageFormatFromMagicJPEG(t *testing.T) {
	data := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10}
	if got := imageFormatFromMagic(data); got != imageFormatJPEG {
		t.Errorf("JPEG magic: got %v", got)
	}
}

func TestImageFormatFromMagicGIF(t *testing.T) {
	for _, sig := range []string{"GIF87a", "GIF89a"} {
		data := append([]byte(sig), 0x00, 0x00)
		if got := imageFormatFromMagic(data); got != imageFormatGIF {
			t.Errorf("GIF magic %q: got %v", sig, got)
		}
	}
}

func TestImageFormatFromMagicWebP(t *testing.T) {
	data := []byte{'R', 'I', 'F', 'F', 0, 0, 0, 0, 'W', 'E', 'B', 'P', 0, 0}
	if got := imageFormatFromMagic(data); got != imageFormatWebP {
		t.Errorf("WebP magic: got %v", got)
	}
}

func TestImageFormatFromMagicNone(t *testing.T) {
	if got := imageFormatFromMagic([]byte("hello world")); got != imageFormatNone {
		t.Errorf("non-image magic: got %v", got)
	}
	if got := imageFormatFromMagic(nil); got != imageFormatNone {
		t.Errorf("nil magic: got %v", got)
	}
}

func TestDetectFileTypeImage(t *testing.T) {
	for _, name := range []string{"a.png", "B.JPG", "c.jpeg", "d.gif", "e.webp"} {
		if got := detectFileType(name); got != fileTypeImage {
			t.Errorf("detectFileType(%q) = %v, want fileTypeImage", name, got)
		}
	}
}

func TestFitImageCellsAspectRatio(t *testing.T) {
	// 200x100 image in an 8x16 cell, max 20 cols, no row limit.
	// pxW = 160, pxH = 80, rows = ceil(80/16) = 5.
	cols, rows := fitImageCells(200, 100, 8, 16, 20, 0)
	if cols != 20 {
		t.Errorf("cols = %d, want 20", cols)
	}
	if rows != 5 {
		t.Errorf("rows = %d, want 5", rows)
	}
}

func TestFitImageCellsRowCapped(t *testing.T) {
	// Tall image: 100x500, cell 8x16. Unconstrained rows would be huge.
	// Cap rows at 10 -> pxH = 160, pxW = 32, cols = 4.
	cols, rows := fitImageCells(100, 500, 8, 16, 80, 10)
	if rows != 10 {
		t.Errorf("rows = %d, want 10", rows)
	}
	if cols != 4 {
		t.Errorf("cols = %d, want 4", cols)
	}
}

func TestFitImageCellsZeros(t *testing.T) {
	if c, r := fitImageCells(0, 0, 8, 16, 80, 0); c != 0 || r != 0 {
		t.Errorf("zero dims: cols=%d rows=%d, want 0,0", c, r)
	}
}

func TestKittySupportEnv(t *testing.T) {
	saved := map[string]string{}
	keys := []string{"TERM", "KITTY_WINDOW_ID", "TERM_PROGRAM", "GHOSTTY_RESOURCES_DIR"}
	for _, k := range keys {
		saved[k] = os.Getenv(k)
		os.Unsetenv(k)
	}
	t.Cleanup(func() {
		for k, v := range saved {
			if v == "" {
				os.Unsetenv(k)
			} else {
				os.Setenv(k, v)
			}
		}
	})

	if kittySupportEnv() {
		t.Errorf("expected no support with all env unset")
	}
	os.Setenv("TERM", "xterm-kitty")
	if !kittySupportEnv() {
		t.Errorf("expected support for TERM=xterm-kitty")
	}
	os.Unsetenv("TERM")
	os.Setenv("KITTY_WINDOW_ID", "1")
	if !kittySupportEnv() {
		t.Errorf("expected support for KITTY_WINDOW_ID")
	}
	os.Unsetenv("KITTY_WINDOW_ID")
	os.Setenv("TERM_PROGRAM", "ghostty")
	if !kittySupportEnv() {
		t.Errorf("expected support for TERM_PROGRAM=ghostty")
	}
	os.Setenv("TERM_PROGRAM", "Apple_Terminal")
	if kittySupportEnv() {
		t.Errorf("did not expect support for Apple_Terminal")
	}
}

// makePNG returns a tiny valid PNG (w x h, solid color) as bytes.
func makePNG(t *testing.T, w, h int) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, color.RGBA{200, 100, 50, 255})
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func TestRenderImageKittyPNGPassthrough(t *testing.T) {
	data := makePNG(t, 40, 20)
	var buf bytes.Buffer
	err := renderImageKitty(&buf, data, imageFormatPNG, imageRenderOpts{
		cellW:    8,
		cellH:    16,
		termCols: 80,
		maxCols:  10,
	})
	if err != nil {
		t.Fatalf("renderImageKitty: %v", err)
	}
	out := buf.String()
	if !strings.HasPrefix(out, "\x1b_G") {
		t.Errorf("expected output to start with kitty APC prefix, got %q", out[:min(20, len(out))])
	}
	if !strings.Contains(out, "f=100") {
		t.Errorf("expected f=100 (PNG) in control data")
	}
	if !strings.Contains(out, "c=10") {
		t.Errorf("expected c=10 in control data, got: %s", out)
	}
	if !strings.HasSuffix(out, "\n") {
		t.Errorf("expected trailing newline")
	}
}

func TestRenderImageKittyRGBADecode(t *testing.T) {
	// Use a PNG but force the RGBA path by claiming it's JPEG.
	data := makePNG(t, 8, 8)
	var buf bytes.Buffer
	err := renderImageKitty(&buf, data, imageFormatJPEG, imageRenderOpts{
		cellW:    8,
		cellH:    16,
		termCols: 80,
		maxCols:  4,
	})
	if err != nil {
		t.Fatalf("renderImageKitty: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "f=32") {
		t.Errorf("expected f=32 (RGBA) in control data")
	}
	if !strings.Contains(out, "s=8") || !strings.Contains(out, "v=8") {
		t.Errorf("expected s=8,v=8 in control data, got: %s", out)
	}
}
