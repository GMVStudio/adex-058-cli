package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestParseFormat(t *testing.T) {
	cases := map[string]Format{
		"json":    FormatJSON,
		"pretty":  FormatPretty,
		"table":   FormatTable,
		"TABLE":   FormatTable,
		"unknown": FormatJSON,
		"":        FormatJSON,
	}
	for in, want := range cases {
		if got := ParseFormat(in); got != want {
			t.Errorf("ParseFormat(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestPrintTableSummaryGoesToErrOut(t *testing.T) {
	var out, errOut bytes.Buffer
	data := map[string]interface{}{
		"items": []interface{}{
			map[string]interface{}{"id": "1", "name": "a"},
			map[string]interface{}{"id": "2", "name": "b"},
		},
		"total":    "2",
		"page":     float64(1),
		"pageSize": float64(20),
	}
	if err := Print(&out, &errOut, data, FormatTable, []string{"id", "name"}); err != nil {
		t.Fatalf("Print: %v", err)
	}

	// stdout must contain only the table, no summary line.
	if strings.Contains(out.String(), "total:") {
		t.Errorf("summary leaked into stdout:\n%s", out.String())
	}
	if !strings.Contains(errOut.String(), "total: 2") {
		t.Errorf("summary missing from errOut:\n%s", errOut.String())
	}
	if !strings.Contains(out.String(), "ID") || !strings.Contains(out.String(), "NAME") {
		t.Errorf("table header missing:\n%s", out.String())
	}
}

func TestPrintTableNilErrOutSuppressesSummary(t *testing.T) {
	var out bytes.Buffer
	data := map[string]interface{}{
		"items": []interface{}{map[string]interface{}{"id": "1"}},
		"total": "1",
	}
	if err := Print(&out, nil, data, FormatTable, []string{"id"}); err != nil {
		t.Fatalf("Print: %v", err)
	}
	if strings.Contains(out.String(), "total:") {
		t.Errorf("summary leaked into stdout with nil errOut:\n%s", out.String())
	}
}

func TestPrintJSON(t *testing.T) {
	var out, errOut bytes.Buffer
	if err := Print(&out, &errOut, map[string]interface{}{"a": 1}, FormatJSON, nil); err != nil {
		t.Fatalf("Print: %v", err)
	}
	if !strings.Contains(out.String(), `"a":1`) {
		t.Errorf("json output = %q", out.String())
	}
}

func TestAutoColumnsDeterministic(t *testing.T) {
	items := []interface{}{
		map[string]interface{}{"c": 1, "a": 2, "b": 3},
	}
	got := autoColumns(items)
	want := []string{"a", "b", "c"}
	if len(got) != len(want) {
		t.Fatalf("autoColumns len = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("autoColumns[%d] = %q, want %q (should be sorted)", i, got[i], want[i])
		}
	}
}

func TestExtractFieldNested(t *testing.T) {
	m := map[string]interface{}{
		"metrics": map[string]interface{}{"charge": 12.5},
		"id":      "x",
	}
	if got := extractField(m, "metrics.charge"); got != "12.5" {
		t.Errorf("extractField nested = %q, want 12.5", got)
	}
	if got := extractField(m, "id"); got != "x" {
		t.Errorf("extractField = %q, want x", got)
	}
	if got := extractField(m, "missing"); got != "" {
		t.Errorf("extractField missing = %q, want empty", got)
	}
}
