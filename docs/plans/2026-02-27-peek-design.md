# peek — pretty file viewer for the terminal

A Go CLI that renders any text file with color and formatting. Markdown gets full rendering via glamour; source code gets syntax highlighting via chroma; everything else gets clean plain text with line numbers.

## CLI Interface

```
peek [flags] [file...]

Flags:
  -p, --pager         Open in pager mode
  -l, --lang LANG     Force language (for stdin or override detection)
  -t, --theme THEME   Override color theme (default: from PEEK_THEME env)
  -n, --no-lines      Hide line numbers on code files
      --force-color   Force color output even when piped

Env:
  PEEK_THEME          Default theme (e.g. dracula, monokai, github-dark)
  NO_COLOR            Disable all color output
```

Examples:
```
peek README.md                    # render markdown
peek main.go                     # syntax-highlighted code
peek -p main.go                  # same, but in pager
cat data.json | peek -l json     # stdin with explicit lang
PEEK_THEME=dracula peek app.py   # themed output
```

## Architecture

Three internal layers:

1. **CLI layer** — parses args via stdlib `flag`, reads stdin or file paths, detects TTY
2. **Detector** — determines file type from extension, maps to chroma lexer name
3. **Renderer** — routes to the right renderer:
   - Markdown -> glamour
   - Source code -> chroma (syntax highlighting + line numbers)
   - JSON/YAML -> chroma (with their respective lexers)
   - Plain text -> line numbers + optional soft wrap

## Rendering Pipeline

```
file arg or stdin
       |
       v
  detect type (extension map -> chroma lexer name)
       |
       v
  markdown?   -> glamour render
  known lang? -> chroma highlight + line numbers
  otherwise   -> plain text + line numbers
       |
       v
  TTY check:
    TTY + no -p  -> print to stdout
    TTY + -p     -> pipe to $PAGER or less -R
    not TTY      -> strip color (unless --force-color)
```

## Pager

Check `$PAGER` env first. If unset, default to `less -R`. No custom pager — shelling out to less is what bat and glow both do.

## Theme System

Single `PEEK_THEME` env var maps to both glamour style and chroma style. Override per-invocation with `--theme`. Chroma ships ~40 built-in styles.

## Project Structure

```
peek/
  main.go           # CLI entry point, flag parsing
  detect.go         # file type detection from extension
  render.go         # renderer dispatch (markdown, code, plain)
  render_md.go      # glamour wrapper
  render_code.go    # chroma wrapper + line numbers
  pager.go          # pager output logic
  theme.go          # PEEK_THEME -> glamour/chroma style mapping
  go.mod
  go.sum
```

Dependencies: glamour, chroma, lipgloss (for line number gutter styling).

## Error Handling

- File not found -> stderr message, exit 1
- Permission denied -> stderr message, exit 1
- Binary file detected -> "peek: binary file, not rendering", exit 0
- Unknown extension + no `-l` -> render as plain text (never fail)

## Testing

- Unit tests for detect.go (extension -> language mapping)
- Unit tests for theme.go (env var -> style resolution)
- Golden file tests: render known .py, .md, .json files and snapshot output
