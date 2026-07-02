# ks units — 广告组

> **前置条件：** 先阅读 [`../adex-shared/SKILL.md`](../../adex-shared/SKILL.md) 了解安装、认证、共享 flags。

广告组是快手广告的第二级投放对象，隶属于计划（campaign）。包含列表、Top-N 排名、详情三个子命令。

- `ks units` — 组列表（`GET /v1/ks/units`）
- `ks units top` — 按指标排名 Top-N（`GET /v1/ks/units/top`）
- `ks units get <id>` — 组详情（`GET /v1/ks/units/{id}`）

## 列表

```bash
# 基本列表
adex ks units --tenant 6 --page-size 20

# 按计划 ID 筛选（最常见用法）
adex ks units --tenant 6 --campaign 9899931248 --format table

# 按广告主筛选
adex ks units --tenant 6 --advertiser 1234567890 --format table

# 按组名模糊匹配
adex ks units --tenant 6 --unit-name "信息流" --format table

# 按投放状态筛选
adex ks units --tenant 6 --put-status 1 --format table

# 聚合所有页
adex ks units --tenant 6 --campaign 9899931248 --page-all --jq '.items[].unitId'
```

### 列表 Flags

| Flag | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--advertiser` | string | — | 广告主 ID 过滤 |
| `--unit` | string | — | 组 ID 精确过滤 |
| `--campaign` | string | — | 计划 ID 过滤（用于计划下钻） |
| `--unit-name` | string | — | 组名模糊匹配 |
| `--put-status` | int | 0 | 投放状态：1=启用 / 2=暂停 / 3=删除（0=全部） |
| `--status` | int | 0 | 组状态（0=全部） |

### 列表共享 Flags

`--tenant`（必需）、`--page-size`、`--page-token`、`--page-all`、`--order-by`（默认 `id`）、`--order-desc`、`--jq`、`--format`、`--dry-run`

### Table 列

| 列 | 字段 |
|----|------|
| ID | `id` |
| Unit ID | `unitId` |
| Unit Name | `unitName` |
| Campaign ID | `campaignId` |
| Advertiser ID | `advertiserId` |
| Put Status | `putStatus` |
| Status | `status` |

## Top-N 排名

按指定指标在日期范围内对组排名。日期范围为**必需**参数。

```bash
# 按消耗排名 Top 20
adex ks units top --tenant 6 --range 30d --metric charge --limit 20

# 按转化数排名
adex ks units top --tenant 6 --range 7d --metric conversion_num --limit 20

# 在特定广告主下排名
adex ks units top --tenant 6 --range 30d --metric charge --advertiser 1234567890 --limit 10

# 升序排名（消耗最少的组）
adex ks units top --tenant 6 --range 30d --metric charge --order-desc=false --limit 10
```

### Top Flags

| Flag | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--advertiser` | string | — | 广告主 ID 过滤 |
| `--metric` | string | `charge` | 排名指标（`charge` / `conversion_num` / `active` ...） |
| `--source` | string | — | 数据源过滤 |
| `--limit` | int | 20 | 返回行数（最大 100） |
| `--order-desc` | bool | true | 降序（true=Top 值最大优先） |
| `--range` | string | — | 相对日期范围如 `7d` / `4w` / `1m`（**必需**） |
| `--begin` | string | — | 起始日期（YYYY-MM-DD） |
| `--end` | string | — | 结束日期（YYYY-MM-DD） |

> `--range` 优先于 `--begin` / `--end`。日期范围至少需要一种指定方式。

### Top 输出

```json
{
  "hasMore": false,
  "nextPageToken": "",
  "items": [
    {
      "groupKey": "29638466721",
      "groupName": "信息流广告组A",
      "charge": 15000.00,
      "rowCount": 30
    }
  ]
}
```

## 详情

```bash
# 查看组详情
adex ks units get 29638466721 --tenant 6

# pretty 格式输出
adex ks units get 29638466721 --tenant 6 --format pretty

# 提取特定字段
adex ks units get 29638466721 --tenant 6 --jq '.unitName'
```

### 详情参数

| 参数 | 位置 | 必填 | 说明 |
|------|------|------|------|
| `<unit_id>` | positional | 是 | 组 ID（路径参数） |
| `--tenant` | flag | 是 | 租户 ID |
| `--jq` | flag | 否 | jq 表达式过滤输出 |
| `--format` | flag | 否 | 输出格式（默认 `json`） |

## 使用场景

- **计划下钻查看组**：`--campaign <CAMPAIGN_ID> --format table`
- **找消耗最高的组**：`top --range 30d --metric charge --limit 20`
- **查看组完整信息**：`get <unit_id> --format pretty`
- **获取组 ID 列表**：`--page-all --jq '.items[].unitId'` 供后续 `creatives` 命令使用

## 层级关系

- 组隶属于计划（`campaignId`），可用 `--campaign` 筛选
- 组归属于广告账户（`advertiserId`），可用 `--advertiser` 筛选
- 组下包含创意（creatives），可用 `adex ks creatives --unit <ID>` 查询
- 组报表可用 `adex ks unit-reports daily/summary` 查询

## 参考

- [adex-ks](../SKILL.md) — 快手广告全部命令
- [adex-ks-campaigns](adex-ks-campaigns.md) — 广告计划命令（上级）
- [adex-ks-creatives](adex-ks-creatives.md) — 创意命令（下钻）
- [adex-ks-reports](adex-ks-reports.md) — 组报表
- [adex-shared](../../adex-shared/SKILL.md) — 认证和全局参数
