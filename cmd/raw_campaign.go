package cmd

import (
	"github.com/spf13/cobra"
)

// newRawCampaignCmd creates the "raw campaign" subcommand group.
func newRawCampaignCmd(f *Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "campaign",
		Short: "Campaign-level reports",
		Long: `Campaign-level reports.

Subcommands:
  daily    Daily campaign report

Examples:
  adex raw campaign daily --tenant 6 --campaign C-618-001-619 --range 1d`,
	}

	cmd.AddCommand(newRawCampaignDailyCmd(f))

	return cmd
}
