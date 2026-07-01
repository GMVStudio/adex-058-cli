package cmd

import (
	"github.com/spf13/cobra"
)

// newOeUnitsCmd creates "adex oe units" with list (default), top, get.
func newOeUnitsCmd(f *Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "units",
		Short: "List Oceanengine units (单元/广告)",
		Long: `List Oceanengine unit snapshots (GET /v1/oe/units).

The "get" subcommand takes the promotion_id path parameter.

Examples:
  adex oe units --tenant 6 --page-size 20
  adex oe units --tenant 6 --project 7650479670059647030 --format table
  adex oe units top --tenant 6 --range 30d --metric charge --limit 10
  adex oe units get 7650483929670156288 --tenant 6`,
		RunE: func(cmd *cobra.Command, args []string) error {
			tenant, err := requireTenant(cmd)
			if err != nil {
				return err
			}
			params := map[string]interface{}{"tenant_id": tenant}

			advertiserID, _ := cmd.Flags().GetString("advertiser")
			promotionID, _ := cmd.Flags().GetString("promotion")
			projectID, _ := cmd.Flags().GetString("project")
			name, _ := cmd.Flags().GetString("name")
			optStatus, _ := cmd.Flags().GetString("opt-status")
			statusFirst, _ := cmd.Flags().GetString("status-first")
			learningPhase, _ := cmd.Flags().GetString("learning-phase")

			setString(params, "advertiser_id", advertiserID)
			setString(params, "promotion_id", promotionID)
			setString(params, "project_id", projectID)
			setString(params, "name", name)
			setString(params, "opt_status", optStatus)
			setString(params, "status_first", statusFirst)
			setString(params, "learning_phase", learningPhase)

			applyOrder(cmd, params)
			applyPaging(cmd, params)

			return f.runList(cmd, "/v1/oe/units", params, colOeUnits)
		},
	}

	addTenantFlag(cmd)
	cmd.Flags().String("advertiser", "", "advertiser ID filter")
	cmd.Flags().String("promotion", "", "promotion (unit) ID filter")
	cmd.Flags().String("project", "", "project ID filter")
	cmd.Flags().String("name", "", "unit name fuzzy match")
	cmd.Flags().String("opt-status", "", "operation status ENABLE/DISABLE")
	cmd.Flags().String("status-first", "", "first-level status filter")
	cmd.Flags().String("learning-phase", "", "learning phase filter")
	addOrderFlags(cmd, "id")
	addPagingFlags(cmd)
	addJQFlag(cmd)

	cmd.AddCommand(
		newTopCmd(f, "oe", "units", "/v1/oe/units/top"),
		newGetCmd(f, "oe", "unit", "promotion_id", func(id string) string {
			return "/v1/oe/units/" + id
		}),
	)

	return cmd
}
