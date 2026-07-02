#!/usr/bin/env bash
# scripts/oss-sync.sh — Generate index.json and upload skills/ to Aliyun OSS.
#
# Usage:
#   ./scripts/oss-sync.sh                  # upload to default bucket
#   OSS_BUCKET=other-bucket ./scripts/oss-sync.sh
#
# Prerequisites:
#   - ossutil installed (v1.7.x) — https://help.aliyun.com/zh/oss/developer-reference/install-ossutil
#   - .env file in project root with OSS_ACCESS_KEY_ID / OSS_ACCESS_KEY_SECRET
#     (or export them as environment variables)
#
# The script:
#   1. Scans skills/*/SKILL.md for frontmatter (name, description)
#   2. Walks each skill directory to collect file paths
#   3. Generates .well-known/skills/index.json
#   4. Uploads everything to oss://<bucket>/.well-known/skills/

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
SKILLS_DIR="$ROOT_DIR/skills"
ENV_FILE="$ROOT_DIR/.env"
STAGING_DIR="$(mktemp -d)"
OSS_PREFIX=".well-known/skills"

# --- 0. Load credentials from .env if present ---
if [ -f "$ENV_FILE" ]; then
  set -a
  # shellcheck disable=SC1090
  source "$ENV_FILE"
  set +a
fi

OSS_BUCKET="${OSS_BUCKET:-adex-skills}"
OSS_ENDPOINT="${OSS_ENDPOINT:-oss-cn-hangzhou.aliyuncs.com}"
OSS_AK_ID="${OSS_ACCESS_KEY_ID:-}"
OSS_AK_SECRET="${OSS_ACCESS_KEY_SECRET:-}"

if [ -z "$OSS_AK_ID" ] || [ -z "$OSS_AK_SECRET" ]; then
  echo "ERROR: OSS_ACCESS_KEY_ID and OSS_ACCESS_KEY_SECRET must be set." >&2
  echo "   Put them in .env or export as environment variables." >&2
  exit 1
fi

cleanup() {
  rm -rf "$STAGING_DIR"
}
trap cleanup EXIT

echo "==> Staging skills in: $STAGING_DIR"

# --- 1. Copy skills/ into staging/.well-known/skills/ ---
STAGING_SKILLS="$STAGING_DIR/$OSS_PREFIX"
mkdir -p "$STAGING_SKILLS"

if [ ! -d "$SKILLS_DIR" ]; then
  echo "ERROR: skills/ directory not found at $SKILLS_DIR" >&2
  exit 1
fi

cp -R "$SKILLS_DIR/." "$STAGING_SKILLS/"

# --- 2. Generate index.json ---
echo "==> Generating index.json"

INDEX_FILE="$STAGING_SKILLS/index.json"

# Use a temporary Go program or Python to generate JSON.
# Prefer Python3 (ubiquitous on dev machines).
python3 - "$STAGING_SKILLS" "$INDEX_FILE" <<'PYEOF'
import json, os, sys, re

skills_dir = sys.argv[1]
index_path = sys.argv[2]

def parse_frontmatter(content):
    """Extract name and description from YAML frontmatter."""
    if not content.startswith("---"):
        return None, None
    end = content.find("\n---", 3)
    if end == -1:
        return None, None
    fm = content[3:end]
    name = None
    desc = None
    for line in fm.split("\n"):
        m = re.match(r'^name:\s*"?([^"\n]+)"?\s*$', line)
        if m:
            name = m.group(1).strip()
        m = re.match(r'^description:\s*"?(.+?)"?\s*$', line)
        if m:
            desc = m.group(1).strip()
    return name, desc

def collect_files(base_dir, skill_name):
    """Walk skill directory and collect relative file paths."""
    files = []
    skill_root = os.path.join(base_dir, skill_name)
    for dirpath, dirnames, filenames in os.walk(skill_root):
        dirnames.sort()
        for fname in sorted(filenames):
            full = os.path.join(dirpath, fname)
            rel = os.path.relpath(full, skill_root)
            # Normalize to forward slashes
            rel = rel.replace(os.sep, "/")
            files.append(rel)
    return files

skills = []
for entry in sorted(os.listdir(skills_dir)):
    skill_path = os.path.join(skills_dir, entry, "SKILL.md")
    if not os.path.isfile(skill_path):
        continue
    with open(skill_path, "r", encoding="utf-8") as f:
        content = f.read()
    name, desc = parse_frontmatter(content)
    if not name:
        name = entry
    if not desc:
        desc = ""
    files = collect_files(skills_dir, entry)
    skills.append({
        "name": name,
        "description": desc,
        "files": files,
    })

index = {"skills": skills}
with open(index_path, "w", encoding="utf-8") as f:
    json.dump(index, f, ensure_ascii=False, indent=2)
    f.write("\n")

print(f"   Generated index.json with {len(skills)} skills:")
for s in skills:
    print(f"   - {s['name']} ({len(s['files'])} files)")
PYEOF

if [ $? -ne 0 ]; then
  echo "ERROR: Failed to generate index.json" >&2
  exit 1
fi

# --- 3. Upload to OSS ---
echo "==> Uploading to oss://$OSS_BUCKET/$OSS_PREFIX/"

# Find ossutil: check PATH, then ~/bin/ossutil
OSSUTIL=""
if command -v ossutil &>/dev/null; then
  OSSUTIL="ossutil"
elif [ -x "$HOME/bin/ossutil" ]; then
  OSSUTIL="$HOME/bin/ossutil"
else
  echo "ERROR: ossutil not found in PATH or ~/bin/ossutil" >&2
  echo "   Install: https://help.aliyun.com/zh/oss/developer-reference/install-ossutil" >&2
  exit 1
fi

# Upload the entire staging directory (ossutil v1.7.x syntax)
"$OSSUTIL" cp -r "$STAGING_SKILLS/" "oss://$OSS_BUCKET/$OSS_PREFIX/" \
  -e "$OSS_ENDPOINT" \
  -i "$OSS_AK_ID" \
  -k "$OSS_AK_SECRET" \
  --update \
  --force

echo ""
echo "==> Upload complete!"
echo "    Index URL: https://$OSS_BUCKET.$OSS_ENDPOINT/$OSS_PREFIX/index.json"
echo "    Skills base: https://$OSS_BUCKET.$OSS_ENDPOINT/$OSS_PREFIX/"
echo ""
echo "    Test with:"
echo "    curl -s https://$OSS_BUCKET.$OSS_ENDPOINT/$OSS_PREFIX/index.json | python3 -m json.tool"
echo "    npx -y skills add https://$OSS_BUCKET.$OSS_ENDPOINT -y"
