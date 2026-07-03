# ks report-metric-meta — 报表指标元数据

> **前置条件：** 先阅读 [`../adex-shared/SKILL.md`](../../adex-shared/SKILL.md) 了解安装、认证、共享 flags。

查询快手报表可用指标字段的元数据，包括字段名、标签、分组、聚合方式、是否可排序等。对应 API：`GET /v1/ks/report-metric-meta`。

> **注意：** 此命令**不需要 `--tenant`**，是唯一不需要租户 ID 的 ks 命令。

## 命令

```bash
# 查询 account 层级所有指标
adex ks report-metric-meta --level account --page-size 50

# 查询 campaign 层级已启用的指标
adex ks report-metric-meta --level campaign --enabled 1 --page-size 50

# 按指标组名筛选
adex ks report-metric-meta --level campaign --group-name "消耗" --page-size 50

# 按字段名模糊匹配
adex ks report-metric-meta --level unit --field "charge" --page-size 50

# 只看可排序字段
adex ks report-metric-meta --level creative --sortable 1 --page-size 50

# 聚合所有页
adex ks report-metric-meta --level campaign --enabled 1 --page-all --jq '.items[].field'

# 表格输出
adex ks report-metric-meta --level account --enabled 1 --format table
```

## Flags

| Flag | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--level` | string | — | 维度：`account` / `campaign` / `unit` / `creative` |
| `--group-name` | string | — | 指标组名过滤 |
| `--field` | string | — | 字段名模糊匹配 |
| `--enabled` | int | 0 | 0=全部 / 1=启用 / 2=禁用 |
| `--sortable` | int | 0 | 0=全部 / 1=可排序 / 2=不可排序 |
| `--page-size` | int | 20 | 每页条数 |
| `--page-token` | string | — | 游标分页 token |
| `--page-all` | bool | false | 聚合所有页 |
| `--order-by` | string | `sort_order` | 排序字段 |
| `--order-desc` | bool | true | 降序排序 |
| `--jq` | string | — | jq 表达式过滤输出 |
| `--format` | enum | `json` | `json` / `pretty` / `table` |
| `--dry-run` | bool | false | 打印请求但不执行 |

> 此命令**没有** `--tenant` flag。

## 输出

```json
{
  "hasMore": true,
  "nextPageToken": "abc123",
  "items": [
    {
      "id": 1,
      "level": "campaign",
      "field": "charge",
      "label": "总消耗",
      "groupName": "消耗指标",
      "agg": "sum",
      "valueType": "number",
      "sortOrder": 1,
      "enabled": 1,
      "sortable": 1
    }
  ]
}
```

### Table 列

| 列 | 字段 |
|----|------|
| ID | `id` |
| Level | `level` |
| Field | `field` |
| Label | `label` |
| Group Name | `groupName` |
| Agg | `agg` |
| Value Type | `valueType` |
| Sort Order | `sortOrder` |
| Enabled | `enabled` |
| Sortable | `sortable` |

## 使用场景

- **查看可用指标**：`--level campaign --enabled 1` 查看计划报表有哪些可用字段
- **找排序字段**：`--level account --sortable 1` 查看哪些字段可用于 `--order-by`
- **找排名指标**：`--level campaign --enabled 1 --jq '.items[].field'` 获取字段名列表，用于 `top --metric <field>`
- **按组浏览**：`--level unit --group-name "消耗"` 按指标组分类查看

## 与其他命令的配合

1. 先查指标元数据，获取可用字段名：
   ```bash
   adex ks report-metric-meta --level campaign --enabled 1 --page-all --jq '.items[].field'
   ```

2. 用字段名作为报表的 `--order-by` 或 Top-N 的 `--metric`：
   ```bash
   adex ks campaign-reports summary --range 30d --group-by campaign_id --order-by charge --order-desc
   adex ks campaigns top --range 30d --metric charge --limit 10
   ```

## 参考

- [adex-ks](../SKILL.md) — 快手广告全部命令
- [adex-ks-reports](adex-ks-reports.md) — 日报表和汇总报表
- [adex-shared](../../adex-shared/SKILL.md) — 认证和全局参数
