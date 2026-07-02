# 工作流：Top-N 消耗分析

> **前置条件：** 先阅读 [`../adex-shared/SKILL.md`](../../adex-shared/SKILL.md) 了解安装、认证、共享 flags。

从大盘概览到逐层下钻的完整消耗分析工作流。适用于用户想了解"最近一段时间哪些投放消耗最多、效果如何"的场景。

## 场景

用户说：
- "最近 30 天投放情况怎么样？"
- "哪个计划消耗最多？下钻看看"
- "帮我分析一下最近的广告投放数据"

## 步骤 1：大盘概览

先用 `dashboard` 快速了解整体投放情况：

```bash
adex ks dashboard --tenant 6 --range 30d
```

关注返回结果中的：
- 总消耗金额
- 活跃账户数量
- 账户消耗排名（`rankings`）

## 步骤 2：账户日趋势

查看账户层级的日报表，了解消耗趋势：

```bash
adex ks account-reports daily --tenant 6 --range 30d --format table
```

如果需要按小时粒度查看某一天的波动：

```bash
adex ks account-reports daily --tenant 6 --range 1d --stat-hour 12 --format table
```

## 步骤 3：计划 Top-N 排名

找出消耗最高的计划：

```bash
adex ks campaigns top --tenant 6 --range 30d --metric charge --limit 10
```

返回结果中 `groupKey` 是计划 ID，`groupName` 是计划名称，`charge` 是总消耗。

## 步骤 4：下钻到组级别

对消耗最高的计划，下钻查看其下的广告组：

```bash
# 方法 A：查看该计划下的组列表
adex ks units --tenant 6 --campaign <CAMPAIGN_ID> --format table

# 方法 B：直接看组级 Top-N（不限定计划）
adex ks units top --tenant 6 --range 30d --metric charge --limit 20
```

## 步骤 5：下钻到创意级别

对消耗最高的组，下钻查看其下的创意：

```bash
# 查看该组下的创意列表
adex ks creatives --tenant 6 --unit <UNIT_ID> --format table

# 或直接看创意级 Top-N
adex ks creatives top --tenant 6 --range 30d --metric charge --limit 10
```

## 步骤 6：查看明细日报表

对关注的计划/组/创意，查看其日报表明细：

```bash
# 计划日报表
adex ks campaign-reports daily --tenant 6 --range 30d --campaign <CAMPAIGN_ID> --format table

# 组日报表
adex ks unit-reports daily --tenant 6 --range 30d --unit <UNIT_ID> --format table

# 创意日报表
adex ks creative-reports daily --tenant 6 --range 30d --creative <CREATIVE_ID> --format table
```

## 步骤 7：汇总对比

按维度分组汇总，对比各计划/组/创意的总消耗：

```bash
# 计划汇总排名
adex ks campaign-reports summary --tenant 6 --range 30d --group-by campaign_id --order-by charge --order-desc

# 组汇总排名
adex ks unit-reports summary --tenant 6 --range 30d --group-by unit_id --order-by charge --order-desc

# 创意汇总排名
adex ks creative-reports summary --tenant 6 --range 30d --group-by creative_id --order-by charge --order-desc
```

## 变体：按转化指标分析

如果用户关注转化而非消耗，将 `--metric charge` 替换为转化指标：

```bash
# 先查可用指标
adex ks report-metric-meta --level campaign --enabled 1 --jq '.items[].field'

# 用转化指标排名
adex ks campaigns top --tenant 6 --range 30d --metric conversion_num --limit 10
adex ks campaign-reports summary --tenant 6 --range 30d --group-by campaign_id --order-by conversion_num --order-desc
```

## 变体：按特定广告主分析

如果用户只关注某个广告主的数据，全程加 `--advertiser`：

```bash
adex ks dashboard --tenant 6 --range 30d
adex ks campaigns top --tenant 6 --range 30d --metric charge --advertiser 1234567890 --limit 10
adex ks campaign-reports summary --tenant 6 --range 30d --group-by campaign_id --advertiser 1234567890
```

## 完整示例

```bash
# 1. 大盘
adex ks dashboard --tenant 6 --range 30d

# 2. 账户趋势
adex ks account-reports daily --tenant 6 --range 30d --format table

# 3. 计划 Top 10
adex ks campaigns top --tenant 6 --range 30d --metric charge --limit 10

# 4. 下钻 Top 计划的组
adex ks units --tenant 6 --campaign 9899931248 --format table

# 5. 下钻 Top 组的创意
adex ks creatives --tenant 6 --unit 29638466721 --format table

# 6. Top 计划日报表
adex ks campaign-reports daily --tenant 6 --range 30d --campaign 9899931248 --format table

# 7. 汇总对比
adex ks campaign-reports summary --tenant 6 --range 30d --group-by campaign_id --order-by charge --order-desc
```

## 参考

- [adex-ks](../SKILL.md) — 快手广告全部命令
- [adex-ks-dashboard](adex-ks-dashboard.md) — 租户级概览
- [adex-ks-campaigns](adex-ks-campaigns.md) — 计划命令
- [adex-ks-units](adex-ks-units.md) — 组命令
- [adex-ks-creatives](adex-ks-creatives.md) — 创意命令
- [adex-ks-reports](adex-ks-reports.md) — 日报表和汇总报表
- [adex-ks-metric-meta](adex-ks-metric-meta.md) — 指标元数据
- [adex-shared](../../adex-shared/SKILL.md) — 认证和全局参数
