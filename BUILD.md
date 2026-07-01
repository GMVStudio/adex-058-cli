# ADEX CLI 编译发布指南

## 环境要求

| 工具 | 版本 | 用途 |
|------|------|------|
| Go | v1.23+ | 编译 CLI 二进制 |
| Node.js | v16+ | npm 包发布 |
| goreleaser | v2+ | 跨平台交叉编译 + GitHub Release 发布 |
| git | — | 版本管理、tag |

## Go 编译

### 本地快速编译

```bash
# 编译到当前目录
go build -o adex .

# 验证
./adex --help
```

### 通过 Makefile 编译（推荐）

```bash
# 完整流程：vet + fmt-check + test
make test

# 仅编译（带版本信息注入）
make build

# 安装到 /usr/local/bin
make install

# 安装到自定义路径
make install PREFIX=$HOME/.local

# 卸载
make uninstall

# 清理构建产物
make clean
```

Makefile 会通过 `git describe --tags --always --dirty` 自动提取版本号，注入到二进制的 `internal/build.Version` 和 `internal/build.Date` 变量中。

### 代码质量检查

```bash
# go vet
make vet

# gofmt 检查（不修改文件，仅报告）
make fmt-check

# 自动格式化
gofmt -w .

# 完整测试（vet + fmt-check + go test）
make test
```

## 跨平台发布（goreleaser）

### 配置文件

`.goreleaser.yml` 定义了跨平台构建规则：

| OS | Arch | 归档格式 |
|----|------|---------|
| darwin | amd64, arm64 | tar.gz |
| linux | amd64, arm64 | tar.gz |
| windows | amd64, arm64 | zip |

归档命名格式：`adex-{version}-{os}-{arch}.tar.gz`（Windows 为 `.zip`）

### 发布步骤

```bash
# 1. 确保工作区干净
git status

# 2. 打 tag（使用语义化版本）
git tag v0.1.0
git push origin v0.1.0

# 3. goreleaser 会由 GitHub Actions 自动触发
#    或手动本地执行：
goreleaser release --clean

# 4. 仅测试构建（不发布，不推送）
goreleaser release --snapshot --clean
```

发布后 GitHub Releases 会包含：
- `adex-{version}-darwin-amd64.tar.gz`
- `adex-{version}-darwin-arm64.tar.gz`
- `adex-{version}-linux-amd64.tar.gz`
- `adex-{version}-linux-arm64.tar.gz`
- `adex-{version}-windows-amd64.zip`
- `adex-{version}-windows-arm64.zip`
- `checksums.txt`（SHA-256 校验和）

## npm 包发布

### 包结构

npm 包 `@gmvstudio/adex-cli` 不包含 Go 源码，而是通过 postinstall 钩子从 GitHub Releases 下载预编译的二进制文件。

```
package.json 中 files 字段定义了发布内容：
├── scripts/install.js          # postinstall: 下载二进制 + 校验
├── scripts/install-wizard.js   # 交互式安装向导
├── scripts/run.js              # bin 入口包装器
├── checksums.txt               # SHA-256 校验和（与 goreleaser 产出一致）
└── CHANGELOG.md
```

### 发布流程

```bash
# 1. 确保 goreleaser 已发布对应版本的 Release
#    npm 包版本必须与 GitHub Release tag 一致

# 2. 更新 package.json 版本号
#    版本号必须与 git tag 一致（不含 v 前缀）
#    例如 git tag v0.1.0 → package.json version "0.1.0"

# 3. 更新 checksums.txt
#    从 GitHub Release 下载 checksums.txt，替换本地文件
#    或 goreleaser 产出后直接复制

# 4. 登录 npm（首次需要）
npm login

# 5. 发布
make npm-publish
# 或
npm publish

# 6. 验证
npm view @gmvstudio/adex-cli version
```

### npm install 工作原理

```
用户执行 npm install -g @gmvstudio/adex-cli
  │
  ├─ npm 下载包文件到 node_modules
  │
  ├─ 触发 postinstall → scripts/install.js
  │    ├─ 检测平台 (os/arch)
  │    ├─ 从 GitHub Releases 下载对应归档
  │    │   GitHub → npmmirror 镜像回退
  │    ├─ SHA-256 校验
  │    └─ 解压到 bin/adex
  │
  └─ 注册 bin/adex → 全局 adex 命令
```

### npx 临时执行

```bash
# npx 会跳过 postinstall 二进制下载
# 首次执行命令时 run.js 会自动触发下载
npx @gmvstudio/adex-cli --help

# 交互式安装向导
npx @gmvstudio/adex-cli install
```

## 完整发布流程（Go + npm 联合发布）

```bash
# 1. 确保所有测试通过
make test

# 2. 更新 CHANGELOG.md

# 3. 更新 package.json 版本号
#    编辑 package.json，使 version 与即将打的 tag 一致

# 4. 提交变更
git add -A
git commit -m "chore: release v0.1.0"

# 5. 打 tag 并推送
git tag v0.1.0
git push origin main --tags

# 6. goreleaser 自动构建并发布 GitHub Release
#    （或手动：goreleaser release --clean）

# 7. 下载 checksums.txt 到项目根目录
#    从 GitHub Release 页面下载 checksums.txt
#    替换本地 checksums.txt

# 8. 发布到 npm
npm publish

# 9. 验证安装
npm install -g @gmvstudio/adex-cli
adex --help
```

## 注意事项

- **版本号一致性**：git tag、`package.json` version、goreleaser 产出三者必须一致
- **checksums.txt**：npm 包中的 `checksums.txt` 必须与 GitHub Release 中的文件内容一致，否则 postinstall 校验会失败
- **镜像回退**：`scripts/install.js` 依次尝试 GitHub → npmmirror，确保中国大陆用户也能正常下载
- **npx 兼容**：npx 环境下 postinstall 会被跳过，`run.js` 在首次执行时自动补下载
- **CGO_ENABLED=0**：goreleaser 构建使用纯 Go，无 C 依赖，确保交叉编译兼容性