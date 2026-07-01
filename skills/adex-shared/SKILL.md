---
name: adex-shared
version: 0.1.0
description: "Use when first setting up adex CLI, configuring API endpoint, or troubleshooting connection issues."
---

# adex CLI 共享规则

本技能指导你如何通过 adex CLI 查询广告投放数据。

## 安装

### 通过 npm 安装（推荐）

```bash
# 安装 CLI
npm install -g @gmvstudio/adex-cli

# 安装 CLI SKILL（必需）
npx -y skills add https://open.feishu.cn --skill -y
```

### 从源码安装

```bash
git clone https://github.com/gmvstudio/adex-cli.git
cd adex-cli
make install
```

## 配置

adex CLI 通过环境变量 `ADEX_API_BASE_URL` 配置 API 端点，默认值为 `http://47.99.131.55:8000`。

```bash
# 设置 API 端点（可选，默认使用内置地址）
export ADEX_API_BASE_URL=http://your-api-host:8000
```

也可以在每条命令中通过 `--base-url` 标志覆盖：

```bash
adex ks dashboard --tenant 6 --base-url http://your-api-host:8000
```

## 验证

```bash
# 使用 --help 验证 CLI 是否正常工作
adex --help
```

## 输出格式

支持三种输出格式，通过 `--format` 标志控制：

- `json`（默认）：紧凑 JSON
- `pretty`：格式化 JSON
- `table`：表格输出

```bash
adex ks dashboard --tenant 6 --format table
```

## 错误处理

错误以 JSON 信封格式输出到 stderr，包含 `type`、`subtype`、`message` 字段：

```json
{"ok":false,"error":{"type":"validation","subtype":"invalid_argument","message":"--tenant must be a positive integer"}}
```

| Exit Code | Category | Description |
|-----------|----------|-------------|
| 0 | — | Success |
| 2 | validation | Invalid input arguments |
| 3 | unauthorized | Authentication failure |
| 4 | network | Network/transport error |
| 5 | api | API returned non-2xx |
| 1 | internal | Unexpected internal error |
