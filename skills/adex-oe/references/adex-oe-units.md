# oe units — 单元/广告

> **前置条件：** 先阅读 [`../adex-shared/SKILL.md`](../../adex-shared/SKILL.md) 了解安装、认证、共享 flags。

单元（promotion）是巨量引擎广告的第二级投放对象，隶属于项目（project）。包含列表、Top-N 排名、详情三个子命令。

- `oe units` — 单元列表（`GET /v1/oe/units`）
- `oe units top` — 按指标排名 Top-N（`GET /v1/oe/units/top`）
- `oe units get <id>` — 单元详情（`GET /v1/oe/units/{id}`）

> **注意：** 巨量引擎中"单元"也称为"promotion"，路径参数和 API 字段使用 `promotion_id`。

## 列表

```bash
# 基本列表
adex oe units --page-size 20

# 按项目 ID 筛选（最常见用法）
adex oe units --project 7650479670059647030 --format table

# 按广告主筛选
adex oe units --advertiser 1866874042754522 --format table

# 按单元名模糊匹配
adex oe units --name "信息流" --format table

# 按操作状态筛选
adex oe units --opt-status ENABLE --format table

# 聚合所有页
adex oe units --project 7650479670059647030 --page-all --jq '.items[].promotionId'
```

### 列表 Flags

| Flag | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--advertiser` | string | — | 广告主 ID 过滤 |
| `--promotion` | string | — | 单元（promotion）ID 精确过滤 |
| `--project` | string | — | 项目 ID 过滤（用于项目下钻） |
| `--name` | string | — | 单元名模糊匹配 |
| `--opt-status` | string | — | 操作状态 `ENABLE` / `DISABLE` |
| `--status-first` | string | — | 一级状态过滤 |
| `--learning-phase` | string | — | 学习阶段过滤 |

### 列表共享 Flags

`--tenant`（可选）、`--page-size`、`--page-token`、`--page-all`、`--order-by`（默认 `id`）、`--order-desc`、`--jq`、`--format`、`--dry-run`

### Table 列

| 列 | 字段 |
|----|------|
| ID | `id` |
| Promotion ID | `promotionId` |
| Name | `name` |
| Project ID | `projectId` |
| Advertiser ID | `advertiserId` |
| Opt Status | `optStatus` |
| Status First | `statusFirst` |
| Learning Phase | `learningPhase` |

## Top-N 排名

按指定指标在日期范围内对单元排名。日期范围为**必需**参数。

```bash
# 按消耗排名 Top 20
adex oe units top --range 30d --metric charge --limit 20

# 按转化数排名
adex oe units top --range 7d --metric convert_cnt --limit 20

# 在特定广告主下排名
adex oe units top --range 30d --metric charge --advertiser 1866874042754522 --limit 10

# 升序排名（消耗最少的单元）
adex oe units top --range 30d --metric charge --order-desc=false --limit 10
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
      "groupKey": "7650483929670156288",
      "groupName": "信息流广告单元A",
      "charge": 15000.00,
      "rowCount": 30
    }
  ]
}
```

## 详情

```bash
# 查看单元详情
adex oe units get 7650483929670156288

# pretty 格式输出
adex oe units get 7650483929670156288 --format pretty

# 提取特定字段
adex oe units get 7650483929670156288 --jq '.name'
```

### 详情参数

| 参数 | 位置 | 必填 | 说明 |
|------|------|------|------|
| `<promotion_id>` | positional | 是 | 单元（promotion）ID（路径参数） |
| `--tenant` | flag | 否 | 租户 ID（可选；缺省使用默认租户） |
| `--jq` | flag | 否 | jq 表达式过滤输出 |
| `--format` | flag | 否 | 输出格式（默认 `json`） |

## 使用场景

- **项目下钻查看单元**：`--project <PROJECT_ID> --format table`
- **找消耗最高的单元**：`top --range 30d --metric charge --limit 20`
- **查看单元完整信息**：`get <promotion_id> --format pretty`
- **获取单元 ID 列表**：`--page-all --jq '.items[].promotionId'` 供后续报表命令使用
- **按学习阶段筛选**：`--learning-phase "LEARNING"` 查看处于学习期的单元

## 层级关系

- 单元隶属于项目（`projectId`），可用 `--project` 筛选
- 单元归属于广告账户（`advertiserId`），可用 `--advertiser` 筛选
- 单元报表可用 `adex oe unit-reports daily/summary` 查询
- 巨量引擎没有创意层级，单元是最细粒度的投放对象

## 参考

- [adex-oe](../SKILL.md) — 巨量引擎广告全部命令
- [adex-oe-projects](adex-oe-projects.md) — 项目命令（上级）
- [adex-oe-reports](adex-oe-reports.md) — 单元报表
- [adex-shared](../../adex-shared/SKILL.md) — 认证和全局参数
