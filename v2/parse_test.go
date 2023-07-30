package testmtx_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/takuoki/testmtx/v2"
)

func TestParser_Parse(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		sheet   testmtx.DocSheet
		want    *testmtx.Sheet
		wantErr string
	}{
		"success": {
			sheet: sampleDocSheet(),
			want:  sampleParsedSheet(),
		},
		"success: custom type": {
			sheet: sampleDocSheet().modify("C", 23, "time").modify("G", 23, "unixtime").modify("H", 23, "2023-01-23 12:34:56"),
			want: func() *testmtx.Sheet {
				sheet := sampleParsedSheet()
				want, _ := sheet.Collections["want"].(*testmtx.ObjectCollection)
				want.PropertyNames = append(want.PropertyNames, "time")
				want.Properties["time"] = &testmtx.SimpleCollection{
					ImplicitNils: map[testmtx.ColumnName]bool{
						"case2": true,
						"case3": true,
						"case4": true,
					},
					ExplicitNils: map[testmtx.ColumnName]bool{},
					Values: map[testmtx.ColumnName]testmtx.SimpleValue{
						"case1": &unixtimeValue{Value: "1674477296"},
					},
				}
				return sheet
			}(),
		},
		"failure: empty first column": {
			sheet:   sampleDocSheet().modify("H", 3, ""),
			wantErr: `first column name is empty (sheet="sample", cell="H3")`,
		},
		"failure: duplicated column name": {
			sheet:   sampleDocSheet().modify("I", 3, "case1"),
			wantErr: `column name ("case1") is duplicated (sheet="sample", cell="I3")`,
		},
		"failure: not belong root": {
			sheet:   sampleDocSheet().modify("B", 20, "").modify("C", 20, "want"),
			wantErr: `must not exist property that does not belong to the root property (sheet="sample", cell="C20")`,
		},
		"failure: duplicated root name": {
			sheet:   sampleDocSheet().modify("B", 20, "in"),
			wantErr: `root property name ("in") is duplicated (sheet="sample", cell="B20")`,
		},
		"failure: invalid object": {
			sheet:   sampleDocSheet().modify("H", 8, "invalid"),
			wantErr: `invalid object value ("invalid") (sheet="sample", cell="H8")`,
		},
		"failure: invalid level of object property": {
			sheet:   sampleDocSheet().modify("D", 9, "").modify("E", 9, "key1"),
			wantErr: `invalid level of object property (sheet="sample", cell="E9")`,
		},
		"failure: invalid array": {
			sheet:   sampleDocSheet().modify("H", 11, "invalid"),
			wantErr: `invalid array value ("invalid") (sheet="sample", cell="H11")`,
		},
		"failure: invalid level of array element": {
			sheet:   sampleDocSheet().modify("D", 12, "").modify("E", 12, "* 0"),
			wantErr: `invalid level of array element (sheet="sample", cell="E12")`,
		},
		"failure: invalid type": {
			sheet:   sampleDocSheet().modify("G", 5, "invalid"),
			wantErr: `invalid type ("invalid") (sheet="sample", cell="G5")`,
		},
		"failure: invalid number": {
			sheet:   sampleDocSheet().modify("H", 13, "invalid"),
			wantErr: `invalid number value ("invalid") (sheet="sample", cell="H13")`,
		},
		"failure: invalid bool": {
			sheet:   sampleDocSheet().modify("H", 7, "invalid"),
			wantErr: `invalid bool value ("invalid") (sheet="sample", cell="H7")`,
		},
	}

	parser, err := testmtx.NewParser(
		testmtx.PropLevel(5),
		testmtx.AdditionalSimpleValues(map[string]testmtx.ConvertValueFunc{
			"unixtime": convertUnixtimeValue,
		}),
	)
	if err != nil {
		t.Fatalf("fail to create parser: %v", err)
	}

	for name, tt := range cases {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := parser.Parse(tt.sheet)

			if tt.wantErr == "" {
				if assert.Nil(t, err) {
					assert.Equal(t, tt.want, got)
				} else {
					fmt.Println(err.Error())
				}
			} else {
				if assert.NotNil(t, err) {
					assert.Equal(t, tt.wantErr, err.Error())
				}
			}
		})
	}
}
