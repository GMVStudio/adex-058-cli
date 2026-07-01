# adex — ADEX CLI

ADEX CLI tool for querying advertising data. Built for humans and AI Agents.

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
git clone https://github.com/GMVStudio/adex-058-cli.git
cd adex-058-cli
make install
```

#### Configure & Use

```bash
# 1. Set API base URL (one-time, or use --base-url flag per command)
export ADEX_API_BASE_URL=http://your-api-host:8000

# 2. Start using
adex ks dashboard --tenant 6
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
adex --help
```

> `--help` prints available commands, confirming the CLI is installed and configured correctly.

## Agent Skills

The CLI embeds AI Agent skills at build time, serving them via `adex skills` commands.

| Skill | Description |
|-------|-------------|
| `adex-shared` | Setup, API credentials, shared flags reference (pagination, jq, date range, output format, error handling) |
| `adex-ks` | Kuaishou (快手) advertising data: accounts, campaigns, units, creatives, reports, top-N, metric meta, dashboard |
| `adex-oe` | Oceanengine (巨量) advertising data: accounts, projects, units, reports, top-N, metric meta, dashboard, budget vs actual |
| `adex-tenant-user` | Tenant listing with filters and current user info |

```bash
# List all skills
adex skills list

# Read a skill's SKILL.md
adex skills read adex-shared
adex skills read adex-ks
adex skills read adex-oe
adex skills read adex-tenant-user

# Read as JSON envelope
adex skills read adex-shared --json

# List files under a skill
adex skills list adex-ks
```

## Usage

```bash
# Show available commands
adex --help

# Query KS dashboard
adex ks dashboard --tenant 6 --range 30d

# Query OE dashboard
adex oe dashboard --tenant 6 --range 30d

# List tenants
adex tenant --status active --format table

# Get current user
adex user

# Budget vs actual comparison
adex oe account-budget-vs-actual --tenant 6 --range 30d --format table
```

## Commands

Run `adex --help` to see all available commands and subcommands.

### init

Bind API credentials (one-time setup).

```bash
adex init --authorization "Bearer adex_xxx"
```

### ks — Kuaishou (快手)

| Subcommand | Description |
|------------|-------------|
| `accounts` | List ad accounts |
| `campaigns` | List campaigns (supports `top` and `get` subcommands) |
| `units` | List ad units (supports `top` and `get`) |
| `creatives` | List creatives (supports `top` and `get`) |
| `account-reports daily/summary` | Account-level reports |
| `campaign-reports daily/summary` | Campaign-level reports |
| `unit-reports daily/summary` | Unit-level reports |
| `creative-reports daily/summary` | Creative-level reports |
| `report-metric-meta` | Report metric metadata |
| `dashboard` | Tenant-level overview |

### oe — Oceanengine (巨量)

| Subcommand | Description |
|------------|-------------|
| `accounts` | List ad accounts |
| `projects` | List projects (supports `top` and `get`) |
| `units` | List units (supports `top` and `get`) |
| `account-reports daily/summary` | Account-level reports |
| `project-reports daily/summary` | Project-level reports |
| `unit-reports daily/summary` | Unit-level reports |
| `report-metric-meta` | Report metric metadata |
| `dashboard` | Tenant-level overview |
| `account-budget-vs-actual` | Budget vs actual spend comparison |

### tenant

List tenants with optional name/status filters and pagination.

### user

Get current authenticated user info (resolved from Bearer API key).

### skills

Read embedded AI Agent skill content.

```bash
adex skills list                    # list all skills
adex skills read adex-ks            # read a skill's SKILL.md
adex skills list adex-ks            # list files under a skill
```

### Global Flags

| Flag | Description |
|------|-------------|
| `--format` | Output format: `json` (default), `pretty`, `table` |
| `--dry-run` | Print request without executing |
| `--base-url` | API base URL (overrides `ADEX_API_BASE_URL` env) |

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `ADEX_API_BASE_URL` | `http://47.99.131.55:8000` | API base URL |
| `ADEX_AUTHORIZATION` | — | API key (Bearer prefix auto-stripped) |
| `ADEX_CONFIG_DIR` | `~/.adex` | Config directory (for testing) |

## Architecture

```
adex-058-cli/
├── main.go                         # 入口
├── skills_embed.go                 # Go embed for skills/*/SKILL.md + references/
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
│   ├── adex-shared/SKILL.md        # shared setup, config & flags reference
│   ├── adex-ks/SKILL.md            # Kuaishou advertising data skill
│   ├── adex-oe/SKILL.md            # Oceanengine advertising data skill
│   └── adex-tenant-user/SKILL.md   # tenant & user management skill
├── cmd/
│   ├── root.go                     # 根命令、全局 flags、类型化错误处理
│   ├── init.go                     # "init" 命令 — 绑定 API Key
│   ├── skill.go                    # "skills" 命令组 (list/read)
│   ├── ks.go                       # "ks" 命令组 (快手)
│   ├── ks_accounts.go              # ks accounts
│   ├── ks_campaigns.go             # ks campaigns (list/top/get)
│   ├── ks_units.go                 # ks units (list/top/get)
│   ├── ks_creatives.go             # ks creatives (list/top/get)
│   ├── ks_reports.go               # ks *-reports (daily/summary)
│   ├── ks_meta.go                  # ks report-metric-meta
│   ├── ks_dashboard.go             # ks dashboard
│   ├── oe.go                       # "oe" 命令组 (巨量)
│   ├── oe_accounts.go              # oe accounts
│   ├── oe_projects.go              # oe projects (list/top/get)
│   ├── oe_units.go                 # oe units (list/top/get)
│   ├── oe_reports.go               # oe *-reports (daily/summary)
│   ├── oe_meta.go                  # oe report-metric-meta
│   ├── oe_dashboard.go             # oe dashboard
│   ├── oe_budget.go                # oe account-budget-vs-actual
│   ├── tenant.go                   # tenant list
│   ├── user.go                     # user (current user info)
│   ├── common.go                   # 共享 flag 注册与执行 helpers
│   ├── report_shared.go            # 共享 top/get/daily/summary 命令构建器
│   └── columns.go                  # 表格列定义
├── errs/
│   └── errs.go                     # 类型化错误分类体系 (RFC 7807 对齐)
└── internal/
    ├── build/build.go              # 版本信息
    ├── client/client.go            # HTTP 客户端 + 类型化错误分类
    ├── config/config.go            # 环境变量 + 文件配置
    ├── output/output.go            # JSON / pretty / table 输出
    ├── output/jq.go                # jq 过滤支持
    ├── daterange/daterange.go      # 日期范围解析 (相对/显式)
    ├── paginate/paginate.go        # 分页聚合 (page-all)
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
