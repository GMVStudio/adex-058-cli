package cmd

import "github.com/spf13/cobra"

// newKsCmd creates the "ks" (Kuaishou) command group.
func newKsCmd(f *Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ks",
		Short: "Kuaishou (快手) advertising data",
		Long: `Kuaishou advertising queries: accounts, campaigns, units, creatives,
daily/summary reports, top-N rankings, metric metadata, and dashboards.

Common flags (shared across commands):
  --tenant       tenant ID (optional; uses default from 'adex tenant use')
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
  adex ks accounts --page-size 20
  adex ks campaign-reports summary --range 30d --group-by campaign_id
  adex ks campaigns top --range 30d --metric charge --limit 10
  adex ks dashboard --range 30d`,
	}

	cmd.AddCommand(
		newKsAccountsCmd(f),
		newKsCampaignsCmd(f),
		newKsUnitsCmd(f),
		newKsCreativesCmd(f),
		newKsAccountReportsCmd(f),
		newKsCampaignReportsCmd(f),
		newKsUnitReportsCmd(f),
		newKsCreativeReportsCmd(f),
		newKsMetricMetaCmd(f),
		newKsDashboardCmd(f),
	)

	return cmd
}
