package testmtx

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/sync/errgroup"
)

var createFile func(name string) (io.WriteCloser, error)

func init() {
	createFile = func(name string) (io.WriteCloser, error) {
		return os.Create(name)
	}
}

// Formatter is an interface for formatting.
// This interface has private methods, so cannot create an original formatter outside of this package.
type Formatter interface {
	fprint(w io.Writer, v value, cn casename, indent int)
	extension() string
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

// Output outputs files using Formatter.
func Output(f Formatter, s *Sheet, outdir string) error {

	if s == nil {
		return errors.New("Sheet is nil")
	}

	for k, v := range s.valueMap {
		dir := fmt.Sprintf("%s/%s", outdir, k)
		if err := os.MkdirAll(dir, 0777); err != nil {
			return fmt.Errorf("Unable to make directory: %v", err)
		}
		eg := &errgroup.Group{}
		for _, cn := range s.cases {
			cn := cn
			eg.Go(func() (er error) {
				defer func() {
					if e := recover(); e != nil {
						er = fmt.Errorf("panic recovered in goroutine: %v", e)
					}
				}()

				filepath := fmt.Sprintf("%s/%s_%s.%s", dir, s.name, cn, f.extension())

				file, err := createFile(filepath)
				if err != nil {
					return fmt.Errorf("Unable to create file: %v", err)
				}
				defer file.Close()

				f.fprint(file, v, cn, 0)
				return nil
			})
		}
		if err := eg.Wait(); err != nil {
			return err
		}
	}

	return nil
}
