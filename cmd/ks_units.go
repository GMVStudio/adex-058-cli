package cmd

import (
	"github.com/spf13/cobra"
)

// newKsUnitsCmd creates "adex ks units" with list (default), top, get.
func newKsUnitsCmd(f *Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "units",
		Short: "List Kuaishou ad units (广告组)",
		Long: `List Kuaishou ad unit snapshots (GET /v1/ks/units).

Examples:
  adex ks units --tenant 6 --page-size 20
  adex ks units --tenant 6 --campaign 9899931248 --format table
  adex ks units top --tenant 6 --range 30d --metric conversion_num --limit 20
  adex ks units get 29638466721 --tenant 6`,
		RunE: func(cmd *cobra.Command, args []string) error {
			tenant, err := requireTenant(cmd)
			if err != nil {
				return err
			}
			params := map[string]interface{}{"tenant_id": tenant}

			advertiserID, _ := cmd.Flags().GetString("advertiser")
			unitID, _ := cmd.Flags().GetString("unit")
			campaignID, _ := cmd.Flags().GetString("campaign")
			unitName, _ := cmd.Flags().GetString("unit-name")
			putStatus, _ := cmd.Flags().GetInt("put-status")
			status, _ := cmd.Flags().GetInt("status")

			setString(params, "advertiser_id", advertiserID)
			setString(params, "unit_id", unitID)
			setString(params, "campaign_id", campaignID)
			setString(params, "unit_name", unitName)
			setInt(params, "put_status", putStatus, 0)
			setInt(params, "status", status, 0)

			applyOrder(cmd, params)
			applyPaging(cmd, params)

			return f.runList(cmd, "/v1/ks/units", params, colKsUnits)
		},
	}

	addTenantFlag(cmd)
	cmd.Flags().String("advertiser", "", "advertiser ID filter")
	cmd.Flags().String("unit", "", "unit ID filter")
	cmd.Flags().String("campaign", "", "campaign ID filter")
	cmd.Flags().String("unit-name", "", "unit name fuzzy match")
	cmd.Flags().Int("put-status", 0, "put status (0=all)")
	cmd.Flags().Int("status", 0, "unit status (0=all)")
	addOrderFlags(cmd, "id")
	addPagingFlags(cmd)
	addJQFlag(cmd)

	cmd.AddCommand(
		newKsTopCmd(f, "units", "/v1/ks/units/top"),
		newKsGetCmd(f, "unit", "unit_id", func(id string) string {
			return "/v1/ks/units/" + id
		}),
	)

	return cmd
}
