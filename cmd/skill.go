package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"

	"github.com/gmvstudio/adex-cli/errs"
	"github.com/gmvstudio/adex-cli/internal/output"
	"github.com/gmvstudio/adex-cli/internal/skillcontent"
	"github.com/spf13/cobra"
)

var embeddedSkillContent fs.FS

// SetEmbeddedSkillContent registers the embedded skill tree.
func SetEmbeddedSkillContent(fsys fs.FS) { embeddedSkillContent = fsys }

func newSkillReader() (*skillcontent.Reader, error) {
	if embeddedSkillContent == nil {
		return nil, errs.NewInternalError(errs.SubtypeFileIO,
			"skill content not embedded in this build")
	}
	return skillcontent.New(embeddedSkillContent), nil
}

type skillReadEnvelope struct {
	Skill    string `json:"skill"`
	Path     string `json:"path"`
	Content  string `json:"content"`
	Guidance string `json:"guidance,omitempty"`
}

type skillListEnvelope struct {
	OK     bool                     `json:"ok"`
	Skills []skillcontent.SkillInfo `json:"skills"`
	Count  int                      `json:"count"`
}

type skillListPathEnvelope struct {
	OK      bool                    `json:"ok"`
	Path    string                  `json:"path"`
	Entries []skillcontent.DirEntry `json:"entries"`
	Count   int                     `json:"count"`
}

func newSkillCmd(f *Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skills",
		Short: "Read embedded skill content (list / read)",
		Long: "Read agent-readable skill content (SKILL.md and reference files) embedded in " +
			"the CLI binary at build time, so it stays in sync with the CLI version.",
	}
	cmd.AddCommand(newSkillListCmd(f), newSkillReadCmd(f))
	return cmd
}

func newSkillListCmd(f *Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [name[/path]]",
		Short: "List skills, or list one layer under a skill path (like ls)",
		Example: `  adex skills list                      # all skills: name, description, version
  adex skills list adex-shared          # one layer under a skill (like ls)`,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return errs.NewValidationError(errs.SubtypeInvalidArgument,
					"list takes at most 1 argument: [name[/path]]").
					WithHint("run 'adex skills list --help'")
			}
			r, err := newSkillReader()
			if err != nil {
				return err
			}
			if len(args) == 0 {
				skills, err := r.List()
				if err != nil {
					return err
				}
				printJSON(f.Out, skillListEnvelope{OK: true, Skills: skills, Count: len(skills)})
				return nil
			}
			entries, listed, err := r.ListPath(args[0])
			if err != nil {
				return err
			}
			printJSON(f.Out, skillListPathEnvelope{OK: true, Path: listed, Entries: entries, Count: len(entries)})
			return nil
		},
	}
	return cmd
}

func newSkillReadCmd(f *Factory) *cobra.Command {
	var asJSON bool
	cmd := &cobra.Command{
		Use:   "read <name>[/<path>] [path]",
		Short: "Print a skill's SKILL.md, or a file under the skill (raw markdown by default)",
		Example: `  adex skills read adex-shared                             # the skill's SKILL.md
  adex skills read adex-shared/references/example.md        # a file under the skill
  adex skills read adex-shared --json                       # JSON envelope`,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			name, relpath, err := parseSkillReadTarget(args)
			if err != nil {
				return err
			}
			r, err := newSkillReader()
			if err != nil {
				return err
			}

			var content []byte
			var pathOut string
			if relpath == "" {
				content, err = r.ReadSkill(name)
				pathOut = "SKILL.md"
			} else {
				content, pathOut, err = r.ReadReference(name, relpath)
			}
			if err != nil {
				return err
			}

			if asJSON {
				env := skillReadEnvelope{Skill: name, Path: pathOut, Content: string(content)}
				printJSON(f.Out, env)
				return nil
			}
			if _, err := f.Out.Write(content); err != nil {
				return errs.NewInternalError(errs.SubtypeFileIO, "failed to write output: %v", err)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "output as a JSON envelope instead of raw markdown")
	return cmd
}

func parseSkillReadTarget(args []string) (name, relpath string, err error) {
	switch len(args) {
	case 1:
		name, relpath = skillcontent.SplitArg(args[0])
		return name, relpath, nil
	case 2:
		return args[0], args[1], nil
	default:
		return "", "", errs.NewValidationError(errs.SubtypeInvalidArgument,
			"read requires 1 or 2 arguments: <name>[/<path>] [path]").
			WithHint("run 'adex skills read --help'")
	}
}

func printJSON(w io.Writer, v interface{}) {
	b, _ := json.Marshal(v)
	if output.GetNotice() != nil {
		var m map[string]interface{}
		if json.Unmarshal(b, &m) == nil {
			output.InjectNotice(m)
			b, _ = json.Marshal(m)
		}
	}
	fmt.Fprintln(w, string(b))
}
