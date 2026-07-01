package cmd

import (
	"github.com/spf13/cobra"
)

// newKsAccountReportsCmd -> "adex ks account-reports {daily,summary}".
func newKsAccountReportsCmd(f *Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account-reports",
		Short: "Account-level daily & summary reports",
	}
	daily := newDailyCmd(f, "ks", "account-reports", "/v1/ks/account-reports/daily", colKsAccountReport, dailyHook{
		register: func(c *cobra.Command) {
			c.Flags().String("account-name", "", "account name fuzzy match")
		},
		collect: func(c *cobra.Command, params map[string]interface{}) {
			accountName, _ := c.Flags().GetString("account-name")
			setString(params, "account_name", accountName)
		},
	})
	summary := newSummaryCmd(f, "ks", "account-reports", "/v1/ks/account-reports/summary", "advertiser_id")
	cmd.AddCommand(daily, summary)
	return cmd
}

// newKsCampaignReportsCmd -> "adex ks campaign-reports {daily,summary}".
func newKsCampaignReportsCmd(f *Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "campaign-reports",
		Short: "Campaign-level daily & summary reports",
	}
	daily := newDailyCmd(f, "ks", "campaign-reports", "/v1/ks/campaign-reports/daily", colKsCampaignReport, dailyHook{
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
	summary := newSummaryCmd(f, "ks", "campaign-reports", "/v1/ks/campaign-reports/summary", "campaign_id")
	cmd.AddCommand(daily, summary)
	return cmd
}

// newKsUnitReportsCmd -> "adex ks unit-reports {daily,summary}".
func newKsUnitReportsCmd(f *Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unit-reports",
		Short: "Unit-level daily & summary reports",
	}
	daily := newDailyCmd(f, "ks", "unit-reports", "/v1/ks/unit-reports/daily", colKsUnitReport, dailyHook{
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
	summary := newSummaryCmd(f, "ks", "unit-reports", "/v1/ks/unit-reports/summary", "unit_id")
	cmd.AddCommand(daily, summary)
	return cmd
}

// newKsCreativeReportsCmd -> "adex ks creative-reports {daily,summary}".
func newKsCreativeReportsCmd(f *Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "creative-reports",
		Short: "Creative-level daily & summary reports",
	}
	daily := newDailyCmd(f, "ks", "creative-reports", "/v1/ks/creative-reports/daily", colKsCreativeReport, dailyHook{
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
	summary := newSummaryCmd(f, "ks", "creative-reports", "/v1/ks/creative-reports/summary", "creative_id")
	cmd.AddCommand(daily, summary)
	return cmd
}
