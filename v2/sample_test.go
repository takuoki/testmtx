package testmtx_test

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/takuoki/testmtx/v2"
)

// propLevel = 5
func sampleDocSheet() *mockDocSheet {
	raw := `
|     |            |            |      |      |     |        | Data    |                      |           |         |
|     | Properties |            |      |      |     | Type   | case1   | case2                | case3     | case4   |
|     | in         |            |      |      |     | object | *new    | *new                 | *new      | *new    |
|     |            | num_key    |      |      |     | number | 101     | 102                  | 103       | *null   |
|     |            | string_key |      |      |     | string | value 1 | value 2<br>multiline |           | *null   |
|     |            | bool_key   |      |      |     | bool   | true    |                      | false     | *null   |
|     |            | object_key |      |      |     | object | *new    | *new                 |           | *null   |
|     |            |            | key1 |      |     | number | 201     | 202                  |           |         |
|     |            |            | key2 |      |     | string | value 2 | *empty               |           |         |
|     |            | array_key  |      |      |     | array  | *new    |                      | *new      | *null   |
|     |            |            | * 0  |      |     | object | *new    |                      | *new      |         |
|     |            |            |      | key3 |     | number | 301     |                      | 303       |         |
|     |            |            |      | key4 |     | string | value 3 |                      | value "3" |         |
|     |            |            | * 1  |      |     | object | *new    |                      |           |         |
|     |            |            |      | key3 |     | number | 401     |                      |           |         |
|     |            |            |      | key4 |     | string | value 4 |                      |           |         |


|     | want       |            |      |      |     | object | *new    | *new                 | *new      | *new    |
|     |            | status     |      |      |     | string | success | failure              | failure   | failure |
|     |            | code       |      |      |     | number | 200     | 401                  | 404       | 500     |
`

	res := [][]string{}
	for _, line := range strings.Split(raw, "\n") {
		values := []string{}
		for i, v := range strings.Split(line, "|") {
			if i == 0 || i == len(strings.Split(line, "|"))-1 {
				continue
			}
			values = append(values, strings.Replace(strings.TrimSpace(v), "<br>", "\n", -1))
		}
		res = append(res, values)
	}

	fmt.Println(len(res))

	return newMockSheet("sample", res)
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
							"case1": &testmtx.StringValue{Value: "value 1"},
							"case2": &testmtx.StringValue{Value: "value 2\nmultiline"},
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
									"case1": &testmtx.StringValue{Value: "value 2"},
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
											"case1": &testmtx.StringValue{Value: "value 3"},
											"case3": &testmtx.StringValue{Value: "value \"3\""},
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
											"case1": &testmtx.StringValue{Value: "value 4"},
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
