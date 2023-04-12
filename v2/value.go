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

type stringValue struct {
	v string
}

func convertStringValue(s string) (SimpleValue, error) {
	if s == strEmpty {
		return &stringValue{v: ""}, nil
	}
	return &stringValue{v: s}, nil
}

func (v *stringValue) StringJSON() string {
	return fmt.Sprintf("%q", strings.Replace(v.v, "\n", "\\n", -1))
}

func (v *stringValue) StringYAML() string {
	return strings.Replace(v.v, "\n", "\\n", -1)
}

type numberValue struct {
	v string
}

func convertNumberValue(s string) (SimpleValue, error) {
	_, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid number value (%q)", s)
	}
	return &numberValue{v: s}, nil
}

func (v *numberValue) StringJSON() string {
	return v.v
}

func (v *numberValue) StringYAML() string {
	return v.v
}

type boolValue struct {
	v bool
}

func convertBoolValue(s string) (SimpleValue, error) {
	switch s {
	case "true", "TRUE", "True":
		return &boolValue{v: true}, nil
	case "false", "FALSE", "False":
		return &boolValue{v: false}, nil
	default:
		return nil, fmt.Errorf("invalid bool value (%q)", s)
	}
}

func (v *boolValue) StringJSON() string {
	return fmt.Sprintf("%t", v.v)
}

func (v *boolValue) StringYAML() string {
	return fmt.Sprintf("%t", v.v)
}
