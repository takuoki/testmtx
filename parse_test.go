package testmtx_test

import (
	"strings"
	"testing"

	"github.com/takuoki/testmtx"
)

func TestParseError(t *testing.T) {

	t.Run("non-positive-prop-level", func(t *testing.T) {
		expectedError := "property level should be positive value"

		_, err := testmtx.NewParser(testmtx.PropLevel(0))

		if err == nil {
			t.Error("error must be occurred")
		}
		if !strings.Contains(err.Error(), expectedError) {
			t.Errorf("error message doesn't contain expected (expected=%q, err=%q)", expectedError, err)
		}
	})

	t.Run("non-initialized-parser", func(t *testing.T) {
		expectedError := "parser is not initilized"

		var p *testmtx.Parser
		_, err := p.Parse(newTestSheet([][]string{}), "test_sheet")

		if err == nil {
			t.Error("error must be occurred")
		}
		if !strings.Contains(err.Error(), expectedError) {
			t.Errorf("error message doesn't contain expected (expected=%q, err=%q)", expectedError, err)
		}
	})
}

func TestParseSheetCheck(t *testing.T) {
	testcase := map[string]struct {
		data          [][]string
		expectedError string
	}{
		"failure-invalid-sheet-format": {
			data:          [][]string{},
			expectedError: "invalid sheet format",
		},
		"failure-duplicated-case-name": {
			data: [][]string{
				{"", "", "", "", "casename", "casename"},
			},
			expectedError: "case name is duplicated",
		},
		"failure-root-level-first": {
			data: [][]string{
				{"", "", "", "", "casename"},
				{"", "key", "", "string", "abc"},
			},
			expectedError: "must not exist property that does not belong to the root property",
		},
		"failure-root-level-second": {
			data: [][]string{
				{"", "", "", "", "casename"},
				{"root", "", "", "object", "*new"},
				{"", "key", "", "string", "abc"},
				{"", "", "", "", ""},
				{"", "key", "", "string", "abc"},
			},
			expectedError: "must not exist property that does not belong to the root property",
		},
		"failure-duplicated-root-property-name": {
			data: [][]string{
				{"", "", "", "", "casename"},
				{"root", "", "", "string", "abc"},
				{"root", "", "", "string", "abc"},
			},
			expectedError: "root property name is duplicated",
		},
		"success-object-child": {
			data: [][]string{
				{"", "", "", "", "casename"},
				{"root", "", "", "object", "*new"},
				{"", "key", "", "string", "abc"},
			},
		},
		"failure-object-child-level": {
			data: [][]string{
				{"", "", "", "", "casename"},
				{"root", "", "", "object", "*new"},
				{"", "", "key", "string", "abc"},
			},
			expectedError: "invalid level of object child",
		},
		"failure-object-child-element": {
			data: [][]string{
				{"", "", "", "", "casename"},
				{"root", "", "", "object", "*new"},
				{"", "key", "", "number", "abc"},
			},
			expectedError: "unable to convert numeric value",
		},
		"success-object-child-duplicated-name": {
			data: [][]string{
				{"", "", "", "", "casename"},
				{"root", "", "", "object", "*new"},
				{"", "key", "", "string", "abc"},
				{"", "key", "", "string", "abc"},
			},
			// Duplicates are allowed because the property names are stored in an array.
			expectedError: "",
		},
		"success-array-child": {
			data: [][]string{
				{"", "", "", "", "casename"},
				{"root", "", "", "array", "*new"},
				{"", "key", "", "string", "abc"},
			},
		},
		"failure-array-child-level": {
			data: [][]string{
				{"", "", "", "", "casename"},
				{"root", "", "", "array", "*new"},
				{"", "", "key", "string", "abc"},
			},
			expectedError: "invalid level of array child",
		},
		"failure-array-child-element": {
			data: [][]string{
				{"", "", "", "", "casename"},
				{"root", "", "", "array", "*new"},
				{"", "key", "", "number", "abc"},
			},
			expectedError: "unable to convert numeric value",
		},
	}

	p, err := testmtx.NewParser(testmtx.PropLevel(3))
	if err != nil {
		t.Fatalf("unable to create new parser: %v", err)
	}
	for name, c := range testcase {
		t.Run(name, func(t *testing.T) {
			_, err := p.Parse(newTestSheet(c.data), "test_sheet")
			if err != nil {
				if c.expectedError == "" {
					t.Errorf("error must not be occurred (err=%q)", err)
				} else {
					if !strings.Contains(err.Error(), c.expectedError) {
						t.Errorf("error message doesn't contain expected (expected=%q, err=%q)", c.expectedError, err)
					}
				}
			} else {
				if c.expectedError != "" {
					t.Errorf("error must be occurred (expected=%q)", c.expectedError)
				}
			}
		})
	}
}

func TestParseValueCheck(t *testing.T) {
	testcase := map[string]struct {
		typ           string
		value         string
		expectedError string
	}{
		"success-object":         {typ: "object", value: "*new"},
		"success-object-nil":     {typ: "object", value: ""},
		"failure-object":         {typ: "object", value: "abc", expectedError: "unable to convert object or array value"},
		"success-array":          {typ: "array", value: "*new"},
		"success-array-nil":      {typ: "array", value: ""},
		"failure-array":          {typ: "array", value: "abc", expectedError: "unable to convert object or array value"},
		"success-string":         {typ: "string", value: "abc"},
		"success-string-empty":   {typ: "string", value: "*empty"},
		"success-string-nil":     {typ: "string", value: ""},
		"success-number-integer": {typ: "number", value: "1"},
		"success-number-float":   {typ: "number", value: "1.2"},
		"success-number-nil":     {typ: "number", value: ""},
		"failure-number":         {typ: "number", value: "abc", expectedError: "unable to convert numeric value"},
		"success-bool-true":      {typ: "bool", value: "true"},
		"success-bool-false":     {typ: "bool", value: "false"},
		"success-bool-nil":       {typ: "bool", value: ""},
		"failure-bool":           {typ: "bool", value: "error", expectedError: "unable to convert bool value"},
		"failure-invalid-type":   {typ: "dummy", value: "", expectedError: "invalid type"},
	}

	p, err := testmtx.NewParser()
	if err != nil {
		t.Fatalf("unable to create new parser: %v", err)
	}
	for name, c := range testcase {
		t.Run(name, func(t *testing.T) {
			_, err := p.Parse(newSingleValueTestSheet(c.typ, c.value), "test_sheet")
			if err != nil {
				if c.expectedError == "" {
					t.Errorf("error must not be occurred (err=%q)", err)
				} else {
					if !strings.Contains(err.Error(), c.expectedError) {
						t.Errorf("error message doesn't contain expected (expected=%q, err=%q)", c.expectedError, err)
					}
				}
			} else {
				if c.expectedError != "" {
					t.Errorf("error must be occurred (expected=%q)", c.expectedError)
				}
			}
		})
	}
}
