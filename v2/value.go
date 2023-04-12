package testmtx

import (
	"fmt"
	"strconv"
	"strings"
)

// TODO: コメント
type SimpleValue interface {
	StringJSON() string
	StringYAML() string
}

// TODO: コメント
type ConvertValueFunc func(s string) (SimpleValue, error)

// TODO: デフォルト値。カスタムが追加される。
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
	return fmt.Sprintf("%q", strings.Replace(v.Value, "\n", "\\n", -1))
}

func (v *StringValue) StringYAML() string {
	return strings.Replace(v.Value, "\n", "\\n", -1)
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
