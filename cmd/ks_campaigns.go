package cmd

import (
	"github.com/spf13/cobra"
)

// newKsCampaignsCmd creates "adex ks campaigns" with list (default), top, get.
func newKsCampaignsCmd(f *Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "campaigns",
		Short: "List Kuaishou campaigns",
		Long: `List Kuaishou campaign snapshots (GET /v1/ks/campaigns).

Examples:
  adex ks campaigns --tenant 6 --page-size 20
  adex ks campaigns --tenant 6 --put-status 1 --format table
  adex ks campaigns top --tenant 6 --range 30d --metric charge --limit 10
  adex ks campaigns get 9899931248 --tenant 6`,
		RunE: func(cmd *cobra.Command, args []string) error {
			tenant, err := requireTenant(cmd)
			if err != nil {
				return err
			}
			params := map[string]interface{}{"tenant_id": tenant}

			advertiserID, _ := cmd.Flags().GetString("advertiser")
			campaignID, _ := cmd.Flags().GetString("campaign")
			campaignName, _ := cmd.Flags().GetString("campaign-name")
			putStatus, _ := cmd.Flags().GetInt("put-status")
			status, _ := cmd.Flags().GetInt("status")
			campaignType, _ := cmd.Flags().GetInt("campaign-type")

			setString(params, "advertiser_id", advertiserID)
			setString(params, "campaign_id", campaignID)
			setString(params, "campaign_name", campaignName)
			setInt(params, "put_status", putStatus, 0)
			setInt(params, "status", status, 0)
			setInt(params, "campaign_type", campaignType, 0)

			applyOrder(cmd, params)
			applyPaging(cmd, params)

			return f.runList(cmd, "/v1/ks/campaigns", params, colKsCampaigns)
		},
	}

	addTenantFlag(cmd)
	cmd.Flags().String("advertiser", "", "advertiser ID filter")
	cmd.Flags().String("campaign", "", "campaign ID filter")
	cmd.Flags().String("campaign-name", "", "campaign name fuzzy match")
	cmd.Flags().Int("put-status", 0, "put status 1=on/2=paused/3=deleted (0=all)")
	cmd.Flags().Int("status", 0, "campaign status (0=all)")
	cmd.Flags().Int("campaign-type", 0, "campaign type (0=all)")
	addOrderFlags(cmd, "id")
	addPagingFlags(cmd)
	addJQFlag(cmd)

	cmd.AddCommand(
		newTopCmd(f, "ks", "campaigns", "/v1/ks/campaigns/top"),
		newGetCmd(f, "ks", "campaign", "campaign_id", func(id string) string {
			return "/v1/ks/campaigns/" + id
		}),
	)

	return cmd
}
