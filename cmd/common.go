package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gmvstudio/adex-cli/errs"
	"github.com/gmvstudio/adex-cli/internal/client"
	"github.com/gmvstudio/adex-cli/internal/config"
	"github.com/gmvstudio/adex-cli/internal/daterange"
	"github.com/gmvstudio/adex-cli/internal/output"
	"github.com/gmvstudio/adex-cli/internal/paginate"
	"github.com/spf13/cobra"
)

// ---------------------------------------------------------------------------
// Common flag registration
//
// These helpers keep flag definitions consistent across every command instead
// of re-declaring the same flags in each RunE.
// ---------------------------------------------------------------------------

// addTenantFlag registers the shared, required --tenant flag.
func addTenantFlag(cmd *cobra.Command) {
	cmd.Flags().Int("tenant", 0, "tenant ID (required)")
	_ = cmd.MarkFlagRequired("tenant")
}

// addPagingFlags registers shared pagination flags.
func addPagingFlags(cmd *cobra.Command) {
	cmd.Flags().Int("page-size", config.DefaultPageSize, "page size")
	cmd.Flags().String("page-token", "", "page token for a specific page")
	cmd.Flags().Bool("page-all", false, "fetch and aggregate every page")
}

// addOrderFlags registers shared ordering flags with a per-command default field.
func addOrderFlags(cmd *cobra.Command, defaultOrderBy string) {
	cmd.Flags().String("order-by", defaultOrderBy, "sort field")
	cmd.Flags().Bool("order-desc", true, "sort descending")
}

// addDateRangeFlags registers shared stat-date range flags.
func addDateRangeFlags(cmd *cobra.Command) {
	cmd.Flags().String("begin", "", "stat date begin (YYYY-MM-DD)")
	cmd.Flags().String("end", "", "stat date end (YYYY-MM-DD)")
	cmd.Flags().String("range", "", "relative range like 7d/4w/1m (overrides --begin/--end)")
}

// addJQFlag registers the shared --jq output filter flag.
func addJQFlag(cmd *cobra.Command) {
	cmd.Flags().String("jq", "", "filter JSON output with a jq expression")
}

// ---------------------------------------------------------------------------
// Flag resolution helpers
// ---------------------------------------------------------------------------

// requireTenant reads and validates the --tenant flag.
func requireTenant(cmd *cobra.Command) (int, error) {
	tenant, _ := cmd.Flags().GetInt("tenant")
	if tenant <= 0 {
		return 0, errs.NewValidationError(errs.SubtypeInvalidArgument,
			"--tenant must be a positive integer").WithParam("--tenant")
	}
	return tenant, nil
}

// applyPaging writes the page_size query param from flags. The page token is
// owned by runList so it can drive --page-all aggregation.
func applyPaging(cmd *cobra.Command, params map[string]interface{}) {
	if pageSize, _ := cmd.Flags().GetInt("page-size"); pageSize > 0 {
		params["page_size"] = pageSize
	}
}

// applyOrder writes order_by / order_desc query params from flags.
func applyOrder(cmd *cobra.Command, params map[string]interface{}) {
	if orderBy, _ := cmd.Flags().GetString("order-by"); orderBy != "" {
		params["order_by"] = orderBy
	}
	orderDesc, _ := cmd.Flags().GetBool("order-desc")
	params["order_desc"] = orderDesc
}

// applyDateRange resolves and writes stat_date_begin / stat_date_end.
// When required is true, both bounds must resolve to a non-empty value.
func applyDateRange(cmd *cobra.Command, params map[string]interface{}, required bool) error {
	rangeStr, _ := cmd.Flags().GetString("range")
	begin, _ := cmd.Flags().GetString("begin")
	end, _ := cmd.Flags().GetString("end")

	b, e, err := daterange.Resolve(rangeStr, begin, end)
	if err != nil {
		return err
	}
	if required && (b == "" || e == "") {
		return errs.NewValidationError(errs.SubtypeInvalidArgument,
			"a date range is required").
			WithHint("pass --range 7d, or both --begin and --end (YYYY-MM-DD)")
	}
	if b != "" {
		params["stat_date_begin"] = b
	}
	if e != "" {
		params["stat_date_end"] = e
	}
	return nil
}

// setString adds a query param only when the value is non-empty.
func setString(params map[string]interface{}, key, val string) {
	if val != "" {
		params[key] = val
	}
}

// setInt adds a query param only when the value differs from skip.
func setInt(params map[string]interface{}, key string, val, skip int) {
	if val != skip {
		params[key] = val
	}
}

// ---------------------------------------------------------------------------
// Request execution
// ---------------------------------------------------------------------------

// dryRun prints the request without executing it, returning true when the
// --dry-run flag is set.
func (f *Factory) dryRun(cmd *cobra.Command, path string, params map[string]interface{}) bool {
	on, _ := cmd.Flags().GetBool("dry-run")
	if !on {
		return false
	}
	fmt.Fprintf(f.ErrOut, "DRY RUN: GET %s%s\n", f.Config.BaseURL, path)
	pretty, _ := json.MarshalIndent(params, "", "  ")
	fmt.Fprintf(f.ErrOut, "Params:\n%s\n", string(pretty))
	return true
}

// output renders a result honoring --jq first, otherwise --format.
func (f *Factory) output(cmd *cobra.Command, result interface{}, columns []string) error {
	if jqExpr, _ := cmd.Flags().GetString("jq"); jqExpr != "" {
		return output.JqFilter(f.Out, result, jqExpr)
	}
	format := output.ParseFormat(f.resolveFormat(cmd))
	return output.Print(f.Out, result, format, columns)
}

// runList executes a list/report query: it validates jq, handles --dry-run,
// aggregates all pages when --page-all is set, and renders the result.
func (f *Factory) runList(cmd *cobra.Command, path string, params map[string]interface{}, columns []string) error {
	if jqExpr, _ := cmd.Flags().GetString("jq"); jqExpr != "" {
		if err := output.ValidateJqExpression(jqExpr); err != nil {
			return err
		}
	}
	if f.dryRun(cmd, path, params) {
		return nil
	}

	c := f.resolveClient(cmd)
	ctx := context.Background()

	fetch := func(token string) (map[string]interface{}, error) {
		p := cloneParams(params)
		if token != "" {
			p["page_token"] = token
		}
		data, err := c.Do(ctx, client.Request{Method: "GET", Path: path, Params: p})
		if err != nil {
			return nil, err
		}
		var m map[string]interface{}
		if err := json.Unmarshal(data, &m); err != nil {
			return nil, errs.NewInternalError(errs.SubtypeUnknown,
				"failed to parse response JSON: %v", err).WithCause(err)
		}
		return m, nil
	}

	pageAll, _ := cmd.Flags().GetBool("page-all")
	var result map[string]interface{}
	var err error
	if pageAll {
		result, err = paginate.All(fetch)
	} else {
		token, _ := cmd.Flags().GetString("page-token")
		result, err = fetch(token)
	}
	if err != nil {
		return err
	}
	return f.output(cmd, result, columns)
}

// runSingle executes a single-object query (get/dashboard). Table columns do
// not apply, so it renders JSON/pretty or a jq projection.
func (f *Factory) runSingle(cmd *cobra.Command, path string, params map[string]interface{}, columns []string) error {
	if jqExpr, _ := cmd.Flags().GetString("jq"); jqExpr != "" {
		if err := output.ValidateJqExpression(jqExpr); err != nil {
			return err
		}
	}
	if f.dryRun(cmd, path, params) {
		return nil
	}

	c := f.resolveClient(cmd)
	data, err := c.Do(context.Background(), client.Request{Method: "GET", Path: path, Params: params})
	if err != nil {
		return err
	}
	var result interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		fmt.Fprintln(f.Out, string(data))
		return nil
	}
	return f.output(cmd, result, columns)
}

// cloneParams shallow-copies a params map so per-page mutation is isolated.
func cloneParams(in map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(in)+1)
	for k, v := range in {
		out[k] = v
	}
	return out
}
