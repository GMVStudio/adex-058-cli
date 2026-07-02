# ks reports — 报表（daily / summary）

> **前置条件：** 先阅读 [`../adex-shared/SKILL.md`](../../adex-shared/SKILL.md) 了解安装、认证、共享 flags。

快手广告报表按资源层级分为四组，每组包含 `daily`（日粒度）和 `summary`（汇总）两个子命令。

| 资源 | daily 命令 | summary 命令 | group-by |
|------|-----------|-------------|----------|
| 账户 | `ks account-reports daily` | `ks account-reports summary` | `advertiser_id` |
| 计划 | `ks campaign-reports daily` | `ks campaign-reports summary` | `campaign_id` |
| 组 | `ks unit-reports daily` | `ks unit-reports summary` | `unit_id` |
| 创意 | `ks creative-reports daily` | `ks creative-reports summary` | `creative_id` |

## daily — 日报表

按天粒度返回报表数据，每行对应一个统计日期。日期范围为**可选**（不传则返回全部历史数据）。

```bash
# 账户日报表 — 最近 30 天
adex ks account-reports daily --tenant 6 --range 30d --page-size 20

# 计划日报表 — 显式日期范围
adex ks campaign-reports daily --tenant 6 --begin 2026-07-01 --end 2026-07-31 --format table

# 按广告主筛选
adex ks account-reports daily --tenant 6 --range 7d --advertiser 1234567890 --format table

# 按计划 ID 筛选
adex ks campaign-reports daily --tenant 6 --range 30d --campaign 9899931248 --format table

# 按小时粒度查看
adex ks account-reports daily --tenant 6 --range 7d --stat-hour 12 --format table

# 聚合所有页
adex ks account-reports daily --tenant 6 --range 30d --page-all --jq '.items[].charge'
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
| `account-reports` | `--account-name`（账户名模糊匹配） |
| `campaign-reports` | `--campaign`（计划 ID）、`--campaign-name`（计划名模糊匹配）、`--status`（计划状态） |
| `unit-reports` | `--unit`（组 ID）、`--campaign`（计划 ID）、`--unit-name`（组名模糊匹配）、`--status`（组状态） |
| `creative-reports` | `--creative`（创意 ID）、`--unit`（组 ID）、`--campaign`（计划 ID）、`--creative-name`（创意名模糊匹配）、`--status`（创意状态） |

### Daily 输出

```json
{
  "hasMore": true,
  "nextPageToken": "abc123",
  "items": [
    {
      "id": 1,
      "advertiserId": "1234567890",
      "accountName": "品牌推广账户",
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
| account-reports | id, advertiserId, accountName, statDate, statHour, charge |
| campaign-reports | id, advertiserId, campaignId, campaignName, statDate, charge |
| unit-reports | id, advertiserId, unitId, unitName, statDate, charge |
| creative-reports | id, advertiserId, creativeId, creativeName, statDate, charge |

## summary — 汇总报表

在日期范围内按维度汇总数据。日期范围为**必需**参数。

- 不传 `--group-by` → 返回单行总计
- 传 `--group-by <dimension>` → 按维度分组返回多行

```bash
# 账户汇总 — 单行总计
adex ks account-reports summary --tenant 6 --range 30d

# 计划汇总 — 按计划 ID 分组
adex ks campaign-reports summary --tenant 6 --range 30d --group-by campaign_id --order-by charge --order-desc

# 组汇总 — 按组 ID 分组
adex ks unit-reports summary --tenant 6 --range 30d --group-by unit_id --order-by charge --order-desc

# 创意汇总 — 按创意 ID 分组
adex ks creative-reports summary --tenant 6 --range 30d --group-by creative_id --order-by charge --order-desc

# 按广告主筛选
adex ks campaign-reports summary --tenant 6 --range 30d --group-by campaign_id --advertiser 1234567890

# 显式日期范围
adex ks account-reports summary --tenant 6 --begin 2026-06-01 --end 2026-06-30
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
| `campaign-reports` | `campaign_id` |
| `unit-reports` | `unit_id` |
| `creative-reports` | `creative_id` |

### Summary 输出

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

## 使用场景

- **看日趋势**：`daily --range 30d --format table` 按天查看消耗变化
- **看汇总排名**：`summary --range 30d --group-by <dimension> --order-by charge --order-desc` 按消耗排名
- **看总消耗**：`summary --range 30d`（不传 `--group-by`）返回单行总计
- **按小时看**：`daily --range 7d --stat-hour 12` 只看每天 12 点的数据
- **下钻分析**：先 `account-reports summary` 看哪个账户消耗多，再 `campaign-reports summary --advertiser <ID>` 看该账户下哪个计划消耗多

## 日期范围说明

| 方式 | 示例 | 说明 |
|------|------|------|
| `--range` | `7d` / `4w` / `1m` | 相对范围，`--range` 优先于 `--begin` / `--end` |
| `--begin` + `--end` | `--begin 2026-06-01 --end 2026-06-30` | 显式日期范围 |
| 不传 | — | daily 返回全部历史数据；summary 必须传日期范围 |

- `daily` 的日期范围是**可选的**，不传则返回全部历史数据
- `summary` 的日期范围是**必需的**，不传会返回验证错误

## 参考

- [adex-ks](../SKILL.md) — 快手广告全部命令
- [adex-ks-metric-meta](adex-ks-metric-meta.md) — 查看报表可用指标字段
- [adex-shared](../../adex-shared/SKILL.md) — 认证和全局参数
