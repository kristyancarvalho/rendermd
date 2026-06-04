# rendermd

[![Release workflow](https://github.com/kristyancarvalho/rendermd/actions/workflows/release.yml/badge.svg)](https://github.com/kristyancarvalho/rendermd/actions/workflows/release.yml)
[![License: MIT](https://img.shields.io/github/license/kristyancarvalho/rendermd)](LICENSE)
[![Go version](https://img.shields.io/github/go-mod/go-version/kristyancarvalho/rendermd)](go.mod)
[![Latest release](https://img.shields.io/github/v/release/kristyancarvalho/rendermd)](https://github.com/kristyancarvalho/rendermd/releases)
[![Active milestones](https://img.shields.io/badge/milestones-active-2563eb)](https://github.com/kristyancarvalho/rendermd/milestones)

`rendermd` is a terminal Markdown renderer built for reading local Markdown files quickly, clearly, and without leaving the command line.

It supports syntax-aware rendering, hot reload, search navigation, configurable themes, and file or stdin input.

## Features

- Markdown rendering for headings, emphasis, inline code, fenced code blocks, blockquotes, lists, task lists, tables, and thematic breaks
- Syntax highlighting for common code blocks
- Hot reload on file save with debounced watcher events
- Search with next and previous result navigation
- Built-in `default` and `light` themes
- TOML configuration at `~/.config/rendermd/config.toml`
- Build-time version, commit, and date metadata

## Installation

### From Source

Requirements:

- Go 1.26 or newer
- Git

```sh
git clone https://github.com/kristyancarvalho/rendermd.git
cd rendermd
make build
install -Dm755 bin/rendermd ~/.local/bin/rendermd
```

### With Go

```sh
go install github.com/kristyancarvalho/rendermd/cmd/rendermd@latest
```

This installs the binary into `$GOBIN`, or `$GOPATH/bin` when `GOBIN` is not set.

### Arch Linux

The AUR package metadata lives in `packaging/aur`.

```sh
cd packaging/aur
makepkg -si
```

After the package is published to AUR:

```sh
paru -S rendermd
```

## Usage

```sh
rendermd README.md
rendermd --watch README.md
rendermd --theme light README.md
rendermd --config ~/.config/rendermd/config.toml README.md
cat README.md | rendermd
rendermd --version
```

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
| `j` / `k` | Scroll down or up |
| `Ctrl+d` / `Ctrl+u` | Half page down or up |
| `gg` | Go to top |
| `G` | Go to bottom |
| `/` | Enter search |
| `n` / `N` | Next or previous result |
| `r` | Reload file |
| `?` | Toggle help |
| `q` | Quit |

## Configuration

Create `~/.config/rendermd/config.toml` or copy `config.example.toml`.

```toml
[ui]
padding           = 2
line_spacing      = 0
scrolloff         = 4
soft_wrap         = true
max_width         = 96
show_line_numbers = false
show_urls         = false
mouse             = true

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

## Rendering Styles

Headings use the configured heading color and keep their text aligned to the document content width. Blockquotes use the quote foreground and background across the full rendered quote line, including wrapped lines and nested quote markers. Code blocks use the code background for both code text and horizontal padding, with language labels shown when syntax display is enabled. Lists keep markers and continuation text aligned within layout segments. Thematic breaks use the border color and span the available content width.

## Search

Search uses an index built from rendered lines and refreshed after layout changes, manual reloads, and watched file updates. The index stores one lowercased string per rendered line, trading a small amount of memory for faster repeated searches in large documents.

## Development

```sh
go test ./...
make build
./bin/rendermd --version
```

The main packages are:

| Path | Purpose |
|------|---------|
| `cmd/rendermd` | CLI entrypoint |
| `internal/app` | Flag parsing and runtime wiring |
| `internal/config` | Defaults, TOML loading, validation, and merge logic |
| `internal/markdown` | Markdown parsing and normalization |
| `internal/layout` | Markdown block layout |
| `internal/render` | Terminal rendering |
| `internal/theme` | Built-in themes and validation |
| `internal/ui` | Bubble Tea model |
| `internal/watch` | File watcher with debounced events |

## Project Planning

Active work is tracked in [GitHub milestones](https://github.com/kristyancarvalho/rendermd/milestones). Use the milestone list to find upcoming releases, open tasks, and the current release focus.

## Contributing

Contributions are welcome. Start with [CONTRIBUTING.md](CONTRIBUTING.md), then open an issue or pull request against the current milestone.

## Releases

Releases are created from `v*` tags by the GitHub Actions workflow in `.github/workflows/release.yml`. See [docs/release.md](docs/release.md) for the release checklist and artifact details.

## License

`rendermd` is released under the [MIT License](LICENSE).
