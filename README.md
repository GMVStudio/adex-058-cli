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
| `adex-shared` | Setup, API endpoint configuration, output formats, error handling reference |

```bash
# List all skills
adex skills list

# Read a skill's SKILL.md
adex skills read adex-shared

# Read as JSON envelope
adex skills read adex-shared --json
```

## Usage

```bash
# Show available commands
adex --help

# Query KS dashboard
adex ks dashboard --tenant 6

# Query OE dashboard
adex oe dashboard --tenant 6
```

## Commands

Run `adex --help` to see all available commands and subcommands.

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
│   └── adex-shared/SKILL.md        # shared setup & config skill
├── cmd/
│   ├── root.go                     # 根命令、全局 flags、类型化错误处理
│   ├── skill.go                    # "skills" 命令组 (list/read)
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
