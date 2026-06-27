# ADEX CLI 发布指南

## 发布架构

```
开发者 git push --tags v*
           │
           ▼
    ┌──────────────────────────────────┐
    │     GitHub Actions: release.yml  │
    └──────────────────────────────────┘
           │
           ├─ Job 1: goreleaser
           │    ├─ checkout (fetch-depth=0)
           │    ├─ setup Go 1.23
           │    ├─ goreleaser release --clean
           │    │    ├─ 交叉编译 6 个平台二进制
           │    │    │  darwin/amd64, darwin/arm64
           │    │    │  linux/amd64,  linux/arm64
           │    │    │  windows/amd64, windows/arm64
           │    │    ├─ 打包归档 (tar.gz / zip)
           │    │    ├─ 生成 checksums.txt (SHA-256)
           │    │    └─ 上传到 GitHub Release
           │    └─ 完成
           │
           └─ Job 2: publish-npm (needs: goreleaser)
                ├─ checkout
                ├─ setup Node 20 + registry
                ├─ gh release download checksums.txt
                │    从 GitHub Release 下载校验和文件
                └─ npm publish --access public
                     使用 NPM_TOKEN (Granular Access Token)
                     发布 @gmvstudio/adex-cli 到 npm registry
```

### 关键组件

| 组件 | 文件 | 作用 |
|------|------|------|
| GitHub Actions | `.github/workflows/release.yml` | 自动化发布流水线 |
| goreleaser | `.goreleaser.yml` | Go 跨平台交叉编译 + GitHub Release |
| npm 包定义 | `package.json` | npm 包元数据、版本号、发布配置 |
| 二进制下载 | `scripts/install.js` | postinstall 时从 GitHub Release 下载二进制 |
| npm 入口 | `scripts/run.js` | bin 包装器，委托给原生二进制 |
| 安装向导 | `scripts/install-wizard.js` | 交互式安装引导 |
| 校验和 | `checksums.txt` | SHA-256 校验，确保二进制完整性 |

### GitHub Secrets

| Secret 名称 | 用途 | 创建方式 |
|-------------|------|---------|
| `NPM_TOKEN` | npm 发布认证 | npm Granular Access Token（Read and write） |
| `GITHUB_TOKEN` | GitHub Release 操作 | Actions 自动提供，无需手动配置 |

## 发布命令

### 自动化发布（推荐）

```bash
# 1. 确保所有测试通过
make test

# 2. 更新 package.json 版本号
#    版本号必须与即将打的 git tag 一致（不含 v 前缀）
#    例如 v0.2.0 → package.json version "0.2.0"

# 3. 更新 CHANGELOG.md（如有）

# 4. 提交变更
git add -A
git commit -m "chore: release v0.2.0"

# 5. 打 tag 并推送（触发 GitHub Actions）
git tag v0.2.0
git push origin main --tags

# 6. 等待 GitHub Actions 完成
#    https://github.com/GMVStudio/adex-058-cli/actions
#    Job 1 goreleaser → Job 2 publish-npm

# 7. 验证发布结果
npm view @gmvstudio/adex-cli version
# 应输出: 0.2.0
```

### 手动发布（备用）

当 GitHub Actions 不可用时，可手动执行完整流程：

```bash
# === 步骤 1: goreleaser 构建并发布 GitHub Release ===

# 确保工作区干净
git status

# 打 tag
git tag v0.2.0
git push origin v0.2.0

# 设置代理（如需要）
export https_proxy=http://127.0.0.1:7890 http_proxy=http://127.0.0.1:7890 all_proxy=http://127.0.0.1:7890

# 设置 GitHub Token
export GITHUB_TOKEN=ghp_your_token

# goreleaser 构建 + 发布
goreleaser release --clean

# 复制 checksums.txt 到项目根目录
cp dist/checksums.txt .

# === 步骤 2: npm 发布 ===

# 登录 npm（如未登录）
npm login

# 发布（scoped 包必须指定 --access public）
npm publish --access public

# 验证
npm view @gmvstudio/adex-cli version
```

### 快捷发布脚本

使用项目根目录的 `scripts/release.sh`：

```bash
# 自动化发布
./scripts/release.sh 0.2.0

# 跳过测试
./scripts/release.sh 0.2.0 --skip-test

# 手动模式（不推送 tag，不触发 Actions，本地执行全流程）
./scripts/release.sh 0.2.0 --manual
```

## 首次配置

### npm Granular Access Token

1. 打开 https://www.npmjs.com/settings/zhijunsoh/tokens/granular-access-tokens/new
2. **Token name**: `adex-cli-publish`
3. **Expiration**: 按需选择（建议 90 天）
4. **Packages and scopes**: Select packages → `@gmvstudio/adex-cli`
5. **Permissions**: Read and write
6. 生成后复制 token

### GitHub Secret 配置

1. 打开 https://github.com/GMVStudio/adex-058-cli/settings/secrets/actions
2. 点击 **New repository secret**
3. **Name**: `NPM_TOKEN`
4. **Secret**: 粘贴上一步生成的 npm token
5. 点击 **Add secret**

## 注意事项

### ⚠️ 版本号一致性

git tag、`package.json` version、goreleaser 产出三者必须完全一致。

```
git tag v0.2.0
package.json: "version": "0.2.0"
goreleaser 产出: adex-0.2.0-darwin-amd64.tar.gz
```

不一致会导致 npm publish 失败或用户安装到错误版本。

### ⚠️ npm registry 必须使用 HTTPS

npm 默认 registry 可能被配置为 `http://`，会导致 `E426 Upgrade Required` 错误。

```bash
# 检查
npm config get registry

# 修复（必须使用 https）
npm config set registry https://registry.npmjs.org/
```

### ⚠️ scoped 包必须指定 --access public

`@gmvstudio/adex-cli` 是 scoped 包，默认发布为私有包（需要付费账户）。
必须在 `npm publish` 时加 `--access public`，或在 `package.json` 中配置：

```json
"publishConfig": {
  "access": "public"
}
```

否则会报 `E402 Payment Required`。

### ⚠️ checksums.txt 不应提交到仓库

`checksums.txt` 是 goreleaser 构建产物，应加入 `.gitignore`，不要提交到 git。
`publish-npm` job 会从 GitHub Release 下载最新的 `checksums.txt`。
如果仓库中存在该文件，`gh release download` 会因文件已存在而失败。

### ⚠️ repository.url 必须与 GitHub 仓库匹配

`package.json` 中的 `repository.url` 必须与实际 GitHub 仓库 URL 完全匹配（包括大小写）。
不匹配会导致 npm Trusted Publishing（OIDC）返回 404，以及 npm 网站上的仓库链接错误。

```json
// 正确
"repository": {
  "url": "git+https://github.com/GMVStudio/adex-058-cli.git"
}

// 错误（仓库名不匹配）
"repository": {
  "url": "git+https://github.com/gmvstudio/adex-cli.git"
}
```

### ⚠️ scripts/install.js 中的 REPO 必须与实际仓库匹配

`scripts/install.js` 中的 `REPO` 常量用于从 GitHub Releases 下载二进制文件，必须与实际仓库匹配。

```javascript
// 正确
const REPO = "GMVStudio/adex-058-cli";

// 错误
const REPO = "gmvstudio/adex-cli";
```

### ⚠️ goreleaser 归档需要 LICENSE 文件

`.goreleaser.yml` 的 `archives.files` 中列出了 `LICENSE`，如果仓库中不存在该文件，
goreleaser 会在归档阶段报 `file does not exist` 错误。确保仓库根目录有 `LICENSE` 文件。

### ⚠️ goreleaser archives 配置格式

goreleaser v2 中 `format_overrides.format` 已废弃，应使用 `formats`（复数）：

```yaml
# 正确 (v2)
archives:
  - formats:
      - tar.gz
    format_overrides:
      - goos: windows
        formats:
          - zip

# 错误 (v1 风格，已废弃)
archives:
  - format_overrides:
      - goos: windows
        format: zip
```

### ⚠️ dist/ 目录不应提交到 git

goreleaser 的构建产物在 `dist/` 目录，不应提交到 git。确保 `.gitignore` 中包含 `dist/`。

### ⚠️ npm Token 类型选择

- **Classic Token**: 受账户 2FA 设置约束，无法在 CI 中绕过 2FA
- **Granular Access Token**: 可独立配置权限，支持 CI/CD 自动化发布

CI/CD 场景必须使用 **Granular Access Token**，不要使用 Classic Token。

### ⚠️ npm Trusted Publishing（OIDC）限制

npm Trusted Publishing 使用 GitHub Actions OIDC 进行无 token 发布，但：
- 包必须已存在于 npm registry（首次发布无法使用）
- 需要在 npm 网站手动配置 Trusted Publisher
- `repository.url` 必须与 GitHub 仓库完全匹配
- 需要 `id-token: write` 权限

**推荐方案**：首次发布使用 Granular Access Token，后续可选择切换到 Trusted Publishing。

### ⚠️ 安全最佳实践

- 不要在聊天、日志、代码中暴露 npm token 或 GitHub token
- npm Granular Access Token 设置合理的过期时间
- 定期轮换 token
- 使用 GitHub Secrets 存储 token，不要硬编码在 workflow 中
- 仓库设置为 private 时，npm provenance 不会生成
