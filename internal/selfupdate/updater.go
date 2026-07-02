// Package selfupdate handles installation detection, npm-based updates,
// skills updates, and platform-specific binary replacement for the CLI
// self-update flow.
package selfupdate

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/gmvstudio/adex-cli/internal/vfs"
)

var execLookPath = exec.LookPath

type InstallMethod int

const (
	InstallNpm InstallMethod = iota
	InstallManual
)

const (
	NpmPackage = "@gmvstudio/adex-cli"
)

const (
	npmInstallTimeout      = 10 * time.Minute
	skillsUpdateTimeout    = 2 * time.Minute
	skillsIndexMaxBodySize = 1 << 20
	verifyTimeout          = 10 * time.Second
)

var (
	skillsIndexFetchTimeout = 10 * time.Second
	officialSkillsIndexURL  = "https://adex-skills.oss-cn-hangzhou.aliyuncs.com/.well-known/skills/index.json"
)

// skillsSource is the primary URL for npx skills add.
const skillsSource = "https://adex-skills.oss-cn-hangzhou.aliyuncs.com"

// skillsSourceFallback is used when the primary source fails.
const skillsSourceFallback = "GMVStudio/adex-058-cli"

type DetectResult struct {
	Method       InstallMethod
	ResolvedPath string
	NpmAvailable bool
}

func (d DetectResult) CanAutoUpdate() bool {
	return d.Method == InstallNpm && d.NpmAvailable
}

func (d DetectResult) ManualReason() string {
	if d.Method == InstallNpm && !d.NpmAvailable {
		return "installed via npm, but npm is not available in PATH"
	}
	return "not installed via npm"
}

type NpmResult struct {
	Stdout bytes.Buffer
	Stderr bytes.Buffer
	Err    error
}

func (r *NpmResult) CombinedOutput() string {
	return r.Stdout.String() + r.Stderr.String()
}

type Updater struct {
	DetectOverride           func() DetectResult
	NpmInstallOverride       func(version string) *NpmResult
	SkillsIndexFetchOverride func() *NpmResult
	SkillsCommandOverride    func(args ...string) *NpmResult
	VerifyOverride           func(expectedVersion string) error
	RestoreAvailableOverride func() bool

	backupCreated bool
}

func New() *Updater { return &Updater{} }

func (u *Updater) DetectInstallMethod() DetectResult {
	if u.DetectOverride != nil {
		return u.DetectOverride()
	}
	exe, err := vfs.Default.Executable()
	if err != nil {
		return DetectResult{Method: InstallManual}
	}
	resolved, err := vfs.Default.EvalSymlinks(exe)
	if err != nil {
		return DetectResult{Method: InstallManual, ResolvedPath: exe}
	}

	method := InstallManual
	if strings.Contains(resolved, "node_modules") {
		method = InstallNpm
	}

	npmAvailable := false
	if method == InstallNpm {
		if _, err := exec.LookPath("npm"); err == nil {
			npmAvailable = true
		}
	}

	return DetectResult{
		Method:       method,
		ResolvedPath: resolved,
		NpmAvailable: npmAvailable,
	}
}

func (u *Updater) RunNpmInstall(version string) *NpmResult {
	if u.NpmInstallOverride != nil {
		return u.NpmInstallOverride(version)
	}
	r := &NpmResult{}
	npmPath, err := exec.LookPath("npm")
	if err != nil {
		r.Err = fmt.Errorf("npm not found in PATH: %w", err)
		return r
	}
	ctx, cancel := context.WithTimeout(context.Background(), npmInstallTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, npmPath, "install", "-g", NpmPackage+"@"+version)
	cmd.Stdout = &r.Stdout
	cmd.Stderr = &r.Stderr
	r.Err = cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		r.Err = fmt.Errorf("npm install timed out after %s", npmInstallTimeout)
	}
	return r
}

func (u *Updater) ListOfficialSkillsIndex() *NpmResult {
	if u.SkillsIndexFetchOverride != nil {
		return u.SkillsIndexFetchOverride()
	}

	r := &NpmResult{}
	ctx, cancel := context.WithTimeout(context.Background(), skillsIndexFetchTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, officialSkillsIndexURL, nil)
	if err != nil {
		r.Err = err
		return r
	}

	client := &http.Client{}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if req.URL.Scheme != "https" {
			return fmt.Errorf("official skills index redirected to non-HTTPS URL: %s", req.URL.Redacted())
		}
		return nil
	}
	resp, err := client.Do(req)
	if err != nil {
		r.Err = err
		return r
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		r.Err = fmt.Errorf("official skills index returned HTTP %d", resp.StatusCode)
		return r
	}

	limited := io.LimitReader(resp.Body, skillsIndexMaxBodySize+1)
	if _, err := io.Copy(&r.Stdout, limited); err != nil {
		r.Err = err
		return r
	}
	if r.Stdout.Len() > skillsIndexMaxBodySize {
		r.Stdout.Reset()
		r.Err = fmt.Errorf("official skills index exceeds %d bytes", skillsIndexMaxBodySize)
		return r
	}
	return r
}

func (u *Updater) ListOfficialSkills() *NpmResult {
	r := u.runSkillsListOfficial(skillsSource)
	if r.Err != nil {
		r = u.runSkillsListOfficial(skillsSourceFallback)
	}
	return r
}

func (u *Updater) ListGlobalSkills() *NpmResult {
	return u.runSkillsListGlobal()
}

func (u *Updater) ListGlobalSkillsJSON() *NpmResult {
	return u.runSkillsCommand("-y", "skills", "ls", "-g", "--json")
}

func (u *Updater) InstallSkill(nameList []string) *NpmResult {
	r := u.runSkillsInstall(skillsSource, nameList)
	if r.Err != nil {
		r = u.runSkillsInstall(skillsSourceFallback, nameList)
	}
	return r
}

func (u *Updater) InstallAllSkills() *NpmResult {
	r := u.runSkillsAdd(skillsSource)
	if r.Err != nil {
		r = u.runSkillsAdd(skillsSourceFallback)
	}
	return r
}

func (u *Updater) runSkillsAdd(source string) *NpmResult {
	return u.runSkillsCommand("-y", "skills", "add", source, "-g", "-y")
}

func (u *Updater) runSkillsListOfficial(source string) *NpmResult {
	return u.runSkillsCommand("-y", "skills", "add", source, "--list")
}

func (u *Updater) runSkillsListGlobal() *NpmResult {
	return u.runSkillsCommand("-y", "skills", "ls", "-g")
}

func (u *Updater) runSkillsInstall(source string, nameList []string) *NpmResult {
	args := []string{"-y", "skills", "add", source, "-s"}
	args = append(args, nameList...)
	args = append(args, "-g", "-y")
	return u.runSkillsCommand(args...)
}

func (u *Updater) runSkillsCommand(args ...string) *NpmResult {
	if u.SkillsCommandOverride != nil {
		return u.SkillsCommandOverride(args...)
	}
	r := &NpmResult{}
	npxPath, err := exec.LookPath("npx")
	if err != nil {
		r.Err = fmt.Errorf("npx not found in PATH: %w", err)
		return r
	}
	ctx, cancel := context.WithTimeout(context.Background(), skillsUpdateTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, npxPath, args...)
	cmd.Stdout = &r.Stdout
	cmd.Stderr = &r.Stderr
	r.Err = cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		r.Err = fmt.Errorf("skills update timed out after %s", skillsUpdateTimeout)
	}
	return r
}

func (u *Updater) VerifyBinary(expectedVersion string) error {
	if u.VerifyOverride != nil {
		return u.VerifyOverride(expectedVersion)
	}
	exe, err := execLookPath("adex")
	if err != nil {
		exe, err = vfs.Default.Executable()
		if err != nil {
			return fmt.Errorf("cannot locate binary: %w", err)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), verifyTimeout)
	defer cancel()
	out, err := exec.CommandContext(ctx, exe, "--version").Output()
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("binary verification timed out after %s", verifyTimeout)
	}
	if err != nil {
		return fmt.Errorf("binary not executable: %w", err)
	}
	fields := strings.Fields(strings.TrimSpace(string(out)))
	if len(fields) == 0 {
		return fmt.Errorf("empty version output")
	}
	actual := strings.TrimPrefix(fields[len(fields)-1], "v")
	expected := strings.TrimPrefix(expectedVersion, "v")
	if actual != expected {
		return fmt.Errorf("expected version %s, got %q", expectedVersion, actual)
	}
	return nil
}

func Truncate(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) <= maxLen {
		return s
	}
	return string(r[len(r)-maxLen:])
}

func (u *Updater) resolveExe() (string, error) {
	exe, err := vfs.Default.Executable()
	if err != nil {
		return "", err
	}
	return vfs.Default.EvalSymlinks(exe)
}
