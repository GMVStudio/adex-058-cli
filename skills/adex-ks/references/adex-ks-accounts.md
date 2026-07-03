# ks accounts — 广告账户列表

> **前置条件：** 先阅读 [`../adex-shared/SKILL.md`](../../adex-shared/SKILL.md) 了解安装、认证、共享 flags。

查询快手广告账户列表，支持按名称、类型、状态等多维度筛选。对应 API：`GET /v1/ks/ad-accounts`。

## 命令

```bash
# 基本列表
adex ks accounts --page-size 20

# 按余额降序排列，表格输出
adex ks accounts --order-by balance --order-desc --format table

# 聚合所有页，提取广告主 ID
adex ks accounts --page-all --jq '.items[].advertiserId'

# 按账户名模糊匹配
adex ks accounts --account-name "品牌" --format table

# 按授权状态和投放状态筛选
adex ks accounts --auth-status active --delivery-status active --format table

# 预览请求但不执行
adex ks accounts --dry-run
```

## Flags

### 专属 Flags

| Flag | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--advertiser` | string | — | 广告主 ID 精确过滤 |
| `--account-name` | string | — | 账户名模糊匹配 |
| `--account-type` | string | — | 账户类型过滤 |
| `--auth-status` | string | — | 授权状态过滤 |
| `--delivery-status` | string | — | 投放状态过滤 |
| `--active-status` | string | — | 活跃状态过滤 |
| `--owner-user` | int | 0 | 归属用户 ID 过滤（0=不过滤） |

### 共享 Flags

| Flag | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--tenant` | int | — | 租户 ID（可选；缺省使用 `adex tenant use` 设定的默认租户） |
| `--page-size` | int | 20 | 每页条数 |
| `--page-token` | string | — | 游标分页 token |
| `--page-all` | bool | false | 聚合所有页 |
| `--order-by` | string | `id` | 排序字段 |
| `--order-desc` | bool | true | 降序排序 |
| `--jq` | string | — | jq 表达式过滤输出 |
| `--format` | enum | `json` | `json` / `pretty` / `table` |
| `--dry-run` | bool | false | 打印请求但不执行 |

## 输出

```json
{
  "hasMore": true,
  "nextPageToken": "abc123",
  "items": [
    {
      "id": 1,
      "advertiserId": "1234567890",
      "accountName": "品牌推广账户",
      "accountType": "NORMAL",
      "authStatus": "active",
      "deliveryStatus": "active",
      "activeStatus": "active",
      "balance": 50000.00
    }
  ]
}
```

### Table 列

| 列 | 字段 |
|----|------|
| ID | `id` |
| Advertiser ID | `advertiserId` |
| Account Name | `accountName` |
| Account Type | `accountType` |
| Auth Status | `authStatus` |
| Delivery Status | `deliveryStatus` |
| Active Status | `activeStatus` |
| Balance | `balance` |

## 使用场景

- **盘点账户**：`--page-all --format table` 查看全部账户
- **查找特定账户**：`--account-name "关键词"` 模糊匹配
- **筛选活跃账户**：`--delivery-status active --auth-status active`
- **获取广告主 ID**：`--jq '.items[].advertiserId'` 提取 ID 供后续命令使用
- **按余额排序**：`--order-by balance --order-desc` 找出余额最高/最低的账户

## 注意事项

- `advertiserId` 是后续命令（`--advertiser`）的关键参数，可先用本命令获取
- 账户名模糊匹配为服务端匹配，不是客户端过滤
- `--page-all` 会自动翻页直到 `hasMore=false`，大量账户时注意耗时

## 参考

- [adex-ks](../SKILL.md) — 快手广告全部命令
- [adex-shared](../../adex-shared/SKILL.md) — 认证和全局参数
