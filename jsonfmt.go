package testmtx

import (
	"fmt"
	"io"
	"strings"
)

// JSONFormatter is ... TODO
type JSONFormatter struct{}

func (f *JSONFormatter) fprint(w io.Writer, v value, cn casename, i int) {
	switch val := v.(type) {
	case *vObject:
		f.fprintObject(w, val, cn, i)
	case *vArray:
		f.fprintArray(w, val, cn, i)
	case *vString:
		f.fprintString(w, val, cn)
	case *vNum:
		f.fprintNum(w, val, cn)
	case *vBool:
		f.fprintBool(w, val, cn)
	}
}

func (f *JSONFormatter) fprintObject(w io.Writer, v *vObject, cn casename, i int) {
	if !v.isNil(cn) {
		fmt.Fprintln(w, "{")
		for _, pn := range v.propertyNames {
			if !v.properties[pn].isNil(cn) {
				fmt.Fprintf(w, "%s\"%s\": ", indents(i+1), pn)
				f.fprint(w, v.properties[pn], cn, i+1)
				if !v.lastProperty(cn, pn) {
					fmt.Fprintln(w, ",")
				} else {
					fmt.Fprintln(w)
				}
			}
		}
		fmt.Fprintf(w, "%s}", indents(i))
	}
}

func (f *JSONFormatter) fprintArray(w io.Writer, v *vArray, cn casename, i int) {
	if !v.isNil(cn) {
		fmt.Fprintln(w, "[")
		for j, e := range v.elements {
			if !e.isNil(cn) {
				fmt.Fprintf(w, "%s", indents(i+1))
				f.fprint(w, e, cn, i+1)
				if !v.lastElement(cn, j) {
					fmt.Fprintln(w, ",")
				} else {
					fmt.Fprintln(w)
				}
			}
		}
		fmt.Fprintf(w, "%s]", indents(i))
	}
}

func (f *JSONFormatter) fprintString(w io.Writer, v *vString, cn casename) {
	if !v.isNil(cn) {
		fmt.Fprintf(w, "\"%s\"", f.escapeString(*v.values[cn]))
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

func (f *JSONFormatter) escapeString(s string) string {
	s = strings.Replace(s, "\n", "\\n", -1)
	return s
}

func (f *JSONFormatter) extention() string {
	return "json"
}
