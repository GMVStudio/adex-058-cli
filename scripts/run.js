#!/usr/bin/env node

const { execFileSync } = require("child_process");
const fs = require("fs");
const path = require("path");

const ext = process.platform === "win32" ? ".exe" : "";
const bin = path.join(__dirname, "..", "bin", "adex" + ext);

// Intercept "install" subcommand — run the setup wizard directly,
// bypassing the native binary (which may not exist yet under npx).
const args = process.argv.slice(2);
if (args[0] === "install") {
  require("./install-wizard.js");
} else {
  // Auto-download binary if missing (e.g. npx skipped postinstall).
  if (!fs.existsSync(bin)) {
    try {
      execFileSync(process.execPath, [path.join(__dirname, "install.js")], {
        stdio: "inherit",
        env: { ...process.env, ADEX_CLI_RUN: "true" },
      });
    } catch (_) {
      console.error(
        `\nFailed to auto-install adex binary.\n` +
        `To fix, run the install script manually:\n` +
        `  node "${path.join(__dirname, "install.js")}"\n`
      );
      process.exit(1);
    }
  }

  try {
    execFileSync(bin, args, { stdio: "inherit" });
  } catch (e) {
    process.exit(e.status || 1);
  }
}
