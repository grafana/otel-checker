# Releasing otel-checker

Releases are prepared automatically from semantic pull request titles. After
changes land on `main`, [release-please](https://github.com/googleapis/release-please)
opens or updates a draft release pull request. Merging that pull request creates
the next version tag and draft GitHub release. GoReleaser then builds and uploads
archives for Linux, macOS, and Windows on amd64 and arm64, along with SHA-256
checksums and build provenance. Linux artifacts disable cgo and are therefore
usable on both glibc- and musl-based distributions. The workflow publishes the
release after the assets are uploaded.

The release workflow is `.github/workflows/release.yml`, and its build settings
are in `.goreleaser.yml`. To test a release build locally without publishing:

```bash
mise run lint
mise run test
mise exec -- goreleaser release --snapshot --clean
```

Maintainers can republish an existing tag by manually running the Release
workflow from that tag and providing the tag name as the workflow input.
