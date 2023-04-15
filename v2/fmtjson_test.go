package testmtx_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/takuoki/testmtx/v2"
)

func TestJSONFormatter_Write(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		col  testmtx.Collection
		cn   testmtx.ColumnName
		want string
	}{
		"case1: in": {
			col:  sampleParsedSheet().Collections["in"],
			cn:   "case1",
			want: "testdata/sample/case1/in.json",
		},
		"case1: want": {
			col:  sampleParsedSheet().Collections["want"],
			cn:   "case1",
			want: "testdata/sample/case1/want.json",
		},
		"case2: in": {
			col:  sampleParsedSheet().Collections["in"],
			cn:   "case2",
			want: "testdata/sample/case2/in.json",
		},
		"case2: want": {
			col:  sampleParsedSheet().Collections["want"],
			cn:   "case2",
			want: "testdata/sample/case2/want.json",
		},
		"case3: in": {
			col:  sampleParsedSheet().Collections["in"],
			cn:   "case3",
			want: "testdata/sample/case3/in.json",
		},
		"case3: want": {
			col:  sampleParsedSheet().Collections["want"],
			cn:   "case3",
			want: "testdata/sample/case3/want.json",
		},
		"case4: in": {
			col:  sampleParsedSheet().Collections["in"],
			cn:   "case4",
			want: "testdata/sample/case4/in.json",
		},
		"case4: want": {
			col:  sampleParsedSheet().Collections["want"],
			cn:   "case4",
			want: "testdata/sample/case4/want.json",
		},
	}

	f, err := testmtx.NewJSONFormatter()
	if err != nil {
		t.Fatalf("fail to create formatter: %v", err)
	}

	for name, tt := range cases {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			file, err := os.Open(tt.want)
			if err != nil {
				t.Fatalf("fail to open file (%q): %v", tt.want, err)
			}
			defer file.Close()

			wantStr, err := io.ReadAll(file)
			if err != nil {
				t.Fatalf("fail to read file (%q): %v", tt.want, err)
			}

			got := &bytes.Buffer{}
			f.Write(got, tt.col, tt.cn)

			assert.Equal(t, string(wantStr), got.String())
		})
	}
}
