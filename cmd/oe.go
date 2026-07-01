package cmd

import "github.com/spf13/cobra"

// newOeCmd creates the "oe" (Oceanengine / 巨量) command group.
func newOeCmd(f *Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "oe",
		Short: "Oceanengine (巨量) advertising data",
		Long: `Oceanengine advertising queries: accounts, projects, units,
daily/summary reports, top-N rankings, metric metadata, dashboards, and
budget-vs-actual comparison.

Common flags (shared across commands):
  --tenant       tenant ID (required for most commands)
  --page-size    page size (default 20)
  --page-token   fetch a specific page
  --page-all     aggregate every page
  --order-by     sort field
  --order-desc   sort descending (default true)
  --range        relative date range like 7d/4w/1m
  --begin/--end  explicit stat date range (YYYY-MM-DD)
  --jq           filter JSON output with a jq expression
  --format       json (default) | pretty | table

Examples:
  adex oe projects --tenant 6 --page-size 20
  adex oe project-reports summary --tenant 6 --range 30d --group-by project_id
  adex oe units top --tenant 6 --range 30d --metric convert_cnt --limit 20
  adex oe dashboard --tenant 6 --range 30d
  adex oe account-budget-vs-actual --tenant 6 --range 30d`,
	}

	cmd.AddCommand(
		newOeAccountsCmd(f),
		newOeProjectsCmd(f),
		newOeUnitsCmd(f),
		newOeAccountReportsCmd(f),
		newOeProjectReportsCmd(f),
		newOeUnitReportsCmd(f),
		newOeMetricMetaCmd(f),
		newOeDashboardCmd(f),
		newOeBudgetVsActualCmd(f),
	)

	return cmd
}
