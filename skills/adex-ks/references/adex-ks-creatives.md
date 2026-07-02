# ks creatives — 创意

> **前置条件：** 先阅读 [`../adex-shared/SKILL.md`](../../adex-shared/SKILL.md) 了解安装、认证、共享 flags。

创意是快手广告的第三级投放对象，隶属于广告组（unit）。包含列表、Top-N 排名、详情三个子命令。

> **注意：** `get` 子命令使用 `biz_key` 作为路径参数（如 `p:29637782154`），不是 creative_id。

- `ks creatives` — 创意列表（`GET /v1/ks/creatives`）
- `ks creatives top` — 按指标排名 Top-N（`GET /v1/ks/creatives/top`）
- `ks creatives get <biz_key>` — 创意详情（`GET /v1/ks/creatives/{biz_key}`）

## 列表

```bash
# 基本列表
adex ks creatives --tenant 6 --page-size 20

# 按组 ID 筛选（最常见用法）
adex ks creatives --tenant 6 --unit 29638466721 --format table

# 按计划 ID 筛选
adex ks creatives --tenant 6 --campaign 9899931248 --format table

# 按创意名模糊匹配
adex ks creatives --tenant 6 --creative-name "视频" --format table

# 按创意类型筛选
adex ks creatives --tenant 6 --creative-type "VIDEO" --format table

# 聚合所有页
adex ks creatives --tenant 6 --unit 29638466721 --page-all --jq '.items[].creativeId'
```

### 列表 Flags

| Flag | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--advertiser` | string | — | 广告主 ID 过滤 |
| `--unit` | string | — | 组 ID 过滤（用于组下钻） |
| `--campaign` | string | — | 计划 ID 过滤 |
| `--creative` | string | — | 创意 ID 精确过滤 |
| `--creative-name` | string | — | 创意名模糊匹配 |
| `--creative-type` | string | — | 创意类型过滤 |
| `--put-status` | int | 0 | 投放状态：1=启用 / 2=暂停 / 3=删除（0=全部） |
| `--status` | int | 0 | 创意状态（0=全部） |

### 列表共享 Flags

`--tenant`（必需）、`--page-size`、`--page-token`、`--page-all`、`--order-by`（默认 `id`）、`--order-desc`、`--jq`、`--format`、`--dry-run`

### Table 列

| 列 | 字段 |
|----|------|
| ID | `id` |
| Creative ID | `creativeId` |
| Creative Name | `creativeName` |
| Unit ID | `unitId` |
| Campaign ID | `campaignId` |
| Advertiser ID | `advertiserId` |
| Put Status | `putStatus` |
| Status | `status` |

## Top-N 排名

按指定指标在日期范围内对创意排名。日期范围为**必需**参数。

```bash
# 按消耗排名 Top 10
adex ks creatives top --tenant 6 --range 30d --metric charge --limit 10

# 按转化数排名
adex ks creatives top --tenant 6 --range 7d --metric convert_cnt --limit 20

# 在特定广告主下排名
adex ks creatives top --tenant 6 --range 30d --metric charge --advertiser 1234567890 --limit 10
```

### Top Flags

| Flag | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--advertiser` | string | — | 广告主 ID 过滤 |
| `--metric` | string | `charge` | 排名指标（`charge` / `convert_cnt` / `active` ...） |
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
      "groupKey": "29637782154",
      "groupName": "品牌视频创意A",
      "charge": 8000.00,
      "rowCount": 30
    }
  ]
}
```

## 详情

`get` 子命令使用 `biz_key` 作为路径参数，不是 creative_id。`biz_key` 格式通常为 `p:<creative_id>`（如 `p:29637782154`）。

```bash
# 查看创意详情
adex ks creatives get p:29637782154 --tenant 6

# pretty 格式输出
adex ks creatives get p:29637782154 --tenant 6 --format pretty

# 提取特定字段
adex ks creatives get p:29637782154 --tenant 6 --jq '.creativeName'
```

### 详情参数

| 参数 | 位置 | 必填 | 说明 |
|------|------|------|------|
| `<biz_key>` | positional | 是 | 创意业务键（路径参数，格式如 `p:29637782154`） |
| `--tenant` | flag | 是 | 租户 ID |
| `--jq` | flag | 否 | jq 表达式过滤输出 |
| `--format` | flag | 否 | 输出格式（默认 `json`） |

> [!IMPORTANT]
> `get` 的路径参数是 `biz_key`（如 `p:29637782154`），**不是** `creativeId`。如果只有 `creativeId`，需要先从列表命令的返回结果中找到对应的 `biz_key`，或直接使用列表命令查询。

## 使用场景

- **组下钻查看创意**：`--unit <UNIT_ID> --format table`
- **计划下钻查看创意**：`--campaign <CAMPAIGN_ID> --format table`
- **找消耗最高的创意**：`top --range 30d --metric charge --limit 10`
- **查看创意完整信息**：`get p:<creative_id> --format pretty`
- **按类型筛选创意**：`--creative-type "VIDEO"` 只看视频创意

## 层级关系

- 创意隶属于组（`unitId`），可用 `--unit` 筛选
- 创意隶属于计划（`campaignId`），可用 `--campaign` 筛选
- 创意归属于广告账户（`advertiserId`），可用 `--advertiser` 筛选
- 创意报表可用 `adex ks creative-reports daily/summary` 查询

## 参考

- [adex-ks](../SKILL.md) — 快手广告全部命令
- [adex-ks-units](adex-ks-units.md) — 广告组命令（上级）
- [adex-ks-reports](adex-ks-reports.md) — 创意报表
- [adex-shared](../../adex-shared/SKILL.md) — 认证和全局参数
