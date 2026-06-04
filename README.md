# mdp

`mdp` is a terminal Markdown pager with syntax-aware rendering, hot reload, fuzzy search, and configurable themes.

## Features

- Markdown rendering for headings, emphasis, inline code, fenced code blocks, blockquotes, lists, task lists, tables, and thematic breaks
- Hot reload on file save with debounced watcher events
- Search with next and previous result navigation
- Built-in `default` and `light` themes
- TOML configuration at `~/.config/mdp/config.toml`
- File and stdin input
- Build-time version, commit, and date metadata

## Installation

### From Source

Requirements:

- Go 1.26 or newer
- Git

```sh
git clone https://github.com/kristyancarvalho/mdp.git
cd mdp
make build
install -Dm755 bin/mdp ~/.local/bin/mdp
```

### With `go install`

```sh
go install github.com/kristyancarvalho/mdp/cmd/mdp@latest
```

This installs the binary into `$GOBIN`, or `$GOPATH/bin` when `GOBIN` is not set.

### Arch Linux

Arch users can build the AUR package metadata from `packaging/aur`.

```sh
cd packaging/aur
makepkg -si
```

After the package is submitted to AUR:

```sh
paru -S mdp-pager
```

The AUR package name is `mdp-pager` because Arch already ships a different package named `mdp`. The package installs this project as `/usr/bin/mdp` and conflicts with Arch's `mdp` package.

## Usage

```sh
mdp README.md
mdp --watch README.md
mdp -w README.md
mdp --theme light README.md
mdp --config ~/.config/mdp/config.toml README.md
cat README.md | mdp
mdp --version
```

`--version` prints the application version, commit hash, and build date. Builds created with `make build`, `make install`, or the release workflow inject those values with linker flags.

## Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--watch` | `-w` | Reload the file after debounced filesystem changes |
| `--config <path>` | `-c <path>` | Load a TOML configuration file |
| `--theme <name>` | `-t <name>` | Select a built-in theme and override the config theme |
| `--version` | | Print version metadata and exit |
| `--help` | | Print CLI help |

## Keybindings

| Key | Action |
|-----|--------|
| `j` / `k` | Scroll down / up |
| `Ctrl+d` / `Ctrl+u` | Half page down / up |
| `gg` | Go to top |
| `G` | Go to bottom |
| `/` | Enter search |
| `n` / `N` | Next / previous result |
| `r` | Reload file |
| `?` | Toggle help |
| `q` | Quit |

## Configuration

Create `~/.config/mdp/config.toml` or copy `config.example.toml`.

```toml
[ui]
padding           = 2
line_spacing      = 0
scrolloff         = 4
soft_wrap         = true
max_width         = 96
show_line_numbers = false
show_urls         = false

[theme]
name       = "default"
background = ""
text       = ""
muted      = ""
heading    = ""
accent     = ""
link       = ""
link_url   = ""
code_bg    = ""
quote_bg   = ""
border     = ""

[markdown]
hide_syntax       = true
render_emphasis   = true
render_strong     = true
render_links      = true
render_images     = false
render_tables     = true
render_task_lists = true

[keys]
up        = "k"
down      = "j"
half_up   = "ctrl+u"
half_down = "ctrl+d"
top       = "g"
bottom    = "G"
search    = "/"
next_hit  = "n"
prev_hit  = "N"
reload    = "r"
quit      = "q"
help      = "?"

[watch]
enabled     = true
debounce_ms = 150
```

Theme colors accept hex colors or ANSI 256-color values. Invalid theme names fall back to `default`; invalid color overrides are ignored with a warning.

## Development

```sh
go test ./...
make build
./bin/mdp --version
```

The main packages are:

| Path | Purpose |
|------|---------|
| `cmd/mdp` | CLI entrypoint |
| `internal/app` | Flag parsing and runtime wiring |
| `internal/config` | Defaults, TOML loading, validation, and merge logic |
| `internal/markdown` | Markdown parsing and normalization |
| `internal/layout` | Markdown block layout |
| `internal/render` | Terminal rendering |
| `internal/theme` | Built-in themes and validation |
| `internal/ui` | Bubble Tea model |
| `internal/watch` | File watcher with debounced events |

## Releases

Releases are created from `v*` tags by the GitHub Actions workflow in `.github/workflows/release.yml`. See `docs/release.md` for the release checklist and artifact details.
