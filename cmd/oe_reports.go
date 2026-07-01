package cmd

import (
	"github.com/spf13/cobra"
)

// newOeAccountReportsCmd -> "adex oe account-reports {daily,summary}".
func newOeAccountReportsCmd(f *Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account-reports",
		Short: "Account-level daily & summary reports",
	}
	daily := newDailyCmd(f, "oe", "account-reports", "/v1/oe/account-reports/daily", colOeAccountReport, dailyHook{})
	summary := newSummaryCmd(f, "oe", "account-reports", "/v1/oe/account-reports/summary", "advertiser_id")
	cmd.AddCommand(daily, summary)
	return cmd
}

// newOeProjectReportsCmd -> "adex oe project-reports {daily,summary}".
func newOeProjectReportsCmd(f *Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project-reports",
		Short: "Project-level daily & summary reports",
	}
	daily := newDailyCmd(f, "oe", "project-reports", "/v1/oe/project-reports/daily", colOeProjectReport, dailyHook{
		register: func(c *cobra.Command) {
			c.Flags().String("project", "", "project ID filter")
			c.Flags().String("project-name", "", "project name fuzzy match")
		},
		collect: func(c *cobra.Command, params map[string]interface{}) {
			projectID, _ := c.Flags().GetString("project")
			projectName, _ := c.Flags().GetString("project-name")
			setString(params, "project_id", projectID)
			setString(params, "project_name", projectName)
		},
	})
	summary := newSummaryCmd(f, "oe", "project-reports", "/v1/oe/project-reports/summary", "project_id")
	cmd.AddCommand(daily, summary)
	return cmd
}

// newOeUnitReportsCmd -> "adex oe unit-reports {daily,summary}".
func newOeUnitReportsCmd(f *Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unit-reports",
		Short: "Unit-level daily & summary reports",
	}
	daily := newDailyCmd(f, "oe", "unit-reports", "/v1/oe/unit-reports/daily", colOeUnitReport, dailyHook{
		register: func(c *cobra.Command) {
			c.Flags().String("project", "", "project ID filter")
			c.Flags().String("promotion", "", "promotion (unit) ID filter")
			c.Flags().String("promotion-name", "", "unit name fuzzy match")
		},
		collect: func(c *cobra.Command, params map[string]interface{}) {
			projectID, _ := c.Flags().GetString("project")
			promotionID, _ := c.Flags().GetString("promotion")
			promotionName, _ := c.Flags().GetString("promotion-name")
			setString(params, "project_id", projectID)
			setString(params, "promotion_id", promotionID)
			setString(params, "promotion_name", promotionName)
		},
	})
	summary := newSummaryCmd(f, "oe", "unit-reports", "/v1/oe/unit-reports/summary", "promotion_id")
	cmd.AddCommand(daily, summary)
	return cmd
}
