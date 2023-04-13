package testmtx_test

import (
	"fmt"
	"strconv"
	"time"

	"github.com/takuoki/testmtx/v2"
)

// propLevel = 5
func sampleDocSheet() *mockDocSheet {
	return newMockSheet("sample", [][]string{
		{},
		{"", "", "", "", "", "", "", "Data"},
		{"", "Properties", "", "", "", "", "Type", "case1", "case2", "case3", "case4"},
		{"", "in", "", "", "", "", "object", "*new", "*new", "*new", "*new"},
		{"", "", "num_key", "", "", "", "number", "101", "102", "103", "*null"},
		{"", "", "string_key", "", "", "", "string", "string value 101", "string value 201", "", "*null"},
		{"", "", "bool_key", "", "", "", "bool", "true", "", "false", "*null"},
		{"", "", "object_key", "", "", "", "object", "*new", "*new", "", "*null"},
		{"", "", "", "key1", "", "", "number", "201", "202", "", ""},
		{"", "", "", "key2", "", "", "string", "string value 201", "*empty", "", ""},
		{"", "", "array_key", "", "", "", "array", "*new", "", "*new", "*null"},
		{"", "", "", "* 0", "", "", "object", "*new", "", "*new", ""},
		{"", "", "", "", "key3", "", "number", "301", "", "303", ""},
		{"", "", "", "", "key4", "", "string", "string value 301", "", "string value 303", ""},
		{"", "", "", "* 1", "", "", "object", "*new", "", "", ""},
		{"", "", "", "", "key3", "", "number", "401", "", "", ""},
		{"", "", "", "", "key4", "", "string", "string value 401", "", "", ""},
		{},
		{},
		{"", "want", "", "", "", "", "object", "*new", "*new", "*new", "*new"},
		{"", "", "status", "", "", "", "string", "success", "failure", "failure", "failure"},
		{"", "", "code", "", "", "", "number", "200", "401", "404", "500"},
	})
}

func sampleParsedSheet() *testmtx.Sheet {
	return &testmtx.Sheet{
		Name:        "sample",
		ColumnNames: []testmtx.ColumnName{"case1", "case2", "case3", "case4"},
		Collections: map[testmtx.PropName]testmtx.Collection{
			"in": &testmtx.ObjectCollection{
				ImplicitNils:  map[testmtx.ColumnName]bool{},
				ExplicitNils:  map[testmtx.ColumnName]bool{},
				PropertyNames: []testmtx.PropName{"num_key", "string_key", "bool_key", "object_key", "array_key"},
				Properties: map[testmtx.PropName]testmtx.Collection{
					"num_key": &testmtx.SimpleCollection{
						ImplicitNils: map[testmtx.ColumnName]bool{},
						ExplicitNils: map[testmtx.ColumnName]bool{
							"case4": true,
						},
						Values: map[testmtx.ColumnName]testmtx.SimpleValue{
							"case1": &testmtx.NumberValue{Value: "101"},
							"case2": &testmtx.NumberValue{Value: "102"},
							"case3": &testmtx.NumberValue{Value: "103"},
						},
					},
					"string_key": &testmtx.SimpleCollection{
						ImplicitNils: map[testmtx.ColumnName]bool{
							"case3": true,
						},
						ExplicitNils: map[testmtx.ColumnName]bool{
							"case4": true,
						},
						Values: map[testmtx.ColumnName]testmtx.SimpleValue{
							"case1": &testmtx.StringValue{Value: "string value 101"},
							"case2": &testmtx.StringValue{Value: "string value 201"},
						},
					},
					"bool_key": &testmtx.SimpleCollection{
						ImplicitNils: map[testmtx.ColumnName]bool{
							"case2": true,
						},
						ExplicitNils: map[testmtx.ColumnName]bool{
							"case4": true,
						},
						Values: map[testmtx.ColumnName]testmtx.SimpleValue{
							"case1": &testmtx.BoolValue{Value: true},
							"case3": &testmtx.BoolValue{Value: false},
						},
					},
					"object_key": &testmtx.ObjectCollection{
						ImplicitNils: map[testmtx.ColumnName]bool{
							"case3": true,
						},
						ExplicitNils: map[testmtx.ColumnName]bool{
							"case4": true,
						},
						PropertyNames: []testmtx.PropName{"key1", "key2"},
						Properties: map[testmtx.PropName]testmtx.Collection{
							"key1": &testmtx.SimpleCollection{
								ImplicitNils: map[testmtx.ColumnName]bool{
									"case3": true,
									"case4": true,
								},
								ExplicitNils: map[testmtx.ColumnName]bool{},
								Values: map[testmtx.ColumnName]testmtx.SimpleValue{
									"case1": &testmtx.NumberValue{Value: "201"},
									"case2": &testmtx.NumberValue{Value: "202"},
								},
							},
							"key2": &testmtx.SimpleCollection{
								ImplicitNils: map[testmtx.ColumnName]bool{
									"case3": true,
									"case4": true,
								},
								ExplicitNils: map[testmtx.ColumnName]bool{},
								Values: map[testmtx.ColumnName]testmtx.SimpleValue{
									"case1": &testmtx.StringValue{Value: "string value 201"},
									"case2": &testmtx.StringValue{Value: ""},
								},
							},
						},
					},
					"array_key": &testmtx.ArrayCollection{
						ImplicitNils: map[testmtx.ColumnName]bool{
							"case2": true,
						},
						ExplicitNils: map[testmtx.ColumnName]bool{
							"case4": true,
						},
						Elements: []testmtx.Collection{
							&testmtx.ObjectCollection{
								ImplicitNils: map[testmtx.ColumnName]bool{
									"case2": true,
									"case4": true,
								},
								ExplicitNils:  map[testmtx.ColumnName]bool{},
								PropertyNames: []testmtx.PropName{"key3", "key4"},
								Properties: map[testmtx.PropName]testmtx.Collection{
									"key3": &testmtx.SimpleCollection{
										ImplicitNils: map[testmtx.ColumnName]bool{
											"case2": true,
											"case4": true,
										},
										ExplicitNils: map[testmtx.ColumnName]bool{},
										Values: map[testmtx.ColumnName]testmtx.SimpleValue{
											"case1": &testmtx.NumberValue{Value: "301"},
											"case3": &testmtx.NumberValue{Value: "303"},
										},
									},
									"key4": &testmtx.SimpleCollection{
										ImplicitNils: map[testmtx.ColumnName]bool{
											"case2": true,
											"case4": true,
										},
										ExplicitNils: map[testmtx.ColumnName]bool{},
										Values: map[testmtx.ColumnName]testmtx.SimpleValue{
											"case1": &testmtx.StringValue{Value: "string value 301"},
											"case3": &testmtx.StringValue{Value: "string value 303"},
										},
									},
								},
							},
							&testmtx.ObjectCollection{
								ImplicitNils: map[testmtx.ColumnName]bool{
									"case2": true,
									"case3": true,
									"case4": true,
								},
								ExplicitNils:  map[testmtx.ColumnName]bool{},
								PropertyNames: []testmtx.PropName{"key3", "key4"},
								Properties: map[testmtx.PropName]testmtx.Collection{
									"key3": &testmtx.SimpleCollection{
										ImplicitNils: map[testmtx.ColumnName]bool{
											"case2": true,
											"case3": true,
											"case4": true,
										},
										ExplicitNils: map[testmtx.ColumnName]bool{},
										Values: map[testmtx.ColumnName]testmtx.SimpleValue{
											"case1": &testmtx.NumberValue{Value: "401"},
										},
									},
									"key4": &testmtx.SimpleCollection{
										ImplicitNils: map[testmtx.ColumnName]bool{
											"case2": true,
											"case3": true,
											"case4": true,
										},
										ExplicitNils: map[testmtx.ColumnName]bool{},
										Values: map[testmtx.ColumnName]testmtx.SimpleValue{
											"case1": &testmtx.StringValue{Value: "string value 401"},
										},
									},
								},
							},
						},
					},
				},
			},
			"want": &testmtx.ObjectCollection{
				ImplicitNils:  map[testmtx.ColumnName]bool{},
				ExplicitNils:  map[testmtx.ColumnName]bool{},
				PropertyNames: []testmtx.PropName{"status", "code"},
				Properties: map[testmtx.PropName]testmtx.Collection{
					"status": &testmtx.SimpleCollection{
						ImplicitNils: map[testmtx.ColumnName]bool{},
						ExplicitNils: map[testmtx.ColumnName]bool{},
						Values: map[testmtx.ColumnName]testmtx.SimpleValue{
							"case1": &testmtx.StringValue{Value: "success"},
							"case2": &testmtx.StringValue{Value: "failure"},
							"case3": &testmtx.StringValue{Value: "failure"},
							"case4": &testmtx.StringValue{Value: "failure"},
						},
					},
					"code": &testmtx.SimpleCollection{
						ImplicitNils: map[testmtx.ColumnName]bool{},
						ExplicitNils: map[testmtx.ColumnName]bool{},
						Values: map[testmtx.ColumnName]testmtx.SimpleValue{
							"case1": &testmtx.NumberValue{Value: "200"},
							"case2": &testmtx.NumberValue{Value: "401"},
							"case3": &testmtx.NumberValue{Value: "404"},
							"case4": &testmtx.NumberValue{Value: "500"},
						},
					},
				},
			},
		},
	}
}

type unixtimeValue struct {
	Value string
}

func convertUnixtimeValue(s string) (testmtx.SimpleValue, error) {
	const timeFormat = "2006-01-02 15:04:05"

	t, err := time.Parse(timeFormat, s)
	if err != nil {
		return nil, fmt.Errorf("invalid unixtime value (%q)", s)
	}
	ut := strconv.FormatInt(t.Unix(), 10)
	return &unixtimeValue{Value: ut}, nil
}

func (v *unixtimeValue) StringJSON() string {
	return fmt.Sprintf("%q", v.Value)
}

func (v *unixtimeValue) StringYAML() string {
	return v.Value
}
