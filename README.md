# peek

`cat` with syntax highlighting, line numbers, markdown rendering, and inline
image previews in terminals that speak the
[kitty graphics protocol](https://sw.kovidgoyal.net/kitty/graphics-protocol/).

## Install

### Prebuilt binary (no Go required)

```sh
curl -fsSL https://raw.githubusercontent.com/tatum/peek/main/install.sh | sh
```

Downloads a static binary for your platform (linux/amd64, linux/arm64,
darwin/arm64) from the rolling [`latest`
release](https://github.com/tatum/peek/releases/tag/latest) into
`~/.local/bin` — no sudo, no toolchain. Set `PEEK_INSTALL_DIR` to install
elsewhere. Binaries are rebuilt by CI on every push to `main`.

### From source

Requires Go 1.25+.

```sh
go install github.com/tatum/peek@latest
```

### From a clone

```sh
git clone https://github.com/tatum/peek.git
cd peek
make install
```

This runs `go install .`, which puts the binary in `$(go env GOPATH)/bin`
(usually `~/go/bin`) — make sure that's on your `$PATH`. No sudo needed.
`make uninstall` removes it.

To just build a local binary without installing, run `make` (or
`go build -o peek .`).

## Usage

```sh
peek file.go              # syntax-highlighted source with line numbers
peek README.md            # rendered markdown
cat data.json | peek      # read from stdin
peek -l python script     # force language detection
peek -p main.go           # open in pager
peek -n file.go           # hide line numbers
peek -t dracula file.go   # choose a color theme
peek -z                   # pick a file with fzf, then peek it
peek -zp                  # pick with fzf, then open in pager
peek logo.png             # render image inline (kitty / ghostty / WezTerm)
peek --width 40 logo.png  # cap image width at 40 terminal columns
```

### Viewing images

`peek` recognises `.png`, `.jpg`, `.jpeg`, `.gif`, and `.webp` files (also via
magic-byte sniffing for stdin and extension-less paths) and renders them
inline using the kitty graphics protocol. This works in
[kitty](https://sw.kovidgoyal.net/kitty/), [ghostty](https://ghostty.org/), and
WezTerm. PNGs are sent as-is; other formats are decoded to RGBA via the Go
standard library and sent as raw pixels.

Images are sized to fit the terminal while preserving aspect ratio. Use
`--width N` and `--height N` to cap the cell footprint, and `--force-kitty`
when running under tmux passthrough or another setup the env sniff doesn't
recognise. Images are skipped (with a message) in pager mode and when stdout
is not a TTY.

### Picking files with fzf

If [fzf](https://github.com/junegunn/fzf) is installed, `-z` / `--fzf` drops
you into an fzf picker with a live preview pane rendered by `peek` itself.
Arrow through the list to see each file syntax-highlighted in the preview,
press enter to peek the selection, or Esc to cancel.

Other flags are forwarded into the preview, so `peek -z -t dracula -n` shows
the preview in the dracula theme with line numbers hidden, and `peek -zp`
opens the selected file in your pager after you pick it — handy for large
files.

The file list comes from fzf's default source, so `$FZF_DEFAULT_COMMAND` is
respected if set (e.g. `fd --type f` or `rg --files`).

### Flags

| Flag | Description |
|------|-------------|
| `-l`, `--lang` | Force language for syntax highlighting |
| `-t`, `--theme` | Color theme (overrides `PEEK_THEME`) |
| `-p`, `--pager` | Open output in a pager |
| `-n`, `--no-lines` | Hide line numbers |
| `-z`, `--fzf` | Pick a file with fzf, then peek it (requires `fzf` on `$PATH`) |
| `--force-color` | Force color output when piped |
| `--force-kitty` | Render images via kitty protocol even when env detection fails |
| `--width` | Max image width in terminal columns (0 = full terminal width) |
| `--height` | Max image height in terminal rows (0 = unbounded) |

### Environment variables

| Variable | Description |
|----------|-------------|
| `PEEK_THEME` | Default color theme (e.g. `dracula`, `monokai`, `dark`, `light`) |
| `NO_COLOR` | Disable all color output when set |
