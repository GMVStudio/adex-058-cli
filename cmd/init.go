package cmd

import (
	"github.com/gmvstudio/adex-cli/errs"
	"github.com/gmvstudio/adex-cli/internal/config"
	"github.com/spf13/cobra"
)

// newInitCmd creates the "init" command that binds the API key (one-time setup).
func newInitCmd(f *Factory) *cobra.Command {
	var authorization string
	var tenantID int

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Bind API credentials (one-time setup)",
		Long: `Bind your ADEX API key so subsequent commands are authenticated.
Optionally set a default tenant so --tenant is not needed on every command.

The credential is written to ` + config.Path() + ` with 0600 permissions.
You may pass the raw key or the full "Bearer <key>" header value.

Examples:
  adex init --authorization "Bearer adex_c93462599a6246a89f55a11b024b1a1a"
  adex init --authorization adex_c93462599a6246a89f55a11b024b1a1a --tenant 6
  adex init --authorization adex_c93462599a6246a89f55a11b024b1a1a --base-url http://47.99.131.55:8000`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if authorization == "" {
				return errs.NewValidationError(errs.SubtypeInvalidArgument,
					"--authorization is required").
					WithParam("--authorization").
					WithHint(`pass your API key, e.g. --authorization "Bearer adex_xxx"`)
			}

			token := config.NormalizeToken(authorization)
			if token == "" {
				return errs.NewValidationError(errs.SubtypeInvalidArgument,
					"--authorization must contain a non-empty API key").
					WithParam("--authorization")
			}

			f.Config.Authorization = token
			if baseURL, _ := cmd.Flags().GetString("base-url"); baseURL != "" {
				f.Config.BaseURL = baseURL
			}
			if tenantID > 0 {
				f.Config.TenantID = tenantID
			}

			if err := config.Save(f.Config); err != nil {
				// config.Save already returns a typed *errs.* error.
				return err
			}

			result := map[string]interface{}{
				"ok":       true,
				"message":  "credentials saved",
				"path":     config.Path(),
				"base_url": f.Config.BaseURL,
			}
			if f.Config.TenantID > 0 {
				result["tenant_id"] = f.Config.TenantID
			} else {
				result["next_step"] = "run 'adex tenant --page-all' to list tenants, then 'adex tenant use <ID>' to set a default tenant"
			}
			printJSON(f.Out, result)
			return nil
		},
	}

	cmd.Flags().StringVar(&authorization, "authorization", "", "API key or 'Bearer <key>' header value (required)")
	cmd.Flags().IntVar(&tenantID, "tenant", 0, "default tenant ID (optional; saves typing --tenant on every command)")

	return cmd
}
