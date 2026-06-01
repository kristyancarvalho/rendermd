# mdp — Markdown Pager

Terminal pager for Markdown files with syntax-aware rendering, hot reload, and fuzzy search.

## Features

- Renders headings, bold, italic, inline code, fenced code blocks, blockquotes, lists, task lists, tables, and thematic breaks
- Hot-reload on file save (`--watch`)
- `/` search with `n`/`N` navigation
- Configurable themes (`default` dark, `light`)
- TOML config file at `~/.config/mdp/config.toml`
- Reads from file or stdin

## Keybindings

| Key | Action |
|-----|--------|
| `j` / `k` | Scroll down / up |
| `Ctrl+d` / `Ctrl+u` | Half page down / up |
| `gg` | Top |
| `G` | Bottom |
| `/` | Enter search |
| `n` / `N` | Next / prev result |
| `r` | Reload file |
| `?` | Help overlay |
| `q` | Quit |

## Usage

```sh
# View a file
mdp README.md

# Hot-reload on save
mdp --watch README.md
mdp -w README.md

# Light theme
mdp --theme light README.md

# Custom config
mdp --config ~/.config/mdp/config.toml README.md

# Pipe from stdin
cat README.md | mdp

# Show version
mdp --version
```

## Build & Run

### Prerequisites

- Go 1.22+
- Internet access (for `go mod download`) OR pre-downloaded modules

### Quick Start

```sh
# Clone or unzip the project
cd mdp

# Download dependencies and build
go mod tidy
go build -o mdp ./cmd/mdp

# Run
./mdp README.md
```

### Install globally

```sh
go install ./cmd/mdp
# Binary goes to $GOPATH/bin/mdp or $HOME/go/bin/mdp
```

### Module path

The default module path is `github.com/kristyancarvalho/mdp`. Before publishing,
replace it with your actual GitHub username:

```sh
# macOS/Linux
find . -type f -name "*.go" | xargs sed -i 's|github.com/kristyancarvalho/mdp|github.com/YOU/mdp|g'
sed -i 's|github.com/kristyancarvalho/mdp|github.com/YOU/mdp|g' go.mod
```

## Configuration

Create `~/.config/mdp/config.toml`:

```toml
[ui]
padding = 2
max_width = 96
soft_wrap = true
show_line_numbers = false
show_urls = false

[theme]
name = "default"   # "default" | "light"

# Override individual colours (hex or ANSI 256):
# heading  = "#89b4fa"
# text     = "#cdd6f4"
# code_bg  = "#313244"

[markdown]
hide_syntax      = true
render_tables    = true
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

## Project Structure

```
mdp/
├── cmd/mdp/main.go              # Entry point
├── internal/
│   ├── app/app.go               # CLI flags, wiring
│   ├── config/
│   │   ├── config.go            # TOML loader & merge
│   │   └── defaults.go          # Default values
│   ├── input/
│   │   ├── keymap.go            # Key→Action resolver
│   │   └── state.go             # UI state enum
│   ├── layout/
│   │   └── engine.go            # Markdown→[]Line layout engine
│   ├── markdown/
│   │   ├── parser.go            # goldmark entry point
│   │   └── normalize.go         # AST→model converter
│   ├── model/
│   │   ├── block.go             # Block node types
│   │   └── span.go              # Inline span types
│   ├── render/
│   │   └── renderer.go          # lipgloss renderer + viewport
│   ├── theme/
│   │   ├── theme.go             # Theme struct
│   │   └── builtin.go           # default / light palettes
│   ├── ui/
│   │   └── ui.go                # Bubbletea model (Update/View)
│   └── watch/
│       └── watcher.go           # fsnotify debounced watcher
├── go.mod
└── README.md
```

## Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/charmbracelet/bubbletea` | TUI framework |
| `github.com/charmbracelet/lipgloss` | Terminal styling |
| `github.com/yuin/goldmark` | Markdown parser (CommonMark + GFM) |
| `github.com/fsnotify/fsnotify` | File-system watcher |
| `github.com/BurntSushi/toml` | Config file parser |
