package cmd

import (
	"github.com/spf13/cobra"
)

// newUserCmd creates "adex user" (GetCurrentUser).
func newUserCmd(f *Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Get current user info",
		Long: `Get the current authenticated user's info, resolved from the Bearer API key
(GET /v1/users/me).

Examples:
  adex user
  adex user --format table
  adex user --jq '.username'`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return f.runSingle(cmd, "/v1/users/me", nil, colUser)
		},
	}

	addJQFlag(cmd)

	return cmd
}
