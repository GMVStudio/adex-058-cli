package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"

	"github.com/gmvstudio/adex-cli/cmd"
)

//go:embed skills/*/SKILL.md skills/*/references
var skillsEmbedFS embed.FS

func init() {
	sub, err := fs.Sub(skillsEmbedFS, "skills")
	if err != nil {
		fmt.Fprintln(os.Stderr, "warning: skills embed assembly failed, skills commands disabled:", err)
		return
	}
	cmd.SetEmbeddedSkillContent(sub)
}
