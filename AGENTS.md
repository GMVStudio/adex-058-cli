# AGENTS.md

## What this is

`adex` is a CLI for querying ADEX advertising data (Kuaishou `ks` and
Oceanengine `oe`). Its primary consumers are **AI agents** and shell scripts,
so output format, flag design, and error structure directly affect agent
success rates.

The one rule to internalize: **every error message you write will be parsed by
a machine to decide its next action.** Make errors structured, actionable, and
specific.

## Build & Test

```bash
make build        # build the adex binary
make unit-test    # go test -race ./... (required before PR)
make test         # vet + fmt-check + unit-test
make lint         # golangci-lint (pinned to CI version)
make tidy-check   # fails if go.mod/go.sum are not tidy
```

## Pre-PR checks (match CI gates)

1. `make unit-test`
2. `go vet ./...`
3. `gofmt -l .` — must produce no output
4. `go mod tidy` — must not change `go.mod`/`go.sum`
5. `make lint`

## Source layout

| Path | What it does |
|------|-------------|
| `main.go` | Entry point → `cmd.Execute()` |
| `skills_embed.go` | Embeds `skills/*` and wires them into the CLI |
| `cmd/root.go` | Root command, `Factory`, `handleError`, client/format resolution |
| `cmd/common.go` | Shared flag registration + `runList` / `runSingle` / `dryRun` |
| `cmd/report_shared.go` | Platform-agnostic `top` / `get` / `daily` / `summary` builders |
| `cmd/columns.go` | Table column definitions per resource |
| `cmd/ks_*.go`, `cmd/oe_*.go` | Per-resource command wiring |
| `errs/errs.go` | Typed error taxonomy + JSON envelope + exit codes |
| `internal/client/client.go` | HTTP client: `Do`, `DoTyped`, typed error classification |
| `internal/config/config.go` | Config load/save (file + env overlay) |
| `internal/selfupdate/` | Installation detection, npm update, skills sync via npx, binary replacement |
| `internal/skillscheck/` | Skills version check, stale notice, incremental sync planning, state persistence |
| `internal/daterange/` | `--range` / `--begin` / `--end` resolution |
| `internal/output/` | `--format` (json/pretty/table) and `--jq` rendering |
| `internal/paginate/` | `--page-all` token aggregation |
| `internal/vfs/` | Filesystem abstraction (`vfs.Default` / `vfs.OS`); test-mockable |

## Code conventions

### Structured errors

Command-facing failures must be typed `errs.*` errors, never a bare
`fmt.Errorf` / `errors.New`. Agents branch on the stderr envelope's
`type` / `subtype` / `param` / `hint` fields.

| Failure | Constructor |
|---------|-------------|
| User flag/arg fails validation | `errs.NewValidationError(errs.SubtypeInvalidArgument, ...).WithParam("--flag")` |
| Missing local config / not bound | `errs.NewValidationError(errs.SubtypeMissingConfig, ...)` |
| API returned non-2xx | handled in `client.Do` → `*errs.APIError` / `*errs.AuthError` |
| Network / transport failure | `errs.NewNetworkError(errs.SubtypeNetworkTransport, ...)` |
| Local file I/O failure | `errs.NewInternalError(errs.SubtypeFileIO, ...)` |
| Unclassified lower-layer error | `errs.NewInternalError(errs.SubtypeUnknown, ...).WithCause(err)` |
| Lower layer already returned a typed error | return it unchanged |

Authoring rules:

- Always preserve the underlying error with `.WithCause(err)` so `errors.Is` /
  `errors.As` keep working (every typed error's embedded `Problem` implements
  `Unwrap`).
- `.WithParam` names only the single user input that failed. Recovery guidance
  goes in `.WithHint(...)`.
- Never call `fmt.Sprintf` yourself for the message — pass the format + args to
  the constructor; it applies `formatMessage`, which is `%`-safe when no args
  are given.

### stdout is data, stderr is everything else

Program output (JSON / table rows) goes to `f.Out`. Progress, dry-run traces,
table summary lines, warnings, and error envelopes go to `f.ErrOut`. Never
write diagnostics to stdout — it corrupts `--jq` and pipe chains. Do not call
`os.Stdout` / `os.Stderr` directly in command or output code; thread the
`Factory` writers instead.

### Use `vfs.*` instead of `os.*` for filesystem access

All filesystem operations (ReadFile, WriteFile, MkdirAll, Stat, Remove, Rename,
Executable, EvalSymlinks, UserHomeDir) must go through `internal/vfs` via
`vfs.Default`. This enables test mocking without touching the real disk.

**Never** import `os` for file I/O in `internal/` or `cmd/` packages. The only
acceptable `os` usage is `os.Getenv` (environment variable reads) and
`os/exec` (subprocess execution). If a new filesystem operation is needed,
extend the `vfs.FS` interface and `vfs.OS` implementation first.

### Production-grade requirements

This CLI ships to real users and AI agents in production. Code must be:

- **Testable**: every behavior change ships with a test. Use `vfs.Default`
  reassignment or override fields (e.g. `SkillsCommandOverride`) to inject
  fakes — never make real network or subprocess calls in unit tests.
- **Cross-platform**: all platform-specific code uses `//go:build` constraints
  (`updater_unix.go`, `updater_windows.go`). Verify with `GOOS=windows go
  build ./...` before PR.
- **Timeout-bounded**: every network call and subprocess invocation must have
  a `context.WithTimeout`. Never block indefinitely.
- **Error-preserving**: use `.WithCause(err)` on typed errors so
  `errors.Is` / `errors.As` keep working through the error chain.
- **No panics in production paths**: recover from unexpected states with
  typed errors, not panics.

### Adding a command

- Register shared flags via `addTenantFlag` / `addPagingFlags` /
  `addOrderFlags` / `addDateRangeFlags` / `addJQFlag`.
- Build params with `setString` / `setInt` so zero values are omitted.
- Call `f.runList` (lists/reports) or `f.runSingle` (single object). Both honor
  `--dry-run`, `--jq`, `--format`, and `--page-all`.
- Add a column set to `cmd/columns.go` for `--format table`.

### Tests

- Every behavior change needs a test alongside the change.
- Command tests use the in-process harness in `cmd/testutil_test.go`
  (`runCmd(t, args...)`), driving commands with `--dry-run` to assert request
  path/params without a network.
- Client tests use `net/http/httptest`.
- Assert typed errors with `errors.As(err, &ve)` and check
  `category` / `subtype` / `param`, not message substrings.

## Commit & PR

- Conventional Commits in English: `feat:`, `fix:`, `docs:`, `test:`,
  `refactor:`, `chore:`, `ci:`.
- Never commit secrets, tokens, or real credentials.
