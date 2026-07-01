---
name: adex-oe
version: 0.3.0
description: "巨量引擎广告数据查询：账户、项目、单元的列表/详情/Top-N 排名，日/汇总报表，指标元数据，租户级概览，预算 vs 实际消耗对比。当用户需要查询巨量/Oceanengine 广告投放数据、消耗、报表或排名时使用。"
metadata:
  requires:
    bins: ["adex"]
  cliHelp: "adex oe --help"
---

# oe — 巨量引擎广告数据

开始前先读 [`../adex-shared/SKILL.md`](../adex-shared/SKILL.md)（安装、认证、共享 flags）。

## 命令总览

| 命令 | 说明 | API |
|------|------|-----|
| `oe accounts` | 广告账户列表 | `GET /v1/oe/ad-accounts` |
| `oe projects` | 项目列表 | `GET /v1/oe/projects` |
| `oe projects top` | 项目 Top-N 排名 | `GET /v1/oe/projects/top` |
| `oe projects get <id>` | 项目详情 | `GET /v1/oe/projects/{id}` |
| `oe units` | 单元列表 | `GET /v1/oe/units` |
| `oe units top` | 单元 Top-N 排名 | `GET /v1/oe/units/top` |
| `oe units get <id>` | 单元详情 | `GET /v1/oe/units/{id}` |
| `oe account-reports daily` | 账户日报表 | `GET /v1/oe/account-reports/daily` |
| `oe account-reports summary` | 账户汇总报表 | `GET /v1/oe/account-reports/summary` |
| `oe project-reports daily` | 项目日报表 | `GET /v1/oe/project-reports/daily` |
| `oe project-reports summary` | 项目汇总报表 | `GET /v1/oe/project-reports/summary` |
| `oe unit-reports daily` | 单元日报表 | `GET /v1/oe/unit-reports/daily` |
| `oe unit-reports summary` | 单元汇总报表 | `GET /v1/oe/unit-reports/summary` |
| `oe report-metric-meta` | 报表指标元数据 | `GET /v1/oe/report-metric-meta` |
| `oe dashboard` | 租户级概览 | `GET /v1/oe/dashboard` |
| `oe account-budget-vs-actual` | 预算 vs 实际消耗 | `GET /v1/oe/account-budget-vs-actual` |

## accounts — 广告账户

```bash
adex oe accounts --tenant 6 --page-size 20
adex oe accounts --tenant 6 --order-by balance --order-desc --format table
adex oe accounts --tenant 6 --page-all --jq '.items[].advertiserId'
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

## projects — 项目

列表、Top-N 排名、详情三个子命令。

```bash
# 列表
adex oe projects --tenant 6 --page-size 20
adex oe projects --tenant 6 --opt-status ENABLE --format table

# Top-N（按消耗排名）
adex oe projects top --tenant 6 --range 30d --metric charge --limit 10

# 详情
adex oe projects get 7650479670059647030 --tenant 6
```

列表 Flags：

| Flag | 说明 |
|------|------|
| `--advertiser` | 广告主 ID 过滤 |
| `--project` | 项目 ID 过滤 |
| `--name` | 项目名模糊匹配 |
| `--opt-status` | 操作状态 ENABLE/DISABLE |
| `--status-first` | 一级状态过滤 |
| `--delivery-mode` | 投放模式过滤 |
| `--landing-type` | 落地页类型过滤 |

Top Flags：

| Flag | 说明 |
|------|------|
| `--advertiser` | 广告主 ID 过滤 |
| `--metric` | 排名指标（charge/convert_cnt/active...） |
| `--source` | 数据源过滤 |
| `--limit` | 返回行数（最大 100） |
| `--order-desc` | 降序（默认 true） |
| `--range` / `--begin` / `--end` | 日期范围（必需） |

## units — 单元/广告

```bash
adex oe units --tenant 6 --project 7650479670059647030 --format table
adex oe units top --tenant 6 --range 30d --metric charge --limit 10
adex oe units get 7650483929670156288 --tenant 6
```

列表 Flags：

| Flag | 说明 |
|------|------|
| `--advertiser` | 广告主 ID 过滤 |
| `--promotion` | 单元（promotion）ID 过滤 |
| `--project` | 项目 ID 过滤 |
| `--name` | 单元名模糊匹配 |
| `--opt-status` | 操作状态 ENABLE/DISABLE |
| `--status-first` | 一级状态过滤 |
| `--learning-phase` | 学习阶段过滤 |

## 报表（daily / summary）

每个资源层级有 `daily`（日粒度）和 `summary`（汇总）两个子命令。

```bash
# 日报表
adex oe account-reports daily --tenant 6 --range 30d --page-size 20
adex oe project-reports daily --tenant 6 --begin 2026-07-01 --end 2026-07-31 --format table

# 汇总报表（不分组 = 单行总计）
adex oe project-reports summary --tenant 6 --range 30d
adex oe project-reports summary --tenant 6 --range 30d --group-by project_id --order-by charge --order-desc
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
| account-reports | — |
| project-reports | `--project`, `--project-name` |
| unit-reports | `--project`, `--promotion`, `--promotion-name` |

各资源 summary 的 `--group-by` 支持值：

| 资源 | group-by |
|------|----------|
| account-reports | `advertiser_id` |
| project-reports | `project_id` |
| unit-reports | `promotion_id` |

## report-metric-meta — 指标元数据

不需要 `--tenant`。

```bash
adex oe report-metric-meta --level account --page-size 50
adex oe report-metric-meta --level project --enabled 1 --page-size 50
```

| Flag | 说明 |
|------|------|
| `--level` | 维度：account/project/unit |
| `--group-name` | 指标组名过滤 |
| `--field` | 字段名模糊匹配 |
| `--enabled` | 0=全部 / 1=启用 / 2=禁用 |

## dashboard — 租户级概览

```bash
adex oe dashboard --tenant 6 --range 30d
adex oe dashboard --tenant 6 --begin 2026-06-01 --end 2026-06-30
```

| Flag | 说明 |
|------|------|
| `--range` / `--begin` / `--end` | 日期范围（必需） |

## account-budget-vs-actual — 预算 vs 实际消耗

对比每个账户的日预算与实际日均消耗。

```bash
adex oe account-budget-vs-actual --tenant 6 --range 30d
adex oe account-budget-vs-actual --tenant 6 --begin 2026-06-01 --end 2026-06-30 --format table
adex oe account-budget-vs-actual --tenant 6 --advertiser 1866874042754522 --range 30d
```

| Flag | 说明 |
|------|------|
| `--advertiser` | 广告主 ID 过滤（单个账户） |
| `--range` / `--begin` / `--end` | 日期范围（必需） |
