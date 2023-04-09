package testmtx

import (
	"fmt"
	"io"
)

// JSONFormatter is a struct to format the sheet object as JSON.
// Create it using NewJSONFormatter function.
type JSONFormatter struct {
	formatter
}

// NewJSONFormatter creates a new JSONFormatter.
// You can change some parameters of the JSONFormatter with JSONFormatOption.
func NewJSONFormatter(options ...JSONFormatOption) (*JSONFormatter, error) {
	f := JSONFormatter{
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

// JSONFormatOption changes some parameters of the JSONFormatter.
type JSONFormatOption func(*JSONFormatter) error

// JSONIndentStr changes the indent string in JSON file.
func JSONIndentStr(s string) JSONFormatOption {
	return func(f *JSONFormatter) error {
		f.setIndentStr(s)
		return nil
	}
}

func (f *JSONFormatter) Write(w io.Writer, col Collection, cn ColumnName, indent int) {
	switch c := col.(type) {
	case *ObjectCollection:
		f.writeObject(w, c, cn, indent)
	case *ArrayCollection:
		f.writeArray(w, c, cn, indent)
	case *SimpleCollection:
		f.writeSimple(w, c, cn)
	}
}

func (f *JSONFormatter) writeObject(w io.Writer, col *ObjectCollection, cn ColumnName, indent int) {
	if col.ImplicitNil(cn) {
		// do nothing
	} else if col.ExplicitNil(cn) {
		fmt.Fprint(w, "null")
	} else {
		fmt.Fprintln(w, "{")
		for _, pn := range col.PropertyNames {
			if col.Properties[pn].ImplicitNil(cn) {
				// do nothing
			} else {
				fmt.Fprintf(w, "%s\"%s\": ", f.indents(indent+1), pn)
				if col.Properties[pn].ExplicitNil(cn) {
					fmt.Fprint(w, "null")
				} else {
					f.Write(w, col.Properties[pn], cn, indent+1)
				}
				if !col.LastProperty(cn, pn) {
					fmt.Fprintln(w, ",")
				} else {
					fmt.Fprintln(w)
				}
			}
		}
		fmt.Fprintf(w, "%s}", f.indents(indent))
	}
}

func (f *JSONFormatter) writeArray(w io.Writer, col *ArrayCollection, cn ColumnName, indent int) {
	if col.ImplicitNil(cn) {
		// do nothing
	} else if col.ExplicitNil(cn) {
		fmt.Fprint(w, "null")
	} else {
		fmt.Fprintln(w, "[")
		for i, e := range col.Elements {
			if e.ImplicitNil(cn) {
				// do nothing
			} else {
				fmt.Fprintf(w, "%s", f.indents(indent+1))
				if e.ExplicitNil(cn) {
					fmt.Fprint(w, "null")
				} else {
					f.Write(w, e, cn, indent+1)
				}
				if !col.LastElement(cn, i) {
					fmt.Fprintln(w, ",")
				} else {
					fmt.Fprintln(w)
				}
			}
		}
		fmt.Fprintf(w, "%s]", f.indents(indent))
	}
}

func (f *JSONFormatter) writeSimple(w io.Writer, col *SimpleCollection, cn ColumnName) {
	if col.ImplicitNil(cn) {
		// do nothing
	} else if col.ExplicitNil(cn) {
		fmt.Fprint(w, "null")
	} else {
		fmt.Fprint(w, col.Values[cn].StringJSON())
	}
}

func (f *JSONFormatter) Extension() string {
	return "json"
}
