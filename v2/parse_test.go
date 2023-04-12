package testmtx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/takuoki/testmtx/v2"
)

func TestParser_Parse(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		sheet   testmtx.DocSheet
		want    *testmtx.Sheet
		wantErr string
	}{
		"success": {
			sheet: sampleDocSheet(),
			want:  sampleParsedSheet(),
		},
		"failure: empty first column": {
			sheet:   sampleDocSheet().modify("H", 3, ""),
			wantErr: `first column name is empty (sheet="sample", cell="H3")`,
		},
	}

	parser, err := testmtx.NewParser(testmtx.PropLevel(5))
	if err != nil {
		t.Fatalf("fail to create parser: %v", err)
	}

	for name, tc := range testcases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			res, err := parser.Parse(tc.sheet)

			if tc.wantErr == "" {
				if assert.Nil(t, err) {
					assert.Equal(t, tc.want, res)
				}
			} else {
				if assert.NotNil(t, err) {
					assert.Equal(t, tc.wantErr, err.Error())
				}
			}
		})
	}
}
