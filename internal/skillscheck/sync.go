package skillscheck

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/gmvstudio/adex-cli/internal/selfupdate"
)

var (
	skillNamePattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_:-]*(@[^\s]+)?$`)
	ansiPattern      = regexp.MustCompile(`\x1b\[[0-?]*[ -/]*[@-~]`)
)

type SyncInput struct {
	Version        string
	OfficialSkills []string
	LocalSkills    []string
	PreviousState  *SkillsState
	StateReadable  bool
	Force          bool
}

type SyncPlan struct {
	Version        string
	OfficialSkills []string
	ToUpdate       []string
	Added          []string
	SkippedDeleted []string
}

func stripANSI(s string) string {
	return ansiPattern.ReplaceAllString(s, "")
}

func ParseSkillsList(text string) []string {
	text = stripANSI(text)
	lines := strings.Split(text, "\n")

	hasGlobalSkills := strings.Contains(text, "Global Skills")
	hasAvailableSkills := strings.Contains(text, "Available Skills")

	if hasGlobalSkills {
		return parseGlobalSkillsList(lines)
	} else if hasAvailableSkills {
		return parseOfficialSkillsList(lines)
	}
	return nil
}

func ParseGlobalSkillsJSON(text string) []string {
	type globalSkill struct {
		Name string `json:"name"`
	}

	var skills []globalSkill
	if err := json.Unmarshal([]byte(text), &skills); err != nil {
		return nil
	}

	seen := map[string]bool{}
	for _, skill := range skills {
		candidate := strings.TrimSpace(skill.Name)
		if candidate == "" || !skillNamePattern.MatchString(candidate) {
			continue
		}
		seen[candidate] = true
	}

	return sortedKeys(seen)
}

func ParseOfficialSkillsIndexJSON(text string) ([]string, error) {
	type officialSkill struct {
		Name string `json:"name"`
	}
	type officialIndex struct {
		Skills []officialSkill `json:"skills"`
	}

	var index officialIndex
	if err := json.Unmarshal([]byte(text), &index); err != nil {
		return nil, err
	}

	seen := map[string]bool{}
	for _, skill := range index.Skills {
		candidate := strings.TrimSpace(skill.Name)
		if skillNamePattern.MatchString(candidate) {
			seen[candidate] = true
		}
	}

	return sortedKeys(seen), nil
}

func parseGlobalSkillsList(lines []string) []string {
	seen := map[string]bool{}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "Global Skills") {
			continue
		}
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "Tip:") {
			continue
		}
		if strings.HasPrefix(trimmed, "Agents:") {
			continue
		}
		if isGlobalSkillsSectionHeader(trimmed) {
			continue
		}

		parts := strings.Fields(trimmed)
		if len(parts) == 0 {
			continue
		}

		candidate := parts[0]
		if candidate == "" || !skillNamePattern.MatchString(candidate) {
			continue
		}
		seen[candidate] = true
	}

	return sortedKeys(seen)
}

func isGlobalSkillsSectionHeader(line string) bool {
	switch line {
	case "General", "Project", "Local":
		return true
	default:
		return false
	}
}

func parseOfficialSkillsList(lines []string) []string {
	seen := map[string]bool{}
	inAvailableSection := false

	for _, line := range lines {
		if strings.Contains(line, "Available Skills") {
			inAvailableSection = true
			continue
		}

		if !inAvailableSection {
			continue
		}

		if strings.Contains(line, "│") {
			parts := strings.FieldsFunc(line, func(r rune) bool {
				return r == '│' || r == ' '
			})

			if len(parts) > 0 {
				candidate := parts[0]
				if skillNamePattern.MatchString(candidate) {
					seen[candidate] = true
				}
			}
		}
	}

	return sortedKeys(seen)
}

func PlanSync(input SyncInput) SyncPlan {
	official := uniqueSorted(input.OfficialSkills)
	if input.Force {
		return SyncPlan{
			Version:        input.Version,
			OfficialSkills: official,
			ToUpdate:       official,
			Added:          []string{},
			SkippedDeleted: []string{},
		}
	}

	officialSet := toSet(official)
	installedOfficial := intersection(input.LocalSkills, officialSet)

	previousOfficial := []string{}
	if input.StateReadable && input.PreviousState != nil {
		previousOfficial = input.PreviousState.OfficialSkills
	}
	previousSet := toSet(previousOfficial)

	newAddedOfficial := []string{}
	for _, skill := range official {
		if !previousSet[skill] {
			newAddedOfficial = append(newAddedOfficial, skill)
		}
	}

	updateSet := toSet(installedOfficial)
	for _, skill := range newAddedOfficial {
		updateSet[skill] = true
	}
	toUpdate := sortedKeys(updateSet)
	updateSet = toSet(toUpdate)

	skipped := []string{}
	for _, skill := range official {
		if !updateSet[skill] {
			skipped = append(skipped, skill)
		}
	}

	return SyncPlan{
		Version:        input.Version,
		OfficialSkills: official,
		ToUpdate:       toUpdate,
		Added:          uniqueSorted(newAddedOfficial),
		SkippedDeleted: skipped,
	}
}

type SkillsRunner interface {
	ListOfficialSkillsIndex() *selfupdate.NpmResult
	ListOfficialSkills() *selfupdate.NpmResult
	ListGlobalSkillsJSON() *selfupdate.NpmResult
	ListGlobalSkills() *selfupdate.NpmResult
	InstallSkill(nameList []string) *selfupdate.NpmResult
	InstallAllSkills() *selfupdate.NpmResult
}

type SyncOptions struct {
	Version string
	Force   bool
	Runner  SkillsRunner
	Now     func() time.Time
}

type SyncResult struct {
	Action         string
	Official       []string
	Updated        []string
	Added          []string
	SkippedDeleted []string
	Failed         []string
	Err            error
	Detail         string
	Force          bool
}

func SyncSkills(opts SyncOptions) *SyncResult {
	if opts.Now == nil {
		opts.Now = time.Now
	}
	if opts.Runner == nil {
		return &SyncResult{Action: "failed", Err: fmt.Errorf("skills runner is nil")}
	}

	official, reason, ok := listOfficialSkills(opts.Runner)
	if !ok {
		return fallbackFullInstall(opts, reason, nil)
	}

	local, ok := listLocalSkills(opts.Runner)
	if !ok {
		return fallbackFullInstall(opts, "local skills list failed or parsed as empty", official)
	}

	previous, readable, err := ReadState()
	if err != nil {
		readable = false
		previous = nil
	}

	plan := PlanSync(SyncInput{
		Version:        opts.Version,
		OfficialSkills: official,
		LocalSkills:    local,
		PreviousState:  previous,
		StateReadable:  readable,
		Force:          opts.Force,
	})

	result := &SyncResult{
		Action:         "synced",
		Official:       plan.OfficialSkills,
		Updated:        plan.ToUpdate,
		Added:          plan.Added,
		SkippedDeleted: plan.SkippedDeleted,
		Force:          opts.Force,
	}

	if len(plan.ToUpdate) == 0 {
		return fallbackFullInstall(opts, "toUpdate skills empty fallback", official)
	}

	if len(plan.ToUpdate) > 0 {
		installResult := opts.Runner.InstallSkill(plan.ToUpdate)
		if installResult == nil || installResult.Err != nil {
			return fallbackFullInstall(opts, resultDetail(installResult), official)
		}
	}

	state := SkillsState{
		Version:              opts.Version,
		OfficialSkills:       plan.OfficialSkills,
		UpdatedSkills:        plan.ToUpdate,
		AddedOfficialSkills:  plan.Added,
		SkippedDeletedSkills: plan.SkippedDeleted,
		UpdatedAt:            opts.Now().UTC().Format(time.RFC3339),
	}
	if err := WriteState(state); err != nil {
		result.Action = "failed"
		result.Err = fmt.Errorf("skills synced but state not written: %w", err)
		return result
	}

	return result
}

func listOfficialSkills(runner SkillsRunner) ([]string, string, bool) {
	reasons := []string{}

	indexResult := runner.ListOfficialSkillsIndex()
	if indexResult == nil || indexResult.Err != nil {
		reasons = append(reasons, "official skills index failed: "+resultDetail(indexResult))
	} else {
		official, err := ParseOfficialSkillsIndexJSON(indexResult.Stdout.String())
		if err != nil {
			reasons = append(reasons, "official skills index JSON invalid: "+err.Error())
		} else if len(official) > 0 {
			return official, "", true
		} else {
			reasons = append(reasons, "official skills index contains no skills")
		}
	}

	officialResult := runner.ListOfficialSkills()
	if officialResult == nil || officialResult.Err != nil {
		reasons = append(reasons, "official skills list failed: "+resultDetail(officialResult))
		return nil, strings.Join(reasons, "; "), false
	}
	official := ParseSkillsList(officialResult.Stdout.String())
	if len(official) > 0 {
		return official, "", true
	}
	if strings.TrimSpace(officialResult.Stdout.String()) != "" {
		reasons = append(reasons, "official skills list parsed as empty despite non-empty stdout")
	} else {
		reasons = append(reasons, "official skills list returned no skills")
	}
	return nil, strings.Join(reasons, "; "), false
}

func listLocalSkills(runner SkillsRunner) ([]string, bool) {
	jsonResult := runner.ListGlobalSkillsJSON()
	if jsonResult != nil && jsonResult.Err == nil {
		if local := ParseGlobalSkillsJSON(jsonResult.Stdout.String()); len(local) > 0 {
			return local, true
		}
	}

	textResult := runner.ListGlobalSkills()
	if textResult != nil && textResult.Err == nil {
		if local := ParseSkillsList(textResult.Stdout.String()); len(local) > 0 {
			return local, true
		}
	}

	return nil, false
}

func fallbackFullInstall(opts SyncOptions, reason string, official []string) *SyncResult {
	installResult := opts.Runner.InstallAllSkills()
	if installResult == nil {
		return &SyncResult{
			Action: "fallback_failed",
			Err:    fmt.Errorf("full skills install failed: empty result (reason: %s)", reason),
			Detail: reason,
			Force:  opts.Force,
		}
	}
	if installResult.Err != nil {
		return &SyncResult{
			Action: "fallback_failed",
			Err:    fmt.Errorf("full skills install failed: %w (reason: %s)", installResult.Err, reason),
			Detail: reason + "\n" + resultDetail(installResult),
			Force:  opts.Force,
		}
	}

	state := SkillsState{
		Version:              opts.Version,
		OfficialSkills:       official,
		UpdatedSkills:        official,
		AddedOfficialSkills:  official,
		SkippedDeletedSkills: []string{},
		UpdatedAt:            opts.Now().UTC().Format(time.RFC3339),
	}
	if writeErr := WriteState(state); writeErr != nil {
		return &SyncResult{
			Action:         "fallback_synced",
			Official:       official,
			Updated:        official,
			Added:          official,
			SkippedDeleted: []string{},
			Detail:         reason + "\nstate write failed: " + writeErr.Error(),
			Force:          opts.Force,
		}
	}

	return &SyncResult{
		Action:         "fallback_synced",
		Official:       official,
		Updated:        official,
		Added:          official,
		SkippedDeleted: []string{},
		Detail:         reason,
		Force:          opts.Force,
	}
}

func resultDetail(result *selfupdate.NpmResult) string {
	if result == nil {
		return ""
	}
	parts := []string{}
	if output := strings.TrimSpace(result.CombinedOutput()); output != "" {
		parts = append(parts, output)
	}
	if result.Err != nil {
		parts = append(parts, result.Err.Error())
	}
	return strings.Join(parts, "\n")
}

func uniqueSorted(values []string) []string {
	return sortedKeys(toSet(values))
}

func toSet(values []string) map[string]bool {
	out := map[string]bool{}
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			out[value] = true
		}
	}
	return out
}

func intersection(values []string, allowed map[string]bool) []string {
	out := map[string]bool{}
	for _, value := range values {
		if allowed[value] {
			out[value] = true
		}
	}
	return sortedKeys(out)
}

func sortedKeys(values map[string]bool) []string {
	out := make([]string, 0, len(values))
	for value := range values {
		out = append(out, value)
	}
	sort.Strings(out)
	return out
}
