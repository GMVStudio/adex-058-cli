package paginate

import "testing"

func TestAllAggregatesPages(t *testing.T) {
	pages := []map[string]interface{}{
		{"items": []interface{}{1, 2}, "hasMore": true, "nextPageToken": "t1"},
		{"items": []interface{}{3, 4}, "hasMore": true, "nextPageToken": "t2"},
		{"items": []interface{}{5}, "hasMore": false, "nextPageToken": ""},
	}
	var seenTokens []string
	i := 0
	fn := func(token string) (map[string]interface{}, error) {
		seenTokens = append(seenTokens, token)
		p := pages[i]
		i++
		return p, nil
	}

	result, err := All(fn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := result["total"].(int); got != 5 {
		t.Errorf("total = %d, want 5", got)
	}
	if got := len(result["items"].([]interface{})); got != 5 {
		t.Errorf("items len = %d, want 5", got)
	}
	if result["hasMore"].(bool) {
		t.Errorf("aggregated result should have hasMore=false")
	}
	want := []string{"", "t1", "t2"}
	for i, w := range want {
		if seenTokens[i] != w {
			t.Errorf("token[%d] = %q, want %q", i, seenTokens[i], w)
		}
	}
}

func TestAllSinglePage(t *testing.T) {
	fn := func(string) (map[string]interface{}, error) {
		return map[string]interface{}{"items": []interface{}{1}}, nil
	}
	result, err := All(fn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["total"].(int) != 1 {
		t.Errorf("total = %v, want 1", result["total"])
	}
}

func TestAllStopsWhenTokenEmptyEvenIfHasMore(t *testing.T) {
	fn := func(string) (map[string]interface{}, error) {
		return map[string]interface{}{"items": []interface{}{1}, "hasMore": true, "nextPageToken": ""}, nil
	}
	if _, err := All(fn); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
