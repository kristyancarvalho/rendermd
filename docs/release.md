# Release Process

`rendermd` publishes GitHub Releases from version tags.

## Requirements

- Go 1.26 or newer
- A clean `dev` branch
- Passing local tests
- Permission to push version tags

## Steps

1. Update release notes or documentation as needed.
2. Run `go test ./...`.
3. Run `make build`.
4. Confirm `./bin/rendermd --version` prints the expected version metadata.
5. Create and push a version tag:

```sh
git tag v1.2.3
git push origin v1.2.3
```

The `release` workflow runs on `v*` tags. It tests the project, builds Linux, macOS, and Windows artifacts, generates `checksums.txt`, and publishes the files to a GitHub Release.

After publishing a release, update `packaging/aur/PKGBUILD` and `packaging/aur/.SRCINFO` if the AUR package needs a version bump.
