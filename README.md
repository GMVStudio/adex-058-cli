# adex — ADEX CLI

ADEX CLI tool for querying campaign reports and advertising data. Built for humans and AI Agents.

## Installation & Quick Start

### Requirements

Before you start, make sure you have:

- Node.js (`npm`/`npx`) — for npm installation
- Go `v1.23`+ — only required for building from source

### Quick Start (Human Users)

#### Install

Choose **one** of the following methods:

**Option 1 — From npm (recommended):**

```bash
npx @gmvstudio/adex-cli@latest install
```

**Option 2 — From source:**

```bash
git clone https://github.com/gmvstudio/adex-cli.git
cd adex-cli
make install
```

#### Configure & Use

```bash
# 1. Set API base URL (one-time, or use --base-url flag per command)
export ADEX_API_BASE_URL=http://your-api-host:8000

# 2. Start using
adex raw campaign daily --tenant 6 --range 1d
```

### Quick Start (AI Agent)

> The following steps are for AI Agents. Some steps may require the user to provide configuration values.

**Step 1 — Install**

```bash
# 安装 CLI
npm install -g @gmvstudio/adex-cli

# 安装 CLI SKILL（必需）
npx -y skills add https://open.feishu.cn --skill -y
```

**Step 2 — Configure API endpoint**

> Set the environment variable to point to the ADEX API server. The user should provide the correct API host.

```bash
export ADEX_API_BASE_URL=http://your-api-host:8000
```

**Step 3 — Verify**

```bash
adex raw campaign daily --tenant 6 --range 1d --dry-run
```

> `--dry-run` prints the request without executing, confirming the CLI is installed and configured correctly.

## Agent Skills

The CLI embeds AI Agent skills at build time, serving them via `adex skills` commands.

| Skill | Description |
|-------|-------------|
| `adex-shared` | Setup, API endpoint configuration, output formats, error handling reference |
| `adex-campaign` | Query campaign daily reports with filtering, sorting, and pagination |

```bash
# List all skills
adex skills list

# Read a skill's SKILL.md
adex skills read adex-campaign

# Read as JSON envelope
adex skills read adex-campaign --json
```

## Usage

```bash
# Query campaign daily report
adex raw campaign daily --tenant 6 --campaign C-618-001-619 --range 1d

# With pretty JSON output
adex raw campaign daily --tenant 6 --range 1d --format pretty

# With table output
adex raw campaign daily --tenant 6 --range 1d --format table

# Dry run (print request without executing)
adex raw campaign daily --tenant 6 --range 1d --dry-run

# Custom API base URL
adex raw campaign daily --tenant 6 --range 1d --base-url http://api.example.com
```

## Commands

### `adex raw campaign daily`

Query daily campaign reports.

| Flag | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `--tenant` | int | yes | — | Tenant ID |
| `--range` | string | yes | — | Time range: `1d`, `7d`, `1h`, `30m` |
| `--campaign` | string | no | — | Campaign ID (numeric) or name pattern |
| `--page` | int | no | 1 | Page number |
| `--page-size` | int | no | 20 | Page size |
| `--order-by` | string | no | charge | Sort field |
| `--order-desc` | bool | no | true | Sort descending |
| `--stat-hour` | int | no | -1 | Stat hour (-1 for latest) |

### Global Flags

| Flag | Description |
|------|-------------|
| `--format` | Output format: `json` (default), `pretty`, `table` |
| `--dry-run` | Print request without executing |
| `--base-url` | API base URL (overrides `ADEX_API_BASE_URL` env) |

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `ADEX_API_BASE_URL` | `http://localhost:8000` | API base URL |

## Architecture

```
adex-058-cli/
├── main.go                         # 入口
├── skills_embed.go                 # Go embed for skills/*/SKILL.md
├── Makefile                        # build/vet/fmt-check/test/install
├── package.json                    # npm distribution
├── .goreleaser.yml                 # cross-platform release builds
├── README.md                       # 使用文档
├── go.mod / go.sum
├── scripts/
│   ├── install.js                  # postinstall: download pre-built binary
│   ├── install-wizard.js           # interactive setup wizard (npx adex install)
│   └── run.js                      # npm bin wrapper → delegates to binary
├── skills/
│   ├── adex-shared/SKILL.md        # shared setup & config skill
│   └── adex-campaign/SKILL.md      # campaign report query skill
├── cmd/
│   ├── root.go                     # 根命令、全局 flags、类型化错误处理
│   ├── skill.go                    # "skills" 命令组 (list/read)
│   ├── raw.go                      # "raw" 命令组
│   ├── raw_campaign.go             # "raw campaign" 子命令组
│   └── raw_campaign_daily.go       # "raw campaign daily" 命令实现
├── errs/
│   └── errs.go                     # 类型化错误分类体系 (RFC 7807 对齐)
└── internal/
    ├── build/build.go              # 版本信息
    ├── client/client.go            # HTTP 客户端 + 类型化错误分类
    ├── config/config.go            # 环境变量配置
    ├── output/output.go            # JSON / pretty / table 输出
    └── skillcontent/reader.go      # 嵌入式 skill 内容读取器
```

## Error Handling

Errors are typed and output as JSON envelopes on stderr:

```json
{"ok":false,"error":{"type":"validation","subtype":"invalid_argument","message":"--tenant must be a positive integer"}}
```

| Exit Code | Category | Description |
|-----------|----------|-------------|
| 0 | — | Success |
| 2 | validation | Invalid input arguments |
| 3 | unauthorized | Authentication failure |
| 4 | network | Network/transport error |
| 5 | api | API returned non-2xx |
| 1 | internal | Unexpected internal error |
