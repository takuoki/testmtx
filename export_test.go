package testmtx

import (
	"io"
	"testing"
)

// Fprint prints sheet value.
// This method is a test method, and an error will occur
// if both the number of cases and the number of root elements are not one.
func (f *JSONFormatter) Fprint(t *testing.T, w io.Writer, s *Sheet, indent int) {
	t.Helper()

	// root要素一つ、ケース数一つであること前提
	if len(s.cases) != 1 {
		t.Error("the number of cases must be one")
		return
	}
	if len(s.valueMap) != 1 {
		t.Error("the number of root elements must be one")
		return
	}
	for _, v := range s.valueMap {
		f.fprint(w, v, s.cases[0], indent)
	}
}

func (f *JSONFormatter) Extension() string {
	return f.extension()
}

// Fprint prints sheet value.
// This method is a test method, and an error will occur
// if both the number of cases and the number of root elements are not one.
func (f *YamlFormatter) Fprint(t *testing.T, w io.Writer, s *Sheet, indent int) {
	t.Helper()

	// root要素一つ、ケース数一つであること前提
	if len(s.cases) != 1 {
		t.Error("the number of cases must be one")
		return
	}
	if len(s.valueMap) != 1 {
		t.Error("the number of root elements must be one")
		return
	}
	for _, v := range s.valueMap {
		f.fprint(w, v, s.cases[0], indent)
	}
}

func (f *YamlFormatter) Extension() string {
	return f.extension()
}
