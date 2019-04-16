package testmtx

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var createFile func(name string) (io.WriteCloser, error)

func init() {
	createFile = func(name string) (io.WriteCloser, error) {
		return os.Create(name)
	}
}

// Formatter is an interface for formatting.
type Formatter interface {
	fprint(w io.Writer, v value, cn casename, indent int)
	extension() string
}

// Output outputs files using Formatter.
func Output(f Formatter, ss []*Sheet, outdir string) error {

	if ss == nil {
		return errors.New("Sheet array is nil")
	}

	for _, s := range ss {
		for k, v := range s.valueMap {
			dir := fmt.Sprintf("%s/%s", outdir, k)
			if err := os.MkdirAll(dir, 0777); err != nil {
				return fmt.Errorf("Unable to make directory: %v", err)
			}
			for _, cn := range s.cases {
				filepath := fmt.Sprintf("%s/%s_%s.%s", dir, s.name, cn, f.extension())
				if err := output(f, filepath, v, cn); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func output(f Formatter, filepath string, v value, cn casename) error {
	file, err := createFile(filepath)
	if err != nil {
		return fmt.Errorf("Unable to create file: %v", err)
	}
	defer file.Close()
	f.fprint(file, v, cn, 0)
	return nil
}

func indents(i int) string {
	return strings.Repeat("  ", i)
}
