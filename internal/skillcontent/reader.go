package skillcontent

import (
	"io/fs"
	"path"
	"sort"
	"strings"

	"github.com/gmvstudio/adex-cli/errs"
	"gopkg.in/yaml.v3"
)

type Reader struct {
	fsys fs.FS
}

func New(fsys fs.FS) *Reader { return &Reader{fsys: fsys} }

type SkillInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version,omitempty"`
}

type DirEntry struct {
	Path  string `json:"path"`
	IsDir bool   `json:"is_dir"`
}

func (r *Reader) List() ([]SkillInfo, error) {
	entries, err := fs.ReadDir(r.fsys, ".")
	if err != nil {
		return nil, errs.NewInternalError(errs.SubtypeFileIO, "failed to read embedded skills: %v", err)
	}
	out := make([]SkillInfo, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if info, ok := r.skillInfo(e.Name()); ok {
			out = append(out, info)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

func (r *Reader) skillInfo(name string) (SkillInfo, bool) {
	data, err := fs.ReadFile(r.fsys, name+"/SKILL.md")
	if err != nil {
		return SkillInfo{}, false
	}
	desc, version := parseFrontmatter(data)
	return SkillInfo{Name: name, Description: desc, Version: version}, true
}

func (r *Reader) ListPath(arg string) ([]DirEntry, string, error) {
	name, sub := SplitArg(arg)
	if err := r.ensureSkill(name); err != nil {
		return nil, "", err
	}
	dir := name
	if sub != "" {
		cleaned, err := cleanSubPath(sub)
		if err != nil {
			return nil, "", err
		}
		dir = name + "/" + cleaned
		info, err := fs.Stat(r.fsys, dir)
		if err != nil {
			return nil, "", errs.NewValidationError(errs.SubtypeInvalidArgument,
				"path %q not found in skill %q", sub, name).
				WithHint("run 'adex skills list " + name + "' to see files in this skill")
		}
		if !info.IsDir() {
			return nil, "", errs.NewValidationError(errs.SubtypeInvalidArgument,
				"path %q is a file, not a directory; use 'adex skills read %s/%s' to read it", sub, name, cleaned)
		}
	}
	entries, err := fs.ReadDir(r.fsys, dir)
	if err != nil {
		return nil, "", errs.NewInternalError(errs.SubtypeFileIO,
			"failed to read embedded skill content: %v", err)
	}
	out := make([]DirEntry, 0, len(entries))
	for _, e := range entries {
		out = append(out, DirEntry{Path: dir + "/" + e.Name(), IsDir: e.IsDir()})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out, dir, nil
}

func SplitArg(arg string) (name, rest string) {
	name, rest, _ = strings.Cut(arg, "/")
	return name, rest
}

func parseFrontmatter(skillMD []byte) (description, version string) {
	lines := strings.Split(string(skillMD), "\n")
	if strings.TrimRight(lines[0], "\r") != "---" {
		return "", ""
	}
	block := make([]string, 0, len(lines))
	closed := false
	for _, ln := range lines[1:] {
		if strings.TrimRight(ln, "\r") == "---" {
			closed = true
			break
		}
		block = append(block, ln)
	}
	if !closed {
		return "", ""
	}
	var fm struct {
		Description string `yaml:"description"`
		Version     string `yaml:"version"`
	}
	if err := yaml.Unmarshal([]byte(strings.Join(block, "\n")), &fm); err != nil {
		return "", ""
	}
	return fm.Description, fm.Version
}

func (r *Reader) ReadSkill(name string) ([]byte, error) {
	if err := r.ensureSkill(name); err != nil {
		return nil, err
	}
	data, err := fs.ReadFile(r.fsys, name+"/SKILL.md")
	if err != nil {
		return nil, errs.NewInternalError(errs.SubtypeFileIO,
			"failed to read embedded skill content: %v", err)
	}
	return data, nil
}

func (r *Reader) ensureSkill(name string) error {
	if name == "" || strings.ContainsAny(name, `/\`) || name == "." || name == ".." {
		return unknownSkill(name)
	}
	info, err := fs.Stat(r.fsys, name)
	if err != nil || !info.IsDir() {
		return unknownSkill(name)
	}
	return nil
}

func unknownSkill(name string) error {
	return errs.NewValidationError(errs.SubtypeInvalidArgument, "unknown skill %q", name).
		WithHint("run 'adex skills list' to see available skills")
}

func cleanSubPath(relpath string) (string, error) {
	cleaned := path.Clean(relpath)
	if relpath == "" || path.IsAbs(relpath) || cleaned == "." ||
		cleaned == ".." || strings.HasPrefix(cleaned, "../") || strings.HasPrefix(cleaned, `..\`) {
		return "", errs.NewValidationError(errs.SubtypeInvalidArgument,
			"invalid path %q: must be a relative path without '..'", relpath)
	}
	return cleaned, nil
}

func (r *Reader) ReadReference(name, relpath string) ([]byte, string, error) {
	if err := r.ensureSkill(name); err != nil {
		return nil, "", err
	}
	cleaned, err := cleanSubPath(relpath)
	if err != nil {
		return nil, "", err
	}
	full := name + "/" + cleaned
	info, err := fs.Stat(r.fsys, full)
	if err != nil {
		return nil, "", errs.NewValidationError(errs.SubtypeInvalidArgument,
			"reference %q not found in skill %q", relpath, name).
			WithHint("run 'adex skills list " + name + "' to see files in this skill")
	}
	if info.IsDir() {
		return nil, "", errs.NewValidationError(errs.SubtypeInvalidArgument,
			"reference %q is a directory, not a file", relpath)
	}
	data, err := fs.ReadFile(r.fsys, full)
	if err != nil {
		return nil, "", errs.NewInternalError(errs.SubtypeFileIO,
			"failed to read embedded skill content: %v", err)
	}
	return data, cleaned, nil
}
