package testmtx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/takuoki/testmtx/v2"
)

func TestExcelDoc(t *testing.T) {

	// test doc
	_, err := testmtx.NewExcelDoc("testdata/dummy.xlsx")
	if assert.NotNil(t, err) {
		assert.Regexp(t, "^fail to open excel file: ", err.Error())
	}

	doc, err := testmtx.NewExcelDoc("testdata/sample.xlsx")
	if !assert.Nil(t, err) {
		assert.Fail(t, "fail to create excel doc")
	}

	sheetNames, err := doc.GetSheetNames()
	if !assert.Nil(t, err) {
		assert.Fail(t, "fail to get sheet names")
	}
	assert.Equal(t, []string{"sample"}, sheetNames)

	// test sheet
	_, err = doc.GetSheet("dummy")
	if assert.NotNil(t, err) {
		assert.Equal(t, `sheet not found (name="dummy")`, err.Error())
	}

	sheet, err := doc.GetSheet("sample")
	if !assert.Nil(t, err) {
		assert.Fail(t, "fail to get sheet")
	}
	assert.Equal(t, "sample", sheet.Name())
	assert.Equal(t, "Properties", sheet.Value(2, 1)) // B3
	assert.Equal(t, "", sheet.Value(2, 100))

	// test row
	rows := sheet.Rows()
	assert.Equal(t, 3, rows[2].Number())
	assert.Equal(t, "Data", rows[1].Value(12)) // M2
	assert.Equal(t, "", rows[1].Value(100))
}
