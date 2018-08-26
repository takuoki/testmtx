package sheet_test

import (
	"reflect"
	"testing"

	"github.com/takuoki/gostr"
	"github.com/takuoki/testmtx/sheet"
)

var (
	bTrue  = true
	bFalse = false
)

func ptrString(s string) *string {
	return &s
}

var (
	dataMap = map[string]sheet.Data{
		"request": &sheet.DObject{
			Values:        map[sheet.Casename]bool{"case1": true, "case2": true, "case_3": true},
			PropertyNames: []string{"num_key", "string_key", "bool_key", "object_key", "array_key"},
			Properties: map[string]sheet.Data{
				"num_key": &sheet.DNum{
					Values: map[sheet.Casename]*string{
						"case1":  ptrString("101"),
						"case2":  ptrString("102"),
						"case_3": ptrString("103"),
					},
				},
				"string_key": &sheet.DString{
					Values: map[sheet.Casename]*string{
						"case1":  ptrString("string value 101"),
						"case2":  ptrString("string value 102"),
						"case_3": nil,
					},
				},
				"bool_key": &sheet.DBool{
					Values: map[sheet.Casename]*bool{
						"case1":  &bTrue,
						"case2":  nil,
						"case_3": &bFalse,
					},
				},
				"object_key": &sheet.DObject{
					Values:        map[sheet.Casename]bool{"case1": true, "case2": true, "case_3": false},
					PropertyNames: []string{"key1", "key2"},
					Properties: map[string]sheet.Data{
						"key1": &sheet.DNum{
							Values: map[sheet.Casename]*string{
								"case1":  ptrString("201"),
								"case2":  ptrString("202"),
								"case_3": nil,
							},
						},
						"key2": &sheet.DString{
							Values: map[sheet.Casename]*string{
								"case1":  ptrString("string value 201"),
								"case2":  ptrString("string value 202"),
								"case_3": nil,
							},
						},
					},
				},
				"array_key": &sheet.DArray{
					Values: map[sheet.Casename]bool{"case1": true, "case2": false, "case_3": true},
					Elements: []sheet.Data{
						&sheet.DObject{
							Values:        map[sheet.Casename]bool{"case1": true, "case2": false, "case_3": true},
							PropertyNames: []string{"key3", "key4"},
							Properties: map[string]sheet.Data{
								"key3": &sheet.DNum{
									Values: map[sheet.Casename]*string{
										"case1":  ptrString("301"),
										"case2":  nil,
										"case_3": ptrString("303"),
									},
								},
								"key4": &sheet.DString{
									Values: map[sheet.Casename]*string{
										"case1":  ptrString("string value 301"),
										"case2":  nil,
										"case_3": ptrString("string value 303"),
									},
								},
							},
						},
						&sheet.DObject{
							Values:        map[sheet.Casename]bool{"case1": true, "case2": false, "case_3": false},
							PropertyNames: []string{"key3", "key4"},
							Properties: map[string]sheet.Data{
								"key3": &sheet.DNum{
									Values: map[sheet.Casename]*string{
										"case1":  ptrString("401"),
										"case2":  nil,
										"case_3": nil,
									},
								},
								"key4": &sheet.DString{
									Values: map[sheet.Casename]*string{
										"case1":  ptrString("string value 401"),
										"case2":  nil,
										"case_3": nil,
									},
								},
							},
						},
					},
				},
			},
		},
		"expected": &sheet.DObject{
			Values:        map[sheet.Casename]bool{"case1": true, "case2": true, "case_3": true},
			PropertyNames: []string{"status", "code"},
			Properties: map[string]sheet.Data{
				"status": &sheet.DString{
					Values: map[sheet.Casename]*string{
						"case1":  ptrString("success"),
						"case2":  ptrString("failure"),
						"case_3": ptrString("failure"),
					},
				},
				"code": &sheet.DNum{
					Values: map[sheet.Casename]*string{
						"case1":  ptrString("200"),
						"case2":  ptrString("401"),
						"case_3": ptrString("404"),
					},
				},
			},
		},
	}

	noCaseDataMap = map[string]sheet.Data{
		"request": &sheet.DObject{
			Values:        map[sheet.Casename]bool{},
			PropertyNames: []string{"num_key", "string_key", "bool_key", "object_key", "array_key"},
			Properties: map[string]sheet.Data{
				"num_key":    &sheet.DNum{Values: map[sheet.Casename]*string{}},
				"string_key": &sheet.DString{Values: map[sheet.Casename]*string{}},
				"bool_key":   &sheet.DBool{Values: map[sheet.Casename]*bool{}},
				"object_key": &sheet.DObject{
					Values:        map[sheet.Casename]bool{},
					PropertyNames: []string{"key1", "key2"},
					Properties: map[string]sheet.Data{
						"key1": &sheet.DNum{Values: map[sheet.Casename]*string{}},
						"key2": &sheet.DString{Values: map[sheet.Casename]*string{}},
					},
				},
				"array_key": &sheet.DArray{
					Values: map[sheet.Casename]bool{},
					Elements: []sheet.Data{
						&sheet.DObject{
							Values:        map[sheet.Casename]bool{},
							PropertyNames: []string{"key3", "key4"},
							Properties: map[string]sheet.Data{
								"key3": &sheet.DNum{Values: map[sheet.Casename]*string{}},
								"key4": &sheet.DString{Values: map[sheet.Casename]*string{}},
							},
						},
						&sheet.DObject{
							Values:        map[sheet.Casename]bool{},
							PropertyNames: []string{"key3", "key4"},
							Properties: map[string]sheet.Data{
								"key3": &sheet.DNum{Values: map[sheet.Casename]*string{}},
								"key4": &sheet.DString{Values: map[sheet.Casename]*string{}},
							},
						},
					},
				},
			},
		},
		"expected": &sheet.DObject{
			Values:        map[sheet.Casename]bool{},
			PropertyNames: []string{"status", "code"},
			Properties: map[string]sheet.Data{
				"status": &sheet.DString{Values: map[sheet.Casename]*string{}},
				"code":   &sheet.DNum{Values: map[sheet.Casename]*string{}},
			},
		},
	}
)

func TestGet(t *testing.T) {

	authFile := "credentials.json"
	spreadsheetID := "1C4BMTwvGVfiLpcGN2VDqsofalQtpORhpTlZtMxPjsr8"
	exceptSheetSet := map[string]struct{}{"overview": {}}

	expected := []*sheet.Sheet{
		{Name: "sheet", Cases: []sheet.Casename{"case1", "case2", "case_3"}, DataMap: dataMap},
		{Name: "sheet2", Cases: []sheet.Casename{}, DataMap: noCaseDataMap},
	}

	shList, err := sheet.Get(authFile, spreadsheetID, exceptSheetSet)

	if err != nil {
		t.Fatalf("error occurred (err=%s)", err.Error())
	}

	for i, sh := range shList {
		if sh.Name != expected[i].Name {
			t.Errorf("sheet name is not match (expected=%s, actual=%s)", expected[i].Name, sh.Name)
		}
		if !reflect.DeepEqual(sh.Cases, expected[i].Cases) {
			t.Errorf("casename list is not match (expected=%v, actual=%v)", expected[i].Cases, sh.Cases)
		}
		if !reflect.DeepEqual(sh.DataMap, expected[i].DataMap) {
			t.Errorf("data structure is not match (expected=%s, actual=%s)",
				gostr.Stringify(expected[i].DataMap), gostr.Stringify(sh.DataMap))
		}
	}
}
