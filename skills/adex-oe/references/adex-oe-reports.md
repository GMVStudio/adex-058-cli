# oe reports — 报表（daily / summary）

> **前置条件：** 先阅读 [`../adex-shared/SKILL.md`](../../adex-shared/SKILL.md) 了解安装、认证、共享 flags。

巨量引擎报表按资源层级分为三组，每组包含 `daily`（日粒度）和 `summary`（汇总）两个子命令。

| 资源 | daily 命令 | summary 命令 | group-by |
|------|-----------|-------------|----------|
| 账户 | `oe account-reports daily` | `oe account-reports summary` | `advertiser_id` |
| 项目 | `oe project-reports daily` | `oe project-reports summary` | `project_id` |
| 单元 | `oe unit-reports daily` | `oe unit-reports summary` | `promotion_id` |

## daily — 日报表

按天粒度返回报表数据，每行对应一个统计日期。日期范围为**可选**（不传则返回全部历史数据）。

```bash
# 账户日报表 — 最近 30 天
adex oe account-reports daily --tenant 6 --range 30d --page-size 20

# 项目日报表 — 显式日期范围
adex oe project-reports daily --tenant 6 --begin 2026-07-01 --end 2026-07-31 --format table

# 按广告主筛选
adex oe account-reports daily --tenant 6 --range 7d --advertiser 1866874042754522 --format table

# 按项目 ID 筛选
adex oe project-reports daily --tenant 6 --range 30d --project 7650479670059647030 --format table

# 按小时粒度查看
adex oe account-reports daily --tenant 6 --range 7d --stat-hour 12 --format table

# 聚合所有页
adex oe account-reports daily --tenant 6 --range 30d --page-all --jq '.items[].charge'
```

### Daily 共享 Flags

| Flag | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--advertiser` | string | — | 广告主 ID 过滤 |
| `--source` | string | — | 数据源过滤 |
| `--stat-hour` | int | -1 | 小时粒度（-1=全部，0-23=指定小时） |
| `--range` | string | — | 相对日期范围如 `7d` / `4w` / `1m`（可选） |
| `--begin` | string | — | 起始日期（YYYY-MM-DD，可选） |
| `--end` | string | — | 结束日期（YYYY-MM-DD，可选） |
| `--order-by` | string | `stat_date` | 排序字段 |
| `--order-desc` | bool | true | 降序排序 |

### 各资源 daily 额外 Flags

| 资源 | 额外 Flags |
|------|-----------|
| `account-reports` | — |
| `project-reports` | `--project`（项目 ID）、`--project-name`（项目名模糊匹配） |
| `unit-reports` | `--project`（项目 ID）、`--promotion`（单元 ID）、`--promotion-name`（单元名模糊匹配） |

### Daily 输出

```json
{
  "hasMore": true,
  "nextPageToken": "abc123",
  "items": [
    {
      "id": 1,
      "advertiserId": "1866874042754522",
      "statDate": "2026-06-15",
      "statHour": -1,
      "charge": 5000.00
    }
  ]
}
```

### Daily Table 列

| 资源 | 列 |
|------|----|
| account-reports | id, advertiserId, statDate, statHour, charge |
| project-reports | id, advertiserId, projectId, projectName, statDate, charge |
| unit-reports | id, advertiserId, promotionId, promotionName, statDate, charge |

## summary — 汇总报表

在日期范围内按维度汇总数据。日期范围为**必需**参数。

- 不传 `--group-by` → 返回单行总计
- 传 `--group-by <dimension>` → 按维度分组返回多行

```bash
# 账户汇总 — 单行总计
adex oe account-reports summary --tenant 6 --range 30d

# 项目汇总 — 按项目 ID 分组
adex oe project-reports summary --tenant 6 --range 30d --group-by project_id --order-by charge --order-desc

# 单元汇总 — 按单元 ID 分组
adex oe unit-reports summary --tenant 6 --range 30d --group-by promotion_id --order-by charge --order-desc

# 按广告主筛选
adex oe project-reports summary --tenant 6 --range 30d --group-by project_id --advertiser 1866874042754522

# 显式日期范围
adex oe account-reports summary --tenant 6 --begin 2026-06-01 --end 2026-06-30
```

### Summary 共享 Flags

| Flag | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--advertiser` | string | — | 广告主 ID 过滤 |
| `--group-by` | string | — | 分组维度（留空=单行总计） |
| `--source` | string | — | 数据源过滤 |
| `--range` | string | — | 相对日期范围如 `7d` / `4w` / `1m`（**必需**） |
| `--begin` | string | — | 起始日期（YYYY-MM-DD，与 `--end` 配合） |
| `--end` | string | — | 结束日期（YYYY-MM-DD） |
| `--order-by` | string | `charge` | 排序字段 |
| `--order-desc` | bool | true | 降序排序 |

### Summary group-by 支持值

| 资源 | `--group-by` 值 |
|------|-----------------|
| `account-reports` | `advertiser_id` |
| `project-reports` | `project_id` |
| `unit-reports` | `promotion_id` |

### Summary 输出

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

## 使用场景

- **看日趋势**：`daily --range 30d --format table` 按天查看消耗变化
- **看汇总排名**：`summary --range 30d --group-by <dimension> --order-by charge --order-desc` 按消耗排名
- **看总消耗**：`summary --range 30d`（不传 `--group-by`）返回单行总计
- **按小时看**：`daily --range 7d --stat-hour 12` 只看每天 12 点的数据
- **下钻分析**：先 `account-reports summary` 看哪个账户消耗多，再 `project-reports summary --advertiser <ID>` 看该账户下哪个项目消耗多

## 日期范围说明

| 方式 | 示例 | 说明 |
|------|------|------|
| `--range` | `7d` / `4w` / `1m` | 相对范围，`--range` 优先于 `--begin` / `--end` |
| `--begin` + `--end` | `--begin 2026-06-01 --end 2026-06-30` | 显式日期范围 |
| 不传 | — | daily 返回全部历史数据；summary 必须传日期范围 |

- `daily` 的日期范围是**可选的**，不传则返回全部历史数据
- `summary` 的日期范围是**必需的**，不传会返回验证错误

## 参考

- [adex-oe](../SKILL.md) — 巨量引擎广告全部命令
- [adex-oe-metric-meta](adex-oe-metric-meta.md) — 查看报表可用指标字段
- [adex-shared](../../adex-shared/SKILL.md) — 认证和全局参数
