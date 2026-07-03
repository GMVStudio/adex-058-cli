# oe projects — 项目

> **前置条件：** 先阅读 [`../adex-shared/SKILL.md`](../../adex-shared/SKILL.md) 了解安装、认证、共享 flags。

项目是巨量引擎广告的第一级投放对象，包含列表、Top-N 排名、详情三个子命令。

- `oe projects` — 项目列表（`GET /v1/oe/projects`）
- `oe projects top` — 按指标排名 Top-N（`GET /v1/oe/projects/top`）
- `oe projects get <id>` — 项目详情（`GET /v1/oe/projects/{id}`）

## 列表

```bash
# 基本列表
adex oe projects --page-size 20

# 按操作状态筛选（ENABLE=启用）
adex oe projects --opt-status ENABLE --format table

# 按广告主筛选
adex oe projects --advertiser 1866874042754522 --format table

# 按项目名模糊匹配
adex oe projects --name "品牌" --format table

# 聚合所有页
adex oe projects --page-all --jq '.items[].projectId'
```

### 列表 Flags

| Flag | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--advertiser` | string | — | 广告主 ID 过滤 |
| `--project` | string | — | 项目 ID 精确过滤 |
| `--name` | string | — | 项目名模糊匹配 |
| `--opt-status` | string | — | 操作状态 `ENABLE` / `DISABLE` |
| `--status-first` | string | — | 一级状态过滤 |
| `--delivery-mode` | string | — | 投放模式过滤 |
| `--landing-type` | string | — | 落地页类型过滤 |

### 列表共享 Flags

`--tenant`（可选）、`--page-size`、`--page-token`、`--page-all`、`--order-by`（默认 `id`）、`--order-desc`、`--jq`、`--format`、`--dry-run`

### Table 列

| 列 | 字段 |
|----|------|
| ID | `id` |
| Project ID | `projectId` |
| Name | `name` |
| Advertiser ID | `advertiserId` |
| Opt Status | `optStatus` |
| Status First | `statusFirst` |
| Delivery Mode | `deliveryMode` |
| Landing Type | `landingType` |

## Top-N 排名

按指定指标在日期范围内对项目排名。日期范围为**必需**参数。

```bash
# 按消耗排名 Top 10
adex oe projects top --range 30d --metric charge --limit 10

# 按转化数排名 Top 20
adex oe projects top --range 7d --metric convert_cnt --limit 20

# 按广告主筛选后排名
adex oe projects top --range 30d --metric charge --advertiser 1866874042754522 --limit 10

# 升序排名（消耗最少的项目）
adex oe projects top --range 30d --metric charge --order-desc=false --limit 10

# 显式日期范围
adex oe projects top --begin 2026-06-01 --end 2026-06-30 --metric charge --limit 10
```

### Top Flags

| Flag | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--advertiser` | string | — | 广告主 ID 过滤 |
| `--metric` | string | `charge` | 排名指标（`charge` / `convert_cnt` / `active` ...） |
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
      "groupKey": "7650479670059647030",
      "groupName": "品牌推广项目A",
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
# 查看项目详情
adex oe projects get 7650479670059647030

# pretty 格式输出
adex oe projects get 7650479670059647030 --format pretty

# 提取特定字段
adex oe projects get 7650479670059647030 --jq '.name'
```

### 详情参数

| 参数 | 位置 | 必填 | 说明 |
|------|------|------|------|
| `<project_id>` | positional | 是 | 项目 ID（路径参数） |
| `--tenant` | flag | 否 | 租户 ID（可选；缺省使用默认租户） |
| `--jq` | flag | 否 | jq 表达式过滤输出 |
| `--format` | flag | 否 | 输出格式（默认 `json`） |

### 详情输出

返回单个项目对象的完整信息，包含基础属性和扩展元数据。

## 使用场景

- **查看所有启用项目**：`--opt-status ENABLE --format table`
- **找消耗最高的项目**：`top --range 30d --metric charge --limit 10`
- **找转化最好的项目**：`top --range 7d --metric convert_cnt --limit 20`
- **查看项目完整信息**：`get <project_id> --format pretty`
- **获取项目 ID 列表**：`--page-all --jq '.items[].projectId'` 供后续 `units` 命令使用

## 层级关系

- 项目归属于广告账户（`advertiserId`），可用 `--advertiser` 筛选
- 项目下包含单元（units），可用 `adex oe units --project <ID>` 查询
- 项目报表可用 `adex oe project-reports daily/summary` 查询

## 参考

- [adex-oe](../SKILL.md) — 巨量引擎广告全部命令
- [adex-oe-units](adex-oe-units.md) — 单元命令（下钻）
- [adex-oe-reports](adex-oe-reports.md) — 项目报表
- [adex-shared](../../adex-shared/SKILL.md) — 认证和全局参数
