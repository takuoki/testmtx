package testmtx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/takuoki/testmtx/v2"
	"github.com/tenntenn/golden"
)

func TestOutputter_Output(t *testing.T) {
	t.Parallel()

	// set to true when updating the golden file
	flagUpdate := false

	cases := map[string]struct {
		o      testmtx.Outputter
		sheet  *testmtx.Sheet
		golden string
	}{
		"oneColumnOneCaseOutputter: JSON": {
			o: func() testmtx.Outputter {
				f, _ := testmtx.NewJSONFormatter()
				return testmtx.NewOneColumnOneCaseOutputter(f)
			}(),
			sheet:  sampleParsedSheet(),
			golden: "oneColumnOneCase_JSON",
		},
		"oneSheetOneCaseOutputter: JSON": {
			o: func() testmtx.Outputter {
				f, _ := testmtx.NewJSONFormatter()
				return testmtx.NewOneSheetOneCaseOutputter(f)
			}(),
			sheet:  sampleParsedSheet(),
			golden: "oneSheetOneCase_JSON",
		},
		"oneColumnOneCaseOutputter: JSON (include empty)": {
			o: func() testmtx.Outputter {
				f, _ := testmtx.NewJSONFormatter()
				return testmtx.NewOneColumnOneCaseOutputter(f)
			}(),
			sheet: func() *testmtx.Sheet {
				s := sampleParsedSheet()
				c := s.Collections["want"].(*testmtx.ObjectCollection)
				c.ImplicitNils["case4"] = true
				return s
			}(),
			golden: "oneColumnOneCase_JSON_empty",
		},
		"oneSheetOneCaseOutputter: JSON (include empty)": {
			o: func() testmtx.Outputter {
				f, _ := testmtx.NewJSONFormatter()
				return testmtx.NewOneSheetOneCaseOutputter(f)
			}(),
			sheet: func() *testmtx.Sheet {
				s := sampleParsedSheet()
				c := s.Collections["want"].(*testmtx.ObjectCollection)
				c.ImplicitNils["case4"] = true
				return s
			}(),
			golden: "oneSheetOneCase_JSON_empty",
		},
	}

	for name, tt := range cases {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dir := t.TempDir()
			err := tt.o.Output(dir, tt.sheet)

			if assert.Nil(t, err) {
				got := golden.Txtar(t, dir)
				if diff := golden.Check(t, flagUpdate, "testdata/golden", tt.golden, got); diff != "" {
					t.Error(diff)
				}
			}
		})
	}
}
