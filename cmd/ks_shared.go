package cmd

import (
	"github.com/gmvstudio/adex-cli/errs"
	"github.com/spf13/cobra"
)

// newKsTopCmd builds a "top" subcommand shared by campaigns/units/creatives.
// They all hit a /top endpoint that returns a SumKsReportReply ranked by a
// metric over a date range.
func newKsTopCmd(f *Factory, resource, path string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "top",
		Short: "Top-N " + resource + " by a metric over a date range",
		Long: `Return the top-N ` + resource + ` ranked by a metric (GET ` + path + `).

Examples:
  adex ks ` + resource + ` top --tenant 6 --range 30d --metric charge --limit 10
  adex ks ` + resource + ` top --tenant 6 --range 30d --metric charge --order-desc=false --limit 5`,
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
			metric, _ := cmd.Flags().GetString("metric")
			source, _ := cmd.Flags().GetString("source")
			limit, _ := cmd.Flags().GetInt("limit")
			orderDesc, _ := cmd.Flags().GetBool("order-desc")

			setString(params, "advertiser_id", advertiserID)
			setString(params, "metric", metric)
			setString(params, "source", source)
			setInt(params, "limit", limit, 0)
			params["order_desc"] = orderDesc

			return f.runList(cmd, path, params, colKsSummary)
		},
	}

	addTenantFlag(cmd)
	cmd.Flags().String("advertiser", "", "advertiser ID filter")
	cmd.Flags().String("metric", "charge", "ranking metric (charge / conversion_num / activation ...)")
	cmd.Flags().String("source", "", "data source filter")
	cmd.Flags().Int("limit", 20, "number of rows to return (max 100)")
	cmd.Flags().Bool("order-desc", true, "sort descending (top values first)")
	addDateRangeFlags(cmd)
	addJQFlag(cmd)

	return cmd
}

// newKsGetCmd builds a "get <id>" subcommand shared by campaigns/units/creatives.
// pathFn builds the full request path from the positional ID argument.
func newKsGetCmd(f *Factory, resource, idName string, pathFn func(id string) string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <" + idName + ">",
		Short: "Get a single " + resource + " detail",
		Long: `Get a single ` + resource + ` detail including the full detail JSON.

Examples:
  adex ks ` + resource + ` get 9899931248 --tenant 6`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			if id == "" {
				return errs.NewValidationError(errs.SubtypeInvalidArgument,
					"%s is required", idName).WithParam(idName)
			}
			tenant, err := requireTenant(cmd)
			if err != nil {
				return err
			}
			params := map[string]interface{}{"tenant_id": tenant}
			return f.runSingle(cmd, pathFn(id), params)
		},
	}

	addTenantFlag(cmd)
	addJQFlag(cmd)

	return cmd
}
