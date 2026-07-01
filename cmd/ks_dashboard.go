package cmd

import (
	"github.com/spf13/cobra"
)

// newKsDashboardCmd creates "adex ks dashboard" (GetKsDashboard).
func newKsDashboardCmd(f *Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dashboard",
		Short: "Tenant-level Kuaishou overview",
		Long: `Tenant-level overview: account stats + range summary + account rankings
(GET /v1/ks/dashboard).

Examples:
  adex ks dashboard --tenant 6 --range 30d
  adex ks dashboard --tenant 6 --begin 2026-06-01 --end 2026-06-30 --ranking-limit 10`,
		RunE: func(cmd *cobra.Command, args []string) error {
			tenant, err := requireTenant(cmd)
			if err != nil {
				return err
			}
			params := map[string]interface{}{"tenant_id": tenant}

			if err := applyDateRange(cmd, params, true); err != nil {
				return err
			}

			source, _ := cmd.Flags().GetString("source")
			rankingLimit, _ := cmd.Flags().GetInt("ranking-limit")
			setString(params, "source", source)
			setInt(params, "ranking_limit", rankingLimit, 0)

			return f.runSingle(cmd, "/v1/ks/dashboard", params, nil)
		},
	}

	addTenantFlag(cmd)
	cmd.Flags().String("source", "", "data source filter")
	cmd.Flags().Int("ranking-limit", 10, "account ranking size (max 100)")
	addDateRangeFlags(cmd)
	addJQFlag(cmd)

	return cmd
}
