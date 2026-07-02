package output

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"unicode/utf8"
)

// Format is the output format enum.
type Format string

const (
	FormatJSON   Format = "json"
	FormatPretty Format = "pretty"
	FormatTable  Format = "table"
)

// ParseFormat converts a string to Format, defaulting to JSON.
func ParseFormat(s string) Format {
	switch strings.ToLower(s) {
	case "pretty":
		return FormatPretty
	case "table":
		return FormatTable
	default:
		return FormatJSON
	}
}

// Print writes data to w in the specified format. For table format, columns
// specifies which keys to extract from each item, and the summary line (total /
// page) is written to errOut so it never corrupts the stdout data stream.
// A nil errOut suppresses the summary line.
func Print(w io.Writer, errOut io.Writer, data interface{}, format Format, columns []string) error {
	switch format {
	case FormatPretty:
		return printPretty(w, data)
	case FormatTable:
		return printTable(w, errOut, data, columns)
	default:
		return printJSON(w, data)
	}
}

func printJSON(w io.Writer, data interface{}) error {
	InjectNotice(data)
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, string(b))
	return err
}

func printPretty(w io.Writer, data interface{}) error {
	InjectNotice(data)
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, string(b))
	return err
}

// printTable renders a list response as a text table.
// It expects data to be a map with "items" being a slice of maps. The summary
// line goes to errOut so stdout stays a clean, pipeable table.
func printTable(w io.Writer, errOut io.Writer, data interface{}, columns []string) error {
	raw, ok := data.(map[string]interface{})
	if !ok {
		return printPretty(w, data)
	}

	items, ok := raw["items"].([]interface{})
	if !ok {
		return printPretty(w, data)
	}

	if len(columns) == 0 {
		columns = autoColumns(items)
	}

	// Print summary line to errOut (stderr) so it never pollutes stdout data.
	total, _ := raw["total"].(string)
	page, _ := raw["page"].(float64)
	pageSize, _ := raw["pageSize"].(float64)
	if total != "" && errOut != nil {
		fmt.Fprintf(errOut, "total: %s, page: %.0f, page_size: %.0f\n", total, page, pageSize)
	}

	// Calculate column widths
	widths := make([]int, len(columns))
	for i, col := range columns {
		widths[i] = utf8.RuneCountInString(strings.ToUpper(col))
	}
	rows := make([][]string, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		row := make([]string, len(columns))
		for i, col := range columns {
			val := extractField(m, col)
			row[i] = val
			if w := utf8.RuneCountInString(val); w > widths[i] {
				widths[i] = w
			}
		}
		rows = append(rows, row)
	}

	// Print header
	for i, col := range columns {
		if i > 0 {
			fmt.Fprint(w, "  ")
		}
		fmt.Fprintf(w, "%-*s", widths[i], strings.ToUpper(col))
	}
	fmt.Fprintln(w)

	// Print separator
	for i, col := range columns {
		if i > 0 {
			fmt.Fprint(w, "  ")
		}
		fmt.Fprintf(w, "%s", strings.Repeat("-", widths[i]))
		_ = col
	}
	fmt.Fprintln(w)

	// Print rows
	for _, row := range rows {
		for i, val := range row {
			if i > 0 {
				fmt.Fprint(w, "  ")
			}
			fmt.Fprintf(w, "%-*s", widths[i], val)
		}
		fmt.Fprintln(w)
	}

	return nil
}

// PendingNotice, if set, returns system-level notices to inject as the
// "_notice" field in JSON output envelopes. Set by cmd/root.go.
var PendingNotice func() map[string]interface{}

// GetNotice returns the current pending notice, or nil.
func GetNotice() map[string]interface{} {
	if PendingNotice == nil {
		return nil
	}
	return PendingNotice()
}

// injectNotice adds a "_notice" field into map[string]interface{} data that
// has an "ok" key (envelope-style responses). Non-map data or maps without
// "ok" are left untouched.
func InjectNotice(data interface{}) {
	if PendingNotice == nil {
		return
	}
	m, ok := data.(map[string]interface{})
	if !ok {
		return
	}
	if _, isEnvelope := m["ok"]; !isEnvelope {
		return
	}
	notice := PendingNotice()
	if notice == nil {
		return
	}
	m["_notice"] = notice
}

func autoColumns(items []interface{}) []string {
	if len(items) == 0 {
		return []string{}
	}
	m, ok := items[0].(map[string]interface{})
	if !ok {
		return []string{}
	}
	cols := make([]string, 0, len(m))
	for k := range m {
		cols = append(cols, k)
	}
	sort.Strings(cols)
	return cols
}

// extractField resolves a column key from a map, supporting dot notation
// for nested fields (e.g. "metrics.charge").
func extractField(m map[string]interface{}, key string) string {
	parts := strings.SplitN(key, ".", 2)
	val, ok := m[parts[0]]
	if !ok {
		return ""
	}
	if len(parts) == 1 {
		return fmt.Sprintf("%v", val)
	}
	nested, ok := val.(map[string]interface{})
	if !ok {
		return ""
	}
	return extractField(nested, parts[1])
}
