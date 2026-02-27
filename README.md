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
```

### Flags

| Flag | Description |
|------|-------------|
| `-l`, `--lang` | Force language for syntax highlighting |
| `-t`, `--theme` | Color theme (overrides `PEEK_THEME`) |
| `-p`, `--pager` | Open output in a pager |
| `-n`, `--no-lines` | Hide line numbers |
| `--force-color` | Force color output when piped |

### Environment variables

| Variable | Description |
|----------|-------------|
| `PEEK_THEME` | Default color theme (e.g. `dracula`, `monokai`, `dark`, `light`) |
| `NO_COLOR` | Disable all color output when set |
