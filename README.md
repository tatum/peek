# peek

`cat` with syntax highlighting, line numbers, and markdown rendering.

## Install

### From source

Requires Go 1.25+.

```sh
go install github.com/tatum/peek@latest
```

### Build locally

```sh
git clone https://github.com/tatum/peek.git
cd peek
go build -o peek .
```

Move the binary somewhere on your `$PATH`:

```sh
mv peek /usr/local/bin/
```

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
```

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

### Environment variables

| Variable | Description |
|----------|-------------|
| `PEEK_THEME` | Default color theme (e.g. `dracula`, `monokai`, `dark`, `light`) |
| `NO_COLOR` | Disable all color output when set |
