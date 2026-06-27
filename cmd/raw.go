package cmd

import (
	"github.com/spf13/cobra"
)

// newRawCmd creates the "raw" command group for direct API access.
func newRawCmd(f *Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "raw",
		Short: "Raw API access for reports and data queries",
		Long: `Raw API access — query ADEX reports and data directly.

Subcommands:
  campaign   Campaign-level reports
    daily     Daily campaign report

Examples:
  adex raw campaign daily --tenant 6 --campaign C-618-001-619 --range 1d
  adex raw campaign daily --tenant 6 --range 7d --format table`,
	}

	cmd.AddCommand(newRawCampaignCmd(f))

	return cmd
}
