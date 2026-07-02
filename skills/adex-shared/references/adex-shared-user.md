# user — 当前用户信息

> **前置条件：** 先阅读 [`../SKILL.md`](../SKILL.md) 了解安装、认证、共享 flags。

查询当前认证用户的个人信息，通过 Bearer API Key 自动解析。对应 API：`GET /v1/users/me`。

> **注意：** 此命令**不需要 `--tenant`** flag。

## 命令

```bash
# JSON 输出（默认）
adex user

# 表格输出
adex user --format table

# 提取当前租户 ID（用于其他命令的 --tenant 参数）
adex user --jq '.currentTenantId'

# 提取用户名
adex user --jq '.username'

# pretty 格式输出
adex user --format pretty

# 预览请求但不执行
adex user --dry-run
```

## Flags

| Flag | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--jq` | string | — | jq 表达式过滤输出 |
| `--format` | enum | `json` | `json` / `pretty` / `table` |
| `--dry-run` | bool | false | 打印请求但不执行 |

> 此命令**没有** `--tenant`、`--page-size`、`--page-token`、`--page-all`、`--order-by`、`--order-desc` flags。它返回单个对象，不是列表。

## 输出

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

## 使用场景

- **验证 API Key**：`adex user` 能正常返回说明 API Key 有效
- **获取租户 ID**：`--jq '.currentTenantId'` 提取当前租户 ID，用于 ks/oe 命令的 `--tenant` 参数
- **查看用户信息**：`--format table` 或 `--format pretty` 查看完整用户信息
- **检查用户状态**：`--jq '.status'` 确认账户是否活跃

## 注意事项

- 此命令返回的是单个对象，不是列表，因此**不支持** `--page-size` / `--page-token` / `--page-all`
- `currentTenantId` 是最常用的字段，作为其他命令 `--tenant` 参数的来源
- 如果 API Key 无效或过期，会返回 `unauthorized` 错误（exit code 3）

## 参考

- [adex-shared](../SKILL.md) — 共享规则和 Skill 路由
- [adex-shared-tenant](adex-shared-tenant.md) — 租户列表（也可获取租户 ID）
