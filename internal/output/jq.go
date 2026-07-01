package output

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"

	"github.com/itchyny/gojq"

	"github.com/gmvstudio/adex-cli/errs"
)

// ValidateJqExpression checks whether a jq expression is syntactically valid.
func ValidateJqExpression(expr string) error {
	query, err := gojq.Parse(expr)
	if err != nil {
		return errs.NewValidationError(errs.SubtypeInvalidArgument, "invalid jq expression: %s", err).WithCause(err)
	}
	if _, err := gojq.Compile(query); err != nil {
		return errs.NewValidationError(errs.SubtypeInvalidArgument, "invalid jq expression: %s", err).WithCause(err)
	}
	return nil
}

// JqFilter applies a jq expression to data and writes the results to w.
// Scalar values are printed raw (no quotes for strings), matching jq -r behavior.
// Complex values (maps, arrays) are printed as indented JSON.
func JqFilter(w io.Writer, data interface{}, expr string) error {
	query, err := gojq.Parse(expr)
	if err != nil {
		return errs.NewValidationError(errs.SubtypeInvalidArgument, "invalid jq expression: %s", err).WithCause(err)
	}
	code, err := gojq.Compile(query)
	if err != nil {
		return errs.NewValidationError(errs.SubtypeInvalidArgument, "invalid jq expression: %s", err).WithCause(err)
	}

	normalized := convertNumbers(toGeneric(data))

	iter := code.Run(normalized)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, isErr := v.(error); isErr {
			return errs.NewValidationError(errs.SubtypeInvalidArgument, "jq error: %s", err).WithCause(err)
		}
		if err := writeJqValue(w, v); err != nil {
			return err
		}
	}
	return nil
}

// writeJqValue writes a single jq result value to w. Scalars are printed raw;
// complex values as indented JSON.
func writeJqValue(w io.Writer, v interface{}) error {
	switch val := v.(type) {
	case nil:
		fmt.Fprintln(w, "null")
	case bool:
		fmt.Fprintln(w, val)
	case int:
		fmt.Fprintln(w, val)
	case float64:
		fmt.Fprintf(w, "%g\n", val)
	case *big.Int:
		fmt.Fprintln(w, val.String())
	case string:
		fmt.Fprintln(w, val)
	default:
		b, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			return errs.NewInternalError(errs.SubtypeUnknown, "failed to marshal jq result: %s", err).WithCause(err)
		}
		fmt.Fprintln(w, string(b))
	}
	return nil
}

// toGeneric round-trips data through JSON so typed structs become
// map[string]interface{} / []interface{} that gojq can process.
func toGeneric(data interface{}) interface{} {
	switch data.(type) {
	case map[string]interface{}, []interface{}:
		return data
	}
	b, err := json.Marshal(data)
	if err != nil {
		return data
	}
	var out interface{}
	if err := json.Unmarshal(b, &out); err != nil {
		return data
	}
	return out
}

// convertNumbers recursively converts json.Number values to int or float64
// so that gojq can process them correctly.
func convertNumbers(v interface{}) interface{} {
	switch val := v.(type) {
	case json.Number:
		if i, err := val.Int64(); err == nil {
			return int(i)
		}
		if f, err := val.Float64(); err == nil {
			return f
		}
		return val.String()
	case map[string]interface{}:
		for k, elem := range val {
			val[k] = convertNumbers(elem)
		}
		return val
	case []interface{}:
		for i, elem := range val {
			val[i] = convertNumbers(elem)
		}
		return val
	default:
		return v
	}
}
