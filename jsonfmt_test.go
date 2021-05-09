package testmtx_test

import (
	"bytes"
	"testing"

	"github.com/takuoki/testmtx"
)

func TestJSONFormatFprint(t *testing.T) {
	testcase := map[string]struct {
		data     [][]string
		expected string
	}{
		"object": {
			data: [][]string{
				{"", "", "", "", "casename"},
				{"root", "", "", "object", "*new"},
				{"", "key1", "", "number", "1"},
				{"", "key2", "", "bool", "true"},
			},
			expected: "{\n  \"key1\": 1,\n  \"key2\": true\n}",
		},
		"array": {
			data: [][]string{
				{"", "", "", "", "casename"},
				{"root", "", "", "array", "*new"},
				{"", "*", "", "number", "1"},
				{"", "*", "", "bool", "true"},
			},
			expected: "[\n  1,\n  true\n]",
		},
		"string": {
			data: [][]string{
				{"", "", "", "", "casename"},
				{"root", "", "", "string", "abc"},
			},
			expected: `"abc"`,
		},
		"number": {
			data: [][]string{
				{"", "", "", "", "casename"},
				{"root", "", "", "number", "1"},
			},
			expected: "1",
		},
		"bool": {
			data: [][]string{
				{"", "", "", "", "casename"},
				{"root", "", "", "bool", "true"},
			},
			expected: "true",
		},
	}

	p, err := testmtx.NewParser(testmtx.PropLevel(3))
	if err != nil {
		t.Fatalf("unable to create new parser: %v", err)
	}
	f, err := testmtx.NewJSONFormatter()
	if err != nil {
		t.Fatalf("unable to create new formatter: %v", err)
	}
	for name, c := range testcase {
		t.Run(name, func(t *testing.T) {
			s, err := p.Parse(newTestSheet(c.data), "test_sheet")
			if err != nil {
				t.Fatalf("fail to parse sheet: %v", err)
			}
			buf := &bytes.Buffer{}
			f.Fprint(t, buf, s, 0)
			if buf.String() != c.expected {
				t.Errorf("print string doesn't match expected (expected=%q, actual=%q)", c.expected, buf.String())
			}
		})
	}
}

func TestJSONFormatExtension(t *testing.T) {
	f, err := testmtx.NewJSONFormatter()
	if err != nil {
		t.Fatalf("unable to create new formatter: %v", err)
	}
	result := f.Extension()
	if result != "json" {
		t.Errorf("extension doesn't match expected (expected=%q, actual=%q)", "json", result)
	}
}
