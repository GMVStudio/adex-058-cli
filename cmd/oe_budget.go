package cmd

import (
	"github.com/spf13/cobra"
)

// newOeBudgetVsActualCmd creates "adex oe account-budget-vs-actual"
// (ListOeAccountBudgetVsActual). This endpoint returns an unpaginated list.
func newOeBudgetVsActualCmd(f *Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account-budget-vs-actual",
		Short: "Per-account daily budget vs. actual average spend",
		Long: `Compare each account's daily budget against its average daily actual spend
over a date range (GET /v1/oe/account-budget-vs-actual).

Examples:
  adex oe account-budget-vs-actual --tenant 6 --range 30d
  adex oe account-budget-vs-actual --tenant 6 --begin 2026-06-01 --end 2026-06-30 --format table
  adex oe account-budget-vs-actual --tenant 6 --advertiser 1866874042754522 --range 30d`,
		RunE: func(cmd *cobra.Command, args []string) error {
			tenant, err := requireTenant(cmd)
			if err != nil {
				return err
			}
			params := map[string]interface{}{"tenant_id": tenant}

			if err := applyDateRange(cmd, params, true); err != nil {
				return err
			}

			advertiserID, _ := cmd.Flags().GetString("advertiser")
			setString(params, "advertiser_id", advertiserID)

			return f.runList(cmd, "/v1/oe/account-budget-vs-actual", params, colOeBudgetVsActual)
		},
	}

	addTenantFlag(cmd)
	cmd.Flags().String("advertiser", "", "advertiser ID filter (single account)")
	addDateRangeFlags(cmd)
	addJQFlag(cmd)

	return cmd
}
