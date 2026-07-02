# oe dashboard — 租户级概览

> **前置条件：** 先阅读 [`../adex-shared/SKILL.md`](../../adex-shared/SKILL.md) 了解安装、认证、共享 flags。

租户级别的巨量引擎广告投放概览，一次调用返回账户统计、日期范围内消耗汇总、以及账户消耗排名。对应 API：`GET /v1/oe/dashboard`。

适合作为分析的第一步，快速了解整体投放情况后再下钻到具体资源。

## 命令

```bash
# 最近 30 天概览
adex oe dashboard --tenant 6 --range 30d

# 显式日期范围
adex oe dashboard --tenant 6 --begin 2026-06-01 --end 2026-06-30

# 最近 7 天概览
adex oe dashboard --tenant 6 --range 7d

# pretty 格式输出
adex oe dashboard --tenant 6 --range 30d --format pretty

# 提取消耗汇总
adex oe dashboard --tenant 6 --range 30d --jq '.charge'

# 预览请求但不执行
adex oe dashboard --tenant 6 --range 30d --dry-run
```

## Flags

| Flag | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--tenant` | int | — | 租户 ID（**必需**） |
| `--range` | string | — | 相对日期范围如 `7d` / `4w` / `1m`（**必需**） |
| `--begin` | string | — | 起始日期（YYYY-MM-DD，与 `--end` 配合） |
| `--end` | string | — | 结束日期（YYYY-MM-DD） |
| `--jq` | string | — | jq 表达式过滤输出 |
| `--format` | enum | `json` | `json` / `pretty` / `table` |
| `--dry-run` | bool | false | 打印请求但不执行 |

> `--range` 优先于 `--begin` / `--end`。日期范围至少需要一种指定方式。

## 输出

返回单个 JSON 对象（非列表），包含租户级汇总信息：

```json
{
  "accountCount": 15,
  "activeAccountCount": 10,
  "charge": 500000.00,
  "range": {
    "begin": "2026-06-01",
    "end": "2026-06-30"
  },
  "rankings": [
    {
      "advertiserId": "1866874042754522",
      "accountName": "品牌推广账户",
      "charge": 50000.00
    }
  ]
}
```

> 实际返回字段以 API 响应为准。以上为典型结构示例。

## 使用场景

- **快速了解大盘**：`--range 30d` 一眼看全局
- **对比不同时间段**：分别用 `--range 7d` 和 `--range 30d` 调用，对比消耗变化
- **找到消耗最多的账户**：查看 `rankings` 部分，或用 `--jq '.rankings[:5]'` 取前 5
- **作为下钻起点**：先 dashboard 看大盘 → 再 `account-reports daily` 看趋势 → 再 `projects top` 看排名

## 注意事项

- `dashboard` 返回的是单个对象，不是列表，因此**不支持** `--page-size` / `--page-token` / `--page-all`
- 日期范围为**必需**参数，不传会返回验证错误

## 参考

- [adex-oe](../SKILL.md) — 巨量引擎广告全部命令
- [adex-oe-accounts](adex-oe-accounts.md) — 账户列表
- [adex-oe-reports](adex-oe-reports.md) — 日报表和汇总报表
- [adex-oe-workflow-top-analysis](adex-oe-workflow-top-analysis.md) — Top-N 分析工作流
- [adex-shared](../../adex-shared/SKILL.md) — 认证和全局参数
