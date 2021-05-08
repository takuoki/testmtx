package testmtx_test

import (
	"github.com/takuoki/gsheets/sheets"
)

type testSheet struct {
	data [][]string
}

func (s *testSheet) Rows() []sheets.Row {
	rs := []sheets.Row{}
	for _, r := range s.data {
		rs = append(rs, &testRow{data: r})
	}
	return rs
}

func (s *testSheet) Value(row, clm int) string {
	if s == nil || len(s.data) <= row || len(s.data[row]) <= clm {
		return ""
	}
	return s.data[row][clm]
}

type testRow struct {
	data []string
}

func (r *testRow) Value(clm int) string {
	if r == nil || len(r.data) <= clm {
		return ""
	}
	return r.data[clm]
}

// newTestSheet returns new sheet object.
// This appends first empty column and two empty rows.
func newTestSheet(data [][]string) sheets.Sheet {
	s := [][]string{{}, {}}
	for _, d := range data {
		s = append(s, append([]string{""}, d...))
	}
	return &testSheet{data: s}
}

// newSingleValueTestSheet returns new sheet object.
// The property level for this sheet is 10 by default.
func newSingleValueTestSheet(typ, value string) sheets.Sheet {
	return &testSheet{
		data: [][]string{
			{}, {},
			{"", "", "", "", "", "", "", "", "", "", "", "", "case1"},
			{"", "key", "", "", "", "", "", "", "", "", "", typ, value},
		},
	}
}
