package testmtx

import (
	"io"
	"strings"
)

// Formatter is an interface for formatting.
// This interface has private methods, so cannot create an original formatter outside of this package.
type Formatter interface {
	Write(w io.Writer, col Collection, cn ColumnName)
	Extension() string
	setIndentStr(s string)
}

type NewFormatterFunc func(options ...FormatOption) (Formatter, error)

// FormatOption changes some parameters of the JSONFormatter.
type FormatOption func(Formatter) error

// IndentStr changes the indent string in JSON file.
func IndentStr(s string) FormatOption {
	return func(f Formatter) error {
		f.setIndentStr(s)
		return nil
	}
}

const defaultIndentStr = "  "

type formatter struct {
	indentStr string
}

func newFormmater() *formatter {
	return &formatter{indentStr: defaultIndentStr}
}

func (f *formatter) setIndentStr(s string) {
	if f == nil {
		return
	}
	f.indentStr = s
}

func (f *formatter) indents(i int) string {
	if f == nil {
		return defaultIndentStr
	}
	return strings.Repeat(f.indentStr, i)
}
