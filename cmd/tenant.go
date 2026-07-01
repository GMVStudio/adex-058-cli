package cmd

import (
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

	return cmd
}
