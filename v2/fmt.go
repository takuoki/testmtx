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
}

type formatter struct {
	indentStr string
}

const defaultIndentStr = "  "

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
