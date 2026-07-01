package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/gmvstudio/adex-cli/internal/config"
)

// runResult captures the outcome of executing the root command in a test.
type runResult struct {
	Out string
	Err string
	// DryRunPath is the request path parsed from the "DRY RUN: GET ..." line.
	DryRunPath string
	// DryRunParams is the params object parsed from the dry-run output.
	DryRunParams map[string]interface{}
	// ExecErr is the error returned by root.Execute().
	ExecErr error
}

// newTestFactory builds a Factory with in-memory buffers and a fixed config,
// so command tests never touch the real filesystem or network.
func newTestFactory(out, errOut *bytes.Buffer) *Factory {
	return &Factory{
		Config: &config.Config{BaseURL: "http://test.local", Authorization: "adex_test"},
		Out:    out,
		ErrOut: errOut,
	}
}

// runCmd executes the root command with args and returns a parsed runResult.
// Callers should pass --dry-run to exercise request building without a network.
func runCmd(t *testing.T, args ...string) runResult {
	t.Helper()
	var out, errOut bytes.Buffer
	f := newTestFactory(&out, &errOut)
	root := NewRootCmd(f)
	root.SetOut(&out)
	root.SetErr(&errOut)
	root.SetArgs(args)

	res := runResult{ExecErr: root.Execute()}
	res.Out = out.String()
	res.Err = errOut.String()

	// Parse the dry-run envelope if present.
	for _, line := range strings.Split(res.Err, "\n") {
		if strings.HasPrefix(line, "DRY RUN: GET ") {
			full := strings.TrimPrefix(line, "DRY RUN: GET ")
			res.DryRunPath = strings.TrimPrefix(full, "http://test.local")
		}
	}
	if idx := strings.Index(res.Err, "Params:\n"); idx >= 0 {
		blob := res.Err[idx+len("Params:\n"):]
		dec := json.NewDecoder(strings.NewReader(blob))
		_ = dec.Decode(&res.DryRunParams)
	}
	return res
}

// paramString returns a param rendered as a string for assertions, tolerating
// the JSON number/bool decoding that happens through the dry-run envelope.
func (r runResult) paramString(key string) (string, bool) {
	v, ok := r.DryRunParams[key]
	if !ok {
		return "", false
	}
	switch t := v.(type) {
	case string:
		return t, true
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64), true
	case bool:
		return strconv.FormatBool(t), true
	default:
		return fmt.Sprintf("%v", t), true
	}
}
