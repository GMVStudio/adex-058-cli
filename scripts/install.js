const fs = require("fs");
const path = require("path");
const { execFileSync } = require("child_process");
const os = require("os");
const crypto = require("crypto");

const VERSION = require("../package.json").version.replace(/-.*$/, "");
const REPO = "GMVStudio/adex-058-cli";
const NAME = "adex";
const DEFAULT_MIRROR_HOST = "https://registry.npmmirror.com";
const ALLOWED_HOSTS = new Set([
  "github.com",
  "objects.githubusercontent.com",
  "registry.npmmirror.com",
]);

const PLATFORM_MAP = {
  darwin: "darwin",
  linux: "linux",
  win32: "windows",
};

const ARCH_MAP = {
  x64: "amd64",
  arm64: "arm64",
};

const platform = PLATFORM_MAP[process.platform];
const arch = ARCH_MAP[process.arch];

const isWindows = process.platform === "win32";
const ext = isWindows ? ".zip" : ".tar.gz";
const archiveName = `${NAME}-${VERSION}-${platform}-${arch}${ext}`;
const GITHUB_URL = `https://github.com/${REPO}/releases/download/v${VERSION}/${archiveName}`;

const binDir = path.join(__dirname, "..", "bin");
const dest = path.join(binDir, NAME + (isWindows ? ".exe" : ""));

function resolveMirrorUrls(env, archive, version) {
  const binaryPath = `/-/binary/adex/v${version}/${archive}`;
  const defaultUrl = joinUrl(DEFAULT_MIRROR_HOST, binaryPath);

  const urls = [];
  const registry = (env.npm_config_registry || "").trim();
  if (registry && !isDefaultNpmjsRegistry(registry) && isValidDownloadBase(registry)) {
    const base = new URL(registry);
    urls.push(joinUrl(base.origin + base.pathname, binaryPath));
  }
  if (!urls.includes(defaultUrl)) urls.push(defaultUrl);
  return urls;
}

function joinUrl(base, suffix) {
  return base.replace(/\/+$/, "") + suffix;
}

function isValidDownloadBase(raw) {
  try {
    const parsed = new URL(raw);
    return parsed.protocol === "https:" && !!parsed.hostname;
  } catch (_) {
    return false;
  }
}

function isDefaultNpmjsRegistry(url) {
  try {
    const { hostname } = new URL(url);
    return hostname === "registry.npmjs.org";
  } catch (_) {
    return false;
  }
}

function assertAllowedHost(url) {
  const { hostname } = new URL(url);
  if (!ALLOWED_HOSTS.has(hostname)) {
    throw new Error(`Download host not allowed: ${hostname}`);
  }
}

function getMirrorUrls(env) {
  const urls = resolveMirrorUrls(env, archiveName, VERSION);
  for (const u of urls) ALLOWED_HOSTS.add(new URL(u).hostname);
  return urls;
}

function isCurlVersionSupported(versionOutput) {
  const match = String(versionOutput).match(/^\s*curl\s+(\d+)\.(\d+)\.(\d+)/i);
  if (!match) return false;
  const major = parseInt(match[1], 10);
  const minor = parseInt(match[2], 10);
  return major > 7 || (major === 7 && minor >= 70);
}

let _curlSupportsSslRevokeBestEffort;

function curlSupportsSslRevokeBestEffort() {
  if (_curlSupportsSslRevokeBestEffort !== undefined) {
    return _curlSupportsSslRevokeBestEffort;
  }
  try {
    const output = execFileSync("curl", ["--version"], {
      stdio: ["ignore", "pipe", "ignore"],
      encoding: "utf8",
      timeout: 5000,
    });
    _curlSupportsSslRevokeBestEffort = isCurlVersionSupported(output);
  } catch (_) {
    _curlSupportsSslRevokeBestEffort = false;
  }
  return _curlSupportsSslRevokeBestEffort;
}

function download(url, destPath) {
  assertAllowedHost(url);
  const args = [
    "--fail", "--location", "--silent", "--show-error",
    "--connect-timeout", "10", "--max-time", "120",
    "--max-redirs", "3",
    "--output", destPath,
  ];
  if (isWindows && curlSupportsSslRevokeBestEffort()) {
    args.unshift("--ssl-revoke-best-effort");
  }
  args.push(url);
  execFileSync("curl", args, { stdio: ["ignore", "ignore", "pipe"] });
}

function extractZipWindows(archivePath, destDir) {
  const psOpts = ["-NoProfile", "-ExecutionPolicy", "Bypass", "-Command"];
  const psStdio = ["ignore", "inherit", "inherit"];
  const psEnv = {
    ...process.env,
    ADEX_CLI_ARCHIVE: archivePath,
    ADEX_CLI_DEST: destDir,
  };

  try {
    const dotnet =
      "$ErrorActionPreference='Stop';" +
      "Add-Type -AssemblyName System.IO.Compression.FileSystem;" +
      "[System.IO.Compression.ZipFile]::ExtractToDirectory($env:ADEX_CLI_ARCHIVE,$env:ADEX_CLI_DEST)";
    execFileSync("powershell.exe", [...psOpts, dotnet], { stdio: psStdio, env: psEnv });
  } catch (primaryErr) {
    try {
      const cmdlet =
        "$ErrorActionPreference='Stop';" +
        "Expand-Archive -LiteralPath $env:ADEX_CLI_ARCHIVE -DestinationPath $env:ADEX_CLI_DEST -Force";
      execFileSync("powershell.exe", [...psOpts, cmdlet], { stdio: psStdio, env: psEnv });
    } catch (secondErr) {
      try {
        execFileSync("tar", ["-xf", archivePath, "-C", destDir], { stdio: psStdio });
      } catch (fallbackErr) {
        throw new Error(
          `Failed to extract ${archivePath}. ` +
          `.NET ZipFile attempt: ${primaryErr.message}. ` +
          `Expand-Archive fallback: ${secondErr.message}. ` +
          `tar fallback: ${fallbackErr.message}`
        );
      }
    }
  }
}

function install() {
  const mirrorUrls = getMirrorUrls(process.env);
  const downloadUrls = [GITHUB_URL, ...mirrorUrls];

  fs.mkdirSync(binDir, { recursive: true });

  const tmpDir = fs.mkdtempSync(path.join(os.tmpdir(), "adex-cli-"));
  const archivePath = path.join(tmpDir, archiveName);

  try {
    let lastErr;
    let downloaded = false;
    const failedUrls = [];
    for (const url of downloadUrls) {
      try {
        download(url, archivePath);
        downloaded = true;
        break;
      } catch (e) {
        failedUrls.push({ url, error: e.message });
        lastErr = e;
      }
    }
    if (!downloaded) {
      const details = failedUrls.map(f => `  ${f.url}\n    → ${f.error}`).join("\n");
      throw new Error(`All download sources failed:\n${details}`);
    }

    const expectedHash = getExpectedChecksum(archiveName);
    verifyChecksum(archivePath, expectedHash);

    if (isWindows) {
      extractZipWindows(archivePath, tmpDir);
    } else {
      execFileSync("tar", ["-xzf", archivePath, "-C", tmpDir], {
        stdio: "ignore",
      });
    }

    const binaryName = NAME + (isWindows ? ".exe" : "");
    const extractedBinary = path.join(tmpDir, binaryName);

    fs.copyFileSync(extractedBinary, dest);
    fs.chmodSync(dest, 0o755);
    console.log(`${NAME} v${VERSION} installed successfully`);
  } finally {
    fs.rmSync(tmpDir, { recursive: true, force: true });
  }
}

function getExpectedChecksum(archiveName, checksumsDir) {
  const dir = checksumsDir || path.join(__dirname, "..");
  const checksumsPath = path.join(dir, "checksums.txt");

  if (!fs.existsSync(checksumsPath)) {
    console.error(
      "[WARN] checksums.txt not found, skipping checksum verification"
    );
    return null;
  }

  const content = fs.readFileSync(checksumsPath, "utf8");
  for (const line of content.split("\n")) {
    const trimmed = line.trim();
    if (!trimmed) continue;
    const idx = trimmed.indexOf("  ");
    if (idx === -1) continue;
    const hash = trimmed.slice(0, idx);
    const name = trimmed.slice(idx + 2);
    if (name === archiveName) return hash;
  }

  throw new Error(`Checksum entry not found for ${archiveName}`);
}

function verifyChecksum(archivePath, expectedHash) {
  if (expectedHash === null) return;

  const hash = crypto.createHash("sha256");
  const fd = fs.openSync(archivePath, "r");
  try {
    const buf = Buffer.alloc(64 * 1024);
    let bytesRead;
    while ((bytesRead = fs.readSync(fd, buf, 0, buf.length, null)) > 0) {
      hash.update(buf.subarray(0, bytesRead));
    }
  } finally {
    fs.closeSync(fd);
  }
  const actual = hash.digest("hex");

  if (actual.toLowerCase() !== expectedHash.toLowerCase()) {
    throw new Error(
      `[SECURITY] Checksum mismatch for ${path.basename(archivePath)}: expected ${expectedHash} but got ${actual}`
    );
  }
}

if (require.main === module) {
  if (!platform || !arch) {
    console.error(
      `Unsupported platform: ${process.platform}-${process.arch}`
    );
    process.exit(1);
  }

  const isNpxPostinstall =
    process.env.npm_command === "exec" && !process.env.ADEX_CLI_RUN;

  if (isNpxPostinstall) {
    process.exit(0);
  }

  try {
    install();
  } catch (err) {
    console.error(`Failed to install ${NAME} binary:`, err.message);
    console.error(
      `\nThe binary will be auto-downloaded on first run.\n` +
      `If you prefer to install manually:\n` +
      `  # 1. Use a proxy:\n` +
      `  export https_proxy=http://your-proxy:port\n` +
      `  npm install -g @gmvstudio/adex-cli\n\n` +
      `  # 2. Point to a corporate npm mirror that proxies /-/binary/adex/...:\n` +
      `  npm install -g @gmvstudio/adex-cli --registry=https://your-corp-mirror/`
    );
    process.exit(0);
  }
}

module.exports = { getExpectedChecksum, verifyChecksum, assertAllowedHost, resolveMirrorUrls, curlSupportsSslRevokeBestEffort, isCurlVersionSupported };
