package cmd

import (
	"github.com/spf13/cobra"
)

// dailyHook lets each resource contribute its own daily-report filters while
// sharing the common tenant/date/paging plumbing.
type dailyHook struct {
	register func(cmd *cobra.Command)
	collect  func(cmd *cobra.Command, params map[string]interface{})
}

// newKsDailyCmd builds a shared "daily" report subcommand.
func newKsDailyCmd(f *Factory, resource, path string, columns []string, hook dailyHook) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "daily",
		Short: "Daily " + resource + " report",
		Long: `Daily ` + resource + ` report (GET ` + path + `).

Examples:
  adex ks ` + resource + ` daily --tenant 6 --range 30d --page-size 20
  adex ks ` + resource + ` daily --tenant 6 --begin 2026-07-01 --end 2026-07-31 --format table`,
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

// newKsSummaryCmd builds a shared "summary" report subcommand. groupHint names
// the supported group_by dimension for help text.
func newKsSummaryCmd(f *Factory, resource, path, groupHint string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "summary",
		Short: "Aggregate " + resource + " report over a date range",
		Long: `Aggregate ` + resource + ` report over a date range (GET ` + path + `).

Leave --group-by empty for a single total row, or set it to ` + groupHint + `
to break the total down per dimension.

Examples:
  adex ks ` + resource + ` summary --tenant 6 --range 30d
  adex ks ` + resource + ` summary --tenant 6 --range 30d --group-by ` + groupHint + ` --order-by charge --order-desc`,
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

			return f.runList(cmd, path, params, colKsSummary)
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

// newKsAccountReportsCmd -> "adex ks account-reports {daily,summary}".
func newKsAccountReportsCmd(f *Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account-reports",
		Short: "Account-level daily & summary reports",
	}
	daily := newKsDailyCmd(f, "account-reports", "/v1/ks/account-reports/daily", colKsAccountReport, dailyHook{
		register: func(c *cobra.Command) {
			c.Flags().String("account-name", "", "account name fuzzy match")
		},
		collect: func(c *cobra.Command, params map[string]interface{}) {
			accountName, _ := c.Flags().GetString("account-name")
			setString(params, "account_name", accountName)
		},
	})
	summary := newKsSummaryCmd(f, "account-reports", "/v1/ks/account-reports/summary", "advertiser_id")
	cmd.AddCommand(daily, summary)
	return cmd
}

// newKsCampaignReportsCmd -> "adex ks campaign-reports {daily,summary}".
func newKsCampaignReportsCmd(f *Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "campaign-reports",
		Short: "Campaign-level daily & summary reports",
	}
	daily := newKsDailyCmd(f, "campaign-reports", "/v1/ks/campaign-reports/daily", colKsCampaignReport, dailyHook{
		register: func(c *cobra.Command) {
			c.Flags().String("campaign", "", "campaign ID filter")
			c.Flags().String("campaign-name", "", "campaign name fuzzy match")
			c.Flags().Int("status", 0, "campaign status (0=all)")
		},
		collect: func(c *cobra.Command, params map[string]interface{}) {
			campaignID, _ := c.Flags().GetString("campaign")
			campaignName, _ := c.Flags().GetString("campaign-name")
			status, _ := c.Flags().GetInt("status")
			setString(params, "campaign_id", campaignID)
			setString(params, "campaign_name", campaignName)
			setInt(params, "status", status, 0)
		},
	})
	summary := newKsSummaryCmd(f, "campaign-reports", "/v1/ks/campaign-reports/summary", "campaign_id")
	cmd.AddCommand(daily, summary)
	return cmd
}

// newKsUnitReportsCmd -> "adex ks unit-reports {daily,summary}".
func newKsUnitReportsCmd(f *Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unit-reports",
		Short: "Unit-level daily & summary reports",
	}
	daily := newKsDailyCmd(f, "unit-reports", "/v1/ks/unit-reports/daily", colKsUnitReport, dailyHook{
		register: func(c *cobra.Command) {
			c.Flags().String("unit", "", "unit ID filter")
			c.Flags().String("campaign", "", "campaign ID filter")
			c.Flags().String("unit-name", "", "unit name fuzzy match")
			c.Flags().Int("status", 0, "unit status (0=all)")
		},
		collect: func(c *cobra.Command, params map[string]interface{}) {
			unitID, _ := c.Flags().GetString("unit")
			campaignID, _ := c.Flags().GetString("campaign")
			unitName, _ := c.Flags().GetString("unit-name")
			status, _ := c.Flags().GetInt("status")
			setString(params, "unit_id", unitID)
			setString(params, "campaign_id", campaignID)
			setString(params, "unit_name", unitName)
			setInt(params, "status", status, 0)
		},
	})
	summary := newKsSummaryCmd(f, "unit-reports", "/v1/ks/unit-reports/summary", "unit_id")
	cmd.AddCommand(daily, summary)
	return cmd
}

// newKsCreativeReportsCmd -> "adex ks creative-reports {daily,summary}".
func newKsCreativeReportsCmd(f *Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "creative-reports",
		Short: "Creative-level daily & summary reports",
	}
	daily := newKsDailyCmd(f, "creative-reports", "/v1/ks/creative-reports/daily", colKsCreativeReport, dailyHook{
		register: func(c *cobra.Command) {
			c.Flags().String("creative", "", "creative ID filter")
			c.Flags().String("unit", "", "unit ID filter")
			c.Flags().String("campaign", "", "campaign ID filter")
			c.Flags().String("creative-name", "", "creative name fuzzy match")
			c.Flags().Int("status", 0, "creative status (0=all)")
		},
		collect: func(c *cobra.Command, params map[string]interface{}) {
			creativeID, _ := c.Flags().GetString("creative")
			unitID, _ := c.Flags().GetString("unit")
			campaignID, _ := c.Flags().GetString("campaign")
			creativeName, _ := c.Flags().GetString("creative-name")
			status, _ := c.Flags().GetInt("status")
			setString(params, "creative_id", creativeID)
			setString(params, "unit_id", unitID)
			setString(params, "campaign_id", campaignID)
			setString(params, "creative_name", creativeName)
			setInt(params, "status", status, 0)
		},
	})
	summary := newKsSummaryCmd(f, "creative-reports", "/v1/ks/creative-reports/summary", "creative_id")
	cmd.AddCommand(daily, summary)
	return cmd
}
