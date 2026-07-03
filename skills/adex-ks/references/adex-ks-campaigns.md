# ks campaigns — 广告计划

> **前置条件：** 先阅读 [`../adex-shared/SKILL.md`](../../adex-shared/SKILL.md) 了解安装、认证、共享 flags。

广告计划是快手广告的第一级投放对象，包含列表、Top-N 排名、详情三个子命令。

- `ks campaigns` — 计划列表（`GET /v1/ks/campaigns`）
- `ks campaigns top` — 按指标排名 Top-N（`GET /v1/ks/campaigns/top`）
- `ks campaigns get <id>` — 计划详情（`GET /v1/ks/campaigns/{id}`）

## 列表

```bash
# 基本列表
adex ks campaigns --page-size 20

# 按投放状态筛选（1=启用）
adex ks campaigns --put-status 1 --format table

# 按广告主筛选
adex ks campaigns --advertiser 1234567890 --format table

# 按计划名模糊匹配
adex ks campaigns --campaign-name "品牌" --format table

# 聚合所有页
adex ks campaigns --page-all --jq '.items[].campaignId'
```

### 列表 Flags

| Flag | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--advertiser` | string | — | 广告主 ID 过滤 |
| `--campaign` | string | — | 计划 ID 精确过滤 |
| `--campaign-name` | string | — | 计划名模糊匹配 |
| `--put-status` | int | 0 | 投放状态：1=启用 / 2=暂停 / 3=删除（0=全部） |
| `--status` | int | 0 | 计划状态（0=全部） |
| `--campaign-type` | int | 0 | 计划类型（0=全部） |

### 列表共享 Flags

`--tenant`（可选）、`--page-size`、`--page-token`、`--page-all`、`--order-by`（默认 `id`）、`--order-desc`、`--jq`、`--format`、`--dry-run`

### Table 列

| 列 | 字段 |
|----|------|
| ID | `id` |
| Campaign ID | `campaignId` |
| Campaign Name | `campaignName` |
| Advertiser ID | `advertiserId` |
| Put Status | `putStatus` |
| Status | `status` |
| Campaign Type | `campaignType` |

## Top-N 排名

按指定指标在日期范围内对计划排名。日期范围为**必需**参数。

```bash
# 按消耗排名 Top 10
adex ks campaigns top --range 30d --metric charge --limit 10

# 按转化数排名 Top 20
adex ks campaigns top --range 7d --metric conversion_num --limit 20

# 按广告主筛选后排名
adex ks campaigns top --range 30d --metric charge --advertiser 1234567890 --limit 10

# 升序排名（消耗最少的计划）
adex ks campaigns top --range 30d --metric charge --order-desc=false --limit 10

# 显式日期范围
adex ks campaigns top --begin 2026-06-01 --end 2026-06-30 --metric charge --limit 10
```

### Top Flags

| Flag | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--advertiser` | string | — | 广告主 ID 过滤 |
| `--metric` | string | `charge` | 排名指标（`charge` / `conversion_num` / `active` ...） |
| `--source` | string | — | 数据源过滤 |
| `--limit` | int | 20 | 返回行数（最大 100） |
| `--order-desc` | bool | true | 降序（true=Top 消耗最多优先） |
| `--range` | string | — | 相对日期范围如 `7d` / `4w` / `1m`（**必需**） |
| `--begin` | string | — | 起始日期（YYYY-MM-DD），与 `--end` 配合使用 |
| `--end` | string | — | 结束日期（YYYY-MM-DD） |

> `--range` 优先于 `--begin` / `--end`。日期范围至少需要一种指定方式。

### Top 输出

```json
{
  "hasMore": false,
  "nextPageToken": "",
  "items": [
    {
      "groupKey": "9899931248",
      "groupName": "品牌推广计划A",
      "charge": 50000.00,
      "rowCount": 30
    }
  ]
}
```

| 列 | 字段 |
|----|------|
| Group Key | `groupKey` |
| Group Name | `groupName` |
| Charge | `charge` |
| Row Count | `rowCount` |

## 详情

```bash
# 查看计划详情
adex ks campaigns get 9899931248

# pretty 格式输出
adex ks campaigns get 9899931248 --format pretty

# 提取特定字段
adex ks campaigns get 9899931248 --jq '.campaignName'
```

### 详情参数

| 参数 | 位置 | 必填 | 说明 |
|------|------|------|------|
| `<campaign_id>` | positional | 是 | 计划 ID（路径参数） |
| `--tenant` | flag | 否 | 租户 ID（可选；缺省使用默认租户） |
| `--jq` | flag | 否 | jq 表达式过滤输出 |
| `--format` | flag | 否 | 输出格式（默认 `json`） |

### 详情输出

返回单个计划对象的完整信息，包含基础属性和扩展元数据。

## 使用场景

- **查看所有启用计划**：`--put-status 1 --format table`
- **找消耗最高的计划**：`top --range 30d --metric charge --limit 10`
- **找转化最好的计划**：`top --range 7d --metric conversion_num --limit 20`
- **查看计划完整信息**：`get <campaign_id> --format pretty`
- **获取计划 ID 列表**：`--page-all --jq '.items[].campaignId'` 供后续 `units` 命令使用

## 层级关系

- 计划归属于广告账户（`advertiserId`），可用 `--advertiser` 筛选
- 计划下包含广告组（units），可用 `adex ks units --campaign <ID>` 查询
- 计划报表可用 `adex ks campaign-reports daily/summary` 查询

## 参考

- [adex-ks](../SKILL.md) — 快手广告全部命令
- [adex-ks-units](adex-ks-units.md) — 广告组命令（下钻）
- [adex-ks-reports](adex-ks-reports.md) — 计划报表
- [adex-shared](../../adex-shared/SKILL.md) — 认证和全局参数
