package client

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gmvstudio/adex-cli/errs"
)

func TestDoSuccessAndParams(t *testing.T) {
	var gotAuth, gotQuery, gotUA string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotUA = r.Header.Get("User-Agent")
		gotQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	c := New(srv.URL, WithAPIKey("secret"), WithUserAgent("adex-test/9"))
	data, err := c.Do(context.Background(), Request{
		Method: "GET",
		Path:   "/v1/x",
		Params: map[string]interface{}{"a": "1", "n": 2, "b": true, "skip": nil},
	})
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	if string(data) != `{"ok":true}` {
		t.Errorf("body = %q", data)
	}
	if gotAuth != "Bearer secret" {
		t.Errorf("Authorization = %q", gotAuth)
	}
	if gotUA != "adex-test/9" {
		t.Errorf("User-Agent = %q", gotUA)
	}
	// nil params must be omitted; int/bool must be serialized.
	if want := "a=1&b=true&n=2"; gotQuery != want {
		t.Errorf("query = %q, want %q", gotQuery, want)
	}
}

func TestDoAPIErrorExtractsMessage(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"message":"bad tenant"}`))
	}))
	defer srv.Close()

	c := New(srv.URL)
	_, err := c.Do(context.Background(), Request{Method: "GET", Path: "/x"})
	var apiErr *errs.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error type = %T, want *errs.APIError", err)
	}
	if apiErr.Code != http.StatusBadRequest {
		t.Errorf("code = %d, want 400", apiErr.Code)
	}
	if apiErr.Message == "" || apiErr.Message[len(apiErr.Message)-len("bad tenant"):] != "bad tenant" {
		t.Errorf("message = %q, want to end with 'bad tenant'", apiErr.Message)
	}
}

func TestDo401IsAuthError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"expired"}`))
	}))
	defer srv.Close()

	c := New(srv.URL)
	_, err := c.Do(context.Background(), Request{Method: "GET", Path: "/x"})
	var authErr *errs.AuthError
	if !errors.As(err, &authErr) {
		t.Fatalf("error type = %T, want *errs.AuthError", err)
	}
	if authErr.Subtype != errs.SubtypeAuthRequired {
		t.Errorf("subtype = %q", authErr.Subtype)
	}
}

func TestDoTypedUnmarshalError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`not json`))
	}))
	defer srv.Close()

	c := New(srv.URL)
	var dest map[string]interface{}
	err := c.DoTyped(context.Background(), Request{Method: "GET", Path: "/x"}, &dest)
	var internal *errs.InternalError
	if !errors.As(err, &internal) {
		t.Fatalf("error type = %T, want *errs.InternalError", err)
	}
}

func TestNewTrimsTrailingSlash(t *testing.T) {
	c := New("http://x/")
	if c.BaseURL != "http://x" {
		t.Errorf("BaseURL = %q, want http://x", c.BaseURL)
	}
}
