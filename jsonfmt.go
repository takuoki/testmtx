package testmtx

import (
	"fmt"
	"io"
	"strings"
)

// JSONFormatter is ... TODO
type JSONFormatter struct {
	formatter
}

// NewJSONFormatter creates a new JSONFormatter.
// You can change some parameters of the JSONFormatter with JSONFormatOption.
func NewJSONFormatter(options ...JSONFormatOption) (*JSONFormatter, error) {
	f := JSONFormatter{
		formatter{indentStr: "  "},
	}
	for _, opt := range options {
		err := opt(&f)
		if err != nil {
			return nil, err
		}
	}
	return &f, nil
}

// JSONFormatOption changes some parameters of the JSONFormatter.
type JSONFormatOption func(*JSONFormatter) error

// JSONIndentStr changes the indent string in JSON file.
func JSONIndentStr(s string) JSONFormatOption {
	return func(f *JSONFormatter) error {
		f.setIndentStr(s)
		return nil
	}
}

func (f *JSONFormatter) fprint(w io.Writer, v value, cn casename, indent int) {
	switch val := v.(type) {
	case *vObject:
		f.fprintObject(w, val, cn, indent)
	case *vArray:
		f.fprintArray(w, val, cn, indent)
	case *vString:
		f.fprintString(w, val, cn)
	case *vNum:
		f.fprintNum(w, val, cn)
	case *vBool:
		f.fprintBool(w, val, cn)
	}
}

func (f *JSONFormatter) fprintObject(w io.Writer, v *vObject, cn casename, indent int) {
	if !v.isNil(cn) {
		fmt.Fprintln(w, "{")
		for _, pn := range v.propertyNames {
			if !v.properties[pn].isNil(cn) {
				fmt.Fprintf(w, "%s\"%s\": ", f.indents(indent+1), pn)
				f.fprint(w, v.properties[pn], cn, indent+1)
				if !v.lastProperty(cn, pn) {
					fmt.Fprintln(w, ",")
				} else {
					fmt.Fprintln(w)
				}
			}
		}
		fmt.Fprintf(w, "%s}", f.indents(indent))
	}
}

func (f *JSONFormatter) fprintArray(w io.Writer, v *vArray, cn casename, indent int) {
	if !v.isNil(cn) {
		fmt.Fprintln(w, "[")
		for j, e := range v.elements {
			if !e.isNil(cn) {
				fmt.Fprintf(w, "%s", f.indents(indent+1))
				f.fprint(w, e, cn, indent+1)
				if !v.lastElement(cn, j) {
					fmt.Fprintln(w, ",")
				} else {
					fmt.Fprintln(w)
				}
			}
		}
		fmt.Fprintf(w, "%s]", f.indents(indent))
	}
}

func (f *JSONFormatter) fprintString(w io.Writer, v *vString, cn casename) {
	if !v.isNil(cn) {
		s := strings.Replace(*v.values[cn], "\n", "\\n", -1)
		fmt.Fprintf(w, "\"%s\"", s)
	}
}

func (f *JSONFormatter) fprintNum(w io.Writer, v *vNum, cn casename) {
	if !v.isNil(cn) {
		fmt.Fprintf(w, "%s", *v.values[cn])
	}
}

func (f *JSONFormatter) fprintBool(w io.Writer, v *vBool, cn casename) {
	if !v.isNil(cn) {
		fmt.Fprintf(w, "%t", *v.values[cn])
	}
}

func (f *JSONFormatter) extension() string {
	return "json"
}
