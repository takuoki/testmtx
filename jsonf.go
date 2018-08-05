package main

import (
	"fmt"
	"io"

	"github.com/takuoki/testmtx/sheet"
)

type jsonf struct{}

func (f *jsonf) OutData(out io.Writer, d sheet.Data, c sheet.Casename, i int) error {
	switch t := d.(type) {
	case *sheet.DString:
		f.outString(out, t, c)
	case *sheet.DNum:
		f.outNum(out, t, c)
	case *sheet.DBool:
		f.outBool(out, t, c)
	case *sheet.DObject:
		f.outObject(out, t, c, i)
	case *sheet.DArray:
		f.outArray(out, t, c, i)
	default:
		return fmt.Errorf("no such type (%v)", t)
	}
	return nil
}

func (f *jsonf) outObject(out io.Writer, d *sheet.DObject, c sheet.Casename, i int) {
	if !d.IsNil(c) {
		out.Write([]byte("{\n"))
		for _, pn := range d.PropertyNames {
			if !d.Properties[pn].IsNil(c) {
				out.Write([]byte(fmt.Sprintf("%s\"%s\": ", f.indents(i+1), pn)))
				f.OutData(out, d.Properties[pn], c, i+1)
				if !d.LastProperty(c, pn) {
					out.Write([]byte(","))
				}
				out.Write([]byte("\n"))
			}
		}
		out.Write([]byte(fmt.Sprintf("%s}", f.indents(i))))
	}
}

func (f *jsonf) outArray(out io.Writer, d *sheet.DArray, c sheet.Casename, i int) {
	if !d.IsNil(c) {
		out.Write([]byte("[\n"))
		for j, e := range d.Elements {
			if !e.IsNil(c) {
				out.Write([]byte(fmt.Sprintf("%s", f.indents(i+1))))
				f.OutData(out, e, c, i+1)
				if !d.LastElement(c, j) {
					out.Write([]byte(","))
				}
				out.Write([]byte("\n"))
			}
		}
		out.Write([]byte(fmt.Sprintf("%s]", f.indents(i))))
	}
}

func (f *jsonf) outString(out io.Writer, d *sheet.DString, c sheet.Casename) {
	if !d.IsNil(c) {
		out.Write([]byte(fmt.Sprintf("\"%s\"", *d.Values[c])))
	}
}

func (f *jsonf) outNum(out io.Writer, d *sheet.DNum, c sheet.Casename) {
	if !d.IsNil(c) {
		out.Write([]byte(fmt.Sprintf("%s", *d.Values[c])))
	}
}

func (f *jsonf) outBool(out io.Writer, d *sheet.DBool, c sheet.Casename) {
	if !d.IsNil(c) {
		out.Write([]byte(fmt.Sprintf("%t", *d.Values[c])))
	}
}

func (f *jsonf) indents(i int) string {
	indent := "  "
	str := ""
	for j := 0; j < i; j++ {
		str += indent
	}
	return str
}

func (f *jsonf) Extention() string {
	return "json"
}
