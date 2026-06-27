#!/usr/bin/env bash
set -euo pipefail

# ADEX CLI release script
# Usage:
#   ./scripts/release.sh 0.2.0           # 自动化发布（tag + push → GitHub Actions）
#   ./scripts/release.sh 0.2.0 --skip-test  # 跳过测试
#   ./scripts/release.sh 0.2.0 --manual    # 手动模式（本地执行全流程）
#
# Prerequisites:
#   - git, go, goreleaser installed
#   - For --manual: GITHUB_TOKEN and npm login required

set +u
VERSION="${1:-}"
SKIP_TEST=false
MANUAL=false

for arg in "$@"; do
  case "$arg" in
    --skip-test) SKIP_TEST=true ;;
    --manual)    MANUAL=true ;;
  esac
done
set -u

if [ -z "$VERSION" ]; then
  echo "Usage: ./scripts/release.sh <version> [--skip-test] [--manual]"
  echo "  version: e.g. 0.2.0 (without v prefix)"
  echo ""
  echo "Options:"
  echo "  --skip-test   Skip make test before release"
  echo "  --manual      Run full release locally (no GitHub Actions)"
  exit 1
fi

TAG="v${VERSION}"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

# Verify version matches package.json
PKG_VERSION=$(node -p "require('./package.json').version")
if [ "$PKG_VERSION" != "$VERSION" ]; then
  echo "ERROR: package.json version is \"$PKG_VERSION\", expected \"$VERSION\""
  echo "       Update package.json version first:"
  echo "       jq '.version = \"$VERSION\"' package.json > tmp && mv tmp package.json"
  exit 1
fi

# Check working tree is clean
if [ -n "$(git status --porcelain)" ]; then
  echo "ERROR: working tree is not clean"
  echo "       Commit or stash changes first:"
  echo "       git add -A && git commit -m \"chore: release ${TAG}\""
  exit 1
fi

# Check tag doesn't already exist
if git rev-parse "$TAG" >/dev/null 2>&1; then
  echo "ERROR: tag ${TAG} already exists"
  exit 1
fi

echo "=========================================="
echo "  ADEX CLI Release: ${TAG}"
echo "  Mode: $( $MANUAL && echo 'manual (local)' || echo 'auto (GitHub Actions)' )"
echo "=========================================="

if [ "$SKIP_TEST" = false ]; then
  echo ""
  echo "[1/4] Running tests..."
  make test
  echo "  ✓ Tests passed"
else
  echo ""
  echo "[1/4] Skipping tests (--skip-test)"
fi

echo ""
echo "[2/4] Checking prerequisites..."

if $MANUAL; then
  # Check goreleaser
  if ! command -v goreleaser >/dev/null 2>&1; then
    echo "ERROR: goreleaser not installed"
    echo "       Install: go install github.com/goreleaser/goreleaser/v2@latest"
    exit 1
  fi

  # Check GITHUB_TOKEN
  if [ -z "${GITHUB_TOKEN:-}" ]; then
    echo "ERROR: GITHUB_TOKEN not set"
    echo "       export GITHUB_TOKEN=ghp_your_token"
    exit 1
  fi

  # Check npm auth
  if ! npm whoami >/dev/null 2>&1; then
    echo "ERROR: not logged in to npm"
    echo "       Run: npm login"
    exit 1
  fi

  echo "  ✓ goreleaser: $(goreleaser --version | head -1)"
  echo "  ✓ GITHUB_TOKEN: set"
  echo "  ✓ npm: logged in as $(npm whoami)"
else
  echo "  ✓ Will be handled by GitHub Actions"
fi

echo ""
echo "[3/4] Creating tag ${TAG}..."
git tag "$TAG"
echo "  ✓ Tag created"

if $MANUAL; then
  echo ""
  echo "[4/4] Manual release..."

  echo "  → Running goreleaser..."
  goreleaser release --clean
  echo "  ✓ goreleaser completed"

  echo "  → Copying checksums.txt..."
  cp dist/checksums.txt .
  echo "  ✓ checksums.txt copied"

  echo "  → Publishing to npm..."
  npm publish --access public
  echo "  ✓ npm publish completed"

  echo ""
  echo "  Verifying..."
  npm view @gmvstudio/adex-cli version
else
  echo ""
  echo "[4/4] Pushing tag to trigger GitHub Actions..."
  git push origin "$TAG"
  echo "  ✓ Tag pushed"
  echo ""
  echo "  GitHub Actions will now:"
  echo "    1. goreleaser: build + publish GitHub Release"
  echo "    2. publish-npm: download checksums + npm publish"
  echo ""
  echo "  Monitor: https://github.com/GMVStudio/adex-058-cli/actions"
  echo ""
  echo "  Verify after completion:"
  echo "    npm view @gmvstudio/adex-cli version"
fi

echo ""
echo "=========================================="
echo "  Release ${TAG} initiated successfully!"
echo "=========================================="
