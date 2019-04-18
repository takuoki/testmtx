package testmtx

import (
	"fmt"
	"io"
	"strings"
)

// YamlFormatter is a struct to format the sheet object as YAML.
// Create it using NewYamlFormatter function.
type YamlFormatter struct {
	formatter
}

// NewYamlFormatter creates a new YamlFormatter.
// You can change some parameters of the YamlFormatter with YamlFormatOption.
func NewYamlFormatter(options ...YamlFormatOption) (*YamlFormatter, error) {
	f := YamlFormatter{
		formatter{indentStr: defaultIndentStr},
	}
	for _, opt := range options {
		err := opt(&f)
		if err != nil {
			return nil, err
		}
	}
	return &f, nil
}

// YamlFormatOption changes some parameters of the YamlFormatter.
type YamlFormatOption func(*YamlFormatter) error

// YamlIndentStr changes the indent string in Yaml file.
func YamlIndentStr(s string) YamlFormatOption {
	return func(f *YamlFormatter) error {
		f.setIndentStr(s)
		return nil
	}
}

func (f *YamlFormatter) fprint(w io.Writer, v value, cn casename, indent int) {
	f.fprintWithBr(w, v, cn, indent, false)
}

func (f *YamlFormatter) fprintWithBr(w io.Writer, v value, cn casename, indent int, br bool) {
	switch t := v.(type) {
	case *vObject:
		f.fprintObject(w, t, cn, indent, br)
	case *vArray:
		f.fprintArray(w, t, cn, indent, br)
	case *vString:
		f.fprintString(w, t, cn)
	case *vNum:
		f.fprintNum(w, t, cn)
	case *vBool:
		f.fprintBool(w, t, cn)
	}
}

func (f *YamlFormatter) fprintObject(w io.Writer, v *vObject, cn casename, i int, br bool) {
	if !v.isNil(cn) {
		if br {
			fmt.Fprintln(w)
		}
		for _, pn := range v.propertyNames {
			if !v.properties[pn].isNil(cn) {
				var idt string
				if br || !v.firstProperty(cn, pn) {
					idt = f.indents(i + 1)
				}
				fmt.Fprintf(w, "%s%s: ", idt, pn)
				f.fprintWithBr(w, v.properties[pn], cn, i+1, true)
				if !v.lastProperty(cn, pn) {
					fmt.Fprintln(w)
				}
			}
		}
	}
}

func (f *YamlFormatter) fprintArray(w io.Writer, v *vArray, cn casename, i int, br bool) {
	if !v.isNil(cn) {
		if br {
			fmt.Fprintln(w)
		}
		for j, e := range v.elements {
			if !e.isNil(cn) {
				idt := ""
				if br || !v.firstElement(cn, j) {
					idt = f.indents(i + 1)
				}
				fmt.Fprintf(w, "%s- ", idt)
				f.fprintWithBr(w, e, cn, i+1, false)
				if !v.lastElement(cn, j) {
					fmt.Fprintln(w)
				}
			}
		}
	}
}

func (f *YamlFormatter) fprintString(w io.Writer, v *vString, cn casename) {
	if !v.isNil(cn) {
		s := strings.Replace(*v.values[cn], "\n", "\\n", -1)
		fmt.Fprint(w, s)
	}
}

func (f *YamlFormatter) fprintNum(w io.Writer, v *vNum, cn casename) {
	if !v.isNil(cn) {
		fmt.Fprintf(w, "%s", *v.values[cn])
	}
}

func (f *YamlFormatter) fprintBool(w io.Writer, v *vBool, cn casename) {
	if !v.isNil(cn) {
		fmt.Fprintf(w, "%t", *v.values[cn])
	}
}

func (f *YamlFormatter) extension() string {
	return "yaml"
}
