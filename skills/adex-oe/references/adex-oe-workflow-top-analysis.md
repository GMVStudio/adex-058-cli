# 工作流：Top-N 消耗分析

> **前置条件：** 先阅读 [`../adex-shared/SKILL.md`](../../adex-shared/SKILL.md) 了解安装、认证、共享 flags。

从大盘概览到逐层下钻的完整消耗分析工作流。适用于用户想了解"最近一段时间哪些投放消耗最多、效果如何"的场景。

## 场景

用户说：
- "最近 30 天投放情况怎么样？"
- "哪个项目消耗最多？下钻看看"
- "帮我分析一下最近的广告投放数据"
- "预算使用情况怎么样？"

## 步骤 1：大盘概览

先用 `dashboard` 快速了解整体投放情况：

```bash
adex oe dashboard --tenant 6 --range 30d
```

关注返回结果中的：
- 总消耗金额
- 活跃账户数量
- 账户消耗排名（`rankings`）

## 步骤 2：账户日趋势

查看账户层级的日报表，了解消耗趋势：

```bash
adex oe account-reports daily --tenant 6 --range 30d --format table
```

如果需要按小时粒度查看某一天的波动：

```bash
adex oe account-reports daily --tenant 6 --range 1d --stat-hour 12 --format table
```

## 步骤 3：项目 Top-N 排名

找出消耗最高的项目：

```bash
adex oe projects top --tenant 6 --range 30d --metric charge --limit 10
```

返回结果中 `groupKey` 是项目 ID，`groupName` 是项目名称，`charge` 是总消耗。

## 步骤 4：下钻到单元级别

对消耗最高的项目，下钻查看其下的单元：

```bash
# 方法 A：查看该项目下的单元列表
adex oe units --tenant 6 --project <PROJECT_ID> --format table

# 方法 B：直接看单元级 Top-N（不限定项目）
adex oe units top --tenant 6 --range 30d --metric charge --limit 20
```

## 步骤 5：查看明细日报表

对关注的项目/单元，查看其日报表明细：

```bash
# 项目日报表
adex oe project-reports daily --tenant 6 --range 30d --project <PROJECT_ID> --format table

# 单元日报表
adex oe unit-reports daily --tenant 6 --range 30d --promotion <PROMOTION_ID> --format table
```

## 步骤 6：汇总对比

按维度分组汇总，对比各项目/单元的总消耗：

```bash
# 项目汇总排名
adex oe project-reports summary --tenant 6 --range 30d --group-by project_id --order-by charge --order-desc

# 单元汇总排名
adex oe unit-reports summary --tenant 6 --range 30d --group-by promotion_id --order-by charge --order-desc
```

## 步骤 7：预算 vs 实际消耗

查看各账户的预算使用情况，了解预算执行效率：

```bash
# 全部账户
adex oe account-budget-vs-actual --tenant 6 --range 30d --format table

# 单个账户
adex oe account-budget-vs-actual --tenant 6 --advertiser <ADVERTISER_ID> --range 30d
```

## 变体：按转化指标分析

如果用户关注转化而非消耗，将 `--metric charge` 替换为转化指标：

```bash
# 先查可用指标
adex oe report-metric-meta --level project --enabled 1 --jq '.items[].field'

# 用转化指标排名
adex oe projects top --tenant 6 --range 30d --metric convert_cnt --limit 10
adex oe project-reports summary --tenant 6 --range 30d --group-by project_id --order-by convert_cnt --order-desc
```

## 变体：按特定广告主分析

如果用户只关注某个广告主的数据，全程加 `--advertiser`：

```bash
adex oe dashboard --tenant 6 --range 30d
adex oe projects top --tenant 6 --range 30d --metric charge --advertiser 1866874042754522 --limit 10
adex oe project-reports summary --tenant 6 --range 30d --group-by project_id --advertiser 1866874042754522
```

## 完整示例

```bash
# 1. 大盘
adex oe dashboard --tenant 6 --range 30d

# 2. 账户趋势
adex oe account-reports daily --tenant 6 --range 30d --format table

# 3. 项目 Top 10
adex oe projects top --tenant 6 --range 30d --metric charge --limit 10

# 4. 下钻 Top 项目的单元
adex oe units --tenant 6 --project 7650479670059647030 --format table

# 5. Top 项目日报表
adex oe project-reports daily --tenant 6 --range 30d --project 7650479670059647030 --format table

# 6. 汇总对比
adex oe project-reports summary --tenant 6 --range 30d --group-by project_id --order-by charge --order-desc

# 7. 预算使用情况
adex oe account-budget-vs-actual --tenant 6 --range 30d --format table
```

## 参考

- [adex-oe](../SKILL.md) — 巨量引擎广告全部命令
- [adex-oe-dashboard](adex-oe-dashboard.md) — 租户级概览
- [adex-oe-projects](adex-oe-projects.md) — 项目命令
- [adex-oe-units](adex-oe-units.md) — 单元命令
- [adex-oe-reports](adex-oe-reports.md) — 日报表和汇总报表
- [adex-oe-budget-vs-actual](adex-oe-budget-vs-actual.md) — 预算 vs 实际消耗
- [adex-oe-metric-meta](adex-oe-metric-meta.md) — 指标元数据
- [adex-shared](../../adex-shared/SKILL.md) — 认证和全局参数
