# Changelog

## [0.2.0] - 2026-06-27

### Added
- npm distribution via `@gmvstudio/adex-cli` package
- Skill system with `adex skills list` and `adex skills read` commands
- Interactive install wizard (`adex install`)
- Cross-platform binary distribution (darwin/linux/windows × amd64/arm64)
- GitHub Actions release workflow (goreleaser + npm publish)
- `BUILD.md` and `RELEASE.md` documentation
- `scripts/release.sh` release helper script

### Changed
- Repository URL corrected to match GitHub repo name
- npm registry forced to HTTPS
- goreleaser v2 `formats` syntax adopted
