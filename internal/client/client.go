package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gmvstudio/adex-cli/errs"
)

// Client is the HTTP client for ADEX API calls.
type Client struct {
	BaseURL   string
	HTTP      *http.Client
	ErrOut    io.Writer
	APIKey    string
	UserAgent string
}

// Option configures a Client.
type Option func(*Client)

func WithAPIKey(key string) Option {
	return func(c *Client) { c.APIKey = key }
}

func WithUserAgent(ua string) Option {
	return func(c *Client) { c.UserAgent = ua }
}

func WithErrOut(w io.Writer) Option {
	return func(c *Client) { c.ErrOut = w }
}

// New creates a new API client.
func New(baseURL string, opts ...Option) *Client {
	c := &Client{
		BaseURL: strings.TrimRight(baseURL, "/"),
		HTTP: &http.Client{
			Timeout: 30 * time.Second,
		},
		UserAgent: "adex-cli/1.0",
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Request describes an API request.
type Request struct {
	Method string
	Path   string
	Params map[string]interface{}
	Body   interface{}
}

// Do executes an API request and returns the raw response body.
// On non-2xx responses, it returns a typed *errs.APIError.
func (c *Client) Do(ctx context.Context, req Request) ([]byte, error) {
	fullURL := c.BaseURL + req.Path
	if len(req.Params) > 0 {
		q := url.Values{}
		for k, v := range req.Params {
			switch val := v.(type) {
			case string:
				q.Set(k, val)
			case int:
				q.Set(k, strconv.Itoa(val))
			case int64:
				q.Set(k, strconv.FormatInt(val, 10))
			case bool:
				q.Set(k, strconv.FormatBool(val))
			case nil:
				continue
			default:
				q.Set(k, fmt.Sprintf("%v", v))
			}
		}
		fullURL += "?" + q.Encode()
	}

	var bodyReader io.Reader
	if req.Body != nil {
		b, err := json.Marshal(req.Body)
		if err != nil {
			return nil, errs.NewInternalError(errs.SubtypeUnknown, "failed to marshal request body: %v", err).WithCause(err)
		}
		bodyReader = strings.NewReader(string(b))
	}

	httpReq, err := http.NewRequestWithContext(ctx, req.Method, fullURL, bodyReader)
	if err != nil {
		return nil, errs.NewNetworkError(errs.SubtypeNetworkTransport, "failed to build request: %v", err).WithCause(err)
	}
	if req.Body != nil {
		httpReq.Header.Set("Content-Type", "application/json")
	}
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("User-Agent", c.UserAgent)
	if c.APIKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)
	}

	resp, err := c.HTTP.Do(httpReq)
	if err != nil {
		return nil, errs.NewNetworkError(errs.SubtypeNetworkTransport, "request to %s failed: %v", fullURL, err).WithCause(err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errs.NewNetworkError(errs.SubtypeNetworkTransport, "failed to read response: %v", err).WithCause(err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiMsg string
		var raw map[string]interface{}
		if json.Unmarshal(data, &raw) == nil {
			if msg, ok := raw["message"].(string); ok {
				apiMsg = msg
			} else if msg, ok := raw["error"].(string); ok {
				apiMsg = msg
			}
		}
		if apiMsg == "" {
			apiMsg = string(data)
		}

		if resp.StatusCode == 401 || resp.StatusCode == 403 {
			return nil, errs.NewAuthError(errs.SubtypeAuthRequired, "API returned %d: %s", resp.StatusCode, apiMsg).
				WithHint("check your API key or token configuration")
		}

		return nil, errs.NewAPIError(resp.StatusCode, "API returned %d: %s", resp.StatusCode, apiMsg)
	}

	return data, nil
}

// DoTyped executes a request and unmarshals the response into dest.
func (c *Client) DoTyped(ctx context.Context, req Request, dest interface{}) error {
	data, err := c.Do(ctx, req)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return errs.NewInternalError(errs.SubtypeUnknown, "failed to parse response JSON: %v", err).WithCause(err)
	}
	return nil
}
