package testmtx

import (
	"fmt"
	"strconv"
)

// SimpleValue is a simple value.
// User can define custom SimpleValue and add it to the parser
// by using AdditionalSimpleValues option when generating the parser.
// See `unixtimeValue` in `sample_test.go` for an example.
type SimpleValue interface {
	StringJSON() string
	StringYAML() string
}

// ConvertValueFunc is a function to convert a string to a SimpleValue.
// See `convertUnixtimeValue` in `sample_test.go` for an example.
type ConvertValueFunc func(s string) (SimpleValue, error)

// User can add custom simple values.
// But cannot overwrite or delete default simple values.
var defaultConvertSimpleValueFuncs = map[string]ConvertValueFunc{
	typeString: convertStringValue,
	typeNumber: convertNumberValue,
	typeBool:   convertBoolValue,
}

type StringValue struct {
	Value string
}

func convertStringValue(s string) (SimpleValue, error) {
	if s == strEmpty {
		return &StringValue{Value: ""}, nil
	}
	return &StringValue{Value: s}, nil
}

func (v *StringValue) StringJSON() string {
	return fmt.Sprintf("%q", v.Value)
}

func (v *StringValue) StringYAML() string {
	return fmt.Sprintf("%q", v.Value)
}

type NumberValue struct {
	Value string
}

func convertNumberValue(s string) (SimpleValue, error) {
	_, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid number value (%q)", s)
	}
	return &NumberValue{Value: s}, nil
}

func (v *NumberValue) StringJSON() string {
	return v.Value
}

func (v *NumberValue) StringYAML() string {
	return v.Value
}

type BoolValue struct {
	Value bool
}

func convertBoolValue(s string) (SimpleValue, error) {
	switch s {
	case "true", "TRUE", "True":
		return &BoolValue{Value: true}, nil
	case "false", "FALSE", "False":
		return &BoolValue{Value: false}, nil
	default:
		return nil, fmt.Errorf("invalid bool value (%q)", s)
	}
}

func (v *BoolValue) StringJSON() string {
	return fmt.Sprintf("%t", v.Value)
}

func (v *BoolValue) StringYAML() string {
	return fmt.Sprintf("%t", v.Value)
}
