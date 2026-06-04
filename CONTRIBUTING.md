# Contributing

Thanks for helping improve `rendermd`. This guide explains how to propose changes, prepare code, and keep contributions easy to review.

## Ways To Contribute

- Fix bugs in Markdown parsing, layout, rendering, configuration, input handling, or file watching.
- Improve tests, benchmarks, and validation coverage.
- Improve documentation, examples, release notes, and packaging metadata.
- Propose focused enhancements through GitHub issues before opening larger pull requests.

## Before You Start

1. Check the active milestones at <https://github.com/kristyancarvalho/rendermd/milestones>.
2. Search existing issues and pull requests for related work.
3. Open an issue for non-trivial changes so the scope is clear before implementation.
4. Keep changes small enough to review in one pass.

## Development Setup

Requirements:

- Go 1.26 or newer
- Git
- A terminal that supports ANSI styling

```sh
git clone https://github.com/kristyancarvalho/rendermd.git
cd rendermd
go test ./...
make build
./bin/rendermd --version
```

## Project Structure

| Path | Purpose |
|------|---------|
| `cmd/rendermd` | CLI entrypoint |
| `internal/app` | CLI flags and runtime wiring |
| `internal/config` | Defaults, TOML loading, validation, and merge logic |
| `internal/input` | Key bindings and input state |
| `internal/layout` | Markdown block layout |
| `internal/markdown` | Markdown parsing and normalization |
| `internal/model` | Shared document model |
| `internal/render` | Terminal rendering |
| `internal/syntax` | Syntax tokenization |
| `internal/theme` | Built-in themes and validation |
| `internal/ui` | Bubble Tea model |
| `internal/watch` | Debounced file watcher |

## Coding Guidelines

- Preserve existing behavior unless the issue explicitly calls for a behavior change.
- Prefer simple, local changes over broad rewrites.
- Keep package boundaries consistent with the current `internal` layout.
- Use existing helper types and test patterns before adding new abstractions.
- Keep user-facing errors clear and actionable.
- Do not commit generated binaries, coverage files, local packages, or `/specs`.

## Tests

Run the full test suite before opening a pull request:

```sh
go test ./...
```

For changes that affect build or release behavior, also run:

```sh
make build
./bin/rendermd --version
```

For AUR packaging changes, regenerate `.SRCINFO` from `packaging/aur`:

```sh
cd packaging/aur
makepkg --printsrcinfo > .SRCINFO
```

## Commits

Use concise commit messages in this format:

```text
<type>/<scope> [message]; <action> (issue <id>)
```

Examples:

```text
feat/ui add dark mode support; implements (issue 123)
fix/layout correct line wrapping; fixes (issue 124)
chore/docs update contributor guide; implements (issue 125)
```

Use a type such as `feat`, `fix`, `docs`, `test`, `chore`, or `refactor`. Use a scope that names the affected area, such as `ui`, `layout`, `docs`, `tests`, `config`, `render`, or `packaging`.

## Pull Requests

Before opening a pull request:

1. Rebase or merge the latest target branch.
2. Run relevant tests and builds.
3. Update documentation when behavior, flags, config, packaging, or release steps change.
4. Link the issue being implemented or fixed.
5. Describe validation performed in the pull request body.

Pull requests should be focused, reviewable, and aligned with an open issue or active milestone.

## Documentation

Documentation changes should stay concise and practical. Update `README.md` for user-facing workflows, `docs/release.md` for release steps, `docs/aur.md` for Arch packaging, and package-level tests or examples when behavior changes.

## Reporting Bugs

When reporting a bug, include:

- `rendermd --version` output
- Operating system and terminal
- The command that failed
- A minimal Markdown sample when possible
- Expected behavior and actual behavior

## Security

Do not disclose security-sensitive issues in public issues. Open a private report or contact the maintainer directly with reproduction details and impact.
