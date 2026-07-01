package daterange

import (
	"testing"
	"time"
)

func withFixedNow(t *testing.T, y int, m time.Month, d int) {
	t.Helper()
	orig := now
	now = func() time.Time { return time.Date(y, m, d, 12, 0, 0, 0, time.UTC) }
	t.Cleanup(func() { now = orig })
}

func TestResolveRangeDays(t *testing.T) {
	withFixedNow(t, 2026, time.July, 1)
	begin, end, err := Resolve("7d", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if begin != "2026-06-25" {
		t.Errorf("begin = %q, want 2026-06-25", begin)
	}
	if end != "2026-07-01" {
		t.Errorf("end = %q, want 2026-07-01", end)
	}
}

func TestResolveRangeWeeksMonths(t *testing.T) {
	withFixedNow(t, 2026, time.July, 1)
	if b, _, _ := Resolve("1w", "", ""); b != "2026-06-25" {
		t.Errorf("1w begin = %q, want 2026-06-25", b)
	}
	if b, _, _ := Resolve("1m", "", ""); b != "2026-06-02" {
		t.Errorf("1m begin = %q, want 2026-06-02", b)
	}
}

func TestResolveRangePrecedesExplicit(t *testing.T) {
	withFixedNow(t, 2026, time.July, 1)
	begin, _, err := Resolve("1d", "2000-01-01", "2000-01-02")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if begin != "2026-07-01" {
		t.Errorf("range should override explicit begin, got %q", begin)
	}
}

func TestResolveExplicit(t *testing.T) {
	begin, end, err := Resolve("", "2026-06-01", "2026-06-30")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if begin != "2026-06-01" || end != "2026-06-30" {
		t.Errorf("got (%q, %q)", begin, end)
	}
}

func TestResolveEmpty(t *testing.T) {
	begin, end, err := Resolve("", "", "")
	if err != nil || begin != "" || end != "" {
		t.Errorf("empty inputs should yield empty result, got (%q,%q,%v)", begin, end, err)
	}
}

func TestResolveErrors(t *testing.T) {
	cases := []struct{ name, r, b, e string }{
		{"bad unit", "7y", "", ""},
		{"bad number", "xd", "", ""},
		{"zero", "0d", "", ""},
		{"bad begin", "", "2026-13-01", ""},
		{"begin after end", "", "2026-07-02", "2026-07-01"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, _, err := Resolve(c.r, c.b, c.e); err == nil {
				t.Errorf("expected error for %+v", c)
			}
		})
	}
}
