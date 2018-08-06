package main

import (
	"fmt"
	"io"

	"github.com/takuoki/testmtx/sheet"
)

type yamlf struct{}

func (f *yamlf) OutData(out io.Writer, d sheet.Data, c sheet.Casename, i int) error {
	return f.OutData2(out, d, c, i, false)
}

func (f *yamlf) OutData2(out io.Writer, d sheet.Data, c sheet.Casename, i int, br bool) error {
	switch t := d.(type) {
	case *sheet.DString:
		f.outString(out, t, c)
	case *sheet.DNum:
		f.outNum(out, t, c)
	case *sheet.DBool:
		f.outBool(out, t, c)
	case *sheet.DObject:
		f.outObject(out, t, c, i, br)
	case *sheet.DArray:
		f.outArray(out, t, c, i, br)
	default:
		return fmt.Errorf("no such type (%v)", t)
	}
	return nil
}

func (f *yamlf) outObject(out io.Writer, d *sheet.DObject, c sheet.Casename, i int, br bool) {
	if !d.IsNil(c) {
		if br {
			out.Write([]byte("\n"))
		}
		for _, pn := range d.PropertyNames {
			if !d.Properties[pn].IsNil(c) {
				idt := ""
				if br || !d.FirstProperty(c, pn) {
					idt = f.indents(i + 1)
				}
				out.Write([]byte(fmt.Sprintf("%s%s: ", idt, pn)))
				f.OutData2(out, d.Properties[pn], c, i+1, true)
				if !d.LastProperty(c, pn) {
					out.Write([]byte("\n"))
				}
			}
		}
	}
}

func (f *yamlf) outArray(out io.Writer, d *sheet.DArray, c sheet.Casename, i int, br bool) {
	if !d.IsNil(c) {
		if br {
			out.Write([]byte("\n"))
		}
		for j, e := range d.Elements {
			if !e.IsNil(c) {
				idt := ""
				if br || !d.FirstElement(c, j) {
					idt = f.indents(i + 1)
				}
				out.Write([]byte(fmt.Sprintf("%s- ", idt)))
				f.OutData2(out, e, c, i+1, false)
				if !d.LastElement(c, j) {
					out.Write([]byte("\n"))
				}
			}
		}
	}
}

func (f *yamlf) outString(out io.Writer, d *sheet.DString, c sheet.Casename) {
	if !d.IsNil(c) {
		out.Write([]byte(fmt.Sprintf("%s", *d.Values[c])))
	}
}

func (f *yamlf) outNum(out io.Writer, d *sheet.DNum, c sheet.Casename) {
	if !d.IsNil(c) {
		out.Write([]byte(fmt.Sprintf("%s", *d.Values[c])))
	}
}

func (f *yamlf) outBool(out io.Writer, d *sheet.DBool, c sheet.Casename) {
	if !d.IsNil(c) {
		out.Write([]byte(fmt.Sprintf("%t", *d.Values[c])))
	}
}

func (f *yamlf) indents(i int) string {
	indent := "  "
	str := ""
	for j := 0; j < i-1; j++ {
		str += indent
	}
	return str
}

func (f *yamlf) Extention() string {
	return "yaml"
}
