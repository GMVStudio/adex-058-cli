# tenant — 租户列表 / 切换默认租户

> **前置条件：** 先阅读 [`../SKILL.md`](../SKILL.md) 了解安装、认证、共享 flags。

查询 adex 平台上的租户列表，支持按名称模糊匹配和状态精确过滤。对应 API：`GET /v1/tenants`。

`adex tenant use <ID>` 可设定默认租户，后续所有命令无需再传 `--tenant`。

> **注意：** 此命令**不需要 `--tenant`** flag。

## 命令

```bash
# 列出所有租户
adex tenant --page-size 20

# 按名称模糊过滤
adex tenant --name acme --format table

# 按状态过滤
adex tenant --status active --page-size 50

# 聚合所有页
adex tenant --page-all --jq '.items[].id'

# 预览请求但不执行
adex tenant --dry-run
```

## Flags

| Flag | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--name` | string | — | 租户名称模糊匹配（留空=不过滤） |
| `--status` | string | — | 状态精确匹配：`active` / `disabled`（留空=不过滤） |
| `--page-size` | int | 20 | 每页条数 |
| `--page-token` | string | — | 游标分页 token |
| `--page-all` | bool | false | 聚合所有页 |
| `--jq` | string | — | jq 表达式过滤输出 |
| `--format` | enum | `json` | `json` / `pretty` / `table` |
| `--dry-run` | bool | false | 打印请求但不执行 |

> 此命令**没有** `--tenant`、`--order-by`、`--order-desc` flags。

## 输出

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

## 使用场景

- **切换默认租户**：`adex tenant use 6` 设定默认租户，后续命令无需 `--tenant`
- **获取租户 ID**：`--format table` 查看所有租户及其 ID
- **搜索特定租户**：`--name "关键词"` 按名称模糊匹配
- **筛选活跃租户**：`--status active` 只看状态为 active 的租户
- **提取所有租户 ID**：`--page-all --jq '.items[].id'` 获取完整 ID 列表
- **验证权限范围**：列出当前 API Key 可访问的所有租户

## 注意事项

- `--name` 为服务端模糊匹配，不是客户端过滤
- `--status` 为精确匹配，只支持 `active` 和 `disabled` 两个值
- 此命令不需要 `--tenant`，是少数不需要租户 ID 的命令之一

## tenant use — 设定默认租户

```bash
# 设定默认租户
adex tenant use 6

# 切换到其他租户
adex tenant use 8
```

设定后，所有支持 `--tenant` 的命令会自动使用该租户，无需每次传入 `--tenant`。

### 优先级

`--tenant` 标志 > `ADEX_TENANT_ID` 环境变量 > 配置文件中的 `tenant_id` > 报错

### 输出

```json
{
  "ok": true,
  "message": "default tenant set",
  "tenant_id": 6,
  "path": "/home/user/.adex/config.json"
}
```

## 参考

- [adex-shared](../SKILL.md) — 共享规则和 Skill 路由
- [adex-shared-user](adex-shared-user.md) — 当前用户信息（也可获取租户 ID）
