package cmd

import (
	"github.com/spf13/cobra"
)

// newKsMetricMetaCmd creates "adex ks report-metric-meta" (ListKsReportMetricMeta).
// Note: this endpoint is tenant-agnostic (no --tenant flag).
func newKsMetricMetaCmd(f *Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report-metric-meta",
		Short: "List Kuaishou report metric metadata",
		Long: `List report metric metadata (GET /v1/ks/report-metric-meta).

Examples:
  adex ks report-metric-meta --level account --page-size 50
  adex ks report-metric-meta --level campaign --enabled 1 --page-size 50`,
		RunE: func(cmd *cobra.Command, args []string) error {
			params := map[string]interface{}{}

			level, _ := cmd.Flags().GetString("level")
			groupName, _ := cmd.Flags().GetString("group-name")
			field, _ := cmd.Flags().GetString("field")
			enabled, _ := cmd.Flags().GetInt("enabled")
			sortable, _ := cmd.Flags().GetInt("sortable")

			setString(params, "level", level)
			setString(params, "group_name", groupName)
			setString(params, "field", field)
			setInt(params, "enabled", enabled, 0)
			setInt(params, "sortable", sortable, 0)

			applyOrder(cmd, params)
			applyPaging(cmd, params)

			return f.runList(cmd, "/v1/ks/report-metric-meta", params, colKsMetricMeta)
		},
	}

	cmd.Flags().String("level", "", "dimension: account/campaign/unit/creative")
	cmd.Flags().String("group-name", "", "metric group name filter")
	cmd.Flags().String("field", "", "field name fuzzy match")
	cmd.Flags().Int("enabled", 0, "enabled filter: 0=all / 1=enabled / 2=disabled")
	cmd.Flags().Int("sortable", 0, "sortable filter: 0=all / 1=sortable / 2=not sortable")
	addOrderFlags(cmd, "sort_order")
	addPagingFlags(cmd)
	addJQFlag(cmd)

	return cmd
}
