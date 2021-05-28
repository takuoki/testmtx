package testmtx_test

import (
	"bytes"
	"testing"

	"github.com/takuoki/testmtx"
)

func TestYamlFormatFprint(t *testing.T) {
	testcase := map[string]struct {
		data     [][]string
		expected string
	}{
		"mix": {
			data: [][]string{
				{"", "", "", "", "", "casename"},
				{"root", "", "", "", "object", "*new"},
				{"", "key1", "", "", "number", "1"},
				{"", "key2", "", "", "array", "*new"},
				{"", "", "*", "", "bool", "true"},
				{"", "", "*", "", "array", "*new"},
				{"", "", "", "*", "string", "abc"},
				{"", "", "", "*", "string", "xyz"},
				{"", "key3", "", "", "array", "*new"},
				{"", "", "*", "", "object", "*new"},
				{"", "", "", "key3-1-1", "string", "abc"},
				{"", "key4", "", "", "array", "*new"},
				{"", "", "*", "", "string", ""},
				{"", "", "*", "", "string", ""},
			},
			expected: "key1: 1\nkey2:\n  - true\n  -\n    - abc\n    - xyz\nkey3:\n  -\n    key3-1-1: abc\nkey4:\n",
		},
		"object": {
			data: [][]string{
				{"", "", "", "", "", "casename"},
				{"root", "", "", "", "object", "*new"},
				{"", "key1", "", "", "number", "1"},
				{"", "key2", "", "", "bool", "true"},
			},
			expected: "key1: 1\nkey2: true\n",
		},
		"array": {
			data: [][]string{
				{"", "", "", "", "", "casename"},
				{"root", "", "", "", "array", "*new"},
				{"", "*", "", "", "number", "1"},
				{"", "*", "", "", "bool", "true"},
			},
			expected: "- 1\n- true\n",
		},
		"string": {
			data: [][]string{
				{"", "", "", "", "", "casename"},
				{"root", "", "", "", "string", "abc"},
			},
			expected: "abc\n",
		},
		"number": {
			data: [][]string{
				{"", "", "", "", "", "casename"},
				{"root", "", "", "", "number", "1"},
			},
			expected: "1\n",
		},
		"bool": {
			data: [][]string{
				{"", "", "", "", "", "casename"},
				{"root", "", "", "", "bool", "true"},
			},
			expected: "true\n",
		},
	}

	p, err := testmtx.NewParser(testmtx.PropLevel(4))
	if err != nil {
		t.Fatalf("unable to create new parser: %v", err)
	}
	f, err := testmtx.NewYamlFormatter()
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
			f.Fprint(t, buf, s)
			if buf.String() != c.expected {
				t.Errorf("print string doesn't match expected (expected=%q, actual=%q)", c.expected, buf.String())
			}
		})
	}
}

func TestYamlFormatExtension(t *testing.T) {
	f, err := testmtx.NewYamlFormatter()
	if err != nil {
		t.Fatalf("unable to create new formatter: %v", err)
	}
	result := f.Extension()
	if result != "yaml" {
		t.Errorf("extension doesn't match expected (expected=%q, actual=%q)", "yaml", result)
	}
}
