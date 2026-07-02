#!/usr/bin/env node

const fs = require("fs");
const path = require("path");
const { execFileSync, execFile } = require("child_process");
const p = require("@clack/prompts");

const PKG = "@gmvstudio/adex-cli";
const SKILLS_REPO = "https://adex-skills.oss-cn-hangzhou.aliyuncs.com";
const SKILLS_REPO_FALLBACK = "GMVStudio/adex-058-cli";
const isWindows = process.platform === "win32";

const messages = {
  zh: {
    setup:          "正在设置 ADEX CLI...",
    step1:          "正在安装 %s...",
    step1Upgrade:   "正在升级 %s (v%s → v%s)...",
    step1Skip:      "已安装 (v%s)，跳过",
    step1Done:      "已全局安装",
    step1Upgraded:  "已升级到 v%s",
    step1Fail:      "全局安装失败。运行以下命令重试: npm install -g %s",
    step2:          "安装 AI Skills",
    step2Skip:      "已安装，跳过",
    step2Spinner:   "正在安装 Skills...",
    step2Done:      "Skills 已安装",
    step2Fail:      "Skills 安装失败。运行以下命令重试: npx skills add %s -y",
    step3:          "正在配置 API 端点...",
    step3Prompt:    "请输入 ADEX API 地址",
    step3Default:   "http://47.99.131.55:8000",
    step3Done:      "API 端点已配置",
    step3Skip:      "跳过 API 配置",
    step4:          "正在验证二进制...",
    step4Fail:      "二进制下载失败。你的网络可能无法直连 GitHub。",
    step4ProxyHint: "请设置代理后重试：\n  export https_proxy=http://your-proxy:port http_proxy=http://your-proxy:port all_proxy=http://your-proxy:port\n  adex --help",
    done:           "安装完成！\n现在可以运行: adex --help",
    cancelled:      "安装已取消",
    nonTtyHint:     "要完成配置，请在终端中运行：\n  export ADEX_API_BASE_URL=http://your-api-host:8000",
  },
  en: {
    setup:          "Setting up ADEX CLI...",
    step1:          "Installing %s globally...",
    step1Upgrade:   "Upgrading %s (v%s → v%s)...",
    step1Skip:      "Already installed (v%s). Skipped",
    step1Done:      "Installed globally",
    step1Upgraded:  "Upgraded to v%s",
    step1Fail:      "Failed to install globally. Run manually: npm install -g %s",
    step2:          "Install AI skills",
    step2Skip:      "Already installed. Skipped",
    step2Spinner:   "Installing skills...",
    step2Done:      "Skills installed",
    step2Fail:      "Failed to install skills. Run manually: npx skills add %s -y",
    step3:          "Configuring API endpoint...",
    step3Prompt:    "Enter ADEX API base URL",
    step3Default:   "http://47.99.131.55:8000",
    step3Done:      "API endpoint configured",
    step3Skip:      "Skipped API configuration",
    step4:          "Verifying binary...",
    step4Fail:      "Binary download failed. Your network may not be able to reach GitHub directly.",
    step4ProxyHint: "Please set proxy and retry:\n  export https_proxy=http://your-proxy:port http_proxy=http://your-proxy:port all_proxy=http://your-proxy:port\n  adex --help",
    done:           "You are all set!\nNow try: adex --help",
    cancelled:      "Installation cancelled",
    nonTtyHint:     "To complete setup, run:\n  export ADEX_API_BASE_URL=http://your-api-host:8000",
  },
};

function handleCancel(value, msg) {
  if (p.isCancel(value)) {
    p.cancel(msg.cancelled);
    process.exit(0);
  }
  return value;
}

function execCmd(cmd, args, opts) {
  if (isWindows) {
    return execFileSync("cmd.exe", ["/c", cmd, ...args], opts);
  }
  return execFileSync(cmd, args, opts);
}

function run(cmd, args, opts = {}) {
  execCmd(cmd, args, { stdio: "inherit", ...opts });
}

function runSilent(cmd, args, opts = {}) {
  return execCmd(cmd, args, {
    stdio: ["ignore", "pipe", "pipe"],
    ...opts,
  });
}

function runSilentAsync(cmd, args, opts = {}) {
  const actualCmd = isWindows ? "cmd.exe" : cmd;
  const actualArgs = isWindows ? ["/c", cmd, ...args] : args;
  return new Promise((resolve, reject) => {
    execFile(actualCmd, actualArgs, {
      stdio: ["ignore", "pipe", "pipe"],
      ...opts,
    }, (err, stdout) => {
      if (err) reject(err);
      else resolve(stdout);
    });
  });
}

function fmt(template, ...values) {
  let i = 0;
  return template.replace(/%s/g, () => values[i++] ?? "");
}

function whichAdex() {
  try {
    const prefix = execFileSync("npm", ["prefix", "-g"], {
      stdio: ["ignore", "pipe", "pipe"],
    }).toString().trim();
    const bin = isWindows
      ? path.join(prefix, "adex.cmd")
      : path.join(prefix, "bin", "adex");
    if (fs.existsSync(bin)) return bin;
  } catch (_) {}
  try {
    const cmd = isWindows ? "where" : "which";
    return execFileSync(cmd, ["adex"], { stdio: ["ignore", "pipe", "pipe"] })
      .toString()
      .split("\n")[0]
      .trim();
  } catch (_) {
    return null;
  }
}

function getLatestVersion() {
  try {
    const out = runSilent("npm", ["view", PKG, "version"], { timeout: 15000 });
    const ver = out.toString().trim();
    return /^\d+\.\d+\.\d+/.test(ver) ? ver : null;
  } catch (_) {
    return null;
  }
}

function semverLessThan(a, b) {
  const pa = a.replace(/-.*$/, "").split(".").map(Number);
  const pb = b.replace(/-.*$/, "").split(".").map(Number);
  for (let i = 0; i < 3; i++) {
    if ((pa[i] || 0) < (pb[i] || 0)) return true;
    if ((pa[i] || 0) > (pb[i] || 0)) return false;
  }
  return false;
}

function getGloballyInstalledVersion() {
  try {
    const out = runSilent("npm", ["list", "-g", PKG], { timeout: 15000 });
    const match = out.toString().match(/@(\d+\.\d+\.\d+[^\s]*)/);
    return match ? match[1] : "unknown";
  } catch (_) {
    return null;
  }
}

function parseLangArg() {
  const args = process.argv.slice(2);
  for (let i = 0; i < args.length; i++) {
    if (args[i] === "--lang" && args[i + 1]) {
      const val = args[i + 1].toLowerCase();
      if (val === "zh" || val === "en") return val;
    }
    if (args[i].startsWith("--lang=")) {
      const val = args[i].split("=")[1].toLowerCase();
      if (val === "zh" || val === "en") return val;
    }
  }
  return null;
}

async function stepSelectLang() {
  const fromArg = parseLangArg();
  if (fromArg) return fromArg;

  const lang = await p.select({
    message: "请选择语言 / Select language",
    options: [
      { value: "zh", label: "中文" },
      { value: "en", label: "English" },
    ],
  });
  return handleCancel(lang, messages.zh);
}

async function stepInstallGlobally(msg) {
  const installedVer = getGloballyInstalledVersion();
  const latestVer = getLatestVersion();
  const needsUpgrade = installedVer && latestVer && semverLessThan(installedVer, latestVer);

  if (installedVer && !needsUpgrade) {
    p.log.info(fmt(msg.step1Skip, installedVer));
    return false;
  }

  const s = p.spinner();
  if (needsUpgrade) {
    s.start(fmt(msg.step1Upgrade, PKG, installedVer, latestVer));
  } else {
    s.start(fmt(msg.step1, PKG));
  }
  try {
    await runSilentAsync("npm", ["install", "-g", PKG], { timeout: 120000 });
    s.stop(needsUpgrade ? fmt(msg.step1Upgraded, latestVer) : msg.step1Done);
    return needsUpgrade;
  } catch (_) {
    s.stop(fmt(msg.step1Fail, PKG));
    process.exit(1);
  }
}

async function skillsAlreadyInstalled() {
  try {
    const out = await runSilentAsync("npx", ["-y", "skills", "ls", "-g"], {
      timeout: 120000,
    });
    return /^adex-/m.test(out.toString());
  } catch (_) {
    return false;
  }
}

async function stepInstallSkills(msg) {
  const s = p.spinner();
  s.start(msg.step2Spinner);
  try {
    if (await skillsAlreadyInstalled()) {
      s.stop(msg.step2Skip);
      return;
    }
    try {
      await runSilentAsync("npx", ["-y", "skills", "add", SKILLS_REPO, "-y", "-g"], {
        timeout: 120000,
      });
    } catch (_) {
      await runSilentAsync("npx", ["-y", "skills", "add", SKILLS_REPO_FALLBACK, "-y", "-g"], {
        timeout: 120000,
      });
    }
    s.stop(msg.step2Done);
  } catch (_) {
    s.stop(fmt(msg.step2Fail, SKILLS_REPO_FALLBACK));
    process.exit(1);
  }
}

async function stepConfigApi(msg) {
  const s = p.spinner();
  s.start(msg.step3);

  const adexBin = whichAdex();
  s.stop(msg.step3);

  if (!adexBin) {
    p.log.warn(msg.step3Skip);
    return;
  }

  const baseUrl = await p.text({
    message: msg.step3Prompt,
    placeholder: msg.step3Default,
    defaultValue: msg.step3Default,
  });

  if (handleCancel(baseUrl, msg)) {
    process.env.ADEX_API_BASE_URL = baseUrl;
    p.log.success(fmt(msg.step3Done) + `: ${baseUrl}`);
  }
}

async function stepVerifyBinary(msg) {
  const s = p.spinner();
  s.start(msg.step4);
  try {
    const adexBin = whichAdex();
    if (!adexBin) {
      s.stop(msg.step4Fail);
      p.log.warn(msg.step4ProxyHint);
      return;
    }
    runSilent(adexBin, ["--help"], { timeout: 30000 });
    s.stop();
  } catch (_) {
    s.stop(msg.step4Fail);
    p.log.warn(msg.step4ProxyHint);
  }
}

async function main() {
  const isInteractive = !!process.stdin.isTTY;
  const lang = isInteractive ? await stepSelectLang() : (parseLangArg() || "en");
  const msg = messages[lang];

  if (isInteractive) {
    p.intro(msg.setup);
    await stepInstallGlobally(msg);
    await stepInstallSkills(msg);
    await stepConfigApi(msg);
    await stepVerifyBinary(msg);
    p.outro(msg.done);
  } else {
    console.log(msg.setup);
    await stepInstallGlobally(msg);
    await stepInstallSkills(msg);
    console.log(msg.nonTtyHint);
  }
}

main().catch((err) => {
  p.cancel("Unexpected error: " + (err.message || err));
  process.exit(1);
});
