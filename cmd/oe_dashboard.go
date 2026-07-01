package cmd

import (
	"github.com/spf13/cobra"
)

// newOeDashboardCmd creates "adex oe dashboard" (GetOeDashboard).
func newOeDashboardCmd(f *Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dashboard",
		Short: "Tenant-level Oceanengine overview",
		Long: `Tenant-level overview: account count + balance breakdown + range summary
metrics + account rankings (GET /v1/oe/dashboard).

Examples:
  adex oe dashboard --tenant 6 --range 30d
  adex oe dashboard --tenant 6 --begin 2026-06-01 --end 2026-06-30`,
		RunE: func(cmd *cobra.Command, args []string) error {
			tenant, err := requireTenant(cmd)
			if err != nil {
				return err
			}
			params := map[string]interface{}{"tenant_id": tenant}

			if err := applyDateRange(cmd, params, true); err != nil {
				return err
			}

			return f.runSingle(cmd, "/v1/oe/dashboard", params, nil)
		},
	}

	addTenantFlag(cmd)
	addDateRangeFlags(cmd)
	addJQFlag(cmd)

	return cmd
}
