package testmtx

import (
	"fmt"

	"github.com/tealeg/xlsx"
)

type excelDoc struct {
	file *xlsx.File
}

func NewExcelDoc(filepath string) (Doc, error) {
	file, err := xlsx.OpenFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("fail to open excel file: %w", err)
	}
	return &excelDoc{
		file: file,
	}, nil
}

func (d *excelDoc) GetSheetNames() ([]string, error) {
	ss := make([]string, 0, len(d.file.Sheets))
	for _, s := range d.file.Sheets {
		ss = append(ss, s.Name)
	}
	return ss, nil
}

func (d *excelDoc) GetSheet(sheetName string) (DocSheet, error) {
	if sh, ok := d.file.Sheet[sheetName]; ok {
		return &excelSheet{
			name:  sheetName,
			sheet: sh,
		}, nil
	}
	return nil, fmt.Errorf("sheet not found (name=%q)", sheetName)
}

type excelSheet struct {
	name  string
	sheet *xlsx.Sheet
}

func (s *excelSheet) Name() string {
	return s.name
}

func (s *excelSheet) Rows() []DocRow {
	rs := make([]DocRow, 0, len(s.sheet.Rows))
	for i, r := range s.sheet.Rows {
		rs = append(rs, &excelRow{number: i + 1, row: r})
	}
	return rs
}

func (s *excelSheet) Value(row, clm int) string {
	if len(s.sheet.Rows) < row || len(s.sheet.Rows[row].Cells) < clm {
		return ""
	}
	return s.sheet.Rows[row].Cells[clm].Value
}

type excelRow struct {
	number int
	row    *xlsx.Row
}

func (r *excelRow) Number() int {
	return r.number
}

func (r *excelRow) Value(clm int) string {
	if len(r.row.Cells) < clm {
		return ""
	}
	return r.row.Cells[clm].Value
}
