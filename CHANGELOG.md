# Changelog

## [0.2.5] - 2026-07-02

### Added
- `adex update` command — self-update CLI binary via npm and sync embedded skills
- `internal/update` package with full update lifecycle (npm update, skills sync, binary replacement)
- Update notice (`_notice.update`) in JSON output when a newer version is available
- Skills sync notice (`_notice.skills`) when locally installed skills are out of sync
- `internal/output` support for notice fields in JSON envelope
- Error taxonomy extended with new subtypes for update-related failures
- `adex-shared` skill updated with `adex update` command documentation

### Fixed
- Bug in `runList` / `runSingle` check logic causing incorrect error handling
- CI workflow adjustment for test execution

## [0.2.4] - 2026-07-01

### Added
- `adex tenant` command — list tenants with name/status filters and pagination
- `adex user` command — get current authenticated user info
- `adex init` command — bind API credentials (one-time setup)
- `adex oe` command group — Oceanengine (巨量) advertising data: accounts, projects, units, reports, metric meta, dashboard, budget-vs-actual
- Embedded skills expanded: `adex-ks`, `adex-oe` (full command coverage)
- `--jq` flag for jq expression filtering on all commands
- `--page-all` flag for automatic multi-page aggregation
- `--range` / `--begin` / `--end` flexible date range specification
- `--dry-run` flag for request inspection without API calls
- Config file persistence (`~/.adex/config.json`) with `ADEX_AUTHORIZATION` env override
- `daterange` package for relative (7d/4w/1m) and explicit date parsing
- `paginate` package for page-all aggregation
- `output/jq.go` for gojq-based JSON filtering
- Shared command builders (`report_shared.go`) for KS and OE reuse

### Changed
- `adex-shared` skill updated with full command tree, init flow, shared flags reference, and skill routing
- `skills_embed.go` now embeds `skills/*/references` directories
- README updated with complete command tree, skills table, and architecture diagram
- `runSingle` extended to accept table columns for single-object commands

## [0.2.3] - 2026-06-27

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
