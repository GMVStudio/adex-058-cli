---
name: adex-shared
version: 0.3.0
description: "Use when first setting up adex CLI, configuring API credentials, or needing shared flags reference (pagination, jq, date range, output format, error handling). Also covers tenant listing with filters and current user info query."
metadata:
  requires:
    bins: ["adex"]
  cliHelp: "adex --help"
---

# adex CLI 共享规则

本技能指导你如何通过 adex CLI 查询广告投放数据。开始前必读。

## 安装

### 通过 npm 安装（推荐）

```bash
# 安装 CLI
npm install -g @gmvstudio/adex-cli

# 安装 CLI SKILL（必需）
npx -y skills add https://adex-skills.oss-cn-hangzhou.aliyuncs.com -g -y
```

### 从源码安装

```bash
git clone https://github.com/gmvstudio/adex-cli.git
cd adex-cli
make install
```

## 初始化配置

首次使用前，必须通过 `adex init` 绑定 API Key：

```bash
adex init --authorization "Bearer adex_c93462599a6246a89f55a11b024b1a1a"
```

也可以传入裸 key（自动补 Bearer 前缀）：

```bash
adex init --authorization adex_c93462599a6246a89f55a11b024b1a1a --base-url http://47.99.131.55:8000
```

配置写入 `~/.adex/config.json`（0600 权限）。

### 环境变量覆盖

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `ADEX_API_BASE_URL` | `http://47.99.131.55:8000` | API base URL |
| `ADEX_AUTHORIZATION` | — | API key（自动去 Bearer 前缀） |
| `ADEX_CONFIG_DIR` | `~/.adex` | 配置目录（测试用） |

环境变量优先于配置文件；`--base-url` 标志优先于环境变量。

## 验证

```bash
adex --help
adex user                    # 验证 API Key 是否有效
```

## 命令树总览

```
adex
├── init                      # 绑定 API Key（一次性）
├── ks                        # 快手广告数据
│   ├── accounts              # 广告账户列表
│   ├── campaigns             # 广告计划列表 / top / get
│   ├── units                 # 广告组列表 / top / get
│   ├── creatives             # 创意列表 / top / get
│   ├── account-reports       # 账户报表 daily / summary
│   ├── campaign-reports      # 计划报表 daily / summary
│   ├── unit-reports          # 组报表 daily / summary
│   ├── creative-reports      # 创意报表 daily / summary
│   ├── report-metric-meta    # 报表指标元数据
│   └── dashboard             # 租户级概览
├── oe                        # 巨量引擎广告数据
│   ├── accounts              # 广告账户列表
│   ├── projects              # 项目列表 / top / get
│   ├── units                 # 单元列表 / top / get
│   ├── account-reports       # 账户报表 daily / summary
│   ├── project-reports       # 项目报表 daily / summary
│   ├── unit-reports          # 单元报表 daily / summary
│   ├── report-metric-meta    # 报表指标元数据
│   ├── dashboard             # 租户级概览
│   └── account-budget-vs-actual # 预算 vs 实际消耗
├── tenant                    # 租户列表
├── user                      # 当前用户信息
└── skills                    # 嵌入式 Skill 内容
    ├── list                  # 列出所有 Skill
    └── read                  # 读取 Skill 内容
```

## 共享 Flags

以下 flags 在大多数命令中通用：

| Flag | 说明 | 默认值 |
|------|------|--------|
| `--tenant` | 租户 ID（大多数命令必需） | — |
| `--page-size` | 每页条数 | 20 |
| `--page-token` | 指定页的游标 token | — |
| `--page-all` | 聚合所有页 | false |
| `--order-by` | 排序字段 | 因命令而异 |
| `--order-desc` | 降序排序 | true |
| `--range` | 相对日期范围（如 7d/4w/1m） | — |
| `--begin` | 起始日期（YYYY-MM-DD） | — |
| `--end` | 结束日期（YYYY-MM-DD） | — |
| `--jq` | jq 表达式过滤 JSON 输出 | — |
| `--format` | 输出格式：json / pretty / table | json |
| `--dry-run` | 打印请求但不执行 | false |
| `--base-url` | 覆盖 API base URL | — |

## 分页

列表接口统一使用 `page_token` 游标分页：

```bash
# 单页查询
adex ks accounts --tenant 6 --page-size 20

# 翻页（透传上一次响应的 next_page_token）
adex ks accounts --tenant 6 --page-token "abc123"

# 聚合所有页（自动翻页直到 has_more=false）
adex ks accounts --tenant 6 --page-all
```

## 日期范围

报表和 dashboard 命令支持灵活的日期范围指定：

```bash
# 相对范围（7d=7天, 4w=4周, 1m=1月）
adex ks dashboard --tenant 6 --range 30d

# 显式日期
adex ks dashboard --tenant 6 --begin 2026-06-01 --end 2026-06-30

# --range 优先于 --begin/--end
```

## jq 过滤

所有命令支持 `--jq` 对 JSON 输出进行过滤：

```bash
# 提取所有 advertiserId
adex ks accounts --tenant 6 --page-all --jq '.items[].advertiserId'

# 提取单个字段
adex user --jq '.username'
```

## 输出格式

通过 `--format` 标志控制：

- `json`（默认）：紧凑 JSON
- `pretty`：格式化 JSON
- `table`：表格输出

```bash
adex ks accounts --tenant 6 --format table
adex ks dashboard --tenant 6 --range 30d --format pretty
```

## Dry-Run

`--dry-run` 打印请求路径和参数到 stderr，不实际调用 API：

```bash
adex ks accounts --tenant 6 --dry-run
```

## 错误处理

错误以 JSON 信封格式输出到 stderr，包含 `type`、`subtype`、`message` 字段：

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

## tenant — 租户列表

不需要 `--tenant` flag。支持名称模糊匹配和状态精确过滤。

```bash
# 列出所有租户
adex tenant --page-size 20

# 按名称模糊过滤
adex tenant --name acme --format table

# 按状态过滤
adex tenant --status active --page-size 50

# 聚合所有页
adex tenant --page-all --jq '.items[].id'
```

| Flag | 说明 |
|------|------|
| `--name` | 租户名称模糊匹配（留空=不过滤） |
| `--status` | 状态精确匹配：active / disabled（留空=不过滤） |
| `--page-size` | 每页条数（默认 20，最大 200） |
| `--page-token` | 游标分页 token |
| `--page-all` | 聚合所有页 |
| `--jq` | jq 表达式过滤输出 |

### 响应结构

```json
{
  "hasMore": true,
  "nextPageToken": "abc123",
  "items": [
    {
      "id": 1,
      "name": "Acme Corp",
      "status": "active",
      "createdBy": 100,
      "createdAt": "2026-01-01T00:00:00Z",
      "updatedAt": "2026-06-01T00:00:00Z"
    }
  ]
}
```

### Table 列

| 列 | 字段 |
|----|------|
| ID | `id` |
| Name | `name` |
| Status | `status` |
| Created By | `createdBy` |
| Created At | `createdAt` |
| Updated At | `updatedAt` |

## user — 当前用户信息

不需要 `--tenant` flag。通过 Bearer API Key 自动解析当前用户。

```bash
# JSON 输出（默认）
adex user

# 表格输出
adex user --format table

# 提取单个字段
adex user --jq '.username'
adex user --jq '.currentTenantId'
```

### 响应结构

```json
{
  "id": 100,
  "username": "admin",
  "name": "管理员",
  "status": "active",
  "currentTenantId": 6,
  "createdAt": "2026-01-01T00:00:00Z",
  "updatedAt": "2026-06-01T00:00:00Z"
}
```

### Table 列

| 列 | 字段 |
|----|------|
| ID | `id` |
| Username | `username` |
| Name | `name` |
| Status | `status` |
| Current Tenant | `currentTenantId` |
| Created At | `createdAt` |
| Updated At | `updatedAt` |

## 常见用法

```bash
# 验证 API Key 是否有效
adex user

# 查看当前租户 ID（用于其他命令的 --tenant 参数）
adex user --jq '.currentTenantId'

# 列出所有活跃租户
adex tenant --status active --page-all --format table

# 查找特定租户
adex tenant --name "Acme" --format table
```

## Skill 路由

| 用户意图 | 路由到 Skill |
|----------|-------------|
| 快手广告数据查询 | [`adex-ks`](../adex-ks/SKILL.md) |
| 巨量引擎广告数据查询 | [`adex-oe`](../adex-oe/SKILL.md) |
| 安装、配置、共享 flags、租户管理、用户信息 | 本 Skill |
