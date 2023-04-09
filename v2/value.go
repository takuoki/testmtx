package testmtx

import (
	"fmt"
	"strconv"
	"strings"
)

// TODO: ラッパーである旨をコメント
// → collectionにマージするでもいいかも
// type simpleValue struct {
// 	implicitNil bool
// 	explicitNil bool
// 	value       SimpleValue
// }

// func (c *Client) convertSimpleValue(valueType, s string) (*simpleValue, error) {
// 	if s == "" {
// 		return &simpleValue{implicitNil: true}, nil
// 	}
// 	if s == strNull {
// 		return &simpleValue{explicitNil: true}, nil
// 	}
// 	v, err := c.convertSimpleValueFuncs[valueType](s)
// 	if err != nil {
// 		return nil, fmt.Errorf("fail to convert simple value (type=%q, value=%q): %w", valueType, s, err)
// 	}
// 	return &simpleValue{value: v}, nil
// }

// func (v *simpleValue) stringJSON() string {
// 	if v.implicitNil {
// 		return ""
// 	}
// 	if v.explicitNil {
// 		return "null"
// 	}
// 	return v.value.StringJSON()
// }

// func (v *simpleValue) stringYAML() string {
// 	if v.implicitNil {
// 		return ""
// 	}
// 	if v.explicitNil {
// 		return "null"
// 	}
// 	return v.value.StringYAML()
// }

// TODO: コメント
type SimpleValue interface {
	StringJSON() string
	StringYAML() string
}

// TODO: コメント
type ConvertValueFunc func(s string) (SimpleValue, error)

// TODO: デフォルト値。カスタムが追加される。
// Clientの生成時に使用される
var defaultConvertSimpleValueFuncs = map[string]func(s string) (SimpleValue, error){
	typeString: convertValueString,
	typeNumber: convertValueNumber,
	typeBool:   convertValueBool,
}

type valueString struct {
	v string
}

func convertValueString(s string) (SimpleValue, error) {
	if s == strEmpty {
		return &valueString{v: ""}, nil
	}
	return &valueString{v: s}, nil
}

func (v *valueString) StringJSON() string {
	return fmt.Sprintf("%q", strings.Replace(v.v, "\n", "\\n", -1))
}

func (v *valueString) StringYAML() string {
	return strings.Replace(v.v, "\n", "\\n", -1)
}

type valueNumber struct {
	v string
}

func convertValueNumber(s string) (SimpleValue, error) {
	_, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, err
	}
	return &valueNumber{v: s}, nil
}

func (v *valueNumber) StringJSON() string {
	return v.v
}

func (v *valueNumber) StringYAML() string {
	return v.v
}

type valueBool struct {
	v bool
}

func convertValueBool(s string) (SimpleValue, error) {
	switch s {
	case "true", "TRUE", "True":
		return &valueBool{v: true}, nil
	case "false", "FALSE", "False":
		return &valueBool{v: false}, nil
	default:
		return nil, fmt.Errorf("invalid bool value (%q)", s)
	}
}

func (v *valueBool) StringJSON() string {
	return fmt.Sprintf("%t", v.v)
}

func (v *valueBool) StringYAML() string {
	return fmt.Sprintf("%t", v.v)
}
