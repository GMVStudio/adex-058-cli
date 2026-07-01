package cmd

import (
	"github.com/spf13/cobra"
)

// newOeMetricMetaCmd creates "adex oe report-metric-meta" (ListOeReportMetricMeta).
// Note: this endpoint is tenant-agnostic (no --tenant flag).
func newOeMetricMetaCmd(f *Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report-metric-meta",
		Short: "List Oceanengine report metric metadata",
		Long: `List report metric metadata (GET /v1/oe/report-metric-meta).

Examples:
  adex oe report-metric-meta --level account --page-size 50
  adex oe report-metric-meta --level project --enabled 1 --page-size 50`,
		RunE: func(cmd *cobra.Command, args []string) error {
			params := map[string]interface{}{}

			level, _ := cmd.Flags().GetString("level")
			groupName, _ := cmd.Flags().GetString("group-name")
			field, _ := cmd.Flags().GetString("field")
			enabled, _ := cmd.Flags().GetInt("enabled")

			setString(params, "level", level)
			setString(params, "group_name", groupName)
			setString(params, "field", field)
			setInt(params, "enabled", enabled, 0)

			applyOrder(cmd, params)
			applyPaging(cmd, params)

			return f.runList(cmd, "/v1/oe/report-metric-meta", params, colOeMetricMeta)
		},
	}

	cmd.Flags().String("level", "", "dimension: account/project/unit")
	cmd.Flags().String("group-name", "", "metric group name filter")
	cmd.Flags().String("field", "", "field name fuzzy match")
	cmd.Flags().Int("enabled", 0, "enabled filter: 0=all / 1=enabled / 2=disabled")
	addOrderFlags(cmd, "sort_order")
	addPagingFlags(cmd)
	addJQFlag(cmd)

	return cmd
}
