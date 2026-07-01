---
name: adex-tenant-user
version: 0.3.0
description: "租户管理与当前用户信息查询。当用户需要列出租户、按名称/状态过滤租户，或查询当前登录用户信息时使用。不负责广告数据查询（走 adex-ks / adex-oe）。"
metadata:
  requires:
    bins: ["adex"]
  cliHelp: "adex tenant --help; adex user --help"
---

# tenant & user — 租户与用户

开始前先读 [`../adex-shared/SKILL.md`](../adex-shared/SKILL.md)（安装、认证、共享 flags）。

## 命令总览

| 命令 | 说明 | API |
|------|------|-----|
| `tenant` | 租户列表（支持名称模糊过滤、状态过滤、分页） | `GET /v1/tenants` |
| `user` | 当前登录用户信息（由 Bearer API Key 解析） | `GET /v1/users/me` |

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
