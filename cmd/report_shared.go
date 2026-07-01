package cmd

import (
	"github.com/gmvstudio/adex-cli/errs"
	"github.com/spf13/cobra"
)

// This file holds platform-agnostic command builders shared by both the
// Kuaishou (ks) and Oceanengine (oe) command trees. The platform argument is
// only used to render accurate help examples ("adex ks ..." vs "adex oe ...").

// newTopCmd builds a "top" subcommand shared by campaigns/projects/units/etc.
// They all hit a /top endpoint returning a summary reply ranked by a metric
// over a date range.
func newTopCmd(f *Factory, platform, resource, path string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "top",
		Short: "Top-N " + resource + " by a metric over a date range",
		Long: `Return the top-N ` + resource + ` ranked by a metric (GET ` + path + `).

Examples:
  adex ` + platform + ` ` + resource + ` top --tenant 6 --range 30d --metric charge --limit 10
  adex ` + platform + ` ` + resource + ` top --tenant 6 --range 30d --metric charge --order-desc=false --limit 5`,
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

			return f.runList(cmd, path, params, colSummary)
		},
	}

	addTenantFlag(cmd)
	cmd.Flags().String("advertiser", "", "advertiser ID filter")
	cmd.Flags().String("metric", "charge", "ranking metric (charge / convert_cnt / active ...)")
	cmd.Flags().String("source", "", "data source filter")
	cmd.Flags().Int("limit", 20, "number of rows to return (max 100)")
	cmd.Flags().Bool("order-desc", true, "sort descending (top values first)")
	addDateRangeFlags(cmd)
	addJQFlag(cmd)

	return cmd
}

// newGetCmd builds a "get <id>" subcommand shared across resources.
// pathFn builds the full request path from the positional ID argument.
func newGetCmd(f *Factory, platform, resource, idName string, pathFn func(id string) string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <" + idName + ">",
		Short: "Get a single " + resource + " detail",
		Long: `Get a single ` + resource + ` detail including the full detail/meta JSON.

Examples:
  adex ` + platform + ` ` + resource + ` get <` + idName + `> --tenant 6`,
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

// dailyHook lets each resource contribute its own daily-report filters while
// sharing the common tenant/date/paging plumbing.
type dailyHook struct {
	register func(cmd *cobra.Command)
	collect  func(cmd *cobra.Command, params map[string]interface{})
}

// newDailyCmd builds a shared "daily" report subcommand.
func newDailyCmd(f *Factory, platform, resource, path string, columns []string, hook dailyHook) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "daily",
		Short: "Daily " + resource + " report",
		Long: `Daily ` + resource + ` report (GET ` + path + `).

Examples:
  adex ` + platform + ` ` + resource + ` daily --tenant 6 --range 30d --page-size 20
  adex ` + platform + ` ` + resource + ` daily --tenant 6 --begin 2026-07-01 --end 2026-07-31 --format table`,
		RunE: func(cmd *cobra.Command, args []string) error {
			tenant, err := requireTenant(cmd)
			if err != nil {
				return err
			}
			params := map[string]interface{}{"tenant_id": tenant}

			if err := applyDateRange(cmd, params, false); err != nil {
				return err
			}

			advertiserID, _ := cmd.Flags().GetString("advertiser")
			source, _ := cmd.Flags().GetString("source")
			statHour, _ := cmd.Flags().GetInt("stat-hour")
			setString(params, "advertiser_id", advertiserID)
			setString(params, "source", source)
			setInt(params, "stat_hour", statHour, -1)

			if hook.collect != nil {
				hook.collect(cmd, params)
			}

			applyOrder(cmd, params)
			applyPaging(cmd, params)

			return f.runList(cmd, path, params, columns)
		},
	}

	addTenantFlag(cmd)
	cmd.Flags().String("advertiser", "", "advertiser ID filter")
	cmd.Flags().String("source", "", "data source filter")
	cmd.Flags().Int("stat-hour", -1, "stat hour granularity (-1 = all)")
	if hook.register != nil {
		hook.register(cmd)
	}
	addDateRangeFlags(cmd)
	addOrderFlags(cmd, "stat_date")
	addPagingFlags(cmd)
	addJQFlag(cmd)

	return cmd
}

// newSummaryCmd builds a shared "summary" report subcommand. groupHint names
// the supported group_by dimension for help text.
func newSummaryCmd(f *Factory, platform, resource, path, groupHint string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "summary",
		Short: "Aggregate " + resource + " report over a date range",
		Long: `Aggregate ` + resource + ` report over a date range (GET ` + path + `).

Leave --group-by empty for a single total row, or set it to ` + groupHint + `
to break the total down per dimension.

Examples:
  adex ` + platform + ` ` + resource + ` summary --tenant 6 --range 30d
  adex ` + platform + ` ` + resource + ` summary --tenant 6 --range 30d --group-by ` + groupHint + ` --order-by charge --order-desc`,
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
			groupBy, _ := cmd.Flags().GetString("group-by")
			source, _ := cmd.Flags().GetString("source")
			setString(params, "advertiser_id", advertiserID)
			setString(params, "group_by", groupBy)
			setString(params, "source", source)

			applyOrder(cmd, params)
			applyPaging(cmd, params)

			return f.runList(cmd, path, params, colSummary)
		},
	}

	addTenantFlag(cmd)
	cmd.Flags().String("advertiser", "", "advertiser ID filter")
	cmd.Flags().String("group-by", "", "group dimension ("+groupHint+"); empty = single total row")
	cmd.Flags().String("source", "", "data source filter")
	addDateRangeFlags(cmd)
	addOrderFlags(cmd, "charge")
	addPagingFlags(cmd)
	addJQFlag(cmd)

	return cmd
}
