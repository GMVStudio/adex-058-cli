package cmd

import (
	"strconv"

	"github.com/gmvstudio/adex-cli/errs"
	"github.com/gmvstudio/adex-cli/internal/config"
	"github.com/spf13/cobra"
)

// newTenantCmd creates "adex tenant" (ListTenants).
func newTenantCmd(f *Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tenant",
		Short: "List tenants",
		Long: `List tenants with optional name and status filters
(GET /v1/tenants).

Examples:
  adex tenant --page-size 20
  adex tenant --name acme --status active --format table
  adex tenant --page-all --jq '.items[].id'`,
		RunE: func(cmd *cobra.Command, args []string) error {
			params := map[string]interface{}{}

			name, _ := cmd.Flags().GetString("name")
			status, _ := cmd.Flags().GetString("status")
			setString(params, "name", name)
			setString(params, "status", status)

			applyPaging(cmd, params)

			return f.runList(cmd, "/v1/tenants", params, colTenants)
		},
	}

	cmd.Flags().String("name", "", "tenant name fuzzy match")
	cmd.Flags().String("status", "", "status filter (active / disabled)")
	addPagingFlags(cmd)
	addJQFlag(cmd)

	cmd.AddCommand(newTenantUseCmd(f))

	return cmd
}

// newTenantUseCmd creates "adex tenant use <id>" to set the default tenant.
func newTenantUseCmd(f *Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use <tenant-id>",
		Short: "Set default tenant",
		Long: `Set the default tenant ID so subsequent commands don't need --tenant.
The value is persisted to ` + config.Path() + `.

Examples:
  adex tenant use 6
  adex tenant use 8`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil || id <= 0 {
				return errs.NewValidationError(errs.SubtypeInvalidArgument,
					"tenant ID must be a positive integer, got %q", args[0]).
					WithParam("tenant-id")
			}

			f.Config.TenantID = id
			if err := config.Save(f.Config); err != nil {
				return err
			}

			printJSON(f.Out, map[string]interface{}{
				"ok":        true,
				"message":   "default tenant set",
				"tenant_id": id,
				"path":      config.Path(),
			})
			return nil
		},
	}

	return cmd
}
