// Package paginate aggregates token-paginated ADEX list responses.
//
// ADEX list replies share a common shape:
//
//	{ "items": [...], "hasMore": bool, "nextPageToken": "..." }
//
// PageFn fetches one page given a page token (empty for the first page).
package paginate

import "github.com/gmvstudio/adex-cli/errs"

// maxPages guards against a server that never stops returning hasMore=true.
const maxPages = 1000

// PageFn fetches a single page. pageToken is empty for the first page.
type PageFn func(pageToken string) (map[string]interface{}, error)

// meta extracts the pagination cursor from a decoded reply.
func meta(data map[string]interface{}) (hasMore bool, next string) {
	hasMore, _ = data["hasMore"].(bool)
	next, _ = data["nextPageToken"].(string)
	return
}

// items extracts the items slice from a decoded reply.
func items(data map[string]interface{}) []interface{} {
	if v, ok := data["items"].([]interface{}); ok {
		return v
	}
	return nil
}

// All follows nextPageToken until hasMore is false, concatenating every page's
// items into a single reply map. The returned map has hasMore=false, an empty
// nextPageToken, the aggregated items, and a "total" count.
func All(fn PageFn) (map[string]interface{}, error) {
	var aggregated []interface{}
	token := ""

	for page := 0; page < maxPages; page++ {
		data, err := fn(token)
		if err != nil {
			return nil, err
		}
		aggregated = append(aggregated, items(data)...)

		hasMore, next := meta(data)
		if !hasMore || next == "" {
			return map[string]interface{}{
				"items":         aggregated,
				"hasMore":       false,
				"nextPageToken": "",
				"total":         len(aggregated),
			}, nil
		}
		token = next
	}

	return nil, errs.NewInternalError(errs.SubtypeUnknown,
		"pagination exceeded %d pages; aborting to avoid an infinite loop", maxPages).
		WithHint("narrow the query with filters or a smaller date range")
}
