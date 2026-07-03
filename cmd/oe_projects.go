package cmd

import (
	"github.com/spf13/cobra"
)

// newOeProjectsCmd creates "adex oe projects" with list (default), top, get.
func newOeProjectsCmd(f *Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "projects",
		Short: "List Oceanengine projects (项目)",
		Long: `List Oceanengine project snapshots (GET /v1/oe/projects).

Examples:
  adex oe projects --tenant 6 --page-size 20
  adex oe projects --tenant 6 --opt-status ENABLE --format table
  adex oe projects top --tenant 6 --range 30d --metric charge --limit 10
  adex oe projects get 7650479670059647030 --tenant 6`,
		RunE: func(cmd *cobra.Command, args []string) error {
			tenant, err := f.resolveTenant(cmd)
			if err != nil {
				return err
			}
			params := map[string]interface{}{"tenant_id": tenant}

			advertiserID, _ := cmd.Flags().GetString("advertiser")
			projectID, _ := cmd.Flags().GetString("project")
			name, _ := cmd.Flags().GetString("name")
			optStatus, _ := cmd.Flags().GetString("opt-status")
			statusFirst, _ := cmd.Flags().GetString("status-first")
			deliveryMode, _ := cmd.Flags().GetString("delivery-mode")
			landingType, _ := cmd.Flags().GetString("landing-type")

			setString(params, "advertiser_id", advertiserID)
			setString(params, "project_id", projectID)
			setString(params, "name", name)
			setString(params, "opt_status", optStatus)
			setString(params, "status_first", statusFirst)
			setString(params, "delivery_mode", deliveryMode)
			setString(params, "landing_type", landingType)

			applyOrder(cmd, params)
			applyPaging(cmd, params)

			return f.runList(cmd, "/v1/oe/projects", params, colOeProjects)
		},
	}

	addTenantFlag(cmd)
	cmd.Flags().String("advertiser", "", "advertiser ID filter")
	cmd.Flags().String("project", "", "project ID filter")
	cmd.Flags().String("name", "", "project name fuzzy match")
	cmd.Flags().String("opt-status", "", "operation status ENABLE/DISABLE")
	cmd.Flags().String("status-first", "", "first-level status filter")
	cmd.Flags().String("delivery-mode", "", "delivery mode filter")
	cmd.Flags().String("landing-type", "", "landing type filter")
	addOrderFlags(cmd, "id")
	addPagingFlags(cmd)
	addJQFlag(cmd)

	cmd.AddCommand(
		newTopCmd(f, "oe", "projects", "/v1/oe/projects/top"),
		newGetCmd(f, "oe", "project", "project_id", func(id string) string {
			return "/v1/oe/projects/" + id
		}),
	)

	return cmd
}
