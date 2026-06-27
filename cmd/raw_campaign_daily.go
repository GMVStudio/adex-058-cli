package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gmvstudio/adex-cli/errs"
	"github.com/gmvstudio/adex-cli/internal/client"
	"github.com/gmvstudio/adex-cli/internal/config"
	"github.com/gmvstudio/adex-cli/internal/output"
	"github.com/spf13/cobra"
)

// parseRange converts a range string like "1d", "7d", "1h", "30m" into
// start and end date strings (YYYY-MM-DD). "1d" means the last 1 day from now.
func parseRange(r string) (string, string, error) {
	r = strings.TrimSpace(r)
	if r == "" {
		return "", "", errs.NewValidationError(errs.SubtypeInvalidArgument, "--range is required").WithParam("--range").WithHint("use a value like 1d, 7d, 1h, 30m")
	}

	unit := r[len(r)-1]
	numStr := r[:len(r)-1]
	num, err := strconv.Atoi(numStr)
	if err != nil || num <= 0 {
		return "", "", errs.NewValidationError(errs.SubtypeInvalidArgument, "invalid range value %q: expected format like 1d, 7d, 1h, 30m", r).WithParam("--range")
	}

	now := time.Now()
	var start time.Time
	switch unit {
	case 'd':
		start = now.AddDate(0, 0, -num)
	case 'h':
		start = now.Add(-time.Duration(num) * time.Hour)
	case 'm':
		start = now.Add(-time.Duration(num) * time.Minute)
	default:
		return "", "", errs.NewValidationError(errs.SubtypeInvalidArgument, "unsupported range unit %q in %q: use d (days), h (hours), or m (minutes)", string(unit), r).WithParam("--range")
	}

	return start.Format("2006-01-02"), now.Format("2006-01-02"), nil
}

// newRawCampaignDailyCmd creates the "raw campaign daily" command.
func newRawCampaignDailyCmd(f *Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "daily",
		Short: "Daily campaign report",
		Long: `Query daily campaign reports from the ADEX API.

Required flags:
  --tenant     Tenant ID
  --range      Time range (e.g. 1d, 7d, 1h, 30m)

Optional flags:
  --campaign   Campaign ID (numeric) or name pattern (string)
  --page       Page number (default: 1)
  --page-size  Page size (default: 20)
  --order-by   Sort field (default: charge)
  --order-desc Sort descending (default: true)
  --stat-hour  Stat hour, -1 for latest (default: -1)

Examples:
  adex raw campaign daily --tenant 6 --campaign C-618-001-619 --range 1d
  adex raw campaign daily --tenant 6 --campaign 9455214173 --range 7d
  adex raw campaign daily --tenant 6 --range 7d --format table
  adex raw campaign daily --tenant 6 --range 1d --page 2 --page-size 50`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCampaignDaily(cmd, f)
		},
	}

	cmd.Flags().Int("tenant", 0, "tenant ID (required)")
	cmd.Flags().String("campaign", "", "campaign ID (numeric) or name pattern (string)")
	cmd.Flags().String("range", "", "time range: 1d, 7d, 1h, 30m (required)")
	cmd.Flags().Int("page", config.DefaultPage, "page number")
	cmd.Flags().Int("page-size", config.DefaultPageSize, "page size")
	cmd.Flags().String("order-by", config.DefaultOrderBy, "sort field")
	cmd.Flags().Bool("order-desc", config.DefaultOrderDesc, "sort descending")
	cmd.Flags().Int("stat-hour", config.DefaultStatHour, "stat hour (-1 for latest)")

	_ = cmd.MarkFlagRequired("tenant")
	_ = cmd.MarkFlagRequired("range")

	return cmd
}

func runCampaignDaily(cmd *cobra.Command, f *Factory) error {
	tenantID, _ := cmd.Flags().GetInt("tenant")
	if tenantID <= 0 {
		return errs.NewValidationError(errs.SubtypeInvalidArgument, "--tenant must be a positive integer").WithParam("--tenant")
	}

	campaignID, _ := cmd.Flags().GetString("campaign")
	rangeStr, _ := cmd.Flags().GetString("range")
	page, _ := cmd.Flags().GetInt("page")
	pageSize, _ := cmd.Flags().GetInt("page-size")
	orderBy, _ := cmd.Flags().GetString("order-by")
	orderDesc, _ := cmd.Flags().GetBool("order-desc")
	statHour, _ := cmd.Flags().GetInt("stat-hour")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	startDate, endDate, err := parseRange(rangeStr)
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"tenant_id":  tenantID,
		"page":       page,
		"page_size":  pageSize,
		"order_by":   orderBy,
		"order_desc": orderDesc,
		"stat_hour":  statHour,
		"start_date": startDate,
		"end_date":   endDate,
	}
	if campaignID != "" {
		if _, err := strconv.ParseInt(campaignID, 10, 64); err == nil {
			params["campaign_id"] = campaignID
		} else {
			params["campaign_name"] = campaignID
		}
	}

	if dryRun {
		fmt.Fprintf(f.ErrOut, "DRY RUN: GET %s/v1/campaign-reports/daily\n", f.Config.BaseURL)
		pretty, _ := json.MarshalIndent(params, "", "  ")
		fmt.Fprintf(f.ErrOut, "Params:\n%s\n", string(pretty))
		return nil
	}

	c := f.resolveClient(cmd)
	ctx := context.Background()

	req := client.Request{
		Method: "GET",
		Path:   "/v1/campaign-reports/daily",
		Params: params,
	}

	data, err := c.Do(ctx, req)
	if err != nil {
		return err
	}

	// Parse for structured output
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		// Fall back to raw output
		fmt.Fprintln(f.Out, string(data))
		return nil
	}

	formatStr := f.resolveFormat(cmd)
	format := output.ParseFormat(formatStr)

	// For table format, use sensible default columns
	columns := []string{"id", "tenantId", "campaignId", "statDate", "charge", "metrics.aclick", "metrics.conversion_num"}

	return output.Print(f.Out, result, format, columns)
}
