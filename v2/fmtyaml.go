package testmtx

import (
	"fmt"
	"io"
)

// YAMLFormatter is a formatter for YAML.
type YAMLFormatter struct {
	formatter
}

// NewYAMLFormatter creates a new YAMLFormatter.
// You can change some parameters of the YAMLFormatter with YAMLFormatOption.
func NewYAMLFormatter(options ...YAMLFormatOption) (*YAMLFormatter, error) {
	f := YAMLFormatter{
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

// YAMLFormatOption changes some parameters of the YAMLFormatter.
type YAMLFormatOption func(*YAMLFormatter) error

// YAMLIndentStr changes the indent string in YAML file.
func YAMLIndentStr(s string) YAMLFormatOption {
	return func(f *YAMLFormatter) error {
		f.setIndentStr(s)
		return nil
	}
}

func (f *YAMLFormatter) Write(w io.Writer, col Collection, cn ColumnName) {
	f.write(w, col, cn, 0, false)
}

func (f *YAMLFormatter) write(w io.Writer, col Collection, cn ColumnName, indent int, isArrayChild bool) {
	switch c := col.(type) {
	case *ObjectCollection:
		f.writeObject(w, c, cn, indent, isArrayChild)
	case *ArrayCollection:
		f.writeArray(w, c, cn, indent, isArrayChild)
	case *SimpleCollection:
		f.writeSimple(w, c, cn)
	}
}

func (f *YAMLFormatter) writeObject(w io.Writer, col *ObjectCollection, cn ColumnName, indent int, isArrayChild bool) {
	if col.ImplicitNil(cn) {
		// do nothing
	} else if col.ExplicitNil(cn) {
		fmt.Fprint(w, "null") // pass only if the root is null
	} else {
		for _, pn := range col.PropertyNames {
			if col.Properties[pn].ImplicitNil(cn) {
				// do nothing
			} else {
				if !isArrayChild || !col.FirstProperty(cn, pn) {
					fmt.Fprint(w, f.indents(indent))
				}
				fmt.Fprintf(w, "%s:", pn)
				if col.Properties[pn].ExplicitNil(cn) {
					fmt.Fprintln(w, " null")
				} else {
					switch col.Properties[pn].(type) {
					case *ObjectCollection, *ArrayCollection:
						fmt.Fprintln(w)
					default:
						fmt.Fprint(w, " ")
					}
					f.write(w, col.Properties[pn], cn, indent+1, false)
				}
			}
		}
	}
}

func (f *YAMLFormatter) writeArray(w io.Writer, col *ArrayCollection, cn ColumnName, indent int, isArrayChild bool) {
	if col.ImplicitNil(cn) {
		// do nothing
	} else if col.ExplicitNil(cn) {
		fmt.Fprint(w, "null") // pass only if the root is null
	} else {
		for i, e := range col.Elements {
			if e.ImplicitNil(cn) {
				// do nothing
			} else {
				if !isArrayChild || !col.FirstElement(cn, i) {
					fmt.Fprint(w, f.indents(indent))
				}
				fmt.Fprint(w, "- ")
				if e.ExplicitNil(cn) {
					fmt.Fprintln(w, "null")
				} else {
					f.write(w, e, cn, indent+1, true)
				}
			}
		}
	}
}

func (f *YAMLFormatter) writeSimple(w io.Writer, col *SimpleCollection, cn ColumnName) {
	if col.ImplicitNil(cn) {
		// do nothing
	} else if col.ExplicitNil(cn) {
		fmt.Fprintln(w, "null") // pass only if the root is null
	} else {
		fmt.Fprintln(w, col.Values[cn].StringYAML())
	}
}

func (f *YAMLFormatter) Extension() string {
	return "yaml"
}
