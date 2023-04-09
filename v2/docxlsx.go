package testmtx

import (
	"errors"
	"fmt"

	"github.com/tealeg/xlsx"
)

type xlsxDoc struct {
	file *xlsx.File
}

func NewXlsxDoc(filepath string) (Doc, error) {
	xlFile, err := xlsx.OpenFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("fail to open excel file: %w", err)
	}
	return &xlsxDoc{
		file: xlFile,
	}, nil
}

func (d *xlsxDoc) GetSheetNames() ([]string, error) {
	if d == nil {
		return nil, errors.New("xlsxDoc is not initialized")
	}
	ss := []string{}
	for _, s := range d.file.Sheets {
		ss = append(ss, s.Name)
	}
	return ss, nil
}

func (d *xlsxDoc) GetSheet(sheetName string) (DocSheet, error) {
	if d == nil {
		return nil, errors.New("xlsxDoc is not initialized")
	}
	return &xlsxSheet{sheet: d.file.Sheet[sheetName]}, nil
}

type xlsxSheet struct {
	sheet *xlsx.Sheet
}

func (s *xlsxSheet) Rows() []DocRow {
	if s == nil || s.sheet == nil {
		return nil
	}
	rs := []DocRow{}
	for _, r := range s.sheet.Rows {
		rs = append(rs, &xlsxRow{row: r})
	}
	return rs
}

func (s *xlsxSheet) Value(row, clm int) string {
	if s == nil || s.sheet == nil || len(s.sheet.Rows) < row || len(s.sheet.Rows[row].Cells) < clm {
		return ""
	}
	return s.sheet.Rows[row].Cells[clm].Value
}

type xlsxRow struct {
	row *xlsx.Row
}

func (r *xlsxRow) Value(clm int) string {
	if r == nil || r.row == nil || len(r.row.Cells) < clm {
		return ""
	}
	return r.row.Cells[clm].Value
}
