---
name: adex-ks
version: 0.3.0
description: "快手广告数据查询：账户、计划、组、创意的列表/详情/Top-N 排名，日/汇总报表，指标元数据，租户级概览。当用户需要查询快手广告投放数据、消耗、报表或排名时使用。"
metadata:
  requires:
    bins: ["adex"]
  cliHelp: "adex ks --help"
---

# ks — 快手广告数据

开始前先读 [`../adex-shared/SKILL.md`](../adex-shared/SKILL.md)（安装、认证、共享 flags）。

## 命令总览

| 命令 | 说明 | API |
|------|------|-----|
| `ks accounts` | 广告账户列表 | `GET /v1/ks/ad-accounts` |
| `ks campaigns` | 广告计划列表 | `GET /v1/ks/campaigns` |
| `ks campaigns top` | 计划 Top-N 排名 | `GET /v1/ks/campaigns/top` |
| `ks campaigns get <id>` | 计划详情 | `GET /v1/ks/campaigns/{id}` |
| `ks units` | 广告组列表 | `GET /v1/ks/units` |
| `ks units top` | 组 Top-N 排名 | `GET /v1/ks/units/top` |
| `ks units get <id>` | 组详情 | `GET /v1/ks/units/{id}` |
| `ks creatives` | 创意列表 | `GET /v1/ks/creatives` |
| `ks creatives top` | 创意 Top-N 排名 | `GET /v1/ks/creatives/top` |
| `ks creatives get <biz_key>` | 创意详情 | `GET /v1/ks/creatives/{biz_key}` |
| `ks account-reports daily` | 账户日报表 | `GET /v1/ks/account-reports/daily` |
| `ks account-reports summary` | 账户汇总报表 | `GET /v1/ks/account-reports/summary` |
| `ks campaign-reports daily` | 计划日报表 | `GET /v1/ks/campaign-reports/daily` |
| `ks campaign-reports summary` | 计划汇总报表 | `GET /v1/ks/campaign-reports/summary` |
| `ks unit-reports daily` | 组日报表 | `GET /v1/ks/unit-reports/daily` |
| `ks unit-reports summary` | 组汇总报表 | `GET /v1/ks/unit-reports/summary` |
| `ks creative-reports daily` | 创意日报表 | `GET /v1/ks/creative-reports/daily` |
| `ks creative-reports summary` | 创意汇总报表 | `GET /v1/ks/creative-reports/summary` |
| `ks report-metric-meta` | 报表指标元数据 | `GET /v1/ks/report-metric-meta` |
| `ks dashboard` | 租户级概览 | `GET /v1/ks/dashboard` |

## accounts — 广告账户

```bash
adex ks accounts --tenant 6 --page-size 20
adex ks accounts --tenant 6 --order-by balance --order-desc --format table
adex ks accounts --tenant 6 --page-all --jq '.items[].advertiserId'
```

| Flag | 说明 |
|------|------|
| `--advertiser` | 广告主 ID 过滤 |
| `--account-name` | 账户名模糊匹配 |
| `--account-type` | 账户类型过滤 |
| `--auth-status` | 授权状态过滤 |
| `--delivery-status` | 投放状态过滤 |
| `--active-status` | 活跃状态过滤 |
| `--owner-user` | 归属用户 ID 过滤 |

## campaigns — 广告计划

列表、Top-N 排名、详情三个子命令。

```bash
# 列表
adex ks campaigns --tenant 6 --page-size 20
adex ks campaigns --tenant 6 --put-status 1 --format table

# Top-N（按消耗排名）
adex ks campaigns top --tenant 6 --range 30d --metric charge --limit 10

# 详情
adex ks campaigns get 9899931248 --tenant 6
```

列表 Flags：

| Flag | 说明 |
|------|------|
| `--advertiser` | 广告主 ID 过滤 |
| `--campaign` | 计划 ID 过滤 |
| `--campaign-name` | 计划名模糊匹配 |
| `--put-status` | 投放状态 1=启用/2=暂停/3=删除（0=全部） |
| `--status` | 计划状态（0=全部） |
| `--campaign-type` | 计划类型（0=全部） |

Top Flags：

| Flag | 说明 |
|------|------|
| `--advertiser` | 广告主 ID 过滤 |
| `--metric` | 排名指标（charge/convert_cnt/active...） |
| `--source` | 数据源过滤 |
| `--limit` | 返回行数（最大 100） |
| `--order-desc` | 降序（默认 true） |
| `--range` / `--begin` / `--end` | 日期范围（必需） |

## units — 广告组

```bash
adex ks units --tenant 6 --campaign 9899931248 --format table
adex ks units top --tenant 6 --range 30d --metric conversion_num --limit 20
adex ks units get 29638466721 --tenant 6
```

列表 Flags：

| Flag | 说明 |
|------|------|
| `--advertiser` | 广告主 ID 过滤 |
| `--unit` | 组 ID 过滤 |
| `--campaign` | 计划 ID 过滤 |
| `--unit-name` | 组名模糊匹配 |
| `--put-status` | 投放状态（0=全部） |
| `--status` | 组状态（0=全部） |

## creatives — 创意

`get` 子命令使用 `biz_key` 作为路径参数（如 `p:29637782154`）。

```bash
adex ks creatives --tenant 6 --unit 29638466721 --format table
adex ks creatives top --tenant 6 --range 30d --metric charge --limit 10
adex ks creatives get p:29637782154 --tenant 6
```

列表 Flags：

| Flag | 说明 |
|------|------|
| `--advertiser` | 广告主 ID 过滤 |
| `--unit` | 组 ID 过滤 |
| `--campaign` | 计划 ID 过滤 |
| `--creative` | 创意 ID 过滤 |
| `--creative-name` | 创意名模糊匹配 |
| `--creative-type` | 创意类型过滤 |
| `--put-status` | 投放状态（0=全部） |
| `--status` | 创意状态（0=全部） |

## 报表（daily / summary）

每个资源层级有 `daily`（日粒度）和 `summary`（汇总）两个子命令。

```bash
# 日报表
adex ks account-reports daily --tenant 6 --range 30d --page-size 20
adex ks campaign-reports daily --tenant 6 --begin 2026-07-01 --end 2026-07-31 --format table

# 汇总报表（不分组 = 单行总计）
adex ks campaign-reports summary --tenant 6 --range 30d
adex ks campaign-reports summary --tenant 6 --range 30d --group-by campaign_id --order-by charge --order-desc
```

Daily 共享 Flags：

| Flag | 说明 |
|------|------|
| `--advertiser` | 广告主 ID 过滤 |
| `--source` | 数据源过滤 |
| `--stat-hour` | 小时粒度（-1=全部） |
| `--range` / `--begin` / `--end` | 日期范围（可选） |
| `--order-by` | 排序字段（默认 stat_date） |

Summary 共享 Flags：

| Flag | 说明 |
|------|------|
| `--advertiser` | 广告主 ID 过滤 |
| `--group-by` | 分组维度（留空=单行总计） |
| `--source` | 数据源过滤 |
| `--range` / `--begin` / `--end` | 日期范围（必需） |
| `--order-by` | 排序字段（默认 charge） |

各资源 daily 额外 Flags：

| 资源 | 额外 Flags |
|------|-----------|
| account-reports | `--account-name` |
| campaign-reports | `--campaign`, `--campaign-name`, `--status` |
| unit-reports | `--unit`, `--campaign`, `--unit-name`, `--status` |
| creative-reports | `--creative`, `--unit`, `--campaign`, `--creative-name`, `--status` |

各资源 summary 的 `--group-by` 支持值：

| 资源 | group-by |
|------|----------|
| account-reports | `advertiser_id` |
| campaign-reports | `campaign_id` |
| unit-reports | `unit_id` |
| creative-reports | `creative_id` |

## report-metric-meta — 指标元数据

不需要 `--tenant`。

```bash
adex ks report-metric-meta --level account --page-size 50
adex ks report-metric-meta --level campaign --enabled 1 --page-size 50
```

| Flag | 说明 |
|------|------|
| `--level` | 维度：account/campaign/unit/creative |
| `--group-name` | 指标组名过滤 |
| `--field` | 字段名模糊匹配 |
| `--enabled` | 0=全部 / 1=启用 / 2=禁用 |
| `--sortable` | 0=全部 / 1=可排序 / 2=不可排序 |

## dashboard — 租户级概览

```bash
adex ks dashboard --tenant 6 --range 30d
adex ks dashboard --tenant 6 --begin 2026-06-01 --end 2026-06-30 --ranking-limit 10
```

| Flag | 说明 |
|------|------|
| `--source` | 数据源过滤 |
| `--ranking-limit` | 账户排名数量（默认 10，最大 100） |
| `--range` / `--begin` / `--end` | 日期范围（必需） |
