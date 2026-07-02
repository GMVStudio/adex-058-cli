package cmd

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/gmvstudio/adex-cli/errs"
	"github.com/gmvstudio/adex-cli/internal/build"
	"github.com/gmvstudio/adex-cli/internal/output"
	"github.com/gmvstudio/adex-cli/internal/selfupdate"
	"github.com/gmvstudio/adex-cli/internal/skillscheck"
	"github.com/gmvstudio/adex-cli/internal/update"
	"github.com/spf13/cobra"
)

const (
	repoURL      = "https://github.com/GMVStudio/adex-058-cli"
	maxNpmOutput = 2000
	osWindows    = "windows"
)

var (
	fetchLatest    = update.FetchLatest
	currentVersion = func() string { return build.Version }
	currentOS      = runtime.GOOS
	newUpdater     = func() *selfupdate.Updater { return selfupdate.New() }
	syncSkills     = func(opts skillscheck.SyncOptions) *skillscheck.SyncResult { return skillscheck.SyncSkills(opts) }
)

func isWindows() bool { return currentOS == osWindows }

func normalizeVersion(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "v")
	return strings.TrimPrefix(s, "V")
}

func releaseURL(version string) string {
	return repoURL + "/releases/tag/v" + strings.TrimPrefix(version, "v")
}

func changelogURL() string { return repoURL + "/blob/main/CHANGELOG.md" }

func symOK() string {
	if isWindows() {
		return "[OK]"
	}
	return "✓"
}

func symFail() string {
	if isWindows() {
		return "[FAIL]"
	}
	return "✗"
}

func symWarn() string {
	if isWindows() {
		return "[WARN]"
	}
	return "⚠"
}

func symArrow() string {
	if isWindows() {
		return "->"
	}
	return "→"
}

func newUpdateCmd(f *Factory) *cobra.Command {
	var jsonOut bool
	var force bool
	var check bool

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update adex CLI to the latest version",
		Long: `Update adex CLI to the latest version.

Detects the installation method automatically:
  - npm install: runs npm install -g @gmvstudio/adex-cli@<version>
  - manual/other: shows GitHub Releases download URL

Use --json for structured output (for AI agents and scripts).
Use --check to only check for updates without installing.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return updateRun(f, jsonOut, force, check)
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "structured JSON output")
	cmd.Flags().BoolVar(&force, "force", false, "force reinstall even if already up to date")
	cmd.Flags().BoolVar(&check, "check", false, "only check for updates, do not install")
	return cmd
}

func updateRun(f *Factory, jsonOut, force, check bool) error {
	cur := currentVersion()
	updater := newUpdater()

	if !check {
		updater.CleanupStaleFiles()
	}
	output.PendingNotice = nil

	latest, err := fetchLatest()
	if err != nil {
		return errs.NewNetworkError(errs.SubtypeNetworkTransport,
			"failed to check latest version: %s", err).WithCause(err)
	}

	if update.ParseVersion(latest) == nil {
		return errs.NewInternalError(errs.SubtypeInvalidResponse,
			"invalid version from registry: %s", latest)
	}

	if !force && !update.IsNewer(latest, cur) {
		var skillsResult *skillscheck.SyncResult
		if !check {
			skillsResult = runSkillsAndState(updater, f, cur, force)
		}
		return reportAlreadyUpToDate(f, jsonOut, cur, latest, skillsResult, check)
	}

	detect := updater.DetectInstallMethod()

	if check {
		return reportCheckResult(f, jsonOut, cur, latest, detect.CanAutoUpdate())
	}

	if !detect.CanAutoUpdate() {
		return doManualUpdate(f, jsonOut, cur, latest, detect, updater)
	}
	return doNpmUpdate(f, jsonOut, cur, latest, updater)
}

func doManualUpdate(f *Factory, jsonOut bool, cur, latest string, detect selfupdate.DetectResult, updater *selfupdate.Updater) error {
	skillsResult := runSkillsAndState(updater, f, cur, true)

	reason := detect.ManualReason()
	if jsonOut {
		out := map[string]interface{}{
			"ok": true, "previous_version": cur, "latest_version": latest,
			"action":  "manual_required",
			"message": fmt.Sprintf("Automatic update unavailable: %s (path: %s)", reason, detect.ResolvedPath),
			"url":     releaseURL(latest), "changelog": changelogURL(),
		}
		applySkillsResult(out, skillsResult)
		printJSON(f.Out, out)
		return nil
	}
	fmt.Fprintf(f.ErrOut, "Automatic update unavailable: %s (path: %s).\n\n", reason, detect.ResolvedPath)
	fmt.Fprintf(f.ErrOut, "To update manually, download the latest release:\n")
	fmt.Fprintf(f.ErrOut, "  Release:   %s\n", releaseURL(latest))
	fmt.Fprintf(f.ErrOut, "  Changelog: %s\n", changelogURL())
	fmt.Fprintf(f.ErrOut, "\nOr install via npm:\n  npm install -g %s@%s\n  npx skills add https://adex-skills.oss-cn-hangzhou.aliyuncs.com -y -g   # sync skills separately\n", selfupdate.NpmPackage, latest)
	emitSkillsTextHints(f, skillsResult)
	return nil
}

func doNpmUpdate(f *Factory, jsonOut bool, cur, latest string, updater *selfupdate.Updater) error {
	restore, err := updater.PrepareSelfReplace()
	if err != nil {
		return errs.NewInternalError(errs.SubtypeUnknown,
			"failed to prepare update: %s", err).WithCause(err)
	}

	if !jsonOut {
		fmt.Fprintf(f.ErrOut, "Updating adex %s %s %s via npm ...\n", cur, symArrow(), latest)
	}

	npmResult := updater.RunNpmInstall(latest)
	if npmResult.Err != nil {
		restore()
		combined := npmResult.CombinedOutput()
		if jsonOut {
			printJSON(f.Out, map[string]interface{}{
				"ok": false, "error": map[string]interface{}{
					"type": "update_error", "message": fmt.Sprintf("npm install failed: %s", npmResult.Err),
					"detail": selfupdate.Truncate(combined, maxNpmOutput),
					"hint":   permissionHint(combined),
				},
			})
			return errs.NewInternalError(errs.SubtypeUnknown, "npm install failed")
		}
		if npmResult.Stdout.Len() > 0 {
			fmt.Fprint(f.ErrOut, npmResult.Stdout.String())
		}
		if npmResult.Stderr.Len() > 0 {
			fmt.Fprint(f.ErrOut, npmResult.Stderr.String())
		}
		fmt.Fprintf(f.ErrOut, "\n%s Update failed: %s\n", symFail(), npmResult.Err)
		if hint := permissionHint(combined); hint != "" {
			fmt.Fprintf(f.ErrOut, "  %s\n", hint)
		}
		return errs.NewInternalError(errs.SubtypeUnknown, "npm install failed: %w", npmResult.Err)
	}

	if err := updater.VerifyBinary(latest); err != nil {
		restore()
		msg := fmt.Sprintf("new binary verification failed: %s", err)
		if jsonOut {
			printJSON(f.Out, map[string]interface{}{
				"ok":    false,
				"error": map[string]interface{}{"type": "update_error", "message": msg},
			})
			return errs.NewInternalError(errs.SubtypeUnknown, msg)
		}
		fmt.Fprintf(f.ErrOut, "\n%s %s\n", symFail(), msg)
		return errs.NewInternalError(errs.SubtypeUnknown, msg)
	}

	skillsResult := runSkillsAndState(updater, f, latest, false)

	if jsonOut {
		result := map[string]interface{}{
			"ok": true, "previous_version": cur, "current_version": latest,
			"latest_version": latest, "action": "updated",
			"message": fmt.Sprintf("adex updated from %s to %s", cur, latest),
			"url":     releaseURL(latest), "changelog": changelogURL(),
		}
		applySkillsResult(result, skillsResult)
		printJSON(f.Out, result)
		return nil
	}

	fmt.Fprintf(f.ErrOut, "\n%s Successfully updated adex from %s to %s\n", symOK(), cur, latest)
	fmt.Fprintf(f.ErrOut, "  Changelog: %s\n", changelogURL())
	if skillsResult != nil {
		fmt.Fprintf(f.ErrOut, "\nUpdating skills ...\n")
	}
	emitSkillsTextHints(f, skillsResult)
	return nil
}

func runSkillsAndState(updater *selfupdate.Updater, f *Factory, stateVersion string, force bool) *skillscheck.SyncResult {
	if !force {
		if existing, ok := skillscheck.ReadSyncedVersion(); ok && normalizeVersion(existing) == normalizeVersion(stateVersion) {
			return nil
		}
	}
	result := syncSkills(skillscheck.SyncOptions{
		Version: stateVersion,
		Force:   force,
		Runner:  updater,
	})
	if result.Err != nil && strings.Contains(result.Err.Error(), "state not written") {
		fmt.Fprintf(f.ErrOut, "warning: %v\n", result.Err)
	}
	return result
}

func reportAlreadyUpToDate(f *Factory, jsonOut bool, cur, latest string, skillsResult *skillscheck.SyncResult, check bool) error {
	if jsonOut {
		out := map[string]interface{}{
			"ok": true, "previous_version": cur, "current_version": cur,
			"latest_version": latest, "action": "already_up_to_date",
			"message": fmt.Sprintf("adex %s is already up to date", cur),
		}
		if check {
			applySkillsStatus(out, cur)
		} else {
			applySkillsResult(out, skillsResult)
		}
		printJSON(f.Out, out)
		return nil
	}
	fmt.Fprintf(f.ErrOut, "%s adex %s is already up to date\n", symOK(), cur)
	if !check {
		emitSkillsTextHints(f, skillsResult)
	}
	return nil
}

func reportCheckResult(f *Factory, jsonOut bool, cur, latest string, canAutoUpdate bool) error {
	if jsonOut {
		out := map[string]interface{}{
			"ok": true, "previous_version": cur, "current_version": cur,
			"latest_version": latest, "action": "update_available",
			"auto_update": canAutoUpdate,
			"message":     fmt.Sprintf("adex %s %s %s available", cur, symArrow(), latest),
			"url":         releaseURL(latest), "changelog": changelogURL(),
		}
		applySkillsStatus(out, cur)
		printJSON(f.Out, out)
		return nil
	}
	fmt.Fprintf(f.ErrOut, "Update available: %s %s %s\n", cur, symArrow(), latest)
	fmt.Fprintf(f.ErrOut, "  Release:   %s\n", releaseURL(latest))
	fmt.Fprintf(f.ErrOut, "  Changelog: %s\n", changelogURL())
	if canAutoUpdate {
		fmt.Fprintf(f.ErrOut, "\nRun `adex update` to install.\n")
	} else {
		fmt.Fprintf(f.ErrOut, "\nDownload the release above to update manually.\n")
	}
	return nil
}

func permissionHint(npmOutput string) string {
	if strings.Contains(npmOutput, "EACCES") && !isWindows() {
		return "Permission denied. Try: sudo adex update, or adjust your npm global prefix: https://docs.npmjs.com/resolving-eacces-permissions-errors"
	}
	return ""
}

func applySkillsStatus(env map[string]interface{}, target string) {
	state, readable, err := skillscheck.ReadState()
	if err != nil || !readable || state.Version == "" {
		return
	}
	status := map[string]interface{}{
		"current": state.Version,
		"target":  target,
		"in_sync": normalizeVersion(state.Version) == normalizeVersion(target),
	}
	if len(state.OfficialSkills) > 0 {
		status["official"] = len(state.OfficialSkills)
	}
	if len(state.UpdatedSkills) > 0 {
		status["updated"] = len(state.UpdatedSkills)
	}
	if len(state.SkippedDeletedSkills) > 0 {
		status["skipped_deleted"] = state.SkippedDeletedSkills
	}
	env["skills_status"] = status
}

func applySkillsResult(env map[string]interface{}, r *skillscheck.SyncResult) {
	switch {
	case r == nil:
		env["skills_action"] = "in_sync"
	case r.Err != nil:
		env["skills_action"] = "failed"
		env["skills_warning"] = fmt.Sprintf("skills update failed: %s", r.Err)
		env["skills_summary"] = skillsSummary(r)
	default:
		env["skills_action"] = "synced"
		env["skills_summary"] = skillsSummary(r)
	}
}

func skillsSummary(r *skillscheck.SyncResult) map[string]interface{} {
	summary := map[string]interface{}{
		"official":        len(r.Official),
		"updated":         len(r.Updated),
		"added":           len(r.Added),
		"skipped_deleted": len(r.SkippedDeleted),
	}
	if len(r.Failed) > 0 {
		summary["failed"] = r.Failed
	}
	return summary
}

func emitSkillsTextHints(f *Factory, r *skillscheck.SyncResult) {
	switch {
	case r == nil:
	case r.Err != nil:
		fmt.Fprintf(f.ErrOut, "%s Skills update failed: %v\n", symWarn(), r.Err)
		if len(r.Failed) > 0 {
			fmt.Fprintf(f.ErrOut, "  Failed skills: %s\n", strings.Join(r.Failed, ", "))
		}
		fmt.Fprintf(f.ErrOut, "  To retry all official skills: adex update --force\n")
	case r.Force:
		fmt.Fprintf(f.ErrOut, "%s Skills updated: restored all %d official skills\n", symOK(), len(r.Official))
	default:
		fmt.Fprintf(f.ErrOut, "%s Skills updated: %d official, %d updated, %d added, %d skipped because deleted locally\n", symOK(), len(r.Official), len(r.Updated), len(r.Added), len(r.SkippedDeleted))
		if len(r.SkippedDeleted) > 0 {
			fmt.Fprintf(f.ErrOut, "  To restore all official skills: adex update --force\n")
		}
	}
}
