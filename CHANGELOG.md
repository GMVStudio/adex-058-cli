# Changelog

## [0.2.7] - 2026-07-03

### Added
- `adex tenant use <ID>` command — set a default tenant so `--tenant` is no longer required on every command
- `--tenant` flag on `adex init` — optionally set the default tenant during initial credential binding
- `TenantID` field in config (`~/.adex/config.json`) with `ADEX_TENANT_ID` env var override
- `Factory.resolveTenant` — resolves `--tenant` flag with fallback to config default, replacing the hard-required `requireTenant`
- Skills updated with 3-step init workflow and `tenant use` documentation

### Changed
- `--tenant` flag is now optional on all KS/OE commands (was previously required); falls back to config default tenant
- Help text and examples for `adex ks`, `adex oe`, and report commands updated to omit `--tenant` when default is set
- `adex init` output includes `tenant_id` when set, or a `next_step` hint when no default tenant is configured
- `adex-shared` skill: init flow rewritten as 3-step process, tenant reference expanded with `tenant use` section, user reference updated

## [0.2.6] - 2026-07-02

### Added
- `ADEX_CONFIG_DIR` env var — override config directory for sandbox/CI environments
- `ADEX_AUTHORIZATION` env var documented in root help text and error hints
- Config dir fallback to `.adex` when `UserHomeDir()` returns empty
- Actionable recovery hints on config save failures (mention `ADEX_CONFIG_DIR` and env var alternatives)
- Test asserting hint mentions `ADEX_AUTHORIZATION` for sandbox recovery
- Test asserting config save errors include `ADEX_CONFIG_DIR` hint and preserve cause

### Fixed
- `fallbackFullInstall` now installs only official skills via `InstallSkill` instead of blindly calling `InstallAllSkills`
- Error wrapping in `skillscheck/state.go` uses `%w: %w` for proper error chain preservation
- `resolveExe` moved to platform-specific files (was incorrectly in generic `updater.go`)
- Removed stale `-g` flag from `npx skills add` hint in update command
- CI workflow version pins updated

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
