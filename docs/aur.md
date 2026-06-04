# Arch User Repository

The Arch package metadata lives in `packaging/aur`.

## Package Name

The package is named `rendermd` and installs the `rendermd` binary at `/usr/bin/rendermd`.

## Local Build

```sh
cd packaging/aur
makepkg -si
```

`makepkg` downloads the pinned upstream source, builds the Go binary, installs the license, and installs `config.example.toml` under `/usr/share/doc/rendermd`.

## AUR Publication

The upstream repository and the release tag used by the package must be publicly fetchable before submission. AUR builders cannot fetch private GitHub repositories.

Copy `PKGBUILD` and `.SRCINFO` into the AUR package repository for `rendermd`, commit them there, and push to `ssh://aur@aur.archlinux.org/rendermd.git`.
