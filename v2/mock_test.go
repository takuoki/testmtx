package testmtx_test

import (
	"errors"

	"github.com/takuoki/clmconv"
	"github.com/takuoki/testmtx/v2"
)

type mockDoc struct {
	sheetNames []string
	sheets     map[string]testmtx.DocSheet
}

func (d *mockDoc) addSheet(name string, sheet testmtx.DocSheet) {
	d.sheetNames = append(d.sheetNames, name)
	d.sheets[name] = sheet
}

func (d *mockDoc) GetSheetNames() ([]string, error) {
	return d.sheetNames, nil
}

func (d *mockDoc) GetSheet(sheetName string) (testmtx.DocSheet, error) {
	s, ok := d.sheets[sheetName]
	if !ok {
		return nil, errors.New("sheet not found")
	}
	return s, nil
}

type mockDocSheet struct {
	name string
	rows [][]string
}

func newMockSheet(name string, rows [][]string) *mockDocSheet {
	return &mockDocSheet{name: name, rows: rows}
}

func (s *mockDocSheet) modify(clmLetter string, rowNumber int, value string) *mockDocSheet {
	clm := clmconv.MustAtoi(clmLetter)
	row := rowNumber - 1
	if row < len(s.rows) && clm < len(s.rows[row]) {
		s.rows[row][clm] = value
	}
	return s
}

func (s *mockDocSheet) Name() string {
	return s.name
}

func (s *mockDocSheet) Rows() []testmtx.DocRow {
	rs := make([]testmtx.DocRow, 0, len(s.rows))
	for i, row := range s.rows {
		rs = append(rs, &mockDocRow{number: i + 1, values: row})
	}
	return rs
}

func (s *mockDocSheet) Value(row, clm int) string {
	if row < len(s.rows) && clm < len(s.rows[row]) {
		return s.rows[row][clm]
	}
	return ""
}

type mockDocRow struct {
	number int
	values []string
}

func (r *mockDocRow) Number() int {
	return r.number
}

func (r *mockDocRow) Value(clm int) string {
	if clm < len(r.values) {
		return r.values[clm]
	}
	return ""
}
