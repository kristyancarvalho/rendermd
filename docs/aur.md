# Arch User Repository

The Arch package metadata lives in `packaging/aur`.

## Package Name

The package is named `mdp-pager` because Arch Linux already provides a different `mdp` package for a Markdown presentation tool. The package installs this project as `/usr/bin/mdp` and declares a conflict with Arch's `mdp` package.

## Local Build

```sh
cd packaging/aur
makepkg -si
```

`makepkg` downloads the pinned upstream source, builds the Go binary, installs the license, and installs `config.example.toml` under `/usr/share/doc/mdp-pager`.

## AUR Publication

The upstream repository must be publicly fetchable before submission. AUR builders cannot fetch private GitHub repositories.

Copy `PKGBUILD` and `.SRCINFO` into the AUR package repository for `mdp-pager`, commit them there, and push to `ssh://aur@aur.archlinux.org/mdp-pager.git`.
