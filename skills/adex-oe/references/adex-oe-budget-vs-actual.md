# oe account-budget-vs-actual — 预算 vs 实际消耗

> **前置条件：** 先阅读 [`../adex-shared/SKILL.md`](../../adex-shared/SKILL.md) 了解安装、认证、共享 flags。

对比每个巨量引擎广告账户的日预算与实际日均消耗，帮助用户了解预算使用效率。对应 API：`GET /v1/oe/account-budget-vs-actual`。

> **注意：** 此命令是巨量引擎特有的，快手没有对应命令。

## 命令

```bash
# 全部账户的预算使用情况
adex oe account-budget-vs-actual --tenant 6 --range 30d

# 表格输出，直观对比
adex oe account-budget-vs-actual --tenant 6 --range 30d --format table

# 显式日期范围
adex oe account-budget-vs-actual --tenant 6 --begin 2026-06-01 --end 2026-06-30 --format table

# 单个账户的预算使用情况
adex oe account-budget-vs-actual --tenant 6 --advertiser 1866874042754522 --range 30d

# 提取预算使用率
adex oe account-budget-vs-actual --tenant 6 --range 30d --jq '.items[].budgetUsageRate'

# 预览请求但不执行
adex oe account-budget-vs-actual --tenant 6 --range 30d --dry-run
```

## Flags

| Flag | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--tenant` | int | — | 租户 ID（**必需**） |
| `--advertiser` | string | — | 广告主 ID 过滤（单个账户） |
| `--range` | string | — | 相对日期范围如 `7d` / `4w` / `1m`（**必需**） |
| `--begin` | string | — | 起始日期（YYYY-MM-DD，与 `--end` 配合） |
| `--end` | string | — | 结束日期（YYYY-MM-DD） |
| `--jq` | string | — | jq 表达式过滤输出 |
| `--format` | enum | `json` | `json` / `pretty` / `table` |
| `--dry-run` | bool | false | 打印请求但不执行 |

> `--range` 优先于 `--begin` / `--end`。日期范围至少需要一种指定方式。

## 输出

```json
{
  "hasMore": false,
  "nextPageToken": "",
  "items": [
    {
      "advertiserId": "1866874042754522",
      "accountName": "品牌推广账户",
      "budgetMode": "BUDGET_MODE_DAY",
      "budget": 10000.00,
      "totalCharge": 8500.00,
      "days": 30,
      "avgDailyCharge": 283.33,
      "budgetUsageRate": 0.028,
      "balance": 50000.00
    }
  ]
}
```

### Table 列

| 列 | 字段 |
|----|------|
| Advertiser ID | `advertiserId` |
| Account Name | `accountName` |
| Budget Mode | `budgetMode` |
| Budget | `budget` |
| Total Charge | `totalCharge` |
| Days | `days` |
| Avg Daily Charge | `avgDailyCharge` |
| Budget Usage Rate | `budgetUsageRate` |
| Balance | `balance` |

## 使用场景

- **预算执行概览**：`--range 30d --format table` 查看所有账户的预算使用情况
- **单个账户分析**：`--advertiser <ID> --range 30d` 深入分析特定账户
- **找出预算消耗过快的账户**：`--jq '.items | sort_by(.budgetUsageRate) | reverse'` 按预算使用率排序
- **对比不同时间段**：分别用 `--range 7d` 和 `--range 30d` 调用，对比预算使用趋势
- **结合 dashboard 使用**：先 `dashboard` 看大盘消耗，再用本命令看预算执行情况

## 注意事项

- 此命令返回的是列表（非分页），不支持 `--page-size` / `--page-token` / `--page-all`
- `budgetUsageRate` 为日均消耗占日预算的比例，越接近 1 表示预算几乎用完
- `budgetMode` 表示预算模式（如日预算 `BUDGET_MODE_DAY`）
- 日期范围为**必需**参数，不传会返回验证错误

## 参考

- [adex-oe](../SKILL.md) — 巨量引擎广告全部命令
- [adex-oe-dashboard](adex-oe-dashboard.md) — 租户级概览
- [adex-oe-accounts](adex-oe-accounts.md) — 账户列表
- [adex-shared](../../adex-shared/SKILL.md) — 认证和全局参数
