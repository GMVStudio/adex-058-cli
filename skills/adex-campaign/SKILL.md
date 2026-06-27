---
name: adex-campaign
version: 0.1.0
description: "Query ADEX campaign daily reports with filtering, sorting, and pagination. Use when the user wants to check advertising campaign performance metrics."
---

# adex campaign 技能

查询广告投放活动日报表数据。

## 命令

```bash
adex raw campaign daily [flags]
```

## 必需参数

| Flag | Type | Description |
|------|------|-------------|
| `--tenant` | int | 租户 ID（正整数） |
| `--range` | string | 时间范围：`1d`（1天）、`7d`（7天）、`1h`（1小时）、`30m`（30分钟） |

## 可选参数

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--campaign` | string | — | 活动 ID（数字）或名称模式 |
| `--page` | int | 1 | 页码 |
| `--page-size` | int | 20 | 每页条数 |
| `--order-by` | string | charge | 排序字段 |
| `--order-desc` | bool | true | 是否降序排列 |
| `--stat-hour` | int | -1 | 统计小时（-1 表示最新） |

## 全局参数

| Flag | Description |
|------|-------------|
| `--format` | 输出格式：`json`（默认）、`pretty`、`table` |
| `--dry-run` | 打印请求但不实际执行 |
| `--base-url` | API 地址（覆盖 `ADEX_API_BASE_URL` 环境变量） |

## 示例

```bash
# 查询最近1天的活动日报
adex raw campaign daily --tenant 6 --range 1d

# 查询指定活动
adex raw campaign daily --tenant 6 --campaign C-618-001-619 --range 1d

# 查询最近7天，表格输出
adex raw campaign daily --tenant 6 --range 7d --format table

# 预览请求（不执行）
adex raw campaign daily --tenant 6 --range 1d --dry-run

# 自定义 API 地址
adex raw campaign daily --tenant 6 --range 1d --base-url http://api.example.com
```

## 注意事项

- `--tenant` 必须是正整数
- `--range` 只接受 `1d`、`7d`、`1h`、`30m` 四种值
- 使用 `--dry-run` 可以在不调用 API 的情况下验证参数是否正确
- 表格模式会在 stderr 输出汇总信息（total、page、page_size）
