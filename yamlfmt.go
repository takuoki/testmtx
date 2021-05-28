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
	switch t := v.(type) {
	case *vObject:
		f.fprintObject(w, t, cn, indent)
	case *vArray:
		f.fprintArray(w, t, cn, indent)
	case *vString:
		f.fprintString(w, t, cn)
	case *vNum:
		f.fprintNum(w, t, cn)
	case *vBool:
		f.fprintBool(w, t, cn)
	}
}

func (f *YamlFormatter) fprintObject(w io.Writer, v *vObject, cn casename, i int) {
	if !v.isNil(cn) {
		for _, pn := range v.propertyNames {
			if !v.properties[pn].isNil(cn) {
				idt := f.indents(i)
				switch v.properties[pn].(type) {
				case *vObject, *vArray:
					fmt.Fprintf(w, "%s%s:\n", idt, pn)
				case *vString, *vNum, *vBool:
					fmt.Fprintf(w, "%s%s: ", idt, pn)
				}
				f.fprint(w, v.properties[pn], cn, i+1)
			}
		}
	}
}

func (f *YamlFormatter) fprintArray(w io.Writer, v *vArray, cn casename, i int) {
	if !v.isNil(cn) {
		for _, e := range v.elements {
			if !e.isNil(cn) {
				idt := f.indents(i)
				switch e.(type) {
				case *vObject, *vArray:
					fmt.Fprintf(w, "%s-\n", idt)
				case *vString, *vNum, *vBool:
					fmt.Fprintf(w, "%s- ", idt)
				}
				f.fprint(w, e, cn, i+1)
			}
		}
	}
}

func (f *YamlFormatter) fprintString(w io.Writer, v *vString, cn casename) {
	if !v.isNil(cn) {
		s := strings.Replace(*v.values[cn], "\n", "\\n", -1)
		fmt.Fprintln(w, s)
	}
}

func (f *YamlFormatter) fprintNum(w io.Writer, v *vNum, cn casename) {
	if !v.isNil(cn) {
		fmt.Fprintf(w, "%s\n", *v.values[cn])
	}
}

func (f *YamlFormatter) fprintBool(w io.Writer, v *vBool, cn casename) {
	if !v.isNil(cn) {
		fmt.Fprintf(w, "%t\n", *v.values[cn])
	}
}

func (f *YamlFormatter) extension() string {
	return "yaml"
}
