package cmd

import (
	"github.com/spf13/cobra"
)

// newKsCreativesCmd creates "adex ks creatives" with list (default), top, get.
func newKsCreativesCmd(f *Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "creatives",
		Short: "List Kuaishou creatives (创意)",
		Long: `List Kuaishou creative snapshots (GET /v1/ks/creatives).

The "get" subcommand takes the biz_key path parameter (e.g. p:29637782154).

Examples:
  adex ks creatives --tenant 6 --page-size 20
  adex ks creatives --tenant 6 --unit 29638466721 --format table
  adex ks creatives top --tenant 6 --range 30d --metric charge --limit 10
  adex ks creatives get p:29637782154 --tenant 6`,
		RunE: func(cmd *cobra.Command, args []string) error {
			tenant, err := f.resolveTenant(cmd)
			if err != nil {
				return err
			}
			params := map[string]interface{}{"tenant_id": tenant}

			advertiserID, _ := cmd.Flags().GetString("advertiser")
			unitID, _ := cmd.Flags().GetString("unit")
			campaignID, _ := cmd.Flags().GetString("campaign")
			creativeID, _ := cmd.Flags().GetString("creative")
			creativeName, _ := cmd.Flags().GetString("creative-name")
			creativeType, _ := cmd.Flags().GetString("creative-type")
			putStatus, _ := cmd.Flags().GetInt("put-status")
			status, _ := cmd.Flags().GetInt("status")

			setString(params, "advertiser_id", advertiserID)
			setString(params, "unit_id", unitID)
			setString(params, "campaign_id", campaignID)
			setString(params, "creative_id", creativeID)
			setString(params, "creative_name", creativeName)
			setString(params, "creative_type", creativeType)
			setInt(params, "put_status", putStatus, 0)
			setInt(params, "status", status, 0)

			applyOrder(cmd, params)
			applyPaging(cmd, params)

			return f.runList(cmd, "/v1/ks/creatives", params, colKsCreatives)
		},
	}

	addTenantFlag(cmd)
	cmd.Flags().String("advertiser", "", "advertiser ID filter")
	cmd.Flags().String("unit", "", "unit ID filter")
	cmd.Flags().String("campaign", "", "campaign ID filter")
	cmd.Flags().String("creative", "", "creative ID filter")
	cmd.Flags().String("creative-name", "", "creative name fuzzy match")
	cmd.Flags().String("creative-type", "", "creative type filter")
	cmd.Flags().Int("put-status", 0, "put status (0=all)")
	cmd.Flags().Int("status", 0, "creative status (0=all)")
	addOrderFlags(cmd, "id")
	addPagingFlags(cmd)
	addJQFlag(cmd)

	cmd.AddCommand(
		newTopCmd(f, "ks", "creatives", "/v1/ks/creatives/top"),
		newGetCmd(f, "ks", "creative", "biz_key", func(id string) string {
			return "/v1/ks/creatives/" + id
		}),
	)

	return cmd
}
